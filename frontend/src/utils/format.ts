import dayjs from 'dayjs'

/**
 * 格式化日期时间
 */
export function formatDateTime(date: string | Date): string {
  return dayjs(date).format('YYYY-MM-DD HH:mm:ss')
}

/**
 * 格式化日期
 */
export function formatDate(date: string | Date): string {
  return dayjs(date).format('YYYY-MM-DD')
}

/**
 * 格式化金额（财务格式，带千位分隔符）
 * 例如：51443.02 -> 51,443.02
 */
export function formatAmount(amount: number): string {
  // 使用 toLocaleString 实现千位分隔符，并保留两位小数
  return amount.toLocaleString('zh-CN', {
    minimumFractionDigits: 2,
    maximumFractionDigits: 2
  })
}

/**
 * 判断是否为收入
 */
export function isIncome(amount: number): boolean {
  return amount > 0
}

/**
 * 判断是否为支出
 */
export function isExpense(amount: number): boolean {
  return amount < 0
}

/**
 * 获取收支类型文本
 */
export function getAmountTypeText(amount: number): string {
  if (amount > 0) return '收入'
  if (amount < 0) return '支出'
  return '其他'
}

/**
 * 格式化文件大小
 */
export function formatFileSize(bytes: number): string {
  if (bytes === 0) return '0 B'
  const k = 1024
  const sizes = ['B', 'KB', 'MB', 'GB']
  const i = Math.floor(Math.log(bytes) / Math.log(k))
  return (bytes / Math.pow(k, i)).toFixed(2) + ' ' + sizes[i]
}

/**
 * 获取文件类型文本
 */
export function getFileTypeText(type: string): string {
  const typeMap: Record<string, string> = {
    image: '图片',
    pdf: 'PDF'
  }
  return typeMap[type] || '未知'
}

/**
 * 下载文件
 */
export function downloadFile(url: string, filename: string) {
  const link = document.createElement('a')
  link.href = url
  link.download = filename
  document.body.appendChild(link)
  link.click()
  document.body.removeChild(link)
}
