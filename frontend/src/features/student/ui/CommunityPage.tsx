import { useGetCommunityQuery } from "@/shared/api/meApi";
import { lastSeenLabel } from "@/shared/lib/time";
import type { CurrentActivity, LeaderboardEntry } from "@/shared/types";
import { Badge, Card, EmptyState, PageHeader, Progress, Spinner, StatCard } from "@/shared/ui";
import { Award, BookOpen, GraduationCap, Trophy, Users, Wifi } from "lucide-react";

// Достижения выводим на фронте из безопасных полей рейтинга — так их видят все студенты.
type Achievement = {
  key: string;
  emoji: string;
  label: string;
  desc: string;
  earned: (e: LeaderboardEntry) => boolean;
};

const ACHIEVEMENTS: Achievement[] = [
  { key: "starter", emoji: "📚", label: "Первые шаги", desc: "Пройден первый урок", earned: (e) => e.lessonsCompleted >= 1 },
  { key: "active", emoji: "🔥", label: "Активность", desc: "7+ дней на платформе", earned: (e) => e.daysVisited >= 7 },
  { key: "half", emoji: "🏁", label: "Половина пути", desc: "Пройдено ≥ 50% курса", earned: (e) => e.progress >= 50 },
  { key: "excellent", emoji: "⭐", label: "Отличник", desc: "Средний балл в квизах ≥ 90%", earned: (e) => e.quizzesPassed >= 3 && e.avgQuizScore >= 90 },
  { key: "finisher", emoji: "🏆", label: "Курс пройден", desc: "Пройдено 100% курса", earned: (e) => e.progress >= 100 },
  { key: "certified", emoji: "🎓", label: "Сертификат", desc: "Получен сертификат", earned: (e) => e.certificates >= 1 },
];

const RANK_BADGE = ["🥇", "🥈", "🥉"];

function initialsOf(name: string): string {
  return (name || "?")
    .split(" ")
    .filter(Boolean)
    .slice(0, 2)
    .map((p) => p[0]?.toUpperCase())
    .join("");
}

// Короткая подпись «чем занят»: онлайн и in_progress → «проходит», иначе последний урок.
function activityText(cur: CurrentActivity | null, online: boolean): string {
  if (!cur || !cur.lesson) return online ? "на платформе" : "";
  const verb = online && cur.status === "in_progress" ? "проходит" : "последний урок:";
  return `${verb} ${cur.lesson}`;
}

type Leader = { key: string; emoji: string; title: string; entry: LeaderboardEntry; detail: string };

// Лучшие в разных номинациях — считаем из безопасных полей рейтинга.
function computeLeaders(entries: LeaderboardEntry[]): Leader[] {
  const out: Leader[] = [];
  const best = (
    key: string,
    emoji: string,
    title: string,
    value: (e: LeaderboardEntry) => number,
    detail: (e: LeaderboardEntry) => string,
    min = 1,
  ) => {
    let top: LeaderboardEntry | null = null;
    for (const e of entries) {
      if (value(e) >= min && (!top || value(e) > value(top))) top = e;
    }
    if (top) out.push({ key, emoji, title, entry: top, detail: detail(top) });
  };

  best("top", "🏆", "Лидер рейтинга", (e) => e.lessonsCompleted, (e) => `${e.lessonsCompleted} уроков пройдено`);
  best("active", "🔥", "Самый активный", (e) => e.minutesSpent, (e) => `${Math.floor(e.minutesSpent / 60)} ч ${e.minutesSpent % 60} мин · ${e.daysVisited} дн.`);
  best(
    "scorer",
    "⭐",
    "Лучший в квизах",
    (e) => (e.quizzesPassed >= 3 ? e.avgQuizScore : 0),
    (e) => `средний балл ${Math.round(e.avgQuizScore)}%`,
  );
  best(
    "closest",
    "🚀",
    "Ближе всех к финишу",
    (e) => (e.progress < 100 ? e.progress : 0),
    (e) => `${Math.round(e.progress)}% курса`,
  );
  return out;
}

