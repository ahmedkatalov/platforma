package domain

import (
	"encoding/json"
	"regexp"
	"strings"
)

// Содержимое урока хранится в JSONB. Здесь — его разбор, очистка от «ответов»
// перед отдачей студенту и проверка присланных решений.

// --- Квиз ---

type QuizOption struct {
	ID      string `json:"id"`
	Text    string `json:"text"`
	Correct bool   `json:"correct,omitempty"`
}

type QuizQuestion struct {
	ID          string `json:"id"`
	Text        string `json:"text"`
	Hint        string `json:"hint,omitempty"`
	Explanation string `json:"explanation,omitempty"`
	Multiple    bool   `json:"multiple,omitempty"`
	// Review отмечает вопрос на повторение ранее пройденной темы.
	Review  bool         `json:"review,omitempty"`
	Options []QuizOption `json:"options"`
}

type Quiz struct {
	Intro        string         `json:"intro,omitempty"`
	PassScore    float64        `json:"passScore,omitempty"`
	TimeLimitSec int            `json:"timeLimitSec,omitempty"`
	Shuffle      bool           `json:"shuffle,omitempty"`
	Questions    []QuizQuestion `json:"questions"`
}

func ParseQuiz(raw json.RawMessage) (*Quiz, error) {
	var quiz Quiz
	if err := json.Unmarshal(raw, &quiz); err != nil {
		return nil, err
	}
	if quiz.PassScore <= 0 {
		quiz.PassScore = 70
	}
	return &quiz, nil
}

// QuestionResult — что вернуть студенту по каждому вопросу после проверки.
type QuestionResult struct {
	QuestionID  string   `json:"questionId"`
	Correct     bool     `json:"correct"`
	CorrectIDs  []string `json:"correctOptionIds"`
	ChosenIDs   []string `json:"chosenOptionIds"`
	Explanation string   `json:"explanation,omitempty"`
}

type QuizResult struct {
	Score        float64          `json:"score"`
	Passed       bool             `json:"passed"`
	CorrectCount int              `json:"correctCount"`
	TotalCount   int              `json:"totalCount"`
	PassScore    float64          `json:"passScore"`
	Questions    []QuestionResult `json:"questions"`
}

// QuizAnswer — ответ студента на один вопрос.
type QuizAnswer struct {
	QuestionID   string   `json:"questionId"`
	OptionIDs    []string `json:"optionIds"`
	SecondsSpent int      `json:"secondsSpent"`
}

// GradeQuiz сверяет ответы с эталоном. Вопрос засчитан, только если выбраны
// ровно все правильные варианты и ни одного лишнего.
func GradeQuiz(quiz *Quiz, answers []QuizAnswer) QuizResult {
	chosen := make(map[string]map[string]bool, len(answers))
	for _, answer := range answers {
		set := make(map[string]bool, len(answer.OptionIDs))
		for _, id := range answer.OptionIDs {
			set[id] = true
		}
		chosen[answer.QuestionID] = set
	}

	result := QuizResult{
		TotalCount: len(quiz.Questions),
		PassScore:  quiz.PassScore,
		Questions:  make([]QuestionResult, 0, len(quiz.Questions)),
	}

	for _, question := range quiz.Questions {
		picked := chosen[question.ID]
		correctIDs := make([]string, 0, 2)
		chosenIDs := make([]string, 0, 2)
		ok := true

		for _, option := range question.Options {
			if option.Correct {
				correctIDs = append(correctIDs, option.ID)
			}
			if picked[option.ID] {
				chosenIDs = append(chosenIDs, option.ID)
			}
			if option.Correct != picked[option.ID] {
				ok = false
			}
		}

		if ok && len(chosenIDs) > 0 {
			result.CorrectCount++
		} else {
			ok = false
		}

		result.Questions = append(result.Questions, QuestionResult{
			QuestionID:  question.ID,
			Correct:     ok,
			CorrectIDs:  correctIDs,
			ChosenIDs:   chosenIDs,
			Explanation: question.Explanation,
		})
	}

	if result.TotalCount > 0 {
		result.Score = float64(result.CorrectCount) / float64(result.TotalCount) * 100
	}
	result.Passed = result.Score >= quiz.PassScore

	return result
}

// --- Терминал ---

type TerminalTask struct {
	ID       string   `json:"id"`
	Prompt   string   `json:"prompt"`
	Expected []string `json:"expected"`          // допустимые команды
	Pattern  string   `json:"pattern,omitempty"` // либо регулярное выражение
	Hint     string   `json:"hint,omitempty"`
	Success  string   `json:"success,omitempty"`
}

