import { useGetAuditQuery } from "@/features/admin/api/adminApi";
import { Badge, Card, EmptyState, PageHeader, Spinner } from "@/shared/ui";
import { Shield } from "lucide-react";

const dateFmt = new Intl.DateTimeFormat("ru-RU", {
  day: "2-digit",
  month: "2-digit",
  year: "numeric",
  hour: "2-digit",
  minute: "2-digit",
});

const ACTION_LABEL: Record<string, string> = {
  "user.create": "Создан аккаунт",
  "user.update": "Изменён аккаунт",
  "user.delete": "Удалён аккаунт",
  "user.reset_password": "Сброшен пароль",
  "course.create": "Создан курс",
  "course.update": "Изменён курс",
  "course.delete": "Удалён курс",
  "enrollment.create": "Назначен курс",
  "enrollment.delete": "Снят курс",
  "theme.update": "Изменено оформление",
};

export default function AuditPage() {
  const { data: entries = [], isLoading } = useGetAuditQuery(200);

  return (
    <>
      <PageHeader title="Журнал действий" subtitle="Кто и что менял на платформе" />

      <Card className="overflow-hidden">
        {isLoading ? (
          <div className="grid place-items-center py-16 text-accent">
            <Spinner size={28} />
          </div>
        ) : entries.length === 0 ? (
          <EmptyState title="Записей пока нет" icon={<Shield size={32} />} />
        ) : (
          <div className="overflow-x-auto">
            <table className="tbl min-w-[40rem]">
              <thead>
                <tr>
                  <th>Когда</th>
                  <th>Кто</th>
                  <th>Действие</th>
                  <th>Объект</th>
                </tr>
              </thead>
              <tbody>
                {entries.map((entry) => (
                  <tr key={entry.id}>
                    <td className="whitespace-nowrap">
                      {dateFmt.format(new Date(entry.createdAt))}
                    </td>
                    <td className="font-medium text-fg">{entry.actorName || "система"}</td>
                    <td>
                      <Badge tone="accent">{ACTION_LABEL[entry.action] ?? entry.action}</Badge>
                    </td>
                    <td className="font-mono text-xs text-faint">
                      {entry.entity}
                      {entry.entityId && `: ${entry.entityId.slice(0, 8)}`}
                    </td>
                  </tr>
                ))}
              </tbody>
            </table>
          </div>
        )}
      </Card>
    </>
  );
}
