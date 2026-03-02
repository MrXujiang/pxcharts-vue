import { defineStore } from 'pinia'
import { ref, computed } from 'vue'
import { 
  setAuthData, 
  getToken, 
  getUserInfo as getStoredUserInfo, 
  clearAuthData,
  updateUserInfo as updateStoredUserInfo,
  type UserInfo 
} from '@/utils/auth'

export const useUserStore = defineStore('user', () => {
  // 状态
  const token = ref<string | null>(getToken())
  const userInfo = ref<UserInfo | null>(getStoredUserInfo())

  // 计算属性
  const isLoggedIn = computed(() => !!token.value)
  const userId = computed(() => userInfo.value?.id)
  const userEmail = computed(() => userInfo.value?.email)
  const userName = computed(() => userInfo.value?.username || userInfo.value?.email)

  /**
   * 登录
   */
  const login = (authToken: string, user: UserInfo) => {
    token.value = authToken
    userInfo.value = user
    setAuthData({ token: authToken, userInfo: user })
  }

  /**
   * 更新用户信息
   */
  const updateUserInfo = (info: Partial<UserInfo>) => {
    if (userInfo.value) {
      userInfo.value = { ...userInfo.value, ...info }
      updateStoredUserInfo(info)
    }
  }

  /**
   * 退出登录
   */
  const logout = () => {
    token.value = null
    userInfo.value = null
    clearAuthData()
  }

  /**
   * 初始化用户信息（从 localStorage 恢复）
   */
  const initUserInfo = () => {
    token.value = getToken()
    userInfo.value = getStoredUserInfo()
  }

  return {
    // 状态
    token,
    userInfo,
    // 计算属性
    isLoggedIn,
    userId,
    userEmail,
    userName,
    // 方法
    login,
    updateUserInfo,
    logout,
    initUserInfo,
  }
})
