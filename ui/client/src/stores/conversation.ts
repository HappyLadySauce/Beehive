import { defineStore } from 'pinia'
import { ref, computed } from 'vue'
import type { ConversationItem } from '@/types/ws'

export const useConversationStore = defineStore('conversation', () => {
  const conversations = ref<ConversationItem[]>([])
  const currentConversationId = ref<string | null>(null)

  const currentConversation = computed(() => {
    const id = currentConversationId.value
    if (!id) return null
    return conversations.value.find((c) => c.id === id) ?? null
  })

  function setConversations(list: ConversationItem[]) {
    conversations.value = list
  }

  function setCurrentConversation(id: string | null) {
    currentConversationId.value = id
  }

  function updateLastMessage(conversationId: string, lastMessage: ConversationItem['lastMessage']) {
    const c = conversations.value.find((x) => x.id === conversationId)
    if (c) c.lastMessage = lastMessage
  }

  function incrementUnread(conversationId: string) {
    const c = conversations.value.find((x) => x.id === conversationId)
    if (c) c.unreadCount = (c.unreadCount || 0) + 1
  }

  function clearUnread(conversationId: string) {
    const c = conversations.value.find((x) => x.id === conversationId)
    if (c) c.unreadCount = 0
  }

  return {
    conversations,
    currentConversationId,
    currentConversation,
    setConversations,
    setCurrentConversation,
    updateLastMessage,
    incrementUnread,
    clearUnread,
  }
})
