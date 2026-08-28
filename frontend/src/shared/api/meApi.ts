import { baseApi } from "@/shared/api/baseApi";
import type {
  ActivityDay,
  Attempt,
  Enrollment,
  QuizCard,
  QuizStats,
  StudentSummary,
  User,
} from "@/shared/types";

export type MeResponse = { user: User; enrollments: Enrollment[] };
export type MyStats = {
  summary: StudentSummary;
  activity: ActivityDay[];
  streak: number;
  quiz: QuizStats;
};

export const meApi = baseApi.injectEndpoints({
  endpoints: (builder) => ({
    getMe: builder.query<MeResponse, void>({
      query: () => "/me",
      providesTags: ["Me"],
    }),
    getMyStats: builder.query<MyStats, number | void>({
      query: (days) => `/me/stats?days=${days ?? 30}`,
      providesTags: ["Progress"],
    }),
    getMyAttempts: builder.query<Attempt[], number | void>({
      query: (limit) => `/me/attempts?limit=${limit ?? 20}`,
      providesTags: ["Attempts"],
    }),
    getMyQuizzes: builder.query<QuizCard[], void>({
      query: () => "/me/quizzes",
      providesTags: ["Progress"],
    }),
    trackActivity: builder.mutation<{ message: string }, { seconds: number }>({
      query: (body) => ({ url: "/me/activity", method: "POST", body }),
    }),
    getPreferences: builder.query<{ theme: unknown }, void>({
      query: () => "/me/preferences",
      providesTags: ["Theme"],
    }),
    savePreferences: builder.mutation<{ theme: unknown }, unknown>({
      query: (theme) => ({ url: "/me/preferences", method: "PUT", body: { theme } }),
    }),
    resetPreferences: builder.mutation<{ theme: unknown }, void>({
      query: () => ({ url: "/me/preferences", method: "DELETE" }),
    }),
    getPublicTheme: builder.query<{ settings: unknown }, void>({
      query: () => "/theme",
    }),
  }),
});

export const {
  useGetMeQuery,
  useGetMyStatsQuery,
  useGetMyAttemptsQuery,
  useGetMyQuizzesQuery,
  useTrackActivityMutation,
  useGetPreferencesQuery,
  useSavePreferencesMutation,
  useResetPreferencesMutation,
  useGetPublicThemeQuery,
} = meApi;
