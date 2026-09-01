import { Link, useParams } from "react-router-dom";

import { useVerifyCertificateQuery } from "@/features/certificates/api/certificatesApi";
import { useTheme } from "@/shared/theme/ThemeProvider";
import { Badge, Button, Card, Spinner } from "@/shared/ui";
import { Check, X, Moon, Sun, Terminal } from "lucide-react";

const dateFmt = new Intl.DateTimeFormat("ru-RU", {
  day: "numeric",
  month: "long",
  year: "numeric",
});

// Публичная страница сертификата: открывается по номеру без авторизации,
// её можно отправить работодателю или распечатать.
export default function CertificatePage() {
  const { serial = "" } = useParams();
  const { mode, toggleMode } = useTheme();
  const { data, isLoading } = useVerifyCertificateQuery(serial, { skip: !serial });

  if (isLoading) {
    return (
      <div className="grid min-h-screen place-items-center text-accent">
        <Spinner size={36} />
      </div>
    );
  }

  const cert = data?.certificate;
  const valid = Boolean(data?.valid);

  return (
    <div className="min-h-screen px-4 py-10">
      <button
        className="btn btn-secondary fixed right-4 top-4 h-9 w-9 !p-0 print:hidden"
        onClick={toggleMode}
        aria-label="Сменить тему"
      >
        {mode === "dark" ? <Sun size={18} /> : <Moon size={18} />}
      </button>

      <div className="mx-auto max-w-3xl">
        <Link to="/" className="mb-6 flex items-center justify-center gap-3 print:hidden">
          <span
            className="grid h-10 w-10 place-items-center rounded-[var(--radius-md)] text-accent-fg"
            style={{ background: "var(--gradient)" }}
          >
            <Terminal size={22} />
          </span>
          <span className="text-lg font-extrabold tracking-tight">
            Okvion <span className="gradient-text">Learning</span>
          </span>
        </Link>

        {!cert ? (
          <Card className="p-8 text-center">
            <span className="mx-auto mb-4 grid h-14 w-14 place-items-center rounded-full bg-[var(--danger-soft)] text-danger">
              <X size={28} />
            </span>
            <h1 className="text-xl font-bold text-fg">Сертификат не найден</h1>
            <p className="mt-2 text-sm text-muted">
              Проверьте номер: он выглядит как <span className="font-mono">DP-2026-XXXXXX</span>
            </p>
          </Card>
        ) : (
          <>
            <Card className="overflow-hidden">
              {/* Верхняя полоса-градиент как «печать» документа */}
              <div className="h-2 w-full" style={{ background: "var(--gradient)" }} />

              <div className="p-8 text-center sm:p-12">
                <p className="text-xs font-bold uppercase tracking-[0.3em] text-faint">
                  Сертификат о прохождении курса
                </p>

                <h1 className="mt-6 text-3xl font-extrabold tracking-tight text-fg sm:text-4xl">
                  {cert.holderName}
                </h1>

                <p className="mt-4 text-sm text-muted">успешно прошёл курс</p>

                <p className="mt-2 text-xl font-bold gradient-text sm:text-2xl">
                  {cert.courseTitle}
                </p>

                <div className="mx-auto mt-8 grid max-w-md grid-cols-3 gap-3">
                  <div className="card-flat py-3">
                    <p className="text-xl font-bold text-fg">{Math.round(cert.score)}%</p>
                    <p className="text-[11px] text-faint">средний балл</p>
                  </div>
                  <div className="card-flat py-3">
                    <p className="text-xl font-bold text-fg">
                      {cert.lessonsCompleted}/{cert.lessonsTotal}
                    </p>
                    <p className="text-[11px] text-faint">уроков</p>
                  </div>
                  <div className="card-flat py-3">
                    <p className="text-xl font-bold text-fg">
                      {new Date(cert.issuedAt).getFullYear()}
                    </p>
                    <p className="text-[11px] text-faint">год выдачи</p>
                  </div>
                </div>

                <div className="mt-8 flex flex-col items-center gap-2">
                  <p className="font-mono text-sm font-bold tracking-wider text-accent">
                    {cert.serial}
                  </p>
                  <p className="text-xs text-faint">
                    выдан {dateFmt.format(new Date(cert.issuedAt))}
                  </p>

                  {valid ? (
                    <Badge tone="success">
                      <Check size={12} /> Сертификат действителен
                    </Badge>
                  ) : (
                    <Badge tone="danger">
                      <X size={12} /> {data?.message ?? "Сертификат отозван"}
                    </Badge>
                  )}
                </div>
              </div>
            </Card>

            <div className="mt-4 flex flex-wrap justify-center gap-2 print:hidden">
              <Button onClick={() => window.print()}>Распечатать</Button>
              <Button
                onClick={() => {
                  void navigator.clipboard.writeText(window.location.href);
                }}
              >
                Скопировать ссылку
              </Button>
              <Link to="/learn" className="btn btn-primary">
                На платформу
              </Link>
            </div>

            <p className="mt-4 text-center text-xs text-faint print:hidden">
              Подлинность сертификата проверяется по этой ссылке — номер уникален.
            </p>
          </>
        )}
      </div>
    </div>
  );
}
