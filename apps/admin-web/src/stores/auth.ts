import { defineStore } from 'pinia'
import { ref, computed } from 'vue'

export const useAuthStore = defineStore('auth', () => {
  const token = ref<string | null>(localStorage.getItem('adminToken'))
  const user = ref<{ id: string; username: string; role: string; permissions: string[] } | null>(null)

  const isLoggedIn = computed(() => token.value !== null)

  function login(tokenValue: string, userData: { id: string; username: string; role: string; permissions: string[] }) {
    token.value = tokenValue
    user.value = userData
    localStorage.setItem('adminToken', tokenValue)
  }

  function logout() {
    token.value = null
    user.value = null
    localStorage.removeItem('adminToken')
  }

  function hasPermission(permission: string): boolean {
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