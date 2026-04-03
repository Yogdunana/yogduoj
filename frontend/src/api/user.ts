import api from './index'
import type { ApiResponse, User, Submission, Contest, PaginatedResponse, PaginationParams } from '@/types'

export function getUserList(params?: PaginationParams & { keyword?: string }) {
  return api.get<ApiResponse<PaginatedResponse<User>>>('/users', { params })
}

export function getUserProfile(id: number) {
  return api.get<ApiResponse<User>>(`/users/${id}`)
}

export function getUserPublicProfile(id: number) {
  return api.get<ApiResponse<User>>(`/users/${id}/public`)
}

export function updateUserProfile(id: number, data: Partial<User>) {
  return api.put<ApiResponse<User>>(`/users/${id}`, data)
}

export function updateProfile(data: { username?: string; email?: string; avatar?: string }) {
  return api.put<ApiResponse<User>>('/user/profile', data)
}

export function changePassword(data: { oldPassword: string; newPassword: string }) {
  return api.put<ApiResponse<null>>('/user/password', data)
}

export function deleteUser(id: number) {
  return api.delete<ApiResponse<null>>(`/users/${id}`)
}

export function getUserSubmissions(params?: PaginationParams & { userId?: number }) {
  return api.get<ApiResponse<PaginatedResponse<Submission>>>('/user/submissions', { params })
}

export function getUserContests(params?: PaginationParams & { userId?: number }) {
  return api.get<ApiResponse<PaginatedResponse<Contest>>>('/user/contests', { params })
}
