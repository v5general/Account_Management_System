<template>
  <div class="project-manage">
    <el-card>
      <template #header>
        <div class="card-header">
          <span>项目管理</span>
          <el-button type="primary" @click="handleAdd">
            <el-icon><Plus /></el-icon>
            新增项目
          </el-button>
        </div>
      </template>

      <el-table :data="tableData" stripe v-loading="loading">
        <el-table-column prop="name" label="项目名称" width="200" />
        <el-table-column prop="description" label="描述" />
        <el-table-column label="所属部门" width="150">
          <template #default="{ row }">
            {{ getDepartmentName(row.department_id) }}
          </template>
        </el-table-column>
        <el-table-column label="状态" width="100">
          <template #default="{ row }">
            <el-tag :type="row.status === 1 ? 'success' : 'info'">
              {{ row.status === 1 ? '进行中' : '已结束' }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="create_time" label="创建时间" width="180" />
        <el-table-column label="操作" width="180" fixed="right">
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
      :title="isEdit ? '编辑项目' : '新增项目'"
      width="500px"
    >
      <el-form :model="form" :rules="rules" ref="formRef" label-width="100px">
        <el-form-item label="项目名称" prop="name">
          <el-input v-model="form.name" placeholder="请输入项目名称" />
        </el-form-item>
        <el-form-item label="所属部门">
          <el-select v-model="form.department_id" placeholder="请选择部门" clearable>
            <el-option
              v-for="dept in departments"
              :key="dept.department_id"
              :label="dept.name"
              :value="dept.department_id"
            />
          </el-select>
        </el-form-item>
        <el-form-item label="描述">
          <el-input v-model="form.description" type="textarea" :rows="3" placeholder="请输入项目描述" />
        </el-form-item>
        <el-form-item label="状态">
          <el-radio-group v-model="form.status">
            <el-radio :value="1">进行中</el-radio>
            <el-radio :value="0">已结束</el-radio>
          </el-radio-group>
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
import {
  getProjectList,
  createProject,
  updateProject,
  deleteProject,
  type Project
} from '@/api/project'
import { getDepartmentList, type Department } from '@/api/department'

const loading = ref(false)
const tableData = ref<Project[]>([])
const departments = ref<Department[]>([])
const dialogVisible = ref(false)
const isEdit = ref(false)
const formRef = ref<FormInstance>()

const form = reactive({
  project_id: '',
  name: '',
  description: '',
  department_id: '',
  status: 1
})

const rules: FormRules = {
  name: [{ required: true, message: '请输入项目名称', trigger: 'blur' }]
}

function getDepartmentName(departmentId: string) {
  const dept = departments.value.find(d => d.department_id === departmentId)
  return dept ? dept.name : '-'
}

async function loadData() {
  loading.value = true
  try {
    const res = await getProjectList()
    tableData.value = res.data.list || []
  } catch (error) {
    console.error('Failed to load projects:', error)
    ElMessage.error('加载项目列表失败')
  } finally {
    loading.value = false
  }
}

async function loadDepartments() {
  try {
    const res = await getDepartmentList()
    departments.value = res.data || []
  } catch (error) {
    console.error('Failed to load departments:', error)
  }
}

function handleAdd() {
  isEdit.value = false
  Object.assign(form, { name: '', description: '', department_id: '', status: 1 })
  dialogVisible.value = true
}

function handleEdit(row: Project) {
  isEdit.value = true
  Object.assign(form, {
    project_id: row.project_id,
    name: row.name,
    description: row.description,
    department_id: row.department_id,
    status: row.status
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
        await updateProject(form.project_id, {
          name: form.name,
          description: form.description,
          department_id: form.department_id,
          status: form.status
        })
        ElMessage.success('更新成功')
      } else {
        await createProject({
          name: form.name,
          description: form.description,
          department_id: form.department_id
        })
        ElMessage.success('创建成功')
      }
      dialogVisible.value = false
      loadData()
    } catch (error: any) {
      console.error('Failed to save project:', error)
      ElMessage.error(error.response?.data?.message || '操作失败')
    } finally {
      loading.value = false
    }
  })
}

async function handleDelete(row: Project) {
  try {
    await ElMessageBox.confirm(
      `确定要删除项目"${row.name}"吗？`,
      '提示',
      { type: 'warning' }
    )
    await deleteProject(row.project_id)
    ElMessage.success('删除成功')
    loadData()
  } catch (error: any) {
    if (error !== 'cancel') {
      console.error('Failed to delete project:', error)
      ElMessage.error(error.response?.data?.message || '删除失败')
    }
  }
}

onMounted(() => {
  loadData()
  loadDepartments()
})
</script>

<style scoped>
.project-manage {
  padding: 20px;
}

.card-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
}
</style>
