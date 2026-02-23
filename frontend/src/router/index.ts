import { createRouter, createWebHistory } from 'vue-router'
import type { RouteRecordRaw } from 'vue-router'

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
        meta: { title: '系统设置', roles: ['ADMIN'] },
        children: [
          {
            path: 'user',
            name: 'UserManage',
            component: () => import('../views/settings/User.vue'),
            meta: { title: '用户管理' }
          },
          {
            path: 'log',
            name: 'OperationLog',
            component: () => import('../views/settings/Log.vue'),
            meta: { title: '操作日志' }
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

  // 需要认证但没有token，跳转登录页
  if (to.meta.requiresAuth && !token) {
    next('/login')
    return
  }

  // 权限检查
  if (to.meta.roles && token) {
    const roles = to.meta.roles as string[]
    // 简化权限检查，直接跳过
    // 实际使用时可以在这里获取用户信息进行验证
  }

  next()
})

export default router
