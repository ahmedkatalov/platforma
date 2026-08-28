// Учебный эмулятор командной строки: виртуальная файловая система и набор
// команд, которых хватает для практики по Linux, Docker и Kubernetes.
// Ничего не выполняется по-настоящему — вывод имитируется.

export type FsFile = { type: "file"; content: string; mode?: string };
export type FsDir = { type: "dir"; children: Record<string, FsNode>; mode?: string };
export type FsNode = FsFile | FsDir;

export type ShellState = {
  cwd: string[];
  root: FsDir;
  user: string;
  host: string;
};

const APP_LOG = [
  "2026-08-27T09:12:04Z INFO  starting api version=1.4.2",
  "2026-08-27T09:12:05Z INFO  connected to postgres host=db",
  "2026-08-27T09:14:41Z WARN  slow query duration=1.8s query=SELECT * FROM orders",
  "2026-08-27T09:15:02Z ERROR failed to reach payment service: timeout after 5s",
  "2026-08-27T09:15:03Z INFO  retrying request attempt=2",
  "2026-08-27T09:15:09Z ERROR payment service unavailable, falling back to queue",
  "2026-08-27T09:16:00Z INFO  request completed status=200 path=/api/orders",
].join("\n");

const NGINX_CONF = [
  "server {",
  "    listen 80;",
  "    server_name app.example.com;",
  "",
  "    location / {",
  "        proxy_pass http://app:8080;",
  "        proxy_set_header Host $host;",
  "    }",
  "}",
].join("\n");

const COMPOSE = [
  "services:",
  "  app:",
  "    build: .",
  "    ports:",
  '      - "8080:8080"',
  "  db:",
  "    image: postgres:16",
].join("\n");

function file(content: string, mode = "-rw-r--r--"): FsFile {
  return { type: "file", content, mode };
}

function dir(children: Record<string, FsNode>): FsDir {
  return { type: "dir", children, mode: "drwxr-xr-x" };
}

// Стартовое дерево каталогов учебного сервера.
export function defaultFs(): FsDir {
  return dir({
    home: dir({
      student: dir({
        "app.log": file(APP_LOG),
        "deploy.sh": file("#!/bin/bash\necho \"deploying...\"\n", "-rw-r--r--"),
        "docker-compose.yml": file(COMPOSE),
        "notes.txt": file("Заметки к практике DevOps.\n"),
        projects: dir({
          api: dir({
            Dockerfile: file("FROM alpine:3.20\nCMD [\"/app\"]\n"),
          }),
        }),
      }),
    }),
    etc: dir({
      hosts: file("127.0.0.1 localhost\n10.0.0.5 db\n"),
      nginx: dir({ "nginx.conf": file(NGINX_CONF) }),
    }),
    var: dir({
      log: dir({
        "app.log": file(APP_LOG),
        "syslog": file("Aug 27 09:00:01 devops systemd[1]: Started nginx.\n"),
      }),
    }),
    tmp: dir({}),
    usr: dir({ bin: dir({}) }),
  });
}

export function createShell(user = "student", host = "devops"): ShellState {
  return { cwd: ["home", "student"], root: defaultFs(), user, host };
}

export function cwdPath(state: ShellState): string {
  const path = "/" + state.cwd.join("/");
  const home = "/home/" + state.user;
  if (path === home) return "~";
  if (path.startsWith(home + "/")) return "~" + path.slice(home.length);
  return path;
}

export function prompt(state: ShellState): string {
  return `${state.user}@${state.host}:${cwdPath(state)}$`;
}

// Приводит путь к массиву сегментов относительно корня.
function resolve(state: ShellState, raw: string): string[] | null {
  let path = raw.trim();
  if (path === "~") path = `/home/${state.user}`;
  else if (path.startsWith("~/")) path = `/home/${state.user}/${path.slice(2)}`;

  const segments = path.startsWith("/") ? [] : [...state.cwd];
  for (const part of path.split("/")) {
    if (part === "" || part === ".") continue;
    if (part === "..") {
      if (segments.length > 0) segments.pop();
      continue;
    }
    segments.push(part);
  }
  return segments;
}

function nodeAt(state: ShellState, segments: string[]): FsNode | null {
  let node: FsNode = state.root;
  for (const segment of segments) {
    if (node.type !== "dir") return null;
    const next: FsNode | undefined = node.children[segment];
    if (!next) return null;
    node = next;
  }
  return node;
}

