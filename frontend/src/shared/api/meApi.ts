import { baseApi } from "@/shared/api/baseApi";
import type {
  ActivityDay,
  Attempt,
  ContactSettings,
  Enrollment,
  Note,
  QuizCard,
  QuizStats,
  StudentSummary,
  User,
} from "@/shared/types";

export type MeResponse = {
  user: User;
  enrollments: Enrollment[];
  sandboxAvailable?: boolean;
};
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
    getMyNotes: builder.query<Note[], void>({
      query: () => "/me/notes",
      providesTags: ["Notes"],
    }),
    createNote: builder.mutation<Note, { lessonId: string; quote: string; body?: string }>({
      query: (body) => ({ url: "/me/notes", method: "POST", body }),
      invalidatesTags: ["Notes"],
    }),
    updateNote: builder.mutation<Note, { id: string; body: string }>({
      query: ({ id, body }) => ({ url: `/me/notes/${id}`, method: "PATCH", body: { body } }),
      invalidatesTags: ["Notes"],
    }),
    deleteNote: builder.mutation<{ message: string }, string>({
      query: (id) => ({ url: `/me/notes/${id}`, method: "DELETE" }),
      invalidatesTags: ["Notes"],
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
    getPublicContacts: builder.query<{ contacts: ContactSettings | null }, void>({
      query: () => "/contacts",
      providesTags: ["Contacts"],
    }),
  }),
});

export const {
  useGetMeQuery,
  useGetMyStatsQuery,
  useGetMyAttemptsQuery,
  useGetMyNotesQuery,
  useCreateNoteMutation,
  useUpdateNoteMutation,
  useDeleteNoteMutation,
  useGetMyQuizzesQuery,
  useTrackActivityMutation,
  useGetPreferencesQuery,
  useSavePreferencesMutation,
  useResetPreferencesMutation,
  useGetPublicThemeQuery,
  useGetPublicContactsQuery,
} = meApi;
