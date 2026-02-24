<template>
  <div class="department-manage">
    <el-card>
      <template #header>
        <div class="card-header">
          <span>部门管理</span>
          <el-button type="primary" @click="handleAdd">
            <el-icon><Plus /></el-icon>
            新增部门
          </el-button>
        </div>
      </template>

      <el-table :data="tableData" stripe v-loading="loading">
        <el-table-column prop="name" label="部门名称" width="200" />
        <el-table-column prop="description" label="描述" />
        <el-table-column prop="sort_order" label="排序" width="100" />
        <el-table-column prop="create_time" label="创建时间" width="180">
  <template #default="{ row }">
    {{ formatDateTime(row.create_time) }}
  </template>
</el-table-column>
        <el-table-column label="操作" width="150" fixed="right">
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
      :title="isEdit ? '编辑部门' : '新增部门'"
      width="500px"
    >
      <el-form :model="form" :rules="rules" ref="formRef" label-width="100px">
        <el-form-item label="部门名称" prop="name">
          <el-input v-model="form.name" placeholder="请输入部门名称" />
        </el-form-item>
        <el-form-item label="描述">
          <el-input v-model="form.description" type="textarea" placeholder="请输入部门描述" />
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
import type { FormInstance, FormRules } from 'element-plus'
import { ElMessage, ElMessageBox } from 'element-plus'
import { Plus } from '@element-plus/icons-vue'
import dayjs from 'dayjs'
import {
  getDepartmentList,
  createDepartment,
  updateDepartment,
  deleteDepartment,
  type Department
} from '@/api/department'

// 格式化时间显示
const formatDateTime = (dateTime: string) => {
  return dayjs(dateTime).format('YYYY-MM-DD HH:mm:ss')
}

const loading = ref(false)
const tableData = ref<Department[]>([])
const dialogVisible = ref(false)
const isEdit = ref(false)
const formRef = ref<FormInstance>()

const form = reactive({
  department_id: '',
  name: '',
  description: '',
  sort_order: 0
})

const rules: FormRules = {
  name: [{ required: true, message: '请输入部门名称', trigger: 'blur' }]
}

async function loadData() {
  loading.value = true
  try {
    const res = await getDepartmentList()
    tableData.value = res.data || []
  } catch (error) {
    console.error('Failed to load departments:', error)
    ElMessage.error('加载部门列表失败')
  } finally {
    loading.value = false
  }
}

function handleAdd() {
  isEdit.value = false
  Object.assign(form, { name: '', description: '', sort_order: 0 })
  dialogVisible.value = true
}

function handleEdit(row: Department) {
  isEdit.value = true
  Object.assign(form, {
    department_id: row.department_id,
    name: row.name,
    description: row.description,
    sort_order: row.sort_order
  })
  dialogVisible.value = true
}

async function handleSubmit() {
  if (!formRef.value) return

  await formRef.value.validate(async (valid) => {
    if (!valid) return

    loading.value = true
    try {
      if (isEdit.value) {
        await updateDepartment(form.department_id, {
          name: form.name,
          description: form.description,
          sort_order: form.sort_order
        })
        ElMessage.success('更新成功')
      } else {
        await createDepartment({
          name: form.name,
          description: form.description,
          sort_order: form.sort_order
        })
        ElMessage.success('创建成功')
      }
      dialogVisible.value = false
      loadData()
    } catch (error: any) {
      console.error('Failed to save department:', error)
      ElMessage.error(error.response?.data?.message || '操作失败')
    } finally {
      loading.value = false
    }
  })
}

async function handleDelete(row: Department) {
  try {
    await ElMessageBox.confirm(
      `确定要删除部门"${row.name}"吗？`,
      '提示',
      { type: 'warning' }
    )
    await deleteDepartment(row.department_id)
    ElMessage.success('删除成功')
    loadData()
  } catch (error: any) {
    if (error !== 'cancel') {
      console.error('Failed to delete department:', error)
      ElMessage.error(error.response?.data?.message || '删除失败')
    }
  }
}

onMounted(() => {
  loadData()
})
</script>

<style scoped>
.department-manage {
  padding: 20px;
}

.card-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
}
</style>
