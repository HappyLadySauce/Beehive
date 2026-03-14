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
