/**
 * Message 模块 WebSocket 接口测试：message.send / history / read。
 */
import { describe, it, expect, beforeEach, afterEach } from "vitest";
import { createClient } from "../src/client.js";
import { WS_URL } from "../src/config.js";
import type {
  ConversationCreateOkPayload,
  MessageSendOkPayload,
  MessageHistoryOkPayload,
} from "../src/types.js";

describe("message", () => {
  const client = createClient(WS_URL);

  beforeEach(async () => {
    await client.connect();
  });

  afterEach(() => {
    client.disconnect();
  });

  it("message.send without login returns unauthorized", async () => {
    await expect(
      client.sendAndWait("message.send", {
        clientMsgId: "noauth-1",
        conversationId: "conv_x",
        body: { type: "text", text: "hello" },
      })
    ).rejects.toThrow(/unauthorized|user not logged in/i);
  });

  it("message.send / history / read after login work as expected", async () => {
    const loginEnv = await client.sendAndWait("auth.login", {
      username: process.env.TEST_USER ?? "testuser",
      password: process.env.TEST_PASSWORD ?? "password123",
      deviceId: "test-dev",
    });
    const loginPayload = loginEnv.payload as { userId?: string };
    const userId = loginPayload.userId!;

    // 1) 创建一个 group 会话，以便发送消息和拉历史
    const createEnv = await client.sendAndWait<ConversationCreateOkPayload>("conversation.create", {
      type: "group",
      name: "Message Test Group",
      memberIds: [userId],
    });
    const convId = (createEnv.payload as ConversationCreateOkPayload).conversationId;

    // 2) 发送一条文本消息
    const sendEnv = await client.sendAndWait<MessageSendOkPayload>("message.send", {
      clientMsgId: `test-${Date.now()}`,
      conversationId: convId,
      body: { type: "text", text: "hello from test" },
    });
    expect(sendEnv.type).toBe("message.send.ok");
    expect(sendEnv.error).toBeFalsy();
    const sendPayload = sendEnv.payload as MessageSendOkPayload;
    expect(typeof sendPayload.serverMsgId).toBe("string");
    expect(sendPayload.serverMsgId.length).toBeGreaterThan(0);
    expect(typeof sendPayload.serverTime).toBe("number");
    expect(sendPayload.conversationId).toBe(convId);

    // 3) 拉取历史消息，应至少包含刚刚发送的那条
    const historyEnv = await client.sendAndWait<MessageHistoryOkPayload>("message.history", {
      conversationId: convId,
      limit: 50,
    });
    expect(historyEnv.type).toBe("message.history.ok");
    expect(historyEnv.error).toBeFalsy();
    const historyPayload = historyEnv.payload as MessageHistoryOkPayload;
    expect(Array.isArray(historyPayload.items)).toBe(true);
    expect(typeof historyPayload.hasMore).toBe("boolean");
    const found = historyPayload.items.some((m) => m.serverMsgId === sendPayload.serverMsgId);
    expect(found).toBe(true);

    // 4) 已读回执 message.read
    const readEnv = await client.sendAndWait("message.read", {
      conversationId: convId,
      serverMsgId: sendPayload.serverMsgId,
    });
    expect(readEnv.type).toBe("message.read.ok");
    expect(readEnv.error).toBeFalsy();
  });
});

