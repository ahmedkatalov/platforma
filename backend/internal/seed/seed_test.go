package seed

import "testing"

// TestDevOpsCourseLoads гарантирует, что все content/*.json парсятся в
// заготовку курса (loadModules паникует на битом JSON) и базово целостны.
func TestDevOpsCourseLoads(t *testing.T) {
	c := DevOpsCourse()
	if c.Slug == "" || c.Title == "" {
		t.Fatal("нет метаданных курса")
	}
	if len(c.Modules) == 0 {
		t.Fatal("нет модулей")
	}
	lessons := 0
	for _, m := range c.Modules {
		if m.Title == "" {
			t.Errorf("модуль без title")
		}
		if len(m.Lessons) == 0 {
			t.Errorf("модуль %q без уроков", m.Title)
		}
		for _, l := range m.Lessons {
			if l.Title == "" || l.Kind == "" {
				t.Errorf("урок без title/kind в модуле %q", m.Title)
			}
			if l.Content == nil {
				t.Errorf("урок %q без content", l.Title)
			}
			lessons++
		}
	}
	t.Logf("курс загружен: модулей=%d, уроков=%d", len(c.Modules), lessons)
}
