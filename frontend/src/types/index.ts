// User types
export interface User {
  id: number
  username: string
  email: string
  role: UserRole
  avatar?: string
  bio?: string
  school?: string
  createdAt: string
  updatedAt: string
}

export type UserRole = 'admin' | 'user' | 'team_leader'

// Team types
export interface Team {
  id: number
  name: string
  description: string
  leaderId: number
  leader?: User
  memberCount: number
  createdAt: string
  updatedAt: string
}

export interface TeamMember {
  id: number
  teamId: number
  userId: number
  user?: User
  role: 'leader' | 'member'
  joinedAt: string
}

// Problem types
export type ProblemType = 'programming' | 'algorithm' | 'ctf'

export interface Problem {
  id: number
  title: string
  description: string
  difficulty: ProblemDifficulty
  type: ProblemType
  tags: string[]
  timeLimit: number
  memoryLimit: number
  totalSubmit: number
  acceptedCount: number
  isPublic: boolean
  inputFormat?: string
  outputFormat?: string
  samples?: ProblemSample[]
  hints?: string[]
  attachments?: string[]
  userStatus?: ProblemUserStatus
  createdAt: string
  updatedAt: string
}

export type ProblemDifficulty = 'easy' | 'medium' | 'hard' | 'expert'
export type ProblemUserStatus = 'unsubmitted' | 'submitted' | 'accepted'

export interface ProblemSample {
  input: string
  output: string
}

export interface ProblemTestCase {
  id: number
  input: string
  output: string
  isSample: boolean
}

// Submission types
export interface Submission {
  id: number
  problemId: number
  problem?: Problem
  userId: number
  user?: User
  language: string
  status: SubmissionStatus
  timeUsed?: number
  memoryUsed?: number
  score?: number
  totalScore?: number
  code?: string
  errorMessage?: string
  judgeDetail?: JudgeTestCase[]
  createdAt: string
}

export type SubmissionStatus =
  | 'pending'
  | 'judging'
  | 'accepted'
  | 'wrong_answer'
  | 'time_limit_exceeded'
  | 'memory_limit_exceeded'
  | 'runtime_error'
  | 'compilation_error'
  | 'presentation_error'
  | 'system_error'

export interface JudgeTestCase {
  id: number
  status: SubmissionStatus
  timeUsed?: number
  memoryUsed?: number
  score?: number
  signal?: string
  errorMessage?: string
}

// Contest types
export interface Contest {
  id: number
  title: string
  description: string
  startTime: string
  endTime: string
  type: ContestType
  status: ContestStatus
  ruleType: ContestRuleType
  problemIds: number[]
  problems?: Problem[]
  participantCount: number
  createdBy: number
  creator?: User
  createdAt: string
  updatedAt: string
}

export type ContestType = 'icpc' | 'ioi' | 'ctf'
export type ContestStatus = 'upcoming' | 'running' | 'ended'
export type ContestRuleType = 'acm' | 'oi' | 'cf'

export interface ContestRanking {
  userId: number
  user?: User
  rank: number
  solvedCount: number
  penalty: number
  problemResults: Record<number, ProblemResult>
}

export interface ProblemResult {
  status: SubmissionStatus
  timeUsed: number
  memoryUsed: number
  submitCount: number
  acceptedAt?: string
}

// Announcement types
export interface Announcement {
  id: number
  title: string
  content: string
  authorId: number
  author?: User
  isPinned: boolean
  createdAt: string
  updatedAt: string
}

// Generic API response types
export interface ApiResponse<T = unknown> {
  code: number
  message: string
  data: T
}

export interface PaginatedResponse<T> {
  items: T[]
  total: number
  page: number
  pageSize: number
  totalPages: number
}

// Auth types
export interface LoginRequest {
  username: string
  password: string
}

export interface RegisterRequest {
  username: string
  email: string
  password: string
}

export interface AuthResponse {
  token: string
  refreshToken: string
  user: User
}

// CTF types
export interface CTFChallenge {
  id: number
  title: string
  description: string
  category: CTFCategory
  points: number
  solvedCount: number
  isSolved: boolean
  hint?: string
  attachments?: string[]
}

export type CTFCategory =
  | 'web'
  | 'crypto'
  | 'pwn'
  | 'reverse'
  | 'misc'
  | 'forensics'
  | 'recon'
  | 'vuln-reproduce'

// Pagination params
export interface PaginationParams {
  page?: number
  pageSize?: number
}
