<template>
  <div class="category-manage">
    <el-card>
      <template #header>
        <div class="card-header">
          <span>费用分类管理</span>
          <el-button type="primary" @click="handleAdd">
            <el-icon><Plus /></el-icon>
            新增分类
          </el-button>
        </div>
      </template>

      <!-- 分类类型切换 -->
      <el-tabs v-model="activeType" @tab-change="handleTypeChange">
        <el-tab-pane label="收入分类" name="INCOME"></el-tab-pane>
        <el-tab-pane label="支出分类" name="EXPENSE"></el-tab-pane>
      </el-tabs>

      <el-table :data="filteredData" stripe v-loading="loading">
        <el-table-column prop="name" label="分类名称" width="200" />
        <el-table-column prop="type" label="类型" width="100">
          <template #default="{ row }">
            <el-tag :type="row.type === 'INCOME' ? 'success' : 'warning'">
              {{ row.type === 'INCOME' ? '收入' : '支出' }}
            </el-tag>
          </template>
        </el-table-column>
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
      :title="isEdit ? '编辑分类' : '新增分类'"
      width="500px"
    >
      <el-form :model="form" :rules="rules" ref="formRef" label-width="100px">
        <el-form-item label="分类类型" prop="type">
          <el-radio-group v-model="form.type">
            <el-radio label="INCOME">收入</el-radio>
            <el-radio label="EXPENSE">支出</el-radio>
          </el-radio-group>
        </el-form-item>
        <el-form-item label="分类名称" prop="name">
          <el-input v-model="form.name" placeholder="请输入分类名称" />
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
import { ref, reactive, computed, onMounted } from 'vue'
import { getCategoryList, createCategory, updateCategory, deleteCategory, type Category, type CategoryType } from '@/api/category'
import type { FormInstance, FormRules } from 'element-plus'
import { ElMessage, ElMessageBox } from 'element-plus'
import { Plus } from '@element-plus/icons-vue'

const loading = ref(false)
const tableData = ref<Category[]>([])
const dialogVisible = ref(false)
const isEdit = ref(false)
const formRef = ref<FormInstance>()
const activeType = ref<CategoryType>('INCOME')

const form = reactive<{
  category_id?: string
  name: string
  type: CategoryType
  description: string
  sort_order: number
}>({
  name: '',
  type: 'INCOME',
  description: '',
  sort_order: 0
})

const rules: FormRules = {
  type: [{ required: true, message: '请选择分类类型', trigger: 'change' }],
  name: [{ required: true, message: '请输入分类名称', trigger: 'blur' }]
}

// 根据当前激活的类型过滤数据
const filteredData = computed(() => {
  return tableData.value.filter(item => item.type === activeType.value)
})

async function loadData() {
  loading.value = true
  try {
    const res = await getCategoryList({ page: 1, page_size: 100 })
    tableData.value = res.data.list || []
  } catch (error) {
    console.error('Failed to load categories:', error)
    ElMessage.error('加载分类列表失败')
  } finally {
    loading.value = false
  }
}

function handleTypeChange(type: CategoryType) {
  activeType.value = type
}

function handleAdd() {
  isEdit.value = false
  Object.assign(form, {
    name: '',
    type: activeType.value,
    description: '',
    sort_order: 0
  })
  dialogVisible.value = true
}

function handleEdit(row: Category) {
  isEdit.value = true
  Object.assign(form, {
    category_id: row.category_id,
    name: row.name,
    type: row.type,
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
        await updateCategory(form.category_id!, form)
        ElMessage.success('更新成功')
      } else {
        await createCategory(form)
        ElMessage.success('创建成功')
      }
      dialogVisible.value = false
      loadData()
    } catch (error: any) {
      console.error('Failed to save category:', error)
      ElMessage.error(error.response?.data?.message || '操作失败')
    } finally {
      loading.value = false
    }
  })
}

async function handleDelete(row: Category) {
  try {
    await ElMessageBox.confirm(`确定要删除分类"${row.name}"吗？`, '提示', { type: 'warning' })
    await deleteCategory(row.category_id)
    ElMessage.success('删除成功')
    loadData()
  } catch (error: any) {
    if (error !== 'cancel') {
      console.error('Failed to delete category:', error)
      ElMessage.error(error.response?.data?.message || '删除失败')
    }
  }
}

onMounted(() => {
  loadData()
})
</script>

<style scoped>
.category-manage {
  padding: 20px;
}

.card-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
}
</style>
