// Человекочитаемое «был в сети» по времени последней активности.

const MIN = 60_000;
const HOUR = 60 * MIN;
const DAY = 24 * HOUR;

/**
 * lastSeenLabel возвращает подпись вида «онлайн», «5 минут назад», «вчера».
 * online приходит с сервера (окно 5 минут) — ему верим в первую очередь.
 */
export function lastSeenLabel(iso: string | null | undefined, online?: boolean): string {
  if (online) return "онлайн";
  if (!iso) return "не заходил";

  const then = new Date(iso).getTime();
  if (Number.isNaN(then)) return "не заходил";

  const diff = Date.now() - then;
  if (diff < 0) return "онлайн";
  if (diff < MIN) return "только что";

  if (diff < HOUR) {
    const m = Math.floor(diff / MIN);
    return `${m} ${plural(m, "минуту", "минуты", "минут")} назад`;
  }
  if (diff < DAY) {
    const h = Math.floor(diff / HOUR);
    return `${h} ${plural(h, "час", "часа", "часов")} назад`;
  }
  if (diff < 2 * DAY) return "вчера";

  const d = Math.floor(diff / DAY);
  if (d < 30) return `${d} ${plural(d, "день", "дня", "дней")} назад`;

  return new Intl.DateTimeFormat("ru-RU", { day: "2-digit", month: "short" }).format(then);
}

// Русское склонение числительных: 1 минуту, 2 минуты, 5 минут.
function plural(n: number, one: string, few: string, many: string): string {
  const mod10 = n % 10;
  const mod100 = n % 100;
  if (mod10 === 1 && mod100 !== 11) return one;
  if (mod10 >= 2 && mod10 <= 4 && (mod100 < 10 || mod100 >= 20)) return few;
  return many;
}
