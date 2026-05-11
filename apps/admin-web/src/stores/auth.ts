import { defineStore } from 'pinia'
import { ref, computed } from 'vue'

const tokenStorageKey = 'adminToken'
const userStorageKey = 'adminUser'

interface AuthUser {
  id: string
  username: string
  role: string
  permissions: string[]
}

function loadStoredUser(): AuthUser | null {
  const raw = localStorage.getItem(userStorageKey)
  if (!raw) {
    return null
  }

  try {
    return JSON.parse(raw) as AuthUser
  } catch (error) {
    console.error('Failed to parse stored admin user:', error)
    localStorage.removeItem(userStorageKey)
    return null
  }
}

export const useAuthStore = defineStore('auth', () => {
  const token = ref<string | null>(localStorage.getItem(tokenStorageKey))
  const user = ref<AuthUser | null>(loadStoredUser())

  const isLoggedIn = computed(() => token.value !== null)

  function login(tokenValue: string, userData: AuthUser) {
    token.value = tokenValue
    user.value = userData
    localStorage.setItem(tokenStorageKey, tokenValue)
    localStorage.setItem(userStorageKey, JSON.stringify(userData))
  }

  function logout() {
    token.value = null
    user.value = null
    localStorage.removeItem(tokenStorageKey)
    localStorage.removeItem(userStorageKey)
  }

  function hasPermission(permission: string): boolean {
    if (!permission) {
      return true
    }
    return user.value?.permissions.includes(permission) ?? false
  }

  return {
    token,
    user,
    isLoggedIn,
    login,
    logout,
    hasPermission
  }
})
