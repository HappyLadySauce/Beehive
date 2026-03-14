<template>
  <n-config-provider :theme="theme" :theme-overrides="themeOverrides">
    <n-dialog-provider>
      <n-message-provider>
        <AppContent />
      </n-message-provider>
    </n-dialog-provider>
  </n-config-provider>
</template>

<script setup lang="ts">
import { NConfigProvider, NMessageProvider, NDialogProvider, darkTheme, lightTheme, type GlobalThemeOverrides } from 'naive-ui'
import { computed } from 'vue'
import { useThemeStore } from '@/stores/theme'
import AppContent from './AppContent.vue'

const themeStore = useThemeStore()
const theme = computed(() => (themeStore.theme === 'dark' ? darkTheme : themeStore.theme === 'light' ? lightTheme : null))
const themeOverrides: GlobalThemeOverrides = {}
</script>
