package repository

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"strings"
	"time"

	"platforma/backend/internal/domain"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

const courseColumns = `id, slug, title, subtitle, description, cover_url, level, tags,
	status, position, created_at, updated_at`

type CourseRepo struct{ db *pgxpool.Pool }

func NewCourseRepo(db *pgxpool.Pool) *CourseRepo { return &CourseRepo{db: db} }

func scanCourse(row pgx.Row) (*domain.Course, error) {
	var c domain.Course
	err := row.Scan(&c.ID, &c.Slug, &c.Title, &c.Subtitle, &c.Description, &c.CoverURL,
		&c.Level, &c.Tags, &c.Status, &c.Position, &c.CreatedAt, &c.UpdatedAt)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, ErrNotFound
	}
	if err != nil {
		return nil, err
	}
	return &c, nil
}

type CourseInput struct {
	Slug        string   `json:"slug"`
	Title       string   `json:"title"`
	Subtitle    string   `json:"subtitle"`
	Description string   `json:"description"`
	CoverURL    string   `json:"coverUrl"`
	Level       string   `json:"level"`
	Tags        []string `json:"tags"`
	Status      string   `json:"status"`
	Position    int      `json:"position"`
}

func (r *CourseRepo) Create(ctx context.Context, in CourseInput, createdBy string) (*domain.Course, error) {
	if in.Tags == nil {
		in.Tags = []string{}
	}
	var author *string
	if createdBy != "" {
		author = &createdBy
	}
	c, err := scanCourse(r.db.QueryRow(ctx, `
		INSERT INTO courses (slug, title, subtitle, description, cover_url, level, tags, status, position, created_by)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)
		RETURNING `+courseColumns,
		in.Slug, in.Title, in.Subtitle, in.Description, in.CoverURL, in.Level,
		in.Tags, in.Status, in.Position, author))
	if err != nil && isUniqueViolation(err) {
		return nil, errors.New("курс с таким slug уже существует")
	}
	return c, err
}

func (r *CourseRepo) Update(ctx context.Context, id string, in CourseInput) (*domain.Course, error) {
	if in.Tags == nil {
		in.Tags = []string{}
	}
	c, err := scanCourse(r.db.QueryRow(ctx, `
		UPDATE courses
		   SET slug = $2, title = $3, subtitle = $4, description = $5, cover_url = $6,
		       level = $7, tags = $8, status = $9, position = $10, updated_at = now()
		 WHERE id = $1
		RETURNING `+courseColumns,
		id, in.Slug, in.Title, in.Subtitle, in.Description, in.CoverURL,
		in.Level, in.Tags, in.Status, in.Position))
	if err != nil && isUniqueViolation(err) {
		return nil, errors.New("курс с таким slug уже существует")
	}
	return c, err
}

func (r *CourseRepo) Delete(ctx context.Context, id string) error {
	tag, err := r.db.Exec(ctx, `DELETE FROM courses WHERE id = $1`, id)
	if err != nil {
		return err
	}
	if tag.RowsAffected() == 0 {
		return ErrNotFound
	}
	return nil
}

func (r *CourseRepo) GetByID(ctx context.Context, id string) (*domain.Course, error) {
	return scanCourse(r.db.QueryRow(ctx, `SELECT `+courseColumns+` FROM courses WHERE id = $1`, id))
}

func (r *CourseRepo) GetBySlug(ctx context.Context, slug string) (*domain.Course, error) {
	return scanCourse(r.db.QueryRow(ctx, `SELECT `+courseColumns+` FROM courses WHERE slug = $1`, slug))
}

