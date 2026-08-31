package service

import (
	"context"
	"errors"
	"log"
	"strings"
	"time"

	"platforma/backend/internal/auth"
	"platforma/backend/internal/config"
	"platforma/backend/internal/domain"
	"platforma/backend/internal/mailer"
	"platforma/backend/internal/repository"
)

var (
	ErrInvalidCredentials = errors.New("Неверная почта или пароль")
	ErrUserBlocked        = errors.New("Аккаунт заблокирован, обратитесь к администратору")
	ErrInvalidCode        = errors.New("Неверный код подтверждения")
	ErrCodeExpired        = errors.New("Код истёк, запросите новый")
	ErrTooManyAttempts    = errors.New("Слишком много попыток, запросите новый код")
	ErrTooManyRequests    = errors.New("Слишком часто — подождите минуту и попробуйте снова")
	ErrEmailTaken         = repository.ErrEmailTaken
)

// AuthService — регистрация с подтверждением почты, вход, refresh и сброс пароля.
type AuthService struct {
	cfg      *config.Config
	users    *repository.UserRepo
	tokens   *repository.TokenRepo
	codes    *repository.CodeRepo
	activity *repository.ActivityRepo
	tm       *auth.TokenManager
	mail     *mailer.Mailer
}

func NewAuthService(
	cfg *config.Config,
	users *repository.UserRepo,
	tokens *repository.TokenRepo,
	codes *repository.CodeRepo,
	activity *repository.ActivityRepo,
	tm *auth.TokenManager,
	mail *mailer.Mailer,
) *AuthService {
	return &AuthService{cfg: cfg, users: users, tokens: tokens, codes: codes,
		activity: activity, tm: tm, mail: mail}
}

// Session — то, что уходит на фронтенд после успешного входа.
type Session struct {
	AccessToken  string      `json:"accessToken"`
	RefreshToken string      `json:"refreshToken"`
	ExpiresAt    time.Time   `json:"expiresAt"`
	User         domain.User `json:"user"`
}

// SendCode отправляет одноразовый код на почту.
// purpose: registration — почта должна быть свободна; password_reset — аккаунт должен существовать.
func (s *AuthService) SendCode(ctx context.Context, email, purpose string) error {
	email = normalizeEmail(email)
	if !looksLikeEmail(email) {
		return errors.New("укажите корректный адрес почты")
	}

	exists, err := s.users.EmailExists(ctx, email)
	if err != nil {
		return err
	}

	switch purpose {
	case domain.PurposeRegistration:
		if exists {
			return ErrEmailTaken
		}
	case domain.PurposePasswordReset:
		// Не раскрываем, зарегистрирована ли почта.
		if !exists {
			return nil
		}
	default:
		return errors.New("неизвестное назначение кода")
	}

	// Не больше 3 писем за 10 минут на один адрес.
	recent, err := s.codes.SentRecently(ctx, email, purpose, 10*time.Minute)
	if err != nil {
		return err
	}
	if recent >= 3 {
		return ErrTooManyRequests
	}

	code, err := auth.NewNumericCode(6)
	if err != nil {
		return err
	}
	if err := s.codes.Create(ctx, email, purpose, auth.HashCode(code),
		time.Now().Add(s.cfg.VerificationCodeTTL)); err != nil {
		return err
	}

	if !s.mail.Enabled() {
		// Локальная разработка без ключей EmailJS: код виден в логе сервера.
		log.Printf("[DEV] код подтверждения для %s (%s): %s", email, purpose, code)
		return nil
	}
	return s.mail.SendCode(ctx, email, code, purpose)
}

// verifyCode проверяет код и помечает его использованным.
func (s *AuthService) verifyCode(ctx context.Context, email, purpose, code string) error {
	rec, err := s.codes.GetActive(ctx, email, purpose)
	if errors.Is(err, repository.ErrNotFound) {
		return ErrCodeExpired
	}
	if err != nil {
		return err
	}
	if rec.Attempts >= repository.MaxCodeAttempts {
		return ErrTooManyAttempts
	}
	if rec.CodeHash != auth.HashCode(strings.TrimSpace(code)) {
		_ = s.codes.IncAttempts(ctx, rec.ID)
		return ErrInvalidCode
	}
	return s.codes.Consume(ctx, rec.ID)
}

type RegisterInput struct {
	Email    string `json:"email"`
	FullName string `json:"fullName"`
	Password string `json:"password"`
	Code     string `json:"code"`
}

// Register завершает самостоятельную регистрацию студента: почта уже подтверждена кодом.
func (s *AuthService) Register(ctx context.Context, in RegisterInput, userAgent, ip string) (*Session, error) {
	email := normalizeEmail(in.Email)
	if !looksLikeEmail(email) {
		return nil, errors.New("укажите корректный адрес почты")
	}
	if strings.TrimSpace(in.FullName) == "" {
		return nil, errors.New("укажите имя")
	}
	if err := auth.ValidatePassword(in.Password); err != nil {
		return nil, err
	}
	if err := s.verifyCode(ctx, email, domain.PurposeRegistration, in.Code); err != nil {
		return nil, err
	}

	hash, err := auth.HashPassword(in.Password)
	if err != nil {
		return nil, err
	}

	user, err := s.users.Create(ctx, repository.CreateUserInput{
		Email:         email,
		FullName:      strings.TrimSpace(in.FullName),
		PasswordHash:  hash,
		Role:          domain.RoleStudent,
		Status:        domain.StatusActive,
		EmailVerified: true,
	})
	if err != nil {
		return nil, err
	}

	return s.issueSession(ctx, user, userAgent, ip)
}

