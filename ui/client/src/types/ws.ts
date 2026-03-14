/**
 * WebSocket 协议类型，与 docs/API/websocket-client-api.md 对齐。
 * 所有 JSON 使用 camelCase。
 */

/** 错误体，响应/推送中的 error 字段 */
export interface WsErrorBody {
  code: string
  message: string
}

/** 统一 Envelope：type / tid / payload / error */
export interface WsEnvelope<T = unknown> {
  type: string
  tid?: string
  payload?: T
  error?: WsErrorBody | null
}

/** 常见错误码 */
export type WsErrorCode =
  | 'bad_request'
  | 'unauthorized'
  | 'forbidden'
  | 'rate_limited'
  | 'not_found'
  | 'internal_error'
  | 'unavailable'

// --- Auth ---
export interface AuthLoginPayload {
  username: string
  password: string
  deviceId?: string
}

export interface AuthTokenLoginPayload {
  accessToken: string
  deviceId?: string
}

export interface AuthLoginOkPayload {
  userId: string
  accessToken?: string
  refreshToken?: string
  expiresIn?: number
}

// --- Presence ---
export interface PresencePingPayload {
  clientTime?: number
}

export interface PresencePingOkPayload {
  serverTime: number
}

// --- User ---
export interface UserMeOkPayload {
  id: string
  nickname: string
  avatarUrl: string
  bio: string
  status: string
}

// --- Message ---
export interface MessageBody {
  type: 'text' | 'image' | 'system'
  text?: string
  [key: string]: unknown
}

export interface MessageSendPayload {
  clientMsgId: string
  conversationId?: string
  toUserId?: string
  body: MessageBody
}

export interface MessageSendOkPayload {
  serverMsgId: string
  serverTime: number
  conversationId: string
}

export interface MessagePushPayload {
  serverMsgId: string
  clientMsgId?: string
  conversationId: string
  fromUserId: string
  toUserId?: string
  body: MessageBody
  serverTime: number
}

export interface MessageHistoryPayload {
  conversationId: string
  before?: number
  limit?: number
}

export interface MessageHistoryItem {
  serverMsgId: string
  clientMsgId?: string | null
  conversationId: string
  fromUserId: string
  toUserId?: string
  body: MessageBody
  serverTime: number
}

export interface MessageHistoryOkPayload {
  items: MessageHistoryItem[]
  hasMore: boolean
}

export interface MessageReadPayload {
  conversationId: string
  serverMsgId: string
}

// --- Conversation ---
export interface ConversationItem {
  id: string
  name: string
  avatar: string
  type: 'single' | 'group' | 'channel'
  memberCount?: number
  lastMessage?: {
    serverMsgId: string
    preview: string
    serverTime: number
  }
  unreadCount: number
  lastActiveAt: number
}

export interface ConversationListPayload {
  cursor?: string | null
  limit?: number
}

export interface ConversationListOkPayload {
  items: ConversationItem[]
  nextCursor: string | null
}

export interface ConversationCreatePayload {
  type: 'single' | 'group' | 'channel'
  name?: string
  memberIds: string[]
}

export interface ConversationCreateOkPayload {
  conversationId: string
}

export interface ConversationAddMemberPayload {
  conversationId: string
  userId: string
  role?: string
}

export interface ConversationRemoveMemberPayload {
  conversationId: string
  userId: string
}
