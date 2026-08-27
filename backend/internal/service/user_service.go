package service

import (
	"context"
	"crypto/rand"
	"errors"
	"log"
	"math/big"
	"strings"

	"platforma/backend/internal/auth"
	"platforma/backend/internal/domain"
	"platforma/backend/internal/mailer"
	"platforma/backend/internal/repository"
)

// UserService — управление аккаунтами студентов со стороны администратора.
type UserService struct {
	users  *repository.UserRepo
	tokens *repository.TokenRepo
	mail   *mailer.Mailer
	audit  *repository.AuditRepo
}

func NewUserService(
	users *repository.UserRepo,
	tokens *repository.TokenRepo,
	mail *mailer.Mailer,
	audit *repository.AuditRepo,
) *UserService {
	return &UserService{users: users, tokens: tokens, mail: mail, audit: audit}
}

type CreateStudentInput struct {
	Email    string `json:"email"`
	FullName string `json:"fullName"`
	Password string `json:"password"` // необязательно: пустой = сгенерировать
	Role     string `json:"role"`     // student по умолчанию
	SendMail bool   `json:"sendMail"`
}

type CreatedStudent struct {
	User         domain.User `json:"user"`
	TempPassword string      `json:"tempPassword"`
	MailSent     bool        `json:"mailSent"`
	MailError    string      `json:"mailError,omitempty"`
}

// CreateStudent заводит аккаунт и (по желанию) отправляет доступы на почту.
func (s *UserService) CreateStudent(ctx context.Context, actorID string, in CreateStudentInput) (*CreatedStudent, error) {
	email := normalizeEmail(in.Email)
	if !looksLikeEmail(email) {
		return nil, errors.New("укажите корректный адрес почты")
	}
	if strings.TrimSpace(in.FullName) == "" {
		return nil, errors.New("укажите имя студента")
	}

	role := domain.RoleStudent
	if in.Role == domain.RoleAdmin {
		role = domain.RoleAdmin
	}

	password := strings.TrimSpace(in.Password)
	if password == "" {
		generated, err := generatePassword(12)
		if err != nil {
			return nil, err
		}
		password = generated
	} else if err := auth.ValidatePassword(password); err != nil {
		return nil, err
	}

	hash, err := auth.HashPassword(password)
	if err != nil {
		return nil, err
	}

	user, err := s.users.Create(ctx, repository.CreateUserInput{
		Email:        email,
		FullName:     strings.TrimSpace(in.FullName),
		PasswordHash: hash,
		Role:         role,
		Status:       domain.StatusInvited,
		CreatedBy:    &actorID,
	})
	if err != nil {
		return nil, err
	}

	out := &CreatedStudent{User: *user, TempPassword: password}

	if in.SendMail {
		if !s.mail.Enabled() {
			out.MailError = "EmailJS не настроен — передайте пароль студенту вручную"
			log.Printf("[DEV] аккаунт %s создан, временный пароль: %s", email, password)
		} else if err := s.mail.SendInvite(ctx, email, user.FullName, password); err != nil {
			out.MailError = err.Error()
		} else {
			out.MailSent = true
		}
	}

	s.audit.Log(ctx, actorID, "user.create", "user", user.ID, map[string]any{
		"email": email, "role": role,
	})
	return out, nil
}

func (s *UserService) Update(ctx context.Context, actorID, userID string, in repository.UpdateUserInput) (*domain.User, error) {
	// Не даём заблокировать или разжаловать последнего администратора.
	if in.Role != nil || in.Status != nil {
		current, err := s.users.GetByID(ctx, userID)
		if err != nil {
			return nil, err
		}
		if current.Role == domain.RoleAdmin {
			losingAdmin := (in.Role != nil && *in.Role != domain.RoleAdmin) ||
				(in.Status != nil && *in.Status == domain.StatusBlocked)
			if losingAdmin {
				if err := s.ensureAnotherAdmin(ctx); err != nil {
					return nil, err
				}
			}
		}
	}

	user, err := s.users.Update(ctx, userID, in)
	if err != nil {
		return nil, err
	}

	// Заблокированного пользователя выкидываем из всех сессий.
	if in.Status != nil && *in.Status == domain.StatusBlocked {
		_ = s.tokens.RevokeAllForUser(ctx, userID)
	}

	s.audit.Log(ctx, actorID, "user.update", "user", userID, in)
	return user, nil
}

// ResetPasswordByAdmin выдаёт новый пароль вместо забытого.
func (s *UserService) ResetPasswordByAdmin(ctx context.Context, actorID, userID string, sendMail bool) (*CreatedStudent, error) {
	user, err := s.users.GetByID(ctx, userID)
	if err != nil {
		return nil, err
	}

	password, err := generatePassword(12)
	if err != nil {
		return nil, err
	}
	hash, err := auth.HashPassword(password)
	if err != nil {
		return nil, err
	}
	if err := s.users.SetPassword(ctx, userID, hash); err != nil {
		return nil, err
	}
	_ = s.tokens.RevokeAllForUser(ctx, userID)

	out := &CreatedStudent{User: *user, TempPassword: password}
	if sendMail {
		if !s.mail.Enabled() {
			out.MailError = "EmailJS не настроен — передайте пароль вручную"
		} else if err := s.mail.SendInvite(ctx, user.Email, user.FullName, password); err != nil {
			out.MailError = err.Error()
		} else {
			out.MailSent = true
		}
	}

	s.audit.Log(ctx, actorID, "user.reset_password", "user", userID, nil)
	return out, nil
}

func (s *UserService) Delete(ctx context.Context, actorID, userID string) error {
	if actorID == userID {
		return errors.New("нельзя удалить собственный аккаунт")
	}
	user, err := s.users.GetByID(ctx, userID)
	if err != nil {
		return err
	}
	if user.Role == domain.RoleAdmin {
		if err := s.ensureAnotherAdmin(ctx); err != nil {
			return err
		}
	}
	if err := s.users.Delete(ctx, userID); err != nil {
		return err
	}
	s.audit.Log(ctx, actorID, "user.delete", "user", userID, map[string]any{"email": user.Email})
	return nil
}

// ensureAnotherAdmin не даёт остаться без администратора.
func (s *UserService) ensureAnotherAdmin(ctx context.Context) error {
	admins, err := s.users.CountByRole(ctx, domain.RoleAdmin)
	if err != nil {
		return err
	}
	if admins <= 1 {
		return repository.ErrLastAdminGone
	}
	return nil
}

const passwordAlphabet = "abcdefghijkmnpqrstuvwxyzABCDEFGHJKLMNPQRSTUVWXYZ23456789"

// generatePassword — читаемый временный пароль без похожих символов.
func generatePassword(length int) (string, error) {
	if length < 8 {
		length = 8
	}
	buf := make([]byte, length)
	max := big.NewInt(int64(len(passwordAlphabet)))
	for i := range buf {
		n, err := rand.Int(rand.Reader, max)
		if err != nil {
			return "", err
		}
		buf[i] = passwordAlphabet[n.Int64()]
	}
	// Гарантируем цифру — требование ValidatePassword.
	buf[length-1] = "23456789"[int(buf[0])%8]
	return string(buf), nil
}
