import { createRouter, createWebHistory } from 'vue-router'
import type { RouteRecordRaw } from 'vue-router'
import { useAuthStore } from '@/stores/auth'

const routes: RouteRecordRaw[] = [
  { path: '/auth/login', name: 'Login', component: () => import('@/views/LoginPage.vue'), meta: { guest: true } },
  { path: '/auth/register', name: 'Register', component: () => import('@/views/RegisterPage.vue'), meta: { guest: true } },
  {
    path: '/app',
    component: () => import('@/views/AppShell.vue'),
    meta: { requiresAuth: true },
    children: [
      { path: 'chats', name: 'Chats', component: () => import('@/views/ChatPage.vue'), children: [] },
      { path: 'chats/:conversationId', name: 'Chat', component: () => import('@/views/ChatPage.vue') },
      { path: 'contacts', name: 'Contacts', component: () => import('@/views/ContactsPage.vue') },
      { path: 'settings', name: 'Settings', component: () => import('@/views/SettingsPage.vue') },
    ],
  },
  { path: '/', redirect: '/app/chats' },
  { path: '/:pathMatch(.*)*', redirect: '/auth/login' },
]

const router = createRouter({
  history: createWebHistory(import.meta.env.BASE_URL),
  routes,
})

router.beforeEach((to) => {
  const auth = useAuthStore()
  if (to.meta.requiresAuth && !auth.loggedIn) {
    return { path: '/auth/login', query: { redirect: to.fullPath } }
  }
  if (to.meta.guest && auth.loggedIn) {
    const redirect = (to.query.redirect as string) || '/app/chats'
    return { path: redirect }
  }
  return true
})

export default router
