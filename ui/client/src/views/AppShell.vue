<template>
  <n-layout has-sider class="app-shell">
    <n-layout-sider v-model:collapsed="siderCollapsed" bordered collapse-mode="width" :collapsed-width="64" :width="240" show-trigger>
      <div class="sidebar-header" :class="{ collapsed: siderCollapsed }">
        <span v-show="!siderCollapsed" class="logo">Beehive</span>
        <n-dropdown :options="userOptions" trigger="click" @select="handleUserSelect">
          <n-button quaternary circle class="sidebar-avatar">
            <n-avatar round size="small">{{ userInitial }}</n-avatar>
          </n-button>
        </n-dropdown>
      </div>
      <div class="sidebar-nav">
        <n-menu
          :options="mainMenuOptions"
          :value="currentMenu"
          @update:value="handleMenuSelect"
        />
        <n-menu
          :options="bottomMenuOptions"
          :value="currentMenu"
          @update:value="handleMenuSelect"
          class="bottom-menu"
        />
      </div>
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
import { computed, ref } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { NLayout, NLayoutSider, NLayoutHeader, NLayoutContent, NMenu, NButton, NAvatar, NDropdown } from 'naive-ui'
import type { MenuOption } from 'naive-ui'
import { useAuthStore } from '@/stores/auth'

const route = useRoute()
const router = useRouter()
const authStore = useAuthStore()

const siderCollapsed = ref(false)

const mainMenuOptions: MenuOption[] = [
  { label: '消息', key: 'chats', path: '/app/chats' },
  { label: '联系人', key: 'contacts', path: '/app/contacts' },
]
const bottomMenuOptions: MenuOption[] = [
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
  return '消息'
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
  const opt = [...mainMenuOptions, ...bottomMenuOptions].find((o) => o.key === key)
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
  display: flex;
  align-items: center;
  gap: 12px;
}
.sidebar-header.collapsed {
  justify-content: center;
  padding: 12px;
}
.sidebar-header .sidebar-avatar {
  flex-shrink: 0;
}
.logo {
  font-size: 1.25rem;
  flex: 1;
}
.sidebar-nav {
  display: flex;
  flex-direction: column;
  flex: 1;
  min-height: 0;
}
.sidebar-nav .n-menu {
  flex: 1;
}
.sidebar-nav .bottom-menu {
  margin-top: auto;
  flex: 0;
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
