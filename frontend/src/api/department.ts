import request from '@/utils/request'

export interface Department {
  department_id: string
  name: string
  description: string
  sort_order: number
  is_deleted: number
  create_time: string
  update_time: string
}

export interface CreateDepartmentParams {
  name: string
  description?: string
  sort_order?: number
}

export interface UpdateDepartmentParams {
  name?: string
  description?: string
  sort_order?: number
}

// 创建部门
export function createDepartment(data: CreateDepartmentParams) {
  return request.post('/departments', data)
}

// 获取部门列表
export function getDepartmentList() {
  return request.get<Department[]>('/departments')
}

// 获取部门详情
export function getDepartment(id: string) {
  return request.get<Department>(`/departments/${id}`)
}

// 更新部门
export function updateDepartment(id: string, data: UpdateDepartmentParams) {
  return request.put(`/departments/${id}`, data)
}

// 删除部门
export function deleteDepartment(id: string) {
  return request.delete(`/departments/${id}`)
}

// 获取部门下的用户列表
export function getDepartmentUsers(id: string, params?: { page?: number; page_size?: number }) {
  return request.get(`/departments/${id}/users`, { params })
}
