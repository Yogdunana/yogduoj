import api from './index'
import type { ApiResponse, Submission, PaginatedResponse, PaginationParams } from '@/types'

export function getSubmissionList(params?: PaginationParams & { problemId?: number; userId?: number; status?: string }) {
  return api.get<ApiResponse<PaginatedResponse<Submission>>>('/submissions', { params })
}

export function getSubmissionDetail(id: number) {
  return api.get<ApiResponse<Submission>>(`/submissions/${id}`)
}

export function submitCode(data: { problemId: number; language: string; code: string }) {
  return api.post<ApiResponse<Submission>>('/submissions', data)
}

export function rejudgeSubmission(id: number) {
  return api.post<ApiResponse<Submission>>(`/submissions/${id}/rejudge`)
}
