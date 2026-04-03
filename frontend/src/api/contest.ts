import api from './index'
import type { ApiResponse, PaginatedResponse, Contest, ContestRanking, PaginationParams } from '@/types'

export function getContestList(params?: PaginationParams & { status?: string }) {
  return api.get<ApiResponse<PaginatedResponse<Contest>>>('/contests', { params })
}

export function getContestDetail(id: number) {
  return api.get<ApiResponse<Contest>>(`/contests/${id}`)
}

export function getContestRanking(id: number, params?: PaginationParams) {
  return api.get<ApiResponse<PaginatedResponse<ContestRanking>>>(`/contests/${id}/ranking`, { params })
}

export function createContest(data: Partial<Contest>) {
  return api.post<ApiResponse<Contest>>('/contests', data)
}

export function updateContest(id: number, data: Partial<Contest>) {
  return api.put<ApiResponse<Contest>>(`/contests/${id}`, data)
}

export function deleteContest(id: number) {
  return api.delete<ApiResponse<null>>(`/contests/${id}`)
}
