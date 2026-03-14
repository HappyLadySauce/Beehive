<template>
  <div class="settings-page">
    <n-card title="个人设置">
      <n-descriptions :column="1" bordered>
        <n-descriptions-item label="用户 ID">{{ authStore.user?.id ?? '-' }}</n-descriptions-item>
        <n-descriptions-item label="昵称">{{ authStore.user?.nickname ?? '-' }}</n-descriptions-item>
      </n-descriptions>
      <n-divider />
      <n-space>
        <span>主题</span>
        <n-radio-group :value="themeStore.theme" @update:value="(v: 'light'|'dark'|null) => themeStore.setTheme(v)">
          <n-radio value="light">浅色</n-radio>
          <n-radio value="dark">深色</n-radio>
        </n-radio-group>
      </n-space>
    </n-card>
  </div>
</template>

<script setup lang="ts">
import { onMounted } from 'vue'
import { NCard, NDescriptions, NDescriptionsItem, NDivider, NSpace, NRadioGroup, NRadio } from 'naive-ui'
import { useAuthStore } from '@/stores/auth'
import { useThemeStore } from '@/stores/theme'
import { useWebSocket } from '@/composables/useWebSocket'

const authStore = useAuthStore()
const themeStore = useThemeStore()
const { send, connected } = useWebSocket()

onMounted(() => {
  if (connected.value) {
    send('user.me', {}, (env) => {
      if (env.type === 'user.me.ok' && env.payload) {
        const u = env.payload as { id: string; nickname: string; avatarUrl: string; bio: string }
        authStore.setUser({ id: u.id, nickname: u.nickname, avatarUrl: u.avatarUrl })
      }
    })
  }
})
</script>

<style scoped>
.settings-page {
  padding: 16px;
}
</style>
