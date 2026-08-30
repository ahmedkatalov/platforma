import { baseApi } from "@/shared/api/baseApi";
import type {
  Certificate,
  CodeCheckResult,
  LessonView,
  QuizResult,
  TerminalCheckResult,
} from "@/shared/types";

export type QuizAnswerPayload = {
  questionId: string;
  optionIds?: string[]; // choice
  order?: string[]; // order: id шагов в порядке студента
  text?: string; // blank: введённый ответ
  secondsSpent: number;
};

export const lessonApi = baseApi.injectEndpoints({
  endpoints: (builder) => ({
    getLesson: builder.query<LessonView, string>({
      query: (id) => `/lessons/${id}`,
      providesTags: (_r, _e, id) => [{ type: "Lesson", id }],
    }),
    startLesson: builder.mutation<{ message: string }, string>({
      query: (id) => ({ url: `/lessons/${id}/start`, method: "POST" }),
    }),
    completeLesson: builder.mutation<
      { message: string; certificate?: Certificate | null },
      { id: string; seconds: number }
    >({
      query: ({ id, seconds }) => ({
        url: `/lessons/${id}/complete`,
        method: "POST",
        body: { seconds },
      }),
      invalidatesTags: (_r, _e, { id }) => [{ type: "Lesson", id }, "Course", "Progress", "Me", "Attempts"],
    }),
    submitQuiz: builder.mutation<
      QuizResult,
      { id: string; answers: QuizAnswerPayload[]; seconds: number }
    >({
      query: ({ id, answers, seconds }) => ({
        url: `/lessons/${id}/quiz`,
        method: "POST",
        body: { answers, seconds },
      }),
      invalidatesTags: (_r, _e, { id }) => [{ type: "Lesson", id }, "Course", "Progress", "Me", "Attempts"],
    }),
    checkTerminal: builder.mutation<
      TerminalCheckResult,
      { id: string; taskId: string; command: string; usedHint: boolean; seconds: number }
    >({
      query: ({ id, ...body }) => ({ url: `/lessons/${id}/terminal`, method: "POST", body }),
    }),
    checkCode: builder.mutation<CodeCheckResult, { id: string; code: string; seconds: number }>({
      query: ({ id, ...body }) => ({ url: `/lessons/${id}/code`, method: "POST", body }),
      invalidatesTags: (_r, _e, { id }) => [{ type: "Lesson", id }, "Course", "Progress", "Me", "Attempts"],
    }),
  }),
});

export const {
  useGetLessonQuery,
  useStartLessonMutation,
  useCompleteLessonMutation,
  useSubmitQuizMutation,
  useCheckTerminalMutation,
  useCheckCodeMutation,
} = lessonApi;
