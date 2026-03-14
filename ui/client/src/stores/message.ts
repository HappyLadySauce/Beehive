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
    // 去重：服务端消息按 serverMsgId，本地 pending 消息按 clientMsgId（pending 时 serverMsgId 为空，多条会误判为重复）
    if (message.serverMsgId && list.some((m) => m.serverMsgId === message.serverMsgId)) return
    if (message.clientMsgId && list.some((m) => (m as MessageWithStatus).clientMsgId === message.clientMsgId)) return
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
