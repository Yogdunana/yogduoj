import api from './index'
import type { ApiResponse, PaginatedResponse, Problem, PaginationParams, ProblemType, ProblemDifficulty } from '@/types'

export interface ProblemListParams extends PaginationParams {
  keyword?: string
  difficulty?: ProblemDifficulty
  type?: ProblemType
  tag?: string
  sortBy?: 'id' | 'acceptanceRate' | 'submissions' | 'difficulty'
  sortOrder?: 'asc' | 'desc'
}

export function listProblems(params?: ProblemListParams) {
  return api.get<ApiResponse<PaginatedResponse<Problem>>>('/problems', { params })
}

export function getProblem(id: number) {
  return api.get<ApiResponse<Problem>>(`/problems/${id}`)
}

export function getProblemSamples(id: number) {
  return api.get<ApiResponse<{ input: string; output: string }[]>>(`/problems/${id}/samples`)
}

// Admin APIs
export function adminCreateProblem(data: Partial<Problem>) {
  return api.post<ApiResponse<Problem>>('/admin/problems', data)
}

export function adminUpdateProblem(id: number, data: Partial<Problem>) {
  return api.put<ApiResponse<Problem>>(`/admin/problems/${id}`, data)
}

export function adminDeleteProblem(id: number) {
  return api.delete<ApiResponse<null>>(`/admin/problems/${id}`)
}

export function adminUploadTestData(id: number, formData: FormData) {
  return api.post<ApiResponse<null>>(`/admin/problems/${id}/test-data`, formData, {
    headers: { 'Content-Type': 'multipart/form-data' },
  })
}

export function adminDeleteTestData(id: number, dataId: number) {
  return api.delete<ApiResponse<null>>(`/admin/problems/${id}/test-data/${dataId}`)
}
