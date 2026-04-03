import api from './index'
import type { ApiResponse, Submission, PaginatedResponse, PaginationParams, SubmissionStatus } from '@/types'

export interface SubmissionListParams extends PaginationParams {
  problemId?: number
  userId?: number
  status?: SubmissionStatus
  language?: string
  startTime?: string
  endTime?: string
}

export function createSubmission(data: { problemId: number; language: string; code: string }) {
  return api.post<ApiResponse<Submission>>('/submissions', data)
}

export function getSubmission(id: number) {
  return api.get<ApiResponse<Submission>>(`/submissions/${id}`)
}

export function getSubmissionCode(id: number) {
  return api.get<ApiResponse<{ code: string }>>(`/submissions/${id}/code`)
}

export function listSubmissions(params?: SubmissionListParams) {
  return api.get<ApiResponse<PaginatedResponse<Submission>>>('/submissions', { params })
}

export function adminRejudge(id: number) {
  return api.post<ApiResponse<Submission>>(`/admin/submissions/${id}/rejudge`)
}
