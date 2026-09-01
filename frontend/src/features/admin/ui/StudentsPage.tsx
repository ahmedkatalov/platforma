import { useState, type FormEvent } from "react";
import { Link } from "react-router-dom";

import {
  useCreateUserMutation,
  useDeleteUserMutation,
  useGetUsersQuery,
  useUpdateUserMutation,
} from "@/features/admin/api/adminApi";
import { apiErrorMessage } from "@/shared/api/baseApi";
import { downloadFile } from "@/shared/lib/download";
import { lastSeenLabel } from "@/shared/lib/time";
import type { CreatedStudent, Role, User, UserStatus } from "@/shared/types";
import {
  Badge,
  Button,
  Card,
  EmptyState,
  Field,
  Input,
  Modal,
  PageHeader,
  Select,
  Spinner,
} from "@/shared/ui";
import { Key, Plus, Search, Trash2, Users } from "lucide-react";
import { useToast } from "@/shared/ui/ToastProvider";

const STATUS_TONE: Record<UserStatus, "success" | "warning" | "danger"> = {
  active: "success",
  invited: "warning",
  blocked: "danger",
};

const STATUS_LABEL: Record<UserStatus, string> = {
  active: "Активен",
  invited: "Приглашён",
  blocked: "Заблокирован",
};

const dateFmt = new Intl.DateTimeFormat("ru-RU", { day: "2-digit", month: "2-digit", year: "numeric" });

