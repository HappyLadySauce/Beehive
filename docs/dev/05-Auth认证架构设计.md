# Auth 认证架构设计

## 概述

本文档分析在 Beehive IM 微服务架构中，认证（Auth）功能应该放置在哪个位置，以及各种方案的优缺点。

## 当前架构状态

根据现有文档，当前架构中：

- **User Service** 负责：
  - 用户注册
  - 用户登录和认证
  - JWT Token 生成和验证
  - 用户信息查询和更新

- **Gateway** 在建立 WebSocket 连接时：
  - 从请求头提取 Token
  - 调用 User Service 的 `ValidateToken` gRPC 方法验证 Token
  - 验证通过后建立连接

## 可选方案分析

### 方案一：保持在 User Service 中（当前方案）

**架构图：**
```
客户端 → Gateway → User Service (验证Token) → 返回用户信息
```

**优点：**
1. ✅ **职责清晰**：认证逻辑与用户数据紧密相关，符合单一职责原则
2. ✅ **数据一致性**：Token 验证可能需要查询用户状态（如用户是否被禁用），放在 User Service 可以保证数据一致性
3. ✅ **安全性高**：JWT Secret 只存储在 User Service，减少泄露风险
4. ✅ **易于扩展**：未来如果需要 Token 黑名单、刷新 Token 等复杂功能，都在一个服务中实现
5. ✅ **架构简单**：不需要额外的服务或复杂的配置

**缺点：**
1. ❌ **性能开销**：每次 Token 验证都需要 gRPC 调用，增加延迟
2. ❌ **Gateway 依赖**：Gateway 必须依赖 User Service 才能工作
3. ❌ **网络开销**：即使 Token 验证是轻量级操作，也需要网络往返

**适用场景：**
- 中小型系统
- Token 验证需要查询数据库的场景
- 需要集中管理认证逻辑的场景

---

### 方案二：在 Gateway 层实现认证中间件

**架构图：**
```
客户端 → Gateway (本地验证Token) → 业务服务
```

**实现方式：**
- Gateway 持有 JWT Secret，本地验证 Token
- 验证通过后，将用户信息注入到请求上下文
- 业务服务从上下文获取用户信息，无需再次验证

**优点：**
1. ✅ **性能最优**：Token 验证在 Gateway 本地完成，无网络开销
2. ✅ **减少依赖**：Gateway 可以独立验证 Token，减少对 User Service 的依赖
3. ✅ **统一入口**：所有认证逻辑集中在 Gateway，便于统一管理
4. ✅ **降低延迟**：WebSocket 连接建立更快

**缺点：**
1. ❌ **安全风险**：JWT Secret 需要在 Gateway 和 User Service 之间共享，增加泄露风险
2. ❌ **功能受限**：如果 Token 验证需要查询数据库（如检查用户状态），仍需要调用 User Service
3. ❌ **配置复杂**：需要同步 JWT Secret 配置
4. ❌ **职责混合**：Gateway 承担了认证职责，可能违反单一职责原则

**适用场景：**
- 高性能要求的系统
- Token 验证是纯 JWT 验证，不需要查询数据库
- Gateway 和 User Service 可以安全共享 Secret 的场景

---

### 方案三：独立的 Auth Service

**架构图：**
```
客户端 → Gateway → Auth Service (验证Token) → 返回用户信息
         ↓
    业务服务层
```

**实现方式：**
- 创建独立的 Auth Service，专门负责认证和授权
- User Service 只负责用户数据管理
- Gateway 调用 Auth Service 验证 Token

**优点：**
1. ✅ **职责分离**：认证和用户管理完全分离
2. ✅ **独立扩展**：Auth Service 可以独立扩展和优化
3. ✅ **灵活性高**：可以支持多种认证方式（JWT、OAuth2、API Key 等）
4. ✅ **易于维护**：认证逻辑集中在一个服务中

**缺点：**
1. ❌ **过度设计**：对于当前系统规模，可能过于复杂
2. ❌ **增加服务**：增加了一个微服务，提高了运维复杂度
3. ❌ **网络开销**：仍然需要 gRPC 调用
4. ❌ **数据依赖**：Auth Service 可能仍需要查询 User Service 获取用户信息

**适用场景：**
- 大型系统，需要多种认证方式
- 需要独立的认证授权服务
- 认证逻辑非常复杂，需要独立团队维护

---

### 方案四：混合方案（推荐）

**架构图：**
```
客户端 → Gateway (快速验证) → User Service (深度验证) → 业务服务
```

