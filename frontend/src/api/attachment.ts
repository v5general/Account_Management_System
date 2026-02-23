import request from './index'

export interface UploadResponse {
  attachment_id: string
  file_name: string
  file_size: number
  file_type: string
}

// 上传附件
export function uploadAttachment(file: File) {
  const formData = new FormData()
  formData.append('file', file)
  return request.post<UploadResponse>('/attachments', formData, {
    headers: {
      'Content-Type': 'multipart/form-data'
    }
  })
}

// 获取附件列表
export function getAttachmentList(recordId: string) {
  return request.get('/attachments', { params: { record_id: recordId } })
}

// 获取附件下载URL
export function getAttachmentUrl(id: string) {
  return `/api/v1/attachments/${id}/download`
}

// 删除附件
export function deleteAttachment(id: string) {
  return request.delete(`/attachments/${id}`)
}
