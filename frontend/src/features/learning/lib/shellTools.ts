// Имитация инструментов DevOps для учебного терминала: вывод статичный,
// но узнаваемый — студент учится читать реальный формат ответа.

const GIT_LOG = [
  "c3f9a21 (HEAD -> feature/healthcheck) feat: add healthcheck endpoint",
  "9d1b70e (origin/main, main) chore: bump dependencies",
  "4a2c8e1 fix: retry payment requests on timeout",
  "7b51d3f docs: describe deploy process",
].join("\n");

const GIT_STATUS_DIRTY = [
  "On branch feature/healthcheck",
  "Your branch is ahead of 'origin/main' by 1 commit.",
  "",
  "Changes not staged for commit:",
  '  (use "git add <file>..." to update what will be committed)',
  "        modified:   internal/handler/health.go",
  "",
  "Untracked files:",
  '  (use "git add <file>..." to include in what will be committed)',
  "        deploy/healthcheck.yml",
  "",
  'no changes added to commit (use "git add" and/or "git commit -a")',
].join("\n");

export function runGit(args: string[]): string {
  const [sub, ...rest] = args;

  switch (sub) {
    case "status":
      return GIT_STATUS_DIRTY;

    case "log":
      if (rest.includes("--oneline")) return GIT_LOG;
      return [
        "commit c3f9a21b8e4d1a77f0c9e2d3b4a5f6071829304a (HEAD -> feature/healthcheck)",
        "Author: Student <student@devops.local>",
        "Date:   Thu Aug 27 09:12:04 2026 +0300",
        "",
        "    feat: add healthcheck endpoint",
      ].join("\n");

    case "branch":
      if (rest[0] && !rest[0].startsWith("-")) return `Switched to a new branch '${rest[0]}'`;
      return "* feature/healthcheck\n  main";

    case "checkout":
      if (rest[0] === "-b") return `Switched to a new branch '${rest[1] ?? "new-branch"}'`;
      return `Switched to branch '${rest[0] ?? "main"}'`;

    case "switch":
      if (rest[0] === "-c") return `Switched to a new branch '${rest[1] ?? "new-branch"}'`;
      return `Switched to branch '${rest[0] ?? "main"}'`;

    case "add":
      return "";

    case "commit":
      return "[feature/healthcheck 1f4d9c2] " +
        (extractMessage(rest) || "commit message") +
        "\n 2 files changed, 34 insertions(+), 3 deletions(-)";

    case "push":
      return [
        "Enumerating objects: 9, done.",
        "Writing objects: 100% (5/5), 612 bytes | 612.00 KiB/s, done.",
        "To github.com:team/api.git",
        "   9d1b70e..1f4d9c2  feature/healthcheck -> feature/healthcheck",
      ].join("\n");

    case "pull":
      return "Already up to date.";

    case "fetch":
      return "From github.com:team/api\n   9d1b70e..e5a7c31  main       -> origin/main";

    case "merge":
      return `Updating 9d1b70e..1f4d9c2\nFast-forward\n 2 files changed, 34 insertions(+)`;

    case "diff":
      return [
        "diff --git a/internal/handler/health.go b/internal/handler/health.go",
        "@@ -12,6 +12,10 @@ func Health(w http.ResponseWriter, r *http.Request) {",
        "+\tif err := db.Ping(r.Context()); err != nil {",
        '+\t\thttp.Error(w, "db unavailable", http.StatusServiceUnavailable)',
        "+\t\treturn",
        "+\t}",
      ].join("\n");

    case "reset":
      return rest.includes("--hard")
        ? "HEAD is now at 9d1b70e chore: bump dependencies"
        : "Unstaged changes after reset:\nM\tinternal/handler/health.go";

    case "revert":
      return '[main 8c2e5f9] Revert "feat: add healthcheck endpoint"\n 1 file changed, 4 deletions(-)';

    case "stash":
      if (rest[0] === "pop") return "On branch feature/healthcheck\nDropped refs/stash@{0}";
      return "Saved working directory and index state WIP on feature/healthcheck";

    case "remote":
      return rest.includes("-v")
        ? "origin\tgit@github.com:team/api.git (fetch)\norigin\tgit@github.com:team/api.git (push)"
        : "origin";

    case "clone":
      return `Cloning into '${(rest[0] ?? "repo").split("/").pop()?.replace(/\.git$/, "")}'...\nremote: Enumerating objects: 142, done.`;

    case "tag":
      if (rest[0]) return "";
      return "v1.3.0\nv1.4.0\nv1.4.2";

    case "init":
      return "Initialized empty Git repository in /home/student/projects/.git/";

    default:
      return `git: '${sub ?? ""}' is not a git command. See 'git --help'.`;
  }
}

function extractMessage(args: string[]): string {
  const index = args.findIndex((arg) => arg === "-m" || arg === "--message");
  return index >= 0 ? (args[index + 1] ?? "") : "";
}

// --- Сеть ---

