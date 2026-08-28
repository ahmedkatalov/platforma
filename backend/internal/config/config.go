package config

import (
	"errors"
	"os"
	"strconv"
	"strings"
	"time"
)

// Config — вся конфигурация сервиса, читается из окружения (.env подхватывается
// в main через loadDotEnv).
type Config struct {
	Port        string
	AppEnv      string
	DatabaseURL string

	JWTSecret     string
	AccessTTL     time.Duration
	RefreshTTL    time.Duration
	CorsOrigins   []string
	PublicBaseURL string
	UploadDir     string
	ReminderEvery time.Duration

	// EmailJS — отправка кодов подтверждения на реальную почту.
	EmailJSServiceID  string
	EmailJSTemplateID string
	EmailJSPublicKey  string
	EmailJSPrivateKey string
	EmailJSFromName   string
	// Шаблон для уведомлений (сертификаты, дедлайны). По умолчанию — тот же, что для кодов.
	EmailJSNoticeTemplateID string
	VerificationCodeTTL     time.Duration
}

func Load() *Config {
	defaultCors := []string{
		"http://localhost:5173",
		"http://127.0.0.1:5173",
		"http://localhost:4173",
		"http://127.0.0.1:4173",
	}

	return &Config{
		Port:        getEnv("PORT", "8090"),
		AppEnv:      getEnv("APP_ENV", "development"),
		DatabaseURL: getEnv("DATABASE_URL", "postgres://platforma:platforma@localhost:5433/platforma?sslmode=disable"),

		JWTSecret:     strings.TrimSpace(os.Getenv("JWT_SECRET")),
		AccessTTL:     getEnvDuration("ACCESS_TOKEN_TTL", 30*time.Minute),
		RefreshTTL:    getEnvDuration("REFRESH_TOKEN_TTL", 30*24*time.Hour),
		CorsOrigins:   getEnvSlice("CORS_ORIGINS", defaultCors),
		PublicBaseURL: getEnv("PUBLIC_BASE_URL", "http://localhost:5173"),
		UploadDir:     getEnv("UPLOAD_DIR", "uploads"),
		ReminderEvery: getEnvDuration("REMINDER_INTERVAL", 6*time.Hour),

		EmailJSServiceID:  getEnv("EMAILJS_SERVICE_ID", ""),
		EmailJSTemplateID: getEnv("EMAILJS_TEMPLATE_ID", ""),
		EmailJSPublicKey:  getEnv("EMAILJS_PUBLIC_KEY", ""),
		EmailJSPrivateKey: getEnv("EMAILJS_PRIVATE_KEY", ""),
		EmailJSFromName:   getEnv("EMAILJS_FROM_NAME", "DevOps Platform"),
		EmailJSNoticeTemplateID: getEnv("EMAILJS_NOTICE_TEMPLATE_ID",
			getEnv("EMAILJS_TEMPLATE_ID", "")),
		VerificationCodeTTL: getEnvDuration("VERIFICATION_CODE_TTL", 15*time.Minute),
	}
}

func (c *Config) Validate() error {
	if strings.TrimSpace(c.DatabaseURL) == "" {
		return errors.New("DATABASE_URL is required")
	}
	if strings.TrimSpace(c.JWTSecret) == "" {
		return errors.New("JWT_SECRET is required")
	}
	if len(c.JWTSecret) < 32 {
		return errors.New("JWT_SECRET must be at least 32 characters")
	}
	return nil
}

// EmailJSEnabled — если ключи не заданы, коды подтверждения печатаются в лог
// (удобно на локальной разработке).
func (c *Config) EmailJSEnabled() bool {
	return strings.TrimSpace(c.EmailJSServiceID) != "" &&
		strings.TrimSpace(c.EmailJSTemplateID) != "" &&
		strings.TrimSpace(c.EmailJSPublicKey) != ""
}

func (c *Config) IsProduction() bool {
	return strings.EqualFold(c.AppEnv, "production")
}

func getEnv(key, fallback string) string {
	if v := strings.TrimSpace(os.Getenv(key)); v != "" {
		return v
	}
	return fallback
}

func getEnvDuration(key string, fallback time.Duration) time.Duration {
	raw := strings.TrimSpace(os.Getenv(key))
	if raw == "" {
		return fallback
	}
	if d, err := time.ParseDuration(raw); err == nil {
		return d
	}
	// Поддержка простого числа = секунды.
	if n, err := strconv.Atoi(raw); err == nil {
		return time.Duration(n) * time.Second
	}
	return fallback
}

func getEnvSlice(key string, fallback []string) []string {
	raw := strings.TrimSpace(os.Getenv(key))
	if raw == "" {
		return fallback
	}
	out := make([]string, 0, 4)
	for _, p := range strings.Split(raw, ",") {
		if s := strings.TrimSpace(p); s != "" {
			out = append(out, s)
		}
	}
	if len(out) == 0 {
		return fallback
	}
	return out
}
