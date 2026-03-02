/**
 * 用户认证信息管理工具
 */

const TOKEN_KEY = 'token'
const USER_INFO_KEY = 'user_info'

// 用户信息接口
export interface UserInfo {
  id: string | number
  email: string
  username?: string
  avatar?: string
  [key: string]: any
}

// 认证数据接口
export interface AuthData {
  token: string
  userInfo: UserInfo
}

/**
 * 保存认证信息
 */
export const setAuthData = (data: AuthData): void => {
  try {
    localStorage.setItem(TOKEN_KEY, data.token)
    localStorage.setItem(USER_INFO_KEY, JSON.stringify(data.userInfo))
  } catch (error) {
    console.error('保存认证信息失败:', error)
  }
}

/**
 * 获取 Token
 */
export const getToken = (): string | null => {
  return localStorage.getItem(TOKEN_KEY)
}

/**
 * 获取用户信息
 */
export const getUserInfo = (): UserInfo | null => {
  try {
    const userInfoStr = localStorage.getItem(USER_INFO_KEY)
    return userInfoStr ? JSON.parse(userInfoStr) : null
  } catch (error) {
    console.error('获取用户信息失败:', error)
    return null
  }
}

/**
 * 更新用户信息
 */
export const updateUserInfo = (userInfo: Partial<UserInfo>): void => {
  try {
    const currentInfo = getUserInfo()
    if (currentInfo) {
      const newInfo = { ...currentInfo, ...userInfo }
      localStorage.setItem(USER_INFO_KEY, JSON.stringify(newInfo))
    }
  } catch (error) {
    console.error('更新用户信息失败:', error)
  }
}

/**
 * 清除认证信息
 */
export const clearAuthData = (): void => {
  localStorage.removeItem(TOKEN_KEY)
  localStorage.removeItem(USER_INFO_KEY)
}

/**
 * 检查是否已登录
 */
export const isAuthenticated = (): boolean => {
  return !!getToken()
}
