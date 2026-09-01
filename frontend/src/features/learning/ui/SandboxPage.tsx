import { useCallback, useEffect, useRef, useState } from "react";
import { Terminal } from "@xterm/xterm";
import { FitAddon } from "@xterm/addon-fit";
import "@xterm/xterm/css/xterm.css";
import {
  BookOpen,
  CheckCircle2,
  Circle,
  Download,
  FileUp,
  Maximize2,
  Minimize2,
  Play,
  RotateCcw,
  TerminalSquare,
  Upload,
} from "lucide-react";

import { Badge, Button, Card, PageHeader, Spinner } from "@/shared/ui";
import { useToast } from "@/shared/ui/ToastProvider";
import { loadV86, V86_ASSETS, type V86Ctor, type V86Instance } from "@/features/learning/lib/v86";

// Учебный проект, который автоматически появляется в /root/app после загрузки.
// Это «сервис»: фоновый процесс, который пишет логи — как настоящее приложение.
const SETUP_SCRIPT = String.raw`
ifconfig lo 127.0.0.1 up 2>/dev/null
mkdir -p /root/app
cat > /root/app/service.sh <<'SH'
#!/bin/sh
i=1
while true; do
  echo "$(date '+%Y-%m-%d %H:%M:%S') [app] tick #$i — сервис работает" >> /root/app/app.log
  i=$((i+1))
  sleep 2
done
SH
cat > /root/app/start.sh <<'SH'
#!/bin/sh
if [ -f /root/app/app.pid ] && kill -0 "$(cat /root/app/app.pid)" 2>/dev/null; then
  echo "Сервис уже запущен (PID=$(cat /root/app/app.pid))"; exit 0
fi
sh /root/app/service.sh &
echo $! > /root/app/app.pid
echo "Сервис запущен, PID=$(cat /root/app/app.pid)."
echo "Логи вживую:  tail -f /root/app/app.log   (Ctrl+C — выйти)"
SH
cat > /root/app/status.sh <<'SH'
#!/bin/sh
P=$(cat /root/app/app.pid 2>/dev/null)
if [ -n "$P" ] && kill -0 "$P" 2>/dev/null; then echo "* Сервис РАБОТАЕТ (PID=$P)"; else echo "o Сервис остановлен"; fi
echo "--- последние строки лога ---"
tail -n 5 /root/app/app.log 2>/dev/null
SH
cat > /root/app/stop.sh <<'SH'
#!/bin/sh
P=$(cat /root/app/app.pid 2>/dev/null)
if [ -n "$P" ] && kill "$P" 2>/dev/null; then echo "Сервис остановлен (PID=$P)"; rm -f /root/app/app.pid; else echo "Сервис не запущен"; fi
SH
chmod +x /root/app/*.sh
cat > /root/app/access.log <<'LOG'
10.0.0.1 - - [10/Oct/2024:13:55:36] "GET /index.html HTTP/1.1" 200 1043
10.0.0.2 - - [10/Oct/2024:13:55:37] "GET /style.css HTTP/1.1" 200 233
10.0.0.1 - - [10/Oct/2024:13:55:39] "GET /api/users HTTP/1.1" 500 51
10.0.0.3 - - [10/Oct/2024:13:55:40] "GET /missing HTTP/1.1" 404 21
10.0.0.2 - - [10/Oct/2024:13:55:41] "POST /api/login HTTP/1.1" 200 88
10.0.0.1 - - [10/Oct/2024:13:55:42] "GET /api/users HTTP/1.1" 500 51
10.0.0.4 - - [10/Oct/2024:13:55:44] "GET /index.html HTTP/1.1" 200 1043
LOG
cat > /root/README <<'TXT'
Песочница Okvion — настоящий Linux (busybox) прямо в браузере.

Учебный сервис (имитация работающего приложения):
  Запустить:   sh /root/app/start.sh
  Статус+лог:  sh /root/app/status.sh
  Логи вживую: tail -f /root/app/app.log     (Ctrl+C выйти)
  Остановить:  sh /root/app/stop.sh

Анализ логов (пример /root/app/access.log):
  grep " 500 " /root/app/access.log
  awk '{print $1}' /root/app/access.log | sort | uniq -c | sort -rn

Свобода действий: vi, ps, top (q выход), ip a, df -h, find, tar ...
Ничего не сломаешь: кнопка "Перезапустить" вернёт чистую систему.
TXT
clear
echo '=================================================='
echo ' Песочница Okvion готова. Проект: /root/app'
echo '=================================================='
cat /root/README
`;

