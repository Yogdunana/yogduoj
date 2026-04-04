import api from './index'
import type { ApiResponse } from '@/types'

export interface PlatformStats {
  totalUsers: number
  totalProblems: number
  totalSubmissions: number
  totalContests: number
  acceptedSubmissions: number
}

export function getStats() {
  return api.get<ApiResponse<PlatformStats>>('/stats')
}
