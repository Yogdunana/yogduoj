import api from './index'
import type { ApiResponse, Team, TeamMember, PaginatedResponse, PaginationParams } from '@/types'

export function getTeamList(params?: PaginationParams & { keyword?: string }) {
  return api.get<ApiResponse<PaginatedResponse<Team>>>('/teams', { params })
}

export function getTeamDetail(id: number) {
  return api.get<ApiResponse<Team>>(`/teams/${id}`)
}

export function getTeamMembers(teamId: number, params?: PaginationParams) {
  return api.get<ApiResponse<PaginatedResponse<TeamMember>>>(`/teams/${teamId}/members`, { params })
}

export function createTeam(data: Partial<Team>) {
  return api.post<ApiResponse<Team>>('/teams', data)
}

export function updateTeam(id: number, data: Partial<Team>) {
  return api.put<ApiResponse<Team>>(`/teams/${id}`, data)
}

export function deleteTeam(id: number) {
  return api.delete<ApiResponse<null>>(`/teams/${id}`)
}

export function joinTeam(teamId: number) {
  return api.post<ApiResponse<TeamMember>>(`/teams/${teamId}/join`)
}

export function leaveTeam(teamId: number) {
  return api.post<ApiResponse<null>>(`/teams/${teamId}/leave`)
}
