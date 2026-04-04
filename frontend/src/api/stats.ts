import api from './index'
import type { ApiResponse } from '@/types'

export interface PlatformStats {
  total_users: number
  total_problems: number
  total_submissions: number
  total_contests: number
  accepted_submissions: number
}

export function getStats() {
  return api.get<ApiResponse<PlatformStats>>('/v1/stats')
}