export default function StudentsPage() {
  const [search, setSearch] = useState("");
  const [status, setStatus] = useState<UserStatus | "">("");
  const [role, setRole] = useState<Role | "">("");
  const [createOpen, setCreateOpen] = useState(false);
  const [created, setCreated] = useState<CreatedStudent | null>(null);

  const { data, isLoading, isFetching } = useGetUsersQuery(
    { search, status, role, limit: 100 },
    { pollingInterval: 60_000 },
  );
  const [updateUser] = useUpdateUserMutation();
  const [deleteUser] = useDeleteUserMutation();
  const toast = useToast();

  const users = data?.items ?? [];

  const toggleBlock = async (user: User) => {
    try {
      await updateUser({
        id: user.id,
        status: user.status === "blocked" ? "active" : "blocked",
      }).unwrap();
      toast.success(user.status === "blocked" ? "Аккаунт разблокирован" : "Аккаунт заблокирован");
    } catch (err) {
      toast.error(apiErrorMessage(err));
    }
  };

  const exportCsv = async () => {
    try {
      await downloadFile("/api/admin/reports/students.csv", "students.csv");
    } catch (err) {
      toast.error(apiErrorMessage(err, "Не удалось скачать отчёт"));
    }
  };

  const removeUser = async (user: User) => {
    if (!window.confirm(`Удалить аккаунт ${user.email}? Действие необратимо.`)) return;
    try {
      await deleteUser(user.id).unwrap();
      toast.success("Аккаунт удалён");
    } catch (err) {
      toast.error(apiErrorMessage(err));
    }
  };

  return (
    <>
      <PageHeader
        title="Студенты"
        subtitle="Аккаунты, доступы и статусы"
        actions={
          <>
            <Button onClick={exportCsv}>Выгрузить CSV</Button>
            <Button variant="primary" icon={<Plus size={18} />} onClick={() => setCreateOpen(true)}>
              Создать аккаунт
            </Button>
          </>
        }
      />

      <Card className="mb-[var(--gap)] flex flex-wrap items-end gap-3 p-[var(--pad)]">
        <div className="min-w-[14rem] flex-1">
          <Field label="Поиск">
            <div className="relative">
              <Search
                size={16}
                className="pointer-events-none absolute left-3 top-1/2 -translate-y-1/2 text-faint"
              />
              <Input
                value={search}
                onChange={(e) => setSearch(e.target.value)}
                placeholder="Имя или почта"
                className="pl-9"
              />
            </div>
          </Field>
        </div>

        <div className="w-40">
          <Field label="Статус">
            <Select value={status} onChange={(e) => setStatus(e.target.value as UserStatus | "")}>
              <option value="">Все</option>
              <option value="active">Активные</option>
              <option value="invited">Приглашённые</option>
              <option value="blocked">Заблокированные</option>
            </Select>
          </Field>
        </div>

        <div className="w-40">
          <Field label="Роль">
            <Select value={role} onChange={(e) => setRole(e.target.value as Role | "")}>
              <option value="">Все</option>
              <option value="student">Студенты</option>
              <option value="admin">Администраторы</option>
            </Select>
          </Field>
        </div>

        <div className="ml-auto text-sm text-muted">
          Найдено: <span className="font-bold text-fg">{data?.total ?? 0}</span>
        </div>
      </Card>

      <Card className="overflow-hidden">
        {isLoading ? (
          <div className="grid place-items-center py-16 text-accent">
            <Spinner size={28} />
          </div>
        ) : users.length === 0 ? (
          <EmptyState
            title="Аккаунтов не найдено"
            description="Измените фильтры или создайте первый аккаунт студента"
            icon={<Users size={32} />}
            action={
              <Button variant="primary" onClick={() => setCreateOpen(true)}>
                Создать аккаунт
              </Button>
            }
          />
        ) : (
          <div className="overflow-x-auto">
            <table className="w-full min-w-[52rem] text-sm">
              <thead>
                <tr className="border-b border-line text-left text-xs uppercase tracking-wide text-faint">
                  <th className="px-4 py-3 font-semibold">Студент</th>
                  <th className="px-4 py-3 font-semibold">Роль</th>
                  <th className="px-4 py-3 font-semibold">Статус</th>
                  <th className="px-4 py-3 font-semibold">Активность</th>
                  <th className="px-4 py-3 font-semibold">Последний вход</th>
                  <th className="px-4 py-3 font-semibold">Создан</th>
                  <th className="px-4 py-3" />
                </tr>
              </thead>
              <tbody className={isFetching ? "opacity-60 transition-opacity" : undefined}>
                {users.map((user) => (
                  <tr key={user.id} className="border-b border-line/60 last:border-0 hover:bg-surface-2">
                    <td className="px-4 py-3">
                      <Link to={`/admin/students/${user.id}`} className="flex items-start gap-2">
                        <span
                          className={`mt-1.5 h-2 w-2 shrink-0 rounded-full ${
                            user.online ? "bg-[var(--success)]" : "bg-[var(--border)]"
                          }`}
                          title={user.online ? "Онлайн" : "Не в сети"}
                        />
                        <span className="min-w-0">
                          <span className="block font-semibold text-fg hover:text-accent">
                            {user.fullName || "Без имени"}
                          </span>
                          <span className="block text-xs text-faint">{user.email}</span>
                        </span>
                      </Link>
                    </td>
                    <td className="px-4 py-3">
                      <Badge tone={user.role === "admin" ? "accent" : "default"}>
                        {user.role === "admin" ? "Админ" : "Студент"}
                      </Badge>
                    </td>
                    <td className="px-4 py-3">
                      <Badge tone={STATUS_TONE[user.status]}>{STATUS_LABEL[user.status]}</Badge>
                    </td>
                    <td className="px-4 py-3">
                      <span className={user.online ? "font-semibold text-success" : "text-muted"}>
                        {lastSeenLabel(user.lastSeenAt, user.online)}
                      </span>
                    </td>
                    <td className="px-4 py-3 text-muted">
                      {user.lastLoginAt ? dateFmt.format(new Date(user.lastLoginAt)) : "—"}
                    </td>
                    <td className="px-4 py-3 text-muted">
                      {dateFmt.format(new Date(user.createdAt))}
                    </td>
                    <td className="px-4 py-3">
                      <div className="flex justify-end gap-1">
                        <Button
                          variant="ghost"
                          className="h-8 !px-2"
                          onClick={() => toggleBlock(user)}
                          title={user.status === "blocked" ? "Разблокировать" : "Заблокировать"}
                        >
                          <Key size={16} />
                        </Button>
                        <Button
                          variant="ghost"
                          className="h-8 !px-2 text-danger"
                          onClick={() => removeUser(user)}
                          title="Удалить"
                        >
                          <Trash2 size={16} />
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

      <CreateStudentModal
        open={createOpen}
        onClose={() => setCreateOpen(false)}
        onCreated={(result) => {
          setCreateOpen(false);
          setCreated(result);
        }}
      />

      <Modal
        open={Boolean(created)}
        onClose={() => setCreated(null)}
        title="Аккаунт создан"
        footer={
          <Button variant="primary" onClick={() => setCreated(null)}>
            Готово
          </Button>
        }
      >
        {created && (
          <div className="space-y-4 text-sm">
            <p className="text-muted">
              Передайте студенту эти данные для первого входа. Пароль показывается один раз.
            </p>
            <div className="card-flat space-y-2 p-3 font-mono text-sm">
              <p>
                <span className="text-faint">Логин: </span>
                <span className="font-bold text-fg">{created.user.email}</span>
              </p>
              <p>
                <span className="text-faint">Пароль: </span>
                <span className="font-bold text-accent">{created.tempPassword}</span>
              </p>
            </div>
            {created.mailSent && (
              <p className="rounded-[var(--radius-md)] bg-[var(--success-soft)] px-3 py-2 text-success">
                Доступы отправлены на почту студента
              </p>
            )}
            {created.mailError && (
              <p className="rounded-[var(--radius-md)] bg-[var(--warning-soft)] px-3 py-2 text-warning">
                Письмо не ушло: {created.mailError}
              </p>
            )}
          </div>
        )}
      </Modal>
    </>
  );
}

function CreateStudentModal({
  open,
  onClose,
  onCreated,
}: {
  open: boolean;
  onClose: () => void;
  onCreated: (result: CreatedStudent) => void;
}) {
  const [email, setEmail] = useState("");
  const [fullName, setFullName] = useState("");
  const [password, setPassword] = useState("");
  const [role, setRole] = useState<Role>("student");
  const [sendMail, setSendMail] = useState(true);
  const [error, setError] = useState("");

  const [createUser, { isLoading }] = useCreateUserMutation();

  const submit = async (event: FormEvent) => {
    event.preventDefault();
    setError("");

    try {
      const result = await createUser({ email, fullName, password, role, sendMail }).unwrap();
      setEmail("");
      setFullName("");
      setPassword("");
      onCreated(result);
    } catch (err) {
      setError(apiErrorMessage(err, "Не удалось создать аккаунт"));
    }
  };

  return (
    <Modal
      open={open}
      onClose={onClose}
      title="Новый аккаунт"
      footer={
        <>
          <Button onClick={onClose}>Отмена</Button>
          <Button variant="primary" form="create-student" type="submit" loading={isLoading}>
            Создать
          </Button>
        </>
      }
    >
      <form id="create-student" onSubmit={submit} className="space-y-4">
        <Field label="Почта" hint="Реальный адрес — на него уйдут доступы">
          <Input
            type="email"
            value={email}
            onChange={(e) => setEmail(e.target.value)}
            placeholder="student@example.com"
            required
          />
        </Field>

        <Field label="Имя и фамилия">
          <Input
            value={fullName}
            onChange={(e) => setFullName(e.target.value)}
            placeholder="Иван Иванов"
            required
          />
        </Field>

        <Field label="Пароль" hint="Оставьте пустым — сгенерируем автоматически">
          <Input
            value={password}
            onChange={(e) => setPassword(e.target.value)}
            placeholder="Сгенерировать автоматически"
          />
        </Field>

        <Field label="Роль">
          <Select value={role} onChange={(e) => setRole(e.target.value as Role)}>
            <option value="student">Студент</option>
            <option value="admin">Администратор</option>
          </Select>
        </Field>

        <label className="flex cursor-pointer items-center gap-2 text-sm text-muted">
          <input
            type="checkbox"
            checked={sendMail}
            onChange={(e) => setSendMail(e.target.checked)}
            className="h-4 w-4 accent-[var(--accent)]"
          />
          Отправить доступы на почту
        </label>

        {error && (
          <p className="rounded-[var(--radius-md)] bg-[var(--danger-soft)] px-3 py-2 text-sm text-danger">
            {error}
          </p>
        )}
      </form>
    </Modal>
  );
}
