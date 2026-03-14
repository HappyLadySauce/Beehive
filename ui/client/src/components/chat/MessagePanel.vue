<template>
  <div ref="scrollRef" class="message-panel" @scroll="onScroll">
    <div v-if="loadingMore" class="load-more">
      <n-spin size="small" />
    </div>
    <div class="messages">
      <div
        v-for="m in messages"
        :key="m.serverMsgId || (m as any).clientMsgId"
        :class="['message-row', isSelf(m) ? 'self' : 'other']"
      >
        <n-avatar v-if="!isSelf(m)" round size="small">{{ (m.fromUserId || '?').slice(-1) }}</n-avatar>
        <div class="bubble">
          <div v-if="m.body?.type === 'text'" class="text">{{ m.body.text }}</div>
          <div v-else class="text">[ 未知类型 ]</div>
          <div class="meta">
            <n-time :time="m.serverTime * 1000" type="datetime" format="HH:mm" />
            <span v-if="(m as MessageWithStatus).status === 'pending'" class="status">发送中</span>
            <template v-else-if="(m as MessageWithStatus).status === 'failed'">
              <span class="status failed">发送失败</span>
              <n-button v-if="isSelf(m)" quaternary size="tiny" @click="retrySend(m)">重试</n-button>
            </template>
          </div>
        </div>
        <n-avatar v-if="isSelf(m)" round size="small">{{ (m.fromUserId || '?').slice(-1) }}</n-avatar>
      </div>
    </div>
    <n-empty v-if="messages.length === 0 && !loadingMore" description="暂无消息" style="margin-top: 48px" />
  </div>
</template>

<script setup lang="ts">
import { ref, computed, watch } from 'vue'
import { NAvatar, NTime, NSpin, NEmpty, NButton } from 'naive-ui'
import { useAuthStore } from '@/stores/auth'
import { useConversationStore } from '@/stores/conversation'
import { useMessageStore, type MessageWithStatus } from '@/stores/message'
import { useWebSocket } from '@/composables/useWebSocket'

const props = defineProps<{ conversationId: string }>()
const authStore = useAuthStore()
const conversationStore = useConversationStore()
const messageStore = useMessageStore()
const { send, connected } = useWebSocket()

const scrollRef = ref<HTMLElement | null>(null)
const loadingMore = ref(false)

const messages = computed(() => messageStore.getMessages(props.conversationId))

function isSelf(m: MessageWithStatus) {
  return m.fromUserId === authStore.user?.id
}

function markRead() {
  const list = messageStore.getMessages(props.conversationId)
  const last = list[list.length - 1]
  if (last?.serverMsgId) {
    send('message.read', { conversationId: props.conversationId, serverMsgId: last.serverMsgId })
  }
  conversationStore.clearUnread(props.conversationId)
}

function loadHistory() {
  if (!connected.value || loadingMore.value) return
  const list = messageStore.getMessages(props.conversationId)
  const before = list.length > 0 ? Math.min(...list.map((m) => m.serverTime)) : undefined
  if (before === undefined && list.length > 0) return
  loadingMore.value = true
  send(
    'message.history',
    { conversationId: props.conversationId, before, limit: 50 },
    (env) => {
      loadingMore.value = false
      if (env.type === 'message.history.ok' && env.payload) {
        const p = env.payload as { items: MessageWithStatus[]; hasMore: boolean }
        messageStore.prependHistory(props.conversationId, p.items)
        if (p.items.length) messageStore.setCursor(props.conversationId, p.items[p.items.length - 1].serverTime)
        markRead()
      }
    }
  )
}

function onScroll() {
  const el = scrollRef.value
  if (!el) return
  if (el.scrollTop < 80) loadHistory()
}

function retrySend(m: MessageWithStatus) {
  const clientMsgId = (m as MessageWithStatus).clientMsgId
  if (!clientMsgId || !m.body) return
  messageStore.updateMessageStatus(props.conversationId, clientMsgId, 'pending')
  send(
    'message.send',
    {
      clientMsgId,
      conversationId: props.conversationId,
      body: m.body,
    },
    (env) => {
      if (env.type === 'message.send.ok' && env.payload) {
        const p = env.payload as { serverMsgId: string }
        messageStore.updateMessageStatus(props.conversationId, clientMsgId, 'sent', p.serverMsgId)
      } else if (env.error) {
        messageStore.updateMessageStatus(props.conversationId, clientMsgId, 'failed')
      }
    }
  )
}

watch(
  () => props.conversationId,
  (id) => {
    if (!id) return
    const list = messageStore.getMessages(id)
    if (list.length === 0) loadHistory()
    else markRead()
  },
  { immediate: true }
)
</script>

<style scoped>
.message-panel {
  flex: 1;
  overflow: auto;
  padding: 16px;
  display: flex;
  flex-direction: column;
}
.load-more {
  text-align: center;
  padding: 8px;
}
.messages {
  display: flex;
  flex-direction: column;
  gap: 12px;
}
.message-row {
  display: flex;
  gap: 8px;
  align-items: flex-end;
}
.message-row.self {
  flex-direction: row-reverse;
}
.bubble {
  max-width: 70%;
  padding: 8px 12px;
  border-radius: 12px;
  background: var(--n-color);
}
.self .bubble {
  background: var(--n-primaryColor);
  color: var(--n-primaryColorSuppl);
}
.text {
  word-break: break-word;
}
.meta {
  font-size: 12px;
  opacity: 0.8;
  margin-top: 4px;
}
.status.failed {
  color: var(--n-errorColor);
}
</style>
