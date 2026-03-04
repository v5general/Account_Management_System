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

            <!-- 下拉式详情 -->
            <el-collapse v-if="version.features.length || version.fixes.length" class="version-collapse">
              <el-collapse-item>
                <template #title>
                  <div class="collapse-title">
                    <el-icon><Document /></el-icon>
                    <span>查看更新详情</span>
                    <el-tag size="small" type="info">
                      {{ version.features.length + version.fixes.length }}项
                    </el-tag>
                  </div>
                </template>

                <div class="version-details">
                  <div class="detail-section" v-if="version.features.length">
                    <div class="section-title">
                      <el-icon><Setting /></el-icon>
                      <span>功能优化</span>
                      <el-tag size="small" type="success">{{ version.features.length }}</el-tag>
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
                      <el-tag size="small" type="warning">{{ version.fixes.length }}</el-tag>
                    </div>
                    <ul class="feature-list">
                      <li v-for="(fix, idx) in version.fixes" :key="'x' + idx">
                        {{ fix }}
                      </li>
                    </ul>
                  </div>
                </div>
              </el-collapse-item>
            </el-collapse>
          </div>
        </div>
      </div>
    </el-card>
  </div>
</template>

<script setup lang="ts">
import { ref } from 'vue'
import { Document } from '@element-plus/icons-vue'

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
    version: 'v1.0.2',
    date: '2026-03-05',
    status: 'stable',
    description: '功能增强与显示优化版本',
    features: [
      '支付方式 - 新增支付方式字段，支持现金、微信、支付宝、银行转账等多种支付方式',
      '重新提交 - 实现被驳回记录的重新提交功能，提升审核流程效率',
      '金额显示 - 优化金额格式化，采用千位分隔符显示（如：51,443.02）',
      '格式化工具 - 统一使用格式化工具函数，提升代码一致性',
      '表格优化 - 调整表格列宽，优化信息展示效果'
    ],
    fixes: []
  },
  {
    version: 'v1.0.1',
    date: '2026-03-01',
    status: 'stable',
    description: '界面优化与问题修复版本',
    features: [
      '登录界面 - 采用左右分栏布局，左侧品牌展示，右侧登录表单',
      '登录界面 - 添加背景装饰元素和渐变色按钮悬浮效果',
      '登录界面 - 响应式设计，小屏幕自动隐藏左侧区域',
      '登录界面 - 底部显示当前版本号',
      '收支列表 - 界面标题从"收支记录"优化为"收支列表"'
    ],
    fixes: [
      '首页记录数 - 修复点击记录数后显示所有记录的问题，改为只显示审核通过的记录',
      '用户模型 - 修复外键关联显式指定 references',
      '部署文档 - 更新构建说明，添加交叉编译指南'
    ]
  },
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
  margin-bottom: 12px;
}

/* 折叠面板样式 */
.version-collapse {
  margin-top: 12px;
  border: none;
}

.version-collapse :deep(.el-collapse-item__header) {
  background-color: #f0f2f5;
  border: 1px solid #e4e7ed;
  border-radius: 6px;
  padding: 0 16px;
  height: 40px;
  line-height: 40px;
  font-size: 14px;
  color: #606266;
  transition: all 0.3s;
}

.version-collapse :deep(.el-collapse-item__header:hover) {
  background-color: #e8eaed;
  color: #409eff;
}

.version-collapse :deep(.el-collapse-item__wrap) {
  border: none;
  background-color: transparent;
}

.version-collapse :deep(.el-collapse-item__content) {
  padding-top: 16px;
}

.collapse-title {
  display: flex;
  align-items: center;
  gap: 8px;
  font-weight: 500;
}

.collapse-title .el-icon {
  font-size: 16px;
}

.version-details {
  background-color: #fff;
  border-radius: 6px;
  padding: 16px;
  border: 1px solid #ebeef5;
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
  gap: 8px;
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
