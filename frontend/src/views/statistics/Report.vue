<template>
  <div class="statistics-report">
    <el-card v-loading="loading">
      <template #header>
        <div class="card-header">
          <span>统计报表</span>
        </div>
      </template>

      <el-form :inline="true" :model="queryForm" class="query-form">
        <el-form-item label="统计维度">
          <el-select v-model="queryForm.dimension">
            <el-option label="按项目" value="project" />
            <el-option label="按人员" value="person" />
            <el-option label="按分类" value="category" />
          </el-select>
        </el-form-item>
        <el-form-item label="统计周期">
          <el-select v-model="queryForm.cycle">
            <el-option label="月度" value="month" />
            <el-option label="季度" value="quarter" />
            <el-option label="年度" value="year" />
          </el-select>
        </el-form-item>
        <el-form-item label="时间范围">
          <el-date-picker
            v-model="dateRange"
            type="daterange"
            range-separator="至"
            start-placeholder="开始日期"
            end-placeholder="结束日期"
            value-format="YYYY-MM-DD"
          />
        </el-form-item>
        <el-form-item>
          <el-button type="primary" @click="handleQuery">查询</el-button>
          <el-button @click="handleReset">重置</el-button>
        </el-form-item>
      </el-form>

      <!-- 汇总数据 -->
      <el-row :gutter="20" class="summary-row">
        <el-col :span="6">
          <div class="summary-card">
            <p class="label">总收入</p>
            <p class="value income">¥{{ formatAmount(statistics.summary?.total_income || 0) }}</p>
          </div>
        </el-col>
        <el-col :span="6">
          <div class="summary-card">
            <p class="label">总支出</p>
            <p class="value expense">¥{{ formatAmount(Math.abs(statistics.summary?.total_expense || 0)) }}</p>
          </div>
        </el-col>
        <el-col :span="6">
          <div class="summary-card">
            <p class="label">净收支</p>
            <p class="value" :class="(statistics.summary?.net_amount || 0) >= 0 ? 'income' : 'expense'">
              ¥{{ formatAmount(Math.abs(statistics.summary?.net_amount || 0)) }}
            </p>
          </div>
        </el-col>
        <el-col :span="6">
          <div class="summary-card">
            <p class="label">记录数</p>
            <p class="value">{{ statistics.summary?.record_count || 0 }}</p>
          </div>
        </el-col>
      </el-row>

      <!-- 明细数据 -->
      <el-table :data="statistics.details || []" stripe class="detail-table">
        <el-table-column prop="key" :label="dimensionLabel" width="200" />
        <el-table-column prop="income" label="收入" width="150">
          <template #default="{ row }">
            <span class="income">¥{{ formatAmount(row.income) }}</span>
          </template>
        </el-table-column>
        <el-table-column prop="expense" label="支出" width="150">
          <template #default="{ row }">
            <span class="expense">¥{{ formatAmount(Math.abs(row.expense)) }}</span>
          </template>
        </el-table-column>
        <el-table-column prop="net_amount" label="净收支" width="150">
          <template #default="{ row }">
            <span :class="row.net_amount >= 0 ? 'income' : 'expense'">
              ¥{{ formatAmount(Math.abs(row.net_amount)) }}
            </span>
          </template>
        </el-table-column>
        <el-table-column prop="record_count" label="记录数" width="120" />
        <el-table-column prop="percentage" label="占比" width="120">
          <template #default="{ row }">
            {{ row.percentage?.toFixed(2) }}%
          </template>
        </el-table-column>
      </el-table>
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
import { getStatistics } from '@/api/statistics'
import type { StatisticsResponse, StatisticsParams } from '@/api/statistics'
import { ElMessage } from 'element-plus'

const loading = ref(false)
const dateRange = ref<string[]>()
const statistics = ref<StatisticsResponse>({
  dimension: '',
  cycle: '',
  start_time: '',
  end_time: '',
  summary: {
    total_income: 0,
    total_expense: 0,
    net_amount: 0,
    record_count: 0
  },
  details: []
})

const queryForm = reactive<StatisticsParams>({
  dimension: 'project',
  cycle: 'month'
})

const dimensionLabel = computed(() => {
  const labels = {
    project: '项目',
    person: '人员',
    category: '分类'
  }
  return labels[queryForm.dimension] || '项目'
})

function formatAmount(amount: number) {
  return amount.toFixed(2)
}

async function handleQuery() {
  console.log('handleQuery 开始执行')
  console.log('queryForm:', queryForm)
  console.log('dateRange:', dateRange.value)

  loading.value = true
  try {
    const params: StatisticsParams = {
      dimension: queryForm.dimension,
      cycle: queryForm.cycle
    }
    // 只有选择了时间范围才添加时间参数
    if (dateRange.value && dateRange.value.length === 2) {
      params.start_time = dateRange.value[0]
      params.end_time = dateRange.value[1]
    }
    console.log('请求参数:', params)

    const res = await getStatistics(params)
    console.log('API 响应:', res)
    console.log('统计数据:', res.data)
    console.log('details:', res.data.details)

    // 确保 details 不为 null
    if (res.data.details === null) {
      console.log('details 是 null，设置为空数组')
      res.data.details = []
    }

    statistics.value = res.data
    console.log('设置 statistics.value 后:', statistics.value)
    console.log('statistics.value.details:', statistics.value.details)
  } catch (error) {
    console.error('Failed to load statistics:', error)
    ElMessage.error('加载统计数据失败')
  } finally {
    loading.value = false
  }
}

function handleReset() {
  dateRange.value = undefined
  queryForm.dimension = 'project'
  queryForm.cycle = 'month'
  // 重置后重新加载所有数据
  handleQuery()
}

// 页面加载时自动查询所有数据
onMounted(() => {
  handleQuery()
})
</script>

<style scoped>
.statistics-report {
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

.summary-row {
  margin-bottom: 20px;
}

.summary-card {
  padding: 20px;
  background-color: #f5f7fa;
  border-radius: 4px;
  text-align: center;
}

.summary-card .label {
  margin: 0 0 10px 0;
  color: #909399;
  font-size: 14px;
}

.summary-card .value {
  margin: 0;
  font-size: 24px;
  font-weight: bold;
  color: #303133;
}

.summary-card .value.income {
  color: #67c23a;
}

.summary-card .value.expense {
  color: #f56c6c;
}

.detail-table .income {
  color: #67c23a;
}

.detail-table .expense {
  color: #f56c6c;
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
