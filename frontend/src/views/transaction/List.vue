<template>
  <div class="transaction-list">
    <el-card>
      <template #header>
        <div class="card-header">
          <span>{{ userStore.isEmployee() ? '我的收支' : '收支记录' }}</span>
          <div class="header-actions" v-if="!userStore.isEmployee()">
            <el-button type="primary" @click="$router.push('/transaction/income')" v-if="userStore.isFinance()">
              <el-icon><Plus /></el-icon>
              收入登记
            </el-button>
            <el-button type="danger" @click="$router.push('/transaction/expense')" v-if="userStore.isFinance()">
              <el-icon><Plus /></el-icon>
              支出登记
            </el-button>
          </div>
        </div>
      </template>

      <el-form :inline="true" :model="queryForm" class="query-form">
        <el-form-item label="时间范围">
          <el-date-picker
            v-model="dateRange"
            type="daterange"
            range-separator="至"
            start-placeholder="开始日期"
            end-placeholder="结束日期"
            @change="handleDateChange"
          />
        </el-form-item>
        <el-form-item label="费用分类">
          <el-select v-model="queryForm.category_id" placeholder="请选择" clearable>
            <el-option
              v-for="cat in categories"
              :key="cat.category_id"
              :label="cat.name"
              :value="cat.category_id"
            />
          </el-select>
        </el-form-item>
        <el-form-item label="项目">
          <el-select v-model="queryForm.project_id" placeholder="请选择项目" clearable>
            <el-option
              v-for="proj in projects"
              :key="proj.project_id"
              :label="proj.name"
              :value="proj.project_id"
            />
          </el-select>
        </el-form-item>
        <el-form-item label="类型">
          <el-select v-model="queryForm.type" placeholder="请选择" clearable>
            <el-option label="全部" value="all" />
            <el-option label="收入" value="income" />
            <el-option label="支出" value="expense" />
          </el-select>
        </el-form-item>
        <el-form-item label="审核状态">
          <el-select v-model="queryForm.status" placeholder="请选择" clearable>
            <el-option label="全部" value="" />
            <el-option label="待审核" :value="0" />
            <el-option label="已审核" :value="1" />
            <el-option label="已驳回" :value="2" />
          </el-select>
        </el-form-item>
        <el-form-item>
          <el-button type="primary" @click="handleQuery">查询</el-button>
          <el-button @click="handleReset">重置</el-button>
        </el-form-item>
      </el-form>

      <el-table :data="tableData" stripe v-loading="loading">
        <el-table-column label="交易时间" width="120">
          <template #default="{ row }">
            {{ formatDate(row.transaction_time) }}
          </template>
        </el-table-column>
        <el-table-column label="项目" width="180">
          <template #default="{ row }">
            <span v-if="row.project_name || row.project?.name">
              {{ row.amount >= 0 ? '来源项目：' : '关联项目：' }}
              {{ row.project_name || row.project?.name }}
            </span>
            <span v-else>-</span>
          </template>
        </el-table-column>
        <el-table-column label="审核状态" width="100">
          <template #default="{ row }">
            <el-tag :type="getStatusType(row.status)">
              {{ getStatusText(row.status) }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column label="金额" width="150">
          <template #default="{ row }">
            <span :class="row.amount > 0 ? 'income' : 'expense'">
              {{ row.amount > 0 ? '+' : '' }}{{ formatAmount(row.amount) }}
            </span>
          </template>
        </el-table-column>
        <el-table-column prop="remark" label="备注" min-width="150" show-overflow-tooltip />
        <el-table-column label="操作" width="150" fixed="right">
          <template #default="{ row }">
            <el-button link type="primary" @click="handleView(row)">查看</el-button>
            <el-button link type="primary" @click="handleViewAttachments(row)" v-if="row.attachments?.length">
              凭证({{ row.attachments.length }})
            </el-button>
          </template>
        </el-table-column>
      </el-table>

      <el-pagination
        v-model:current-page="queryForm.page"
        v-model:page-size="queryForm.page_size"
        :total="total"
        :page-sizes="[10, 20, 50, 100]"
        layout="total, sizes, prev, pager, next, jumper"
        @size-change="loadData"
        @current-change="loadData"
      />
    </el-card>

    <!-- 返回按钮 -->
    <div class="back-button-container">
      <el-button @click="$router.back()" circle size="large" class="back-button">
        <el-icon><Back /></el-icon>
      </el-button>
    </div>

    <!-- 凭证查看对话框 -->
    <el-dialog v-model="attachmentDialogVisible" title="凭证附件" width="800px">
      <div class="attachment-list">
        <div v-for="att in currentAttachments" :key="att.attachment_id" class="attachment-item">
          <div class="attachment-info">
            <el-icon><Document /></el-icon>
            <span>{{ att.file_name }}</span>
            <el-tag size="small">{{ formatFileSize(att.file_size) }}</el-tag>
          </div>
          <div class="attachment-actions">
            <el-button link type="primary" @click="handlePreviewAttachment(att)" v-if="isImage(att.file_name || '')">
              <el-icon><View /></el-icon>
              查看
            </el-button>
            <el-button link type="primary" @click="handleDownloadAttachment(att)">
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
          style="width: 100%; max-height: 500px"
        >
          <template #error>
            <div class="image-error">
              <el-icon><Picture /></el-icon>
              <span>图片加载失败</span>
            </div>
          </template>
        </el-image>
      </div>
    </el-dialog>

    <!-- 收支详情对话框 -->
    <el-dialog v-model="detailDialogVisible" title="收支详情" width="700px">
      <el-descriptions :column="2" border v-if="currentTransaction">
        <el-descriptions-item label="交易时间">
          {{ formatDate(currentTransaction.transaction_time) }}
        </el-descriptions-item>
        <el-descriptions-item label="金额">
          <span :class="currentTransaction.amount > 0 ? 'income' : 'expense'" style="font-size: 16px; font-weight: bold;">
            {{ currentTransaction.amount > 0 ? '+' : '' }}{{ formatAmount(currentTransaction.amount) }}
          </span>
        </el-descriptions-item>
        <el-descriptions-item label="费用分类">
          {{ currentTransaction.category?.name || '-' }}
        </el-descriptions-item>
        <el-descriptions-item>
          <template #label>
            {{ currentTransaction.amount >= 0 ? '来源项目' : '关联项目' }}
          </template>
          {{ currentTransaction.project?.name || '-' }}
        </el-descriptions-item>
        <el-descriptions-item label="关联人员">
          {{ currentTransaction.person?.real_name || '-' }}
        </el-descriptions-item>
        <el-descriptions-item label="审核状态">
          <el-tag :type="getStatusType(currentTransaction.status)">
            {{ getStatusText(currentTransaction.status) }}
          </el-tag>
        </el-descriptions-item>
        <el-descriptions-item label="录入人">
          {{ currentTransaction.creator?.real_name || '-' }}
        </el-descriptions-item>
        <el-descriptions-item label="创建时间">
          {{ formatDateTime(currentTransaction.create_time) }}
        </el-descriptions-item>
        <el-descriptions-item label="修改时间">
          {{ formatDateTime(currentTransaction.update_time) }}
        </el-descriptions-item>
        <el-descriptions-item label="备注" :span="2">
          {{ currentTransaction.remark || '-' }}
        </el-descriptions-item>
      </el-descriptions>
      <template #footer v-if="currentTransaction?.attachments?.length">
        <div class="detail-attachments">
          <div class="attachment-title">凭证附件：</div>
          <div class="attachment-list">
            <div v-for="att in currentTransaction.attachments" :key="att.attachment_id" class="attachment-item">
              <div class="attachment-info">
                <el-icon><Document /></el-icon>
                <span>{{ att.file_name }}</span>
                <el-tag size="small">{{ formatFileSize(att.file_size) }}</el-tag>
              </div>
              <div class="attachment-actions">
                <el-button link type="primary" @click="handlePreviewAttachment(att)" v-if="isImage(att.file_name || '')">
                  <el-icon><View /></el-icon>
                  查看
                </el-button>
                <el-button link type="primary" @click="handleDownloadAttachment(att)">
                  <el-icon><Download /></el-icon>
                  下载
                </el-button>
              </div>
            </div>
          </div>
        </div>
      </template>
    </el-dialog>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, onMounted } from 'vue'