**实现方式：**
1. **Gateway 层**：实现轻量级 Token 验证
   - 验证 JWT 签名和过期时间
   - 缓存验证结果（Redis）
   - 快速拒绝无效 Token

2. **User Service**：提供深度验证
   - 验证用户状态（是否被禁用）
   - Token 黑名单检查
   - 刷新 Token 逻辑

3. **验证流程**：
   ```
   Gateway 收到请求
   ↓
   检查 Redis 缓存（用户ID、状态）
   ↓
   缓存命中 → 直接使用
   缓存未命中 → 调用 User Service 验证 → 缓存结果
   ```

**优点：**
1. ✅ **性能优化**：大部分请求通过缓存快速验证
2. ✅ **功能完整**：支持复杂的验证逻辑
3. ✅ **安全性高**：深度验证仍在 User Service
4. ✅ **可扩展**：可以根据需要调整验证策略

**缺点：**
1. ❌ **实现复杂**：需要实现缓存逻辑
2. ❌ **一致性**：需要处理缓存和数据库的一致性

**适用场景：**
- 中大型系统
- 需要平衡性能和功能的场景
- 有 Redis 基础设施

---

## 推荐方案

### 短期推荐：方案一（保持在 User Service）

**理由：**
1. **符合当前架构**：与现有设计一致，无需重构
2. **简单可靠**：架构清晰，易于理解和维护
3. **功能完整**：支持所有认证需求，包括用户状态检查
4. **安全性好**：JWT Secret 集中管理

**优化建议：**
- 在 Gateway 和 User Service 之间使用连接池，减少连接开销
- 考虑在 Gateway 层添加 Token 缓存（Redis），减少 gRPC 调用
- 使用 gRPC 拦截器实现统一的认证逻辑

### 长期推荐：方案四（混合方案）

**理由：**
1. **性能优化**：通过缓存大幅减少 gRPC 调用
2. **功能完整**：支持复杂的认证需求
3. **可扩展性**：为未来扩展预留空间

**实施步骤：**
1. 第一阶段：保持当前架构，添加 Redis 缓存层
2. 第二阶段：在 Gateway 实现轻量级验证，缓存验证结果
3. 第三阶段：根据实际需求，考虑是否需要独立的 Auth Service

---

## 实施建议

### 如果选择方案一（当前方案）

**优化措施：**

1. **添加 Token 缓存**
```go
// Gateway 层添加缓存
func (h *Handler) validateTokenWithCache(token string) (*Claims, error) {
    // 1. 检查 Redis 缓存
    cacheKey := fmt.Sprintf("token:%s", token)
    cached, err := h.redis.Get(ctx, cacheKey).Result()
    if err == nil {
        // 缓存命中，解析并返回
        return parseClaims(cached)
    }
    
    // 2. 缓存未命中，调用 User Service
    claims, err := h.userService.ValidateToken(token)
    if err != nil {
        return nil, err
    }
    
    // 3. 缓存结果（设置过期时间略小于 Token 过期时间）
    h.redis.Set(ctx, cacheKey, claims, 23*time.Hour)
    return claims, nil
}
```

2. **使用 gRPC 连接池**
```go
// 复用 gRPC 连接
var userServiceConn *grpc.ClientConn

func init() {
    conn, err := grpc.Dial(userServiceAddr, grpc.WithInsecure())
    if err != nil {
        log.Fatal(err)
    }
    userServiceConn = conn
}
```

3. **批量验证优化**（如果未来需要）
```go
// User Service 支持批量验证
rpc ValidateTokens(ValidateTokensRequest) returns (ValidateTokensResponse);
```

### 如果选择方案四（混合方案）

**实施步骤：**

1. **第一步：添加 Redis 缓存层**
   - 在 Gateway 中集成 Redis
   - 实现 Token 验证结果缓存
   - 设置合理的缓存过期时间

2. **第二步：Gateway 轻量级验证**
   - Gateway 持有 JWT Secret（通过配置中心或环境变量）
   - 实现本地 JWT 签名验证
   - 对于需要深度验证的场景，仍调用 User Service

3. **第三步：监控和优化**
   - 监控缓存命中率
   - 监控 gRPC 调用频率
   - 根据实际数据调整策略

---

## 性能对比

