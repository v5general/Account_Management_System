<template>
  <div class="income-form">
    <el-card>
      <template #header>
        <span>收入登记</span>
      </template>

      <el-form :model="form" :rules="rules" ref="formRef" label-width="120px">
        <el-form-item label="来源项目" prop="project_id">
          <el-select v-model="project_id" placeholder="请选择来源项目" clearable style="width: 100%">
            <el-option
              v-for="proj in projects"
              :key="proj.project_id"
              :label="proj.name"
              :value="proj.project_id"
            />
          </el-select>
        </el-form-item>

        <el-form-item label="费用分类" prop="category_id">
          <el-select v-model="category_id" placeholder="请选择费用分类" clearable style="width: 100%">
            <el-option
              v-for="cat in categories"
              :key="cat.category_id"
              :label="cat.name"
              :value="cat.category_id"
            />
          </el-select>
        </el-form-item>

        <el-form-item label="金额" prop="amount">
          <el-input-number v-model="amount" :min="0.01" :precision="2" :step="100" />
          <span class="unit">元</span>
        </el-form-item>

        <el-form-item label="交易时间" prop="transaction_time">
          <el-date-picker
            v-model="transaction_time"
            type="date"
            placeholder="选择日期"
            format="YYYY-MM-DD"
          />
        </el-form-item>

        <el-form-item label="备注">
          <el-input v-model="remark" type="textarea" :rows="3" placeholder="请输入备注" />
        </el-form-item>

        <el-form-item label="凭证附件" prop="attachment_ids" required>
          <el-upload
            :action="uploadAction"
            :headers="uploadHeaders"
            :on-success="handleUploadSuccess"
            :on-remove="handleUploadRemove"
            :file-list="fileList"
            :before-upload="beforeUpload"
            :limit="10"
          >
            <el-button type="primary">
              <el-icon><Upload /></el-icon>
              上传凭证
            </el-button>
            <template #tip>
              <div class="el-upload__tip">支持jpg/png/pdf格式，单个文件不超过100MB，最多10个</div>
            </template>
          </el-upload>
        </el-form-item>

        <el-form-item>
          <el-button type="primary" @click="handleSubmit" :loading="loading">提交</el-button>
          <el-button @click="handleReset">重置</el-button>
        </el-form-item>
      </el-form>
    </el-card>

    <!-- 返回按钮 -->
    <div class="back-button-container">
      <el-button @click="$router.back()" circle size="large" class="back-button">
        <el-icon><Back /></el-icon>
      </el-button>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, computed, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { createTransaction } from '@/api/transaction'
import { getCategoryList } from '@/api/category'
import { getProjectList } from '@/api/project'
import type { FormInstance, FormRules, UploadUserFile, UploadProps } from 'element-plus'
import { ElMessage } from 'element-plus'

const router = useRouter()
const formRef = ref<FormInstance>()
const loading = ref(false)

const categories = ref([])
const projects = ref([])
const fileList = ref<UploadUserFile[]>([])
const attachmentIds = ref<string[]>([])

const uploadAction = '/api/v1/attachments'
const uploadHeaders = {
  Authorization: `Bearer ${localStorage.getItem('token')}`
}

// 使用 ref 确保响应性
const project_id = ref('')
const category_id = ref('')
const amount = ref(0)
const transaction_time = ref(new Date())
const remark = ref('')
const attachment_ids = ref<string[]>([])

// 用于提交的表单对象
const form = computed(() => ({
  project_id: project_id.value,
  category_id: category_id.value,
  amount: amount.value,
  transaction_time: transaction_time.value,
  remark: remark.value,
  attachment_ids: attachment_ids.value
}))

const rules: FormRules = {
  project_id: [{ required: true, message: '请选择来源项目', trigger: 'change' }],
  amount: [{ required: true, message: '请输入金额', trigger: 'blur' }],
  transaction_time: [{ required: true, message: '请选择交易时间', trigger: 'change' }],
  attachment_ids: [
    { required: true, message: '请上传至少一个凭证附件', trigger: 'change' }
  ]
}

async function loadCategories() {
  try {
    const res = await getCategoryList({ page: 1, page_size: 100, type: 'INCOME' })
    categories.value = res.data.list || []
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

const beforeUpload: UploadProps['beforeUpload'] = (file) => {
  const isValidType = ['image/jpeg', 'image/png', 'application/pdf'].includes(file.type)
  if (!isValidType) {
    ElMessage.error('只能上传jpg/png/pdf格式的文件')
    return false
  }
  const isValidSize = file.size <= 100 * 1024 * 1024
  if (!isValidSize) {
    ElMessage.error('文件大小不能超过100MB')
    return false
  }
  return true
}

function handleUploadSuccess(response: any) {
  if (response.code === 0) {
    attachment_ids.value.push(response.data.attachment_id)
    ElMessage.success('上传成功')
  }
}

function handleUploadRemove() {
  // 重新获取attachment_ids
  attachment_ids.value = fileList.value
    .filter(f => f.status === 'success' && f.response?.code === 0)
    .map(f => f.response.data.attachment_id)
}

async function handleSubmit() {
  if (!formRef.value) return

  await formRef.value.validate(async (valid) => {
    if (!valid) return

    if (attachment_ids.value.length === 0) {
      ElMessage.warning('请上传至少一个凭证附件')
      return
    }

    loading.value = true
    try {
      const data = {
        ...form.value,
        amount: Math.abs(amount.value),
        transaction_time: (transaction_time.value as Date).toISOString().slice(0, 10) + ' 00:00:00'
      }
      await createTransaction(data)
      ElMessage.success('登记成功')
      router.push('/transaction/list')
    } catch (error) {
      console.error('Failed to create transaction:', error)
    } finally {
      loading.value = false
    }
  })
}

function handleReset() {
  formRef.value?.resetFields()
  fileList.value = []
  attachment_ids.value = []
  project_id.value = ''
  category_id.value = ''
  amount.value = 0
  transaction_time.value = new Date()
  remark.value = ''
}

onMounted(() => {
  loadCategories()
  loadProjects()
})
</script>

<style scoped>
.income-form {
  padding: 20px;
}

.unit {
  margin-left: 10px;
  color: #909399;
}

:deep(.el-input-number) {
  width: 200px;
}

/* 修复 el-select 显示问题 */
:deep(.el-select) {
  width: 100%;
}

:deep(.el-select__wrapper) {
  min-height: 32px;
  width: 100%;
}

:deep(.el-select__selection) {
  display: flex;
  flex-wrap: wrap;
  gap: 4px;
  min-height: 20px;
}

:deep(.el-select__selected-item) {
  display: inline-flex !important;
  color: var(--el-text-color-regular) !important;
  visibility: visible !important;
  opacity: 1 !important;
}

:deep(.el-select__placeholder) {
  color: var(--el-text-color-placeholder);
}

:deep(.el-select__input) {
  color: var(--el-text-color-regular) !important;
  visibility: visible !important;
}

:deep(.el-select__wrapper .el-select__selection .el-select__selected-item span) {
  display: inline !important;
  visibility: visible !important;
  color: var(--el-text-color-regular) !important;
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
