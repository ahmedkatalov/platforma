#!/usr/bin/env python3
"""check_course.py — проверки курса «DevOps с нуля до практики» (раздел 6 промпта).

Использование:
  python3 check_course.py devops-engineer.course.json        # экспорт platforma-course v1
  python3 check_course.py backend/internal/seed/content/     # папка seed (001.json … 018.json)
  добавьте --net, чтобы проверять HTTP-статусы ссылок (нужна сеть).
"""
import json, re, sys, os, glob, uuid
from collections import Counter, defaultdict

NET = '--net' in sys.argv
args = [a for a in sys.argv[1:] if not a.startswith('--')]
if not args:
    print(__doc__); sys.exit(2)
path = args[0]

# ---------- загрузка: экспорт или seed ----------
def load(path):
    if os.path.isdir(path):
        mods = []
        for f in sorted(glob.glob(os.path.join(path, '*.json'))):
            d = json.load(open(f, encoding='utf-8'))
            mods.append({'id': None, 'title': d.get('title'), 'summary': d.get('summary'), 'lessons': d['lessons'], '_file': os.path.basename(f)})
        return {'format': 'seed', 'course': {}, 'modules': mods}, 'seed'
    d = json.load(open(path, encoding='utf-8'))
    return d, 'export'

data, mode = load(path)
mods = data['modules']
problems = defaultdict(list)   # проверка -> список нарушений
def bad(check, msg): problems[check].append(msg)
def ref(mi, li): return f"M{mi+1}.L{li+1}"

# ---------- 1. схема ----------
TOP = {'format', 'version', 'course', 'modules'}
COURSE = {'id', 'slug', 'title', 'subtitle', 'description', 'coverUrl', 'level', 'tags'}
MODULE = {'id', 'title', 'summary', 'lessons'}
LESSON = {'id', 'title', 'kind', 'summary', 'content', 'durationMin'}
TEXT = {'body', 'resources'}
QUIZ = {'intro', 'passScore', 'questions', 'resources', 'shuffle', 'timeLimitSec'}
Q_SINGLE = {'id', 'text', 'options', 'explanation', 'multiple', 'review', 'hint', 'type'}
Q_MATCH = {'id', 'text', 'type', 'pairs', 'explanation', 'review', 'hint'}
Q_BLANK = {'id', 'hint', 'text', 'type', 'accept', 'explanation', 'review'}
Q_ORDER = {'id', 'text', 'type', 'items', 'explanation', 'review', 'hint'}
TERM = {'debug', 'intro', 'shell', 'tasks', 'challenge', 'resources'}
TASK = {'id', 'hint', 'hints', 'prompt', 'predict', 'success', 'expected', 'pattern'}
CODE = {'hint', 'task', 'checks', 'starter', 'language', 'solution', 'resources'}
CHECK = {'type', 'value', 'message'}
RES = {'url', 'note', 'title'}
KINDS = {'text', 'quiz', 'terminal', 'code'}
LANGS = {'yaml', 'bash', 'hcl', 'dockerfile', 'nginx'}
CTYPES = {'contains', 'regex', 'notContains'}

def unknown(obj, allowed, where):
    extra = set(obj.keys()) - allowed
    if extra: bad('1.схема', f"{where}: неизвестные ключи {sorted(extra)}")

if mode == 'export':
    unknown(data, TOP, 'верхний уровень')
    if data.get('format') != 'platforma-course' or data.get('version') != 1: bad('1.схема', 'format/version не platforma-course v1')
    unknown(data.get('course', {}), COURSE, 'course')

# ---------- 2. идентификаторы ----------
ids = Counter()
def is_uuid(s):
    try: uuid.UUID(str(s)); return True
    except Exception: return False
if mode == 'export':
    cid = data['course'].get('id'); ids[cid] += 1
    if not is_uuid(cid): bad('2.id', f"course.id не UUID: {cid}")

def ok_check(text, c):
    t, v = c.get('type'), c.get('value', '')
    if t == 'contains': return v in text
    if t == 'notContains': return v not in text
    if t == 'regex':
        try: return re.search(v, text) is not None
        except re.error: return False
    return False

