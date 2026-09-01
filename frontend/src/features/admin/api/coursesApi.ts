import { baseApi } from "@/shared/api/baseApi";
import type {
  Course,
  CourseLevel,
  CourseStatus,
  Lesson,
  LessonKind,
  LessonProgress,
  Module,
} from "@/shared/types";

export type CoursePayload = {
  slug: string;
  title: string;
  subtitle: string;
  description: string;
  coverUrl: string;
  level: CourseLevel;
  tags: string[];
  status: CourseStatus;
  position: number;
};

export type ModulePayload = { title: string; summary: string; position: number };

export type LessonPayload = {
  title: string;
  kind: LessonKind;
  summary: string;
  content: Record<string, unknown>;
  durationMin: number;
  position: number;
};

export const coursesApi = baseApi.injectEndpoints({
  endpoints: (builder) => ({
    // --- Витрина для студента ---
    getStudentCourses: builder.query<
      {
        course: Course;
        enrolled: boolean;
        completedLessons: number;
        requestStatus?: string;
      }[],
      void
    >({
      query: () => "/courses",
      providesTags: ["Courses", "Progress", "CourseRequests"],
    }),
    requestCourseAccess: builder.mutation<{ message: string }, { courseId: string }>({
      query: (body) => ({ url: "/courses/request-enroll", method: "POST", body }),
      invalidatesTags: ["Courses", "CourseRequests"],
    }),
    getStudentCourse: builder.query<
      {
        course: Course;
        enrolled: boolean;
        progress: LessonProgress[];
        moduleAccess?: Record<string, boolean>;
        requests?: Record<string, string>;
      },
      string
    >({
      query: (slug) => `/courses/${slug}`,
      providesTags: (_r, _e, slug) => [{ type: "Course", id: slug }, "Access"],
    }),
    requestModuleAccess: builder.mutation<{ message: string }, { slug: string; moduleId: string }>({
      query: ({ moduleId }) => ({
        url: "/courses/request-access",
        method: "POST",
        body: { moduleId },
      }),
      invalidatesTags: (_r, _e, { slug }) => [{ type: "Course", id: slug }, "Access"],
    }),

    // --- Редактор администратора ---
    getAdminCourses: builder.query<Course[], CourseStatus | void>({
      query: (status) => (status ? `/admin/courses?status=${status}` : "/admin/courses"),
      providesTags: ["Courses"],
    }),
    getAdminCourse: builder.query<Course, string>({
      query: (id) => `/admin/courses/${id}`,
      providesTags: (_r, _e, id) => [{ type: "Course", id }],
    }),
    createCourse: builder.mutation<Course, CoursePayload>({
      query: (body) => ({ url: "/admin/courses", method: "POST", body }),
      invalidatesTags: ["Courses", "Overview"],
    }),
    updateCourse: builder.mutation<Course, { id: string } & CoursePayload>({
      query: ({ id, ...body }) => ({ url: `/admin/courses/${id}`, method: "PUT", body }),
      invalidatesTags: (_r, _e, { id }) => ["Courses", { type: "Course", id }],
    }),
    deleteCourse: builder.mutation<{ message: string }, string>({
      query: (id) => ({ url: `/admin/courses/${id}`, method: "DELETE" }),
      invalidatesTags: ["Courses", "Overview"],
    }),
    importCourse: builder.mutation<
      { course: Course; modules: number; lessons: number; message: string },
      { raw: string; replace: boolean }
    >({
      query: ({ raw, replace }) => ({
        url: `/admin/courses/import${replace ? "?replace=true" : ""}`,
        method: "POST",
        body: raw, // сырой текст файла — разбирает сервер (надёжнее и с точной ошибкой)
        headers: { "Content-Type": "application/json" },
      }),
      // Обновление на месте затрагивает структуру, прогресс и доступы.
      invalidatesTags: ["Courses", "Overview", "Course", "Progress", "Access"],
    }),

    createModule: builder.mutation<Module, { courseId: string } & ModulePayload>({
      query: ({ courseId, ...body }) => ({
        url: `/admin/courses/${courseId}/modules`,
        method: "POST",
        body,
      }),
      invalidatesTags: (_r, _e, { courseId }) => [{ type: "Course", id: courseId }, "Courses"],
    }),
    updateModule: builder.mutation<Module, { moduleId: string; courseId: string } & ModulePayload>({
      query: ({ moduleId, courseId: _courseId, ...body }) => ({
        url: `/admin/courses/modules/${moduleId}`,
        method: "PUT",
        body,
      }),
      invalidatesTags: (_r, _e, { courseId }) => [{ type: "Course", id: courseId }],
    }),
    deleteModule: builder.mutation<{ message: string }, { moduleId: string; courseId: string }>({
      query: ({ moduleId }) => ({ url: `/admin/courses/modules/${moduleId}`, method: "DELETE" }),
      invalidatesTags: (_r, _e, { courseId }) => [{ type: "Course", id: courseId }, "Courses"],
    }),

    createLesson: builder.mutation<Lesson, { moduleId: string; courseId: string } & LessonPayload>({
      query: ({ moduleId, courseId: _courseId, ...body }) => ({
        url: `/admin/courses/modules/${moduleId}/lessons`,
        method: "POST",
        body,
      }),
      invalidatesTags: (_r, _e, { courseId }) => [{ type: "Course", id: courseId }, "Courses"],
    }),
    updateLesson: builder.mutation<Lesson, { lessonId: string; courseId: string } & LessonPayload>({
      query: ({ lessonId, courseId: _courseId, ...body }) => ({
        url: `/admin/courses/lessons/${lessonId}`,
        method: "PUT",
        body,
      }),
      invalidatesTags: (_r, _e, { courseId }) => [{ type: "Course", id: courseId }],
    }),
    deleteLesson: builder.mutation<{ message: string }, { lessonId: string; courseId: string }>({
      query: ({ lessonId }) => ({ url: `/admin/courses/lessons/${lessonId}`, method: "DELETE" }),
      invalidatesTags: (_r, _e, { courseId }) => [{ type: "Course", id: courseId }, "Courses"],
    }),
  }),
});

export const {
  useGetStudentCoursesQuery,
  useRequestCourseAccessMutation,
  useGetStudentCourseQuery,
  useRequestModuleAccessMutation,
  useGetAdminCoursesQuery,
  useGetAdminCourseQuery,
  useCreateCourseMutation,
  useUpdateCourseMutation,
  useDeleteCourseMutation,
  useImportCourseMutation,
  useCreateModuleMutation,
  useUpdateModuleMutation,
  useDeleteModuleMutation,
  useCreateLessonMutation,
  useUpdateLessonMutation,
  useDeleteLessonMutation,
} = coursesApi;
