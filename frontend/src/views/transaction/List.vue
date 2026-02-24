<template>
  <div class="transaction-list">
    <el-card>
      <template #header>
        <div class="card-header">
          <span>收支记录</span>
          <div class="header-actions">
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
          <el-input v-model="queryForm.project_name" placeholder="请输入项目名称" clearable />
        </el-form-item>
        <el-form-item label="类型">
          <el-select v-model="queryForm.type" placeholder="请选择" clearable>
            <el-option label="全部" value="all" />
            <el-option label="收入" value="income" />
            <el-option label="支出" value="expense" />
          </el-select>
        </el-form-item>
        <el-form-item>
          <el-button type="primary" @click="handleQuery">查询</el-button>
          <el-button @click="handleReset">重置</el-button>
        </el-form-item>
      </el-form>

      <el-table :data="tableData" stripe v-loading="loading">
        <el-table-column prop="transaction_time" label="交易时间" width="180" />
        <el-table-column prop="project_name" label="项目" width="150" />
        <el-table-column prop="category.name" label="费用分类" width="120" />
        <el-table-column prop="person.username" label="关联人员" width="120" />
        <el-table-column prop="amount" label="金额" width="120">
          <template #default="{ row }">
            <span :class="row.amount > 0 ? 'income' : 'expense'">
              {{ row.amount > 0 ? '+' : '' }}{{ row.amount }}
            </span>
          </template>
        </el-table-column>
        <el-table-column prop="remark" label="备注" show-overflow-tooltip />
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

    <!-- 凭证查看对话框 -->
    <el-dialog v-model="attachmentDialogVisible" title="凭证附件" width="600px">
      <div class="attachment-list">
        <div v-for="att in currentAttachments" :key="att.attachment_id" class="attachment-item">
          <el-icon><Document /></el-icon>
          <span>{{ att.file_name }}</span>
          <el-button link type="primary" @click="handleDownloadAttachment(att)">下载</el-button>
        </div>
      </div>
    </el-dialog>

    <!-- 收支详情对话框 -->
    <el-dialog v-model="detailDialogVisible" title="收支详情" width="700px">
      <el-descriptions :column="2" border v-if="currentTransaction">
        <el-descriptions-item label="交易时间">
          {{ currentTransaction.transaction_time }}
        </el-descriptions-item>
        <el-descriptions-item label="金额">
          <span :class="currentTransaction.amount > 0 ? 'income' : 'expense'">
            {{ currentTransaction.amount > 0 ? '+' : '' }}{{ currentTransaction.amount }}
          </span>
        </el-descriptions-item>
        <el-descriptions-item label="费用分类">
          {{ currentTransaction.category?.name || '-' }}
        </el-descriptions-item>
        <el-descriptions-item label="项目">
          {{ currentTransaction.project_name || '-' }}
        </el-descriptions-item>
        <el-descriptions-item label="关联人员">
          {{ currentTransaction.person?.username || '-' }}
        </el-descriptions-item>
        <el-descriptions-item label="状态">
          <el-tag v-if="currentTransaction.status === 0" type="warning">待审核</el-tag>
          <el-tag v-else-if="currentTransaction.status === 1" type="success">已审核</el-tag>
          <el-tag v-else type="danger">已驳回</el-tag>
        </el-descriptions-item>
        <el-descriptions-item label="录入人">
          {{ currentTransaction.creator?.username || '-' }}
        </el-descriptions-item>
        <el-descriptions-item label="创建时间">
          {{ currentTransaction.create_time }}
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
              <el-icon><Document /></el-icon>
              <span>{{ att.file_name }}</span>
              <el-button link type="primary" @click="handleDownloadAttachment(att)">下载</el-button>
            </div>
          </div>
        </div>
      </template>
    </el-dialog>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, onMounted } from 'vue'
import { getTransactionList } from '@/api/transaction'
import { getCategoryList } from '@/api/category'
import { getAttachmentUrl } from '@/api/attachment'
import { useUserStore } from '@/store/user'
import type { Transaction, Attachment } from '@/api/transaction'

const userStore = useUserStore()

const loading = ref(false)
const tableData = ref<Transaction[]>([])
const total = ref(0)
const categories = ref([])
const dateRange = ref<[Date, Date]>()
const attachmentDialogVisible = ref(false)
const currentAttachments = ref<Attachment[]>([])
const detailDialogVisible = ref(false)
const currentTransaction = ref<Transaction | null>(null)

const queryForm = reactive({
  page: 1,
  page_size: 20,
  start_time: '',
  end_time: '',
  category_id: '',
  project_name: '',
  type: 'all'
})

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

async function loadData() {
  loading.value = true
  try {
    const res = await getTransactionList(queryForm)
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
    project_name: '',
    type: 'all'
  })
  loadData()
}

function handleView(row: Transaction) {
  currentTransaction.value = row
  detailDialogVisible.value = true
}

function handleViewAttachments(row: Transaction) {
  currentAttachments.value = row.attachments || []
  attachmentDialogVisible.value = true
}

function handleDownloadAttachment(att: Attachment) {
  window.open(getAttachmentUrl(att.attachment_id), '_blank')
}

onMounted(() => {
  loadCategories()
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
  gap: 10px;
  padding: 10px;
  border: 1px solid #ebeef5;
  border-radius: 4px;
}

.detail-attachments {
  padding: 10px 0;
}

.attachment-title {
  font-weight: bold;
  margin-bottom: 10px;
}
</style>