// v86 отправляет в порт по одному байту на символ (charCodeAt & 0xFF),
// поэтому кириллицу нужно заранее закодировать в UTF-8 и отдать как байты.
function toUtf8Bytes(text: string): string {
  const bytes = new TextEncoder().encode(text);
  let out = "";
  for (let i = 0; i < bytes.length; i++) out += String.fromCharCode(bytes[i]);
  return out;
}

type Task = { id: string; text: string };

const TASKS: Task[] = [
  { id: "look", text: "Осмотрись: sh /root/app/status.sh и cat /root/README" },
  { id: "start", text: "Запусти сервис: sh /root/app/start.sh" },
  { id: "logs", text: "Смотри логи вживую: tail -f /root/app/app.log (Ctrl+C — выйти)" },
  { id: "grep", text: "Найди ошибки 500: grep ' 500 ' /root/app/access.log" },
  { id: "awk", text: "Топ IP: awk '{print $1}' /root/app/access.log | sort | uniq -c | sort -rn" },
  { id: "vi", text: "В vi /root/app/service.sh поменяй sleep 2 на sleep 1 и перезапусти сервис" },
  { id: "stop", text: "Останови сервис: sh /root/app/stop.sh, проверь статус" },
  { id: "own", text: "Напиши свой скрипт в /root и запусти его" },
  { id: "explore", text: "Изучи систему: ps, top (q), df -h, ip a, uname -a" },
];

const TASKS_KEY = "okvion.sandbox.tasks";

const TERM_THEME = {
  background: "#0b0f17",
  foreground: "#d5dae5",
  cursor: "#6ee7b7",
  selectionBackground: "#334155",
  black: "#1e293b",
  brightBlack: "#475569",
};