import { getTransactionList, getTransactionDetail } from '@/api/transaction'
import { getCategoryList } from '@/api/category'
import { getProjectList } from '@/api/project'
import { getAttachmentUrl } from '@/api/attachment'
import { useUserStore } from '@/store/user'
import { ElMessage } from 'element-plus'
import type { Transaction, Attachment } from '@/api/transaction'

const userStore = useUserStore()

const loading = ref(false)
const tableData = ref<Transaction[]>([])
const total = ref(0)
const categories = ref([])
const projects = ref([])
const dateRange = ref<[Date, Date]>()
const attachmentDialogVisible = ref(false)
const currentAttachments = ref<Attachment[]>([])
const detailDialogVisible = ref(false)
const currentTransaction = ref<Transaction | null>(null)

const previewImage = ref('')
const previewImages = ref<string[]>([])
const previewIndex = ref(0)

const queryForm = reactive({
  page: 1,
  page_size: 20,
  start_time: '',
  end_time: '',
  category_id: '',
  project_id: '',
  type: 'all',
  status: ''
})

// 判断是否为图片文件
function isImage(filename: string): boolean {
  const ext = filename.toLowerCase().split('.').pop()
  return ['jpg', 'jpeg', 'png', 'gif', 'bmp', 'webp'].includes(ext || '')
}

// 格式化日期（只显示到日）
function formatDate(dateStr: string): string {
  if (!dateStr) return '-'
  return dateStr.split(' ')[0]
}

