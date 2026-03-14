/**
 * WebSocket Envelope 与各 type 的 payload 类型，与 docs/API/websocket-client-api.md 对齐。
 */

export interface ErrBody {
  code: string;
  message: string;
}

export interface Envelope<T = unknown> {
  type: string;
  tid?: string;
  payload?: T;
  error?: ErrBody | null;
}

// --- Auth ---
export interface AuthLoginPayload {
  username: string;
  password: string;
  deviceId?: string;
}

export interface AuthTokenLoginPayload {
  accessToken: string;
  deviceId?: string;
}

export interface AuthLoginOkPayload {
  userId: string;
  accessToken: string;
  refreshToken: string;
  expiresIn: number;
}

// --- Presence ---
export interface PresencePingPayload {
  clientTime?: number;
}

export interface PresencePingOkPayload {
  serverTime: number;
}

// --- User ---
export interface UserMeOkPayload {
  id: string;
  nickname: string;
  avatarUrl: string;
  bio: string;
  status: string;
}

// --- Conversation ---
export interface ConversationListPayload {
  cursor?: string | null;
  limit?: number;
}

export interface ConversationListItem {
  id: string;
  name: string;
  avatar: string;
  type: "single" | "group" | "channel";
  unreadCount: number;
  lastActiveAt: number;
  lastMessage?: {
    serverMsgId: string;
    preview: string;
    serverTime: number;
  };
}

export interface ConversationListOkPayload {
  items: ConversationListItem[];
  nextCursor: string | null;
}

export interface ConversationCreatePayload {
  type: "single" | "group" | "channel";
  name?: string;
  memberIds: string[];
}

export interface ConversationCreateOkPayload {
  conversationId: string;
}

export interface ConversationAddMemberPayload {
  conversationId: string;
  userId: string;
  role?: string;
}

export interface ConversationRemoveMemberPayload {
  conversationId: string;
  userId: string;
}

// --- Message ---
export interface MessageBody {
  type: string;
  text?: string;
}

export interface MessageSendPayload {
  clientMsgId: string;
  conversationId?: string;
  toUserId?: string;
  body: MessageBody;
}

export interface MessageSendOkPayload {
  serverMsgId: string;
  serverTime: number;
  conversationId: string;
}

export interface MessageHistoryPayload {
  conversationId: string;
  before?: number;
  limit?: number;
}

export interface MessageHistoryItem {
  serverMsgId: string;
  clientMsgId: string | null;
  conversationId: string;
  fromUserId: string;
  toUserId?: string | null;
  body: MessageBody;
  serverTime: number;
}

export interface MessageHistoryOkPayload {
  items: MessageHistoryItem[];
  hasMore: boolean;
}

export interface MessageReadPayload {
  conversationId: string;
  serverMsgId: string;
}

export interface MessageReadOkPayload {}

