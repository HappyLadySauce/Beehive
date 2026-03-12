## IM Web 客户端设计（Vue 3 + Naive UI）

本节描述 Beehive IM Web 客户端的页面结构、路由设计以及核心组件与状态模型。

---

### 1. 页面与路由结构

#### 1.1 顶层路由

建议的路由结构（仅示意）：

- `/auth/login`：登录页。
- `/auth/register`：注册页（可选）。
- `/app`：主应用框架。
  - `/app/chats`：默认展示会话列表 + 当前会话。
  - `/app/contacts`：联系人/群组列表与管理。
  - `/app/settings`：个人设置（账号、安全、外观等）。

可能的 `vue-router` 配置草图：

```ts
// 仅为示意代码，实际实现可在前端项目中调整
const routes = [
  { path: '/auth/login', name: 'Login', component: LoginPage },
  { path: '/auth/register', name: 'Register', component: RegisterPage },
  {
    path: '/app',
    component: AppShell,
    children: [
      { path: 'chats/:conversationId?', name: 'Chats', component: ChatPage },
      { path: 'contacts', name: 'Contacts', component: ContactsPage },
      { path: 'settings', name: 'Settings', component: SettingsPage },
    ],
  },
  { path: '/:pathMatch(.*)*', redirect: '/auth/login' },
];
```

---

### 2. 主要页面布局

#### 2.1 登录/注册页

- 使用 Naive UI 的表单组件构建：
  - 输入：用户名/邮箱、密码（+ 确认密码）、可选验证码；
  - 提交时调用后端 AuthService 的 HTTP 接口或通过 WebSocket `auth.login` 消息完成登录；
  - 成功后保存 token 与用户信息到全局 store，跳转 `/app/chats`。

#### 2.2 主聊天页（AppShell + ChatPage）

- `AppShell` 布局：
  - 顶部：Logo、当前用户头像与菜单（个人资料、退出登录）。
  - 左列：会话列表导航，可切换到联系人/群组视图。
  - 中间：当前会话消息面板。
  - 右侧（可选显示）：会话详情或用户资料面板。

- `ChatPage` 内部区域划分：
  - `ConversationList`：左侧会话列表；
  - `MessagePanel`：中间消息时间线；
  - `MessageInput`：底部输入框与发送按钮；
  - `ConversationSidePanel`：右侧会话详情（可收起）。

---

### 3. 核心组件设计

#### 3.1 AppShell

- 负责页面框架与导航布局。
- 与认证 store 集成：
  - 当未登录时重定向到 `/auth/login`；
  - 展示当前用户信息与登出入口。

#### 3.2 ConversationList

- 展示最近会话列表：
  - 会话名称/头像；
  - 最后一条消息摘要；
  - 未读计数；
  - 最近活跃时间。
- 支持：
  - 按名称/关键字搜索；
  - 置顶/收藏会话；
  - 点击切换当前会话并路由到 `/app/chats/:conversationId`。

#### 3.3 MessagePanel

- 展示选中会话的消息时间线：
  - 按时间排序分组，可根据发送者对齐左右（自己在右，其他人在左）；
  - 每条消息显示：头像、昵称、发送时间、正文；
  - 支持基本的消息类型：
    - 文本；
    - Emoji（后续可扩展）；
    - 图片/文件占位（可后续实现上传）。
- 支持滚动加载历史：
  - 滚动接近顶部时触发拉取更早的记录；
  - 使用游标或时间戳分页，调用 MessageService 的历史接口。

#### 3.4 MessageInput

- 底部输入区：
  - 多行输入框；
  - 发送按钮；
  - 支持 Enter 发送、Shift+Enter 换行；
  - 显示发送状态（发送中、失败可重试）。
- 发送逻辑：
  - 调用 WebSocket `message.send` 消息，将 `clientMsgId`、`conversationId`、`body` 发送给 Gateway；
  - 本地先插入一条「待确认」状态的消息气泡；
  - 收到 `message.send.ok` 或 `message.push` 后更新状态。

#### 3.5 ContactsPage / ContactAndGroup

- 提供联系人列表与群组列表管理：
  - 按字母或搜索筛选联系人；
  - 查看某人详情，发起单聊；
  - 管理自己所属的群组，创建新群聊、添加/移除成员等。

#### 3.6 SettingsPage / Profile & Preferences

- 个人资料：
  - 昵称、头像、个人签名等；
  - 调用 UserService 接口保存。
- 安全与账号：
  - 修改密码；
  - 管理多设备登录（可与 PresenceService 集成）。
- 外观设置：
  - 亮/暗主题切换；
  - 字号、语言选择（如接入 i18n）。

---

### 4. 状态管理模型（Pinia）

#### 4.1 用户与认证 store

- `useAuthStore` 示例字段：
  - `user`: 当前登录用户信息（id、昵称、头像等）；
  - `accessToken` / `refreshToken`；
  - `loggedIn`：布尔值；
  - 操作：
    - `login(credentials)`；
    - `logout()`；
    - `refreshToken()`。

#### 4.2 会话 store

- `useConversationStore`：
  - `conversations`: 会话数组，每个包含：
    - `id`, `name`, `avatar`, `lastMessage`, `lastActiveAt`, `unreadCount` 等；
  - `currentConversationId`；
  - 操作：
    - `setConversations(list)`；
    - `setCurrentConversation(id)`；
    - `updateLastMessage(conversationId, message)`；
    - `incrementUnread(conversationId)` / `clearUnread(conversationId)`。

#### 4.3 消息 store

- `useMessageStore`：
  - `messagesByConversation`: 以 `conversationId` 为 key 的消息数组；
  - `paginationByConversation`: 保存每个会话历史分页游标；
  - 操作：
    - `appendMessage(conversationId, message)`；
    - `prependHistory(conversationId, messages)`；
    - `updateMessageStatus(conversationId, clientMsgId, status)`。

---

### 5. WebSocket 集成设计

#### 5.1 useWebSocket Hook（草案）

职责：

- 管理与 Gateway 的 WebSocket 连接：
  - 建立连接（带上 token 或在连接后发送 `auth.tokenLogin`）；
  - 自动重连（带退避策略）；
  - 定期发送心跳（如 `presence.ping`），配合后端维持在线状态。
- 提供统一的 `send` 接口：
  - 自动封装为约定的 JSON Envelope（`type/tid/payload`）。
- 派发消息：
  - 根据 `type` 将消息派发到对应回调，例如：
    - `auth.login.ok` → 触发登录成功处理；
    - `message.push` → 写入 `useMessageStore`；
    - `conversation.updated` → 更新会话信息等。

#### 5.2 与 store 的协作

- 登录：
  - UI 调用 `authStore.login()`，内部通过 HTTP 或 WebSocket `auth.login` 完成认证；
  - 登录成功后，初始化 WebSocket 连接或使用 token 登录；
  - 加载用户会话列表与基础资料。
- 收消息：
  - `useWebSocket` 收到 `message.push`，调用 `messageStore.appendMessage()`；
  - 更新对应会话的 `lastMessage` 与 `unreadCount`。
- 发消息：
  - 由 `MessageInput` 调用 `useWebSocket.send('message.send', payload)`；
  - 发送前在 `messageStore` 中插入一条本地 `pending` 消息；
  - 收到成功 ACK 后更新状态为 `sent`。

