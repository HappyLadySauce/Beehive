/**
 * 从 asyncapi.yaml 按 tag 生成各模块的接口文档（Markdown）。
 * 输出到 docs/API/generated/ws-{tag}.md。
 */
import fs from "fs";
import path from "path";
import { fileURLToPath } from "url";
import yaml from "yaml";

const __dirname = path.dirname(fileURLToPath(import.meta.url));
const repoRoot = path.resolve(__dirname, "../../..");
const asyncapiPath = path.join(repoRoot, "test/interface/asyncapi.yaml");
const outDir = path.join(repoRoot, "docs/API/generated");

const spec = yaml.parse(fs.readFileSync(asyncapiPath, "utf8"));
const messages = spec.components?.messages ?? {};
const tagToMessages = new Map();

for (const [key, msg] of Object.entries(messages)) {
  const tags = msg.tags ?? [];
  const tagNames = tags.map((t) => (typeof t === "string" ? t : t?.name)).filter(Boolean);
  if (tagNames.length === 0) tagNames.push("common");
  for (const tag of tagNames) {
    if (!tagToMessages.has(tag)) tagToMessages.set(tag, []);
    tagToMessages.get(tag).push({ key, ...msg });
  }
}

if (!fs.existsSync(outDir)) fs.mkdirSync(outDir, { recursive: true });

const titleByTag = {
  auth: "Auth（认证与登出）",
  presence: "Presence（心跳与在线）",
  user: "User（用户资料）",
  common: "通用",
};

for (const [tag, list] of tagToMessages) {
  if (tag === "common") continue;
  const title = titleByTag[tag] ?? tag;
  const lines = [
    `# WebSocket 接口：${title}`,
    "",
    "本文档由 `test/interface/asyncapi.yaml` 自动生成，与 [websocket-client-api.md](../websocket-client-api.md) 对齐。",
    "",
    "## 消息类型",
    "",
  ];
  for (const msg of list) {
    const name = msg.name ?? msg.key;
    const summary = msg.title ?? msg.summary ?? "";
    lines.push(`### \`${name}\``);
    if (summary) lines.push("", summary, "");
    lines.push("");
  }
  const outPath = path.join(outDir, `ws-${tag}.md`);
  fs.writeFileSync(outPath, lines.join("\n"), "utf8");
  console.log("Wrote", outPath);
}

const indexPath = path.join(outDir, "README.md");
const indexLines = [
  "# WebSocket 按模块接口文档",
  "",
  "由 `pnpm run docs:generate` 从 [test/interface/asyncapi.yaml](../../test/interface/asyncapi.yaml) 生成。",
  "",
  "| 模块 | 文档 |",
  "|------|------|",
  ...Array.from(tagToMessages.keys())
    .filter((t) => t !== "common")
    .map((t) => `| ${titleByTag[t] ?? t} | [ws-${t}.md](./ws-${t}.md) |`),
  "",
  "完整协议说明见 [websocket-client-api.md](../websocket-client-api.md)。",
  "",
];
fs.writeFileSync(indexPath, indexLines.join("\n"), "utf8");
console.log("Wrote", indexPath);
