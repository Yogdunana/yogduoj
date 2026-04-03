import { defineStore } from 'pinia'
import { ref, computed } from 'vue'
import type { User } from '@/types'
import { login as apiLogin, register as apiRegister, logout as apiLogout, getProfile } from '@/api/auth'

export const useAuthStore = defineStore('auth', () => {
  const user = ref<User | null>(null)
  const token = ref<string>(localStorage.getItem('token') || '')
  const refreshTokenValue = ref<string>(localStorage.getItem('refreshToken') || '')

  const isLoggedIn = computed(() => !!token.value)
  const isAdmin = computed(() => user.value?.role === 'admin')

  async function login(username: string, password: string) {
    const res = await apiLogin({ username, password })
    const data = res.data.data
    token.value = data.token
    refreshTokenValue.value = data.refreshToken
    user.value = data.user
    localStorage.setItem('token', data.token)
    localStorage.setItem('refreshToken', data.refreshToken)
    localStorage.setItem('user', JSON.stringify(data.user))
  }

  async function register(username: string, email: string, password: string) {
    const res = await apiRegister({ username, email, password })
    const data = res.data.data
    token.value = data.token
    refreshTokenValue.value = data.refreshToken
    user.value = data.user
    localStorage.setItem('token', data.token)
    localStorage.setItem('refreshToken', data.refreshToken)
    localStorage.setItem('user', JSON.stringify(data.user))
  }

  async function logout() {
    try {
      await apiLogout()
    } catch {
      // ignore logout errors
    }
    user.value = null
    token.value = ''
    refreshTokenValue.value = ''
    localStorage.removeItem('token')
    localStorage.removeItem('refreshToken')
    localStorage.removeItem('user')
  }

  async function fetchProfile() {
    const res = await getProfile()
    user.value = res.data.data
    localStorage.setItem('user', JSON.stringify(res.data.data))
  }

  // Initialize user from localStorage
  function init() {
    const storedUser = localStorage.getItem('user')
    if (storedUser && token.value) {
      try {
        user.value = JSON.parse(storedUser)
      } catch {
        user.value = null
      }
    }
  }

  init()

  return {
    user,
    token,
    isLoggedIn,
    isAdmin,
    login,
    register,
    logout,
    fetchProfile,
  }
})
