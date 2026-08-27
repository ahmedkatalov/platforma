import { useState, type FormEvent } from "react";
import { Link, useNavigate } from "react-router-dom";

import { useResetPasswordMutation, useSendCodeMutation } from "@/features/auth/api/authApi";
import { apiErrorMessage } from "@/shared/api/baseApi";
import { Button, Field, Input } from "@/shared/ui";
import { useToast } from "@/shared/ui/ToastProvider";

import AuthLayout from "./AuthLayout";

export default function ForgotPasswordPage() {
  const [step, setStep] = useState<"email" | "confirm">("email");
  const [email, setEmail] = useState("");
  const [code, setCode] = useState("");
  const [password, setPassword] = useState("");
  const [error, setError] = useState("");

  const [sendCode, { isLoading: sending }] = useSendCodeMutation();
  const [resetPassword, { isLoading: resetting }] = useResetPasswordMutation();
  const navigate = useNavigate();
  const toast = useToast();

  const requestCode = async (event: FormEvent) => {
    event.preventDefault();
    setError("");

    try {
      await sendCode({ email, purpose: "password_reset" }).unwrap();
      toast.success("Если аккаунт существует, код отправлен на почту");
      setStep("confirm");
    } catch (err) {
      setError(apiErrorMessage(err, "Не удалось отправить код"));
    }
  };

  const submit = async (event: FormEvent) => {
    event.preventDefault();
    setError("");

    try {
      await resetPassword({ email, code, password }).unwrap();
      toast.success("Пароль обновлён, войдите с новым паролем");
      navigate("/login", { replace: true });
    } catch (err) {
      setError(apiErrorMessage(err, "Не удалось сменить пароль"));
    }
  };

  return (
    <AuthLayout
      title="Восстановление пароля"
      subtitle={
        step === "email"
          ? "Отправим код на вашу почту"
          : "Введите код из письма и новый пароль"
      }
      footer={
        <Link to="/login" className="font-semibold text-accent hover:underline">
          Вернуться ко входу
        </Link>
      }
    >
      {step === "email" ? (
        <form onSubmit={requestCode} className="space-y-4">
          <Field label="Почта">
            <Input
              type="email"
              value={email}
              onChange={(e) => setEmail(e.target.value)}
              placeholder="you@example.com"
              autoComplete="email"
              required
            />
          </Field>

          {error && (
            <p className="rounded-[var(--radius-md)] bg-[var(--danger-soft)] px-3 py-2 text-sm text-danger">
              {error}
            </p>
          )}

          <Button type="submit" variant="primary" className="w-full" loading={sending}>
            Отправить код
          </Button>
        </form>
      ) : (
        <form onSubmit={submit} className="space-y-4">
          <Field label="Код из письма">
            <Input
              value={code}
              onChange={(e) => setCode(e.target.value.replace(/\D/g, "").slice(0, 6))}
              placeholder="000000"
              inputMode="numeric"
              className="text-center text-lg font-bold tracking-[0.4em]"
              required
            />
          </Field>

          <Field label="Новый пароль" hint="Минимум 8 символов, буквы и цифры">
            <Input
              type="password"
              value={password}
              onChange={(e) => setPassword(e.target.value)}
              autoComplete="new-password"
              required
            />
          </Field>

          {error && (
            <p className="rounded-[var(--radius-md)] bg-[var(--danger-soft)] px-3 py-2 text-sm text-danger">
              {error}
            </p>
          )}

          <Button type="submit" variant="primary" className="w-full" loading={resetting}>
            Сменить пароль
          </Button>
        </form>
      )}
    </AuthLayout>
  );
}