OLD = ['node:20-slim', 'node:20-alpine', 'node:20', 'node-version: "20"', 'prom/prometheus:v2.53.0', 'grafana/grafana:11.1.0',
       'nginx:1.27', 'nginx:1.28', 'postgres:16', 'golang:1.22', 'alpine:3.20', 'python:3.12-slim', 'python:3.12', 'busybox:1.36',
       'redis:7', 'Terraform v1.7.0', 'required_version = ">= 1.6"', 'version = "~> 5.0"', 'Client Version: v1.30.0', 'ubuntu-22.04',
       'actions/checkout@v4', 'actions/setup-node@v4', 'docker-compose ', 'rm /var/log/app.log']
HEADS = ['Зачем это нужно', 'Проверьте себя', 'Частые ошибки новичка', 'Запомнить']

dur_by_mod = []; kinds = Counter(); resources = []
for mi, m in enumerate(mods):
    if mode == 'export':
        unknown(m, MODULE, f"module {mi+1}"); ids[m.get('id')] += 1
        if not is_uuid(m.get('id')): bad('2.id', f"module {mi+1} id не UUID")
    dsum = 0
    for li, l in enumerate(m['lessons']):
        r = ref(mi, li); k = l.get('kind'); c = l.get('content', {}) or {}
        kinds[k] += 1; dsum += int(l.get('durationMin') or 0)
        if mode == 'export':
            unknown(l, LESSON, r); ids[l.get('id')] += 1
            if not is_uuid(l.get('id')): bad('2.id', f"{r} id не UUID")
        if k not in KINDS: bad('2.id', f"{r} kind={k}")
        for res in (c.get('resources') or []):
            unknown(res, RES, f"{r} resource"); resources.append((r, res.get('title', ''), res.get('url', '')))
        blob = json.dumps(l, ensure_ascii=False)
        for o in OLD:
            if o in blob: bad('7.устаревшее', f"{r}: «{o}»")
        if re.search(r'kubectl logs app(?![-/\w])', blob): bad('7.устаревшее', f"{r}: «kubectl logs app» как отдельная команда")

        if k == 'text':
            unknown(c, TEXT, r); body = c.get('body', '') or ''
            for h in HEADS:
                if h not in body: bad('6.текст', f"{r}: нет заголовка «{h}»")
            if not (5500 <= len(body) <= 10000): bad('6.текст', f"{r}: длина body {len(body)} (норма 5500–10000)")
        elif k == 'quiz':
            unknown(c, QUIZ, r)
            if c.get('passScore') != 70: bad('5.квиз', f"{r}: passScore={c.get('passScore')}")
            for q in c.get('questions', []):
                qt = q.get('type', 'single'); qid = q.get('id')
                if qt == 'match':
                    unknown(q, Q_MATCH, f"{r} {qid}")
                    if len(q.get('pairs', [])) < 3: bad('5.квиз', f"{r} {qid}: match < 3 пар")
                elif qt == 'blank':
                    unknown(q, Q_BLANK, f"{r} {qid}")
                    if not q.get('accept'): bad('5.квиз', f"{r} {qid}: blank без accept")
                elif qt == 'order':
                    unknown(q, Q_ORDER, f"{r} {qid}")
                    if len(q.get('items', [])) < 3: bad('5.квиз', f"{r} {qid}: order < 3 элементов")
                else:
                    unknown(q, Q_SINGLE, f"{r} {qid}")
                    n = sum(1 for o in q.get('options', []) if o.get('correct'))
                    if q.get('multiple'):
                        if n < 2: bad('5.квиз', f"{r} {qid}: multiple, но correct={n}")
                    elif n != 1: bad('5.квиз', f"{r} {qid}: correct={n} (нужно ровно 1)")
        elif k == 'terminal':
            unknown(c, TERM, r)
            for t in c.get('tasks', []):
                unknown(t, TASK, f"{r} {t.get('id')}")
                hs = t.get('hints') or []; ex = t.get('expected') or []
                if len(hs) != 3: bad('4.терминал', f"{r} {t.get('id')}: hints={len(hs)} (нужно 3)")
                if not ex and not t.get('pattern'): bad('4.терминал', f"{r} {t.get('id')}: пустой expected")
                if not t.get('success'): bad('4.терминал', f"{r} {t.get('id')}: пустой success")
                if ex and hs and ex[0] not in (hs[-1] + ' ' + (t.get('hint') or '')):
                    bad('4.терминал', f"{r} {t.get('id')}: expected[0] «{ex[0]}» нет в 3-й подсказке/hint")
        elif k == 'code':
            unknown(c, CODE, r)
            if c.get('language') not in LANGS: bad('2.id', f"{r} language={c.get('language')}")
            chs = c.get('checks', []); sol = c.get('solution', '') or ''; st = c.get('starter', '') or ''
            for ch in chs:
                unknown(ch, CHECK, f"{r} check")
                if ch.get('type') not in CTYPES: bad('2.id', f"{r} check.type={ch.get('type')}")
            fails = [ch.get('value') for ch in chs if not ok_check(sol, ch)]
            if fails: bad('3.код', f"{r}: solution НЕ проходит checks: {fails}")
            if chs and all(ok_check(st, ch) for ch in chs): bad('3.код', f"{r}: starter проходит ВСЕ checks")
    dur_by_mod.append((m.get('title') or m.get('_file'), dsum, len(m['lessons'])))