func (s *AuthService) Login(ctx context.Context, email, password, userAgent, ip string) (*Session, error) {
	user, err := s.users.GetByEmail(ctx, normalizeEmail(email))
	if errors.Is(err, repository.ErrNotFound) {
		return nil, ErrInvalidCredentials
	}
	if err != nil {
		return nil, err
	}
	if !auth.CheckPassword(user.PasswordHash, password) {
		return nil, ErrInvalidCredentials
	}
	if user.Status == domain.StatusBlocked {
		return nil, ErrUserBlocked
	}

	// Приглашённый студент активируется первым входом по временному паролю.
	if user.Status == domain.StatusInvited {
		if _, err := s.users.Update(ctx, user.ID, repository.UpdateUserInput{
			Status: strPtr(domain.StatusActive),
		}); err != nil {
			return nil, err
		}
		user.Status = domain.StatusActive
	}

	return s.issueSession(ctx, user, userAgent, ip)
}

// Refresh обменивает refresh-токен на новую пару (старый отзывается).
func (s *AuthService) Refresh(ctx context.Context, refreshToken, userAgent, ip string) (*Session, error) {
	hash := auth.HashToken(strings.TrimSpace(refreshToken))
	rec, err := s.tokens.GetValid(ctx, hash)
	if errors.Is(err, repository.ErrNotFound) {
		return nil, errors.New("сессия истекла, войдите заново")
	}
	if err != nil {
		return nil, err
	}

	user, err := s.users.GetByID(ctx, rec.UserID)
	if err != nil {
		return nil, err
	}
	if user.Status == domain.StatusBlocked {
		return nil, ErrUserBlocked
	}

	if err := s.tokens.Revoke(ctx, hash); err != nil {
		return nil, err
	}
	return s.issueSession(ctx, user, userAgent, ip)
}

func (s *AuthService) Logout(ctx context.Context, refreshToken string) error {
	if strings.TrimSpace(refreshToken) == "" {
		return nil
	}
	return s.tokens.Revoke(ctx, auth.HashToken(strings.TrimSpace(refreshToken)))
}

type ResetPasswordInput struct {
	Email    string `json:"email"`
	Code     string `json:"code"`
	Password string `json:"password"`
}

// ResetPassword меняет пароль по коду с почты и разлогинивает все сессии.
func (s *AuthService) ResetPassword(ctx context.Context, in ResetPasswordInput) error {
	email := normalizeEmail(in.Email)
	if err := auth.ValidatePassword(in.Password); err != nil {
		return err
	}
	if err := s.verifyCode(ctx, email, domain.PurposePasswordReset, in.Code); err != nil {
		return err
	}

	user, err := s.users.GetByEmail(ctx, email)
	if err != nil {
		return err
	}
	hash, err := auth.HashPassword(in.Password)
	if err != nil {
		return err
	}
	if err := s.users.SetPassword(ctx, user.ID, hash); err != nil {
		return err
	}
	return s.tokens.RevokeAllForUser(ctx, user.ID)
}

// ChangePassword — смена пароля авторизованным пользователем.
func (s *AuthService) ChangePassword(ctx context.Context, userID, current, next string) error {
	user, err := s.users.GetByID(ctx, userID)
	if err != nil {
		return err
	}
	if !auth.CheckPassword(user.PasswordHash, current) {
		return errors.New("текущий пароль указан неверно")
	}
	if err := auth.ValidatePassword(next); err != nil {
		return err
	}
	hash, err := auth.HashPassword(next)
	if err != nil {
		return err
	}
	return s.users.SetPassword(ctx, user.ID, hash)
}

func (s *AuthService) issueSession(ctx context.Context, user *domain.User, userAgent, ip string) (*Session, error) {
	access, expires, err := s.tm.IssueAccess(user.ID, user.Role, user.Email)
	if err != nil {
		return nil, err
	}

	refresh, err := auth.NewRefreshToken()
	if err != nil {
		return nil, err
	}
	if err := s.tokens.Create(ctx, user.ID, auth.HashToken(refresh), userAgent, ip,
		time.Now().Add(s.tm.RefreshTTL())); err != nil {
		return nil, err
	}

	_ = s.users.TouchLogin(ctx, user.ID)
	_ = s.activity.Touch(ctx, user.ID, 0)

	return &Session{
		AccessToken:  access,
		RefreshToken: refresh,
		ExpiresAt:    expires,
		User:         *user,
	}, nil
}

func normalizeEmail(email string) string {
	return strings.ToLower(strings.TrimSpace(email))
}

func looksLikeEmail(email string) bool {
	at := strings.Index(email, "@")
	dot := strings.LastIndex(email, ".")
	return at > 0 && dot > at+1 && dot < len(email)-1 && !strings.ContainsAny(email, " \t\n")
}

func strPtr(s string) *string { return &s }
