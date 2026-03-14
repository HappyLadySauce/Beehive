/**
 * WebSocket 客户端封装：按 tid 配对请求-响应，支持超时。
 */
import WebSocket from "ws";
import type { Envelope } from "./types.js";

const DEFAULT_TIMEOUT_MS = 5000;

type Pending = {
  requestType: string;
  resolve: (env: Envelope) => void;
  reject: (err: Error) => void;
  timer: ReturnType<typeof setTimeout>;
};

export function createClient(url: string, timeoutMs = DEFAULT_TIMEOUT_MS) {
  let ws: WebSocket | null = null;
  const pendingByTid = new Map<string, Pending>();

  function connect(): Promise<void> {
    return new Promise((resolve, reject) => {
      ws = new WebSocket(url);
      ws.on("open", () => resolve());
      ws.on("error", (err) => reject(err));
      ws.on("message", (data: Buffer) => {
        try {
          const env = JSON.parse(data.toString()) as Envelope;
          const tid = env.tid ?? "";
          const p = pendingByTid.get(tid);
          if (!p) return;
          const isOk = env.type === `${p.requestType}.ok`;
          const isError = !!env.error;
          if (isOk || isError) {
            pendingByTid.delete(tid);
            clearTimeout(p.timer);
            if (env.error) {
              p.reject(new Error(`${env.error.code}: ${env.error.message}`));
            } else {
              p.resolve(env);
            }
          }
        } catch {
          // ignore non-JSON or unknown messages
        }
      });
    });
  }

  function disconnect(): void {
    for (const [, p] of pendingByTid) {
      clearTimeout(p.timer);
      p.reject(new Error("disconnected"));
    }
    pendingByTid.clear();
    if (ws) {
      ws.close();
      ws = null;
    }
  }

  function sendAndWait<T = unknown>(
    type: string,
    payload: object,
    tid?: string
  ): Promise<Envelope<T>> {
    return new Promise((resolve, reject) => {
      if (!ws || ws.readyState !== WebSocket.OPEN) {
        reject(new Error("not connected"));
        return;
      }
      const id = tid ?? `req-${Date.now()}-${Math.random().toString(36).slice(2)}`;
      const timer = setTimeout(() => {
        if (pendingByTid.delete(id)) {
          reject(new Error(`timeout waiting for response (tid=${id})`));
        }
      }, timeoutMs);
      pendingByTid.set(id, {
        requestType: type,
        resolve: resolve as (env: Envelope) => void,
        reject,
        timer,
      });
      ws.send(
        JSON.stringify({
          type,
          tid: id,
          payload,
        })
      );
    });
  }

  return { connect, disconnect, sendAndWait };
}

export type WsClient = ReturnType<typeof createClient>;
