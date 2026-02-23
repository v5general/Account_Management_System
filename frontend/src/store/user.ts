import { defineStore } from 'pinia'
import { ref } from 'vue'
import type { UserInfo } from '@/api/auth'
import { getCurrentUser, logout as logoutApi } from '@/api/auth'
import { ElMessage } from 'element-plus'

export const useUserStore = defineStore('user', () => {
  const token = ref<string>(localStorage.getItem('token') || '')
  const userInfo = ref<UserInfo | null>(null)

  // 设置token
  function setToken(newToken: string) {
    token.value = newToken
    localStorage.setItem('token', newToken)
  }

  // 清除token
  function clearToken() {
    token.value = ''
    localStorage.removeItem('token')
  }

  // 设置用户信息
  function setUserInfo(info: UserInfo) {
    userInfo.value = info
  }

  // 清除用户信息
  function clearUserInfo() {
    userInfo.value = null
  }

  // 获取用户信息
  async function fetchUserInfo() {
    try {
      const res = await getCurrentUser()
      setUserInfo(res.data)
    } catch (error) {
      clearToken()
      clearUserInfo()
    }
  }

  // 登出
  async function logout() {
    try {
      await logoutApi()
    } catch (error) {
      console.error('Logout error:', error)
    } finally {
      clearToken()
      clearUserInfo()
      // 延迟导入 router 避免循环依赖
      import('@/router').then(({ default: router }) => {
        router.push('/login')
      })
      ElMessage.success('已退出登录')
    }
  }

  // 检查是否有权限
  function hasRole(role: string | string[]) {
    if (!userInfo.value) return false
    if (Array.isArray(role)) {
      return role.includes(userInfo.value.role)
    }
    return userInfo.value.role === role
  }

  // 是否是管理员
  const isAdmin = () => userInfo.value?.role === 'ADMIN'

  // 是否是财务人员
  const isFinance = () => userInfo.value?.role === 'FINANCE' || userInfo.value?.role === 'ADMIN'

  // 是否是员工
  const isEmployee = () => userInfo.value?.role === 'EMPLOYEE'

  return {
    token,
    userInfo,
    setToken,
    clearToken,
    setUserInfo,
    clearUserInfo,
    fetchUserInfo,
    logout,
    hasRole,
    isAdmin,
    isFinance,
    isEmployee
  }
})
