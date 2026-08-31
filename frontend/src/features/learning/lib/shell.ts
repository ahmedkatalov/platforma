// Учебный эмулятор командной строки: виртуальная файловая система и набор
// команд, которых хватает для практики по Linux, Git, сети, Docker,
// Kubernetes, IaC и мониторингу. Ничего не выполняется по-настоящему —
// вывод имитируется.

import {
  runAnsible,
  runCurl,
  runDate,
  runDig,
  runDu,
  runFree,
  runGit,
  runHelm,
  runNginx,
  runOpenssl,
  runPing,
  runSs,
  runSshKeygen,
  runTerraform,
  runUptime,
} from "./shellTools";

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

const MAIN_TF = [
  'terraform {',
  '  required_providers {',
  '    aws = {',
  '      source  = "hashicorp/aws"',
  '      version = "~> 5.0"',
  '    }',
  '  }',
  '}',
  '',
  'resource "aws_instance" "app" {',
  '  ami           = var.ami_id',
  '  instance_type = "t3.micro"',
  '',
  '  tags = {',
  '    Name = "app-server"',
  '  }',
  '}',
].join("\n");

const PLAYBOOK = [
  "- name: Configure web servers",
  "  hosts: web",
  "  become: true",
  "  tasks:",
  "    - name: Install nginx",
  "      apt:",
  "        name: nginx",
  "        state: present",
  "",
  "    - name: Start nginx service",
  "      service:",
  "        name: nginx",
  "        state: started",
  "        enabled: true",
].join("\n");

const PROMETHEUS_YML = [
  "global:",
  "  scrape_interval: 15s",
  "",
  "scrape_configs:",
  "  - job_name: api",
  "    static_configs:",
  "      - targets: ['app:8080']",
  "",
  "rule_files:",
  "  - /etc/prometheus/alerts.yml",
].join("\n");

const ALERTS_YML = [
  "groups:",
  "  - name: api",
  "    rules:",
  "      - alert: HighErrorRate",
  "        expr: rate(http_requests_total{code=~\"5..\"}[5m]) > 0.05",
  "        for: 10m",
  "        labels:",
  "          severity: critical",
  "        annotations:",
  "          summary: Больше 5% ошибок за пять минут",
].join("\n");

const CI_WORKFLOW = [
  "name: ci",
  "",
  "on:",
  "  push:",
  "    branches: [main]",
  "",
  "jobs:",
  "  build:",
  "    runs-on: ubuntu-latest",
  "    steps:",
  "      - uses: actions/checkout@v4",
  "      - run: go test ./...",
].join("\n");

const NGINX_ACCESS = [
  '10.0.0.9 - - [27/Aug/2026:09:15:01 +0300] "GET /api/orders HTTP/1.1" 200 1832 "-" "curl/8.6.0"',
  '10.0.0.9 - - [27/Aug/2026:09:15:02 +0300] "POST /api/orders HTTP/1.1" 500 214 "-" "curl/8.6.0"',
  '10.0.0.11 - - [27/Aug/2026:09:15:04 +0300] "GET /health HTTP/1.1" 200 15 "-" "kube-probe/1.30"',
  '10.0.0.9 - - [27/Aug/2026:09:15:09 +0300] "GET /api/orders HTTP/1.1" 500 214 "-" "curl/8.6.0"',
  '10.0.0.12 - - [27/Aug/2026:09:15:12 +0300] "GET /static/app.js HTTP/1.1" 304 0 "-" "Mozilla/5.0"',
].join("\n");

const DOCKERFILE = [
  "FROM golang:1.25-alpine AS build",
  "WORKDIR /src",
  "COPY go.mod go.sum ./",
  "RUN go mod download",
  "COPY . .",
  "RUN go build -o /app ./cmd/api",
  "",
  "FROM alpine:3.20",
  "RUN adduser -D -u 10001 app",
  "USER app",
  "COPY --from=build /app /app",
  "EXPOSE 8080",
  'ENTRYPOINT ["/app"]',
].join("\n");

