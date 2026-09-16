package seed

import "testing"

// TestCoursesLoad гарантирует, что все content/*.json и content-go/*.json
// парсятся в заготовки курсов (readModules паникует на битом JSON) и базово
// целостны. Заодно проверяет уникальность slug'ов.
func TestCoursesLoad(t *testing.T) {
	all := AllCourses()
	if len(all) == 0 {
		t.Fatal("нет курсов")
	}
	slugs := map[string]bool{}
	for _, c := range all {
		if c.Slug == "" || c.Title == "" {
			t.Fatalf("курс без slug/title: %+v", c.Title)
		}
		if slugs[c.Slug] {
			t.Errorf("дубликат slug: %s", c.Slug)
		}
		slugs[c.Slug] = true
		if len(c.Modules) == 0 {
			t.Errorf("курс %q без модулей", c.Slug)
		}
		lessons := 0
		for _, m := range c.Modules {
			if m.Title == "" {
				t.Errorf("[%s] модуль без title", c.Slug)
			}
			if len(m.Lessons) == 0 {
				t.Errorf("[%s] модуль %q без уроков", c.Slug, m.Title)
			}
			for _, l := range m.Lessons {
				if l.Title == "" || l.Kind == "" {
					t.Errorf("[%s] урок без title/kind в модуле %q", c.Slug, m.Title)
				}
				if l.Content == nil {
					t.Errorf("[%s] урок %q без content", c.Slug, l.Title)
				}
				lessons++
			}
		}
		t.Logf("курс %q: модулей=%d, уроков=%d", c.Slug, len(c.Modules), lessons)
	}
}
