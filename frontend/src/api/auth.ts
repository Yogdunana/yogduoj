import api from './index'
import type { ApiResponse, AuthResponse, LoginRequest, RegisterRequest, User } from '@/types'

export function login(data: LoginRequest) {
  return api.post<ApiResponse<AuthResponse>>('/auth/login', data)
}

export function register(data: RegisterRequest) {
  return api.post<ApiResponse<AuthResponse>>('/auth/register', data)
}

export function logout() {
  return api.post<ApiResponse<null>>('/auth/logout')
}

export function refreshToken() {
  const refreshTokenValue = localStorage.getItem('refreshToken')
  return api.post<ApiResponse<AuthResponse>>('/auth/refresh', {
    refreshToken: refreshTokenValue,
  })
}

export function getProfile() {
  return api.get<ApiResponse<User>>('/auth/profile')
}