| 方案 | Token 验证延迟 | 网络开销 | 实现复杂度 | 安全性 |
|------|--------------|---------|-----------|--------|
| 方案一（当前） | ~5-10ms | 1 次 gRPC | 低 | 高 |
| 方案一+缓存 | ~1-2ms (缓存命中) | 0-1 次 gRPC | 中 | 高 |
| 方案二（Gateway） | ~0.1ms | 0 次 | 中 | 中 |
| 方案三（独立服务） | ~5-10ms | 1 次 gRPC | 高 | 高 |
| 方案四（混合） | ~0.1-2ms | 0-1 次 gRPC | 高 | 高 |

---

## 安全考虑

### JWT Secret 管理

无论选择哪种方案，都需要注意：

1. **不要硬编码 Secret**：使用环境变量或配置中心
2. **定期轮换 Secret**：实现 Secret 轮换机制
3. **最小权限原则**：Gateway 如果持有 Secret，应该只有验证权限
4. **审计日志**：记录所有认证相关的操作

### Token 验证策略

1. **签名验证**：必须验证 JWT 签名
2. **过期检查**：检查 Token 是否过期
3. **用户状态**：验证用户是否被禁用（需要查询数据库）
4. **Token 黑名单**：支持 Token 撤销（可选）

---

## 结论

**已选择方案：方案三（独立的 Auth Service）**

**实施状态：** ✅ 已确定架构方案，文档已更新

**架构优势：**
1. ✅ **职责分离**：认证和用户管理完全分离，符合单一职责原则
2. ✅ **独立扩展**：Auth Service 可以独立扩展和优化，不影响其他服务
3. ✅ **灵活性高**：可以支持多种认证方式（JWT、OAuth2、API Key 等）
4. ✅ **易于维护**：认证逻辑集中在一个服务中，便于团队协作
5. ✅ **高性能**：通过 Redis 缓存 Token 验证结果，提升性能

**服务职责划分：**

- **Auth Service**：负责所有认证授权相关功能
  - 用户登录验证
  - JWT Token 生成和验证
  - Token 刷新和撤销
  - 权限验证
  - Token 黑名单管理

- **User Service**：专注于用户数据管理
  - 用户注册
  - 用户信息查询和更新
  - 用户资料管理
  - 用户状态管理

**实施路径：**
```
✅ 方案三（独立 Auth Service）
    ↓
实施阶段：
1. 创建 Auth Service 项目结构
2. 实现 Protocol Buffers 定义
3. 实现认证逻辑（登录、Token 生成/验证）
4. 集成 Redis 缓存
5. 更新 Gateway 调用 Auth Service
6. 迁移 User Service 认证逻辑
```

---

## 方案三实施指南

### 1. 创建 Auth Service 项目结构

```bash
# 创建目录结构
mkdir -p cmd/auth-service/cmd
mkdir -p internal/service/auth
mkdir -p api/proto/auth
mkdir -p configs
```

### 2. 定义 Protocol Buffers

```protobuf
// api/proto/auth/auth.proto
syntax = "proto3";

package auth;

option go_package = "github.com/HappyLadySauce/Beehive/api/proto/auth";

service AuthService {
    // 用户登录
    rpc Login(LoginRequest) returns (LoginResponse);
    
    // 验证 Token
    rpc ValidateToken(ValidateTokenRequest) returns (ValidateTokenResponse);
    
    // 刷新 Token
    rpc RefreshToken(RefreshTokenRequest) returns (RefreshTokenResponse);
    
    // 撤销 Token
    rpc RevokeToken(RevokeTokenRequest) returns (RevokeTokenResponse);
}

message LoginRequest {
    string username = 1;
    string password = 2;
}

message LoginResponse {
    string token = 1;
    string refresh_token = 2;
    int64 expires_at = 3;
    UserInfo user = 4;
}

message ValidateTokenRequest {
    string token = 1;
}

message ValidateTokenResponse {
    bool valid = 1;
    string user_id = 2;
    string username = 3;
    UserInfo user = 4;
}

message RefreshTokenRequest {
    string refresh_token = 1;
}

message RefreshTokenResponse {
    string token = 1;
    string refresh_token = 2;
    int64 expires_at = 3;
}

message RevokeTokenRequest {
    string token = 1;
}

message RevokeTokenResponse {
    bool success = 1;
}

message UserInfo {
    string id = 1;
    string username = 2;
    string nickname = 3;
    string avatar = 4;
    string email = 5;
}
```

### 3. 实现 Auth Service