const PASSWD = [
  "root:x:0:0:root:/root:/bin/bash",
  "daemon:x:1:1:daemon:/usr/sbin:/usr/sbin/nologin",
  "www-data:x:33:33:www-data:/var/www:/usr/sbin/nologin",
  "postgres:x:114:120:PostgreSQL admin:/var/lib/postgresql:/bin/bash",
  "student:x:1000:1000:Student:/home/student:/bin/bash",
].join("\n");

const CPUINFO = [
  "processor\t: 0",
  "model name\t: Intel(R) Xeon(R) CPU @ 2.30GHz",
  "cpu MHz\t\t: 2300.000",
  "cores\t\t: 4",
  "",
  "processor\t: 1",
  "model name\t: Intel(R) Xeon(R) CPU @ 2.30GHz",
  "",
  "processor\t: 2",
  "model name\t: Intel(R) Xeon(R) CPU @ 2.30GHz",
  "",
  "processor\t: 3",
  "model name\t: Intel(R) Xeon(R) CPU @ 2.30GHz",
].join("\n");

const MEMINFO = [
  "MemTotal:        4030524 kB",
  "MemFree:          812340 kB",
  "MemAvailable:    2600128 kB",
  "SwapTotal:             0 kB",
].join("\n");

const CI_LOG = [
  "2026-08-27T10:00:01Z INFO  checkout main",
  "2026-08-27T10:00:12Z INFO  running: go test ./...",
  "2026-08-27T10:00:44Z Error: cart_test.go:52: want status 200, got 500",
  "2026-08-27T10:00:44Z FAIL  tests failed (41/42)",
  "2026-08-27T10:00:44Z INFO  build cancelled, prod untouched",
].join("\n");

const USERS_TXT = "anna\nbob\nanna\ncarol\nbob\nanna\ndave\ncarol\nbob\n";

const SSH_CONFIG = [
  "Host web-1",
  "    HostName 10.0.0.21",
  "    User deploy",
  "    IdentityFile ~/.ssh/id_ed25519",
].join("\n");

