import { baseApi } from "@/shared/api/baseApi";
import type { Session, User } from "@/shared/types";

export type SendCodePayload = { email: string; purpose: "registration" | "password_reset" };
export type RegisterPayload = { email: string; fullName: string; password: string; code: string };
export type LoginPayload = { email: string; password: string };
export type ResetPasswordPayload = { email: string; code: string; password: string };

export const authApi = baseApi.injectEndpoints({
  endpoints: (builder) => ({
    sendCode: builder.mutation<{ message: string }, SendCodePayload>({
      query: (body) => ({ url: "/auth/send-code", method: "POST", body }),
    }),
    register: builder.mutation<Session, RegisterPayload>({
      query: (body) => ({ url: "/auth/register", method: "POST", body }),
    }),
    login: builder.mutation<Session, LoginPayload>({
      query: (body) => ({ url: "/auth/login", method: "POST", body }),
    }),
    logout: builder.mutation<{ message: string }, { refreshToken: string }>({
      query: (body) => ({ url: "/auth/logout", method: "POST", body }),
    }),
    resetPassword: builder.mutation<{ message: string }, ResetPasswordPayload>({
      query: (body) => ({ url: "/auth/reset-password", method: "POST", body }),
    }),
    changePassword: builder.mutation<
      { message: string },
      { currentPassword: string; newPassword: string }
    >({
      query: (body) => ({ url: "/me/change-password", method: "POST", body }),
    }),
    updateProfile: builder.mutation<User, { fullName: string }>({
      query: (body) => ({ url: "/me", method: "PATCH", body }),
      invalidatesTags: ["Me"],
    }),
  }),
});

export const {
  useSendCodeMutation,
  useRegisterMutation,
  useLoginMutation,
  useLogoutMutation,
  useResetPasswordMutation,
  useChangePasswordMutation,
  useUpdateProfileMutation,
} = authApi;
