<template>
  <div class="operation-log">
    <el-card>
      <template #header>
        <div class="card-header">
          <span>操作日志</span>
        </div>
      </template>

      <el-form :inline="true" :model="queryForm" class="query-form">
        <el-form-item label="操作类型">
          <el-select v-model="queryForm.operation_type" placeholder="请选择" clearable>
            <el-option label="登录" value="LOGIN" />
            <el-option label="创建" value="CREATE" />
            <el-option label="更新" value="UPDATE" />
            <el-option label="删除" value="DELETE" />
          </el-select>
        </el-form-item>
        <el-form-item label="时间范围">
          <el-date-picker
            v-model="dateRange"
            type="daterange"
            range-separator="至"
            start-placeholder="开始日期"
            end-placeholder="结束日期"
          />
        </el-form-item>
        <el-form-item>
          <el-button type="primary" @click="handleQuery">查询</el-button>
        </el-form-item>
      </el-form>

      <el-table :data="tableData" stripe v-loading="loading">
        <el-table-column prop="create_time" label="操作时间" width="180" />
        <el-table-column prop="user.username" label="操作人" width="120" />
        <el-table-column prop="operation_type" label="操作类型" width="100">
          <template #default="{ row }">
            <el-tag :type="getOperationTypeTag(row.operation_type)">
              {{ getOperationTypeLabel(row.operation_type) }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="module" label="模块" width="120" />
        <el-table-column prop="description" label="描述" show-overflow-tooltip />
        <el-table-column prop="request_ip" label="IP地址" width="150" />
        <el-table-column prop="status" label="状态" width="80">
          <template #default="{ row }">
            <el-tag :type="row.status === 1 ? 'success' : 'danger'">
              {{ row.status === 1 ? '成功' : '失败' }}
            </el-tag>
          </template>
        </el-table-column>
      </el-table>

      <el-pagination
        v-model:current-page="queryForm.page"
        v-model:page-size="queryForm.page_size"
        :total="total"
        layout="total, sizes, prev, pager, next, jumper"
        @size-change="loadData"
        @current-change="loadData"
      />
    </el-card>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, onMounted } from 'vue'

const loading = ref(false)
const tableData = ref([])
const total = ref(0)
const dateRange = ref<[Date, Date]>()

const queryForm = reactive({
  page: 1,
  page_size: 20,
  operation_type: '',
  start_time: '',
  end_time: ''
})

function getOperationTypeLabel(type: string) {
  const labels: Record<string, string> = {
    LOGIN: '登录',
    CREATE: '创建',
    UPDATE: '更新',
    DELETE: '删除'
  }
  return labels[type] || type
}

function getOperationTypeTag(type: string) {
  const tags: Record<string, string> = {
    LOGIN: 'info',
    CREATE: 'success',
    UPDATE: 'warning',
    DELETE: 'danger'
  }
  return tags[type] || 'info'
}

async function loadData() {
  loading.value = true
  try {
    // 这里应该调用获取操作日志的API
    // const res = await getOperationLogList(queryForm)
    // tableData.value = res.data.list
    // total.value = res.data.total
    // 模拟数据
    tableData.value = []
    total.value = 0
  } catch (error) {
    console.error('Failed to load logs:', error)
  } finally {
    loading.value = false
  }
}

function handleQuery() {
  queryForm.page = 1
  loadData()
}

onMounted(() => {
  loadData()
})
</script>

<style scoped>
.operation-log {
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
</style>
