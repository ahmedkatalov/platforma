package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"platforma/backend/internal/auth"
	"platforma/backend/internal/config"
	"platforma/backend/internal/db"
	"platforma/backend/internal/dotenv"
	"platforma/backend/internal/handler"
	"platforma/backend/internal/mailer"
	"platforma/backend/internal/repository"
	"platforma/backend/internal/router"
	"platforma/backend/internal/service"
)

func main() {
	dotenv.Load(".env")

	cfg := config.Load()
	if err := cfg.Validate(); err != nil {
		log.Fatalf("конфигурация: %v", err)
	}

	ctx := context.Background()
	pool, err := db.New(ctx, cfg.DatabaseURL)
	if err != nil {
		log.Fatalf("подключение к базе: %v", err)
	}
	defer pool.Close()

	if err := db.RunMigrations(ctx, pool); err != nil {
		log.Fatalf("миграции: %v", err)
	}

	// Репозитории.
	userRepo := repository.NewUserRepo(pool)
	tokenRepo := repository.NewTokenRepo(pool)
	codeRepo := repository.NewCodeRepo(pool)
	courseRepo := repository.NewCourseRepo(pool)
	activityRepo := repository.NewActivityRepo(pool)
	auditRepo := repository.NewAuditRepo(pool)
	statsRepo := repository.NewStatsRepo(pool)
	themeRepo := repository.NewThemeRepo(pool)

	// Сервисы.
	tokens := auth.NewTokenManager(cfg.JWTSecret, cfg.AccessTTL, cfg.RefreshTTL)
	mail := mailer.New(cfg)
	authSvc := service.NewAuthService(cfg, userRepo, tokenRepo, codeRepo, activityRepo, tokens, mail)
	userSvc := service.NewUserService(userRepo, tokenRepo, mail, auditRepo)

	if !mail.Enabled() {
		log.Println("mailer: EmailJS не настроен — коды подтверждения выводятся в лог")
	}

	// Хендлеры.
	authHandler := handler.NewAuthHandler(authSvc)
	meHandler := handler.NewMeHandler(userRepo, courseRepo, activityRepo, statsRepo, themeRepo, authHandler)
	adminHandler := handler.NewAdminHandler(userRepo, courseRepo, statsRepo, activityRepo, auditRepo, themeRepo, userSvc)
	courseHandler := handler.NewCourseHandler(courseRepo, auditRepo)
	themeHandler := handler.NewThemeHandler(themeRepo)

	go cleanupLoop(ctx, tokenRepo, codeRepo)

	srv := &http.Server{
		Addr: ":" + cfg.Port,
		Handler: router.New(router.Deps{
			Config:  cfg,
			Tokens:  tokens,
			Auth:    authHandler,
			Me:      meHandler,
			Admin:   adminHandler,
			Courses: courseHandler,
			Theme:   themeHandler,
		}),
		ReadHeaderTimeout: 10 * time.Second,
		ReadTimeout:       15 * time.Second,
		WriteTimeout:      30 * time.Second,
		IdleTimeout:       60 * time.Second,
	}

	go func() {
		log.Printf("api: слушаю на http://localhost:%s (env=%s)", cfg.Port, cfg.AppEnv)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("http сервер: %v", err)
		}
	}()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)
	<-stop

	log.Println("api: останавливаюсь...")
	shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := srv.Shutdown(shutdownCtx); err != nil {
		log.Printf("api: некорректная остановка: %v", err)
	}
}

// cleanupLoop раз в час удаляет протухшие токены и коды.
func cleanupLoop(ctx context.Context, tokens *repository.TokenRepo, codes *repository.CodeRepo) {
	ticker := time.NewTicker(time.Hour)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			if n, err := tokens.DeleteExpired(ctx); err == nil && n > 0 {
				log.Printf("cleanup: удалено refresh-токенов: %d", n)
			}
			if n, err := codes.DeleteExpired(ctx); err == nil && n > 0 {
				log.Printf("cleanup: удалено кодов подтверждения: %d", n)
			}
		}
	}
}
