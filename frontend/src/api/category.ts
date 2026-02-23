import request from './index'

// 分类类型：收入/支出
export type CategoryType = 'INCOME' | 'EXPENSE'

export interface Category {
  category_id: string
  name: string
  type: CategoryType
  description: string
  sort_order: number
  is_deleted: number
  create_time: string
  update_time: string
}

export interface CreateCategoryParams {
  name: string
  type: CategoryType
  description?: string
  sort_order?: number
}

// 创建费用分类
export function createCategory(data: CreateCategoryParams) {
  return request.post('/categories', data)
}

// 获取分类列表
export function getCategoryList(params?: { page?: number; page_size?: number; keyword?: string; type?: CategoryType }) {
  return request.get('/categories', { params })
}

// 更新分类
export function updateCategory(id: string, data: Partial<Category>) {
  return request.put(`/categories/${id}`, data)
}

// 删除分类
export function deleteCategory(id: string) {
  return request.delete(`/categories/${id}`)
}