// List отдаёт курсы со счётчиками. Пустой status = все статусы.
func (r *CourseRepo) List(ctx context.Context, status string) ([]domain.Course, error) {
	where := "1 = 1"
	args := []any{}
	if s := strings.TrimSpace(status); s != "" {
		args = append(args, s)
		where = fmt.Sprintf("c.status = $%d", len(args))
	}

	rows, err := r.db.Query(ctx, fmt.Sprintf(`
		SELECT c.id, c.slug, c.title, c.subtitle, c.description, c.cover_url, c.level, c.tags,
		       c.status, c.position, c.created_at, c.updated_at,
		       (SELECT count(*) FROM modules m WHERE m.course_id = c.id),
		       (SELECT count(*) FROM lessons l JOIN modules m ON m.id = l.module_id WHERE m.course_id = c.id),
		       (SELECT count(*) FROM enrollments e WHERE e.course_id = c.id)
		  FROM courses c
		 WHERE %s
		 ORDER BY c.position, c.created_at`, where), args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	out := make([]domain.Course, 0, 8)
	for rows.Next() {
		var c domain.Course
		if err := rows.Scan(&c.ID, &c.Slug, &c.Title, &c.Subtitle, &c.Description, &c.CoverURL,
			&c.Level, &c.Tags, &c.Status, &c.Position, &c.CreatedAt, &c.UpdatedAt,
			&c.ModulesCount, &c.LessonsCount, &c.StudentsCount); err != nil {
			return nil, err
		}
		out = append(out, c)
	}
	return out, rows.Err()
}

// WithContent подгружает модули и уроки курса.
// includeContent=false отдаёт уроки без тела (для списков и навигации).
func (r *CourseRepo) WithContent(ctx context.Context, course *domain.Course, includeContent bool) error {
	modRows, err := r.db.Query(ctx, `
		SELECT id, course_id, title, summary, position, created_at, updated_at
		  FROM modules WHERE course_id = $1 ORDER BY position, created_at`, course.ID)
	if err != nil {
		return err
	}
	defer modRows.Close()

	modules := make([]domain.Module, 0, 8)
	index := map[string]int{}
	for modRows.Next() {
		var m domain.Module
		if err := modRows.Scan(&m.ID, &m.CourseID, &m.Title, &m.Summary, &m.Position,
			&m.CreatedAt, &m.UpdatedAt); err != nil {
			return err
		}
		m.Lessons = []domain.Lesson{}
		index[m.ID] = len(modules)
		modules = append(modules, m)
	}
	if err := modRows.Err(); err != nil {
		return err
	}

	contentExpr := "'{}'::jsonb"
	if includeContent {
		contentExpr = "l.content"
	}

	lessonRows, err := r.db.Query(ctx, fmt.Sprintf(`
		SELECT l.id, l.module_id, l.title, l.kind, l.summary, %s, l.duration_min, l.position,
		       l.created_at, l.updated_at
		  FROM lessons l
		  JOIN modules m ON m.id = l.module_id
		 WHERE m.course_id = $1
		 ORDER BY m.position, l.position, l.created_at`, contentExpr), course.ID)
	if err != nil {
		return err
	}
	defer lessonRows.Close()

	for lessonRows.Next() {
		var l domain.Lesson
		if err := lessonRows.Scan(&l.ID, &l.ModuleID, &l.Title, &l.Kind, &l.Summary, &l.Content,
			&l.DurationMin, &l.Position, &l.CreatedAt, &l.UpdatedAt); err != nil {
			return err
		}
		if i, ok := index[l.ModuleID]; ok {
			modules[i].Lessons = append(modules[i].Lessons, l)
		}
	}
	if err := lessonRows.Err(); err != nil {
		return err
	}

	course.Modules = modules
	course.ModulesCount = len(modules)
	for _, m := range modules {
		course.LessonsCount += len(m.Lessons)
	}
	return nil
}

// --- Модули ---

type ModuleInput struct {
	Title    string `json:"title"`
	Summary  string `json:"summary"`
	Position int    `json:"position"`
}

func (r *CourseRepo) CreateModule(ctx context.Context, courseID string, in ModuleInput) (*domain.Module, error) {
	var m domain.Module
	err := r.db.QueryRow(ctx, `
		INSERT INTO modules (course_id, title, summary, position)
		VALUES ($1, $2, $3, COALESCE(NULLIF($4, 0),
		        (SELECT COALESCE(max(position), 0) + 1 FROM modules WHERE course_id = $1)))
		RETURNING id, course_id, title, summary, position, created_at, updated_at`,
		courseID, in.Title, in.Summary, in.Position).
		Scan(&m.ID, &m.CourseID, &m.Title, &m.Summary, &m.Position, &m.CreatedAt, &m.UpdatedAt)
	if err != nil {
		return nil, err
	}
	return &m, nil
}

func (r *CourseRepo) UpdateModule(ctx context.Context, id string, in ModuleInput) (*domain.Module, error) {
	var m domain.Module
	err := r.db.QueryRow(ctx, `
		UPDATE modules SET title = $2, summary = $3, position = $4, updated_at = now()
		 WHERE id = $1
		RETURNING id, course_id, title, summary, position, created_at, updated_at`,
		id, in.Title, in.Summary, in.Position).
		Scan(&m.ID, &m.CourseID, &m.Title, &m.Summary, &m.Position, &m.CreatedAt, &m.UpdatedAt)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, ErrNotFound
	}
	if err != nil {
		return nil, err
	}
	return &m, nil
}