const TF_LOCK = [
  "# This file is maintained automatically by \"terraform init\".",
  "# Manual edits may be lost in future updates.",
  "",
  'provider "registry.terraform.io/hashicorp/aws" {',
  '  version = "5.60.0"',
  "}",
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
        "users.txt": file(USERS_TXT),
        "access.log": file(NGINX_ACCESS),
        "app.txt": file("первая строка приложения\n"),
        "ci.log": file(CI_LOG),
        ".terraform.lock.hcl": file(TF_LOCK),
        ".github": dir({
          workflows: dir({ "ci.yml": file(CI_WORKFLOW) }),
        }),
        "Hello-World": dir({
          README: file("Hello World!\n"),
          ".git": dir({ HEAD: file("ref: refs/heads/master\n") }),
        }),
        myapp: dir({
          "README.md": file("# myapp\n\nУчебный веб-проект. История коммитов — для практики git.\n"),
          "index.js": file("console.log('hello from myapp');\n"),
          "package.json": file('{\n  "name": "myapp",\n  "version": "1.0.0"\n}\n'),
        }),
        projects: dir({
          api: dir({
            Dockerfile: file(DOCKERFILE),
            ".github": dir({
              workflows: dir({ "ci.yml": file(CI_WORKFLOW) }),
            }),
          }),
        }),
        infra: dir({
          "main.tf": file(MAIN_TF),
          "variables.tf": file('variable "ami_id" {\n  type    = string\n  default = "ami-0c55b159cbfafe1f0"\n}\n'),
        }),
        ansible: dir({
          "playbook.yml": file(PLAYBOOK),
          "inventory.ini": file("[web]\nweb-1 ansible_host=10.0.0.21\n"),
        }),
        ".ssh": dir({
          "known_hosts": file("github.com ssh-ed25519 AAAAC3NzaC1lZDI1NTE5AAAAIOMqqnkVzrm0SdG6UOoqKLsabgH5C9okWi0dh2l9GKJl\n"),
          id_ed25519: file("-----BEGIN OPENSSH PRIVATE KEY-----\n(учебный приватный ключ)\n-----END OPENSSH PRIVATE KEY-----\n", "-rw-------"),
          "id_ed25519.pub": file("ssh-ed25519 AAAAC3NzaC1lZDI1NTE5AAAAIStudentKey student@devops\n"),
          config: file(SSH_CONFIG),
        }),
      }),
    }),
    etc: dir({
      hosts: file("127.0.0.1 localhost\n10.0.0.5 db\n"),
      passwd: file(PASSWD),
      "resolv.conf": file("nameserver 10.0.0.2\nsearch cluster.local\n"),
      nginx: dir({
        "nginx.conf": file(NGINX_CONF),
        "sites-enabled": dir({
          "app.conf": file(NGINX_CONF),
        }),
      }),
      prometheus: dir({
        "prometheus.yml": file(PROMETHEUS_YML),
        "alerts.yml": file(ALERTS_YML),
      }),
    }),
    var: dir({
      log: dir({
        "app.log": file(APP_LOG),
        "syslog": file("Aug 27 09:00:01 devops systemd[1]: Started nginx.\n"),
        nginx: dir({
          "access.log": file(NGINX_ACCESS),
          "error.log": file(
            '2026/08/27 09:15:02 [error] 712#712: *184 connect() failed (111: Connection refused) while connecting to upstream, client: 10.0.0.9, server: app.example.com, upstream: "http://10.0.0.4:8080/api/orders"\n',
          ),
        }),
      }),
    }),
    proc: dir({
      cpuinfo: file(CPUINFO),
      meminfo: file(MEMINFO),
    }),
    sys: dir({
      class: dir({
        net: dir({
          eth0: dir({ address: file("02:42:ac:11:00:02\n") }),
          lo: dir({ address: file("00:00:00:00:00:00\n") }),
        }),
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
  if (path === "~" || path === "$HOME") path = `/home/${state.user}`;
  else if (path.startsWith("~/")) path = `/home/${state.user}/${path.slice(2)}`;
  else if (path.startsWith("$HOME/")) path = `/home/${state.user}/${path.slice(6)}`;

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
    case "inspect":
      return [
        "[",
        "    {",
        '        "Id": "3f1c9a4b7e21c8f0a6d9e4c7b3128f5a",',
        '        "State": { "Status": "running", "Running": true, "ExitCode": 0, "OOMKilled": false },',
        '        "Config": { "Image": "nginx:1.27", "Env": ["PATH=/usr/bin"] },',
        '        "NetworkSettings": { "IPAddress": "172.17.0.2", "Ports": { "80/tcp": "8080" } }',
        "    }",
        "]",
      ].join("\n");
    case "stats":
      return [
        "CONTAINER   CPU %   MEM USAGE / LIMIT   MEM %   NET I/O",
        "web         0.03%   12.4MiB / 512MiB    2.42%   1.2kB / 0B",
        "db          0.15%   48.9MiB / 512MiB    9.55%   3.4kB / 0B",
      ].join("\n");
    case "top":
      return "UID    PID    PPID   C   CMD\nroot   1043   1020   0   nginx: master process";
    default:
      return `docker: '${sub ?? ""}' is not a docker command.\nSee 'docker --help'`;
  }
}

// Имитация вывода kubectl.
function runKubectl(args: string[]): string {
  const [sub, ...rest] = args;
  switch (sub) {
    case "get":
      if (rest[0]?.startsWith("ev")) {
        return [
          "LAST SEEN   TYPE      REASON      OBJECT        MESSAGE",
          "2m          Normal    Scheduled   pod/api-7d9f  Successfully assigned api-7d9f to node-1",
          "2m          Normal    Pulled      pod/api-7d9f  Container image \"api:1.5.0\" already present",
          "90s         Warning   BackOff     pod/api-7d9f  Back-off restarting failed container",
          "80s         Warning   Unhealthy   pod/api-7d9f  Readiness probe failed: connection refused",
        ].join("\n");
      }
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
    case "exec":
      return "/ # (вы внутри пода, для выхода наберите exit)";
    case "top":
      return [
        "NAME          CPU(cores)   MEMORY(bytes)",
        "api-7d9f      42m          128Mi",
        "api-3k2p      38m          121Mi",
      ].join("\n");
    case "port-forward":
      return "Forwarding from 127.0.0.1:8080 -> 8080\nForwarding from [::1]:8080 -> 8080";
    default:
      return `error: unknown command "${sub ?? ""}" for "kubectl"`;
  }
}

const HELP = [
  "Доступные команды учебного терминала:",
  "  файлы:      pwd, ls, cd, cat, head, tail, grep, find, echo, sort, uniq, wc",
  "              mkdir, touch, rm, cp, mv, chmod, chown, du",
  "  система:    ps, whoami, uname, df, free, uptime, date, systemctl, journalctl",
  "  сеть:       curl, wget, dig, ss, ping, nginx",
  "  git:        status, log, branch, checkout, add, commit, push, pull, merge, diff, revert",
  "  контейнеры: docker …, kubectl …, helm …",
  "  инфра:      terraform …, ansible, ansible-playbook",
  "  секреты:    openssl, ssh-keygen",
  "  прочее:     history, clear, help",
  "",
  "Работают конвейеры: например `cat app.log | grep ERROR | wc -l`.",
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
    const cmd = tokens[0];
    const flags = tokens.filter((t) => t.startsWith("-"));
    const rows = result.output ? result.output.split("\n") : [];

    if (cmd === "grep") {
      const pattern = tokens.find((t, i) => i > 0 && !t.startsWith("-")) ?? "";
      const insensitive = flags.some((f) => f.includes("i"));
      const invert = flags.some((f) => f.includes("v"));
      const needle = insensitive ? pattern.toLowerCase() : pattern;
      result = {
        ...result,
        output: rows
          .filter((row) => {
            const has = (insensitive ? row.toLowerCase() : row).includes(needle);
            return invert ? !has : has;
          })
          .join("\n"),
      };
    } else if (cmd === "wc") {
      result = { ...result, output: String(result.output ? rows.length : 0) };
    } else if (cmd === "head" || cmd === "tail") {
      const n = Number(tokens[tokens.indexOf("-n") + 1]) || 10;
      result = { ...result, output: (cmd === "head" ? rows.slice(0, n) : rows.slice(-n)).join("\n") };
    } else if (cmd === "sort") {
      const numeric = flags.some((f) => f.includes("n"));
      const reverse = flags.some((f) => f.includes("r"));
      const sorted = [...rows].sort((a, b) =>
        numeric ? (parseFloat(a) || 0) - (parseFloat(b) || 0) : a.localeCompare(b),
      );
      if (reverse) sorted.reverse();
      result = { ...result, output: sorted.join("\n") };
    } else if (cmd === "uniq") {
      const withCount = flags.some((f) => f.includes("c"));
      const out: string[] = [];
      let prev: string | null = null;
      let count = 0;
      for (const row of rows) {
        if (row === prev) count += 1;
        else {
          if (prev !== null) out.push(withCount ? `${String(count).padStart(7)} ${prev}` : prev);
          prev = row;
          count = 1;
        }
      }
      if (prev !== null) out.push(withCount ? `${String(count).padStart(7)} ${prev}` : prev);
      result = { ...result, output: out.join("\n") };
    } else if (cmd === "cut") {
      const delim = tokens[tokens.indexOf("-d") + 1] ?? "\t";
      const fieldSpec = tokens[tokens.indexOf("-f") + 1] ?? "1";
      const fields = fieldSpec.split(",").map((x) => Number(x) - 1);
      result = {
        ...result,
        output: rows
          .map((row) => {
            const cols = row.split(delim);
            return fields.map((i) => cols[i] ?? "").join(delim);
          })
          .join("\n"),
      };
    } else if (cmd === "tr") {
      const a = tokens[1] ?? "";
      const b = tokens[2] ?? "";
      result = {
        ...result,
        output:
          a === "a-z" && b === "A-Z"
            ? result.output.toUpperCase()
            : a === "A-Z" && b === "a-z"
              ? result.output.toLowerCase()
              : result.output,
      };
    } else if (cmd === "xargs") {
      // Поддержка: xargs wc -l — посчитать строки каждого файла из потока.
      if (tokens[1] === "wc") {
        const out = rows
          .filter(Boolean)
          .map((name) => {
            const segments = resolve(state, name);
            const node = segments && nodeAt(state, segments);
            const n = node && node.type === "file" ? node.content.trimEnd().split("\n").length : 0;
            return `${String(n).padStart(7)} ${name}`;
          });
        result = { ...result, output: out.join("\n") };
      }
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
      const nIndex = args.indexOf("-n");
      const count = nIndex >= 0 ? Number(args[nIndex + 1]) || 10 : 10;
      // Имя файла — операнд, который не является значением флага -n.
      const nValue = nIndex >= 0 ? args[nIndex + 1] : null;
      const target = operands.find((o) => o !== nValue) ?? operands[operands.length - 1];
      if (!target) return { output: `${command}: не указан файл`, state };
      const segments = resolve(state, target);
      const node = segments && nodeAt(state, segments);
      if (!node || node.type !== "file") {
        return { output: `${command}: ${target}: нет такого файла`, state };
      }
      const rows = node.content.trimEnd().split("\n");
      const slice = command === "head" ? rows.slice(0, count) : rows.slice(-count);
      const follow = hasFlag("f") ? "\n(следим за файлом… Ctrl+C для выхода)" : "";
      return { output: slice.join("\n") + follow, state };
    }

    case "grep": {
      const pattern = operands[0] ?? "";
      const target = operands[1];
      const insensitive = hasFlag("i");
      const countOnly = hasFlag("c");
      const numbered = hasFlag("n");
      if (!target) return { output: "grep: не указан файл", state };

      const segments = resolve(state, target);
      const node = segments && nodeAt(state, segments);
      if (!node || node.type !== "file") {
        return { output: `grep: ${target}: нет такого файла или каталога`, state };
      }
      const needle = insensitive ? pattern.toLowerCase() : pattern;
      const matches: string[] = [];
      node.content.split("\n").forEach((row, i) => {
        if ((insensitive ? row.toLowerCase() : row).includes(needle)) {
          matches.push(numbered ? `${i + 1}:${row}` : row);
        }
      });
      if (countOnly) return { output: String(matches.length), state };
      return { output: matches.join("\n"), state };
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

    case "git":
      return { output: runGit(args), state };

    // --- Сеть ---
    case "curl":
      return { output: runCurl(args), state };

    case "wget":
      return {
        output: "--2026-08-27 09:20:14--  " + (operands[0] ?? "") + "\nHTTP request sent, awaiting response... 200 OK\nSaved",
        state,
      };

    case "ip":
    case "ifconfig": {
      const eth0 = [
        "eth0: flags=4163<UP,BROADCAST,RUNNING,MULTICAST>  mtu 1500",
        "        inet 10.0.0.21  netmask 255.255.255.0  broadcast 10.0.0.255",
        "        ether 02:42:ac:11:00:02  txqueuelen 1000  (Ethernet)",
      ].join("\n");
      const lo = [
        "lo: flags=73<UP,LOOPBACK,RUNNING>  mtu 65536",
        "        inet 127.0.0.1  netmask 255.0.0.0",
        "        loop  txqueuelen 1000  (Local Loopback)",
      ].join("\n");
      if (operands.includes("lo")) return { output: lo, state };
      if (operands.includes("eth0")) return { output: eth0, state };
      return { output: eth0 + "\n" + lo, state };
    }

    case "ssh":
      return {
        output: `Welcome to Ubuntu 24.04 LTS (GNU/Linux)\ndeploy@${operands[0] ?? "server"}:~$ (учебная сессия, наберите exit для выхода)`,
        state,
      };

    case "scp": {
      const src = operands[0] ?? "file";
      return { output: `${src}                              100% 1234     1.2KB/s   00:00`, state };
    }

    case "dig":
    case "nslookup":
      return { output: runDig(args), state };

    case "ss":
    case "netstat":
      return { output: runSs(args), state };

    case "ping":
      return { output: runPing(args), state };

    case "nginx":
      return { output: runNginx(args), state };

    // --- Инфраструктура как код ---
    case "terraform":
      return { output: runTerraform(args), state };

    case "ansible":
    case "ansible-playbook":
      return { output: runAnsible(command, args), state };

    case "helm":
      return { output: runHelm(args), state };

    // --- Безопасность ---
    case "openssl":
      return { output: runOpenssl(args), state };

    case "ssh-keygen":
      return { output: runSshKeygen(args), state };

    // --- Ресурсы системы ---
    case "free":
      return { output: runFree(args), state };

    case "du":
      return { output: runDu(args), state };

    case "uptime":
      return { output: runUptime(), state };

    case "date":
      return { output: runDate(), state };

    case "cut": {
      const target = operands[operands.length - 1];
      const segments = target ? resolve(state, target) : null;
      const node = segments && nodeAt(state, segments);
      if (!node || node.type !== "file") {
        return { output: `cut: ${target ?? ""}: нет такого файла`, state };
      }
      const delim = args[args.indexOf("-d") + 1] ?? "\t";
      const fieldSpec = args[args.indexOf("-f") + 1] ?? "1";
      const fields = fieldSpec.split(",").map((x) => Number(x) - 1);
      const out = node.content
        .trimEnd()
        .split("\n")
        .map((row) => {
          const cols = row.split(delim);
          return fields.map((i) => cols[i] ?? "").join(delim);
        });
      return { output: out.join("\n"), state };
    }

    case "sort": {
      const target = operands[0];
      const segments = target ? resolve(state, target) : null;
      const node = segments && nodeAt(state, segments);
      if (!node || node.type !== "file") {
        return { output: `sort: ${target ?? ""}: нет такого файла`, state };
      }
      const rows = node.content.trimEnd().split("\n").sort();
      if (hasFlag("r")) rows.reverse();
      return { output: rows.join("\n"), state };
    }

    case "uniq": {
      const target = operands[0];
      const segments = target ? resolve(state, target) : null;
      const node = segments && nodeAt(state, segments);
      if (!node || node.type !== "file") {
        return { output: `uniq: ${target ?? ""}: нет такого файла`, state };
      }
      const rows = node.content.trimEnd().split("\n");
      const unique = rows.filter((row, index) => index === 0 || row !== rows[index - 1]);
      return { output: unique.join("\n"), state };
    }

    case "wc": {
      const target = operands[0];
      const segments = target ? resolve(state, target) : null;
      const node = segments && nodeAt(state, segments);
      if (!node || node.type !== "file") {
        return { output: `wc: ${target ?? ""}: нет такого файла`, state };
      }
      const content = node.content.trimEnd();
      const lines = content ? content.split("\n").length : 0;
      const words = content ? content.split(/\s+/).filter(Boolean).length : 0;
      if (hasFlag("l")) return { output: `${lines} ${target}`, state };
      if (hasFlag("w")) return { output: `${words} ${target}`, state };
      return { output: `${lines} ${words} ${content.length} ${target}`, state };
    }

    case "exit":
      return { output: "выход из контейнера", state };

    default:
      return { output: `${command}: команда не найдена. Наберите help для списка команд.`, state };
  }
}
