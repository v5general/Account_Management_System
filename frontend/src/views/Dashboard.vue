<template>
  <div class="dashboard">
    <el-row :gutter="20">
      <el-col :span="6">
        <el-card class="stat-card income">
          <div class="stat-content">
            <div class="stat-icon">
              <el-icon><TrendCharts /></el-icon>
            </div>
            <div class="stat-info">
              <p class="stat-label">总收入</p>
              <p class="stat-value">¥{{ formatAmount(summary.totalIncome) }}</p>
            </div>
          </div>
        </el-card>
      </el-col>
      <el-col :span="6">
        <el-card class="stat-card expense">
          <div class="stat-content">
            <div class="stat-icon">
              <el-icon><TrendCharts /></el-icon>
            </div>
            <div class="stat-info">
              <p class="stat-label">总支出</p>
              <p class="stat-value">¥{{ formatAmount(Math.abs(summary.totalExpense)) }}</p>
            </div>
          </div>
        </el-card>
      </el-col>
      <el-col :span="6">
        <el-card class="stat-card net">
          <div class="stat-content">
            <div class="stat-icon">
              <el-icon><DataLine /></el-icon>
            </div>
            <div class="stat-info">
              <p class="stat-label">净收支</p>
              <p class="stat-value" :class="summary.netAmount >= 0 ? 'positive' : 'negative'">
                ¥{{ formatAmount(Math.abs(summary.netAmount)) }}
              </p>
            </div>
          </div>
        </el-card>
      </el-col>
      <el-col :span="6">
        <el-card class="stat-card count">
          <div class="stat-content">
            <div class="stat-icon">
              <el-icon><Document /></el-icon>
            </div>
            <div class="stat-info">
              <p class="stat-label">记录数</p>
              <p class="stat-value">{{ summary.recordCount }}</p>
            </div>
          </div>
        </el-card>
      </el-col>
    </el-row>

    <el-row :gutter="20" style="margin-top: 20px">
      <el-col :span="12">
        <el-card>
          <template #header>
            <div class="card-header">
              <span>最近收支记录</span>
            </div>
          </template>
          <el-table :data="recentTransactions" stripe>
            <el-table-column prop="transaction_time" label="时间" width="180" />
            <el-table-column prop="project_name" label="项目" />
            <el-table-column prop="amount" label="金额" width="120">
              <template #default="{ row }">
                <span :class="row.amount > 0 ? 'income' : 'expense'">
                  {{ row.amount > 0 ? '+' : '' }}{{ row.amount }}
                </span>
              </template>
            </el-table-column>
          </el-table>
        </el-card>
      </el-col>
      <el-col :span="12">
        <el-card>
          <template #header>
            <div class="card-header">
              <span>按项目统计</span>
            </div>
          </template>
          <div ref="chartRef" style="height: 300px"></div>
        </el-card>
      </el-col>
    </el-row>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { getTransactionList } from '@/api/transaction'
import * as echarts from 'echarts'

const summary = ref({
  totalIncome: 0,
  totalExpense: 0,
  netAmount: 0,
  recordCount: 0
})

const recentTransactions = ref<any[]>([])
const chartRef = ref<HTMLElement>()

function formatAmount(amount: number) {
  return amount.toFixed(2)
}

async function loadSummary() {
  try {
    const res = await getTransactionList({ page: 1, page_size: 10 })
    recentTransactions.value = res.data.list || []

    // 计算汇总
    const list = res.data.list || []
    summary.value = {
      totalIncome: list.filter((t: any) => t.amount > 0).reduce((sum: number, t: any) => sum + t.amount, 0),
      totalExpense: list.filter((t: any) => t.amount < 0).reduce((sum: number, t: any) => sum + t.amount, 0),
      netAmount: 0,
      recordCount: res.data.total || 0
    }
    summary.value.netAmount = summary.value.totalIncome + summary.value.totalExpense
  } catch (error) {
    console.error('Failed to load summary:', error)
  }
}

function initChart() {
  if (!chartRef.value) return

  const chart = echarts.init(chartRef.value)
  chart.setOption({
    tooltip: {
      trigger: 'item'
    },
    series: [
      {
        type: 'pie',
        radius: '70%',
        data: [
          { value: 1048, name: '项目A' },
          { value: 735, name: '项目B' },
          { value: 580, name: '项目C' },
          { value: 484, name: '项目D' }
        ]
      }
    ]
  })
}

onMounted(() => {
  loadSummary()
  initChart()
})
</script>

<style scoped>
.dashboard {
  padding: 20px;
}

.stat-card {
  cursor: pointer;
  transition: transform 0.3s;
}

.stat-card:hover {
  transform: translateY(-5px);
}

.stat-content {
  display: flex;
  align-items: center;
  gap: 20px;
}

.stat-icon {
  width: 60px;
  height: 60px;
  display: flex;
  align-items: center;
  justify-content: center;
  border-radius: 8px;
  font-size: 30px;
}

.income .stat-icon {
  background-color: #ecf5ff;
  color: #409eff;
}

.expense .stat-icon {
  background-color: #fef0f0;
  color: #f56c6c;
}

.net .stat-icon {
  background-color: #f0f9ff;
  color: #67c23a;
}

.count .stat-icon {
  background-color: #fdf6ec;
  color: #e6a23c;
}

.stat-info {
  flex: 1;
}

.stat-label {
  margin: 0 0 8px 0;
  color: #909399;
  font-size: 14px;
}

.stat-value {
  margin: 0;
  font-size: 24px;
  font-weight: bold;
  color: #303133;
}

.stat-value.positive {
  color: #67c23a;
}

.stat-value.negative {
  color: #f56c6c;
}

.card-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.income {
  color: #67c23a;
}

.expense {
  color: #f56c6c;
}
</style>