// 格式化日期时间（显示到秒）
function formatDateTime(dateStr: string): string {
  if (!dateStr) return '-'
  return dateStr
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

// 获取审核状态类型
function getStatusType(status: number): string {
  switch (status) {
    case 0: return 'warning'
    case 1: return 'success'
    case 2: return 'danger'
    default: return 'info'
  }
}

// 获取审核状态文本
function getStatusText(status: number): string {
  switch (status) {
    case 0: return '待审核'
    case 1: return '已审核'
    case 2: return '已驳回'
    default: return '未知'
  }
}

// 获取带token的图片URL
function getImageUrl(att: Attachment): string {
  const token = localStorage.getItem('token')
  return `/api/v1/attachments/${att.attachment_id}/download?token=${token}`
}

// 预览附件
async function handlePreviewAttachment(att: Attachment) {
  const token = localStorage.getItem('token')
  if (!token) {
    ElMessage.error('请先登录')
    return
  }

  if (!att.attachment_id) {
    ElMessage.error('附件ID不存在')
    return
  }

  try {
    const response = await fetch(`/api/v1/attachments/${att.attachment_id}/download`, {
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
    const attachments = currentAttachments.value.length > 0
      ? currentAttachments.value
      : currentTransaction.value?.attachments || []
    previewImages.value = []
    let index = 0

    for (let i = 0; i < attachments.length; i++) {
      const a = attachments[i]
      if (a.attachment_id && isImage(a.file_name || '')) {
        if (a.attachment_id === att.attachment_id) {
          index = previewImages.value.length
        }
        try {
          // 为每个图片创建一个临时的 blob URL
          const imgResponse = await fetch(`/api/v1/attachments/${a.attachment_id}/download`, {
            headers: { 'Authorization': `Bearer ${token}` }
          })
          if (imgResponse.ok) {
            const imgBlob = await imgResponse.blob()
            previewImages.value.push(window.URL.createObjectURL(imgBlob))
          }
        } catch (e) {
          // 忽略单个图片加载失败，继续加载其他图片
          console.error('Failed to load image:', a.file_name, e)
        }
      }
    }

    previewIndex.value = index
  } catch (error) {
    ElMessage.error('加载图片失败')
  }
}

function handleDateChange(val: [Date, Date] | null) {
  if (val) {
    queryForm.start_time = val[0].toISOString().split('T')[0]
    queryForm.end_time = val[1].toISOString().split('T')[0]
  } else {
    queryForm.start_time = ''
    queryForm.end_time = ''
  }
}

async function loadCategories() {
  try {
    const res = await getCategoryList({ page: 1, page_size: 100 })
    categories.value = res.data.list
  } catch (error) {
    console.error('Failed to load categories:', error)
  }
}

async function loadProjects() {
  try {
    const res = await getProjectList({ page: 1, page_size: 100 })
    projects.value = res.data.list?.filter((p: any) => p.status === 1) || []
  } catch (error) {
    console.error('Failed to load projects:', error)
  }
}

async function loadData() {
  loading.value = true
  try {
    // 员工只能查看自己的收支记录
    const params = { ...queryForm }
    if (userStore.isEmployee()) {
      params.person_id = userStore.userInfo?.user_id
    }
    const res = await getTransactionList(params)
    tableData.value = res.data.list
    total.value = res.data.total
  } catch (error) {
    console.error('Failed to load transactions:', error)
  } finally {
    loading.value = false
  }
}

function handleQuery() {
  queryForm.page = 1
  loadData()
}

function handleReset() {
  dateRange.value = undefined
  Object.assign(queryForm, {
    page: 1,
    page_size: 20,
    start_time: '',
    end_time: '',
    category_id: '',
    project_id: '',
    type: 'all',
    status: ''
  })
  loadData()
}

async function handleView(row: Transaction) {
  try {
    const res = await getTransactionDetail(row.record_id)
    currentTransaction.value = res.data
    detailDialogVisible.value = true
  } catch (error) {
    console.error('获取详情失败:', error)
    ElMessage.error('获取详情失败')
  }
}

function handleViewAttachments(row: Transaction) {
  currentAttachments.value = row.attachments || []
  previewImage.value = ''
  previewImages.value = []
  attachmentDialogVisible.value = true
}

async function handleDownloadAttachment(att: Attachment) {
  const token = localStorage.getItem('token')
  if (!token) {
    ElMessage.error('请先登录')
    return
  }

  try {
    // 使用原生 fetch 下载文件
    const response = await fetch(`/api/v1/attachments/${att.attachment_id}/download`, {
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
    a.download = att.file_name
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
  loadCategories()
  loadProjects()
  loadData()
})
</script>

<style scoped>
.transaction-list {
  padding: 20px;
}

.card-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.query-form {
  margin-bottom: 20px;
}

.el-pagination {
  margin-top: 20px;
  justify-content: flex-end;
}

.income {
  color: #67c23a;
}

.expense {
  color: #f56c6c;
}

.attachment-list {
  display: flex;
  flex-direction: column;
  gap: 10px;
}

.attachment-item {
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 10px;
  border: 1px solid #ebeef5;
  border-radius: 4px;
}

.attachment-info {
  display: flex;
  align-items: center;
  gap: 10px;
  flex: 1;
}

.attachment-actions {
  display: flex;
  align-items: center;
  gap: 10px;
}

.detail-attachments {
  padding: 10px 0;
}

.attachment-title {
  font-weight: bold;
  margin-bottom: 10px;
}

.image-preview-container {
  margin-top: 20px;
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
