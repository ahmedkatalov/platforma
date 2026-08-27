package router

import (
	"net/http"
	"time"

	"platforma/backend/internal/auth"
	"platforma/backend/internal/config"
	"platforma/backend/internal/handler"
	appmw "platforma/backend/internal/middleware"

	"github.com/go-chi/chi/v5"
	chimw "github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
)

// Deps — всё, что нужно роутеру для сборки дерева маршрутов.
type Deps struct {
	Config  *config.Config
	Tokens  *auth.TokenManager
	Auth    *handler.AuthHandler
	Me      *handler.MeHandler
	Admin   *handler.AdminHandler
	Courses *handler.CourseHandler
	Theme   *handler.ThemeHandler
}

func New(d Deps) http.Handler {
	r := chi.NewRouter()

	r.Use(chimw.RequestID)
	r.Use(chimw.RealIP)
	r.Use(chimw.Logger)
	r.Use(chimw.Recoverer)
	r.Use(chimw.Timeout(30 * time.Second))
	r.Use(cors.Handler(cors.Options{
		AllowedOrigins:   d.Config.CorsOrigins,
		AllowedMethods:   []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type"},
		AllowCredentials: true,
		MaxAge:           300,
	}))

	r.Get("/health", func(w http.ResponseWriter, _ *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"status":"ok"}`))
	})

	r.Route("/api", func(r chi.Router) {
		// Публичное.
		r.Mount("/auth", d.Auth.Routes())
		r.Get("/theme", d.Theme.Public)

		// Требует авторизации.
		r.Group(func(r chi.Router) {
			r.Use(appmw.Auth(d.Tokens))

			r.Mount("/me", d.Me.Routes())
			r.Mount("/courses", d.Courses.StudentRoutes())

			// Только администратор.
			r.Group(func(r chi.Router) {
				r.Use(appmw.RequireAdmin)
				r.Mount("/admin", d.Admin.Routes(d.Courses.AdminRoutes()))
			})
		})
	})

	return r
}
