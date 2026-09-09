import { useMemo, useState } from "react";

import { useGetLiveActivityQuery } from "@/features/admin/api/adminApi";
import { lastSeenLabel } from "@/shared/lib/time";
import type { AdminActivityRow, UserStatus } from "@/shared/types";
import { Badge, Card, EmptyState, PageHeader, Spinner, StatCard } from "@/shared/ui";
import { Activity, BookOpen, Clock, Users, Wifi } from "lucide-react";

const STATUS_TONE: Record<UserStatus, "success" | "warning" | "danger"> = {
  active: "success",
  invited: "warning",
  blocked: "danger",
};

const STATUS_LABEL: Record<UserStatus, string> = {
  active: "Активен",
  invited: "Приглашён",
  blocked: "Заблокирован",
};

function CurrentCell({ row }: { row: AdminActivityRow }) {
  const cur = row.current;
  if (!cur || !cur.lesson) {
    return <span className="text-faint">— ещё не начинал</span>;
  }
  const isNow = row.online && cur.status === "in_progress";
  return (
    <div className="min-w-0">
      <span className="flex items-center gap-1.5">
        {isNow && <span className="h-1.5 w-1.5 shrink-0 rounded-full bg-success" />}
        <span className="truncate font-medium text-fg">{cur.lesson}</span>
      </span>
      <span className="block truncate text-xs text-faint">
        {cur.course}
        {cur.status === "completed" ? " · пройден" : isNow ? " · сейчас" : ""}
      </span>
    </div>
  );
}

export default function ActivityPage() {
  const [onlyOnline, setOnlyOnline] = useState(false);
  // Живая активность: обновляем чаще, чем обычные таблицы.
  const { data = [], isLoading, isFetching } = useGetLiveActivityQuery(300, {
    pollingInterval: 30_000,
  });

  const onlineNow = useMemo(() => data.filter((r) => r.online).length, [data]);
  const activeToday = useMemo(() => data.filter((r) => r.minutesToday > 0).length, [data]);
  const learningNow = useMemo(
    () => data.filter((r) => r.online && r.current?.status === "in_progress").length,
    [data],
  );

  const rows = onlyOnline ? data.filter((r) => r.online) : data;

  return (
    <>
      <PageHeader
        title="Активность"
        subtitle="Кто сейчас онлайн, чем занят и когда был активен"
      />

      <div className="grid gap-[var(--gap)] sm:grid-cols-2 xl:grid-cols-4">
        <StatCard label="Сейчас онлайн" value={onlineNow} icon={<Wifi size={20} />} />
        <StatCard
          label="Учатся прямо сейчас"
          value={learningNow}
          hint="открыт урок"
          icon={<BookOpen size={20} />}
        />
        <StatCard
          label="Были сегодня"
          value={activeToday}
          hint="заходили за день"
          icon={<Clock size={20} />}
        />
        <StatCard label="Всего студентов" value={data.length} icon={<Users size={20} />} />
      </div>

      <Card className="mt-[var(--gap)] overflow-hidden">
        <div className="flex flex-wrap items-center justify-between gap-3 border-b border-line px-[var(--pad)] py-3">
          <div className="flex items-center gap-2">
            <Activity size={18} className="text-accent" />
            <h2 className="text-base font-bold text-fg">Живая активность</h2>
            {isFetching && <Spinner size={14} />}
          </div>
          <label className="flex cursor-pointer items-center gap-2 text-sm text-muted">
            <input
              type="checkbox"
              checked={onlyOnline}
              onChange={(e) => setOnlyOnline(e.target.checked)}
              className="h-4 w-4 accent-[var(--accent)]"
            />
            Только онлайн
          </label>
        </div>

        {isLoading ? (
          <div className="grid place-items-center py-16 text-accent">
            <Spinner size={28} />
          </div>
        ) : rows.length === 0 ? (
          <EmptyState
            title={onlyOnline ? "Сейчас никого нет онлайн" : "Активности пока нет"}
            description="Как только студенты начнут заниматься, здесь появятся их сессии"
            icon={<Activity size={32} />}
          />
        ) : (
          <div className="overflow-x-auto">
            <table className="tbl min-w-[52rem]">
              <thead>
                <tr>
                  <th>Студент</th>
                  <th>Статус</th>
                  <th>Чем занят</th>
                  <th>Активность</th>
                  <th className="num">Сегодня</th>
                  <th className="num">Пройдено</th>
                </tr>
              </thead>
              <tbody className={isFetching ? "opacity-70 transition-opacity" : undefined}>
                {rows.map((row) => (
                  <tr key={row.userId}>
                    <td>
                      <div className="flex items-start gap-2">
                        <span
                          className={`mt-1.5 h-2 w-2 shrink-0 rounded-full ${
                            row.online ? "bg-[var(--success)]" : "bg-[var(--border)]"
                          }`}
                          title={row.online ? "Онлайн" : "Не в сети"}
                        />
                        <span className="min-w-0">
                          <span className="block font-semibold text-fg">
                            {row.fullName || "Без имени"}
                          </span>
                          <span className="block text-xs text-faint">{row.email}</span>
                        </span>
                      </div>
                    </td>
                    <td>
                      <Badge tone={STATUS_TONE[row.status]}>{STATUS_LABEL[row.status]}</Badge>
                    </td>
                    <td className="max-w-[18rem]">
                      <CurrentCell row={row} />
                    </td>
                    <td>
                      <span className={row.online ? "font-semibold text-success" : "text-muted"}>
                        {lastSeenLabel(row.lastSeenAt, row.online)}
                      </span>
                    </td>
                    <td className="num">
                      {row.minutesToday > 0 ? (
                        <span className="font-medium text-fg">{row.minutesToday} мин</span>
                      ) : (
                        <span className="text-faint">—</span>
                      )}
                    </td>
                    <td className="num font-medium text-fg">{row.lessonsCompleted}</td>
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
