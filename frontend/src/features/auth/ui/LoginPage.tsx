import { useState, type FormEvent } from "react";
import { Link, useNavigate } from "react-router-dom";

import { useAppDispatch } from "@/app/store";
import { useLoginMutation } from "@/features/auth/api/authApi";
import { sessionStarted } from "@/features/auth/authSlice";
import { apiErrorMessage } from "@/shared/api/baseApi";
import { Button, Field, Input } from "@/shared/ui";

import AuthLayout from "./AuthLayout";

export default function LoginPage() {
  const [email, setEmail] = useState("");
  const [password, setPassword] = useState("");
  const [error, setError] = useState("");

  const [login, { isLoading }] = useLoginMutation();
  const dispatch = useAppDispatch();
  const navigate = useNavigate();

  const onSubmit = async (event: FormEvent) => {
    event.preventDefault();
    setError("");

    try {
      const session = await login({ email, password }).unwrap();
      dispatch(sessionStarted(session));
      navigate(session.user.role === "admin" ? "/admin" : "/learn", { replace: true });
    } catch (err) {
      setError(apiErrorMessage(err, "Не удалось войти"));
    }
  };

  return (
    <AuthLayout
      title="Вход на платформу"
      subtitle="Введите почту и пароль, выданные администратором"
    >
      <form onSubmit={onSubmit} className="space-y-4">
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

        <Field label="Пароль">
          <Input
            type="password"
            value={password}
            onChange={(e) => setPassword(e.target.value)}
            placeholder="••••••••"
            autoComplete="current-password"
            required
          />
        </Field>

        {error && (
          <p className="rounded-[var(--radius-md)] bg-[var(--danger-soft)] px-3 py-2 text-sm text-danger">
            {error}
          </p>
        )}

        <Button type="submit" variant="primary" className="w-full" loading={isLoading}>
          Войти
        </Button>

        <div className="text-center">
          <Link to="/forgot-password" className="text-sm text-muted hover:text-accent">
            Забыли пароль?
          </Link>
        </div>
      </form>
    </AuthLayout>
  );
}