export function runCurl(args: string[]): string {
  const url = args.find((arg) => !arg.startsWith("-")) ?? "";
  const headOnly = args.includes("-I") || args.includes("--head");
  const verbose = args.includes("-v");

  if (url.includes("/metrics")) {
    return [
      "# HELP http_requests_total Total number of HTTP requests",
      "# TYPE http_requests_total counter",
      'http_requests_total{code="200",path="/api/orders"} 18422',
      'http_requests_total{code="500",path="/api/orders"} 37',
      "# HELP http_request_duration_seconds Request duration",
      "# TYPE http_request_duration_seconds histogram",
      'http_request_duration_seconds_bucket{le="0.5"} 17980',
      'http_request_duration_seconds_bucket{le="1"} 18310',
    ].join("\n");
  }

  if (url.includes("/health")) {
    return headOnly
      ? "HTTP/1.1 200 OK\nContent-Type: application/json\nContent-Length: 15"
      : '{"status":"ok"}';
  }

  const head = [
    "HTTP/1.1 200 OK",
    "Server: nginx/1.27.1",
    "Content-Type: text/html; charset=utf-8",
    "Content-Length: 1256",
    "Connection: keep-alive",
  ].join("\n");

  if (headOnly) return head;
  if (verbose) return `*   Trying 10.0.0.5:80...\n* Connected to ${url || "app"} port 80\n${head}`;
  return "<!doctype html>\n<html><head><title>App</title></head><body>OK</body></html>";
}

export function runDig(args: string[]): string {
  const name = args.find((arg) => !arg.startsWith("+") && !arg.startsWith("-")) ?? "example.com";
  if (args.includes("+short")) return "93.184.216.34";

  return [
    "; <<>> DiG 9.18.24 <<>> " + name,
    ";; ->>HEADER<<- opcode: QUERY, status: NOERROR, id: 41230",
    "",
    ";; QUESTION SECTION:",
    `;${name}.\t\t\tIN\tA`,
    "",
    ";; ANSWER SECTION:",
    `${name}.\t\t300\tIN\tA\t93.184.216.34`,
    "",
    ";; Query time: 12 msec",
  ].join("\n");
}

export function runSs(args: string[]): string {
  const listening = args.some((arg) => arg.includes("l"));
  const rows = [
    "State    Recv-Q   Send-Q   Local Address:Port    Peer Address:Port  Process",
    "LISTEN   0        4096           0.0.0.0:80           0.0.0.0:*      users:((\"nginx\",pid=712))",
    "LISTEN   0        4096           0.0.0.0:8080         0.0.0.0:*      users:((\"app\",pid=1043))",
    "LISTEN   0        244          127.0.0.1:5432         0.0.0.0:*      users:((\"postgres\",pid=880))",
  ];
  if (!listening) {
    rows.push("ESTAB    0        0           10.0.0.4:8080       10.0.0.9:51422  users:((\"app\",pid=1043))");
  }
  return rows.join("\n");
}

export function runPing(args: string[]): string {
  const host = args.find((arg) => !arg.startsWith("-")) ?? "db";
  return [
    `PING ${host} (10.0.0.5) 56(84) bytes of data.`,
    `64 bytes from ${host} (10.0.0.5): icmp_seq=1 ttl=64 time=0.412 ms`,
    `64 bytes from ${host} (10.0.0.5): icmp_seq=2 ttl=64 time=0.377 ms`,
    "",
    `--- ${host} ping statistics ---`,
    "2 packets transmitted, 2 received, 0% packet loss, time 1002ms",
  ].join("\n");
}

export function runNginx(args: string[]): string {
  if (args.includes("-t")) {
    return [
      "nginx: the configuration file /etc/nginx/nginx.conf syntax is ok",
      "nginx: configuration file /etc/nginx/nginx.conf test is successful",
    ].join("\n");
  }
  if (args.includes("-s")) return "";
  if (args.includes("-v") || args.includes("-V")) return "nginx version: nginx/1.27.1";
  return "nginx: invalid option — используйте -t для проверки конфигурации";
}

// --- Инфраструктура как код ---

export function runTerraform(args: string[]): string {
  const [sub] = args;

  switch (sub) {
    case "init":
      return [
        "Initializing the backend...",
        "Initializing provider plugins...",
        '- Finding hashicorp/aws versions matching "~> 5.0"...',
        "- Installing hashicorp/aws v5.62.0...",
        "",
        "Terraform has been successfully initialized!",
      ].join("\n");

    case "fmt":
      return "main.tf";

    case "validate":
      return "Success! The configuration is valid.";

    case "plan":
      return [
        "Terraform will perform the following actions:",
        "",
        "  # aws_instance.app will be created",
        '  + resource "aws_instance" "app" {',
        '      + ami           = "ami-0c55b159cbfafe1f0"',
        '      + instance_type = "t3.micro"',
        "      + id            = (known after apply)",
        "    }",
        "",
        "Plan: 1 to add, 0 to change, 0 to destroy.",
      ].join("\n");

    case "apply":
      return [
        "aws_instance.app: Creating...",
        "aws_instance.app: Creation complete after 21s [id=i-0a1b2c3d4e5f6a7b8]",
        "",
        "Apply complete! Resources: 1 added, 0 changed, 0 destroyed.",
      ].join("\n");

    case "destroy":
      return "Destroy complete! Resources: 1 destroyed.";

    case "state":
      return "aws_instance.app";

    case "output":
      return 'app_ip = "10.0.0.14"';

    default:
      return "Usage: terraform [init|validate|plan|apply|destroy|fmt|output]";
  }
}

