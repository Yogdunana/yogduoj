import api from './index'
import type { ApiResponse, PaginatedResponse, User, Problem, Contest, Submission, PaginationParams } from '@/types'

// Admin - User management
export function adminGetUsers(params?: PaginationParams & { keyword?: string }) {
  return api.get<ApiResponse<PaginatedResponse<User>>>('/admin/users', { params })
}

export function adminUpdateUser(id: number, data: Partial<User>) {
  return api.put<ApiResponse<User>>(`/admin/users/${id}`, data)
}

export function adminDeleteUser(id: number) {
  return api.delete<ApiResponse<null>>(`/admin/users/${id}`)
}

// Admin - Problem management
export function adminGetProblems(params?: PaginationParams) {
  return api.get<ApiResponse<PaginatedResponse<Problem>>>('/admin/problems', { params })
}

export function adminCreateProblem(data: Partial<Problem>) {
  return api.post<ApiResponse<Problem>>('/admin/problems', data)
}

export function adminUpdateProblem(id: number, data: Partial<Problem>) {
  return api.put<ApiResponse<Problem>>(`/admin/problems/${id}`, data)
}

export function adminDeleteProblem(id: number) {
  return api.delete<ApiResponse<null>>(`/admin/problems/${id}`)
}

// Admin - Contest management
export function adminGetContests(params?: PaginationParams) {
  return api.get<ApiResponse<PaginatedResponse<Contest>>>('/admin/contests', { params })
}

export function adminCreateContest(data: Partial<Contest>) {
  return api.post<ApiResponse<Contest>>('/admin/contests', data)
}

export function adminUpdateContest(id: number, data: Partial<Contest>) {
  return api.put<ApiResponse<Contest>>(`/admin/contests/${id}`, data)
}

export function adminDeleteContest(id: number) {
  return api.delete<ApiResponse<null>>(`/admin/contests/${id}`)
}

// Admin - Submission management
export function adminGetSubmissions(params?: PaginationParams) {
  return api.get<ApiResponse<PaginatedResponse<Submission>>>('/admin/submissions', { params })
}

export function adminRejudge(id: number) {
  return api.post<ApiResponse<null>>(`/admin/submissions/${id}/rejudge`)
}

// Admin - Judge monitor
export function getJudgeStatus() {
  return api.get<ApiResponse<unknown>>('/admin/judge/status')
}

// Admin - Announcements
export function adminGetAnnouncements(params?: PaginationParams) {
  return api.get<ApiResponse<PaginatedResponse<unknown>>>('/admin/announcements', { params })
}

// Admin - System config
export function getSystemConfig() {
  return api.get<ApiResponse<Record<string, unknown>>>('/admin/config')
}

export function updateSystemConfig(data: Record<string, unknown>) {
  return api.put<ApiResponse<Record<string, unknown>>>('/admin/config', data)
}

// Admin - Import
export function importProblems(file: File) {
  const formData = new FormData()
  formData.append('file', file)
  return api.post<ApiResponse<null>>('/admin/import/problems', formData, {
    headers: { 'Content-Type': 'multipart/form-data' },
  })
}
