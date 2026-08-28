import { baseApi } from "@/shared/api/baseApi";
import type { Asset, Certificate, CertificateCheck } from "@/shared/types";

export const certificatesApi = baseApi.injectEndpoints({
  endpoints: (builder) => ({
    // Публичная проверка по номеру — работает и без авторизации.
    verifyCertificate: builder.query<CertificateCheck, string>({
      query: (serial) => `/certificates/${encodeURIComponent(serial)}`,
    }),
    getMyCertificates: builder.query<Certificate[], void>({
      query: () => "/me/certificates",
      providesTags: ["Certificates"],
    }),
    getCertificates: builder.query<Certificate[], number | void>({
      query: (limit) => `/admin/certificates?limit=${limit ?? 100}`,
      providesTags: ["Certificates"],
    }),
    revokeCertificate: builder.mutation<{ message: string }, string>({
      query: (id) => ({ url: `/admin/certificates/${id}/revoke`, method: "POST" }),
      invalidatesTags: ["Certificates"],
    }),
    restoreCertificate: builder.mutation<{ message: string }, string>({
      query: (id) => ({ url: `/admin/certificates/${id}/restore`, method: "POST" }),
      invalidatesTags: ["Certificates"],
    }),

    // --- Файлы уроков ---
    getAssets: builder.query<Asset[], void>({
      query: () => "/admin/uploads",
      providesTags: ["Assets"],
    }),
    uploadAsset: builder.mutation<Asset, File>({
      query: (file) => {
        const form = new FormData();
        form.append("file", file);
        return { url: "/admin/uploads", method: "POST", body: form };
      },
      invalidatesTags: ["Assets"],
    }),
    deleteAsset: builder.mutation<{ message: string }, string>({
      query: (id) => ({ url: `/admin/uploads/${id}`, method: "DELETE" }),
      invalidatesTags: ["Assets"],
    }),
  }),
});

export const {
  useVerifyCertificateQuery,
  useGetMyCertificatesQuery,
  useGetCertificatesQuery,
  useRevokeCertificateMutation,
  useRestoreCertificateMutation,
  useGetAssetsQuery,
  useUploadAssetMutation,
  useDeleteAssetMutation,
} = certificatesApi;
