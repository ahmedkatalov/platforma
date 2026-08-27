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
