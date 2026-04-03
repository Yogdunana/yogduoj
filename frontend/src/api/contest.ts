import api from './index'
import type {
  ApiResponse,
  PaginatedResponse,
  Contest,
  ContestRanking,
  Problem,
  PaginationParams,
  ContestType,
  ContestStatus,
  ContestRuleType,
} from '@/types'

export interface ContestListParams extends PaginationParams {
  keyword?: string
  status?: ContestStatus | ''
  type?: ContestType | ''
  category?: string
  ruleType?: ContestRuleType | ''
}

export function listContests(params?: ContestListParams) {
  return api.get<ApiResponse<PaginatedResponse<Contest>>>('/contests', { params })
}

export function getContest(id: number) {
  return api.get<ApiResponse<Contest>>(`/contests/${id}`)
}

export function signupContest(id: number) {
  return api.post<ApiResponse<null>>(`/contests/${id}/signup`)
}

export function withdrawContest(id: number) {
  return api.post<ApiResponse<null>>(`/contests/${id}/withdraw`)
}

export function getContestProblems(id: number) {
  return api.get<ApiResponse<Problem[]>>(`/contests/${id}/problems`)
}

export function submitToContest(contestId: number, data: { problemId: number; language: string; code: string; flag?: string }) {
  return api.post<ApiResponse<{ submissionId: number }>>(`/contests/${contestId}/submit`, data)
}

export function getContestRanking(id: number, params?: PaginationParams) {
  return api.get<ApiResponse<PaginatedResponse<ContestRanking>>>(`/contests/${id}/ranking`, { params })
}

export function getContestFrozenRanking(id: number) {
  return api.get<ApiResponse<ContestRanking[]>>(`/contests/${id}/ranking/frozen`)
}

// Admin APIs
export function adminCreateContest(data: Partial<Contest>) {
  return api.post<ApiResponse<Contest>>('/admin/contests', data)
}

export function adminUpdateContest(id: number, data: Partial<Contest>) {
  return api.put<ApiResponse<Contest>>(`/admin/contests/${id}`, data)
}

export function adminDeleteContest(id: number) {
  return api.delete<ApiResponse<null>>(`/admin/contests/${id}`)
}

export function adminUpdateContestStatus(id: number, data: { status: ContestStatus }) {
  return api.put<ApiResponse<null>>(`/admin/contests/${id}/status`, data)
}

export function adminAddContestProblem(contestId: number, data: { problemId: number; label?: string }) {
  return api.post<ApiResponse<null>>(`/admin/contests/${contestId}/problems`, data)
}

export function adminRemoveContestProblem(contestId: number, problemId: number) {
  return api.delete<ApiResponse<null>>(`/admin/contests/${contestId}/problems/${problemId}`)
}

export function adminGetContestSignups(id: number, params?: PaginationParams) {
  return api.get<ApiResponse<PaginatedResponse<any>>>(`/admin/contests/${id}/signups`, { params })
}

// Aliases for backward compatibility
export const getContestDetail = getContest
export const getContestList = listContests
