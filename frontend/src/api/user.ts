import api from './index'
import type { ApiResponse, User, PaginatedResponse, PaginationParams } from '@/types'

export function getUserList(params?: PaginationParams & { keyword?: string }) {
  return api.get<ApiResponse<PaginatedResponse<User>>>('/users', { params })
}

export function getUserProfile(id: number) {
  return api.get<ApiResponse<User>>(`/users/${id}`)
}

export function updateUserProfile(id: number, data: Partial<User>) {
  return api.put<ApiResponse<User>>(`/users/${id}`, data)
}

export function deleteUser(id: number) {
  return api.delete<ApiResponse<null>>(`/users/${id}`)
}
