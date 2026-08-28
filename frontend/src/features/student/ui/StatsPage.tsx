import { useState } from "react";
import {
  Area,
  AreaChart,
  Bar,
  BarChart,
  CartesianGrid,
  ResponsiveContainer,
  Tooltip,
  XAxis,
  YAxis,
} from "recharts";

import { useGetMyAttemptsQuery, useGetMyStatsQuery } from "@/shared/api/meApi";
import type { Attempt } from "@/shared/types";
import { Badge, Card, EmptyState, PageHeader, Progress, Select, Spinner, StatCard } from "@/shared/ui";
import { IconChart, IconCheck, IconClock, IconFlame, IconTerminal } from "@/shared/ui/icons";

const dayFmt = new Intl.DateTimeFormat("ru-RU", { day: "2-digit", month: "short" });
const stampFmt = new Intl.DateTimeFormat("ru-RU", {
  day: "2-digit",
  month: "2-digit",
  hour: "2-digit",
  minute: "2-digit",
});

const KIND_LABEL: Record<Attempt["kind"], string> = {
  quiz: "Квиз",
  terminal: "Терминал",
  code: "Код",
};

export default function StatsPage() {
  const [days, setDays] = useState(30);
  const { data: stats, isLoading } = useGetMyStatsQuery(days);
  const { data: attempts = [] } = useGetMyAttemptsQuery(20);

  if (isLoading || !stats) {
    return (
      <div className="grid place-items-center py-20 text-accent">
        <Spinner size={32} />
      </div>
    );
  }

  const chartData = stats.activity.map((day) => ({
    day: dayFmt.format(new Date(day.day)),
    Минуты: Math.round(day.secondsSpent / 60),
    Заходов: day.visits,
  }));

  return (
    <>
      <PageHeader
        title="Моя статистика"
        subtitle="Посещаемость, время и успеваемость"
        actions={
          <Select
            value={days}
            onChange={(e) => setDays(Number(e.target.value))}
            className="w-44"
          >
            <option value={7}>Неделя</option>
            <option value={30}>30 дней</option>
            <option value={90}>90 дней</option>
            <option value={365}>Год</option>
          </Select>
        }
      />

      <div className="grid gap-[var(--gap)] sm:grid-cols-2 xl:grid-cols-4">
        <StatCard
          label="Прогресс"
          value={`${Math.round(stats.summary.progress)}%`}
          hint={`${stats.summary.lessonsCompleted} из ${stats.summary.lessonsTotal} уроков`}
          icon={<IconChart size={20} />}
        />
        <StatCard
          label="Дней подряд"
          value={stats.streak}
          hint="Текущая серия"
          icon={<IconFlame size={20} />}
        />
        <StatCard
          label="Дней посещения"
          value={stats.summary.daysVisited}
          hint="За всё время"
          icon={<IconClock size={20} />}
        />
        <StatCard
          label="Время обучения"
          value={`${Math.floor(stats.summary.minutesSpent / 60)} ч`}
          hint={`${stats.summary.minutesSpent} минут`}
          icon={<IconClock size={20} />}
        />
      </div>

      <Card className="mt-[var(--gap)] p-[var(--pad)]">
        <h2 className="mb-4 text-base font-bold text-fg">Время на платформе</h2>
        {chartData.length === 0 ? (
          <EmptyState
            title="Данных пока нет"
            description="Занимайтесь — и здесь появится ваша статистика"
            icon={<IconChart size={32} />}
          />
        ) : (
          <div className="h-64">
            <ResponsiveContainer width="100%" height="100%">
              <AreaChart data={chartData}>
                <defs>
                  <linearGradient id="minutesFill" x1="0" y1="0" x2="0" y2="1">
                    <stop offset="0%" stopColor="var(--accent)" stopOpacity={0.5} />
                    <stop offset="100%" stopColor="var(--accent)" stopOpacity={0} />
                  </linearGradient>
                </defs>
                <CartesianGrid strokeDasharray="3 3" stroke="var(--border)" vertical={false} />
                <XAxis
                  dataKey="day"
                  tick={{ fill: "var(--text-muted)", fontSize: 11 }}
                  axisLine={false}
                  tickLine={false}
                />
                <YAxis
                  tick={{ fill: "var(--text-muted)", fontSize: 11 }}
                  axisLine={false}
                  tickLine={false}
                  allowDecimals={false}
                />
                <Tooltip
                  contentStyle={{
                    background: "var(--surface-solid)",
                    border: "1px solid var(--border)",
                    borderRadius: "var(--radius-md)",
                    color: "var(--text)",
                  }}
                />
                <Area
                  type="monotone"
                  dataKey="Минуты"
                  stroke="var(--accent)"
                  strokeWidth={2}
                  fill="url(#minutesFill)"
                />
              </AreaChart>
            </ResponsiveContainer>
          </div>
        )}
      </Card>

      <div className="mt-[var(--gap)] grid gap-[var(--gap)] lg:grid-cols-2">
        <Card className="p-[var(--pad)]">
          <h2 className="mb-4 text-base font-bold text-fg">Заходы по дням</h2>
          {chartData.length === 0 ? (
            <p className="py-8 text-center text-sm text-muted">Пока пусто</p>
          ) : (
            <div className="h-56">
              <ResponsiveContainer width="100%" height="100%">
                <BarChart data={chartData}>
                  <CartesianGrid strokeDasharray="3 3" stroke="var(--border)" vertical={false} />
                  <XAxis
                    dataKey="day"
                    tick={{ fill: "var(--text-muted)", fontSize: 11 }}
                    axisLine={false}
                    tickLine={false}
                  />
                  <YAxis
                    tick={{ fill: "var(--text-muted)", fontSize: 11 }}
                    axisLine={false}
                    tickLine={false}
                    allowDecimals={false}
                  />
                  <Tooltip
                    cursor={{ fill: "var(--accent-soft)" }}
                    contentStyle={{
                      background: "var(--surface-solid)",
                      border: "1px solid var(--border)",
                      borderRadius: "var(--radius-md)",
                      color: "var(--text)",
                    }}
                  />
                  <Bar dataKey="Заходов" fill="var(--accent-2)" radius={[6, 6, 0, 0]} />
                </BarChart>
              </ResponsiveContainer>
            </div>
          )}
        </Card>

        <Card className="p-[var(--pad)]">
          <h2 className="mb-4 text-base font-bold text-fg">Успеваемость</h2>
          <div className="space-y-4">
            <div>
              <div className="mb-1.5 flex items-center justify-between text-sm">
                <span className="text-muted">Пройдено уроков</span>
                <span className="font-bold text-accent">
                  {stats.summary.lessonsCompleted} / {stats.summary.lessonsTotal}
                </span>
              </div>
              <Progress value={stats.summary.progress} />
            </div>

            <div className="card-flat p-3 text-sm text-muted">
              <p>
                Курсов назначено:{" "}
                <span className="font-bold text-fg">{stats.summary.courses}</span>
              </p>
              <p className="mt-1">
                Средняя активность:{" "}
                <span className="font-bold text-fg">
                  {stats.summary.daysVisited > 0
                    ? Math.round(stats.summary.minutesSpent / stats.summary.daysVisited)
                    : 0}{" "}
                  мин в день
                </span>
              </p>
            </div>

            <div className="card-flat p-3 text-sm text-muted">
              <p>
                Точность в квизах:{" "}
                <span className="font-bold text-fg">{Math.round(stats.quiz.accuracy)}%</span>{" "}
                ({stats.quiz.answeredCorrect} из {stats.quiz.answeredTotal} ответов)
              </p>
              <p className="mt-1">
                Средняя скорость ответа:{" "}
                <span className="font-bold text-fg">
                  {stats.quiz.avgSecondsPerQuestion.toFixed(1)} c
                </span>
                {stats.quiz.fastestSeconds > 0 && ` · лучший ${stats.quiz.fastestSeconds} c`}
              </p>
            </div>
          </div>
        </Card>
      </div>

      <div className="mt-[var(--gap)] grid gap-[var(--gap)] sm:grid-cols-2 xl:grid-cols-4">
        <StatCard
          label="Попыток в квизах"
          value={stats.quiz.attempts}
          hint={`${stats.quiz.passed} успешных`}
          icon={<IconCheck size={20} />}
        />
        <StatCard
          label="Средний балл"
          value={`${Math.round(stats.quiz.averageScore)}%`}
          hint={`лучший результат ${Math.round(stats.quiz.bestScore)}%`}
          icon={<IconChart size={20} />}
        />
        <StatCard
          label="Точность ответов"
          value={`${Math.round(stats.quiz.accuracy)}%`}
          hint={`${stats.quiz.answeredCorrect} из ${stats.quiz.answeredTotal}`}
          icon={<IconCheck size={20} />}
        />
        <StatCard
          label="Скорость ответа"
          value={`${stats.quiz.avgSecondsPerQuestion.toFixed(1)} c`}
          hint="в среднем на вопрос"
          icon={<IconClock size={20} />}
        />
      </div>

      <Card className="mt-[var(--gap)] overflow-hidden">
        <div className="border-b border-line px-[var(--pad)] py-3">
          <h2 className="text-base font-bold text-fg">История попыток</h2>
        </div>

        {attempts.length === 0 ? (
          <EmptyState
            title="Попыток пока нет"
            description="Пройдите квиз или тренажёр — результаты появятся здесь"
            icon={<IconTerminal size={32} />}
          />
        ) : (
          <div className="overflow-x-auto">
            <table className="w-full min-w-[40rem] text-sm">
              <thead>
                <tr className="border-b border-line text-left text-xs uppercase tracking-wide text-faint">
                  <th className="px-4 py-3 font-semibold">Когда</th>
                  <th className="px-4 py-3 font-semibold">Урок</th>
                  <th className="px-4 py-3 font-semibold">Тип</th>
                  <th className="px-4 py-3 font-semibold">Результат</th>
                  <th className="px-4 py-3 font-semibold">Время</th>
                </tr>
              </thead>
              <tbody>
                {attempts.map((attempt) => (
                  <tr key={attempt.id} className="border-b border-line/60 last:border-0">
                    <td className="whitespace-nowrap px-4 py-3 text-muted">
                      {stampFmt.format(new Date(attempt.createdAt))}
                    </td>
                    <td className="px-4 py-3">
                      <span className="block font-medium text-fg">{attempt.lessonTitle}</span>
                      <span className="block text-xs text-faint">{attempt.courseTitle}</span>
                    </td>
                    <td className="px-4 py-3">
                      <Badge tone={attempt.kind === "quiz" ? "accent" : "default"}>
                        {KIND_LABEL[attempt.kind]}
                      </Badge>
                    </td>
                    <td className="px-4 py-3">
                      <span
                        className={`font-bold ${attempt.passed ? "text-success" : "text-warning"}`}
                      >
                        {Math.round(attempt.score)}%
                      </span>
                      {attempt.totalCount > 0 && (
                        <span className="ml-1.5 text-xs text-faint">
                          {attempt.correctCount}/{attempt.totalCount}
                        </span>
                      )}
                    </td>
                    <td className="px-4 py-3 text-muted">
                      {Math.floor(attempt.durationSeconds / 60)}:
                      {String(attempt.durationSeconds % 60).padStart(2, "0")}
                    </td>
                  </tr>
                ))}
              </tbody>
            </table>
          </div>
        )}
      </Card>
    </>
  );
}
