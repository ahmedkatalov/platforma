// Лёгкая подсветка синтаксиса на Prism (без тяжёлого Monaco). Один модуль на
// оба места: примеры кода в уроках и редактор, где студент пишет код.
import Prism from "prismjs";

// Грамматики нужных языков (js/ts/markup/css уже в ядре Prism).
import "prismjs/components/prism-typescript";
import "prismjs/components/prism-jsx";
import "prismjs/components/prism-tsx";
import "prismjs/components/prism-go";
import "prismjs/components/prism-bash";
import "prismjs/components/prism-yaml";
import "prismjs/components/prism-json";
import "prismjs/components/prism-sql";
import "prismjs/components/prism-docker";
import "prismjs/components/prism-hcl";
import "prismjs/components/prism-nginx";
import "prismjs/components/prism-python";
import "prismjs/components/prism-protobuf";
import "prismjs/components/prism-markdown";
import "prismjs/components/prism-toml";
import "prismjs/components/prism-diff";

// Псевдонимы языков → грамматика Prism.
const ALIAS: Record<string, string> = {
  sh: "bash",
  shell: "bash",
  zsh: "bash",
  console: "bash",
  dockerfile: "docker",
  proto: "protobuf",
  terraform: "hcl",
  tf: "hcl",
  yml: "yaml",
  golang: "go",
  js: "javascript",
  ts: "typescript",
  md: "markdown",
  "": "text",
  text: "text",
  txt: "text",
  plaintext: "text",
};

export function resolveLang(lang?: string): string {
  const l = (lang ?? "").toLowerCase().replace(/^language-/, "").trim();
  return ALIAS[l] ?? l;
}

function escapeHtml(s: string): string {
  return s.replace(/&/g, "&amp;").replace(/</g, "&lt;").replace(/>/g, "&gt;");
}

// Возвращает HTML с токенами Prism. Если грамматики нет — просто экранированный
// текст (подсветки не будет, но код читается и ничего не ломается).
export function highlightToHtml(code: string, lang?: string): string {
  const resolved = resolveLang(lang);
  const grammar = Prism.languages[resolved];
  if (!grammar) return escapeHtml(code);
  try {
    return Prism.highlight(code, grammar, resolved);
  } catch {
    return escapeHtml(code);
  }
}
