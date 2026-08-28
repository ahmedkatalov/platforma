import {
  createApi,
  fetchBaseQuery,
  type BaseQueryFn,
  type FetchArgs,
  type FetchBaseQueryError,
} from "@reduxjs/toolkit/query/react";

import { tokenStorage } from "./tokenStorage";

const rawBaseQuery = fetchBaseQuery({
  baseUrl: "/api",
  prepareHeaders: (headers) => {
    const token = tokenStorage.access();
    if (token) headers.set("Authorization", `Bearer ${token}`);
    return headers;
  },
});

// Один общий промис обновления токена, чтобы параллельные запросы
// не дёргали /auth/refresh несколько раз.
let refreshing: Promise<boolean> | null = null;

async function refreshTokens(): Promise<boolean> {
  const refresh = tokenStorage.refresh();
  if (!refresh) return false;

  const response = await fetch("/api/auth/refresh", {
    method: "POST",
    headers: { "Content-Type": "application/json" },
    body: JSON.stringify({ refreshToken: refresh }),
  });

  if (!response.ok) {
    tokenStorage.clear();
    return false;
  }

  const data = (await response.json()) as { accessToken: string; refreshToken: string };
  tokenStorage.save(data.accessToken, data.refreshToken);
  return true;
}

const baseQueryWithReauth: BaseQueryFn<string | FetchArgs, unknown, FetchBaseQueryError> = async (
  args,
  api,
  extraOptions,
) => {
  let result = await rawBaseQuery(args, api, extraOptions);

  if (result.error?.status === 401 && tokenStorage.refresh()) {
    refreshing ??= refreshTokens().finally(() => {
      refreshing = null;
    });

    const ok = await refreshing;
    if (ok) {
      result = await rawBaseQuery(args, api, extraOptions);
    } else {
      tokenStorage.clear();
      window.dispatchEvent(new CustomEvent("platforma:session-expired"));
    }
  }

  return result;
};

export const baseApi = createApi({
  reducerPath: "api",
  baseQuery: baseQueryWithReauth,
  tagTypes: [
    "Me",
    "Users",
    "User",
    "Courses",
    "Course",
    "Lesson",
    "Overview",
    "Theme",
    "Audit",
    "Progress",
    "Attempts",
  ],
  endpoints: () => ({}),
});

// Достаёт человекочитаемое сообщение об ошибке из ответа бэкенда.
export function apiErrorMessage(error: unknown, fallback = "Что-то пошло не так"): string {
  if (typeof error === "object" && error !== null && "data" in error) {
    const data = (error as { data?: unknown }).data;
    if (typeof data === "string" && data.trim()) return data;
    if (typeof data === "object" && data !== null && "message" in data) {
      const message = (data as { message?: unknown }).message;
      if (typeof message === "string" && message.trim()) return message;
    }
  }
  if (typeof error === "object" && error !== null && "error" in error) {
    const message = (error as { error?: unknown }).error;
    if (typeof message === "string") return message;
  }
  return fallback;
}
