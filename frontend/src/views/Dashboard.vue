<template>
  <div class="dashboard">
    <el-row :gutter="20">
      <el-col :span="6">
        <el-card class="stat-card income" @click="showIncomeDetail">
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
        <el-card class="stat-card expense" @click="showExpenseDetail">
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
        <el-card class="stat-card net" @click="showNetDetail">
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
        <el-card class="stat-card count" @click="showCountDetail">
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
          <el-table :data="recentTransactions" stripe @row-click="showTransactionDetail" style="cursor: pointer">
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
          <div ref="chartRef" style="height: 300px; cursor: pointer" @click="showProjectDetail"></div>
        </el-card>
      </el-col>
    </el-row>

    <!-- Transaction Detail Dialog -->
    <el-dialog v-model="transactionDialogVisible" title="交易记录详情" width="600px">
      <el-descriptions :column="2" border v-if="selectedTransaction">
        <el-descriptions-item label="交易ID">{{ selectedTransaction.record_id || '-' }}</el-descriptions-item>
        <el-descriptions-item label="项目名称">{{ selectedTransaction.project_name || '-' }}</el-descriptions-item>
        <el-descriptions-item label="交易金额">
          <span :class="selectedTransaction.amount > 0 ? 'income' : 'expense'">
            {{ selectedTransaction.amount > 0 ? '+' : '' }}{{ formatAmount(selectedTransaction.amount) }}
          </span>
        </el-descriptions-item>
        <el-descriptions-item label="交易时间">{{ formatDate(selectedTransaction.transaction_time) }}</el-descriptions-item>
        <el-descriptions-item label="创建时间">{{ formatDateTime(selectedTransaction.create_time) }}</el-descriptions-item>
        <el-descriptions-item label="更新时间">{{ formatDateTime(selectedTransaction.update_time) }}</el-descriptions-item>
        <el-descriptions-item label="备注" :span="2">{{ selectedTransaction.remark || '-' }}</el-descriptions-item>
      </el-descriptions>
      <template #footer>
        <el-button @click="transactionDialogVisible = false">关闭</el-button>
      </template>
    </el-dialog>

    <!-- 收入明细对话框 -->
    <el-dialog v-model="incomeDialogVisible" title="收入明细" width="700px">
      <el-descriptions :column="1" border>
        <el-descriptions-item label="总收入金额">
          <span class="income">¥{{ formatAmount(summary.totalIncome) }}</span>
        </el-descriptions-item>
        <el-descriptions-item label="收入记录数">
          {{ incomeTransactions.length }} 笔
        </el-descriptions-item>
      </el-descriptions>
      <el-divider />
      <el-table :data="incomeTransactions" stripe max-height="400">
        <el-table-column prop="transaction_time" label="交易时间" width="180">
          <template #default="{ row }">
            {{ formatDate(row.transaction_time) }}
          </template>
        </el-table-column>
        <el-table-column prop="project_name" label="项目" />
        <el-table-column prop="amount" label="金额" width="120">
          <template #default="{ row }">
            <span class="income">+{{ row.amount }}</span>
          </template>
        </el-table-column>
      </el-table>
      <template #footer>
        <el-button @click="incomeDialogVisible = false">关闭</el-button>
      </template>
    </el-dialog>

    <!-- 支出明细对话框 -->
    <el-dialog v-model="expenseDialogVisible" title="支出明细" width="700px">
      <el-descriptions :column="1" border>
        <el-descriptions-item label="总支出金额">
          <span class="expense">¥{{ formatAmount(Math.abs(summary.totalExpense)) }}</span>
        </el-descriptions-item>
        <el-descriptions-item label="支出记录数">
          {{ expenseTransactions.length }} 笔
        </el-descriptions-item>
      </el-descriptions>
      <el-divider />
      <el-table :data="expenseTransactions" stripe max-height="400">
        <el-table-column prop="transaction_time" label="交易时间" width="180">
          <template #default="{ row }">
            {{ formatDate(row.transaction_time) }}
          </template>
        </el-table-column>
        <el-table-column prop="project_name" label="项目" />
        <el-table-column prop="amount" label="金额" width="120">
          <template #default="{ row }">
            <span class="expense">{{ row.amount }}</span>
          </template>
        </el-table-column>
      </el-table>
      <template #footer>
        <el-button @click="expenseDialogVisible = false">关闭</el-button>
      </template>
    </el-dialog>

    <!-- 净收支明细对话框 -->
    <el-dialog v-model="netDialogVisible" title="净收支明细" width="700px">
      <el-descriptions :column="1" border>
        <el-descriptions-item label="总收入">
          <span class="income">¥{{ formatAmount(summary.totalIncome) }}</span>
        </el-descriptions-item>
        <el-descriptions-item label="总支出">
          <span class="expense">¥{{ formatAmount(Math.abs(summary.totalExpense)) }}</span>
        </el-descriptions-item>
        <el-descriptions-item label="净收支">
          <span :class="summary.netAmount >= 0 ? 'income' : 'expense'">
            ¥{{ formatAmount(Math.abs(summary.netAmount)) }}
          </span>
        </el-descriptions-item>
      </el-descriptions>
      <template #footer>
        <el-button @click="netDialogVisible = false">关闭</el-button>
      </template>
    </el-dialog>

    <!-- 记录数明细对话框 -->
    <el-dialog v-model="countDialogVisible" title="记录明细" width="700px">
      <el-descriptions :column="2" border>
        <el-descriptions-item label="总记录数">
          {{ summary.recordCount }} 笔
        </el-descriptions-item>
        <el-descriptions-item label="收入记录">
          {{ incomeTransactions.length }} 笔
        </el-descriptions-item>
        <el-descriptions-item label="支出记录">
          {{ expenseTransactions.length }} 笔
        </el-descriptions-item>
      </el-descriptions>
      <el-divider />
      <el-table :data="recentTransactions" stripe max-height="400">
        <el-table-column prop="transaction_time" label="交易时间" width="180">
          <template #default="{ row }">
            {{ formatDate(row.transaction_time) }}
          </template>
        </el-table-column>
        <el-table-column prop="project_name" label="项目" />
        <el-table-column prop="amount" label="金额" width="120">
          <template #default="{ row }">
            <span :class="row.amount > 0 ? 'income' : 'expense'">
              {{ row.amount > 0 ? '+' : '' }}{{ row.amount }}
            </span>
          </template>
        </el-table-column>
      </el-table>
      <template #footer>
        <el-button @click="countDialogVisible = false">关闭</el-button>
      </template>
    </el-dialog>

    <!-- 项目明细对话框 -->
    <el-dialog v-model="projectDialogVisible" title="项目明细信息" width="700px">
      <el-table :data="projectStatistics" border>
        <el-table-column prop="project_name" label="项目名称" />
        <el-table-column prop="total_amount" label="总金额" width="150">
          <template #default="{ row }">
            <span :class="row.total_amount >= 0 ? 'income' : 'expense'">
              ¥{{ formatAmount(row.total_amount) }}
            </span>
          </template>
        </el-table-column>
        <el-table-column prop="transaction_count" label="交易笔数" width="120" />
      </el-table>
      <template #footer>
        <el-button @click="projectDialogVisible = false">关闭</el-button>
      </template>
    </el-dialog>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { getTransactionList } from '@/api/transaction'