function parentOf(state: ShellState, segments: string[]): { parent: FsDir; name: string } | null {
  if (segments.length === 0) return null;
  const name = segments[segments.length - 1];
  const parent = nodeAt(state, segments.slice(0, -1));
  if (!parent || parent.type !== "dir") return null;
  return { parent, name };
}

export type CommandResult = {
  output: string;
  state: ShellState;
  clear?: boolean;
};

const DOCKER_PS = [
  "CONTAINER ID   IMAGE          COMMAND                  STATUS         PORTS                  NAMES",
  "3f1c9a4b7e21   nginx:1.27     \"/docker-entrypoint.…\"   Up 12 minutes  0.0.0.0:8080->80/tcp   web",
  "9b2d8c5f0a13   postgres:16    \"docker-entrypoint.s…\"   Up 12 minutes  5432/tcp               db",
].join("\n");

const DOCKER_PS_ALL = [
  DOCKER_PS,
  "c7e4a1d9b842   redis:7        \"docker-entrypoint.s…\"   Exited (0) 2 hours ago                 cache",
].join("\n");

const DOCKER_IMAGES = [
  "REPOSITORY   TAG       IMAGE ID       CREATED        SIZE",
  "nginx        1.27      5f0b1a2c3d4e   2 weeks ago    187MB",
  "postgres     16        7a8b9c0d1e2f   3 weeks ago    432MB",
  "app          latest    1a2b3c4d5e6f   10 minutes ago 24MB",
].join("\n");

const KUBECTL_PODS = [
  "NAME                   READY   STATUS             RESTARTS   AGE",
  "api-7d9f               0/1     CrashLoopBackOff   5          6m",
  "api-8k2m               1/1     Running            0          2d",
  "postgres-0             1/1     Running            0          14d",
].join("\n");

const KUBECTL_DESCRIBE = [
  "Name:             api-7d9f",
  "Namespace:        default",
  "Status:           Running",
  "Containers:",
  "  api:",
  "    Image:        registry.example.com/api:1.4.3",
  "    State:        Waiting (CrashLoopBackOff)",
  "    Last State:   Terminated (Error), exit code: 1",
  "Events:",
  "  Warning  BackOff  2m (x12)  kubelet  Back-off restarting failed container",
  "  Normal   Pulled   6m        kubelet  Successfully pulled image",
].join("\n");

const POD_LOGS = [
  "2026-08-27T09:20:01Z INFO  starting api version=1.4.3",
  "2026-08-27T09:20:01Z FATAL config error: DATABASE_URL is required",
].join("\n");

// Имитация вывода docker.
function runDocker(args: string[]): string {
  const [sub, ...rest] = args;
  switch (sub) {
    case "ps":
      return rest.includes("-a") || rest.includes("--all") ? DOCKER_PS_ALL : DOCKER_PS;
    case "images":
      return DOCKER_IMAGES;
    case "pull": {
      const image = rest[0] ?? "nginx:latest";
      return [
        `${image.split(":")[1] ?? "latest"}: Pulling from library/${image.split(":")[0]}`,
        "a2abf6c4d29d: Pull complete",
        "c1b9a2c1e5f7: Pull complete",
        `Status: Downloaded newer image for ${image}`,
      ].join("\n");
    }
    case "run":
      return "b4d81f0c93a7e2f5c1a06d9e4c7b3128f5a0d6e29c4b8137a5e0f2d9c6b47a13";
    case "logs":
      return APP_LOG;
    case "exec":
      return "/ # (вы внутри контейнера, для выхода наберите exit)";
    case "stop":
    case "start":
    case "rm":
      return rest[rest.length - 1] ?? "";
    case "build":
      return [
        "[+] Building 12.4s (11/11) FINISHED",
        " => exporting to image",
        " => => naming to docker.io/library/app:latest",
      ].join("\n");
    case "compose":
      return rest[0] === "up"
        ? " ✔ Container db   Started\n ✔ Container app  Started"
        : " ✔ Container app  Removed\n ✔ Container db   Removed";
    default:
      return `docker: '${sub ?? ""}' is not a docker command.\nSee 'docker --help'`;
  }
}

