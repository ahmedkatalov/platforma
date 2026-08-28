import { useMemo, useState } from "react";
import { Link } from "react-router-dom";

import { useGetMyQuizzesQuery, useGetMyStatsQuery } from "@/shared/api/meApi";
import type { QuizCard } from "@/shared/types";
import {
  Badge,
  Card,
  EmptyState,
  PageHeader,
  Progress,
  Select,
  Spinner,
  StatCard,
} from "@/shared/ui";
import { IconCheck, IconChart, IconClock, IconFlame } from "@/shared/ui/icons";

type Filter = "all" | "todo" | "repeat" | "done";

const dateFmt = new Intl.DateTimeFormat("ru-RU", {
  day: "2-digit",
  month: "2-digit",
  year: "2-digit",
});

const DAY = 86_400_000;

// Квиз стоит повторить, если он сдан давно или со слабым результатом.
function needsRepeat(quiz: QuizCard): boolean {
  if (quiz.status !== "completed") return false;
  if ((quiz.bestScore ?? 0) < 90) return true;
  if (!quiz.lastTriedAt) return false;
  return Date.now() - new Date(quiz.lastTriedAt).getTime() > 14 * DAY;
}

export default function QuizzesPage() {
  const { data: quizzes = [], isLoading } = useGetMyQuizzesQuery();
  const { data: stats } = useGetMyStatsQuery(30);
  const [filter, setFilter] = useState<Filter>("all");

  const groups = useMemo(() => {
    const filtered = quizzes.filter((quiz) => {
      switch (filter) {
        case "todo":
          return quiz.status !== "completed";
        case "repeat":
          return needsRepeat(quiz);
        case "done":
          return quiz.status === "completed";
        default:
          return true;
      }
    });

    // Группируем по модулю, порядок сохраняем как в курсе.
    const byModule = new Map<string, QuizCard[]>();
    for (const quiz of filtered) {
      const list = byModule.get(quiz.moduleTitle) ?? [];
      list.push(quiz);
      byModule.set(quiz.moduleTitle, list);
    }
    return [...byModule.entries()];
  }, [quizzes, filter]);

  if (isLoading) {
    return (
      <div className="grid place-items-center py-20 text-accent">
        <Spinner size={32} />
      </div>
    );
  }

  const done = quizzes.filter((q) => q.status === "completed").length;
  const repeat = quizzes.filter(needsRepeat).length;
  const totalQuestions = quizzes.reduce((sum, q) => sum + q.questions, 0);

  return (
    <>
      <PageHeader
        title="Квизы"
        subtitle="Проверка знаний по темам — можно проходить сколько угодно раз"
        actions={
          <Select
            value={filter}
            onChange={(e) => setFilter(e.target.value as Filter)}
            className="w-52"
          >
            <option value="all">Все квизы</option>
            <option value="todo">Не пройдены</option>
            <option value="repeat">Стоит повторить</option>
            <option value="done">Пройдены</option>
          </Select>
        }
      />

      <div className="grid gap-[var(--gap)] sm:grid-cols-2 xl:grid-cols-4">
        <StatCard
          label="Пройдено"
          value={`${done} из ${quizzes.length}`}
          hint={`${totalQuestions} вопросов всего`}
          icon={<IconCheck size={20} />}
        />
        <StatCard
          label="Средний балл"
          value={stats && stats.quiz.attempts > 0 ? `${Math.round(stats.quiz.averageScore)}%` : "—"}
          hint={stats ? `${stats.quiz.attempts} попыток` : undefined}
          icon={<IconChart size={20} />}
        />
        <StatCard
          label="Точность ответов"
          value={stats && stats.quiz.answeredTotal > 0 ? `${Math.round(stats.quiz.accuracy)}%` : "—"}
          hint={stats ? `${stats.quiz.answeredCorrect} из ${stats.quiz.answeredTotal}` : undefined}
          icon={<IconCheck size={20} />}
        />
        <StatCard
          label="Стоит повторить"
          value={repeat}
          hint="сданы давно или на слабый балл"
          icon={<IconFlame size={20} />}
        />
      </div>

      <div className="mt-[var(--gap)]">
        <Progress value={quizzes.length ? (done / quizzes.length) * 100 : 0} />
      </div>

      {groups.length === 0 ? (
        <Card className="mt-[var(--gap)]">
          <EmptyState
            title={filter === "all" ? "Квизов пока нет" : "Ничего не подходит под фильтр"}
            description={
              filter === "all"
                ? "Как только администратор откроет вам курс, квизы появятся здесь"
                : "Попробуйте выбрать другой фильтр"
            }
            icon={<IconChart size={32} />}
          />
        </Card>
      ) : (
        <div className="mt-[var(--gap)] space-y-[var(--gap)]">
          {groups.map(([moduleTitle, items]) => (
            <Card key={moduleTitle} className="p-[var(--pad)]">
              <h2 className="mb-3 text-base font-bold text-fg">{moduleTitle}</h2>

              <ul className="space-y-2">
                {items.map((quiz) => {
                  const passed = quiz.status === "completed";
                  const repeatNeeded = needsRepeat(quiz);

                  return (
                    <li key={quiz.lessonId}>
                      <Link
                        to={`/learn/courses/${quiz.courseSlug}/lessons/${quiz.lessonId}`}
                        className="card-flat flex flex-wrap items-center gap-3 p-3 transition-colors hover:bg-surface-hover"
                      >
                        <span
                          className={`grid h-8 w-8 shrink-0 place-items-center rounded-full text-xs font-bold ${
                            passed
                              ? "bg-[var(--success-soft)] text-success"
                              : "bg-surface-2 text-muted"
                          }`}
                        >
                          {passed ? <IconCheck size={16} /> : quiz.questions}
                        </span>

                        <span className="min-w-0 flex-1">
                          <span className="block truncate text-sm font-semibold text-fg">
                            {quiz.title}
                          </span>
                          <span className="block truncate text-xs text-faint">
                            {quiz.questions} вопросов · порог {Math.round(quiz.passScore)}% ·{" "}
                            {quiz.durationMin} мин
                            {quiz.lastTriedAt &&
                              ` · последняя попытка ${dateFmt.format(new Date(quiz.lastTriedAt))}`}
                          </span>
                        </span>

                        {repeatNeeded && <Badge tone="warning">повторить</Badge>}

                        {quiz.bestScore != null ? (
                          <span
                            className={`shrink-0 text-sm font-bold ${
                              quiz.bestScore >= quiz.passScore ? "text-success" : "text-warning"
                            }`}
                          >
                            {Math.round(quiz.bestScore)}%
                          </span>
                        ) : (
                          <Badge>
                            <IconClock size={12} /> не пройден
                          </Badge>
                        )}
                      </Link>
                    </li>
                  );
                })}
              </ul>
            </Card>
          ))}
        </div>
      )}

      <Card className="mt-[var(--gap)] p-[var(--pad)] text-sm text-muted">
        Квиз можно перепроходить без ограничений — засчитывается лучший результат.
        Возвращайтесь к пройденным темам раз в пару недель: так материал остаётся в голове.
      </Card>
    </>
  );
}
