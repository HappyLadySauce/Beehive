<template>
  <n-layout has-sider class="app-shell">
    <n-layout-sider bordered collapse-mode="width" :collapsed-width="64" :width="240" show-trigger>
      <div class="sidebar-header">
        <span class="logo">Beehive</span>
      </div>
      <n-menu
        :options="menuOptions"
        :value="currentMenu"
        @update:value="handleMenuSelect"
      />
    </n-layout-sider>
    <n-layout>
      <n-layout-header bordered class="header">
        <span class="header-title">{{ headerTitle }}</span>
        <n-dropdown :options="userOptions" trigger="click" @select="handleUserSelect">
          <n-button quaternary circle>
            <n-avatar round size="small">{{ userInitial }}</n-avatar>
          </n-button>
        </n-dropdown>
      </n-layout-header>
      <n-layout-content content-style="height: 100%; overflow: auto" class="content">
        <router-view />
      </n-layout-content>
    </n-layout>
  </n-layout>
</template>

<script setup lang="ts">
import { computed } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { NLayout, NLayoutSider, NLayoutHeader, NLayoutContent, NMenu, NButton, NAvatar, NDropdown } from 'naive-ui'
import type { MenuOption } from 'naive-ui'
import { useAuthStore } from '@/stores/auth'

const route = useRoute()
const router = useRouter()
const authStore = useAuthStore()

const menuOptions: MenuOption[] = [
  { label: '会话', key: 'chats', path: '/app/chats' },
  { label: '联系人', key: 'contacts', path: '/app/contacts' },
  { label: '设置', key: 'settings', path: '/app/settings' },
]

const currentMenu = computed(() => {
  const p = route.path
  if (p.startsWith('/app/chats')) return 'chats'
  if (p.startsWith('/app/contacts')) return 'contacts'
  if (p.startsWith('/app/settings')) return 'settings'
  return 'chats'
})

const headerTitle = computed(() => {
  if (route.path.startsWith('/app/settings')) return '设置'
  if (route.path.startsWith('/app/contacts')) return '联系人'
  return '会话'
})

const userInitial = computed(() => {
  const u = authStore.user
  if (u?.nickname) return u.nickname.slice(0, 1).toUpperCase()
  if (u?.id) return u.id.slice(-1).toUpperCase()
  return '?'
})

const userOptions = [
  { label: '退出登录', key: 'logout' },
]

function handleMenuSelect(key: string) {
  const opt = menuOptions.find((o) => o.key === key)
  if (opt && 'path' in opt) router.push(opt.path as string)
}

function handleUserSelect(key: string) {
  if (key === 'logout') {
    authStore.logout()
    router.replace('/auth/login')
  }
}
</script>

<style scoped>
.app-shell {
  height: 100vh;
}
.sidebar-header {
  padding: 16px;
  font-weight: 600;
  border-bottom: 1px solid var(--n-border-color);
}
.logo {
  font-size: 1.25rem;
}
.header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 0 16px;
  height: 56px;
}
.header-title {
  font-weight: 600;
}
.content {
  padding: 0;
}
</style>
