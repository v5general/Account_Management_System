<template>
  <div class="expense-form">
    <el-card>
      <template #header>
        <span>支出登记</span>
      </template>

      <el-form :model="form" :rules="rules" ref="formRef" label-width="120px">
        <el-form-item label="费用分类" prop="category_id">
          <el-select v-model="category_id" placeholder="请选择费用分类" clearable filterable style="width: 100%">
            <el-option
              v-for="cat in categories"
              :key="cat.category_id"
              :label="cat.name"
              :value="cat.category_id"
            />
          </el-select>
        </el-form-item>

        <el-form-item label="关联项目" prop="project_id">
          <el-select v-model="project_id" placeholder="请选择关联项目" clearable filterable style="width: 100%">
            <el-option
              v-for="proj in projects"
              :key="proj.project_id"
              :label="proj.name"
              :value="proj.project_id"
            />
          </el-select>
        </el-form-item>

        <el-form-item label="关联人员" prop="person_id" required>
          <el-select v-model="person_id" placeholder="请选择关联人员" filterable style="width: 100%">
            <el-option
              v-for="user in users"
              :key="user.user_id"
              :label="user.real_name"
              :value="user.user_id"
            />
          </el-select>
        </el-form-item>

        <el-form-item label="金额" prop="amount">
          <el-input-number v-model="amount" :min="0.01" :precision="2" :step="100" />
          <span class="unit">元</span>
        </el-form-item>

        <el-form-item label="支付方式" prop="payment_method_id">
          <el-select v-model="payment_method_id" placeholder="请选择支付方式" clearable filterable style="width: 100%">
            <el-option
              v-for="method in paymentMethods"
              :key="method.payment_method_id"
              :label="method.name"
              :value="method.payment_method_id"
            />
            <template #footer>
              <el-button type="primary" link @click="showAddPaymentMethodDialog">
                <el-icon><Plus /></el-icon>
                新增支付方式
              </el-button>
            </template>
          </el-select>
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

    <!-- 新增支付方式对话框 -->
    <el-dialog
      v-model="addPaymentMethodDialogVisible"
      title="新增支付方式"
      width="400px"
    >
      <el-form :model="paymentMethodForm" :rules="paymentMethodRules" ref="paymentMethodFormRef" label-width="80px">
        <el-form-item label="名称" prop="name">
          <el-input v-model="paymentMethodForm.name" placeholder="请输入支付方式名称" />
        </el-form-item>
        <el-form-item label="描述">
          <el-input v-model="paymentMethodForm.description" type="textarea" :rows="2" placeholder="请输入描述（可选）" />
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="addPaymentMethodDialogVisible = false">取消</el-button>
        <el-button type="primary" @click="handleAddPaymentMethod" :loading="addPaymentMethodLoading">确定</el-button>
      </template>
    </el-dialog>

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
import { createTransaction, resubmitTransaction } from '@/api/transaction'
import { getCategoryList } from '@/api/category'
import { getProjectList } from '@/api/project'
import { getUserList } from '@/api/user'
import { getPaymentMethodList, createPaymentMethod } from '@/api/payment_method'
import type { FormInstance, FormRules, UploadUserFile, UploadProps } from 'element-plus'
import { ElMessage } from 'element-plus'
import { Back, Upload, Plus } from '@element-plus/icons-vue'

const router = useRouter()
const formRef = ref<FormInstance>()
const loading = ref(false)
const resubmitRecordId = ref<string>('')

const categories = ref([])
const projects = ref([])
const users = ref([])
const paymentMethods = ref([])
const fileList = ref<UploadUserFile[]>([])
const attachmentIds = ref<string[]>([])

// 新增支付方式相关
const addPaymentMethodDialogVisible = ref(false)
const addPaymentMethodLoading = ref(false)
const paymentMethodFormRef = ref<FormInstance>()
const paymentMethodForm = reactive({
  name: '',
  description: ''
})
const paymentMethodRules: FormRules = {
  name: [{ required: true, message: '请输入支付方式名称', trigger: 'blur' }]
}

const uploadAction = '/api/v1/attachments'
const uploadHeaders = {
  Authorization: `Bearer ${localStorage.getItem('token')}`
}

// 使用 ref 确保响应性
const project_id = ref('')
const category_id = ref('')
const person_id = ref('')
const payment_method_id = ref('')
const amount = ref(0)
const transaction_time = ref(new Date())
const remark = ref('')
const attachment_ids = ref<string[]>([])

// 用于提交的表单对象
const form = computed(() => ({
  project_id: project_id.value,
  category_id: category_id.value,
  person_id: person_id.value,
  payment_method_id: payment_method_id.value,
  amount: amount.value,
  transaction_time: transaction_time.value,
  remark: remark.value,
  attachment_ids: attachment_ids.value
}))

