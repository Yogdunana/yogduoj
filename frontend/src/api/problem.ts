import api from './index'
import type { ApiResponse, PaginatedResponse, Problem, PaginationParams } from '@/types'

export function getProblemList(params?: PaginationParams & { keyword?: string; difficulty?: string; tag?: string }) {
  return api.get<ApiResponse<PaginatedResponse<Problem>>>('/problems', { params })
}

export function getProblemDetail(id: number) {
  return api.get<ApiResponse<Problem>>(`/problems/${id}`)
}

export function createProblem(data: Partial<Problem>) {
  return api.post<ApiResponse<Problem>>('/problems', data)
}

export function updateProblem(id: number, data: Partial<Problem>) {
  return api.put<ApiResponse<Problem>>(`/problems/${id}`, data)
}

export function deleteProblem(id: number) {
  return api.delete<ApiResponse<null>>(`/problems/${id}`)
}
