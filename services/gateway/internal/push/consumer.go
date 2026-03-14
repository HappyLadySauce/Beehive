package push

import (
	"context"
	"encoding/json"
	"sync"

	amqp "github.com/rabbitmq/amqp091-go"
	"github.com/zeromicro/go-zero/core/logx"

	"github.com/HappyLadySauce/Beehive/services/conversation/conversationservice"
	"github.com/HappyLadySauce/Beehive/services/gateway/internal/config"
	"github.com/HappyLadySauce/Beehive/services/gateway/internal/ws"
	"github.com/HappyLadySauce/Beehive/services/presence/presenceservice"
)

// messageCreatedEvent 与 Message 服务 PostMessage 发布的 JSON 一致。
type messageCreatedEvent struct {
	ServerMsgId    string                 `json:"serverMsgId"`
	ClientMsgId   string                 `json:"clientMsgId"`
	ConversationId string                 `json:"conversationId"`
	FromUserId    string                 `json:"fromUserId"`
	ToUserId      string                 `json:"toUserId"`
	Body          map[string]interface{} `json:"body"`
	ServerTime    float64                `json:"serverTime"` // JSON number
}

// Consumer 消费 message.created 并向本实例连接推送 message.push。
type Consumer struct {
	cfg       config.Config
	hub       *ws.Hub
	conv      conversationservice.ConversationService
	pres      presenceservice.PresenceService
	conn      *amqp.Connection
	ch        *amqp.Channel
	queueName string
	closed    bool
	mu        sync.Mutex
}

// NewConsumer 创建推送消费者；调用方需在退出时调用 Close。
func NewConsumer(cfg config.Config, hub *ws.Hub, conv conversationservice.ConversationService, pres presenceservice.PresenceService) (*Consumer, error) {
	exchange := cfg.RabbitMQExchange
	if exchange == "" {
		exchange = "im.events"
	}
	routeKey := cfg.RabbitMQRouteKey
	if routeKey == "" {
		routeKey = "message.created"
	}
	queue := cfg.RabbitMQQueue
	if queue == "" {
		queue = "gateway.push." + cfg.GatewayID
		if queue == "gateway.push." {
			queue = "gateway.push.default"
		}
	}

	conn, err := amqp.Dial(cfg.RabbitMQURL)
	if err != nil {
		return nil, err
	}
	ch, err := conn.Channel()
	if err != nil {
		_ = conn.Close()
		return nil, err
	}
	if err := ch.ExchangeDeclare(exchange, "topic", true, false, false, false, nil); err != nil {
		_ = ch.Close()
		_ = conn.Close()
		return nil, err
	}
	q, err := ch.QueueDeclare(queue, true, false, false, false, nil)
	if err != nil {
		_ = ch.Close()
		_ = conn.Close()
		return nil, err
	}
	if err := ch.QueueBind(q.Name, routeKey, exchange, false, nil); err != nil {
		_ = ch.Close()
		_ = conn.Close()
		return nil, err
	}

	return &Consumer{cfg: cfg, hub: hub, conv: conv, pres: pres, conn: conn, ch: ch, queueName: q.Name}, nil
}

// Run 在调用方 goroutine 中阻塞消费；返回时表示连接关闭或 Close 被调用。
func (c *Consumer) Run(ctx context.Context) {
	deliveries, err := c.ch.Consume(c.queueName, "", false, false, false, false, nil)
	if err != nil {
		logx.Errorf("push consumer: Consume failed: %v", err)
		return
	}
	for {
		select {
		case <-ctx.Done():
			return
		case d, ok := <-deliveries:
			if !ok {
				return
			}
			c.handleDelivery(context.Background(), &d)
		}
	}
}

func (c *Consumer) handleDelivery(ctx context.Context, d *amqp.Delivery) {
	var ev messageCreatedEvent
	if err := json.Unmarshal(d.Body, &ev); err != nil {
		logx.Errorf("push consumer: unmarshal message.created failed: %v", err)
		_ = d.Nack(false, false)
		return
	}
	if ev.ConversationId == "" {
		_ = d.Ack(false)
		return
	}

	membersResp, err := c.conv.ListMembers(ctx, &conversationservice.ListMembersRequest{ConversationId: ev.ConversationId})
	if err != nil {
		logx.Errorf("push consumer: ListMembers failed conversationId=%s: %v", ev.ConversationId, err)
		_ = d.Nack(false, true)
		return
	}

	gatewayID := c.cfg.GatewayID
	if gatewayID == "" {
		gatewayID = "gateway-1"
	}
	serverTime := int64(ev.ServerTime)
	payload := map[string]interface{}{
		"serverMsgId":    ev.ServerMsgId,
		"clientMsgId":   ev.ClientMsgId,
		"conversationId": ev.ConversationId,
		"fromUserId":    ev.FromUserId,
		"toUserId":      ev.ToUserId,
		"body":          ev.Body,
		"serverTime":    serverTime,
	}

	seenConn := make(map[string]struct{})
	for _, m := range membersResp.Items {
		if m == nil || m.UserId == "" {
			continue
		}
		if m.Status != "" && m.Status != "active" {
			continue
		}
		sessResp, err := c.pres.GetOnlineSessions(ctx, &presenceservice.GetOnlineSessionsRequest{UserId: m.UserId})
		if err != nil {
			logx.Errorf("push consumer: GetOnlineSessions failed userId=%s: %v", m.UserId, err)
			continue
		}
		for _, s := range sessResp.Sessions {
			if s == nil || s.GatewayId != gatewayID {
				continue
			}
			if _, ok := seenConn[s.ConnId]; ok {
				continue
			}
			seenConn[s.ConnId] = struct{}{}
			conn := c.hub.Get(s.ConnId)
			if conn == nil {
				continue
			}
			if err := conn.WriteJSON(&ws.Envelope{
				Type:    "message.push",
				Tid:     ev.ClientMsgId,
				Payload: payload,
				Error:   nil,
			}); err != nil {
				logx.Errorf("push consumer: WriteJSON failed connId=%s: %v", s.ConnId, err)
			}
		}
	}
	_ = d.Ack(false)
}

// Close 关闭连接与 channel。
func (c *Consumer) Close() error {
	c.mu.Lock()
	defer c.mu.Unlock()
	if c.closed {
		return nil
	}
	c.closed = true
	if c.ch != nil {
		_ = c.ch.Close()
		c.ch = nil
	}
	if c.conn != nil {
		_ = c.conn.Close()
		c.conn = nil
	}
	return nil
}
