import api from './index'
import type { ApiResponse, Announcement, PaginatedResponse, PaginationParams } from '@/types'

export function getAnnouncementList(params?: PaginationParams) {
  return api.get<ApiResponse<PaginatedResponse<Announcement>>>('/announcements', { params })
}

export function getAnnouncementDetail(id: number) {
  return api.get<ApiResponse<Announcement>>(`/announcements/${id}`)
}

export function createAnnouncement(data: Partial<Announcement>) {
  return api.post<ApiResponse<Announcement>>('/announcements', data)
}

export function updateAnnouncement(id: number, data: Partial<Announcement>) {
  return api.put<ApiResponse<Announcement>>(`/announcements/${id}`, data)
}

export function deleteAnnouncement(id: number) {
  return api.delete<ApiResponse<null>>(`/announcements/${id}`)
}
