<template>
  <div class="conversation-list">
    <n-list hoverable clickable>
      <n-list-item
        v-for="c in conversationStore.conversations"
        :key="c.id"
        @click="select(c.id)"
      >
        <n-thing>
          <template #avatar>
            <n-badge :value="c.unreadCount > 0 ? c.unreadCount : undefined" :max="99">
              <n-avatar round size="medium">{{ (c.name || c.id).slice(0, 1) }}</n-avatar>
            </n-badge>
          </template>
          <template #header>
            {{ c.type === 'group' && c.memberCount != null ? `${c.name || c.id}(${c.memberCount})` : (c.name || c.id) }}
          </template>
          <template #header-extra>
            <n-time :time="c.lastActiveAt * 1000" type="relative" />
          </template>
          <template #default>
            {{ c.lastMessage?.preview ?? '暂无消息' }}
          </template>
        </n-thing>
      </n-list-item>
    </n-list>
    <n-empty v-if="conversationStore.conversations.length === 0" description="暂无会话" style="margin-top: 24px" />
  </div>
</template>

<script setup lang="ts">
import { onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { NList, NListItem, NThing, NAvatar, NBadge, NTime, NEmpty } from 'naive-ui'
import { useConversationStore } from '@/stores/conversation'
import { useWebSocket } from '@/composables/useWebSocket'

const router = useRouter()
const conversationStore = useConversationStore()
const { send, connected } = useWebSocket()

function select(id: string) {
  conversationStore.setCurrentConversation(id)
  router.push(`/app/chats/${id}`)
}

onMounted(() => {
  if (connected.value) {
    send('conversation.list', { limit: 50 }, (env) => {
      if (env.type === 'conversation.list.ok' && env.payload) {
        const p = env.payload as { items: typeof conversationStore.conversations; nextCursor: string | null }
        conversationStore.setConversations(p.items)
      }
    })
  }
})
</script>

<style scoped>
.conversation-list {
  padding: 8px;
  height: 100%;
  overflow: auto;
}
</style>
