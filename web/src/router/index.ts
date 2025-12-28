import { createRouter, createWebHistory } from 'vue-router'

const router = createRouter({
  history: createWebHistory(import.meta.env.BASE_URL),
  routes: [
    {
      path: '/login',
      name: 'login',
      component: () => import('@/views/auth/LoginView.vue'),
      meta: { guest: true }
    },
    {
      path: '/register',
      name: 'register',
      component: () => import('@/views/auth/RegisterView.vue'),
      meta: { guest: true }
    },
    {
      path: '/',
      component: () => import('@/layouts/MainLayout.vue'),
      meta: { requiresAuth: true },
      children: [
        {
          path: '',
          name: 'dashboard',
          component: () => import('@/views/DashboardView.vue')
        },
        {
          path: 'tasks',
          name: 'tasks',
          component: () => import('@/views/TaskListView.vue')
        },
        {
          path: 'tasks/new',
          name: 'task-create',
          component: () => import('@/views/TaskCreateView.vue')
        },
        {
          path: 'tasks/:id',
          name: 'task-detail',
          component: () => import('@/views/TaskDetailView.vue')
        },
        {
          path: 'settings',
          name: 'settings',
          component: () => import('@/views/SettingsView.vue')
        },
        {
          path: 'settings/smtp',
          name: 'settings-smtp',
          component: () => import('@/views/SmtpConfigView.vue')
        },
        {
          path: 'settings/webhooks',
          name: 'settings-webhooks',
          component: () => import('@/views/WebhookConfigView.vue')
        },
        {
          path: 'admin',
          meta: { requiresAdmin: true },
          children: [
            {
              path: 'users',
              name: 'admin-users',
              component: () => import('@/views/admin/UserManagementView.vue')
            },
            {
              path: 'system',
              name: 'admin-system',
              component: () => import('@/views/admin/SystemConfigView.vue')
            },
            {
              path: 'audit',
              name: 'admin-audit',
              component: () => import('@/views/admin/AuditLogView.vue')
            }
          ]
        }
      ]
    },
    {
      path: '/:pathMatch(.*)*',
      name: 'not-found',
      component: () => import('@/views/NotFoundView.vue')
    }
  ]
})

function getAuthFromStorage() {
  try {
    const stored = localStorage.getItem('auth')
    if (stored) {
      return JSON.parse(stored)
    }
  } catch {
    // ignore
  }
  return null
}

router.beforeEach((to, _from, next) => {
  const auth = getAuthFromStorage()
  const isAuthenticated = !!auth?.token

  if (to.meta.requiresAuth && !isAuthenticated) {
    return next({ name: 'login', query: { redirect: to.fullPath } })
  }

  if (to.meta.guest && isAuthenticated) {
    return next({ name: 'dashboard' })
  }

  if (to.meta.requiresAdmin && auth?.user?.role !== 'admin') {
    return next({ name: 'dashboard' })
  }

  next()
})

export default router
