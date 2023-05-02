package tcpinput

import (
	"bufio"
	"github.com/m4n5ter/pstash/stash/config"
	"github.com/m4n5ter/pstash/stash/handler"
	"github.com/zeromicro/go-zero/core/logx"
	"net"
	"strconv"
	"time"
)

type TcpInput struct {
	Timeout  int
	exit     chan struct{}
	listener net.Listener
	handler  *handler.MessageHandler
}

func NewTcpInput(c config.TcpInputConf, handler *handler.MessageHandler) *TcpInput {
	l, e := net.Listen("tcp", c.IP+":"+strconv.Itoa(c.Port))
	if e != nil {
		panic(e)
	}
	return &TcpInput{
		listener: l,
		exit:     make(chan struct{}),
		Timeout:  c.Timeout,
		handler:  handler,
	}
}

func (ti *TcpInput) Start() {
	for {
		select {
		case <-ti.exit:
			return
		default:
		}

		conn, err := ti.listener.Accept()
		if err != nil {
			continue
		}
		err = conn.SetReadDeadline(time.Now().Add(time.Duration(ti.Timeout) * time.Millisecond))
		if err != nil {
			logx.Errorf("set read deadline error: %v", err)
		}

		go func() {
			defer conn.Close()
			r := bufio.NewReader(conn)
			for {
				line, err := r.ReadString('\n')
				if err != nil {
					break
				}
				err = ti.handler.Consume("", line)
				if err != nil {
					logx.Errorf("consume error: %v", err)
				}
			}
		}()
	}
}

func (ti *TcpInput) Stop() {
	close(ti.exit)
	ti.listener.Close()
}