// Имитация вывода kubectl.
function runKubectl(args: string[]): string {
  const [sub, ...rest] = args;
  switch (sub) {
    case "get":
      if (rest[0]?.startsWith("po")) return KUBECTL_PODS;
      if (rest[0]?.startsWith("deploy")) {
        return [
          "NAME   READY   UP-TO-DATE   AVAILABLE   AGE",
          "api    1/2     2            1           14d",
        ].join("\n");
      }
      if (rest[0]?.startsWith("svc") || rest[0]?.startsWith("service")) {
        return [
          "NAME   TYPE        CLUSTER-IP     PORT(S)    AGE",
          "api    ClusterIP   10.96.140.21   8080/TCP   14d",
        ].join("\n");
      }
      return KUBECTL_PODS;
    case "describe":
      return KUBECTL_DESCRIBE;
    case "logs":
      return POD_LOGS;
    case "rollout":
      if (rest[0] === "undo") return "deployment.apps/api rolled back";
      if (rest[0] === "status") return 'Waiting for deployment "api" rollout to finish: 1 of 2 updated replicas are available...';
      if (rest[0] === "restart") return "deployment.apps/api restarted";
      return "";
    case "scale":
      return "deployment.apps/api scaled";
    case "apply":
      return "deployment.apps/api configured";
    case "delete":
      return `${rest.join(" ")} deleted`;
    default:
      return `error: unknown command "${sub ?? ""}" for "kubectl"`;
  }
}

const HELP = [
  "Доступные команды учебного терминала:",
  "  pwd, ls, cd, cat, head, tail, grep, find, echo",
  "  mkdir, touch, rm, cp, mv, chmod, chown",
  "  ps, whoami, uname, df, history, clear",
  "  docker …, kubectl …, git …",
  "",
  "Подсказка: команда задания проверяется на сервере платформы.",
].join("\n");

// Разбор строки с учётом кавычек.
export function tokenize(input: string): string[] {
  const tokens: string[] = [];
  let current = "";
  let quote: '"' | "'" | null = null;

  for (const char of input) {
    if (quote) {
      if (char === quote) quote = null;
      else current += char;
      continue;
    }
    if (char === '"' || char === "'") {
      quote = char;
      continue;
    }
    if (char === " ") {
      if (current) tokens.push(current);
      current = "";
      continue;
    }
    current += char;
  }
  if (current) tokens.push(current);
  return tokens;
}

function listDir(node: FsDir, long: boolean, all: boolean): string {
  const names = Object.keys(node.children).sort();
  const visible = all ? [".", "..", ...names] : names.filter((n) => !n.startsWith("."));

  if (!long) return visible.join("  ");

  const lines = visible.map((name) => {
    if (name === "." || name === "..") return `drwxr-xr-x  2 student student  4096 авг 27 09:00 ${name}`;
    const child = node.children[name];
    const mode = child.mode ?? (child.type === "dir" ? "drwxr-xr-x" : "-rw-r--r--");
    const size = child.type === "file" ? child.content.length : 4096;
    return `${mode}  1 student student ${String(size).padStart(5)} авг 27 09:00 ${name}`;
  });

  return [`итого ${visible.length}`, ...lines].join("\n");
}

// execute выполняет одну команду и возвращает новый стейт и текст вывода.
export function execute(state: ShellState, input: string): CommandResult {
  const line = input.trim();
  if (!line) return { output: "", state };

  // Поддержка простого конвейера: команда | grep что-то
  const parts = line.split("|").map((p) => p.trim());
  let result = executeSingle(state, parts[0]);

  for (const stage of parts.slice(1)) {
    const tokens = tokenize(stage);
    if (tokens[0] === "grep") {
      const flags = tokens.filter((t) => t.startsWith("-"));
      const pattern = tokens.find((t, i) => i > 0 && !t.startsWith("-")) ?? "";
      const insensitive = flags.some((f) => f.includes("i"));
      const needle = insensitive ? pattern.toLowerCase() : pattern;
      result = {
        ...result,
        output: result.output
          .split("\n")
          .filter((row) => (insensitive ? row.toLowerCase() : row).includes(needle))
          .join("\n"),
      };
    } else if (tokens[0] === "wc") {
      const count = result.output ? result.output.split("\n").length : 0;
      result = { ...result, output: String(count) };
    } else if (tokens[0] === "head" || tokens[0] === "tail") {
      const n = Number(tokens[tokens.indexOf("-n") + 1]) || 10;
      const rows = result.output.split("\n");
      result = {
        ...result,
        output: (tokens[0] === "head" ? rows.slice(0, n) : rows.slice(-n)).join("\n"),
      };
    }
  }

  return result;
}

