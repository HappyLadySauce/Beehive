import { ref, onUnmounted, readonly } from 'vue'
import type { WsEnvelope } from '@/types/ws'

type MessageHandler = (env: WsEnvelope) => void

function getWsBaseUrl(): string {
  const base = import.meta.env.VITE_WS_URL
  if (base) return base
  const { protocol, host } = window.location
  const wsProtocol = protocol === 'https:' ? 'wss:' : 'ws:'
  return `${wsProtocol}//${host}`
}

function buildWsUrl(): string {
  const base = getWsBaseUrl()
  const path = base.endsWith('/ws') ? '' : '/ws'
  return base.replace(/\/$/, '') + path
}

let tidCounter = 0
function nextTid(): string {
  tidCounter += 1
  return `c-${Date.now()}-${tidCounter}`
}

// Singleton state
const connected = ref(false)
let ws: WebSocket | null = null
const handlers = new Map<string, MessageHandler>()
const pendingByTid = new Map<string, MessageHandler>()
let reconnectTimer: ReturnType<typeof setTimeout> | null = null
let heartbeatTimer: ReturnType<typeof setInterval> | null = null
let reconnectAttempts = 0
const maxReconnectDelay = 30000
const heartbeatInterval = 30000
let onOpenCallback: (() => void) | null = null

export function setWsOnOpen(cb: () => void) {
  onOpenCallback = cb
}

function clearTimers() {
  if (reconnectTimer) {
    clearTimeout(reconnectTimer)
    reconnectTimer = null
  }
  if (heartbeatTimer) {
    clearInterval(heartbeatTimer)
  }
  heartbeatTimer = null
}

function startHeartbeat() {
  if (heartbeatTimer) return
  heartbeatTimer = setInterval(() => {
    if (ws?.readyState === WebSocket.OPEN) {
      try {
        ws.send(JSON.stringify({
          type: 'presence.ping',
          tid: nextTid(),
          payload: { clientTime: Math.floor(Date.now() / 1000) },
        }))
      } catch {
        // ignore
      }
    }
  }, heartbeatInterval)
}

function handleMessage(event: MessageEvent) {
  try {
    const env = JSON.parse(event.data as string) as WsEnvelope
    const tid = env.tid
    if (tid && pendingByTid.has(tid)) {
      const cb = pendingByTid.get(tid)!
      pendingByTid.delete(tid)
      cb(env)
    }
    const type = env.type
    const handler = handlers.get(type)
    if (handler) handler(env)
    const generic = handlers.get('*')
    if (generic) generic(env)
  } catch {
    // ignore
  }
}

function connect(): Promise<void> {
  return new Promise((resolve, reject) => {
    if (ws?.readyState === WebSocket.OPEN) {
      resolve()
      return
    }
    const url = buildWsUrl()
    const s = new WebSocket(url)
    ws = s

    s.onopen = () => {
      connected.value = true
      reconnectAttempts = 0
      startHeartbeat()
      onOpenCallback?.()
      resolve()
    }

    s.onclose = () => {
      connected.value = false
      ws = null
      clearTimers()
      const delay = Math.min(1000 * 2 ** reconnectAttempts, maxReconnectDelay)
      reconnectAttempts += 1
      reconnectTimer = setTimeout(() => {
        connect().catch(() => {})
      }, delay)
    }

    s.onerror = () => {
      reject(new Error('WebSocket error'))
    }

    s.onmessage = handleMessage
  })
}

function ensureConnected(): Promise<void> {
  if (ws?.readyState === WebSocket.OPEN) return Promise.resolve()
  return connect()
}

function send<T>(type: string, payload?: T, onResponse?: MessageHandler): string {
  const tid = nextTid()
  const env: WsEnvelope<T> = { type, tid, payload }
  if (onResponse) pendingByTid.set(tid, onResponse)
  if (ws?.readyState === WebSocket.OPEN) {
    ws.send(JSON.stringify(env))
  }
  return tid
}

function on(type: string, handler: MessageHandler) {
  handlers.set(type, handler)
}

function off(type: string) {
  handlers.delete(type)
}

export function useWebSocket() {
  const unregister: (() => void)[] = []

  onUnmounted(() => {
    unregister.forEach((f) => f())
  })

  function registerOn(type: string, handler: MessageHandler) {
    on(type, handler)
    unregister.push(() => off(type))
  }

  return {
    connected: readonly(connected),
    connect,
    ensureConnected,
    send,
    on: registerOn,
    off,
  }
}
