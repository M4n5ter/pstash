package natsinput

import (
	"context"
	"github.com/m4n5ter/pstash/stash/config"

	// 标准库
	"errors"
	"fmt"
	"strings"
	"unsafe"

	// 第三方库
	"github.com/nats-io/nats.go"
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
	_ = n.subscription.Drain()
	_ = n.conn.Drain()
	close(n.exit)
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
