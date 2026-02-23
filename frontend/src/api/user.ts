import request from '@/utils/request'

export interface UserInfo {
  user_id: string
  username: string
  role: string
  department: string
  status?: number
  create_time?: string
}

export interface CreateUserParams {
  username: string
  password: string
  role: 'ADMIN' | 'EMPLOYEE' | 'FINANCE'
  department?: string
}

export interface UpdateUserParams {
  role?: 'ADMIN' | 'EMPLOYEE' | 'FINANCE'
  department?: string
  status?: number
}

export interface UserListParams {
  page?: number
  page_size?: number
  role?: string
  keyword?: string
}

export interface UserListResponse {
  list: UserInfo[]
  total: number
  page: number
  page_size: number
}

// 创建用户
export function createUser(data: CreateUserParams) {
  return request.post('/users', data)
}

// 获取用户列表
export function getUserList(params: UserListParams) {
  return request.get<UserListResponse>('/users', { params })
}

// 更新用户
export function updateUser(id: string, data: UpdateUserParams) {
  return request.put(`/users/${id}`, data)
}

// 删除用户
export function deleteUser(id: string) {
  return request.delete(`/users/${id}`)
}

// 重置密码
export function resetPassword(id: string, password: string) {
  return request.post(`/users/${id}/reset-password`, { password })
}
