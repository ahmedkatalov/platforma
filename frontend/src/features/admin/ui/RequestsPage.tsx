import { useState } from "react";
import { Link } from "react-router-dom";
import { Book, Check, Inbox, Layers, X } from "lucide-react";

import {
  useApproveAccessRequestMutation,
  useApproveCourseRequestMutation,
  useGetAccessRequestsQuery,
  useGetCourseRequestsQuery,
  useRejectAccessRequestMutation,
  useRejectCourseRequestMutation,
} from "@/features/admin/api/adminApi";
import { apiErrorMessage } from "@/shared/api/baseApi";
import type { AccessRequestStatus } from "@/shared/types";
import { Badge, Button, Card, EmptyState, Field, PageHeader, Select, Spinner } from "@/shared/ui";
import { useToast } from "@/shared/ui/ToastProvider";

const STATUS_LABEL: Record<AccessRequestStatus, string> = {
  pending: "Ожидает",
  approved: "Одобрена",
  rejected: "Отклонена",
};
const STATUS_TONE: Record<AccessRequestStatus, "warning" | "success" | "danger"> = {
  pending: "warning",
  approved: "success",
  rejected: "danger",
};
const dateFmt = new Intl.DateTimeFormat("ru-RU", {
  day: "2-digit",
  month: "short",
  hour: "2-digit",
  minute: "2-digit",
});

export default function RequestsPage() {
  const [status, setStatus] = useState<AccessRequestStatus | "">("pending");
  const filter = status || undefined;

  const { data: courseReqs = [], isLoading: loadingCourses } = useGetCourseRequestsQuery(filter, {
    pollingInterval: 30_000,
  });
  const { data: chapterReqs = [], isLoading: loadingChapters } = useGetAccessRequestsQuery(filter, {
    pollingInterval: 30_000,
  });

  const [approveCourse, courseApprove] = useApproveCourseRequestMutation();
  const [rejectCourse, courseReject] = useRejectCourseRequestMutation();
  const [approveChapter, chapterApprove] = useApproveAccessRequestMutation();
  const [rejectChapter, chapterReject] = useRejectAccessRequestMutation();
  const toast = useToast();

  const run = async (fn: () => Promise<unknown>, ok: string) => {
    try {
      await fn();
      toast.success(ok);
    } catch (err) {
      toast.error(apiErrorMessage(err));
    }
  };

  const total = courseReqs.length + chapterReqs.length;
  const loading = loadingCourses || loadingChapters;

  return (
    <>
      <PageHeader
        title="Заявки"
        subtitle="Доступ к курсам и к следующим главам — одобряйте или отклоняйте"
      />

      <Card className="mb-[var(--gap)] flex flex-wrap items-end gap-3 p-[var(--pad)]">
        <div className="w-52">
          <Field label="Статус">
            <Select value={status} onChange={(e) => setStatus(e.target.value as AccessRequestStatus | "")}>
              <option value="pending">Ожидают решения</option>
              <option value="">Все</option>
              <option value="approved">Одобренные</option>
              <option value="rejected">Отклонённые</option>
            </Select>
          </Field>
        </div>
        <div className="ml-auto text-sm text-muted">
          Найдено: <span className="font-bold text-fg">{total}</span>
        </div>
      </Card>

      {loading ? (
        <div className="grid place-items-center py-16 text-accent">
          <Spinner size={28} />
        </div>
      ) : total === 0 ? (
        <Card>
          <EmptyState
            title="Заявок нет"
            description="Когда студент запросит курс или следующую главу, заявка появится здесь"
            icon={<Inbox size={32} />}
          />
        </Card>
      ) : (
        <div className="space-y-[var(--gap)]">
          {courseReqs.length > 0 && (
            <section>
              <h2 className="mb-2 flex items-center gap-2 text-sm font-bold text-fg">
                <Book size={16} className="text-accent" /> Заявки на курсы
              </h2>
              <div className="space-y-3">
                {courseReqs.map((req) => (
                  <RequestRow
                    key={req.id}
                    userId={req.userId}
                    name={req.userName || req.userEmail}
                    status={req.status}
                    target={`Курс «${req.courseTitle}»`}
                    note={req.note}
                    createdAt={req.createdAt}
                    approving={courseApprove.isLoading}
                    rejecting={courseReject.isLoading}
                    onApprove={() => run(() => approveCourse(req.id).unwrap(), "Студент записан на курс")}
                    onReject={() => run(() => rejectCourse(req.id).unwrap(), "Заявка отклонена")}
                  />
                ))}
              </div>
            </section>
          )}

          {chapterReqs.length > 0 && (
            <section>
              <h2 className="mb-2 flex items-center gap-2 text-sm font-bold text-fg">
                <Layers size={16} className="text-accent" /> Заявки на главы
              </h2>
              <div className="space-y-3">
                {chapterReqs.map((req) => (
                  <RequestRow
                    key={req.id}
                    userId={req.userId}
                    name={req.userName || req.userEmail}
                    status={req.status}
                    target={`«${req.courseTitle}» → Глава ${req.chapterNo}: ${req.moduleTitle}`}
                    note={req.note}
                    createdAt={req.createdAt}
                    approving={chapterApprove.isLoading}
                    rejecting={chapterReject.isLoading}
                    onApprove={() => run(() => approveChapter(req.id).unwrap(), "Глава открыта студенту")}
                    onReject={() => run(() => rejectChapter(req.id).unwrap(), "Заявка отклонена")}
                  />
                ))}
              </div>
            </section>
          )}
        </div>
      )}
    </>
  );
}

function RequestRow({
  userId,
  name,
  status,
  target,
  note,
  createdAt,
  approving,
  rejecting,
  onApprove,
  onReject,
}: {
  userId: string;
  name: string;
  status: AccessRequestStatus;
  target: string;
  note: string;
  createdAt: string;
  approving: boolean;
  rejecting: boolean;
  onApprove: () => void;
  onReject: () => void;
}) {
  return (
    <Card className="flex flex-wrap items-center gap-4 p-[var(--pad)]">
      <div className="min-w-0 flex-1">
        <div className="mb-1 flex flex-wrap items-center gap-2">
          <Link to={`/admin/students/${userId}`} className="text-sm font-bold text-fg hover:text-accent">
            {name}
          </Link>
          <Badge tone={STATUS_TONE[status]}>{STATUS_LABEL[status]}</Badge>
        </div>
        <p className="text-sm text-muted">{target}</p>
        {note && <p className="mt-1 text-xs text-faint">«{note}»</p>}
        <p className="mt-1 text-[11px] text-faint">{dateFmt.format(new Date(createdAt))}</p>
      </div>

      {status === "pending" && (
        <div className="flex shrink-0 gap-2">
          <Button variant="primary" icon={<Check size={16} />} loading={approving} onClick={onApprove}>
            Открыть
          </Button>
          <Button variant="ghost" className="text-danger" icon={<X size={16} />} loading={rejecting} onClick={onReject}>
            Отклонить
          </Button>
        </div>
      )}
    </Card>
  );
}
