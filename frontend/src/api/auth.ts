import request from './index'

export interface LoginParams {
  username: string
  password: string
}

export interface LoginResponse {
  token: string
  user: {
    user_id: string
    username: string
    role: string
    department: string
  }
}

export interface UserInfo {
  user_id: string
  username: string
  role: string
  department: string
}

// 用户登录
export function login(data: LoginParams) {
  return request.post<LoginResponse>('/auth/login', data)
}

// 用户注销
export function logout() {
  return request.post('/auth/logout')
}

// 获取当前用户信息
export function getCurrentUser() {
  return request.get<{ data: UserInfo }>('/auth/me')
}
