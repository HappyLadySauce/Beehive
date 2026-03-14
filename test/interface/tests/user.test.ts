/**
 * User 模块 WebSocket 接口测试：user.me（未登录返回 unauthorized，已登录返回资料）。
 */
import { describe, it, expect, beforeEach, afterEach } from "vitest";
import { createClient } from "../src/client.js";
import { WS_URL } from "../src/config.js";
import type { UserMeOkPayload } from "../src/types.js";

describe("user", () => {
  const client = createClient(WS_URL);

  beforeEach(async () => {
    await client.connect();
  });

  afterEach(() => {
    client.disconnect();
  });

  it("user.me without login returns unauthorized", async () => {
    await expect(client.sendAndWait("user.me", {})).rejects.toThrow(/unauthorized/);
  });

  it("user.me after login returns user.me.ok with id, nickname, etc.", async () => {
    await client.sendAndWait("auth.login", {
      username: process.env.TEST_USER ?? "testuser",
      password: process.env.TEST_PASSWORD ?? "password123",
      deviceId: "test-dev",
    });

    const env = await client.sendAndWait<UserMeOkPayload>("user.me", {});
    expect(env.type).toBe("user.me.ok");
    expect(env.error).toBeFalsy();
    expect(env.payload).toBeDefined();
    const p = env.payload as UserMeOkPayload;
    expect(p.id).toBeDefined();
    expect(typeof p.nickname).toBe("string");
    expect(typeof p.avatarUrl).toBe("string");
    expect(typeof p.bio).toBe("string");
    expect(typeof p.status).toBe("string");
  });
});