export function runAnsible(command: string, args: string[]): string {
  if (command === "ansible-playbook") {
    const playbook = args.find((arg) => arg.endsWith(".yml")) ?? "playbook.yml";
    if (args.includes("--check")) {
      return [
        `PLAY [Configure web servers] ***`,
        "",
        "TASK [Install nginx] ***",
        "changed: [web-1]",
        "",
        "PLAY RECAP ***",
        "web-1   : ok=2    changed=1    unreachable=0    failed=0    skipped=0",
        "",
        "(режим проверки: изменения не применялись)",
      ].join("\n");
    }
    return [
      `PLAY [Configure web servers] *** (${playbook})`,
      "",
      "TASK [Gathering Facts] ***",
      "ok: [web-1]",
      "",
      "TASK [Install nginx] ***",
      "changed: [web-1]",
      "",
      "TASK [Start nginx service] ***",
      "ok: [web-1]",
      "",
      "PLAY RECAP ***",
      "web-1   : ok=3    changed=1    unreachable=0    failed=0    skipped=0",
    ].join("\n");
  }

  // ansible <hosts> -m <module>
  if (args.includes("-m")) {
    const host = args[0] ?? "all";
    return `${host} | SUCCESS => {\n    "changed": false,\n    "ping": "pong"\n}`;
  }
  return "Usage: ansible <hosts> -m <module> | ansible-playbook playbook.yml";
}

export function runHelm(args: string[]): string {
  const [sub] = args;
  switch (sub) {
    case "list":
      return [
        "NAME  NAMESPACE  REVISION  STATUS    CHART        APP VERSION",
        "api   default    3         deployed  api-0.4.1    1.4.2",
      ].join("\n");
    case "install":
    case "upgrade":
      return "Release \"api\" has been upgraded. Happy Helming!\nREVISION: 4\nSTATUS: deployed";
    case "rollback":
      return "Rollback was a success! Happy Helming!";
    case "uninstall":
      return 'release "api" uninstalled';
    default:
      return "Usage: helm [list|install|upgrade|rollback|uninstall]";
  }
}

// --- Безопасность ---

export function runOpenssl(args: string[]): string {
  const [sub, ...rest] = args;

  if (sub === "rand") {
    return rest.includes("-hex")
      ? "9f2c41b8a7e35d06c1948fbb27ad5e30f6c8a91d4b7e2035c8fa16d940b3e7c2"
      : "случайные байты (используйте -hex или -base64)";
  }
  if (sub === "x509") {
    return [
      "subject=CN = app.example.com",
      "issuer=C = US, O = Let's Encrypt, CN = R11",
      "notBefore=Jul 15 08:00:00 2026 GMT",
      "notAfter=Oct 13 08:00:00 2026 GMT",
    ].join("\n");
  }
  if (sub === "s_client") {
    return [
      "CONNECTED(00000003)",
      "depth=2 C = US, O = Internet Security Research Group, CN = ISRG Root X1",
      "SSL handshake has read 4821 bytes and written 401 bytes",
      "Verify return code: 0 (ok)",
    ].join("\n");
  }
  return "Usage: openssl [rand|x509|s_client] …";
}

export function runSshKeygen(args: string[]): string {
  const type = args[args.indexOf("-t") + 1] ?? "ed25519";
  return [
    `Generating public/private ${type} key pair.`,
    "Your identification has been saved in /home/student/.ssh/id_" + type,
    "Your public key has been saved in /home/student/.ssh/id_" + type + ".pub",
    "The key fingerprint is:",
    "SHA256:9pQ1c4bZm0YfR7kTn2xJ5dLwV8sHe3AoU6yPq1BvCgE student@devops",
  ].join("\n");
}

// --- Мониторинг и система ---

export function runFree(args: string[]): string {
  const human = args.includes("-h");
  return human
    ? [
        "               total        used        free      shared  buff/cache   available",
        "Mem:           3.8Gi       1.6Gi       0.4Gi        18Mi       1.8Gi       1.9Gi",
        "Swap:          1.0Gi          0B       1.0Gi",
      ].join("\n")
    : [
        "               total        used        free      shared  buff/cache   available",
        "Mem:         3999232     1682944      419840       18432     1896448     1994752",
        "Swap:        1048576           0     1048576",
      ].join("\n");
}

export function runDu(args: string[]): string {
  const human = args.includes("-h");
  return human
    ? ["1.2M\t./projects", "620K\t./app.log", "4.0K\t./notes.txt", "1.9M\t."].join("\n")
    : ["1204\t./projects", "620\t./app.log", "4\t./notes.txt", "1908\t."].join("\n");
}

export function runDate(): string {
  return "Thu Aug 27 09:20:14 MSK 2026";
}

export function runUptime(): string {
  return " 09:20:14 up 14 days,  3:07,  1 user,  load average: 0.42, 0.55, 0.61";
}
