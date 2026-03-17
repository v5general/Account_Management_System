<template>
  <el-container class="layout-container">
    <el-aside :width="sidebarWidth">
      <div class="logo">
        <h2 v-if="!appStore.sidebarCollapsed">账务管理系统</h2>
        <h2 v-else>账务</h2>
      </div>
      <el-menu
        :default-active="currentRoute"
        :collapse="appStore.sidebarCollapsed"
        router
      >
        <el-menu-item index="/dashboard">
          <el-icon><HomeFilled /></el-icon>
          <span>首页</span>
        </el-menu-item>
        <el-sub-menu index="/transaction" v-if="!userStore.isEmployee()">
          <template #title>
            <el-icon><Money /></el-icon>
            <span>收支管理</span>
          </template>
          <el-menu-item index="/transaction/income" v-if="userStore.isFinance()">收入登记</el-menu-item>
          <el-menu-item index="/transaction/expense" v-if="userStore.isFinance()">支出登记</el-menu-item>
          <el-menu-item index="/transaction/audit" v-if="userStore.isAdmin()">收支审核</el-menu-item>
          <el-menu-item index="/transaction/list">收支列表</el-menu-item>
        </el-sub-menu>
        <el-menu-item index="/transaction/list" v-if="userStore.isEmployee()">
          <el-icon><Money /></el-icon>
          <span>我的收支</span>
        </el-menu-item>
        <el-menu-item index="/category" v-if="userStore.isFinance()">
          <el-icon><Folder /></el-icon>
          <span>费用分类</span>
        </el-menu-item>
        <el-menu-item index="/statistics" v-if="userStore.isFinance()">
          <el-icon><DataAnalysis /></el-icon>
          <span>统计报表</span>
        </el-menu-item>
        <el-sub-menu index="/settings">
          <template #title>
            <el-icon><Setting /></el-icon>
            <span>系统设置</span>
          </template>
          <!-- 员工和财务显示账号管理 -->
          <el-menu-item index="/settings/account" v-if="!userStore.isAdmin()">账号管理</el-menu-item>
          <!-- 财务和管理员显示支付方式 -->
          <el-menu-item index="/settings/payment-method" v-if="userStore.isFinance()">支付方式</el-menu-item>
          <!-- 管理员显示完整菜单 -->
          <template v-if="userStore.isAdmin()">
            <el-menu-item index="/settings/user">用户管理</el-menu-item>
            <el-menu-item index="/settings/department">部门管理</el-menu-item>
            <el-menu-item index="/settings/project">项目管理</el-menu-item>
            <el-menu-item index="/settings/log">操作日志</el-menu-item>
          </template>
          <el-menu-item index="/settings/version">版本记录</el-menu-item>
        </el-sub-menu>
      </el-menu>
    </el-aside>
    <el-container>
      <el-header>
        <div class="header-left">
          <el-icon class="collapse-icon" @click="appStore.toggleSidebar">
            <Fold v-if="!appStore.sidebarCollapsed" />
            <Expand v-else />
          </el-icon>
          <el-breadcrumb separator="/">
            <el-breadcrumb-item :to="{ path: '/dashboard' }">首页</el-breadcrumb-item>
            <el-breadcrumb-item v-if="currentRouteMeta">{{ currentRouteMeta }}</el-breadcrumb-item>
          </el-breadcrumb>
        </div>
        <div class="header-right">
          <el-dropdown @command="handleCommand">
            <span class="user-info">
              <el-icon><User /></el-icon>
              <span style="margin-left: 8px">{{ displayName }}</span>
            </span>
            <template #dropdown>
              <el-dropdown-menu>
                <el-dropdown-item command="logout" class="logout-dropdown-item">退出登录</el-dropdown-item>
              </el-dropdown-menu>
            </template>
          </el-dropdown>
        </div>
      </el-header>
      <el-main>
        <router-view />
      </el-main>
    </el-container>
  </el-container>
</template>

<script setup lang="ts">
import { computed } from 'vue'
import { useRoute } from 'vue-router'
import { useAppStore } from '@/store/app'
import { useUserStore } from '@/store/user'
import { useRouter } from 'vue-router'
import { ElMessageBox } from 'element-plus'

const route = useRoute()
const router = useRouter()
const appStore = useAppStore()
const userStore = useUserStore()

const sidebarWidth = computed(() => appStore.sidebarCollapsed ? '64px' : '200px')
const currentRoute = computed(() => route.path)
const currentRouteMeta = computed(() => route.meta.title as string)

// 显示用户名和角色
const displayName = computed(() => {
  const user = userStore.userInfo
  if (!user) return '未登录'
  const roleName = {
    ADMIN: '管理员',
    FINANCE: '财务',
    EMPLOYEE: '员工'
  }[user.role] || user.role
  return `${user.real_name}（${roleName}）`
})

async function handleCommand(command: string) {
  if (command === 'logout') {
    try {
      await ElMessageBox.confirm('确定要退出登录吗？', '提示', {
        type: 'warning'
      })
      localStorage.removeItem('token')
      router.push('/login')
    } catch {
      // 取消操作
    }
  }
}
</script>

<style scoped>
.layout-container {
  height: 100%;
}

.el-aside {
  background-color: #304156;
  transition: width 0.3s;
}

.logo {
  height: 60px;
  display: flex;
  align-items: center;
  justify-content: center;
  color: #fff;
  font-size: 18px;
  border-bottom: 1px solid #1f2d3d;
}

.logo h2 {
  margin: 0;
  font-size: 18px;
}

.el-menu {
  border-right: none;
}

.el-header {
  background-color: #fff;
  border-bottom: 1px solid #e6e6e6;
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 0 20px;
}

.header-left {
  display: flex;
  align-items: center;
  gap: 20px;
}

.collapse-icon {
  font-size: 20px;
  cursor: pointer;
}

.header-right {
  display: flex;
  align-items: center;
}

.user-info {
  display: flex;
  align-items: center;
  cursor: pointer;
  border: none;
  outline: none;
}

.user-info:focus {
  outline: none;
}

.el-main {
  background-color: #f0f2f5;
  padding: 20px;
}

/* 移除 dropdown 的边框 */
.header-right :deep(.el-dropdown) {
  border: none;
  outline: none;
}

.header-right :deep(.el-dropdown:focus) {
  outline: none;
}

.header-right :deep(.el-dropdown:focus-visible) {
  outline: none;
}

/* 退出登录红色样式 */
:deep(.logout-dropdown-item) {
  color: #f56c6c !important;
}

:deep(.logout-dropdown-item:hover) {
  color: #f56c6c !important;
  background-color: #fef0f0 !important;
}
</style>
