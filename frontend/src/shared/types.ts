// Общие типы, повторяющие ответы Go-бэкенда.

export type Role = "admin" | "student";
export type UserStatus = "invited" | "active" | "blocked";

export type User = {
  id: string;
  email: string;
  fullName: string;
  role: Role;
  status: UserStatus;
  emailVerified: boolean;
  avatarUrl: string;
  lastLoginAt: string | null;
  createdAt: string;
  updatedAt: string;
};

export type Session = {
  accessToken: string;
  refreshToken: string;
  expiresAt: string;
  user: User;
};

export type CourseLevel = "beginner" | "intermediate" | "advanced";
export type CourseStatus = "draft" | "published" | "archived";
export type LessonKind = "text" | "quiz" | "terminal" | "code";

export type Lesson = {
  id: string;
  moduleId: string;
  title: string;
  kind: LessonKind;
  summary: string;
  content: Record<string, unknown>;
  durationMin: number;
  position: number;
  createdAt: string;
  updatedAt: string;
};

export type Module = {
  id: string;
  courseId: string;
  title: string;
  summary: string;
  position: number;
  lessons?: Lesson[];
  createdAt: string;
  updatedAt: string;
};

export type Course = {
  id: string;
  slug: string;
  title: string;
  subtitle: string;
  description: string;
  coverUrl: string;
  level: CourseLevel;
  tags: string[];
  status: CourseStatus;
  position: number;
  modules?: Module[];
  modulesCount: number;
  lessonsCount: number;
  studentsCount: number;
  createdAt: string;
  updatedAt: string;
};

export type Enrollment = {
  id: string;
  userId: string;
  courseId: string;
  status: "active" | "completed" | "paused";
  dueDate: string | null;
  startedAt: string | null;
  completedAt: string | null;
  createdAt: string;
  courseTitle?: string;
  courseSlug?: string;
};

export type ActivityDay = {
  day: string;
  visits: number;
  secondsSpent: number;
};

export type StudentSummary = {
  userId: string;
  email: string;
  fullName: string;
  status: UserStatus;
  lastLoginAt: string | null;
  courses: number;
  lessonsTotal: number;
  lessonsCompleted: number;
  daysVisited: number;
  minutesSpent: number;
  progress: number;
};

export type AdminOverview = {
  students: number;
  activeStudents: number;
  blockedStudents: number;
  invitedStudents: number;
  admins: number;
  courses: number;
  publishedCourses: number;
  lessons: number;
  enrollments: number;
  activeToday: number;
  activeWeek: number;
};

export type AuditEntry = {
  id: number;
  actorId: string | null;
  actorName: string;
  action: string;
  entity: string;
  entityId: string;
  payload: Record<string, unknown>;
  createdAt: string;
};

export type Paginated<T> = {
  items: T[];
  total: number;
  page: number;
  limit: number;
};

export type CreatedStudent = {
  user: User;
  tempPassword: string;
  mailSent: boolean;
  mailError?: string;
};

// --- Прохождение уроков ---

export type LessonProgress = {
  lessonId: string;
  status: "in_progress" | "completed";
  score: number | null;
  bestScore: number | null;
  attempts: number;
  secondsSpent: number;
  completedAt: string | null;
  updatedAt: string;
};

export type TaskState = {
  taskId: string;
  attempts: number;
  hintsUsed: number;
  completedAt: string | null;
};

export type LessonView = {
  lesson: Lesson;
  courseId: string;
  courseSlug: string;
  courseTitle: string;
  moduleTitle: string;
  prevLessonId: string | null;
  nextLessonId: string | null;
  progress: LessonProgress[];
  tasks: TaskState[];
};

// Содержимое уроков (правильные ответы вырезаны на сервере).
export type QuizOption = { id: string; text: string };

export type QuizQuestion = {
  id: string;
  text: string;
  hint?: string;
  multiple?: boolean;
  /** Вопрос на повторение ранее пройденной темы. */
  review?: boolean;
  options: QuizOption[];
};

export type QuizContent = {
  intro?: string;
  resources?: LessonResource[];
  passScore?: number;
  timeLimitSec?: number;
  shuffle?: boolean;
  questions: QuizQuestion[];
};

export type TerminalTask = {
  id: string;
  prompt: string;
  hint?: string;
  success?: string;
};

export type TerminalContent = {
  intro?: string;
  shell?: string;
  tasks: TerminalTask[];
  resources?: LessonResource[];
};

export type CodeContent = {
  language: string;
  task: string;
  starter: string;
  hint?: string;
  resources?: LessonResource[];
};

export type LessonResource = { title: string; url: string; note?: string };

export type TextContent = { body?: string; resources?: LessonResource[] };

export type QuestionResult = {
  questionId: string;
  correct: boolean;
  correctOptionIds: string[];
  chosenOptionIds: string[];
  explanation?: string;
};

export type QuizResult = {
  score: number;
  passed: boolean;
  correctCount: number;
  totalCount: number;
  passScore: number;
  questions: QuestionResult[];
  certificate?: Certificate | null;
};

export type TerminalCheckResult = {
  solved: boolean;
  message: string;
  hint?: string;
  completedTasks: string[];
  lessonComplete: boolean;
  certificate?: Certificate | null;
};

export type CodeCheckResult = {
  passed: boolean;
  score: number;
  checks: { ok: boolean; message: string }[];
  hint?: string;
  certificate?: Certificate | null;
};

export type Attempt = {
  id: string;
  lessonId: string;
  lessonTitle: string;
  courseTitle: string;
  kind: "quiz" | "terminal" | "code";
  score: number;
  correctCount: number;
  totalCount: number;
  passed: boolean;
  durationSeconds: number;
  details: Record<string, unknown>;
  createdAt: string;
};

export type QuizStats = {
  attempts: number;
  passed: number;
  averageScore: number;
  bestScore: number;
  accuracy: number;
  avgSecondsPerQuestion: number;
  fastestSeconds: number;
  answeredTotal: number;
  answeredCorrect: number;
};

// --- Сертификаты и файлы ---

export type Certificate = {
  id: string;
  serial: string;
  userId: string;
  courseId: string;
  holderName: string;
  courseTitle: string;
  score: number;
  lessonsTotal: number;
  lessonsCompleted: number;
  revokedAt: string | null;
  issuedAt: string;
};

export type CertificateCheck = {
  valid: boolean;
  message?: string;
  certificate?: Certificate;
};

export type Asset = {
  id: string;
  filename: string;
  original: string;
  mime: string;
  sizeBytes: number;
  url: string;
  createdAt: string;
};

export type QuizCard = {
  lessonId: string;
  title: string;
  summary: string;
  courseId: string;
  courseSlug: string;
  courseTitle: string;
  moduleTitle: string;
  position: number;
  questions: number;
  passScore: number;
  durationMin: number;
  status: "not_started" | "in_progress" | "completed";
  bestScore: number | null;
  attempts: number;
  lastTriedAt: string | null;
};

export type Note = {
  id: string;
  lessonId: string;
  lessonTitle: string;
  lessonKind: LessonKind;
  moduleTitle: string;
  courseSlug: string;
  courseTitle: string;
  quote: string;
  body: string;
  createdAt: string;
  updatedAt: string;
};
