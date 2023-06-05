package natsinput

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"unsafe"

	"github.com/nats-io/nats.go"

	"github.com/m4n5ter/pstash/stash/config"
)

var (
	authTypeNone  = "none"
	authTypePass  = "password"
	authTypeToken = "token"
	authTypeNkey  = "nkey"
	authTypeCreds = "creds"
)

var (
	InvalidFormatError    = errors.New("invalid format")
	EmptyAuthContentError = errors.New("empty auth content")
)

type Nats struct {
	conn               *nats.Conn
	subscription       *nats.Subscription
	handler            ConsumeHandler
	needAck            bool
	enableJS           bool
	queueSubscribeName string
	subject            string
	exit               chan struct{}
}

type ConsumeHandler interface {
	Consume(key string, value string) error
}

func NewNats(c config.NatsConf, handler ConsumeHandler) *Nats {
	n := &Nats{}
	n.conn = MustNewNatsConnection(c)
	n.handler = handler
	n.needAck = c.NeedAck
	n.enableJS = c.EnableJetStream
	n.queueSubscribeName = c.QueueSubscribeName
	n.subject = c.Subject
	n.exit = make(chan struct{})

	return n
}

func MustNewNatsConnection(c config.NatsConf) *nats.Conn {
	conn, err := NewNatsConnection(c)
	if err != nil {
		panic(err)
	}
	return conn
}

func NewNatsConnection(c config.NatsConf) (*nats.Conn, error) {
	var opts []nats.Option
	opts = append(opts, nats.Name(c.ConnectionName))

	// if auth type is not none, auth content must not be empty
	if c.AuthType != authTypeNone && c.AuthContent == "" {
		return nil, EmptyAuthContentError
	}

	// TLS
	if c.TlsKeyPath != "" && c.TlsCertPath != "" {
		if c.TlsCaPath != "" {
			opts = append(opts, nats.RootCAs(c.TlsCaPath))
		}
		opts = append(opts, nats.ClientCert(c.TlsCertPath, c.TlsKeyPath))
	}

	switch c.AuthType {
	case authTypeNone:
		// do nothing
	case authTypePass:
		var user, pass string
		if i := strings.Index(c.AuthContent, "@"); i == -1 {
			return nil, fmt.Errorf("can't find '@' in AutyContent : %w", InvalidFormatError)
		} else {
			user = c.AuthContent[:i]
			pass = c.AuthContent[i+1:]
		}

		opts = append(opts, nats.UserInfo(user, pass))
	case authTypeToken:
		opts = append(opts, nats.Token(c.AuthContent))
	case authTypeNkey:
		opt, err := nats.NkeyOptionFromSeed(c.AuthContent)
		if err != nil {
			return nil, fmt.Errorf("get nkey file error : %w", err)
		}
		opts = append(opts, opt)
	case authTypeCreds:
		opts = append(opts, nats.UserCredentials(c.AuthContent))
	}

	return nats.Connect(strings.Join(c.Urls, ","), opts...)
}

func (n *Nats) Start() {
	n.start()
}

func (n *Nats) Stop() {
	n.stop()
}

func (n *Nats) start() {
	n.getSubscription()
	for {
		select {
		case <-n.exit:
			return
		default:
			n.consume()
		}
	}
}

func (n *Nats) stop() {
	close(n.exit)
	_ = n.subscription.Drain()
	_ = n.conn.Drain()
}

func (n *Nats) getSubscription() {
	if n.enableJS {
		durable := nats.Durable("PSTASH")
		js, _ := n.conn.JetStream()

		if n.queueSubscribeName != "" {
			subscription, err := js.QueueSubscribeSync(n.subject, n.queueSubscribeName, durable)
			if err != nil {
				panic(err)
			}
			n.subscription = subscription
		} else {
			subscription, err := js.SubscribeSync(n.subject, durable)
			if err != nil {
				panic(err)
			}
			n.subscription = subscription
		}

	} else {
		subscription, err := n.conn.SubscribeSync(n.subject)
		if err != nil {
			panic(err)
		}
		n.subscription = subscription
	}
}

func (n *Nats) consume() {
	//TODO: 这里无论在消费的时候是否出错，都响应了 ack
	// 但是应该处理一下错误，比如消息格式错误，可以直接响应ack（相当于丢弃了错误格式的消息）
	// 但是如果是其他错误，比如消费者处理消息的时候出错了，或者网络问题（不是消息本身的问题），那么就不能响应ack，这样消息就会一直在队列中
	// 当然如果没用使用 JetStream，那么消息只能看到一次，相当于不论是什么错误，一定会丢失消息
	msg, _ := n.subscription.NextMsgWithContext(context.Background())
	err := n.handler.Consume("", bytes2string(&msg.Data))
	if err != nil {
		fmt.Println(err)
	}

	if n.needAck {
		_ = msg.AckSync()
	}
}

func bytes2string(b *[]byte) string {
	return *(*string)(unsafe.Pointer(b))
}
