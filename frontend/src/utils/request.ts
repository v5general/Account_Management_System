import axios from 'axios'
import { ElMessage } from 'element-plus'

// 创建axios实例
const request = axios.create({
  baseURL: '/api/v1',
  timeout: 30000
})

// 请求拦截器
request.interceptors.request.use(
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

// 响应拦截器
request.interceptors.response.use(
  (response) => {
    const res = response.data
    if (res.code !== 0) {
      // 以下错误码不显示提示：权限不足、数据不存在、令牌过期、账号已禁用等
      // 注意：只有2003（用户名重复）需要显示提示
      const silentCodes = [1002, 1003, 2002, 2004]
      if (!silentCodes.includes(res.code)) {
        ElMessage.error(res.message || '请求失败')
      }

      // 1002: 令牌无效或过期
      if (res.code === 1002) {
        localStorage.removeItem('token')
        localStorage.removeItem('user')
        // 动态导入router避免循环依赖
        import('@/router').then(({ default: router }) => {
          router.push('/login')
        })
      }

      return Promise.reject(new Error(res.message || '请求失败'))
    }
    return res
  },
  (error) => {
    // 所有HTTP错误都不显示提示
    return Promise.reject(error)
  }
)

export default request