func (r *CourseRepo) DeleteModule(ctx context.Context, id string) error {
	tag, err := r.db.Exec(ctx, `DELETE FROM modules WHERE id = $1`, id)
	if err != nil {
		return err
	}
	if tag.RowsAffected() == 0 {
		return ErrNotFound
	}
	return nil
}

// --- Уроки ---

type LessonInput struct {
	Title       string          `json:"title"`
	Kind        string          `json:"kind"`
	Summary     string          `json:"summary"`
	Content     json.RawMessage `json:"content"`
	DurationMin int             `json:"durationMin"`
	Position    int             `json:"position"`
}

func (r *CourseRepo) CreateLesson(ctx context.Context, moduleID string, in LessonInput) (*domain.Lesson, error) {
	if len(in.Content) == 0 {
		in.Content = json.RawMessage(`{}`)
	}
	var l domain.Lesson
	err := r.db.QueryRow(ctx, `
		INSERT INTO lessons (module_id, title, kind, summary, content, duration_min, position)
		VALUES ($1, $2, $3, $4, $5, $6, COALESCE(NULLIF($7, 0),
		        (SELECT COALESCE(max(position), 0) + 1 FROM lessons WHERE module_id = $1)))
		RETURNING id, module_id, title, kind, summary, content, duration_min, position, created_at, updated_at`,
		moduleID, in.Title, in.Kind, in.Summary, in.Content, in.DurationMin, in.Position).
		Scan(&l.ID, &l.ModuleID, &l.Title, &l.Kind, &l.Summary, &l.Content, &l.DurationMin,
			&l.Position, &l.CreatedAt, &l.UpdatedAt)
	if err != nil {
		return nil, err
	}
	return &l, nil
}

func (r *CourseRepo) UpdateLesson(ctx context.Context, id string, in LessonInput) (*domain.Lesson, error) {
	if len(in.Content) == 0 {
		in.Content = json.RawMessage(`{}`)
	}
	var l domain.Lesson
	err := r.db.QueryRow(ctx, `
		UPDATE lessons
		   SET title = $2, kind = $3, summary = $4, content = $5, duration_min = $6,
		       position = $7, updated_at = now()
		 WHERE id = $1
		RETURNING id, module_id, title, kind, summary, content, duration_min, position, created_at, updated_at`,
		id, in.Title, in.Kind, in.Summary, in.Content, in.DurationMin, in.Position).
		Scan(&l.ID, &l.ModuleID, &l.Title, &l.Kind, &l.Summary, &l.Content, &l.DurationMin,
			&l.Position, &l.CreatedAt, &l.UpdatedAt)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, ErrNotFound
	}
	if err != nil {
		return nil, err
	}
	return &l, nil
}

func (r *CourseRepo) GetLesson(ctx context.Context, id string) (*domain.Lesson, error) {
	var l domain.Lesson
	err := r.db.QueryRow(ctx, `
		SELECT id, module_id, title, kind, summary, content, duration_min, position, created_at, updated_at
		  FROM lessons WHERE id = $1`, id).
		Scan(&l.ID, &l.ModuleID, &l.Title, &l.Kind, &l.Summary, &l.Content, &l.DurationMin,
			&l.Position, &l.CreatedAt, &l.UpdatedAt)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, ErrNotFound
	}
	if err != nil {
		return nil, err
	}
	return &l, nil
}

