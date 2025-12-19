import { createRouter, createWebHistory } from 'vue-router'
import { useAuthStore } from '@/stores/auth'

const router = createRouter({
  history: createWebHistory(import.meta.env.BASE_URL),
  routes: [
    {
      path: '/login',
      name: 'login',
      component: () => import('@/views/auth/LoginView.vue'),
      meta: { requiresAuth: false }
    },
    {
      path: '/register',
      name: 'register',
      component: () => import('@/views/auth/RegisterView.vue'),
      meta: { requiresAuth: false }
    },
    {
      path: '/',
      component: () => import('@/components/layout/AppLayout.vue'),
      meta: { requiresAuth: true },
      children: [
        {
          path: '',
          name: 'dashboard',
          component: () => import('@/views/dashboard/DashboardView.vue')
        },
        {
          path: 'chat',
          name: 'chat',
          component: () => import('@/views/chat/ChatView.vue')
        },
        {
          path: 'chat/:contactId',
          name: 'chat-conversation',
          component: () => import('@/views/chat/ChatView.vue'),
          props: true
        },
        {
          path: 'templates',
          name: 'templates',
          component: () => import('@/views/settings/TemplatesView.vue')
        },
        {
          path: 'flows',
          name: 'flows',
          component: () => import('@/views/settings/FlowsView.vue')
        },
        {
          path: 'campaigns',
          name: 'campaigns',
          component: () => import('@/views/settings/CampaignsView.vue')
        },
        {
          path: 'chatbot',
          name: 'chatbot',
          component: () => import('@/views/chatbot/ChatbotView.vue')
        },
        {
          path: 'chatbot/keywords',
          name: 'chatbot-keywords',
          component: () => import('@/views/chatbot/KeywordsView.vue')
        },
        {
          path: 'chatbot/flows',
          name: 'chatbot-flows',
          component: () => import('@/views/chatbot/ChatbotFlowsView.vue')
        },
        {
          path: 'chatbot/ai',
          name: 'chatbot-ai',
          component: () => import('@/views/chatbot/AIContextsView.vue')
        },
        {
          path: 'settings',
          name: 'settings',
          component: () => import('@/views/settings/SettingsView.vue')
        },
        {
          path: 'settings/accounts',
          name: 'accounts',
          component: () => import('@/views/settings/AccountsView.vue')
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

// Navigation guard
router.beforeEach(async (to, from, next) => {
  const authStore = useAuthStore()

  // Check if route requires auth
  if (to.meta.requiresAuth !== false) {
    if (!authStore.isAuthenticated) {
      // Try to restore session from localStorage
      const restored = authStore.restoreSession()
      if (!restored) {
        return next({ name: 'login', query: { redirect: to.fullPath } })
      }
    }
  } else {
    // Redirect to dashboard if already logged in
    if (authStore.isAuthenticated && (to.name === 'login' || to.name === 'register')) {
      return next({ name: 'dashboard' })
    }
  }

  next()
})

export default router