export default function CommunityPage() {
  // Обновляем раз в минуту — чтобы «онлайн» и «чем занят» были свежими.
  const { data, isLoading } = useGetCommunityQuery(undefined, { pollingInterval: 60_000 });

  if (isLoading || !data) {
    return (
      <div className="grid place-items-center py-20 text-accent">
        <Spinner size={32} />
      </div>
    );
  }

  const { overview, entries, me } = data;
  const onlineNow = entries.filter((e) => e.online);
  const leaders = computeLeaders(entries);

  return (
    <>
      <PageHeader
        title="Достижения"
        subtitle="Кто чем занят, рейтинг и общий прогресс платформы"
      />

      <div className="grid gap-[var(--gap)] sm:grid-cols-2 xl:grid-cols-5">
        <StatCard label="Студентов" value={overview.students} icon={<Users size={20} />} />
        <StatCard
          label="Сейчас онлайн"
          value={overview.onlineNow}
          hint={`${overview.activeWeek} за неделю`}
          icon={<Wifi size={20} />}
        />
        <StatCard
          label="Уроков пройдено"
          value={overview.lessonsCompleted}
          hint="всеми студентами"
          icon={<BookOpen size={20} />}
        />
        <StatCard
          label="Сертификатов"
          value={overview.certificates}
          hint="выдано платформой"
          icon={<GraduationCap size={20} />}
        />
        <StatCard
          label="Достижений"
          value={ACHIEVEMENTS.length}
          hint="можно получить"
          icon={<Award size={20} />}
        />
      </div>

      {/* Лидеры в номинациях — кто лучший, самый активный и т.д. */}
      {leaders.length > 0 && (
        <div className="mt-[var(--gap)] grid gap-[var(--gap)] sm:grid-cols-2 xl:grid-cols-4">
          {leaders.map((l) => {
            const isMe = l.entry.userId === me;
            return (
              <Card
                key={l.key}
                className={`flex items-center gap-3 p-4 ${isMe ? "ring-1 ring-accent-border" : ""}`}
              >
                <span className="grid h-12 w-12 shrink-0 place-items-center rounded-full bg-surface-2 text-2xl">
                  {l.emoji}
                </span>
                <div className="min-w-0">
                  <p className="text-xs font-semibold uppercase tracking-wide text-faint">{l.title}</p>
                  <p className="flex items-center gap-1.5 truncate text-sm font-bold text-fg">
                    {l.entry.fullName || "Без имени"}
                    {isMe && <Badge tone="accent">Вы</Badge>}
                  </p>
                  <p className="truncate text-xs text-muted">{l.detail}</p>
                </div>
              </Card>
            );
          })}
        </div>
      )}

      {/* Живая активность: кто сейчас на платформе и чем занят. */}
      <Card className="mt-[var(--gap)] p-[var(--pad)]">
        <div className="mb-4 flex items-center gap-2">
          <span className="h-2.5 w-2.5 rounded-full bg-success" />
          <h2 className="text-base font-bold text-fg">Сейчас на платформе</h2>
          <span className="text-sm text-faint">· {onlineNow.length}</span>
        </div>
        {onlineNow.length === 0 ? (
          <p className="py-2 text-sm text-muted">Сейчас никого нет онлайн — загляните позже.</p>
        ) : (
          <div className="grid gap-3 sm:grid-cols-2 lg:grid-cols-3">
            {onlineNow.map((e) => (
              <div key={e.userId} className="card-flat flex items-center gap-3 p-3">
                <span
                  className="relative grid h-10 w-10 shrink-0 place-items-center rounded-full text-xs font-bold text-accent-fg"
                  style={{ background: "var(--gradient)" }}
                >
                  {initialsOf(e.fullName)}
                  <span className="absolute -bottom-0.5 -right-0.5 h-3 w-3 rounded-full border-2 border-surface-solid bg-success" />
                </span>
                <div className="min-w-0">
                  <p className="flex items-center gap-1.5 truncate text-sm font-semibold text-fg">
                    {e.fullName || "Без имени"}
                    {e.userId === me && <Badge tone="accent">Вы</Badge>}
                  </p>
                  <p className="truncate text-xs text-muted">
                    {e.current?.lesson ? (
                      <>
                        {e.current.status === "in_progress" ? "проходит " : "последний урок: "}
                        <span className="text-fg">{e.current.lesson}</span>
                        {e.current.course && <span className="text-faint"> · {e.current.course}</span>}
                      </>
                    ) : (
                      "на платформе"
                    )}
                  </p>
                </div>
              </div>
            ))}
          </div>
        )}
      </Card>

      {/* Легенда: какие бывают достижения и за что. */}
      <Card className="mt-[var(--gap)] p-[var(--pad)]">
        <h2 className="mb-4 text-base font-bold text-fg">Значки достижений</h2>
        <div className="grid gap-3 sm:grid-cols-2 lg:grid-cols-3">
          {ACHIEVEMENTS.map((a) => (
            <div key={a.key} className="card-flat flex items-start gap-3 p-3">
              <span className="text-2xl leading-none" aria-hidden>
                {a.emoji}
              </span>
              <div className="min-w-0">
                <p className="text-sm font-bold text-fg">{a.label}</p>
                <p className="text-xs text-muted">{a.desc}</p>
              </div>
            </div>
          ))}
        </div>
      </Card>

      {/* Рейтинг студентов. */}
      <Card className="mt-[var(--gap)] overflow-hidden">
        <div className="flex items-center gap-2 border-b border-line px-[var(--pad)] py-3">
          <Trophy size={18} className="text-accent" />
          <h2 className="text-base font-bold text-fg">Рейтинг по пройденным урокам</h2>
        </div>

        {entries.length === 0 ? (
          <EmptyState
            title="Пока пусто"
            description="Как только студенты начнут заниматься, здесь появится рейтинг"
            icon={<Trophy size={32} />}
          />
        ) : (
          <ul className="divide-y divide-line">
            {entries.map((e, i) => {
              const isMe = e.userId === me;
              const earned = ACHIEVEMENTS.filter((a) => a.earned(e));
              return (
                <li
                  key={e.userId}
                  className={`flex flex-wrap items-center gap-x-4 gap-y-3 px-[var(--pad)] py-4 ${
                    isMe ? "bg-accent-soft" : ""
                  }`}
                >
                  {/* Место */}
                  <span className="grid w-8 shrink-0 place-items-center text-lg font-bold text-muted">
                    {i < 3 ? RANK_BADGE[i] : i + 1}
                  </span>

                  {/* Аватар */}
                  <span
                    className="relative grid h-10 w-10 shrink-0 place-items-center rounded-full text-xs font-bold text-accent-fg"
                    style={{ background: "var(--gradient)" }}
                  >
                    {initialsOf(e.fullName)}
                    {e.online && (
                      <span className="absolute -bottom-0.5 -right-0.5 h-3 w-3 rounded-full border-2 border-surface-solid bg-success" />
                    )}
                  </span>

                  {/* Имя, активность и прогресс */}
                  <div className="min-w-[12rem] flex-1">
                    <div className="mb-0.5 flex items-center gap-2">
                      <span className="truncate text-sm font-semibold text-fg">
                        {e.fullName || "Без имени"}
                      </span>
                      {isMe && <Badge tone="accent">Вы</Badge>}
                    </div>
                    <p className="mb-1.5 truncate text-xs text-muted" title={e.current?.course}>
                      <span className={e.online ? "text-success" : "text-faint"}>
                        {lastSeenLabel(e.lastSeenAt, e.online)}
                      </span>
                      {activityText(e.current, e.online) && (
                        <span className="text-faint"> · {activityText(e.current, e.online)}</span>
                      )}
                    </p>
                    <div className="flex items-center gap-2">
                      <Progress value={e.progress} />
                      <span className="shrink-0 text-xs font-medium text-muted">
                        {e.lessonsCompleted}/{e.lessonsTotal}
                      </span>
                    </div>
                  </div>

                  {/* Достижения */}
                  <div className="flex min-w-0 flex-wrap items-center gap-1.5">
                    {earned.length === 0 ? (
                      <span className="text-xs text-faint">пока нет значков</span>
                    ) : (
                      earned.map((a) => (
                        <span
                          key={a.key}
                          title={`${a.label} — ${a.desc}`}
                          className="grid h-7 w-7 place-items-center rounded-full bg-surface-2 text-base"
                        >
                          {a.emoji}
                        </span>
                      ))
                    )}
                  </div>
                </li>
              );
            })}
          </ul>
        )}
      </Card>
    </>
  );
}
