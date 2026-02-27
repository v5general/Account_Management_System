import request from './index'

export interface Transaction {
  record_id: string
  amount: number
  category_id?: string
  project_id?: string
  project_name?: string
  person_id?: string
  transaction_time: string
  remark: string
  status: number
  creator_id: string
  create_time: string
  update_time: string
  category?: {
    category_id: string
    name: string
  }
  project?: {
    project_id: string
    name: string
    description: string
  }
  person?: {
    user_id: string
    username: string
    real_name: string
  }
  creator?: {
    user_id: string
    username: string
    real_name: string
  }
  attachments?: Attachment[]
}

export interface Attachment {
  attachment_id: string
  record_id: string
  file_name: string
  file_path: string
  file_size: number
  file_type: string
  upload_time: string
}

export interface CreateTransactionParams {
  amount: number
  category_id?: string
  project_id?: string
  person_id?: string
  transaction_time: string
  remark?: string
  attachment_ids: string[]
}

export interface TransactionListParams {
  page?: number
  page_size?: number
  start_time?: string
  end_time?: string
  category_id?: string
  project_id?: string
  person_id?: string
  type?: 'income' | 'expense' | 'all'
}

// 创建收支记录
export function createTransaction(data: CreateTransactionParams) {
  return request.post('/transactions', data)
}

// 获取收支记录列表
export function getTransactionList(params: TransactionListParams) {
  return request.get('/transactions', { params })
}

// 获取收支记录详情
export function getTransactionDetail(id: string) {
  return request.get(`/transactions/${id}`)
}

// 更新收支记录
export function updateTransaction(id: string, data: Partial<Transaction>) {
  return request.put(`/transactions/${id}`, data)
}

// 删除收支记录
export function deleteTransaction(id: string) {
  return request.delete(`/transactions/${id}`)
}

// 审核通过
export function approveTransaction(id: string, remark?: string) {
  return request.put(`/transactions/${id}/approve`, { remark })
}

// 驳回
export function rejectTransaction(id: string, reason: string) {
  return request.put(`/transactions/${id}/reject`, { reason })
}
