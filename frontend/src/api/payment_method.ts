import request from './index'

export interface PaymentMethod {
  payment_method_id: string
  name: string
  description: string
  sort_order: number
  is_deleted: number
  create_time: string
  update_time: string
}

export interface CreatePaymentMethodParams {
  name: string
  description?: string
  sort_order?: number
}

// 创建支付方式
export function createPaymentMethod(data: CreatePaymentMethodParams) {
  return request.post('/payment-methods', data)
}

// 获取支付方式列表
export function getPaymentMethodList(params?: { page?: number; page_size?: number; keyword?: string }) {
  return request.get('/payment-methods', { params })
}

// 更新支付方式
export function updatePaymentMethod(id: string, data: Partial<PaymentMethod>) {
  return request.put(`/payment-methods/${id}`, data)
}

// 删除支付方式
export function deletePaymentMethod(id: string) {
  return request.delete(`/payment-methods/${id}`)
}
