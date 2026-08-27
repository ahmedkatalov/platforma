import { baseApi } from "@/shared/api/baseApi";
import type {
  ActivityDay,
  AdminOverview,
  AuditEntry,
  CreatedStudent,
  Enrollment,
  Paginated,
  Role,
  StudentSummary,
  User,
  UserStatus,
} from "@/shared/types";

export type UsersQuery = {
  search?: string;
  role?: Role | "";
  status?: UserStatus | "";
  page?: number;
  limit?: number;
};

export type UserDetails = {
  user: User;
  enrollments: Enrollment[];
  summary: StudentSummary;
  activity: ActivityDay[];
};

export type CreateUserPayload = {
  email: string;
  fullName: string;
  password?: string;
  role?: Role;
  sendMail: boolean;
};

export const adminApi = baseApi.injectEndpoints({
  endpoints: (builder) => ({
    getOverview: builder.query<AdminOverview, void>({
      query: () => "/admin/overview",
      providesTags: ["Overview"],
    }),
    getUsers: builder.query<Paginated<User>, UsersQuery | void>({
      query: (params) => {
        const q = new URLSearchParams();
        if (params?.search) q.set("search", params.search);
        if (params?.role) q.set("role", params.role);
        if (params?.status) q.set("status", params.status);
        q.set("page", String(params?.page ?? 1));
        q.set("limit", String(params?.limit ?? 50));
        return `/admin/users?${q.toString()}`;
      },
      providesTags: ["Users"],
    }),
    getUser: builder.query<UserDetails, string>({
      query: (id) => `/admin/users/${id}`,
      providesTags: (_r, _e, id) => [{ type: "User", id }],
    }),
    createUser: builder.mutation<CreatedStudent, CreateUserPayload>({
      query: (body) => ({ url: "/admin/users", method: "POST", body }),
      invalidatesTags: ["Users", "Overview"],
    }),
    updateUser: builder.mutation<
      User,
      { id: string; fullName?: string; email?: string; role?: Role; status?: UserStatus }
    >({
      query: ({ id, ...body }) => ({ url: `/admin/users/${id}`, method: "PATCH", body }),
      invalidatesTags: (_r, _e, { id }) => ["Users", "Overview", { type: "User", id }],
    }),
    deleteUser: builder.mutation<{ message: string }, string>({
      query: (id) => ({ url: `/admin/users/${id}`, method: "DELETE" }),
      invalidatesTags: ["Users", "Overview"],
    }),
    resetUserPassword: builder.mutation<CreatedStudent, { id: string; sendMail: boolean }>({
      query: ({ id, sendMail }) => ({
        url: `/admin/users/${id}/reset-password`,
        method: "POST",
        body: { sendMail },
      }),
    }),
    enroll: builder.mutation<Enrollment[], { userId: string; courseId: string }>({
      query: ({ userId, courseId }) => ({
        url: `/admin/users/${userId}/enrollments`,
        method: "POST",
        body: { courseId },
      }),
      invalidatesTags: (_r, _e, { userId }) => [{ type: "User", id: userId }, "Courses", "Overview"],
    }),
    unenroll: builder.mutation<{ message: string }, { userId: string; courseId: string }>({
      query: ({ userId, courseId }) => ({
        url: `/admin/users/${userId}/enrollments/${courseId}`,
        method: "DELETE",
      }),
      invalidatesTags: (_r, _e, { userId }) => [{ type: "User", id: userId }, "Courses", "Overview"],
    }),
    getStudentsProgress: builder.query<StudentSummary[], number | void>({
      query: (limit) => `/admin/students-progress?limit=${limit ?? 100}`,
      providesTags: ["Progress"],
    }),
    getAudit: builder.query<AuditEntry[], number | void>({
      query: (limit) => `/admin/audit?limit=${limit ?? 50}`,
      providesTags: ["Audit"],
    }),
    getPlatformTheme: builder.query<{ settings: unknown }, void>({
      query: () => "/admin/theme",
    }),
    savePlatformTheme: builder.mutation<{ settings: unknown }, unknown>({
      query: (settings) => ({ url: "/admin/theme", method: "PUT", body: { settings } }),
    }),
  }),
});

export const {
  useGetOverviewQuery,
  useGetUsersQuery,
  useGetUserQuery,
  useCreateUserMutation,
  useUpdateUserMutation,
  useDeleteUserMutation,
  useResetUserPasswordMutation,
  useEnrollMutation,
  useUnenrollMutation,
  useGetStudentsProgressQuery,
  useGetAuditQuery,
  useGetPlatformThemeQuery,
  useSavePlatformThemeMutation,
} = adminApi;
