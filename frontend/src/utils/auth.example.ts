/**
 * 用户认证系统使用示例
 * 
 * 本文件展示如何在项目中使用认证系统
 */

import { useUserStore } from '@/stores/user'
import { getToken, isAuthenticated } from '@/utils/auth'

// ============================================
// 1. 在组件中使用用户信息
// ============================================

// 在 Vue 组件中
export function componentExample() {
  const userStore = useUserStore()
  
  // 获取用户信息
  console.log('用户ID:', userStore.userId)
  console.log('用户邮箱:', userStore.userEmail)
  console.log('用户名:', userStore.userName)
  console.log('是否登录:', userStore.isLoggedIn)
  console.log('完整用户信息:', userStore.userInfo)
  
  // 更新用户信息
  userStore.updateUserInfo({
    username: '新用户名',
    avatar: 'https://example.com/avatar.jpg'
  })
  
  // 退出登录
  const logout = () => {
    userStore.logout()
    // 跳转到登录页
    // router.push('/user/login')
  }
}

// ============================================
// 2. 在路由守卫中使用
// ============================================

export function routerGuardExample() {
  // 在 router/guard/index.ts 中
  /*
  router.beforeEach((to, from, next) => {
    // 检查是否需要登录
    if (to.meta.requiresAuth) {
      if (isAuthenticated()) {
        next()
      } else {
        next('/user/login')
      }
    } else {
      next()
    }
  })
  */
}

// ============================================
// 3. 在 HTTP 请求中使用 Token
// ============================================

export function httpRequestExample() {
  // 在 utils/req.ts 中
  /*
  import { getToken } from '@/utils/auth'
  
  // 请求拦截器
  axios.interceptors.request.use(config => {
    const token = getToken()
    if (token) {
      config.headers.Authorization = `Bearer ${token}`
    }
    return config
  })
  
  // 响应拦截器
  axios.interceptors.response.use(
    response => response,
    error => {
      // Token 过期或无效
      if (error.response?.status === 401) {
        const userStore = useUserStore()
        userStore.logout()
        router.push('/user/login')
      }
      return Promise.reject(error)
    }
  )
  */
}

// ============================================
// 4. 登录流程示例
// ============================================

export async function loginExample() {
  /*
  // 在登录页面
  import { login as loginApi } from '@/api'
  import { useUserStore } from '@/stores/user'
  
  const handleLogin = async () => {
    try {
      const response = await loginApi({
        email: 'user@example.com',
        password: '123456'
      })
      
      const { token, user } = response.data
      
      const userStore = useUserStore()
      userStore.login(token, {
        id: user.id,
        email: user.email,
        username: user.username,
        ...user
      })
      
      // 跳转到主页
      router.push('/home')
    } catch (error) {
      console.error('登录失败:', error)
    }
  }
  */
}

// ============================================
// 5. 注册流程示例
// ============================================

export async function registerExample() {
  /*
  // 在注册页面
  import { registerByEmailCode } from '@/api'
  import { useUserStore } from '@/stores/user'
  
  const handleRegister = async () => {
    try {
      const response = await registerByEmailCode({
        email: 'user@example.com',
        verifyCode: '123456',
        password: '123456'
      })
      
      const { token, user } = response.data
      
      const userStore = useUserStore()
      userStore.login(token, {
        id: user.id,
        email: user.email,
        ...user
      })
      
      // 跳转到标签选择页面或主页
      router.push('/user/entry')
    } catch (error) {
      console.error('注册失败:', error)
    }
  }
  */
}

// ============================================
// 6. localStorage 数据结构
// ============================================

/*
// 存储的数据结构示例:

localStorage.getItem('multi_table_token')
// "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."

localStorage.getItem('multi_table_user_info')
// {
//   "id": "12345",
//   "email": "user@example.com",
//   "username": "用户名",
//   "avatar": "https://example.com/avatar.jpg"
// }
*/
