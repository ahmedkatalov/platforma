import { Lock, LockOpen } from "lucide-react";

import { useGetAdminCourseQuery } from "@/features/admin/api/coursesApi";
import {
  useGetUserModuleAccessQuery,
  useSetModuleAccessMutation,
} from "@/features/admin/api/adminApi";
import { apiErrorMessage } from "@/shared/api/baseApi";
import { useToast } from "@/shared/ui/ToastProvider";

// Панель открытия/закрытия глав курса для конкретного студента.
export function ChapterAccess({ userId, courseId }: { userId: string; courseId: string }) {
  const { data: course } = useGetAdminCourseQuery(courseId);
  const { data: accessData } = useGetUserModuleAccessQuery({ userId, courseId });
  const [setAccess, { isLoading }] = useSetModuleAccessMutation();
  const toast = useToast();

  const modules = course?.modules ?? [];
  const granted = new Set(accessData?.granted ?? []);
  if (modules.length === 0) return null;

  const toggle = async (moduleId: string, next: boolean) => {
    try {
      await setAccess({ userId, moduleId, granted: next }).unwrap();
    } catch (err) {
      toast.error(apiErrorMessage(err, "Не удалось изменить доступ"));
    }
  };

  return (
    <div className="mt-2 border-t border-line pt-2">
      <p className="mb-1.5 text-xs font-semibold text-faint">
        Доступ к главам ({granted.size}/{modules.length})
      </p>
      <div className="flex flex-wrap gap-1.5">
        {modules.map((m, i) => {
          const open = granted.has(m.id);
          return (
            <button
              key={m.id}
              onClick={() => toggle(m.id, !open)}
              disabled={isLoading}
              title={m.title}
              className={`flex items-center gap-1 rounded-full border px-2 py-1 text-[11px] font-medium transition-colors disabled:opacity-50 ${
                open
                  ? "border-[var(--success)] bg-[var(--success-soft)] text-success"
                  : "border-line text-faint hover:bg-surface-2"
              }`}
            >
              {open ? <LockOpen size={11} /> : <Lock size={11} />}
              Глава {i + 1}
            </button>
          );
        })}
      </div>
      <p className="mt-1 text-[11px] text-faint">Нажмите на главу, чтобы открыть или закрыть её студенту.</p>
    </div>
  );
}
