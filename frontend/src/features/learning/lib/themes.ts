import type { Lesson, Module } from "@/shared/types";

// Тема — это связка внутри главы: несколько страниц (теория, тренажёр, практика)
// и квиз в конце для проверки. Тему собираем из уроков модуля: всё, что идёт
// до квиза включительно, — одна тема. Так «глава → тема → страницы + квиз».
export type Theme = {
  key: string;
  title: string;
  pages: Lesson[]; // страницы темы (не квизы)
  quiz: Lesson | null; // финальная проверка темы
  lessons: Lesson[]; // все уроки темы по порядку
};

function capitalize(s: string): string {
  return s ? s[0].toUpperCase() + s.slice(1) : s;
}

function titleFor(pages: Lesson[], quiz: Lesson | null): string {
  // Квиз обычно назван по теме («Квиз: права доступа») — берём это как заголовок.
  if (quiz && quiz.title.startsWith("Квиз:")) {
    return capitalize(quiz.title.replace(/^Квиз:\s*/, "").trim());
  }
  return pages[0]?.title ?? quiz?.title ?? "Тема";
}

export function groupThemes(module: Module): Theme[] {
  const lessons = module.lessons ?? [];
  const themes: Theme[] = [];
  let bucket: Lesson[] = [];

  const flush = () => {
    if (bucket.length === 0) return;
    const quiz = bucket.find((l) => l.kind === "quiz") ?? null;
    const pages = bucket.filter((l) => l.kind !== "quiz");
    const anchor = pages[0] ?? bucket[0];
    themes.push({
      key: anchor.id,
      title: titleFor(pages, quiz),
      pages,
      quiz,
      lessons: [...bucket],
    });
    bucket = [];
  };

  for (const lesson of lessons) {
    bucket.push(lesson);
    if (lesson.kind === "quiz") flush(); // квиз закрывает тему
  }
  flush();
  return themes;
}

export type ThemeProgress = { done: number; total: number; percent: number };

export function themeProgress(
  theme: Theme,
  isDone: (lessonId: string) => boolean,
): ThemeProgress {
  const total = theme.lessons.length;
  const done = theme.lessons.filter((l) => isDone(l.id)).length;
  return { done, total, percent: total ? (done / total) * 100 : 0 };
}