dups = [i for i, n in ids.items() if n > 1 and i is not None]
if dups: bad('2.id', f"дубликаты id: {dups[:5]}")

# ---------- 8. ссылки ----------
RULES = [(r'\bman\b|man7|\(man\)', ('man7.org', 'die.net', 'gnu.org', 'manpages')), (r'\bMDN\b', ('developer.mozilla.org',)),
         (r'PostgreSQL', ('postgresql.org',)), (r'\bDocker\b', ('docker.com',)), (r'\bKubernetes\b', ('kubernetes.io',)),
         (r'официальн', None)]
byurl = defaultdict(set)
for r, t, u in resources:
    byurl[u].add(t)
    dom = re.sub(r'^https?://(www\.)?', '', u).split('/')[0]
    if 'wikipedia.org' in dom and 'Википедия' not in t and 'Wikipedia' not in t:
        bad('8.ссылки', f"{r}: Википедия подписана как «{t}»")
    for pat, doms in RULES:
        if re.search(pat, t):
            ok = ('wikipedia.org' not in dom) if doms is None else any(x in dom for x in doms)
            if not ok: bad('8.ссылки', f"{r}: «{t}» → {dom} (подпись не соответствует домену)")
            break
multi = {u: sorted(ts) for u, ts in byurl.items() if len(ts) > 1}
if NET:
    import urllib.request
    for u in sorted(byurl):
        try:
            req = urllib.request.Request(u, method='HEAD', headers={'User-Agent': 'Mozilla/5.0'})
            code = urllib.request.urlopen(req, timeout=8).getcode()
            if code >= 400: bad('8.ссылки', f"HTTP {code}: {u}")
        except Exception as e:
            bad('8.ссылки', f"недоступна: {u} ({type(e).__name__})")

# ---------- отчёт ----------
print(f"Проверка: {path} ({mode}) — модулей {len(mods)}, уроков {sum(kinds.values())}, по типам {dict(kinds)}")
total = 0
for name in ['1.схема', '2.id', '3.код', '4.терминал', '5.квиз', '6.текст', '7.устаревшее', '8.ссылки']:
    items = problems.get(name, [])
    total += len(items)
    print(f"\n[{name}] нарушений: {len(items)}")
    for it in items[:60]: print("   -", it)
    if len(items) > 60: print(f"   … ещё {len(items)-60}")
print(f"\n[8.ссылки] всего {len(resources)}, уникальных URL {len(byurl)}, URL с разными подписями: {len(multi)}")
print("\n[9.итог] durationMin по модулям:")
for t, d, n in dur_by_mod: print(f"   {d:5d} мин · {n:2d} ур. · {t}")
print(f"   ИТОГО: {sum(d for _, d, _ in dur_by_mod)} мин, уроков {sum(kinds.values())}: {dict(kinds)}")
print(f"\nИТОГО НАРУШЕНИЙ: {total}")
sys.exit(1 if total else 0)
