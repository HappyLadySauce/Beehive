/**
 * Auth 模块 WebSocket 接口测试：auth.login、auth.tokenLogin、auth.logout。
 */
import { describe, it, expect, beforeEach, afterEach } from "vitest";
import { createClient } from "../src/client.js";
import { WS_URL } from "../src/config.js";
import type { AuthLoginOkPayload, Envelope } from "../src/types.js";

describe("auth", () => {
  const client = createClient(WS_URL);

  beforeEach(async () => {
    await client.connect();
  });

  afterEach(() => {
    client.disconnect();
  });

  it("auth.login with invalid credentials returns error", async () => {
    await expect(
      client.sendAndWait("auth.login", {
        username: "nonexistent",
        password: "wrong",
        deviceId: "test-dev",
      })
    ).rejects.toThrow(/unauthorized|invalid/);
  });

  it("auth.login with valid credentials returns auth.login.ok", async () => {
    const env = await client.sendAndWait("auth.login", {
      username: process.env.TEST_USER ?? "testuser",
      password: process.env.TEST_PASSWORD ?? "password123",
      deviceId: "test-dev",
    });
    expect(env.type).toBe("auth.login.ok");
    expect(env.error).toBeFalsy();
    expect(env.payload).toBeDefined();
    const p = env.payload as { userId?: string; accessToken?: string; expiresIn?: number };
    expect(p.userId).toBeDefined();
    expect(p.accessToken).toBeDefined();
    expect(typeof p.expiresIn).toBe("number");
  });

  it("auth.logout after login returns auth.logout.ok", async () => {
    const loginEnv = await client.sendAndWait("auth.login", {
      username: process.env.TEST_USER ?? "testuser",
      password: process.env.TEST_PASSWORD ?? "password123",
      deviceId: "test-dev",
    }) as Envelope<{ accessToken?: string }>;
    const token = loginEnv.payload?.accessToken;
    expect(token).toBeDefined();

    const logoutEnv = await client.sendAndWait("auth.logout", {
      accessToken: token ?? "",
    });
    expect(logoutEnv.type).toBe("auth.logout.ok");
    expect(logoutEnv.error).toBeFalsy();
  });

  it("auth.tokenLogin with valid accessToken returns auth.tokenLogin.ok", async () => {
    const loginEnv = await client.sendAndWait<AuthLoginOkPayload>("auth.login", {
      username: process.env.TEST_USER ?? "testuser",
      password: process.env.TEST_PASSWORD ?? "password123",
      deviceId: "test-dev",
    });
    expect(loginEnv.type).toBe("auth.login.ok");
    const loginPayload = loginEnv.payload as AuthLoginOkPayload;
    expect(loginPayload.accessToken).toBeDefined();

    const tokenEnv = await client.sendAndWait("auth.tokenLogin", {
      accessToken: loginPayload.accessToken,
      deviceId: "test-dev-2",
    });
    expect(tokenEnv.type).toBe("auth.tokenLogin.ok");
    expect(tokenEnv.error).toBeFalsy();
    const tokenPayload = tokenEnv.payload as AuthLoginOkPayload;
    expect(tokenPayload.userId).toBe(loginPayload.userId);
  });
});
