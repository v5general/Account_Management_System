<template>
  <div class="statistics-report">
    <el-card>
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
          />
        </el-form-item>
        <el-form-item>
          <el-button type="primary" @click="handleQuery">查询</el-button>
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
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, computed } from 'vue'
import { getStatistics } from '@/api/statistics'
import type { StatisticsResponse, StatisticsParams } from '@/api/statistics'
import { ElMessage } from 'element-plus'

const loading = ref(false)
const dateRange = ref<[Date, Date]>()
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
  cycle: 'month',
  start_time: '',
  end_time: ''
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
  if (!dateRange.value) {
    ElMessage.warning('请选择时间范围')
    return
  }

  loading.value = true
  try {
    const params = {
      ...queryForm,
      start_time: dateRange.value[0].toISOString().split('T')[0],
      end_time: dateRange.value[1].toISOString().split('T')[0]
    }
    const res = await getStatistics(params)
    statistics.value = res.data
  } catch (error) {
    console.error('Failed to load statistics:', error)
  } finally {
    loading.value = false
  }
}
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
</style>