export default function SandboxPage() {
  const toast = useToast();
  const [status, setStatus] = useState<"idle" | "loading" | "running" | "error">("idle");
  const [errorMsg, setErrorMsg] = useState("");
  const [fullscreen, setFullscreen] = useState(false);
  const [done, setDone] = useState<Set<string>>(() => {
    try {
      return new Set<string>(JSON.parse(localStorage.getItem(TASKS_KEY) || "[]") as string[]);
    } catch {
      return new Set<string>();
    }
  });

  const termHostRef = useRef<HTMLDivElement>(null);
  const termRef = useRef<Terminal | null>(null);
  const fitRef = useRef<FitAddon | null>(null);
  const emuRef = useRef<V86Instance | null>(null);
  const injectedRef = useRef(false);
  const fileInputRef = useRef<HTMLInputElement>(null);
  const stateInputRef = useRef<HTMLInputElement>(null);

  const toggleTask = (id: string) => {
    setDone((prev) => {
      const next = new Set(prev);
      next.has(id) ? next.delete(id) : next.add(id);
      localStorage.setItem(TASKS_KEY, JSON.stringify([...next]));
      return next;
    });
  };

  // Отправить строку команд в терминал (как будто студент их набрал).
  const send = useCallback((text: string) => {
    emuRef.current?.serial0_send(toUtf8Bytes(text));
  }, []);

  const injectProject = useCallback(() => {
    if (!emuRef.current) return;
    injectedRef.current = true;
    send(SETUP_SCRIPT.trim() + "\n");
    toast.success("Учебный проект загружен в /root/app");
  }, [send, toast]);

  // Полностью останавливает и очищает эмулятор и терминал.
  const teardown = useCallback(async () => {
    try {
      await emuRef.current?.destroy?.();
    } catch {
      /* ignore */
    }
    emuRef.current = null;
    termRef.current?.dispose();
    termRef.current = null;
    fitRef.current = null;
    injectedRef.current = false;
  }, []);

  const boot = useCallback(
    async (initialState?: ArrayBuffer) => {
      if (!termHostRef.current) return;
      setStatus("loading");
      setErrorMsg("");
      await teardown();

      let V86: V86Ctor;
      try {
        V86 = await loadV86();
      } catch (err) {
        setStatus("error");
        setErrorMsg(err instanceof Error ? err.message : "Ошибка загрузки движка");
        return;
      }

      // Терминал.
      const term = new Terminal({
        convertEol: false,
        cursorBlink: true,
        fontSize: 14,
        fontFamily: '"JetBrains Mono", "Fira Code", ui-monospace, SFMono-Regular, Menlo, monospace',
        theme: TERM_THEME,
        scrollback: 5000,
      });
      const fit = new FitAddon();
      term.loadAddon(fit);
      term.open(termHostRef.current);
      fit.fit();
      term.focus();
      termRef.current = term;
      fitRef.current = fit;

      // Буферизуем байты последовательного порта и пишем пачками — так быстрее.
      let buf: number[] = [];
      let scheduled = false;
      let sawOutput = false;
      let tail = "";
      let silence: ReturnType<typeof setTimeout> | null = null;
      const startedAt = Date.now();

      const flush = () => {
        scheduled = false;
        if (buf.length) {
          term.write(new Uint8Array(buf));
          buf = [];
        }
      };

      const options: Record<string, unknown> = {
        wasm_path: V86_ASSETS.wasm_path,
        bios: { url: V86_ASSETS.bios },
        vga_bios: { url: V86_ASSETS.vga_bios },
        cdrom: { url: V86_ASSETS.cdrom },
        autostart: true,
        memory_size: 128 * 1024 * 1024,
        vga_memory_size: 2 * 1024 * 1024,
        disable_speaker: true,
      };
      if (initialState) options.initial_state = { buffer: initialState };

      let emulator: V86Instance;
      try {
        emulator = new V86(options);
      } catch (err) {
        setStatus("error");
        setErrorMsg(err instanceof Error ? err.message : "Не удалось запустить эмулятор");
        return;
      }
      emuRef.current = emulator;

      emulator.add_listener("serial0-output-byte", (data) => {
        const byte = data as number;
        buf.push(byte);
        sawOutput = true;
        if (!scheduled) {
          scheduled = true;
          requestAnimationFrame(flush);
        }

        if (initialState || injectedRef.current) return;

        // Основной сигнал готовности — появилось приглашение оболочки (~%, #, $).
        tail = (tail + String.fromCharCode(byte)).slice(-80);
        if (Date.now() - startedAt > 3000 && /[~/\w][%#$]\s$/.test(tail)) {
          injectProject();
          return;
        }
        // Запасной вариант: вывод «успокоился» на 2 секунды после загрузки.
        if (silence) clearTimeout(silence);
        silence = setTimeout(() => {
          if (!injectedRef.current && sawOutput && Date.now() - startedAt > 4000) injectProject();
        }, 2000);
      });

      term.onData((d) => emulator.serial0_send(toUtf8Bytes(d)));
      setStatus("running");
    },
    [teardown, injectProject],
  );

  // Подгоняем размер терминала под контейнер.
  useEffect(() => {
    const onResize = () => fitRef.current?.fit();
    window.addEventListener("resize", onResize);
    return () => window.removeEventListener("resize", onResize);
  }, []);

  useEffect(() => {
    const t = setTimeout(() => fitRef.current?.fit(), 60);
    return () => clearTimeout(t);
  }, [fullscreen]);

  // Чистим за собой при уходе со страницы.
  useEffect(() => () => void teardown(), [teardown]);

  const onUploadFile = async (file: File) => {
    if (!emuRef.current) return;
    if (file.size > 256 * 1024) {
      toast.error("Файл больше 256 КБ — загрузите файл поменьше");
      return;
    }
    const text = await file.text();
    const safeName = file.name.replace(/[^A-Za-z0-9._-]/g, "_");
    // Пишем файл в /root через heredoc со случайным маркером конца.
    const marker = "EOF_OKVION_UPLOAD";
    send(`cat > /root/${safeName} <<'${marker}'\n${text}\n${marker}\n`);
    send(`echo 'Файл загружен: /root/${safeName}'\n`);
    toast.success(`Файл /root/${safeName} загружен в песочницу`);
  };

  const onSaveState = async () => {
    if (!emuRef.current) return;
    try {
      const state = await emuRef.current.save_state();
      const blob = new Blob([state], { type: "application/octet-stream" });
      const url = URL.createObjectURL(blob);
      const a = document.createElement("a");
      a.href = url;
      a.download = "okvion-sandbox-state.bin";
      a.click();
      URL.revokeObjectURL(url);
      toast.success("Состояние сохранено в файл");
    } catch {
      toast.error("Не удалось сохранить состояние");
    }
  };

  const onRestoreState = async (file: File) => {
    const buffer = await file.arrayBuffer();
    await boot(buffer);
    toast.success("Состояние восстановлено");
  };

  const doneCount = TASKS.filter((t) => done.has(t.id)).length;

  return (
    <div className={fullscreen ? "fixed inset-0 z-50 overflow-auto bg-bg p-4" : undefined}>
      {!fullscreen && (
        <PageHeader
          title="Песочница — Linux-терминал"
          subtitle="Настоящий Linux прямо в браузере: экспериментируйте свободно, ничего не сломаете"
          actions={
            <Badge tone="accent">
              <TerminalSquare size={14} /> beta
            </Badge>
          }
        />
      )}

      <div className={fullscreen ? "grid gap-4" : "grid gap-[var(--gap)] xl:grid-cols-[1fr_20rem]"}>
        {/* Терминал */}
        <Card className="flex flex-col overflow-hidden p-0">
          <div className="flex flex-wrap items-center gap-2 border-b border-line px-3 py-2">
            <span className="mr-1 flex items-center gap-1.5 text-sm font-semibold text-fg">
              <TerminalSquare size={16} className="text-accent" /> student@okvion
            </span>
            {status === "running" && <span className="h-2 w-2 rounded-full bg-[var(--success)]" title="Работает" />}

            <div className="ml-auto flex flex-wrap items-center gap-1.5">
              {status === "running" && (
                <>
                  <Button variant="ghost" className="h-8 !px-2" onClick={injectProject} title="Загрузить учебный проект заново">
                    <Upload size={15} /> Проект
                  </Button>
                  <Button variant="ghost" className="h-8 !px-2" onClick={() => fileInputRef.current?.click()} title="Загрузить свой файл в /root">
                    <FileUp size={15} /> Файл
                  </Button>
                  <Button variant="ghost" className="h-8 !px-2" onClick={onSaveState} title="Сохранить состояние в файл">
                    <Download size={15} />
                  </Button>
                  <Button variant="ghost" className="h-8 !px-2" onClick={() => boot()} title="Перезапустить">
                    <RotateCcw size={15} />
                  </Button>
                </>
              )}
              <Button
                variant="ghost"
                className="h-8 !px-2"
                onClick={() => setFullscreen((v) => !v)}
                title={fullscreen ? "Свернуть" : "На весь экран"}
              >
                {fullscreen ? <Minimize2 size={15} /> : <Maximize2 size={15} />}
              </Button>
            </div>
          </div>

          <div className="relative bg-[#0b0f17]">
            {/* Хост терминала */}
            <div
              ref={termHostRef}
              className="w-full px-2 py-1"
              style={{ height: fullscreen ? "calc(100vh - 7rem)" : "66vh" }}
            />

            {status !== "running" && (
              <div className="absolute inset-0 grid place-items-center bg-[#0b0f17]/95 p-6 text-center">
                {status === "loading" ? (
                  <div className="flex flex-col items-center gap-3 text-accent">
                    <Spinner size={32} />
                    <p className="text-sm text-muted">Загружаю Linux в браузер (~10 МБ, один раз)…</p>
                  </div>
                ) : status === "error" ? (
                  <div className="max-w-md space-y-3">
                    <p className="text-sm text-danger">{errorMsg}</p>
                    <Button variant="primary" icon={<RotateCcw size={16} />} onClick={() => boot()}>
                      Попробовать снова
                    </Button>
                  </div>
                ) : (
                  <div className="max-w-md space-y-4">
                    <TerminalSquare size={40} className="mx-auto text-accent" />
                    <div>
                      <p className="text-base font-bold text-fg">Запустить песочницу</p>
                      <p className="mt-1 text-sm text-muted">
                        Откроется настоящий Linux (busybox) прямо в браузере. Первая загрузка ~10 МБ,
                        дальше — из кэша. Делайте что угодно: файлы, скрипты, сервисы, анализ логов.
                      </p>
                    </div>
                    <Button variant="primary" icon={<Play size={16} />} onClick={() => boot()}>
                      Запустить
                    </Button>
                  </div>
                )}
              </div>
            )}
          </div>
        </Card>

        {/* Инструкция и задания */}
        {!fullscreen && (
          <div className="space-y-[var(--gap)]">
            <Card className="p-[var(--pad)]">
              <h2 className="mb-2 flex items-center gap-2 text-sm font-bold text-fg">
                <BookOpen size={16} className="text-accent" /> Как пользоваться
              </h2>
              <ul className="space-y-1.5 text-[13px] text-muted">
                <li>• Это <b className="text-fg">настоящий Linux</b>, а не имитация — команды выполняются по-настоящему.</li>
                <li>• В <code className="text-accent">/root/app</code> уже лежит учебный «сервис» — приложение, которое пишет логи.</li>
                <li>• Запустить: <code className="text-accent">sh /root/app/start.sh</code>, статус: <code className="text-accent">sh /root/app/status.sh</code></li>
                <li>• Логи вживую: <code className="text-accent">tail -f /root/app/app.log</code> (Ctrl+C — выйти).</li>
                <li>• «Файл» — загрузить свой скрипт/приложение в <code className="text-accent">/root</code>.</li>
                <li>• Сломать ничего нельзя: «Перезапустить» вернёт чистую систему.</li>
              </ul>
            </Card>

            <Card className="p-[var(--pad)]">
              <div className="mb-2 flex items-center justify-between">
                <h2 className="flex items-center gap-2 text-sm font-bold text-fg">
                  <CheckCircle2 size={16} className="text-accent" /> Задания
                </h2>
                <span className="text-xs font-bold text-accent">
                  {doneCount}/{TASKS.length}
                </span>
              </div>
              <ul className="space-y-1.5">
                {TASKS.map((task) => {
                  const isDone = done.has(task.id);
                  return (
                    <li key={task.id}>
                      <button
                        onClick={() => toggleTask(task.id)}
                        className="flex w-full items-start gap-2 rounded-[var(--radius-md)] p-1.5 text-left text-[13px] transition-colors hover:bg-surface-2"
                      >
                        {isDone ? (
                          <CheckCircle2 size={16} className="mt-0.5 shrink-0 text-success" />
                        ) : (
                          <Circle size={16} className="mt-0.5 shrink-0 text-faint" />
                        )}
                        <span className={isDone ? "text-faint line-through" : "text-muted"}>{task.text}</span>
                      </button>
                    </li>
                  );
                })}
              </ul>
              <p className="mt-2 text-[11px] text-faint">Отмечайте выполненное вручную — прогресс сохраняется на этом устройстве.</p>
            </Card>

            <Card className="p-[var(--pad)]">
              <h2 className="mb-2 text-sm font-bold text-fg">Состояние</h2>
              <p className="mb-2 text-[12px] text-muted">
                Сохраните работу в файл и восстановите позже — со всеми созданными файлами.
              </p>
              <div className="flex gap-2">
                <Button variant="secondary" className="h-8 flex-1 !px-2 text-xs" onClick={onSaveState} disabled={status !== "running"}>
                  <Download size={14} /> Сохранить
                </Button>
                <Button variant="secondary" className="h-8 flex-1 !px-2 text-xs" onClick={() => stateInputRef.current?.click()}>
                  <Upload size={14} /> Восстановить
                </Button>
              </div>
            </Card>
          </div>
        )}
      </div>

      {/* Скрытые input для загрузки файлов */}
      <input
        ref={fileInputRef}
        type="file"
        className="hidden"
        onChange={(e) => {
          const f = e.target.files?.[0];
          e.target.value = "";
          if (f) void onUploadFile(f);
        }}
      />
      <input
        ref={stateInputRef}
        type="file"
        accept=".bin"
        className="hidden"
        onChange={(e) => {
          const f = e.target.files?.[0];
          e.target.value = "";
          if (f) void onRestoreState(f);
        }}
      />
    </div>
  );
}