func (r *CourseRepo) DeleteLesson(ctx context.Context, id string) error {
	tag, err := r.db.Exec(ctx, `DELETE FROM lessons WHERE id = $1`, id)
	if err != nil {
		return err
	}
	if tag.RowsAffected() == 0 {
		return ErrNotFound
	}
	return nil
}

// --- Обновление курса на месте (без потери прогресса и доступов) ---

type LessonSync struct {
	ID          string
	Title       string
	Kind        string
	Summary     string
	Content     json.RawMessage
	DurationMin int
}

type ModuleSync struct {
	ID      string
	Title   string
	Summary string
	Lessons []LessonSync
}

type SyncResult struct {
	ModulesAdded, ModulesUpdated, ModulesRemoved int
	LessonsAdded, LessonsUpdated, LessonsRemoved int
}

func norm(s string) string { return strings.ToLower(strings.TrimSpace(s)) }

// SyncContent обновляет модули и уроки курса НА МЕСТЕ: совпавшие по id (а если id
// нет — по названию) правятся с сохранением своего id, поэтому прогресс студентов
// (lesson_progress) и открытые главы (module_access) не теряются. Новые элементы
// добавляются, исчезнувшие — удаляются. Всё в одной транзакции.
func (r *CourseRepo) SyncContent(ctx context.Context, courseID string, modules []ModuleSync) (SyncResult, error) {
	var res SyncResult
	tx, err := r.db.Begin(ctx)
	if err != nil {
		return res, err
	}
	defer tx.Rollback(ctx)

	// Существующие модули курса.
	exMods, err := scanIDTitle(ctx, tx, `SELECT id, title FROM modules WHERE course_id = $1`, courseID)
	if err != nil {
		return res, err
	}
	modByID := map[string]bool{}
	modByTitle := map[string]string{}
	for _, m := range exMods {
		modByID[m.id] = true
		if _, ok := modByTitle[norm(m.title)]; !ok {
			modByTitle[norm(m.title)] = m.id
		}
	}
	usedMods := map[string]bool{}

	for mi, m := range modules {
		pos := mi + 1
		modID := ""
		if m.ID != "" && modByID[m.ID] && !usedMods[m.ID] {
			modID = m.ID
		} else if id, ok := modByTitle[norm(m.Title)]; ok && !usedMods[id] {
			modID = id
		}
		if modID != "" {
			if _, err := tx.Exec(ctx,
				`UPDATE modules SET title=$2, summary=$3, position=$4, updated_at=now() WHERE id=$1`,
				modID, m.Title, m.Summary, pos); err != nil {
				return res, err
			}
			res.ModulesUpdated++
		} else {
			if err := tx.QueryRow(ctx,
				`INSERT INTO modules (course_id, title, summary, position) VALUES ($1,$2,$3,$4) RETURNING id`,
				courseID, m.Title, m.Summary, pos).Scan(&modID); err != nil {
				return res, err
			}
			res.ModulesAdded++
		}
		usedMods[modID] = true

		// Уроки этого модуля.
		exLessons, err := scanIDTitle(ctx, tx, `SELECT id, title FROM lessons WHERE module_id = $1`, modID)
		if err != nil {
			return res, err
		}
		lesByID := map[string]bool{}
		lesByTitle := map[string]string{}
		for _, l := range exLessons {
			lesByID[l.id] = true
			if _, ok := lesByTitle[norm(l.title)]; !ok {
				lesByTitle[norm(l.title)] = l.id
			}
		}
		usedLes := map[string]bool{}

		for li, l := range m.Lessons {
			lpos := li + 1
			content := l.Content
			if len(content) == 0 {
				content = json.RawMessage(`{}`)
			}
			lesID := ""
			if l.ID != "" && lesByID[l.ID] && !usedLes[l.ID] {
				lesID = l.ID
			} else if id, ok := lesByTitle[norm(l.Title)]; ok && !usedLes[id] {
				lesID = id
			}
			if lesID != "" {
				if _, err := tx.Exec(ctx, `
					UPDATE lessons SET title=$2, kind=$3, summary=$4, content=$5,
					       duration_min=$6, position=$7, updated_at=now() WHERE id=$1`,
					lesID, l.Title, l.Kind, l.Summary, content, l.DurationMin, lpos); err != nil {
					return res, err
				}
				res.LessonsUpdated++
			} else {
				if _, err := tx.Exec(ctx, `
					INSERT INTO lessons (module_id, title, kind, summary, content, duration_min, position)
					VALUES ($1,$2,$3,$4,$5,$6,$7)`,
					modID, l.Title, l.Kind, l.Summary, content, l.DurationMin, lpos); err != nil {
					return res, err
				}
				res.LessonsAdded++
			}
			if lesID != "" {
				usedLes[lesID] = true
			}
		}
		for _, l := range exLessons {
			if !usedLes[l.id] {
				if _, err := tx.Exec(ctx, `DELETE FROM lessons WHERE id=$1`, l.id); err != nil {
					return res, err
				}
				res.LessonsRemoved++
			}
		}
	}
	for _, m := range exMods {
		if !usedMods[m.id] {
			if _, err := tx.Exec(ctx, `DELETE FROM modules WHERE id=$1`, m.id); err != nil {
				return res, err
			}
			res.ModulesRemoved++
		}
	}

	if err := tx.Commit(ctx); err != nil {
		return res, err
	}
	return res, nil
}

