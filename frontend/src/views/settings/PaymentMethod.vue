<template>
  <div class="payment-method-manage">
    <el-card>
      <template #header>
        <div class="card-header">
          <span>支付方式管理</span>
          <el-button type="primary" @click="handleAdd">
            <el-icon><Plus /></el-icon>
            新增支付方式
          </el-button>
        </div>
      </template>

      <el-table :data="tableData" stripe v-loading="loading">
        <el-table-column prop="name" label="支付方式名称" width="200" />
        <el-table-column prop="description" label="描述" show-overflow-tooltip />
        <el-table-column prop="sort_order" label="排序" width="100" />
        <el-table-column prop="create_time" label="创建时间" width="180" />
        <el-table-column label="操作" width="200" fixed="right">
          <template #default="{ row }">
            <el-button link type="primary" @click="handleEdit(row)">编辑</el-button>
            <el-button link type="danger" @click="handleDelete(row)">删除</el-button>
          </template>
        </el-table-column>
      </el-table>
    </el-card>

    <!-- 编辑对话框 -->
    <el-dialog
      v-model="dialogVisible"
      :title="isEdit ? '编辑支付方式' : '新增支付方式'"
      width="500px"
    >
      <el-form :model="form" :rules="rules" ref="formRef" label-width="100px">
        <el-form-item label="名称" prop="name">
          <el-input v-model="form.name" placeholder="请输入支付方式名称" />
        </el-form-item>
        <el-form-item label="描述">
          <el-input v-model="form.description" type="textarea" :rows="3" placeholder="请输入描述" />
        </el-form-item>
        <el-form-item label="排序">
          <el-input-number v-model="form.sort_order" :min="0" />
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="dialogVisible = false">取消</el-button>
        <el-button type="primary" @click="handleSubmit" :loading="loading">确定</el-button>
      </template>
    </el-dialog>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, onMounted } from 'vue'
import { getPaymentMethodList, createPaymentMethod, updatePaymentMethod, deletePaymentMethod, type PaymentMethod } from '@/api/payment_method'
import type { FormInstance, FormRules } from 'element-plus'
import { ElMessage, ElMessageBox } from 'element-plus'
import { Plus } from '@element-plus/icons-vue'

const loading = ref(false)
const tableData = ref<PaymentMethod[]>([])
const dialogVisible = ref(false)
const isEdit = ref(false)
const formRef = ref<FormInstance>()

const form = reactive<{
  payment_method_id?: string
  name: string
  description: string
  sort_order: number
}>({
  name: '',
  description: '',
  sort_order: 0
})

const rules: FormRules = {
  name: [{ required: true, message: '请输入支付方式名称', trigger: 'blur' }]
}

async function loadData() {
  loading.value = true
  try {
    const res = await getPaymentMethodList({ page: 1, page_size: 100 })
    tableData.value = res.data.list || []
  } catch (error) {
    console.error('Failed to load payment methods:', error)
    ElMessage.error('加载支付方式列表失败')
  } finally {
    loading.value = false
  }
}

function handleAdd() {
  isEdit.value = false
  Object.assign(form, {
    name: '',
    description: '',
    sort_order: 0
  })
  dialogVisible.value = true
}

function handleEdit(row: PaymentMethod) {
  isEdit.value = true
  Object.assign(form, {
    payment_method_id: row.payment_method_id,
    name: row.name,
    description: row.description,
    sort_order: row.sort_order
  })
  dialogVisible.value = true
}

async function handleSubmit() {
  if (!formRef.value) return

  try {
    await formRef.value.validate()
  } catch {
    return
  }

  loading.value = true
  try {
    if (isEdit.value) {
      await updatePaymentMethod(form.payment_method_id!, form)
      ElMessage.success('更新成功')
    } else {
      await createPaymentMethod(form)
      ElMessage.success('创建成功')
    }
    dialogVisible.value = false
    loadData()
  } catch (error: any) {
    console.error('Failed to save payment method:', error)
    ElMessage.error(error.response?.data?.message || '操作失败')
  } finally {
    loading.value = false
  }
}

async function handleDelete(row: PaymentMethod) {
  try {
    await ElMessageBox.confirm(`确定要删除支付方式"${row.name}"吗？`, '提示', { type: 'warning' })
    await deletePaymentMethod(row.payment_method_id)
    ElMessage.success('删除成功')
    loadData()
  } catch (error: any) {
    if (error !== 'cancel') {
      console.error('Failed to delete payment method:', error)
      ElMessage.error(error.response?.data?.message || '删除失败')
    }
  }
}

onMounted(() => {
  loadData()
})
</script>

<style scoped>
.payment-method-manage {
  padding: 20px;
}

.card-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
}
</style>
