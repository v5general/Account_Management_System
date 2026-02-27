import request from './index'

export interface StatisticsParams {
  dimension: 'project' | 'person' | 'category'
  cycle: 'month' | 'quarter' | 'year'
  start_time?: string
  end_time?: string
}

export interface StatisticsResponse {
  dimension: string
  cycle: string
  start_time: string
  end_time: string
  summary: {
    total_income: number
    total_expense: number
    net_amount: number
    record_count: number
  }
  details: {
    key: string
    income: number
    expense: number
    net_amount: number
    record_count: number
    percentage: number
  }[]
}

// 获取统计数据
export function getStatistics(params: StatisticsParams) {
  return request.get<StatisticsResponse>('/transactions/statistics', { params })
}