```go
// internal/service/auth/service.go
package auth

import (
    "context"
    "crypto/sha256"
    "encoding/hex"
    "fmt"
    "time"
    
    "github.com/golang-jwt/jwt/v5"
    "github.com/redis/go-redis/v9"
    "google.golang.org/grpc/codes"
    "google.golang.org/grpc/status"
    
    pb "github.com/HappyLadySauce/Beehive/api/proto/auth"
    userpb "github.com/HappyLadySauce/Beehive/api/proto/user"
)

type Service struct {
    pb.UnimplementedAuthServiceServer
    config      *config.Config
    redis       *redis.Client
    userService userpb.UserServiceClient
}

func NewService(cfg *config.Config, redisClient *redis.Client, userClient userpb.UserServiceClient) *Service {
    return &Service{
        config:      cfg,
        redis:       redisClient,
        userService: userClient,
    }
}

func (s *Service) Login(ctx context.Context, req *pb.LoginRequest) (*pb.LoginResponse, error) {
    // 1. 调用 User Service 获取用户信息
    userReq := &userpb.GetUserByUsernameRequest{
        Username: req.Username,
    }
    userResp, err := s.userService.GetUserByUsername(ctx, userReq)
    if err != nil {
        return nil, status.Error(codes.Unauthenticated, "invalid credentials")
    }
    
    user := userResp.User
    
    // 2. 验证密码
    passwordHash := s.hashPassword(req.Password, user.Salt)
    if passwordHash != user.PasswordHash {
        return nil, status.Error(codes.Unauthenticated, "invalid credentials")
    }
    
    // 3. 检查用户状态
    if user.Status != "active" {
        return nil, status.Error(codes.PermissionDenied, "user is not active")
    }
    
    // 4. 生成 Token
    token, err := s.generateToken(user.Id, user.Username)
    if err != nil {
        return nil, status.Error(codes.Internal, "failed to generate token")
    }
    
    // 5. 生成 Refresh Token
    refreshToken, err := s.generateRefreshToken(user.Id)
    if err != nil {
        return nil, status.Error(codes.Internal, "failed to generate refresh token")
    }
    
    // 6. 缓存 Token 信息
    s.cacheToken(token, user.Id)
    
    return &pb.LoginResponse{
        Token:        token,
        RefreshToken: refreshToken,
        ExpiresAt:    time.Now().Add(time.Duration(s.config.JWT.ExpireHours) * time.Hour).Unix(),
        User: &pb.UserInfo{
            Id:       user.Id,
            Username: user.Username,
            Nickname: user.Nickname,
            Avatar:   user.Avatar,
            Email:    user.Email,
        },
    }, nil
}

func (s *Service) ValidateToken(ctx context.Context, req *pb.ValidateTokenRequest) (*pb.ValidateTokenResponse, error) {
    // 1. 检查黑名单
    if s.isTokenBlacklisted(req.Token) {
        return &pb.ValidateTokenResponse{Valid: false}, nil
    }
    
    // 2. 验证 JWT
    claims, err := s.validateJWT(req.Token)
    if err != nil {
        return &pb.ValidateTokenResponse{Valid: false}, nil
    }
    
    // 3. 检查缓存
    cachedUserID := s.getCachedUserID(req.Token)
    if cachedUserID == "" {
        // 缓存未命中，查询用户信息
        userReq := &userpb.GetUserRequest{UserId: claims.UserID}
        userResp, err := s.userService.GetUser(ctx, userReq)
        if err != nil {
            return &pb.ValidateTokenResponse{Valid: false}, nil
        }
        
        // 缓存结果
        s.cacheToken(req.Token, claims.UserID)
        
        return &pb.ValidateTokenResponse{
            Valid:    true,
            UserId:   claims.UserID,
            Username: claims.Username,
            User: &pb.UserInfo{
                Id:       userResp.User.Id,
                Username: userResp.User.Username,
                Nickname: userResp.User.Nickname,
                Avatar:   userResp.User.Avatar,
                Email:    userResp.User.Email,
            },
        }, nil
    }
    
    return &pb.ValidateTokenResponse{
        Valid:    true,
        UserId:   claims.UserID,
        Username: claims.Username,
    }, nil
}

func (s *Service) RevokeToken(ctx context.Context, req *pb.RevokeTokenRequest) (*pb.RevokeTokenResponse, error) {
    // 1. 验证 Token 有效性
    claims, err := s.validateJWT(req.Token)
    if err != nil {
        return &pb.RevokeTokenResponse{Success: false}, nil
    }
    
    // 2. 加入黑名单
    blacklistKey := fmt.Sprintf("token:blacklist:%s", req.Token)
    expireTime := time.Duration(s.config.JWT.ExpireHours) * time.Hour
    s.redis.Set(ctx, blacklistKey, "1", expireTime)
    
    // 3. 清除缓存
    cacheKey := fmt.Sprintf("token:%s", req.Token)
    s.redis.Del(ctx, cacheKey)
    
    return &pb.RevokeTokenResponse{Success: true}, nil
}

// 辅助方法
func (s *Service) hashPassword(password, salt string) string {
    hash := sha256.Sum256([]byte(password + salt))
    return hex.EncodeToString(hash[:])
}

func (s *Service) generateToken(userID, username string) (string, error) {
    claims := &jwt.MapClaims{
        "user_id":  userID,
        "username": username,
        "exp":      time.Now().Add(time.Duration(s.config.JWT.ExpireHours) * time.Hour).Unix(),
        "iat":      time.Now().Unix(),
    }
    
    token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
    return token.SignedString([]byte(s.config.JWT.Secret))
}

func (s *Service) validateJWT(tokenString string) (*Claims, error) {
    token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
        return []byte(s.config.JWT.Secret), nil
    })
    
    if err != nil {
        return nil, err
    }
    
    if claims, ok := token.Claims.(*Claims); ok && token.Valid {
        return claims, nil
    }
    
    return nil, fmt.Errorf("invalid token")
}

func (s *Service) cacheToken(token, userID string) {
    cacheKey := fmt.Sprintf("token:%s", token)
    expireTime := time.Duration(s.config.JWT.ExpireHours) * time.Hour
    s.redis.Set(ctx, cacheKey, userID, expireTime)
}

func (s *Service) getCachedUserID(token string) string {
    cacheKey := fmt.Sprintf("token:%s", token)
    userID, err := s.redis.Get(ctx, cacheKey).Result()
    if err != nil {
        return ""
    }
    return userID
}

func (s *Service) isTokenBlacklisted(token string) bool {
    blacklistKey := fmt.Sprintf("token:blacklist:%s", token)
    exists, _ := s.redis.Exists(ctx, blacklistKey).Result()
    return exists > 0
}
```

