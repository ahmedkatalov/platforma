import { useState } from "react";

import {
  useGetCertificatesQuery,
  useRestoreCertificateMutation,
  useRevokeCertificateMutation,
} from "@/features/certificates/api/certificatesApi";
import { apiErrorMessage } from "@/shared/api/baseApi";
import { downloadFile } from "@/shared/lib/download";
import { Badge, Button, Card, EmptyState, Input, PageHeader, Spinner } from "@/shared/ui";
import { Check, X, Search, Shield } from "lucide-react";
import { useToast } from "@/shared/ui/ToastProvider";

const dateFmt = new Intl.DateTimeFormat("ru-RU", {
  day: "2-digit",
  month: "2-digit",
  year: "numeric",
});

export default function CertificatesPage() {
  const [search, setSearch] = useState("");
  const { data: certificates = [], isLoading } = useGetCertificatesQuery(200);
  const [revoke] = useRevokeCertificateMutation();
  const [restore] = useRestoreCertificateMutation();
  const toast = useToast();

  const query = search.trim().toLowerCase();
  const items = query
    ? certificates.filter(
        (cert) =>
          cert.serial.toLowerCase().includes(query) ||
          cert.holderName.toLowerCase().includes(query) ||
          cert.courseTitle.toLowerCase().includes(query),
      )
    : certificates;

  const exportCsv = async () => {
    try {
      await downloadFile("/api/admin/reports/certificates.csv", "certificates.csv");
    } catch (err) {
      toast.error(apiErrorMessage(err, "Не удалось скачать отчёт"));
    }
  };

  const toggle = async (id: string, revoked: boolean) => {
    try {
      if (revoked) {
        await restore(id).unwrap();
        toast.success("Сертификат восстановлен");
      } else {
        if (!window.confirm("Отозвать сертификат? Проверка по номеру покажет, что он недействителен."))
          return;
        await revoke(id).unwrap();
        toast.success("Сертификат отозван");
      }
    } catch (err) {
      toast.error(apiErrorMessage(err));
    }
  };

  return (
    <>
      <PageHeader
        title="Сертификаты"
        subtitle="Выдаются автоматически при прохождении всех уроков курса"
        actions={<Button onClick={exportCsv}>Выгрузить CSV</Button>}
      />

      <Card className="mb-[var(--gap)] p-[var(--pad)]">
        <div className="relative">
          <Search
            size={16}
            className="pointer-events-none absolute left-3 top-1/2 -translate-y-1/2 text-faint"
          />
          <Input
            value={search}
            onChange={(e) => setSearch(e.target.value)}
            placeholder="Номер, студент или курс"
            className="pl-9"
          />
        </div>
      </Card>

      <Card className="overflow-hidden">
        {isLoading ? (
          <div className="grid place-items-center py-16 text-accent">
            <Spinner size={28} />
          </div>
        ) : items.length === 0 ? (
          <EmptyState
            title="Сертификатов пока нет"
            description="Первый сертификат появится, когда студент пройдёт курс целиком"
            icon={<Shield size={32} />}
          />
        ) : (
          <div className="overflow-x-auto">
            <table className="w-full min-w-[52rem] text-sm">
              <thead>
                <tr className="border-b border-line text-left text-xs uppercase tracking-wide text-faint">
                  <th className="px-4 py-3 font-semibold">Номер</th>
                  <th className="px-4 py-3 font-semibold">Студент</th>
                  <th className="px-4 py-3 font-semibold">Курс</th>
                  <th className="px-4 py-3 font-semibold">Балл</th>
                  <th className="px-4 py-3 font-semibold">Выдан</th>
                  <th className="px-4 py-3 font-semibold">Статус</th>
                  <th className="px-4 py-3" />
                </tr>
              </thead>
              <tbody>
                {items.map((cert) => (
                  <tr key={cert.id} className="border-b border-line/60 last:border-0 hover:bg-surface-2">
                    <td className="px-4 py-3">
                      <a
                        href={`/certificates/${cert.serial}`}
                        target="_blank"
                        rel="noreferrer"
                        className="font-mono text-xs font-bold text-accent hover:underline"
                      >
                        {cert.serial}
                      </a>
                    </td>
                    <td className="px-4 py-3 font-medium text-fg">{cert.holderName}</td>
                    <td className="px-4 py-3 text-muted">{cert.courseTitle}</td>
                    <td className="px-4 py-3">
                      <span className="font-bold text-fg">{Math.round(cert.score)}%</span>
                      <span className="ml-1.5 text-xs text-faint">
                        {cert.lessonsCompleted}/{cert.lessonsTotal}
                      </span>
                    </td>
                    <td className="whitespace-nowrap px-4 py-3 text-muted">
                      {dateFmt.format(new Date(cert.issuedAt))}
                    </td>
                    <td className="px-4 py-3">
                      {cert.revokedAt ? (
                        <Badge tone="danger">Отозван</Badge>
                      ) : (
                        <Badge tone="success">Действителен</Badge>
                      )}
                    </td>
                    <td className="px-4 py-3">
                      <div className="flex justify-end">
                        <Button
                          variant="ghost"
                          className={`h-8 !px-2 ${cert.revokedAt ? "text-success" : "text-danger"}`}
                          onClick={() => toggle(cert.id, Boolean(cert.revokedAt))}
                          title={cert.revokedAt ? "Восстановить" : "Отозвать"}
                        >
                          {cert.revokedAt ? <Check size={16} /> : <X size={16} />}
                        </Button>
                      </div>
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
