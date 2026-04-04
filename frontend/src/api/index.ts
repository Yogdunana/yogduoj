import axios from 'axios'
import type { ApiResponse } from '@/types'
import router from '@/router'

// Convert snake_case keys to camelCase recursively
function toCamelCase(obj: unknown): unknown {
  if (obj === null || obj === undefined) return obj
  if (Array.isArray(obj)) return obj.map(toCamelCase)
  if (typeof obj === 'object') {
    const result: Record<string, unknown> = {}
    for (const [key, value] of Object.entries(obj as Record<string, unknown>)) {
      const camelKey = key.replace(/_([a-z])/g, (_, c) => c.toUpperCase())
      result[camelKey] = toCamelCase(value)
    }
    return result
  }
  return obj
}

// Convert camelCase keys to snake_case recursively (for request params)
function toSnakeCase(obj: unknown): unknown {
  if (obj === null || obj === undefined) return obj
  if (Array.isArray(obj)) return obj.map(toSnakeCase)
  if (typeof obj === 'object') {
    const result: Record<string, unknown> = {}
    for (const [key, value] of Object.entries(obj as Record<string, unknown>)) {
      const snakeKey = key.replace(/[A-Z]/g, (c) => `_${c.toLowerCase()}`)
      result[snakeKey] = toSnakeCase(value)
    }
    return result
  }
  return obj
}

const api = axios.create({
  baseURL: import.meta.env.VITE_API_BASE_URL || '/api/v1',
  timeout: 15000,
  headers: {
    'Content-Type': 'application/json',
  },
})

// Request interceptor - add JWT token & convert camelCase to snake_case
api.interceptors.request.use(
  (config) => {
    const token = localStorage.getItem('token')
    if (token) {
      config.headers.Authorization = `Bearer ${token}`
    }
    if (config.params) {
      config.params = toSnakeCase(config.params) as Record<string, unknown>
    }
    if (config.data && typeof config.data === 'object') {
      config.data = toSnakeCase(config.data)
    }
    return config
  },
  (error) => {
    return Promise.reject(error)
  }
)

// Response interceptor - handle errors & convert snake_case to camelCase
api.interceptors.response.use(
  (response) => {
    response.data = toCamelCase(response.data)
    return response
  },
  async (error) => {
    const status = error.response?.status
    const message = error.response?.data?.message || error.message

    if (status === 401) {
      localStorage.removeItem('token')
      localStorage.removeItem('refreshToken')
      localStorage.removeItem('user')
      if (router.currentRoute.value.name !== 'Login') {
        router.push({ name: 'Login', query: { redirect: router.currentRoute.value.fullPath } })
      }
    }

    return Promise.reject(new Error(message))
  }
)

export default api
export type { ApiResponse }
