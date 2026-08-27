// Команда createadmin создаёт (или обновляет) администратора платформы.
//
//	go run ./cmd/createadmin
//	go run ./cmd/createadmin -email admin@example.com -name "Ахмед" -password "Secret123"
package main

import (
	"bufio"
	"context"
	"errors"
	"flag"
	"fmt"
	"log"
	"os"
	"strings"
	"syscall"
	"time"

	"platforma/backend/internal/auth"
	"platforma/backend/internal/config"
	"platforma/backend/internal/db"
	"platforma/backend/internal/domain"
	"platforma/backend/internal/dotenv"
	"platforma/backend/internal/repository"

	"golang.org/x/term"
)

func main() {
	email := flag.String("email", "", "почта администратора")
	name := flag.String("name", "", "имя администратора")
	password := flag.String("password", "", "пароль (если не задан — спросим в терминале)")
	force := flag.Bool("force", false, "если аккаунт существует — обновить его до администратора и сменить пароль")
	flag.Parse()

	dotenv.Load(".env")
	cfg := config.Load()

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	pool, err := db.New(ctx, cfg.DatabaseURL)
	if err != nil {
		log.Fatalf("подключение к базе: %v", err)
	}
	defer pool.Close()

	if err := db.RunMigrations(ctx, pool); err != nil {
		log.Fatalf("миграции: %v", err)
	}

	reader := bufio.NewReader(os.Stdin)

	if strings.TrimSpace(*email) == "" {
		*email = prompt(reader, "Почта администратора: ")
	}
	if strings.TrimSpace(*name) == "" {
		*name = prompt(reader, "Имя администратора: ")
	}
	if strings.TrimSpace(*password) == "" {
		*password = promptPassword("Пароль: ")
		if confirm := promptPassword("Повторите пароль: "); confirm != *password {
			log.Fatal("пароли не совпадают")
		}
	}

	if err := auth.ValidatePassword(*password); err != nil {
		log.Fatalf("пароль: %v", err)
	}

	hash, err := auth.HashPassword(*password)
	if err != nil {
		log.Fatalf("хеширование пароля: %v", err)
	}

	users := repository.NewUserRepo(pool)
	normalized := strings.ToLower(strings.TrimSpace(*email))

	existing, err := users.GetByEmail(ctx, normalized)
	switch {
	case err == nil:
		if !*force {
			log.Fatalf("пользователь %s уже существует — запустите с флагом -force, чтобы обновить его", normalized)
		}
		if err := users.SetPassword(ctx, existing.ID, hash); err != nil {
			log.Fatalf("обновление пароля: %v", err)
		}
		adminRole, active := domain.RoleAdmin, domain.StatusActive
		fullName := strings.TrimSpace(*name)
		if _, err := users.Update(ctx, existing.ID, repository.UpdateUserInput{
			FullName: &fullName,
			Role:     &adminRole,
			Status:   &active,
		}); err != nil {
			log.Fatalf("обновление пользователя: %v", err)
		}
		fmt.Printf("\n✓ Администратор обновлён: %s\n", normalized)

	case errors.Is(err, repository.ErrNotFound):
		user, err := users.Create(ctx, repository.CreateUserInput{
			Email:         normalized,
			FullName:      strings.TrimSpace(*name),
			PasswordHash:  hash,
			Role:          domain.RoleAdmin,
			Status:        domain.StatusActive,
			EmailVerified: true,
		})
		if err != nil {
			log.Fatalf("создание администратора: %v", err)
		}
		fmt.Printf("\n✓ Администратор создан: %s (id=%s)\n", user.Email, user.ID)

	default:
		log.Fatalf("поиск пользователя: %v", err)
	}

	fmt.Println("  Войдите на платформу этой почтой и паролем.")
}

func prompt(reader *bufio.Reader, label string) string {
	fmt.Print(label)
	line, err := reader.ReadString('\n')
	if err != nil {
		log.Fatalf("чтение ввода: %v", err)
	}
	value := strings.TrimSpace(line)
	if value == "" {
		log.Fatal("значение не может быть пустым")
	}
	return value
}

func promptPassword(label string) string {
	fmt.Print(label)
	raw, err := term.ReadPassword(int(syscall.Stdin))
	fmt.Println()
	if err != nil {
		log.Fatalf("чтение пароля: %v", err)
	}
	return strings.TrimSpace(string(raw))
}
