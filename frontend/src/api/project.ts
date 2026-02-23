import request from '@/utils/request'

export interface Project {
  project_id: string
  name: string
  description: string
  department_id: string
  status: number
  is_deleted: number
  create_time: string
  update_time: string
}

export interface CreateProjectParams {
  name: string
  description?: string
  department_id?: string
}

export interface UpdateProjectParams {
  name?: string
  description?: string
  department_id?: string
  status?: number
}

export interface ProjectListParams {
  page?: number
  page_size?: number
  department_id?: string
}

export interface ProjectListResponse {
  list: Project[]
  total: number
}

// 创建项目
export function createProject(data: CreateProjectParams) {
  return request.post('/projects', data)
}

// 获取项目列表
export function getProjectList(params?: ProjectListParams) {
  return request.get<ProjectListResponse>('/projects', { params })
}

// 获取项目详情
export function getProject(id: string) {
  return request.get<Project>(`/projects/${id}`)
}

// 更新项目
export function updateProject(id: string, data: UpdateProjectParams) {
  return request.put(`/projects/${id}`, data)
}

// 删除项目
export function deleteProject(id: string) {
  return request.delete(`/projects/${id}`)
}
