import { defineStore } from 'pinia'
import { ref } from 'vue'
import type { MessageHistoryItem } from '@/types/ws'

export type MessageStatus = 'pending' | 'sent' | 'failed'

export interface MessageWithStatus extends MessageHistoryItem {
  status?: MessageStatus
  clientMsgId?: string
}

export const useMessageStore = defineStore('message', () => {
  const messagesByConversation = ref<Record<string, MessageWithStatus[]>>({})
  const cursorByConversation = ref<Record<string, number>>({})

  function getMessages(conversationId: string): MessageWithStatus[] {
    return messagesByConversation.value[conversationId] ?? []
  }

  function appendMessage(conversationId: string, message: MessageWithStatus) {
    if (!messagesByConversation.value[conversationId]) {
      messagesByConversation.value[conversationId] = []
    }
    const list = messagesByConversation.value[conversationId]
    if (list.some((m) => m.serverMsgId === message.serverMsgId)) return
    list.push(message)
  }

  function prependHistory(conversationId: string, messages: MessageWithStatus[]) {
    const list = messagesByConversation.value[conversationId] ?? []
    const existing = new Set(list.map((m) => m.serverMsgId))
    const toPrepend = messages.filter((m) => !existing.has(m.serverMsgId))
    messagesByConversation.value[conversationId] = [...toPrepend, ...list]
  }

  function updateMessageStatus(
    conversationId: string,
    clientMsgId: string,
    status: MessageStatus,
    serverMsgId?: string
  ) {
    const list = messagesByConversation.value[conversationId]
    if (!list) return
    const m = list.find((x) => (x as MessageWithStatus).clientMsgId === clientMsgId)
    if (m) {
      (m as MessageWithStatus).status = status
      if (serverMsgId) m.serverMsgId = serverMsgId
    }
  }

  function setCursor(conversationId: string, before: number) {
    cursorByConversation.value[conversationId] = before
  }

  function getCursor(conversationId: string): number | undefined {
    return cursorByConversation.value[conversationId]
  }

  return {
    messagesByConversation,
    getMessages,
    appendMessage,
    prependHistory,
    updateMessageStatus,
    setCursor,
    getCursor,
  }
})
