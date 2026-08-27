// Токены живут в localStorage: access — короткий, refresh — на 30 дней.
const ACCESS_KEY = "platforma.access";
const REFRESH_KEY = "platforma.refresh";

export const tokenStorage = {
  access(): string | null {
    return localStorage.getItem(ACCESS_KEY);
  },
  refresh(): string | null {
    return localStorage.getItem(REFRESH_KEY);
  },
  save(access: string, refresh: string) {
    localStorage.setItem(ACCESS_KEY, access);
    localStorage.setItem(REFRESH_KEY, refresh);
  },
  clear() {
    localStorage.removeItem(ACCESS_KEY);
    localStorage.removeItem(REFRESH_KEY);
  },
};
