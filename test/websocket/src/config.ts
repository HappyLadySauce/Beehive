/**
 * 测试用配置，WS_URL 可从环境变量覆盖（如 CI）。
 */
export const WS_URL = process.env.WS_URL ?? "ws://127.0.0.1:8080/ws";
