import { useState } from "react";
import { Link } from "react-router-dom";
import { Check, Inbox, X } from "lucide-react";

import {
  useApproveAccessRequestMutation,
  useGetAccessRequestsQuery,
  useRejectAccessRequestMutation,
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
  const { data: requests = [], isLoading } = useGetAccessRequestsQuery(status || undefined, {
    pollingInterval: 30_000,
  });
  const [approve, { isLoading: approving }] = useApproveAccessRequestMutation();
  const [reject, { isLoading: rejecting }] = useRejectAccessRequestMutation();
  const toast = useToast();

  const decide = async (fn: () => Promise<unknown>, ok: string) => {
    try {
      await fn();
      toast.success(ok);
    } catch (err) {
      toast.error(apiErrorMessage(err));
    }
  };

  return (
    <>
      <PageHeader
        title="Заявки на доступ"
        subtitle="Студенты просят открыть следующую главу — одобрите или отклоните"
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
          Найдено: <span className="font-bold text-fg">{requests.length}</span>
        </div>
      </Card>

      {isLoading ? (
        <div className="grid place-items-center py-16 text-accent">
          <Spinner size={28} />
        </div>
      ) : requests.length === 0 ? (
        <Card>
          <EmptyState
            title="Заявок нет"
            description="Когда студент завершит главу и запросит следующую, заявка появится здесь"
            icon={<Inbox size={32} />}
          />
        </Card>
      ) : (
        <div className="space-y-3">
          {requests.map((req) => (
            <Card key={req.id} className="flex flex-wrap items-center gap-4 p-[var(--pad)]">
              <div className="min-w-0 flex-1">
                <div className="mb-1 flex flex-wrap items-center gap-2">
                  <Link
                    to={`/admin/students/${req.userId}`}
                    className="text-sm font-bold text-fg hover:text-accent"
                  >
                    {req.userName || req.userEmail}
                  </Link>
                  <Badge tone={STATUS_TONE[req.status]}>{STATUS_LABEL[req.status]}</Badge>
                </div>
                <p className="text-sm text-muted">
                  Курс «{req.courseTitle}» → <span className="font-semibold text-fg">Глава {req.chapterNo}: {req.moduleTitle}</span>
                </p>
                {req.note && <p className="mt-1 text-xs text-faint">«{req.note}»</p>}
                <p className="mt-1 text-[11px] text-faint">{dateFmt.format(new Date(req.createdAt))}</p>
              </div>

              {req.status === "pending" && (
                <div className="flex shrink-0 gap-2">
                  <Button
                    variant="primary"
                    icon={<Check size={16} />}
                    loading={approving}
                    onClick={() => decide(() => approve(req.id).unwrap(), "Глава открыта студенту")}
                  >
                    Открыть
                  </Button>
                  <Button
                    variant="ghost"
                    className="text-danger"
                    icon={<X size={16} />}
                    loading={rejecting}
                    onClick={() => decide(() => reject(req.id).unwrap(), "Заявка отклонена")}
                  >
                    Отклонить
                  </Button>
                </div>
              )}
            </Card>
          ))}
        </div>
      )}
    </>
  );
}