const rules: FormRules = {
  category_id: [{ required: true, message: '请选择费用分类', trigger: 'change' }],
  project_id: [{ required: true, message: '请选择关联项目', trigger: 'change' }],
  person_id: [{ required: true, message: '请选择关联人员', trigger: 'change' }],
  amount: [{ required: true, message: '请输入金额', trigger: 'blur' }],
  transaction_time: [{ required: true, message: '请选择交易时间', trigger: 'change' }],
  attachment_ids: [
    { required: true, message: '请上传至少一个凭证附件', trigger: 'change' }
  ]
}

async function loadCategories() {
  try {
    const res = await getCategoryList({ page: 1, page_size: 100, type: 'EXPENSE' })
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

async function loadUsers() {
  try {
    const res = await getUserList({ page: 1, page_size: 100, role: 'EMPLOYEE' })
    users.value = res.data.list || []
  } catch (error) {
    console.error('Failed to load users:', error)
  }
}

async function loadPaymentMethods() {
  try {
    const res = await getPaymentMethodList({ page: 1, page_size: 100 })
    paymentMethods.value = res.data.list || []
  } catch (error) {
    console.error('Failed to load payment methods:', error)
  }
}

// 显示新增支付方式对话框
function showAddPaymentMethodDialog() {
  paymentMethodForm.name = ''
  paymentMethodForm.description = ''
  addPaymentMethodDialogVisible.value = true
}

// 新增支付方式
async function handleAddPaymentMethod() {
  if (!paymentMethodFormRef.value) return

  await paymentMethodFormRef.value.validate(async (valid) => {
    if (!valid) return

    addPaymentMethodLoading.value = true
    try {
      const res = await createPaymentMethod({
        name: paymentMethodForm.name,
        description: paymentMethodForm.description
      })
      ElMessage.success('创建成功')
      addPaymentMethodDialogVisible.value = false

      // 刷新支付方式列表
      await loadPaymentMethods()

      // 自动选中新创建的支付方式
      if (res.data?.payment_method_id) {
        payment_method_id.value = res.data.payment_method_id
      }
    } catch (error: any) {
      console.error('Failed to create payment method:', error)
      ElMessage.error(error.response?.data?.message || '创建失败')
    } finally {
      addPaymentMethodLoading.value = false
    }
  })
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
      const date = transaction_time.value as Date
      const formattedDate = `${date.getFullYear()}-${String(date.getMonth() + 1).padStart(2, '0')}-${String(date.getDate()).padStart(2, '0')} 00:00:00`

      // 判断是否是重新提交
      if (resubmitRecordId.value) {
        // 重新提交被驳回的记录
        const data = {
          project_id: project_id.value,
          category_id: category_id.value,
          person_id: person_id.value,
          payment_method_id: payment_method_id.value || null,
          amount: -Math.abs(amount.value),
          transaction_time: formattedDate,
          remark: remark.value
        }
        await resubmitTransaction(resubmitRecordId.value, data)
        ElMessage.success('重新提交成功')
        sessionStorage.removeItem('resubmitTransaction')
        resubmitRecordId.value = ''
      } else {
        // 创建新记录
        const data = {
          ...form.value,
          amount: -Math.abs(amount.value), // 支出为负数
          transaction_time: formattedDate
        }
        await createTransaction(data)
        ElMessage.success('登记成功')
      }
      router.push('/transaction/list')
    } catch (error) {
      console.error('Failed to submit transaction:', error)
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
  person_id.value = ''
  payment_method_id.value = ''
  amount.value = 0
  transaction_time.value = new Date()
  remark.value = ''
  resubmitRecordId.value = ''
  sessionStorage.removeItem('resubmitTransaction')
}

onMounted(() => {
  loadCategories()
  loadProjects()
  loadUsers()
  loadPaymentMethods()

  // 检查是否有重新提交的数据
  const resubmitDataStr = sessionStorage.getItem('resubmitTransaction')
  if (resubmitDataStr) {
    try {
      const resubmitData = JSON.parse(resubmitDataStr)
      // 填充表单数据
      resubmitRecordId.value = resubmitData.record_id || ''
      project_id.value = resubmitData.project_id || ''
      category_id.value = resubmitData.category_id || ''
      person_id.value = resubmitData.person_id || ''
      payment_method_id.value = resubmitData.payment_method_id || resubmitData.payment_method?.payment_method_id || ''
      amount.value = Math.abs(resubmitData.amount) || 0
      if (resubmitData.transaction_time) {
        transaction_time.value = new Date(resubmitData.transaction_time)
      }
      remark.value = resubmitData.remark || ''

      ElMessage.info('已加载被驳回的数据，请修改后重新提交')
    } catch (error) {
      console.error('Failed to parse resubmit data:', error)
    }
  }
})
</script>

<style scoped>
.expense-form {
  padding: 20px;
}

.unit {
  margin-left: 10px;
  color: #909399;
}

:deep(.el-input-number) {
  width: 200px;
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
