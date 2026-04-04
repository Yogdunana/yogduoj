import axios from 'axios'
import type { ApiResponse } from '@/types'
import router from '@/router'

const api = axios.create({
  baseURL: import.meta.env.VITE_API_BASE_URL || '/api/v1',
  timeout: 15000,
  headers: {
    'Content-Type': 'application/json',
  },
})

// Request interceptor - add JWT token
api.interceptors.request.use(
  (config) => {
    const token = localStorage.getItem('token')
    if (token) {
      config.headers.Authorization = `Bearer ${token}`
    }
    return config
  },
  (error) => {
    return Promise.reject(error)
  }
)

// Response interceptor - handle errors
api.interceptors.response.use(
  (response) => {
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