// UnmarshalJSON допускает expected как строку и как массив строк.
func (t *TerminalTask) UnmarshalJSON(data []byte) error {
	type alias struct {
		ID       string          `json:"id"`
		Prompt   string          `json:"prompt"`
		Expected json.RawMessage `json:"expected"`
		Pattern  string          `json:"pattern"`
		Hint     string          `json:"hint"`
		Success  string          `json:"success"`
	}

	var raw alias
	if err := json.Unmarshal(data, &raw); err != nil {
		return err
	}

	t.ID = raw.ID
	t.Prompt = raw.Prompt
	t.Pattern = raw.Pattern
	t.Hint = raw.Hint
	t.Success = raw.Success
	t.Expected = nil

	if len(raw.Expected) == 0 {
		return nil
	}
	var single string
	if err := json.Unmarshal(raw.Expected, &single); err == nil {
		if strings.TrimSpace(single) != "" {
			t.Expected = []string{single}
		}
		return nil
	}
	return json.Unmarshal(raw.Expected, &t.Expected)
}

type TerminalLesson struct {
	Intro string          `json:"intro,omitempty"`
	Shell string          `json:"shell,omitempty"` // подпись приглашения, по умолчанию student@devops
	Files json.RawMessage `json:"files,omitempty"` // стартовая файловая система
	Tasks []TerminalTask  `json:"tasks"`
}

func ParseTerminal(raw json.RawMessage) (*TerminalLesson, error) {
	var lesson TerminalLesson
	if err := json.Unmarshal(raw, &lesson); err != nil {
		return nil, err
	}
	return &lesson, nil
}

// --- Урок с кодом ---

type CodeCheck struct {
	Type    string `json:"type"` // contains | notContains | regex
	Value   string `json:"value"`
	Message string `json:"message,omitempty"`
}

type CodeLesson struct {
	Language string      `json:"language"`
	Task     string      `json:"task"`
	Starter  string      `json:"starter"`
	Solution string      `json:"solution,omitempty"`
	Hint     string      `json:"hint,omitempty"`
	Checks   []CodeCheck `json:"checks"`
}

func ParseCode(raw json.RawMessage) (*CodeLesson, error) {
	var lesson CodeLesson
	if err := json.Unmarshal(raw, &lesson); err != nil {
		return nil, err
	}
	if lesson.Language == "" {
		lesson.Language = "bash"
	}
	return &lesson, nil
}

// SanitizeContent убирает из содержимого всё, что подсказывает ответ:
// флаги correct, пояснения, ожидаемые команды и эталонное решение.
func SanitizeContent(kind string, raw json.RawMessage) json.RawMessage {
	switch kind {
	case LessonQuiz:
		quiz, err := ParseQuiz(raw)
		if err != nil {
			return json.RawMessage(`{}`)
		}
		for qi := range quiz.Questions {
			quiz.Questions[qi].Explanation = ""
			for oi := range quiz.Questions[qi].Options {
				quiz.Questions[qi].Options[oi].Correct = false
			}
		}
		return mustJSON(quiz)

	case LessonTerminal:
		lesson, err := ParseTerminal(raw)
		if err != nil {
			return json.RawMessage(`{}`)
		}
		for i := range lesson.Tasks {
			lesson.Tasks[i].Expected = nil
			lesson.Tasks[i].Pattern = ""
		}
		return mustJSON(lesson)

	case LessonCode:
		lesson, err := ParseCode(raw)
		if err != nil {
			return json.RawMessage(`{}`)
		}
		lesson.Solution = ""
		lesson.Checks = nil
		return mustJSON(lesson)

	default:
		return raw
	}
}

func mustJSON(v any) json.RawMessage {
	data, err := json.Marshal(v)
	if err != nil {
		return json.RawMessage(`{}`)
	}
	return data
}

// RunCodeCheck проверяет решение студента по одному условию задания.
func RunCodeCheck(check CodeCheck, code string) bool {
	switch check.Type {
	case "notContains":
		return !strings.Contains(code, check.Value)
	case "regex":
		re, err := regexp.Compile(check.Value)
		return err == nil && re.MatchString(code)
	default: // contains
		return strings.Contains(code, check.Value)
	}
}

// MatchTerminalCommand проверяет, подходит ли введённая команда под задание.
func MatchTerminalCommand(task *TerminalTask, command string) bool {
	normalized := NormalizeCommand(command)
	if normalized == "" {
		return false
	}

	for _, expected := range task.Expected {
		if strings.EqualFold(NormalizeCommand(expected), normalized) {
			return true
		}
	}

	if task.Pattern != "" {
		if re, err := regexp.Compile(task.Pattern); err == nil && re.MatchString(normalized) {
			return true
		}
	}
	return false
}

// DescribeCodeCheck — текст проверки, если автор урока не задал свой.
func DescribeCodeCheck(check CodeCheck) string {
	switch check.Type {
	case "notContains":
		return "В решении не должно быть: " + check.Value
	case "regex":
		return "Решение должно соответствовать шаблону: " + check.Value
	default:
		return "Решение должно содержать: " + check.Value
	}
}

// NormalizeCommand приводит команду к виду, пригодному для сравнения:
// схлопывает пробелы и убирает хвостовые разделители.
func NormalizeCommand(cmd string) string {
	return strings.Join(strings.Fields(strings.TrimSpace(cmd)), " ")
}
