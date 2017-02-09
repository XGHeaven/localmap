package server

import (
	"errors"
	"math/rand"
	"net"
	"time"

	"github.com/xgheaven/localmap/logger"
)

const (
	maxTry = 10
)

var (
	random = rand.New(rand.NewSource(time.Now().UnixNano()))
)

func (client *ServerClient) NewRandomListener(min, max int) (*net.TCPListener, uint16, error) {
	var (
		listener *net.TCPListener
		err      error
		try      = 0
		addr     = client.server.Addr
	)

	for ; try < maxTry; try++ {
		addr.Port = random.Intn(max-min) + min
		listener, err = net.ListenTCP("tcp", &addr)
		if err != nil {
			continue
		}
		logger.Infof("generator port at %d for client %s", addr.Port, client.UUID)
		return listener, uint16(addr.Port), nil
	}

	return nil, 0, errors.New("try out")
}

func Push2Queue(listener *net.TCPListener, queue chan *net.TCPConn) {
	for {
		conn, err := listener.AcceptTCP()
		if err != nil {
			logger.Error(err)
			if opErr, ok := err.(*net.OpError); ok && opErr.Op == "accept" {
				logger.Debug("client close")
				return
			}
			continue
		}
		logger.Debug("connected: ", conn.RemoteAddr(), "to", conn.LocalAddr())
		queue <- conn
	}
}