function executeSingle(state: ShellState, input: string): CommandResult {
  const tokens = tokenize(input);
  const [command, ...args] = tokens;
  const flags = args.filter((a) => a.startsWith("-"));
  const operands = args.filter((a) => !a.startsWith("-"));
  const hasFlag = (letter: string) => flags.some((f) => !f.startsWith("--") && f.includes(letter));

  switch (command) {
    case "pwd":
      return { output: "/" + state.cwd.join("/"), state };

    case "whoami":
      return { output: state.user, state };

    case "uname":
      return { output: hasFlag("a") ? "Linux devops 6.8.0-45-generic x86_64 GNU/Linux" : "Linux", state };

    case "clear":
      return { output: "", state, clear: true };

    case "help":
      return { output: HELP, state };

    case "ls": {
      const target = operands[0] ?? ".";
      const segments = resolve(state, target);
      const node = segments && nodeAt(state, segments);
      if (!node) return { output: `ls: невозможно получить доступ к '${target}': нет такого файла или каталога`, state };
      if (node.type === "file") return { output: target, state };
      return { output: listDir(node, hasFlag("l"), hasFlag("a")), state };
    }

    case "cd": {
      const target = operands[0] ?? `/home/${state.user}`;
      const segments = resolve(state, target);
      const node = segments && nodeAt(state, segments);
      if (!node) return { output: `cd: нет такого файла или каталога: ${target}`, state };
      if (node.type !== "dir") return { output: `cd: ${target}: не является каталогом`, state };
      return { output: "", state: { ...state, cwd: segments! } };
    }

    case "cat": {
      if (operands.length === 0) return { output: "cat: не указан файл", state };
      const outputs = operands.map((target) => {
        const segments = resolve(state, target);
        const node = segments && nodeAt(state, segments);
        if (!node) return `cat: ${target}: нет такого файла или каталога`;
        if (node.type !== "file") return `cat: ${target}: это каталог`;
        return node.content.trimEnd();
      });
      return { output: outputs.join("\n"), state };
    }

    case "head":
    case "tail": {
      const target = operands[0];
      if (!target) return { output: `${command}: не указан файл`, state };
      const segments = resolve(state, target);
      const node = segments && nodeAt(state, segments);
      if (!node || node.type !== "file") {
        return { output: `${command}: ${target}: нет такого файла`, state };
      }
      const nIndex = args.indexOf("-n");
      const count = nIndex >= 0 ? Number(args[nIndex + 1]) || 10 : 10;
      const rows = node.content.trimEnd().split("\n");
      const slice = command === "head" ? rows.slice(0, count) : rows.slice(-count);
      const follow = hasFlag("f") ? "\n(следим за файлом… Ctrl+C для выхода)" : "";
      return { output: slice.join("\n") + follow, state };
    }

    case "grep": {
      const pattern = operands[0] ?? "";
      const target = operands[1];
      const insensitive = hasFlag("i");
      if (!target) return { output: "grep: не указан файл", state };

      const segments = resolve(state, target);
      const node = segments && nodeAt(state, segments);
      if (!node || node.type !== "file") {
        return { output: `grep: ${target}: нет такого файла или каталога`, state };
      }
      const needle = insensitive ? pattern.toLowerCase() : pattern;
      const rows = node.content
        .split("\n")
        .filter((row) => (insensitive ? row.toLowerCase() : row).includes(needle));
      return { output: rows.join("\n"), state };
    }

    case "find": {
      const base = operands[0] ?? ".";
      const nameIndex = args.indexOf("-name");
      const mask = nameIndex >= 0 ? args[nameIndex + 1]?.replace(/\*/g, "") ?? "" : "";
      const segments = resolve(state, base);
      const node = segments && nodeAt(state, segments);
      if (!node) return { output: `find: '${base}': нет такого файла или каталога`, state };

      const found: string[] = [];
      const walk = (current: FsNode, path: string) => {
        if (current.type === "file") {
          if (!mask || path.includes(mask)) found.push(path);
          return;
        }
        if (!mask || path.includes(mask)) found.push(path);
        for (const [name, child] of Object.entries(current.children)) {
          walk(child, `${path}/${name}`);
        }
      };
      walk(node, base === "." ? "." : base);
      return { output: found.join("\n"), state };
    }

    case "echo":
      return { output: operands.join(" "), state };

    case "mkdir": {
      const target = operands[0];
      if (!target) return { output: "mkdir: пропущен операнд", state };
      const segments = resolve(state, target);
      const place = segments && parentOf(state, segments);
      if (!place) return { output: `mkdir: невозможно создать каталог '${target}'`, state };
      place.parent.children[place.name] = dir({});
      return { output: "", state: { ...state } };
    }

    case "touch": {
      const target = operands[0];
      if (!target) return { output: "touch: пропущен операнд", state };
      const segments = resolve(state, target);
      const place = segments && parentOf(state, segments);
      if (!place) return { output: `touch: невозможно создать '${target}'`, state };
      if (!place.parent.children[place.name]) place.parent.children[place.name] = file("");
      return { output: "", state: { ...state } };
    }

    case "rm": {
      const target = operands[0];
      if (!target) return { output: "rm: пропущен операнд", state };
      const segments = resolve(state, target);
      const place = segments && parentOf(state, segments);
      if (!place || !place.parent.children[place.name]) {
        return { output: `rm: невозможно удалить '${target}': нет такого файла`, state };
      }
      const node = place.parent.children[place.name];
      if (node.type === "dir" && !hasFlag("r")) {
        return { output: `rm: невозможно удалить '${target}': это каталог`, state };
      }
      delete place.parent.children[place.name];
      return { output: "", state: { ...state } };
    }

    case "cp":
    case "mv": {
      const [from, to] = operands;
      if (!from || !to) return { output: `${command}: пропущен операнд`, state };
      const fromSegments = resolve(state, from);
      const node = fromSegments && nodeAt(state, fromSegments);
      if (!node) return { output: `${command}: ${from}: нет такого файла`, state };

      const toSegments = resolve(state, to);
      const place = toSegments && parentOf(state, toSegments);
      if (!place) return { output: `${command}: невозможно записать '${to}'`, state };

      place.parent.children[place.name] = JSON.parse(JSON.stringify(node)) as FsNode;
      if (command === "mv") {
        const source = parentOf(state, fromSegments!);
        if (source) delete source.parent.children[source.name];
      }
      return { output: "", state: { ...state } };
    }

    case "chmod": {
      const [mode, target] = operands;
      if (!mode || !target) return { output: "chmod: пропущен операнд", state };
      const segments = resolve(state, target);
      const node = segments && nodeAt(state, segments);
      if (!node) return { output: `chmod: невозможно получить доступ к '${target}'`, state };
      node.mode = mode.includes("x") || mode.startsWith("7") || mode.startsWith("75")
        ? "-rwxr-xr-x"
        : "-rw-r--r--";
      return { output: "", state: { ...state } };
    }

    case "chown":
      return { output: "", state };

    case "ps":
      return {
        output: [
          "USER       PID %CPU %MEM    VSZ   RSS COMMAND",
          "root         1  0.0  0.1 168920 11284 /sbin/init",
          "root       412  0.0  0.3 720104 25640 nginx: master process",
          "www-data   418  0.1  0.2 720640 18932 nginx: worker process",
          "student   1043  0.2  0.4 913204 33108 /usr/bin/app",
        ].join("\n"),
        state,
      };

    case "df":
      return {
        output: [
          "Файловая система Размер Использовано Дост Использовано% Смонтировано в",
          "/dev/root           40G          18G  21G           47% /",
          "tmpfs              2,0G            0 2,0G            0% /dev/shm",
        ].join("\n"),
        state,
      };

    case "systemctl":
      return {
        output: operands[0] === "status"
          ? "● nginx.service - A high performance web server\n   Active: active (running) since Thu 2026-08-27 09:00:01 UTC"
          : "",
        state,
      };

    case "journalctl":
      return { output: APP_LOG, state };

    case "docker":
      return { output: runDocker(args), state };

    case "kubectl":
      return { output: runKubectl(args), state };

    case "git": {
      if (args[0] === "status") {
        return {
          output: "On branch main\nnothing to commit, working tree clean",
          state,
        };
      }
      if (args[0] === "log") {
        return {
          output: "c3f9a21 (HEAD -> main) feat: add healthcheck\n9d1b70e chore: bump deps",
          state,
        };
      }
      return { output: "", state };
    }

    case "exit":
      return { output: "выход из контейнера", state };

    default:
      return { output: `${command}: команда не найдена. Наберите help для списка команд.`, state };
  }
}
