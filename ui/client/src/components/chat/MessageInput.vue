<template>
  <div class="message-input">
    <n-input
      v-model:value="text"
      type="textarea"
      placeholder="输入消息，Enter 发送，Shift+Enter 换行"
      :autosize="{ minRows: 1, maxRows: 4 }"
      :disabled="!connected"
      @keydown="onKeydown"
    />
    <n-button type="primary" :disabled="!text.trim() || !connected" :loading="sending" @click="sendMessage">
      发送
    </n-button>
  </div>
</template>

<script setup lang="ts">
import { ref, watch } from 'vue'
import { NInput, NButton } from 'naive-ui'
import { useMessage } from 'naive-ui'
import { useAuthStore } from '@/stores/auth'
import { useMessageStore } from '@/stores/message'
import { useWebSocket } from '@/composables/useWebSocket'
const props = defineProps<{ conversationId: string }>()
const toast = useMessage()
const authStore = useAuthStore()
const messageStore = useMessageStore()
const { send, connected } = useWebSocket()

const text = ref('')
const sending = ref(false)

function genClientMsgId() {
  return `msg-${Date.now()}-${Math.random().toString(36).slice(2, 9)}`
}

function sendMessage() {
  const t = text.value.trim()
  if (!t || !connected.value) return
  const clientMsgId = genClientMsgId()
  const payload = {
    clientMsgId,
    conversationId: props.conversationId,
    body: { type: 'text' as const, text: t },
  }
  messageStore.appendMessage(props.conversationId, {
    serverMsgId: '',
    clientMsgId,
    conversationId: props.conversationId,
    fromUserId: authStore.user!.id,
    body: { type: 'text', text: t },
    serverTime: Math.floor(Date.now() / 1000),
    status: 'pending',
  })
  text.value = ''
  sending.value = true
  send('message.send', payload, (env) => {
    sending.value = false
    if (env.type === 'message.send.ok' && env.payload) {
      const p = env.payload as { serverMsgId: string; conversationId: string }
      messageStore.updateMessageStatus(props.conversationId, clientMsgId, 'sent', p.serverMsgId)
    } else if (env.error) {
      messageStore.updateMessageStatus(props.conversationId, clientMsgId, 'failed')
      if (env.error.code === 'rate_limited') {
        toast.warning('发送过于频繁，请稍后再试')
      } else {
        toast.error(env.error.message || '发送失败')
      }
    }
  })
}

function onKeydown(e: KeyboardEvent) {
  if (e.key === 'Enter' && !e.shiftKey) {
    e.preventDefault()
    sendMessage()
  }
}

watch(
  () => props.conversationId,
  () => {
    text.value = ''
  }
)
</script>

<style scoped>
.message-input {
  display: flex;
  gap: 8px;
  padding: 12px 16px;
  border-top: 1px solid var(--n-border-color);
  align-items: flex-end;
}
.message-input .n-input {
  flex: 1;
}
</style>