// scanIDTitle — вспомогательный: собирает пары (id, title).
func scanIDTitle(ctx context.Context, tx pgx.Tx, sql, arg string) ([]struct{ id, title string }, error) {
	rows, err := tx.Query(ctx, sql, arg)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	out := make([]struct{ id, title string }, 0, 16)
	for rows.Next() {
		var r struct{ id, title string }
		if err := rows.Scan(&r.id, &r.title); err != nil {
			return nil, err
		}
		out = append(out, r)
	}
	return out, rows.Err()
}

// --- Записи на курс ---

// Enroll записывает студента на курс. dueDate необязателен — срок прохождения.
func (r *CourseRepo) Enroll(ctx context.Context, userID, courseID, assignedBy string, dueDate *time.Time) error {
	var by *string
	if assignedBy != "" {
		by = &assignedBy
	}
	_, err := r.db.Exec(ctx, `
		INSERT INTO enrollments (user_id, course_id, assigned_by, started_at, due_date)
		VALUES ($1, $2, $3, now(), $4)
		ON CONFLICT (user_id, course_id) DO UPDATE
		   SET due_date = COALESCE(EXCLUDED.due_date, enrollments.due_date)`,
		userID, courseID, by, dueDate)
	return err
}

// SetDueDate меняет срок прохождения (nil — снять срок).
func (r *CourseRepo) SetDueDate(ctx context.Context, userID, courseID string, dueDate *time.Time) error {
	tag, err := r.db.Exec(ctx, `
		UPDATE enrollments SET due_date = $3, deadline_notified_at = NULL
		 WHERE user_id = $1 AND course_id = $2`, userID, courseID, dueDate)
	if err != nil {
		return err
	}
	if tag.RowsAffected() == 0 {
		return ErrNotFound
	}
	return nil
}

func (r *CourseRepo) Unenroll(ctx context.Context, userID, courseID string) error {
	_, err := r.db.Exec(ctx,
		`DELETE FROM enrollments WHERE user_id = $1 AND course_id = $2`, userID, courseID)
	return err
}

func (r *CourseRepo) IsEnrolled(ctx context.Context, userID, courseID string) (bool, error) {
	var exists bool
	err := r.db.QueryRow(ctx,
		`SELECT EXISTS (SELECT 1 FROM enrollments WHERE user_id = $1 AND course_id = $2)`,
		userID, courseID).Scan(&exists)
	return exists, err
}

