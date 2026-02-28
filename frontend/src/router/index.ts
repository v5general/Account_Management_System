import { createRouter, createWebHistory } from 'vue-router'
import type { RouteRecordRaw } from 'vue-router'
import { useUserStore } from '@/store/user'

const routes: RouteRecordRaw[] = [
  {
    path: '/login',
    name: 'Login',
    component: () => import('../views/Login.vue'),
    meta: { requiresAuth: false }
  },
  {
    path: '/',
    component: () => import('../components/Layout.vue'),
    redirect: '/dashboard',
    meta: { requiresAuth: true },
    children: [
      {
        path: 'dashboard',
        name: 'Dashboard',
        component: () => import('../views/Dashboard.vue'),
        meta: { title: '首页仪表盘' }
      },
      {
        path: 'transaction',
        name: 'Transaction',
        redirect: '/transaction/list',
        meta: { title: '收支管理' },
        children: [
          {
            path: 'list',
            name: 'TransactionList',
            component: () => import('../views/transaction/List.vue'),
            meta: { title: '收支列表' }
          },
          {
            path: 'income',
            name: 'Income',
            component: () => import('../views/transaction/Income.vue'),
            meta: { title: '收入登记', roles: ['ADMIN', 'FINANCE'] }
          },
          {
            path: 'expense',
            name: 'Expense',
            component: () => import('../views/transaction/Expense.vue'),
            meta: { title: '支出登记', roles: ['ADMIN', 'FINANCE'] }
          }
        ]
      },
      {
        path: 'category',
        name: 'Category',
        component: () => import('../views/category/Manage.vue'),
        meta: { title: '费用分类', roles: ['ADMIN', 'FINANCE'] }
      },
      {
        path: 'statistics',
        name: 'Statistics',
        component: () => import('../views/statistics/Report.vue'),
        meta: { title: '统计报表', roles: ['ADMIN', 'FINANCE'] }
      },
      {
        path: 'settings',
        name: 'Settings',
        redirect: '/settings/user',
        meta: { title: '系统设置' },
        children: [
          {
            path: 'user',
            name: 'UserManage',
            component: () => import('../views/settings/User.vue'),
            meta: { title: '用户管理', roles: ['ADMIN'] }
          },
          {
            path: 'department',
            name: 'DepartmentManage',
            component: () => import('../views/settings/Department.vue'),
            meta: { title: '部门管理', roles: ['ADMIN'] }
          },
          {
            path: 'project',
            name: 'ProjectManage',
            component: () => import('../views/settings/Project.vue'),
            meta: { title: '项目管理', roles: ['ADMIN'] }
          },
          {
            path: 'log',
            name: 'OperationLog',
            component: () => import('../views/settings/Log.vue'),
            meta: { title: '操作日志', roles: ['ADMIN', 'FINANCE'] }
          },
          {
            path: 'account',
            name: 'AccountManage',
            component: () => import('../views/settings/Account.vue'),
            meta: { title: '账号管理', roles: ['EMPLOYEE', 'FINANCE'] }
          }
        ]
      }
    ]
  },
  {
    path: '/:pathMatch(.*)*',
    redirect: '/dashboard'
  }
]

const router = createRouter({
  history: createWebHistory(),
  routes
})

// 路由守卫
router.beforeEach(async (to, from, next) => {
  const token = localStorage.getItem('token')
  const userStore = useUserStore()

  // 需要认证但没有token，跳转登录页
  if (to.meta.requiresAuth && !token) {
    next('/login')
    return
  }

  // 如果有token但没有用户信息，先获取用户信息
  if (token && !userStore.userInfo) {
    try {
      await userStore.fetchUserInfo()
    } catch {
      // 获取用户信息失败，清除token并跳转登录页
      userStore.clearToken()
      next('/login')
      return
    }
  }

  // 权限检查
  if (to.meta.roles && userStore.userInfo) {
    const roles = to.meta.roles as string[]
    if (!roles.includes(userStore.userInfo.role)) {
      // 没有权限，跳转首页
      next('/dashboard')
      return
    }
  }

  // 系统设置页面特殊处理
  if (to.path.startsWith('/settings')) {
    if (userStore.isEmployee() && to.path !== '/settings/account') {
      next('/settings/account')
      return
    }
  }

  next()
})

export default router
