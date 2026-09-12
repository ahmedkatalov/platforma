package domain

import (
	"encoding/json"
	"testing"
)

// Раньше детерминированный хеш-порядок для некоторых id (например order q5,
// match q1) совпадал с правильным — вопрос выглядел «решённым заранее».
// Проверяем, что после SanitizeContent порядок НИКОГДА не совпадает с верным.
func TestSanitizeQuizNeverShowsCorrectOrder(t *testing.T) {
	raw := json.RawMessage(`{"questions":[
		{"id":"q5","type":"order","text":"t","items":[
			{"id":"s1","text":"a"},{"id":"s2","text":"b"},
			{"id":"s3","text":"c"},{"id":"s4","text":"d"}]},
		{"id":"q1","type":"match","text":"t","pairs":[
			{"id":"p1","left":"L1","right":"R1"},{"id":"p2","left":"L2","right":"R2"},
			{"id":"p3","left":"L3","right":"R3"},{"id":"p4","left":"L4","right":"R4"}]}
	]}`)

	out := SanitizeContent(LessonQuiz, raw)

	var parsed struct {
		Questions []struct {
			ID     string `json:"id"`
			Type   string `json:"type"`
			Items  []struct{ ID string } `json:"items"`
			Rights []struct{ ID string } `json:"rights"`
		} `json:"questions"`
	}
	if err := json.Unmarshal(out, &parsed); err != nil {
		t.Fatalf("не разобрать выход: %v", err)
	}

	for _, q := range parsed.Questions {
		switch q.Type {
		case "order":
			got := []string{}
			for _, it := range q.Items {
				got = append(got, it.ID)
			}
			if eq(got, []string{"s1", "s2", "s3", "s4"}) {
				t.Errorf("order %s показан в правильном порядке (решён заранее): %v", q.ID, got)
			}
			if len(got) != 4 {
				t.Errorf("order %s: ожидали 4 шага, получили %v", q.ID, got)
			}
		case "match":
			got := []string{}
			for _, r := range q.Rights {
				got = append(got, r.ID)
			}
			if eq(got, []string{"R_p1", "R_p2", "R_p3", "R_p4"}) {
				t.Errorf("match %s: правые части выровнены по левым (решён заранее): %v", q.ID, got)
			}
			if len(got) != 4 {
				t.Errorf("match %s: ожидали 4 правых части, получили %v", q.ID, got)
			}
		}
	}
}

func eq(a, b []string) bool {
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		if a[i] != b[i] {
			return false
		}
	}
	return true
}
