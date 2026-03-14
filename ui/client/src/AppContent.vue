<template>
  <router-view />
</template>

<script setup lang="ts">
import { onMounted } from 'vue'
import { useMessage } from 'naive-ui'
import { useWebSocket } from '@/composables/useWebSocket'
import type { WsEnvelope } from '@/types/ws'

const message = useMessage()
const { on } = useWebSocket()

const errorTypes = [
  'auth.login.error',
  'auth.register.error',
  'auth.tokenLogin.error',
  'message.send.error',
  'conversation.create.error',
  'conversation.list.error',
  'message.history.error',
]

onMounted(() => {
  errorTypes.forEach((type) => {
    on(type, (env: WsEnvelope) => {
      if (env.error?.code === 'rate_limited') {
        message.warning('操作过于频繁，请稍后再试')
      } else if (env.error?.message) {
        message.error(env.error.message)
      }
    })
  })
})
</script>
