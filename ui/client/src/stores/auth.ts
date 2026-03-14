import { defineStore } from 'pinia'
import { ref, computed } from 'vue'

const STORAGE_TOKEN = 'beehive-access-token'
const STORAGE_USER_ID = 'beehive-user-id'

export interface AuthUser {
  id: string
  nickname?: string
  avatarUrl?: string
}

export const useAuthStore = defineStore('auth', () => {
  const token = localStorage.getItem(STORAGE_TOKEN)
  const savedUserId = localStorage.getItem(STORAGE_USER_ID)
  const user = ref<AuthUser | null>(savedUserId ? { id: savedUserId } : null)
  const accessToken = ref<string | null>(token)
  const refreshToken = ref<string | null>(null)

  const loggedIn = computed(() => !!accessToken.value)

  function setFromTokenLogin(payload: { userId: string; accessToken?: string; refreshToken?: string }) {
    user.value = { id: payload.userId }
    if (payload.accessToken != null) {
      accessToken.value = payload.accessToken
      localStorage.setItem(STORAGE_TOKEN, payload.accessToken)
    }
    if (payload.refreshToken != null) refreshToken.value = payload.refreshToken
    if (payload.userId) localStorage.setItem(STORAGE_USER_ID, payload.userId)
  }

  function setUser(profile: AuthUser) {
    user.value = profile
  }

  function logout() {
    user.value = null
    accessToken.value = null
    refreshToken.value = null
    localStorage.removeItem(STORAGE_TOKEN)
    localStorage.removeItem(STORAGE_USER_ID)
  }

  return {
    user,
    accessToken,
    refreshToken,
    loggedIn,
    setFromTokenLogin,
    setUser,
    logout,
  }
})
