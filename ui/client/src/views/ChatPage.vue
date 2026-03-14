<template>
  <div class="chat-page">
    <n-layout has-sider class="chat-layout">
      <n-layout-sider bordered :width="280" class="conversation-sider">
        <ConversationList />
      </n-layout-sider>
      <n-layout>
        <template v-if="currentConversationId">
          <div class="chat-header">
            <span class="chat-title">{{ chatTitle }}</span>
            <n-dropdown v-if="isGroupChat" :options="groupMenuOptions" trigger="click" @select="handleGroupMenuSelect">
              <n-button quaternary>更多</n-button>
            </n-dropdown>
          </div>
          <div class="chat-main">
            <MessagePanel :conversation-id="currentConversationId" />
            <MessageInput :conversation-id="currentConversationId" />
          </div>
        </template>
        <template v-else>
          <div class="empty-chat">
            <n-empty description="选择或开始一段会话" />
          </div>
        </template>
      </n-layout>
      <n-layout-sider v-if="currentConversationId && isGroupChat" bordered :width="260" class="group-sider">
        <GroupSidebar :conversation-id="currentConversationId" />
      </n-layout-sider>
    </n-layout>
  </div>
</template>

<script setup lang="ts">
import { computed, watch } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { NLayout, NLayoutSider, NEmpty, NButton, NDropdown, useDialog } from 'naive-ui'
import ConversationList from '@/components/chat/ConversationList.vue'
import MessagePanel from '@/components/chat/MessagePanel.vue'
import MessageInput from '@/components/chat/MessageInput.vue'
import GroupSidebar from '@/components/chat/GroupSidebar.vue'
import { useConversationStore } from '@/stores/conversation'
import { useAuthStore } from '@/stores/auth'
import { useWebSocket } from '@/composables/useWebSocket'
import { useMessage } from 'naive-ui'

const route = useRoute()
const router = useRouter()
const message = useMessage()
const authStore = useAuthStore()
const conversationStore = useConversationStore()
const { send } = useWebSocket()

const currentConversationId = computed(() => {
  const id = route.params.conversationId as string
  if (id) return id
  return conversationStore.currentConversationId
})

watch(
  () => route.params.conversationId,
  (id) => {
    conversationStore.setCurrentConversation((id as string) ?? null)
  },
  { immediate: true }
)

const currentConv = computed(() => conversationStore.currentConversation)

const isGroupChat = computed(() => currentConv.value?.type === 'group')

const chatTitle = computed(() => {
  const c = currentConv.value
  if (!c) return currentConversationId.value || ''
  if (c.type === 'group' && c.memberCount != null) return `${c.name || c.id}(${c.memberCount})`
  return c.name || c.id
})

const dialog = useDialog()
const groupMenuOptions = [{ label: '退出群聊', key: 'leaveGroup' }]

function handleGroupMenuSelect(key: string) {
  if (key === 'leaveGroup') {
    dialog.warning({
      title: '退出群聊',
      content: '确定要退出该群聊吗？',
      positiveText: '确定',
      negativeText: '取消',
      onPositiveClick: () => leaveGroup(),
    })
  }
}

function leaveGroup() {
  const id = currentConversationId.value
  if (!id || !authStore.user) return
  send('conversation.removeMember', { conversationId: id, userId: authStore.user.id }, (env) => {
    if (env.type === 'conversation.removeMember.ok') {
      conversationStore.setConversations(conversationStore.conversations.filter((c) => c.id !== id))
      conversationStore.setCurrentConversation(null)
      router.push('/app/chats')
      message.success('已退出群聊')
    } else if (env.error) {
      message.error(env.error.message || '退出失败')
    }
  })
}
</script>

<style scoped>
.chat-page,
.chat-layout {
  height: 100%;
}
.conversation-sider {
  height: 100%;
}
.chat-header {
  padding: 12px 16px;
  border-bottom: 1px solid var(--n-border-color);
  font-weight: 600;
  display: flex;
  align-items: center;
  justify-content: space-between;
}
.chat-title {
  font-size: 1rem;
}
.chat-main {
  flex: 1;
  display: flex;
  flex-direction: column;
  min-height: 0;
}
.group-sider {
  height: 100%;
}
.empty-chat {
  height: 100%;
  display: flex;
  align-items: center;
  justify-content: center;
}
</style>