import * as echarts from 'echarts'

const router = useRouter()

const summary = ref({
  totalIncome: 0,
  totalExpense: 0,
  netAmount: 0,
  recordCount: 0
})

const recentTransactions = ref<any[]>([])
const chartRef = ref<HTMLElement>()

// Dialog states
const transactionDialogVisible = ref(false)
const projectDialogVisible = ref(false)
const incomeDialogVisible = ref(false)
const expenseDialogVisible = ref(false)
const netDialogVisible = ref(false)
const countDialogVisible = ref(false)
const selectedTransaction = ref<any>(null)
const projectStatistics = ref<any[]>([])

// 计算收入和支出交易列表
const incomeTransactions = ref<any[]>([])
const expenseTransactions = ref<any[]>([])

function formatAmount(amount: number) {
  return amount.toFixed(2)
}

function formatDate(dateStr: string) {
  if (!dateStr) return '-'
  return dateStr.split(' ')[0]
}

function formatDateTime(dateStr: string) {
  if (!dateStr) return '-'
  return dateStr
}

function navigateTo(path: string) {
  router.push(path)
}

function showIncomeDetail() {
  incomeDialogVisible.value = true
}

function showExpenseDetail() {
  expenseDialogVisible.value = true
}

function showNetDetail() {
  netDialogVisible.value = true
}

