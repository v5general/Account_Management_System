<template>
  <div class="audit-container">
    <el-card>
      <template #header>
        <div class="card-header">
          <span>收支审核</span>
        </div>
      </template>

      <!-- 筛选条件 -->
      <el-form :inline="true" :model="filters" class="filter-form">
        <el-form-item label="收支类型">
          <el-select v-model="filters.type" placeholder="全部" clearable @change="loadData">
            <el-option label="全部" value="all" />
            <el-option label="收入" value="income" />
            <el-option label="支出" value="expense" />
          </el-select>
        </el-form-item>
        <el-form-item label="审核状态">
          <el-select v-model="filters.status" placeholder="待审核" clearable @change="loadData">
            <el-option label="待审核" :value="0" />
            <el-option label="已审核" :value="1" />
            <el-option label="已驳回" :value="2" />
          </el-select>
        </el-form-item>
        <el-form-item label="日期范围">
          <el-date-picker
            v-model="dateRange"
            type="daterange"
            range-separator="至"
            start-placeholder="开始日期"
            end-placeholder="结束日期"
            value-format="YYYY-MM-DD"
            @change="handleDateChange"
          />
        </el-form-item>
      </el-form>

      <!-- 数据表格 -->
      <el-table :data="tableData" v-loading="loading" border>
        <el-table-column prop="transaction_time" label="交易时间" width="180" />
        <el-table-column label="收支类型" width="100">
          <template #default="{ row }">
            <el-tag :type="row.amount > 0 ? 'success' : 'danger'">
              {{ row.amount > 0 ? '收入' : '支出' }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column label="金额" width="140">
          <template #default="{ row }">
            <span :class="row.amount > 0 ? 'income-amount' : 'expense-amount'">
              {{ row.amount > 0 ? '+' : '' }}{{ formatAmount(row.amount) }}
            </span>
          </template>
        </el-table-column>
        <el-table-column prop="category.name" label="费用分类" width="120" />
        <el-table-column label="项目" width="150">
          <template #default="{ row }">
            <span v-if="row.project?.name || row.project_name">
              {{ row.amount >= 0 ? '来源项目：' : '关联项目：' }}
              {{ row.project?.name || row.project_name }}
            </span>
            <span v-else>-</span>
          </template>
        </el-table-column>
        <el-table-column prop="person.real_name" label="关联人员" width="120" />
        <el-table-column prop="creator.real_name" label="创建人" width="120" />
        <el-table-column label="审核状态" width="100">
          <template #default="{ row }">
            <el-tag :type="getStatusType(row.status)">
              {{ getStatusText(row.status) }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="remark" label="备注" min-width="150" show-overflow-tooltip />
        <el-table-column label="附件" width="100">
          <template #default="{ row }">
            <el-tag v-if="row.attachments && row.attachments.length > 0">
              {{ row.attachments.length }}个
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column label="操作" width="200" fixed="right">
          <template #default="{ row }">
            <el-button link type="primary" @click="viewDetail(row)">查看</el-button>
            <el-button v-if="row.status === 0" link type="success" @click="approve(row)">通过</el-button>
            <el-button v-if="row.status === 0" link type="danger" @click="reject(row)">驳回</el-button>
          </template>
        </el-table-column>
      </el-table>

      <!-- 分页 -->
      <el-pagination
        v-model:current-page="pagination.page"
        v-model:page-size="pagination.page_size"
        :total="pagination.total"
        :page-sizes="[10, 20, 50, 100]"
        layout="total, sizes, prev, pager, next, jumper"
        @current-change="loadData"
        @size-change="loadData"
      />
    </el-card>

    <!-- 返回按钮 -->
    <div class="back-button-container">
      <el-button @click="$router.back()" circle size="large" class="back-button">
        <el-icon><Back /></el-icon>
      </el-button>
    </div>

    <!-- 详情弹窗 -->
    <el-dialog v-model="detailVisible" title="收支详情" width="600px">
      <el-descriptions :column="2" border v-if="currentRecord">
        <el-descriptions-item label="交易时间">
          {{ currentRecord.transaction_time }}
        </el-descriptions-item>
        <el-descriptions-item label="收支类型">
          <el-tag :type="currentRecord.amount > 0 ? 'success' : 'danger'">
            {{ currentRecord.amount > 0 ? '收入' : '支出' }}
          </el-tag>
        </el-descriptions-item>
        <el-descriptions-item label="金额" :span="2">
          <span :class="currentRecord.amount > 0 ? 'income-amount' : 'expense-amount'" style="font-size: 18px">
            {{ currentRecord.amount > 0 ? '+' : '' }}{{ formatAmount(currentRecord.amount) }}
          </span>
        </el-descriptions-item>
        <el-descriptions-item label="费用分类">
          {{ currentRecord.category?.name || '-' }}
        </el-descriptions-item>
        <el-descriptions-item>
          <template #label>
            {{ currentRecord.amount >= 0 ? '来源项目' : '关联项目' }}
          </template>
          {{ currentRecord.project?.name || '-' }}
        </el-descriptions-item>
        <el-descriptions-item label="关联人员">
          {{ currentRecord.person?.real_name || '-' }}
        </el-descriptions-item>
        <el-descriptions-item label="创建人">
          {{ currentRecord.creator?.real_name || '-' }}
        </el-descriptions-item>
        <el-descriptions-item label="审核状态">
          <el-tag :type="getStatusType(currentRecord.status)">
            {{ getStatusText(currentRecord.status) }}
          </el-tag>
        </el-descriptions-item>
        <el-descriptions-item label="创建时间">
          {{ currentRecord.create_time }}
        </el-descriptions-item>
        <el-descriptions-item label="备注" :span="2">
          {{ currentRecord.remark || '-' }}
        </el-descriptions-item>
        <el-descriptions-item label="凭证附件" :span="2" v-if="currentRecord.attachments && currentRecord.attachments.length > 0">
          <div class="attachment-list">
            <div v-for="att in currentRecord.attachments" :key="att.attachment_id" class="attachment-item">
              <div class="attachment-info">
                <el-icon><Document /></el-icon>
                <span>{{ att.file_name }}</span>
                <el-tag size="small">{{ formatFileSize(att.file_size) }}</el-tag>
              </div>
              <div class="attachment-actions">
                <el-button link type="primary" @click="previewAttachment(att)" v-if="isImage(att.file_name || '')">
                  <el-icon><View /></el-icon>
                  查看
                </el-button>
                <el-button link type="primary" @click="downloadAttachment(att)">
                  <el-icon><Download /></el-icon>
                  下载
                </el-button>
              </div>
            </div>
          </div>
          <!-- 图片预览区域 -->
          <div v-if="previewImage" class="image-preview-container">
            <el-image
              :src="previewImage"
              :preview-src-list="previewImages"
              :initial-index="previewIndex"
              fit="contain"
              style="width: 100%; max-height: 400px"
            >
              <template #error>
                <div class="image-error">
                  <el-icon><Picture /></el-icon>
                  <span>图片加载失败</span>
                </div>
              </template>
            </el-image>
          </div>
        </el-descriptions-item>
      </el-descriptions>
      <template #footer>
        <el-button @click="detailVisible = false">关闭</el-button>
        <el-button v-if="currentRecord && currentRecord.status === 0" type="success" @click="approve(currentRecord)">通过</el-button>
        <el-button v-if="currentRecord && currentRecord.status === 0" type="danger" @click="reject(currentRecord)">驳回</el-button>
      </template>
    </el-dialog>

    <!-- 审核通过弹窗 -->
    <el-dialog v-model="approveVisible" title="审核通过" width="500px">
      <el-form :model="approveForm" label-width="80px">
        <el-form-item label="审核备注">
          <el-input v-model="approveForm.remark" type="textarea" :rows="3" placeholder="请输入审核备注（选填）" />
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="approveVisible = false">取消</el-button>
        <el-button type="success" @click="confirmApprove" :loading="submitLoading">确认通过</el-button>
      </template>
    </el-dialog>

    <!-- 驳回弹窗 -->
    <el-dialog v-model="rejectVisible" title="驳回审核" width="500px">
      <el-form :model="rejectForm" :rules="rejectRules" ref="rejectFormRef" label-width="80px">
        <el-form-item label="驳回原因" prop="reason">
          <el-input v-model="rejectForm.reason" type="textarea" :rows="3" placeholder="请输入驳回原因" />
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="rejectVisible = false">取消</el-button>
        <el-button type="danger" @click="confirmReject" :loading="submitLoading">确认驳回</el-button>
      </template>
    </el-dialog>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, onMounted } from 'vue'
import { getTransactionList, approveTransaction, rejectTransaction } from '@/api/transaction'
import { ElMessage, type FormInstance, type FormRules } from 'element-plus'
import dayjs from 'dayjs'
import request from '@/utils/request'

const loading = ref(false)
const submitLoading = ref(false)
const tableData = ref<any[]>([])
const dateRange = ref<[string, string]>([])

const filters = reactive({
  type: 'all',
  status: 0 // 默认显示待审核
})

const pagination = reactive({
  page: 1,
  page_size: 20,
  total: 0
})

const detailVisible = ref(false)
const approveVisible = ref(false)
const rejectVisible = ref(false)
const currentRecord = ref<any>(null)

const previewImage = ref('')
const previewImages = ref<string[]>([])
const previewIndex = ref(0)

const approveForm = reactive({
  remark: ''
})

const rejectForm = reactive({
  reason: ''
})

const rejectFormRef = ref<FormInstance>()

const rejectRules: FormRules = {
  reason: [{ required: true, message: '请输入驳回原因', trigger: 'blur' }]
}

// 判断是否为图片文件
function isImage(filename: string): boolean {
  const ext = filename.toLowerCase().split('.').pop()
  return ['jpg', 'jpeg', 'png', 'gif', 'bmp', 'webp'].includes(ext || '')
}

// 格式化文件大小
function formatFileSize(bytes: number | undefined): string {
  if (!bytes || bytes === 0) return '未知大小'
  const k = 1024
  const sizes = ['B', 'KB', 'MB', 'GB']
  const i = Math.floor(Math.log(bytes) / Math.log(k))
  return (bytes / Math.pow(k, i)).toFixed(2) + ' ' + sizes[i]
}

// 格式化金额（千分位分隔符）
function formatAmount(amount: number): string {
  return amount.toLocaleString('en-US', { minimumFractionDigits: 2, maximumFractionDigits: 2 })
}

// 预览附件
async function previewAttachment(attachment: any) {
  const token = localStorage.getItem('token')
  if (!token) {
    ElMessage.error('请先登录')
    return
  }

  if (!attachment.attachment_id) {
    ElMessage.error('附件ID不存在')
    return
  }

  try {
    const response = await fetch(`/api/v1/attachments/${attachment.attachment_id}/download`, {
      headers: {
        'Authorization': `Bearer ${token}`
      }
    })

    if (!response.ok) {
      throw new Error('加载失败')
    }

    const blob = await response.blob()
    const url = window.URL.createObjectURL(blob)
    previewImage.value = url

    // 加载所有图片到预览列表
    const attachments = currentRecord.value?.attachments || []
    previewImages.value = []
    let index = 0

    for (let i = 0; i < attachments.length; i++) {
      const att = attachments[i]
      if (att.attachment_id && isImage(att.file_name || '')) {
        if (att.attachment_id === attachment.attachment_id) {
          index = previewImages.value.length
        }
        try {
          const imgResponse = await fetch(`/api/v1/attachments/${att.attachment_id}/download`, {
            headers: { 'Authorization': `Bearer ${token}` }
          })
          if (imgResponse.ok) {
            const imgBlob = await imgResponse.blob()
            previewImages.value.push(window.URL.createObjectURL(imgBlob))
          }
        } catch (e) {
          console.error('Failed to load image:', att.file_name, e)
        }
      }
    }

    previewIndex.value = index
  } catch (error) {
    ElMessage.error('加载图片失败')
  }
}

function getStatusType(status: number) {
  switch (status) {
    case 0: return 'warning'
    case 1: return 'success'
    case 2: return 'danger'
    default: return 'info'
  }
}

function getStatusText(status: number) {
  switch (status) {
    case 0: return '待审核'
    case 1: return '已审核'
    case 2: return '已驳回'
    default: return '未知'
  }
}

function handleDateChange() {
  if (dateRange.value && dateRange.value.length === 2) {
    filters.start_time = dateRange.value[0] + ' 00:00:00'
    filters.end_time = dateRange.value[1] + ' 23:59:59'
  } else {
    filters.start_time = undefined
    filters.end_time = undefined
  }
  loadData()
}

async function loadData() {
  loading.value = true
  try {
    const params: any = {
      page: pagination.page,
      page_size: pagination.page_size,
      type: filters.type
    }
    if (filters.status !== undefined && filters.status !== null && filters.status !== '') {
      // 不过滤状态，需要在客户端过滤
    }
    if ((filters as any).start_time) {
      params.start_time = (filters as any).start_time
    }
    if ((filters as any).end_time) {
      params.end_time = (filters as any).end_time
    }

    const res = await getTransactionList(params)
    let list = res.data.list

    // 客户端过滤状态
    if (filters.status !== undefined && filters.status !== null && filters.status !== '') {
      list = list.filter((item: any) => item.status === filters.status)
    }

    tableData.value = list
    pagination.total = res.data.total
  } catch (error) {
    console.error('Failed to load audit data:', error)
  } finally {
    loading.value = false
  }
}

function viewDetail(row: any) {
  currentRecord.value = row
  previewImage.value = ''
  previewImages.value = []
  detailVisible.value = true
}

function approve(row: any) {
  currentRecord.value = row
  approveForm.remark = ''
  approveVisible.value = true
  detailVisible.value = false
}

async function confirmApprove() {
  submitLoading.value = true
  try {
    await approveTransaction(currentRecord.value.record_id, approveForm.remark)
    ElMessage.success('审核通过')
    approveVisible.value = false
    loadData()
  } catch (error) {
    // Error handled by request interceptor
  } finally {
    submitLoading.value = false
  }
}

function reject(row: any) {
  currentRecord.value = row
  rejectForm.reason = ''
  rejectVisible.value = true
  detailVisible.value = false
}

async function confirmReject() {
  if (!rejectFormRef.value) return
  await rejectFormRef.value.validate(async (valid) => {
    if (!valid) return
    submitLoading.value = true
    try {
      await rejectTransaction(currentRecord.value.record_id, rejectForm.reason)
      ElMessage.success('已驳回')
      rejectVisible.value = false
      loadData()
    } catch (error) {
      // Error handled by request interceptor
    } finally {
      submitLoading.value = false
    }
  })
}

// 下载附件
async function downloadAttachment(attachment: any) {
  const token = localStorage.getItem('token')
  if (!token) {
    ElMessage.error('请先登录')
    return
  }

  try {
    // 使用原生 fetch 下载文件
    const response = await fetch(`/api/v1/attachments/${attachment.attachment_id}/download`, {
      headers: {
        'Authorization': `Bearer ${token}`
      }
    })

    if (!response.ok) {
      throw new Error('下载失败')
    }

    const blob = await response.blob()
    const url = window.URL.createObjectURL(blob)
    const a = document.createElement('a')
    a.href = url
    a.download = attachment.file_name
    document.body.appendChild(a)
    a.click()
    window.URL.revokeObjectURL(url)
    document.body.removeChild(a)

    ElMessage.success('下载成功')
  } catch (error) {
    ElMessage.error('下载失败')
  }
}

onMounted(() => {
  loadData()
})
</script>

<style scoped>
.audit-container {
  padding: 20px;
}

.card-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.filter-form {
  margin-bottom: 20px;
}

.income-amount {
  color: #67c23a;
  font-weight: bold;
}

.expense-amount {
  color: #f56c6c;
  font-weight: bold;
}

.attachment-list {
  display: flex;
  flex-direction: column;
  gap: 8px;
}

.attachment-item {
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 8px;
  border: 1px solid #ebeef5;
  border-radius: 4px;
}

.attachment-info {
  display: flex;
  align-items: center;
  gap: 8px;
  flex: 1;
}

.attachment-actions {
  display: flex;
  align-items: center;
  gap: 8px;
}

.image-preview-container {
  margin-top: 16px;
  text-align: center;
}

.image-error {
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  height: 200px;
  color: #909399;
}

.image-error .el-icon {
  font-size: 48px;
  margin-bottom: 10px;
}

:deep(.el-pagination) {
  margin-top: 20px;
  justify-content: flex-end;
}

/* 返回按钮 */
.back-button-container {
  position: fixed;
  bottom: 30px;
  right: 30px;
  z-index: 100;
}

.back-button {
  width: 50px;
  height: 50px;
  font-size: 20px;
  box-shadow: 0 2px 12px 0 rgba(0, 0, 0, 0.1);
}
</style>
