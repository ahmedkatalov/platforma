package domain

import (
	"encoding/json"
	"time"
)

// Роли и статусы пользователей.
const (
	RoleAdmin   = "admin"
	RoleStudent = "student"

	StatusInvited = "invited"
	StatusActive  = "active"
	StatusBlocked = "blocked"
)

// Назначение одноразового кода на почту.
const (
	PurposeRegistration  = "registration"
	PurposePasswordReset = "password_reset"
)

// Виды уроков.
const (
	LessonText     = "text"
	LessonQuiz     = "quiz"
	LessonTerminal = "terminal"
	LessonCode     = "code"
)

type User struct {
	ID            string     `json:"id"`
	Email         string     `json:"email"`
	FullName      string     `json:"fullName"`
	Role          string     `json:"role"`
	Status        string     `json:"status"`
	EmailVerified bool       `json:"emailVerified"`
	AvatarURL     string     `json:"avatarUrl"`
	LastLoginAt   *time.Time `json:"lastLoginAt"`
	CreatedAt     time.Time  `json:"createdAt"`
	UpdatedAt     time.Time  `json:"updatedAt"`

	PasswordHash string `json:"-"`
}

func (u *User) IsAdmin() bool  { return u.Role == RoleAdmin }
func (u *User) IsActive() bool { return u.Status == StatusActive }

type Course struct {
	ID          string    `json:"id"`
	Slug        string    `json:"slug"`
	Title       string    `json:"title"`
	Subtitle    string    `json:"subtitle"`
	Description string    `json:"description"`
	CoverURL    string    `json:"coverUrl"`
	Level       string    `json:"level"`
	Tags        []string  `json:"tags"`
	Status      string    `json:"status"`
	Position    int       `json:"position"`
	CreatedAt   time.Time `json:"createdAt"`
	UpdatedAt   time.Time `json:"updatedAt"`

	// Заполняются в развёрнутых ответах.
	Modules       []Module `json:"modules,omitempty"`
	ModulesCount  int      `json:"modulesCount"`
	LessonsCount  int      `json:"lessonsCount"`
	StudentsCount int      `json:"studentsCount"`
}

type Module struct {
	ID        string    `json:"id"`
	CourseID  string    `json:"courseId"`
	Title     string    `json:"title"`
	Summary   string    `json:"summary"`
	Position  int       `json:"position"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`

	Lessons []Lesson `json:"lessons,omitempty"`
}

type Lesson struct {
	ID          string          `json:"id"`
	ModuleID    string          `json:"moduleId"`
	Title       string          `json:"title"`
	Kind        string          `json:"kind"`
	Summary     string          `json:"summary"`
	Content     json.RawMessage `json:"content"`
	DurationMin int             `json:"durationMin"`
	Position    int             `json:"position"`
	CreatedAt   time.Time       `json:"createdAt"`
	UpdatedAt   time.Time       `json:"updatedAt"`
}

type Enrollment struct {
	ID          string     `json:"id"`
	UserID      string     `json:"userId"`
	CourseID    string     `json:"courseId"`
	Status      string     `json:"status"`
	DueDate     *time.Time `json:"dueDate"`
	StartedAt   *time.Time `json:"startedAt"`
	CompletedAt *time.Time `json:"completedAt"`
	CreatedAt   time.Time  `json:"createdAt"`

	CourseTitle string `json:"courseTitle,omitempty"`
	CourseSlug  string `json:"courseSlug,omitempty"`
}

type ActivityDay struct {
	Day          string `json:"day"`
	Visits       int    `json:"visits"`
	SecondsSpent int    `json:"secondsSpent"`
}
