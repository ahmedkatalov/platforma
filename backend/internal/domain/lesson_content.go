package domain

import (
	"encoding/json"
	"fmt"
	"hash/fnv"
	"regexp"
	"sort"
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

// QuizItem — шаг для вопроса-упорядочивания. В seed items идут в ПРАВИЛЬНОМ
// порядке; студенту они отдаются перемешанными.
type QuizItem struct {
	ID   string `json:"id"`
	Text string `json:"text"`
}

// QuizPair — пара для вопроса-сопоставления: левой части соответствует правая.
// В seed пары заданы в ПРАВИЛЬНОМ соответствии; студенту правые части перемешиваются.
type QuizPair struct {
	ID    string `json:"id"`
	Left  string `json:"left"`
	Right string `json:"right"`
}

// Типы вопросов квиза.
const (
	QTypeChoice = "choice" // один или несколько вариантов (по умолчанию)
	QTypeOrder  = "order"  // расставить шаги по порядку
	QTypeBlank  = "blank"  // вписать ответ (короткий текст/команда)
	QTypeMatch  = "match"  // сопоставить левые части с правыми
)

type QuizQuestion struct {
	ID          string `json:"id"`
	Text        string `json:"text"`
	Hint        string `json:"hint,omitempty"`
	Explanation string `json:"explanation,omitempty"`
	Multiple    bool   `json:"multiple,omitempty"`
	// Review отмечает вопрос на повторение ранее пройденной темы.
	Review bool `json:"review,omitempty"`

	// Type: "choice" (по умолчанию), "order" или "blank".
	Type string `json:"type,omitempty"`

	Options []QuizOption `json:"options,omitempty"` // для choice
	Items   []QuizItem   `json:"items,omitempty"`   // для order (в правильном порядке)
	Accept  []string     `json:"accept,omitempty"`  // для blank (допустимые ответы)
	Pairs   []QuizPair   `json:"pairs,omitempty"`   // для match (в правильном соответствии)

	// Lefts и Rights заполняются только при отдаче студенту (match):
	// левые части в исходном порядке, правые — перемешанные.
	Lefts  []QuizItem `json:"lefts,omitempty"`
	Rights []QuizItem `json:"rights,omitempty"`
}

// Kind возвращает тип вопроса с учётом значения по умолчанию.
func (q QuizQuestion) Kind() string {
	switch q.Type {
	case QTypeOrder:
		return QTypeOrder
	case QTypeBlank:
		return QTypeBlank
	case QTypeMatch:
		return QTypeMatch
	default:
		return QTypeChoice
	}
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
	Kind        string   `json:"kind"`
	Correct     bool     `json:"correct"`
	CorrectIDs  []string `json:"correctOptionIds,omitempty"`
	ChosenIDs   []string `json:"chosenOptionIds,omitempty"`
	CorrectText string   `json:"correctText,omitempty"` // для blank и order — эталон текстом
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
	OptionIDs    []string `json:"optionIds,omitempty"` // choice
	Order        []string `json:"order,omitempty"`     // order: id шагов в порядке студента
	Text         string   `json:"text,omitempty"`      // blank: введённый текст
	SecondsSpent int      `json:"secondsSpent"`
}

// normalizeAnswer приводит короткий ответ к сравнимому виду: без регистра,
// лишних пробелов и хвостовых точек.
func normalizeAnswer(s string) string {
	s = strings.ToLower(strings.TrimSpace(s))
	s = strings.Trim(s, ".;")
	return strings.Join(strings.Fields(s), " ")
}

// GradeQuiz сверяет ответы с эталоном. Поддерживает три типа вопросов:
// choice (варианты), order (порядок шагов) и blank (вписать ответ).
func GradeQuiz(quiz *Quiz, answers []QuizAnswer) QuizResult {
	byQuestion := make(map[string]QuizAnswer, len(answers))
	for _, answer := range answers {
		byQuestion[answer.QuestionID] = answer
	}

	result := QuizResult{
		TotalCount: len(quiz.Questions),
		PassScore:  quiz.PassScore,
		Questions:  make([]QuestionResult, 0, len(quiz.Questions)),
	}

	for _, question := range quiz.Questions {
		answer := byQuestion[question.ID]
		res := QuestionResult{
			QuestionID:  question.ID,
			Kind:        question.Kind(),
			Explanation: question.Explanation,
		}

		switch question.Kind() {
		case QTypeOrder:
			// Правильный порядок — как items заданы в seed.
			correct := make([]string, len(question.Items))
			for i, item := range question.Items {
				correct[i] = item.ID
			}
			res.CorrectText = orderText(question.Items, correct)
			res.Correct = len(answer.Order) == len(correct)
			if res.Correct {
				for i := range correct {
					if answer.Order[i] != correct[i] {
						res.Correct = false
						break
					}
				}
			}

		case QTypeBlank:
			given := normalizeAnswer(answer.Text)
			res.CorrectText = strings.Join(question.Accept, " / ")
			for _, acc := range question.Accept {
				if normalizeAnswer(acc) == given && given != "" {
					res.Correct = true
					break
				}
			}

		case QTypeMatch:
			// Правильно, если для каждой левой части выбрана её правая.
			// answer.Order — id правых частей в порядке, выровненном по левым.
			correct := make([]string, len(question.Pairs))
			for i, p := range question.Pairs {
				correct[i] = "R_" + p.ID
			}
			res.CorrectText = matchText(question.Pairs)
			res.Correct = len(answer.Order) == len(correct)
			if res.Correct {
				for i := range correct {
					if answer.Order[i] != correct[i] {
						res.Correct = false
						break
					}
				}
			}

		default: // choice
			picked := make(map[string]bool, len(answer.OptionIDs))
			for _, id := range answer.OptionIDs {
				picked[id] = true
			}
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
			res.CorrectIDs = correctIDs
			res.ChosenIDs = chosenIDs
			res.Correct = ok && len(chosenIDs) > 0
		}

		if res.Correct {
			result.CorrectCount++
		}
		result.Questions = append(result.Questions, res)
	}

	if result.TotalCount > 0 {
		result.Score = float64(result.CorrectCount) / float64(result.TotalCount) * 100
	}
	result.Passed = result.Score >= quiz.PassScore

	return result
}

// shuffleItems детерминированно перемешивает шаги на месте: порядок зависит от
// id вопроса и шага, поэтому он стабилен между загрузками, но не совпадает с
// правильным (который задан в seed).
func shuffleItems(questionID string, items []QuizItem) {
	sort.SliceStable(items, func(i, j int) bool {
		return itemSeed(questionID, items[i].ID) < itemSeed(questionID, items[j].ID)
	})
}

func itemSeed(questionID, itemID string) uint32 {
	h := fnv.New32a()
	_, _ = h.Write([]byte(questionID + ":" + itemID))
	return h.Sum32()
}

// orderText собирает читаемую подпись правильного порядка (для показа после ответа).
func orderText(items []QuizItem, order []string) string {
	byID := make(map[string]string, len(items))
	for _, item := range items {
		byID[item.ID] = item.Text
	}
	parts := make([]string, 0, len(order))
	for i, id := range order {
		parts = append(parts, fmt.Sprintf("%d. %s", i+1, byID[id]))
	}
	return strings.Join(parts, "  →  ")
}

// matchText собирает читаемое правильное соответствие «левое → правое».
func matchText(pairs []QuizPair) string {
	parts := make([]string, 0, len(pairs))
	for _, p := range pairs {
		parts = append(parts, fmt.Sprintf("%s → %s", p.Left, p.Right))
	}
	return strings.Join(parts, "; ")
}

// --- Терминал ---

type TerminalTask struct {
	ID       string   `json:"id"`
	Prompt   string   `json:"prompt"`
	Expected []string `json:"expected"`          // допустимые команды
	Pattern  string   `json:"pattern,omitempty"` // либо регулярное выражение
	Hint     string   `json:"hint,omitempty"`
	Hints    []string `json:"hints,omitempty"` // прогрессивные подсказки: концепт → синтаксис → команда
	Predict  string   `json:"predict,omitempty"` // вопрос «предскажи вывод» перед вводом команды
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
		Hints    []string        `json:"hints"`
		Predict  string          `json:"predict"`
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
	t.Hints = raw.Hints
	t.Predict = raw.Predict
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
	Intro     string          `json:"intro,omitempty"`
	Shell     string          `json:"shell,omitempty"` // подпись приглашения, по умолчанию student@devops
	Files     json.RawMessage `json:"files,omitempty"` // стартовая файловая система
	Tasks     []TerminalTask  `json:"tasks"`
	Challenge string          `json:"challenge,omitempty"` // «измени одну вещь» — задание на модификацию
	Debug     string          `json:"debug,omitempty"`     // «если сломалось» — типовая ошибка и как её чинить
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
			q := &quiz.Questions[qi]
			q.Explanation = ""
			// choice: убираем флаги правильности.
			for oi := range q.Options {
				q.Options[oi].Correct = false
			}
			// order: перемешиваем шаги, чтобы правильный порядок не был подсказкой.
			if q.Kind() == QTypeOrder {
				shuffleItems(q.ID, q.Items)
			}
			// blank: скрываем список допустимых ответов.
			if q.Kind() == QTypeBlank {
				q.Accept = nil
			}
			// match: отдаём левые части по порядку, правые — перемешанными,
			// а само соответствие (Pairs) убираем.
			if q.Kind() == QTypeMatch {
				lefts := make([]QuizItem, 0, len(q.Pairs))
				rights := make([]QuizItem, 0, len(q.Pairs))
				for _, p := range q.Pairs {
					lefts = append(lefts, QuizItem{ID: "L_" + p.ID, Text: p.Left})
					rights = append(rights, QuizItem{ID: "R_" + p.ID, Text: p.Right})
				}
				shuffleItems(q.ID, rights)
				q.Lefts = lefts
				q.Rights = rights
				q.Pairs = nil
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
