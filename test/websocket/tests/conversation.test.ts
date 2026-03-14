/**
 * Conversation 模块 WebSocket 接口测试：conversation.list / create / addMember / removeMember。
 */
import { describe, it, expect, beforeEach, afterEach } from "vitest";
import { createClient } from "../src/client.js";
import { WS_URL } from "../src/config.js";
import type {
  ConversationListOkPayload,
  ConversationCreateOkPayload,
} from "../src/types.js";

describe("conversation", () => {
  const client = createClient(WS_URL);

  beforeEach(async () => {
    await client.connect();
  });

  afterEach(() => {
    client.disconnect();
  });

  it("conversation.list without login returns unauthorized", async () => {
    await expect(
      client.sendAndWait("conversation.list", { cursor: null, limit: 10 })
    ).rejects.toThrow(/unauthorized|user not logged in/i);
  });

  it("conversation.list after login returns conversation.list.ok", async () => {
    const loginEnv = await client.sendAndWait("auth.login", {
      username: process.env.TEST_USER ?? "testuser",
      password: process.env.TEST_PASSWORD ?? "password123",
      deviceId: "test-dev",
    });
    expect(loginEnv.type).toMatch(/auth\.(login|tokenLogin)\.ok/);

    const env = await client.sendAndWait<ConversationListOkPayload>("conversation.list", {
      cursor: null,
      limit: 50,
    });
    expect(env.type).toBe("conversation.list.ok");
    expect(env.error).toBeFalsy();
    const p = env.payload as ConversationListOkPayload;
    expect(p).toBeDefined();
    expect(Array.isArray(p.items)).toBe(true);
    // nextCursor 可为 null 或 string，这里只断言字段存在
    expect("nextCursor" in p).toBe(true);
  });

  it("conversation.create as group with self member returns conversation.create.ok", async () => {
    const loginEnv = await client.sendAndWait("auth.login", {
      username: process.env.TEST_USER ?? "testuser",
      password: process.env.TEST_PASSWORD ?? "password123",
      deviceId: "test-dev",
    });
    const loginPayload = loginEnv.payload as { userId?: string };
    const userId = loginPayload.userId;
    expect(userId).toBeDefined();

    const env = await client.sendAndWait<ConversationCreateOkPayload>("conversation.create", {
      type: "group",
      name: "Test Group",
      memberIds: [userId],
    });
    expect(env.type).toBe("conversation.create.ok");
    expect(env.error).toBeFalsy();
    const p = env.payload as ConversationCreateOkPayload;
    expect(typeof p.conversationId).toBe("string");
    expect(p.conversationId.length).toBeGreaterThan(0);
  });

  it("conversation.addMember with non-existent user returns error", async () => {
    const loginEnv = await client.sendAndWait("auth.login", {
      username: process.env.TEST_USER ?? "testuser",
      password: process.env.TEST_PASSWORD ?? "password123",
      deviceId: "test-dev",
    });
    const loginPayload = loginEnv.payload as { userId?: string };
    const userId = loginPayload.userId!;

    const createEnv = await client.sendAndWait<ConversationCreateOkPayload>("conversation.create", {
      type: "group",
      name: "For AddMember Error",
      memberIds: [userId],
    });
    const convId = (createEnv.payload as ConversationCreateOkPayload).conversationId;

    await expect(
      client.sendAndWait("conversation.addMember", {
        conversationId: convId,
        userId: "00000000-0000-0000-0000-000000000000",
      })
    ).rejects.toThrow(/not_found|bad_request|user already in conversation/i);
  });

  it("conversation.removeMember with non-existent user returns error", async () => {
    const loginEnv = await client.sendAndWait("auth.login", {
      username: process.env.TEST_USER ?? "testuser",
      password: process.env.TEST_PASSWORD ?? "password123",
      deviceId: "test-dev",
    });
    const loginPayload = loginEnv.payload as { userId?: string };
    const userId = loginPayload.userId!;

    const createEnv = await client.sendAndWait<ConversationCreateOkPayload>("conversation.create", {
      type: "group",
      name: "For RemoveMember Error",
      memberIds: [userId],
    });
    const convId = (createEnv.payload as ConversationCreateOkPayload).conversationId;

    await expect(
      client.sendAndWait("conversation.removeMember", {
        conversationId: convId,
        userId: "00000000-0000-0000-0000-000000000000",
      })
    ).rejects.toThrow(/not_found|bad_request/i);
  });
});

