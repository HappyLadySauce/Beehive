<template>
  <div class="chat-page">
    <n-layout has-sider class="chat-layout">
      <n-layout-sider bordered :width="280" class="conversation-sider">
        <ConversationList />
      </n-layout-sider>
      <n-layout>
        <template v-if="currentConversationId">
          <MessagePanel :conversation-id="currentConversationId" />
          <MessageInput :conversation-id="currentConversationId" />
        </template>
        <template v-else>
          <div class="empty-chat">
            <n-empty description="选择或开始一段会话" />
          </div>
        </template>
      </n-layout>
    </n-layout>
  </div>
</template>

<script setup lang="ts">
import { computed } from 'vue'
import { useRoute } from 'vue-router'
import { NLayout, NLayoutSider, NEmpty } from 'naive-ui'
import ConversationList from '@/components/chat/ConversationList.vue'
import MessagePanel from '@/components/chat/MessagePanel.vue'
import MessageInput from '@/components/chat/MessageInput.vue'
import { useConversationStore } from '@/stores/conversation'

const route = useRoute()
const conversationStore = useConversationStore()

const currentConversationId = computed(() => {
  const id = route.params.conversationId as string
  if (id) return id
  return conversationStore.currentConversationId
})
</script>

<style scoped>
.chat-page,
.chat-layout {
  height: 100%;
}
.conversation-sider {
  height: 100%;
}
.empty-chat {
  height: 100%;
  display: flex;
  align-items: center;
  justify-content: center;
}
</style>
