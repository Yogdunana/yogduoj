import api from './index'
import type { ApiResponse, Announcement, PaginatedResponse, PaginationParams } from '@/types'

export interface AnnouncementListParams extends PaginationParams {
  keyword?: string
}

export function listAnnouncements(params?: AnnouncementListParams) {
  return api.get<ApiResponse<PaginatedResponse<Announcement>>>('/announcements', { params })
}

export function getAnnouncement(id: number) {
  return api.get<ApiResponse<Announcement>>(`/announcements/${id}`)
}

// Admin APIs
export function adminCreateAnnouncement(data: Partial<Announcement>) {
  return api.post<ApiResponse<Announcement>>('/admin/announcements', data)
}

export function adminUpdateAnnouncement(id: number, data: Partial<Announcement>) {
  return api.put<ApiResponse<Announcement>>(`/admin/announcements/${id}`, data)
}

export function adminDeleteAnnouncement(id: number) {
  return api.delete<ApiResponse<null>>(`/admin/announcements/${id}`)
}

// Aliases for backward compatibility
export const createAnnouncement = adminCreateAnnouncement
export const updateAnnouncement = adminUpdateAnnouncement
export const deleteAnnouncement = adminDeleteAnnouncement
export const getAnnouncementDetail = getAnnouncement