### 4. 更新 Gateway 调用

```go
// internal/gateway/grpc/client.go
type Client struct {
    authService    pb.AuthServiceClient
    userService    pb.UserServiceClient
    messageService pb.MessageServiceClient
    presenceService pb.PresenceServiceClient
}

func NewClient(cfg *config.Config) (*Client, error) {
    // 连接 Auth Service
    authConn, err := grpc.Dial(cfg.GRPC.AuthServiceAddr, grpc.WithInsecure())
    if err != nil {
        return nil, err
    }
    
    // ... 其他服务连接
    
    return &Client{
        authService: pb.NewAuthServiceClient(authConn),
        // ...
    }, nil
}

// internal/gateway/websocket/handler.go
func (h *Handler) HandleConnection(w http.ResponseWriter, r *http.Request) {
    token := extractToken(r)
    if token == "" {
        http.Error(w, "Missing authorization token", http.StatusUnauthorized)
        return
    }
    
    // 调用 Auth Service 验证 Token
    req := &pb.ValidateTokenRequest{Token: token}
    resp, err := h.authService.ValidateToken(context.Background(), req)
    if err != nil || !resp.Valid {
        http.Error(w, "Invalid token", http.StatusUnauthorized)
        return
    }
    
    // ... 建立 WebSocket 连接
}
```

### 5. 配置文件更新

```yaml
# configs/auth-service.yaml
server:
  auth_service_port: 50050

jwt:
  secret: your-secret-key-change-in-production
  expire_hours: 24
  refresh_expire_hours: 168

redis:
  addr: localhost:6379
  password: ""
  db: 0

grpc:
  user_service_addr: localhost:50051

log:
  level: info
  format: json
```

### 6. 迁移检查清单

- [ ] 创建 Auth Service 项目结构
- [ ] 定义 Protocol Buffers
- [ ] 实现 Auth Service 核心功能
- [ ] 集成 Redis 缓存
- [ ] 更新 Gateway 调用 Auth Service
- [ ] 从 User Service 移除认证逻辑
- [ ] 更新配置文件
- [ ] 更新文档
- [ ] 编写单元测试
- [ ] 集成测试

---

## 参考

- [用户登录与操作逻辑](./01-用户登录与操作逻辑.md)
- [微服务架构与Cobra框架](./04-微服务架构与Cobra框架.md)
- [WebSocket Gateway 设计](./02-WebSocket-Gateway设计.md)
- [完整开发指南](./00-完整开发指南.md)
