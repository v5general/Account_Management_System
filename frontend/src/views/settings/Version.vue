<template>
  <div class="version-info">
    <el-card>
      <template #header>
        <div class="card-header">
          <span>版本记录</span>
        </div>
      </template>

      <div class="version-timeline">
        <div class="timeline-item" v-for="version in versions" :key="version.version">
          <div class="timeline-dot"></div>
          <div class="timeline-content">
            <div class="version-header">
              <span class="version-date">{{ version.date }}</span>
              <span class="version-number">{{ version.version }}</span>
              <el-tag :type="version.status === 'stable' ? 'success' : 'warning'" size="small">
                {{ version.status }}
              </el-tag>
            </div>
            <div class="version-desc">{{ version.description }}</div>

            <div class="version-details" v-if="version.features.length || version.fixes.length">
              <div class="detail-section" v-if="version.features.length">
                <div class="section-title">
                  <el-icon><Setting /></el-icon>
                  <span>功能优化</span>
                </div>
                <ul class="feature-list">
                  <li v-for="(feature, idx) in version.features" :key="'f' + idx">
                    {{ feature }}
                  </li>
                </ul>
              </div>

              <div class="detail-section" v-if="version.fixes.length">
                <div class="section-title fixes">
                  <el-icon><CircleCheck /></el-icon>
                  <span>问题修复</span>
                </div>
                <ul class="feature-list">
                  <li v-for="(fix, idx) in version.fixes" :key="'x' + idx">
                    {{ fix }}
                  </li>
                </ul>
              </div>
            </div>
          </div>
        </div>
      </div>
    </el-card>
  </div>
</template>

<script setup lang="ts">
import { ref } from 'vue'

interface VersionInfo {
  version: string
  date: string
  status: 'stable' | 'pre'
  description: string
  features: string[]
  fixes: string[]
}

const versions = ref<VersionInfo[]>([
  {
    version: 'v1.0.0',
    date: '2026-03-01',
    status: 'stable',
    description: '首个正式版本发布，包含完整的财务管理系统功能',
    features: [
      '用户管理 - 支持管理员、财务、员工三种角色',
      '部门管理 - 部门的增删改查功能',
      '项目管理 - 项目的增删改查功能',
      '收支登记 - 支持收入和支出的登记',
      '收支审核 - 管理员可审核收支记录',
      '收支列表 - 查看和筛选收支记录',
      '费用分类 - 管理费用分类信息',
      '统计报表 - 收支统计和项目统计',
      '操作日志 - 记录用户操作行为',
      '账号管理 - 员工可管理自己的账号信息',
      '首页仪表盘 - 展示收支概览和最近记录',
      '版本记录 - 展示系统版本更新信息'
    ],
    fixes: []
  }
])
</script>

<style scoped>
.version-info {
  padding: 0;
}

.card-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.version-timeline {
  position: relative;
  padding-left: 30px;
}

.version-timeline::before {
  content: '';
  position: absolute;
  left: 6px;
  top: 0;
  bottom: 0;
  width: 2px;
  background-color: #e4e7ed;
}

.timeline-item {
  position: relative;
  padding-bottom: 30px;
}

.timeline-item:last-child {
  padding-bottom: 0;
}

.timeline-dot {
  position: absolute;
  left: -24px;
  top: 4px;
  width: 14px;
  height: 14px;
  border-radius: 50%;
  background-color: #409eff;
  border: 2px solid #fff;
  box-shadow: 0 0 0 2px #409eff;
}

.timeline-content {
  background-color: #fafafa;
  border-radius: 8px;
  padding: 16px 20px;
}

.version-header {
  display: flex;
  align-items: center;
  gap: 12px;
  margin-bottom: 8px;
}

.version-date {
  color: #409eff;
  font-size: 14px;
}

.version-number {
  font-size: 16px;
  font-weight: 600;
  color: #303133;
}

.version-desc {
  color: #606266;
  font-size: 14px;
  margin-bottom: 16px;
}

.version-details {
  border-top: 1px solid #ebeef5;
  padding-top: 16px;
}

.detail-section {
  margin-bottom: 16px;
}

.detail-section:last-child {
  margin-bottom: 0;
}

.section-title {
  display: flex;
  align-items: center;
  gap: 6px;
  color: #409eff;
  font-size: 14px;
  font-weight: 500;
  margin-bottom: 10px;
}

.section-title.fixes {
  color: #e6a23c;
}

.feature-list {
  margin: 0;
  padding-left: 20px;
  color: #606266;
  font-size: 13px;
  line-height: 2;
}

.feature-list li {
  position: relative;
}

.feature-list li::marker {
  color: #909399;
}
</style>
