import { useState, type FormEvent } from "react";
import { Link, useNavigate } from "react-router-dom";

import { useAppDispatch } from "@/app/store";
import { useRegisterMutation, useSendCodeMutation } from "@/features/auth/api/authApi";
import { sessionStarted } from "@/features/auth/authSlice";
import { apiErrorMessage } from "@/shared/api/baseApi";
import { Button, Field, Input } from "@/shared/ui";
import { useToast } from "@/shared/ui/ToastProvider";

import AuthLayout from "./AuthLayout";

// Регистрация в два шага: сначала код на реальную почту, затем пароль и имя.
export default function RegisterPage() {
  const [step, setStep] = useState<"email" | "confirm">("email");
  const [email, setEmail] = useState("");
  const [fullName, setFullName] = useState("");
  const [password, setPassword] = useState("");
  const [code, setCode] = useState("");
  const [error, setError] = useState("");

  const [sendCode, { isLoading: sending }] = useSendCodeMutation();
  const [register, { isLoading: registering }] = useRegisterMutation();
  const dispatch = useAppDispatch();
  const navigate = useNavigate();
  const toast = useToast();

  const requestCode = async (event: FormEvent) => {
    event.preventDefault();
    setError("");

    try {
      await sendCode({ email, purpose: "registration" }).unwrap();
      toast.success(`Код отправлен на ${email}`);
      setStep("confirm");
    } catch (err) {
      setError(apiErrorMessage(err, "Не удалось отправить код"));
    }
  };

  const submit = async (event: FormEvent) => {
    event.preventDefault();
    setError("");

    try {
      const session = await register({ email, fullName, password, code }).unwrap();
      dispatch(sessionStarted(session));
      navigate("/learn", { replace: true });
    } catch (err) {
      setError(apiErrorMessage(err, "Не удалось завершить регистрацию"));
    }
  };

  return (
    <AuthLayout
      title="Регистрация"
      subtitle={
        step === "email"
          ? "Укажите настоящую почту — на неё придёт код подтверждения"
          : `Код отправлен на ${email}. Введите его и придумайте пароль`
      }
      footer={
        <>
          Уже есть аккаунт?{" "}
          <Link to="/login" className="font-semibold text-accent hover:underline">
            Войти
          </Link>
        </>
      }
    >
      {step === "email" ? (
        <form onSubmit={requestCode} className="space-y-4">
          <Field label="Почта" hint="На этот адрес придёт шестизначный код">
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
            Получить код
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

          <Field label="Имя и фамилия">
            <Input
              value={fullName}
              onChange={(e) => setFullName(e.target.value)}
              placeholder="Иван Иванов"
              autoComplete="name"
              required
            />
          </Field>

          <Field label="Пароль" hint="Минимум 8 символов, буквы и цифры">
            <Input
              type="password"
              value={password}
              onChange={(e) => setPassword(e.target.value)}
              placeholder="••••••••"
              autoComplete="new-password"
              required
            />
          </Field>

          {error && (
            <p className="rounded-[var(--radius-md)] bg-[var(--danger-soft)] px-3 py-2 text-sm text-danger">
              {error}
            </p>
          )}

          <Button type="submit" variant="primary" className="w-full" loading={registering}>
            Создать аккаунт
          </Button>

          <button
            type="button"
            className="w-full text-sm text-muted hover:text-accent"
            onClick={() => setStep("email")}
          >
            Изменить почту
          </button>
        </form>
      )}
    </AuthLayout>
  );
}
