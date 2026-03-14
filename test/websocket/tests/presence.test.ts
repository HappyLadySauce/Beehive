/**
 * Presence 模块 WebSocket 接口测试：presence.ping（未登录/已登录）。
 * 默认使用种子用户 testuser/password123（见 db/migrations/003_seed_test_user.sql）。
 */
import { describe, it, expect, beforeEach, afterEach } from "vitest";
import { createClient } from "../src/client.js";
import { WS_URL } from "../src/config.js";
import type { PresencePingOkPayload } from "../src/types.js";

describe("presence", () => {
  const client = createClient(WS_URL);

  beforeEach(async () => {
    await client.connect();
  });

  afterEach(() => {
    client.disconnect();
  });

  it("presence.ping without login returns unauthorized", async () => {
    await expect(
      client.sendAndWait("presence.ping", { clientTime: Math.floor(Date.now() / 1000) })
    ).rejects.toThrow(/unauthorized/);
  });

  it("presence.ping after login returns presence.ping.ok with serverTime", async () => {
    await client.sendAndWait("auth.login", {
      username: process.env.TEST_USER ?? "testuser",
      password: process.env.TEST_PASSWORD ?? "password123",
      deviceId: "test-dev",
    });

    const env = await client.sendAndWait<PresencePingOkPayload>("presence.ping", {
      clientTime: Math.floor(Date.now() / 1000),
    });
    expect(env.type).toBe("presence.ping.ok");
    expect(env.error).toBeFalsy();
    expect(env.payload).toBeDefined();
    expect(typeof (env.payload as PresencePingOkPayload).serverTime).toBe("number");
  });
});
