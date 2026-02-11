package client

import (
	"context"
	"encoding/json"
	"time"

	"github.com/HappyLadySauce/Beehive/common/libnet"
	"github.com/HappyLadySauce/Beehive/imrpc/imrpcclient"

	"github.com/zeromicro/go-zero/core/logx"
)

const heartBeatTimeout = time.Second * 60

// 客户端
type Client struct {
	Session *libnet.Session
	Manager *libnet.Manager
	IMRPC   imrpcclient.Imrpc
	heartbeat	chan *libnet.Message
}

func NewClient(manager *libnet.Manager, session *libnet.Session, imrpc imrpcclient.Imrpc) *Client {
	return &Client{
		Session: session,
		Manager: manager,
		IMRPC: imrpc,
		heartbeat: make(chan *libnet.Message),
	}
}

func (c *Client) Login(msg *libnet.Message) error {
	loginReq, err := makeLoginMessage(msg)
	if err != nil {
		return err
	}

	c.Session.SetToken(loginReq.Token)
	c.Manager.AddSession(c.Session)

	_, err = c.IMRPC.Login(context.Background(), &imrpcclient.LoginRequest{
		Token: loginReq.Token,
		Authorization: loginReq.Authorization,
		SessionId: c.Session.SessionId().String(),
	})

	msg.Status = 0
	msg.Data = []byte("Login Success")
	err = c.Send(*msg)
	if err != nil {
		logx.Errorf("[Login] c.Send msg: %v, error: %v", msg, err)
	}

	return nil
}

func (c *Client) Receive() (*libnet.Message, error) {
	return c.Session.Receive()
}

func (c *Client) Send(msg libnet.Message) error {
	return c.Session.Send(msg)
}

func (c *Client) Close() error {
	return c.Session.Close()
}

func (c *Client) HandlPackage(msg *libnet.Message) error {
	// 消息转发
	req := makePostMessage(c.Session.SessionId().String(), msg)
	if req == nil {
		return nil
	}
	_, err := c.IMRPC.PostMessage(context.Background(), req)
	if err != nil {
		logx.Errorf("[HandlePackage] client.PostMessage error: %v", err)
	}

	return err
}

// 心跳检测
func (c *Client) Heartbeat() {
	timer := time.NewTimer(heartBeatTimeout)
	for {
		select {
		case heartbeat := <-c.heartbeat:
			c.Session.SetReadDeadline(time.Now().Add(heartBeatTimeout * 5))
			c.Send(*heartbeat)
			break
		case <- timer.C:
		}
	}
}

// 生成登录消息
func makeLoginMessage(msg *libnet.Message) (*imrpcclient.LoginRequest, error) {
	// TODO: 添加token和authorization的校验
	var loginReq imrpcclient.LoginRequest
	err := json.Unmarshal(msg.Data, &loginReq)
	if err != nil {
		return nil, err
	}

	return &loginReq, nil
}

// 生成消息转发消息
func makePostMessage(sessionId string, msg *libnet.Message) *imrpcclient.PostMsg {
	// TODO: 添加token和authorization的校验
	var postMessageReq imrpcclient.PostMsg
	err := json.Unmarshal(msg.Data, &postMessageReq)
	if err != nil {
		logx.Errorf("[makePostMessage] json.Unmarshal msg: %v error: %v", msg, err)
		return nil
	}
	postMessageReq.SessionId = sessionId

	return &postMessageReq
}