function showCountDetail() {
  countDialogVisible.value = true
}

async function loadSummary() {
  try {
    const res = await getTransactionList({ page: 1, page_size: 10 })
    recentTransactions.value = res.data.list || []

    // 计算汇总
    const list = res.data.list || []
    incomeTransactions.value = list.filter((t: any) => t.amount >= 0)
    expenseTransactions.value = list.filter((t: any) => t.amount < 0)
    summary.value = {
      totalIncome: incomeTransactions.value.reduce((sum: number, t: any) => sum + t.amount, 0),
      totalExpense: expenseTransactions.value.reduce((sum: number, t: any) => sum + t.amount, 0),
      netAmount: 0,
      recordCount: res.data.total || 0
    }
    summary.value.netAmount = summary.value.totalIncome + summary.value.totalExpense

    // 准备项目统计数据
    calculateProjectStatistics(list)
  } catch (error) {
    console.error('Failed to load summary:', error)
  }
}

function calculateProjectStatistics(transactions: any[]) {
  const projectMap = new Map<string, { total_amount: number; transaction_count: number }>()

  transactions.forEach((t: any) => {
    const projectName = t.project_name || '未分类'
    if (!projectMap.has(projectName)) {
      projectMap.set(projectName, { total_amount: 0, transaction_count: 0 })
    }
    const stats = projectMap.get(projectName)!
    stats.total_amount += t.amount
    stats.transaction_count += 1
  })

  projectStatistics.value = Array.from(projectMap.entries()).map(([project_name, stats]) => ({
    project_name,
    total_amount: stats.total_amount,
    transaction_count: stats.transaction_count
  }))
}

function showTransactionDetail(row: any) {
  selectedTransaction.value = row
  transactionDialogVisible.value = true
}

function showProjectDetail() {
  projectDialogVisible.value = true
}

function initChart() {
  if (!chartRef.value) return

  const chart = echarts.init(chartRef.value)

  // 使用实际的项目统计数据
  const chartData = projectStatistics.value.map((p: any) => ({
    value: Math.abs(p.total_amount),
    name: p.project_name
  }))

  chart.setOption({
    tooltip: {
      trigger: 'item',
      formatter: '{b}: ¥{c} ({d}%)'
    },
    legend: {
      orient: 'vertical',
      left: 'left'
    },
    series: [
      {
        type: 'pie',
        radius: '60%',
        data: chartData.length > 0 ? chartData : [
          { value: 0, name: '暂无数据' }
        ],
        emphasis: {
          itemStyle: {
            shadowBlur: 10,
            shadowOffsetX: 0,
            shadowColor: 'rgba(0, 0, 0, 0.5)'
          }
        }
      }
    ]
  })

  // 响应式调整
  window.addEventListener('resize', () => {
    chart.resize()
  })
}

onMounted(async () => {
  await loadSummary()
  // 在数据加载后初始化图表
  setTimeout(() => {
    initChart()
  }, 100)
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
