package mq

import (
	"encoding/json"
	"sync"

	amqp "github.com/rabbitmq/amqp091-go"
)

// Publisher 向 RabbitMQ 交换机发布消息；URL 为空时不连接，Publish 为 no-op。
type Publisher struct {
	conn    *amqp.Connection
	ch      *amqp.Channel
	exchange string
	routeKey string
	mu      sync.Mutex
}

// NewPublisher 连接 RabbitMQ；若 url 为空则返回 nil（调用方可不发布事件）。
func NewPublisher(url, exchange, routeKey string) (*Publisher, error) {
	if url == "" {
		return nil, nil
	}
	if exchange == "" {
		exchange = "im.events"
	}
	if routeKey == "" {
		routeKey = "message.created"
	}
	conn, err := amqp.Dial(url)
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
	return &Publisher{conn: conn, ch: ch, exchange: exchange, routeKey: routeKey}, nil
}

// PublishJSON 将 v 序列化为 JSON 并发布到配置的 exchange + routeKey。
func (p *Publisher) PublishJSON(v interface{}) error {
	if p == nil {
		return nil
	}
	body, err := json.Marshal(v)
	if err != nil {
		return err
	}
	return p.Publish(body)
}

// Publish 发布原始 body。
func (p *Publisher) Publish(body []byte) error {
	if p == nil {
		return nil
	}
	p.mu.Lock()
	defer p.mu.Unlock()
	return p.ch.Publish(p.exchange, p.routeKey, false, false, amqp.Publishing{
		ContentType: "application/json",
		Body:        body,
	})
}

// Close 关闭连接与 channel。
func (p *Publisher) Close() error {
	if p == nil {
		return nil
	}
	p.mu.Lock()
	defer p.mu.Unlock()
	if p.ch != nil {
		_ = p.ch.Close()
		p.ch = nil
	}
	if p.conn != nil {
		_ = p.conn.Close()
		p.conn = nil
	}
	return nil
}
