package mailer

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"
	"time"

	"platforma/backend/internal/config"
)

const emailJSEndpoint = "https://api.emailjs.com/api/v1.0/email/send"

var ErrNotConfigured = errors.New("EmailJS не настроен")

// Mailer отправляет письма через EmailJS REST API (тот же способ, что в cmf:
// приватный ключ передаётся в accessToken, а Origin обязателен для strict-режима).
type Mailer struct {
	cfg    *config.Config
	client *http.Client
}

func New(cfg *config.Config) *Mailer {
	return &Mailer{cfg: cfg, client: &http.Client{Timeout: 15 * time.Second}}
}

func (m *Mailer) Enabled() bool { return m.cfg.EmailJSEnabled() }

type emailJSRequest struct {
	ServiceID      string         `json:"service_id"`
	TemplateID     string         `json:"template_id"`
	UserID         string         `json:"user_id"`
	AccessToken    string         `json:"accessToken,omitempty"`
	TemplateParams map[string]any `json:"template_params,omitempty"`
}

// SendCode отправляет одноразовый код подтверждения.
// В шаблон уходит несколько синонимов переменных — EmailJS подставит те,
// что реально есть в шаблоне.
func (m *Mailer) SendCode(ctx context.Context, to, code, purpose string) error {
	to = strings.TrimSpace(to)
	if to == "" {
		return errors.New("пустой адрес получателя")
	}
	if strings.TrimSpace(code) == "" {
		return errors.New("пустой код")
	}

	title := "Подтверждение регистрации"
	if purpose == "password_reset" {
		title = "Восстановление пароля"
	}

	minutes := int(m.cfg.VerificationCodeTTL.Minutes())

	return m.send(ctx, m.cfg.EmailJSTemplateID, to, map[string]any{
		"to_email":  to,
		"email":     to,
		"to_name":   to,
		"passcode":  code,
		"code":      code,
		"subject":   title,
		"title":     title,
		"from_name": m.cfg.EmailJSFromName,
		"expires":   fmt.Sprintf("%d", minutes),
		"message":   fmt.Sprintf("%s. Ваш код: %s. Код действует %d минут.", title, code, minutes),
	})
}

// SendInvite отправляет студенту приглашение с логином и временным паролем.
func (m *Mailer) SendInvite(ctx context.Context, to, fullName, tempPassword string) error {
	to = strings.TrimSpace(to)
	if to == "" {
		return errors.New("пустой адрес получателя")
	}

	name := strings.TrimSpace(fullName)
	if name == "" {
		name = to
	}

	message := fmt.Sprintf(
		"%s, для вас создан аккаунт на платформе. Логин: %s, временный пароль: %s. Войдите по адресу %s и смените пароль.",
		name, to, tempPassword, m.cfg.PublicBaseURL,
	)

	return m.send(ctx, m.cfg.EmailJSTemplateID, to, map[string]any{
		"to_email":  to,
		"email":     to,
		"to_name":   name,
		"passcode":  tempPassword,
		"code":      tempPassword,
		"subject":   "Доступ к платформе",
		"title":     "Доступ к платформе",
		"from_name": m.cfg.EmailJSFromName,
		"link":      m.cfg.PublicBaseURL,
		"message":   message,
	})
}

func (m *Mailer) send(ctx context.Context, templateID, to string, params map[string]any) error {
	if !m.Enabled() {
		return ErrNotConfigured
	}

	payload := emailJSRequest{
		ServiceID:      m.cfg.EmailJSServiceID,
		TemplateID:     templateID,
		UserID:         m.cfg.EmailJSPublicKey,
		AccessToken:    m.cfg.EmailJSPrivateKey,
		TemplateParams: params,
	}

	body, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("marshal emailjs payload: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, emailJSEndpoint, bytes.NewReader(body))
	if err != nil {
		return fmt.Errorf("build emailjs request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Origin", "http://localhost")

	resp, err := m.client.Do(req)
	if err != nil {
		return fmt.Errorf("emailjs request failed: %w", err)
	}
	defer resp.Body.Close()

	respBody, _ := io.ReadAll(resp.Body)
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return fmt.Errorf("emailjs send failed: status=%d body=%s", resp.StatusCode, strings.TrimSpace(string(respBody)))
	}

	log.Printf("mailer: письмо отправлено to=%s status=%d", to, resp.StatusCode)
	return nil
}
