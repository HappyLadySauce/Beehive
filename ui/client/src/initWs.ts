/**
 * 初始化 WebSocket：若已有 token 则连接并 tokenLogin，并注册全局消息派发（message.push 等）。
 * 在 main.ts 中 app.mount 之后调用。
 */
import { useAuthStore } from '@/stores/auth'
import { useConversationStore } from '@/stores/conversation'
import { useMessageStore } from '@/stores/message'
import type { WsEnvelope } from '@/types/ws'
import type { MessagePushPayload } from '@/types/ws'

// Use module-level refs to the composable's internal state so we can register once.
// We need to access the singleton send/on - import from a module that exposes them.
import { useWebSocket, setWsOnOpen } from '@/composables/useWebSocket'

export function initWs() {
  const authStore = useAuthStore()
  const conversationStore = useConversationStore()
  const messageStore = useMessageStore()
  const { ensureConnected, send, on } = useWebSocket()

  function doTokenLogin() {
    if (!authStore.accessToken) return
    send('auth.tokenLogin', {
      accessToken: authStore.accessToken,
      deviceId: `web-${typeof navigator !== 'undefined' ? navigator.userAgent.slice(0, 32) : 'unknown'}`,
    }, (env) => {
      if (env.type !== 'auth.tokenLogin.ok' || !env.payload) return
      const p = env.payload as { userId: string }
      authStore.setFromTokenLogin({ userId: p.userId, accessToken: authStore.accessToken!, refreshToken: authStore.refreshToken ?? undefined })
      send('user.me', {}, (meEnv) => {
        if (meEnv.type === 'user.me.ok' && meEnv.payload) {
          const u = meEnv.payload as { id: string; nickname: string; avatarUrl: string }
          authStore.setUser({ id: u.id, nickname: u.nickname, avatarUrl: u.avatarUrl })
        }
      })
    })
  }

  setWsOnOpen(() => doTokenLogin())

  on('message.push', (env: WsEnvelope) => {
    const payload = env.payload as MessagePushPayload | undefined
    if (!payload) return
    const convId = payload.conversationId
    messageStore.appendMessage(convId, {
      serverMsgId: payload.serverMsgId,
      clientMsgId: payload.clientMsgId,
      conversationId: payload.conversationId,
      fromUserId: payload.fromUserId,
      toUserId: payload.toUserId,
      body: payload.body,
      serverTime: payload.serverTime,
    })
    conversationStore.updateLastMessage(convId, {
      serverMsgId: payload.serverMsgId,
      preview: payload.body?.type === 'text' ? payload.body.text?.slice(0, 50) ?? '' : '[消息]',
      serverTime: payload.serverTime,
    })
    const current = conversationStore.currentConversationId
    if (convId !== current) {
      conversationStore.incrementUnread(convId)
    }
  })

  on('auth.login.ok', (env: WsEnvelope) => {
    const p = env.payload as { userId: string; accessToken?: string; refreshToken?: string } | undefined
    if (p) authStore.setFromTokenLogin({ userId: p.userId, accessToken: p.accessToken, refreshToken: p.refreshToken })
  })

  on('auth.register.ok', (env: WsEnvelope) => {
    const p = env.payload as { userId: string; accessToken?: string; refreshToken?: string } | undefined
    if (p) authStore.setFromTokenLogin({ userId: p.userId, accessToken: p.accessToken, refreshToken: p.refreshToken })
  })

  on('auth.tokenLogin.ok', (env: WsEnvelope) => {
    const p = env.payload as { userId: string; accessToken?: string; refreshToken?: string } | undefined
    if (p) authStore.setFromTokenLogin({ userId: p.userId, accessToken: p.accessToken, refreshToken: p.refreshToken })
  })

  if (authStore.accessToken) {
    ensureConnected().then(() => doTokenLogin()).catch(() => {})
  }
}