func (r *CourseRepo) EnrollmentsForUser(ctx context.Context, userID string) ([]domain.Enrollment, error) {
	rows, err := r.db.Query(ctx, `
		SELECT e.id, e.user_id, e.course_id, e.status, e.due_date, e.started_at, e.completed_at,
		       e.created_at, c.title, c.slug
		  FROM enrollments e
		  JOIN courses c ON c.id = e.course_id
		 WHERE e.user_id = $1
		 ORDER BY e.created_at DESC`, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	out := make([]domain.Enrollment, 0, 4)
	for rows.Next() {
		var e domain.Enrollment
		if err := rows.Scan(&e.ID, &e.UserID, &e.CourseID, &e.Status, &e.DueDate, &e.StartedAt,
			&e.CompletedAt, &e.CreatedAt, &e.CourseTitle, &e.CourseSlug); err != nil {
			return nil, err
		}
		out = append(out, e)
	}
	return out, rows.Err()
}

// --- Пошаговый доступ к главам (модулям) ---

// GrantModuleAccess открывает студенту главу.
func (r *CourseRepo) GrantModuleAccess(ctx context.Context, userID, moduleID, grantedBy string) error {
	var by *string
	if grantedBy != "" {
		by = &grantedBy
	}
	_, err := r.db.Exec(ctx, `
		INSERT INTO module_access (user_id, module_id, granted_by)
		VALUES ($1, $2, $3)
		ON CONFLICT (user_id, module_id) DO NOTHING`, userID, moduleID, by)
	return err
}

// RevokeModuleAccess закрывает студенту главу.
func (r *CourseRepo) RevokeModuleAccess(ctx context.Context, userID, moduleID string) error {
	_, err := r.db.Exec(ctx,
		`DELETE FROM module_access WHERE user_id = $1 AND module_id = $2`, userID, moduleID)
	return err
}

// HasModuleAccess — открыта ли студенту глава.
func (r *CourseRepo) HasModuleAccess(ctx context.Context, userID, moduleID string) (bool, error) {
	var exists bool
	err := r.db.QueryRow(ctx,
		`SELECT EXISTS (SELECT 1 FROM module_access WHERE user_id = $1 AND module_id = $2)`,
		userID, moduleID).Scan(&exists)
	return exists, err
}

// ModuleAccessMap — какие главы курса открыты студенту (id главы -> true).
func (r *CourseRepo) ModuleAccessMap(ctx context.Context, userID, courseID string) (map[string]bool, error) {
	rows, err := r.db.Query(ctx, `
		SELECT ma.module_id
		  FROM module_access ma
		  JOIN modules m ON m.id = ma.module_id
		 WHERE ma.user_id = $1 AND m.course_id = $2`, userID, courseID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	out := make(map[string]bool, 16)
	for rows.Next() {
		var id string
		if err := rows.Scan(&id); err != nil {
			return nil, err
		}
		out[id] = true
	}
	return out, rows.Err()
}

// GrantFirstModuleAccess открывает самую первую главу курса — вызывается при записи,
// чтобы студенту было с чего начать.
func (r *CourseRepo) GrantFirstModuleAccess(ctx context.Context, userID, courseID, grantedBy string) error {
	var by *string
	if grantedBy != "" {
		by = &grantedBy
	}
	_, err := r.db.Exec(ctx, `
		INSERT INTO module_access (user_id, module_id, granted_by)
		SELECT $1, m.id, $3
		  FROM modules m
		 WHERE m.course_id = $2
		 ORDER BY m.position, m.created_at
		 LIMIT 1
		ON CONFLICT (user_id, module_id) DO NOTHING`, userID, courseID, by)
	return err
}

// HasTerminalCourse — записан ли студент на курс, где есть уроки-тренажёры
// (терминал). По этому признаку показываем песочницу-терминал.
func (r *CourseRepo) HasTerminalCourse(ctx context.Context, userID string) (bool, error) {
	var exists bool
	err := r.db.QueryRow(ctx, `
		SELECT EXISTS (
			SELECT 1
			  FROM enrollments e
			  JOIN modules m ON m.course_id = e.course_id
			  JOIN lessons l ON l.module_id = m.id
			 WHERE e.user_id = $1 AND l.kind = 'terminal'
		)`, userID).Scan(&exists)
	return exists, err
}

// ModuleCourseID возвращает id курса, которому принадлежит глава.
func (r *CourseRepo) ModuleCourseID(ctx context.Context, moduleID string) (string, error) {
	var courseID string
	err := r.db.QueryRow(ctx,
		`SELECT course_id FROM modules WHERE id = $1`, moduleID).Scan(&courseID)
	if errors.Is(err, pgx.ErrNoRows) {
		return "", ErrNotFound
	}
	return courseID, err
}
