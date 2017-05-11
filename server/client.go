package server

import (
	"io"
	"net"
	"time"

	"github.com/google/uuid"
	"github.com/xgheaven/localmap/connect"
	"github.com/xgheaven/localmap/logger"
	"github.com/xgheaven/localmap/util"
)

const (
	INIT = iota
	CONNECTION
	ESTABLISH
	CLOSING
	END
)

type ServerClient struct {
	*connect.TCPConnect
	server       *Server
	UUID         uuid.UUID
	State        int
	Sport, Cport uint16
	Slsn, Clsn   *net.TCPListener
	Schan, Cchan chan *net.TCPConn
	heart        *time.Ticker
}

func NewServerClient(server *Server, clientConn *connect.TCPConnect) *ServerClient {
	client := &ServerClient{}
	client.TCPConnect = clientConn
	client.UUID, _ = uuid.NewUUID()
	client.server = server
	return client
}

func (client *ServerClient) Start() {
	// set timeout 10s
	client.SetReadDeadline(time.Now().Add(time.Second * 10))
	helloBlock, err := client.ReadHello()
	if err != nil {
		if e, ok := err.(net.Error); ok {
			if e.Timeout() {
				logger.Error("client connect timeout")
			}
		}
		logger.Error(err)
		client.Close()
		return
	}

	logger.Info("client connect success, version", helloBlock.Version)

	// remove timeout
	client.SetReadDeadline(time.Time{})

	sLsn, sPort, sErr := client.NewRandomListener(25000, 30000)

	if sErr != nil {
		logger.Error("listen public port error")
		client.Close()
		return
	}

	cLsn, cPort, cErr := client.NewRandomListener(30000, 35000)

	if cErr != nil {
		logger.Error("listen private port error")
		client.Close()
		return
	}

	client.Slsn = sLsn
	client.Sport = sPort
	client.Clsn = cLsn
	client.Cport = cPort
	client.Schan = make(chan *net.TCPConn, 100)
	client.Cchan = make(chan *net.TCPConn, 100)

	client.WriteHelloReply(sPort, cPort)

	go Push2Queue(client.Slsn, client.Schan)
	go Push2Queue(client.Clsn, client.Cchan)

	// handle connection
	go func() {
		for {
			sConn, sOk := <-client.Schan
			// check schan close
			if !sOk {
				logger.Warning("server channel close")
				break
			}
			client.WriteRequestConnect()
			cConn, cOk := <-client.Cchan
			if !cOk {
				logger.Warning("client channel close")
				break
			}
			go util.LinkConnect(sConn, cConn)
			go util.LinkConnect(cConn, sConn)
		}
	}()

	client.heart = time.NewTicker(time.Second * 55)

	// send heart
	go func() {
		for range client.heart.C {
			_, err := client.WriteHeart()

			if err != nil {
				logger.Info("Heart Close", err)
				break
			}
		}
	}()

loop:
	for {
		client.SetReadDeadline(time.Now().Add(time.Minute))
		block, err := client.ReadBlock()
		if err == io.EOF {
			logger.Error("Client:", client.UUID, "EOF")
			client.EndSelf()
			break
		}

		if e, ok := err.(net.Error); ok {
			if e.Timeout() {
				logger.Error(e)
				break
			}
		}

		if err != nil {
			logger.Error("Client:", client.UUID, err)
			continue
		}

		switch block.Type {
		case connect.CLOSE:
			client.WriteClose()
			client.EndSelf()
			break loop
		case connect.HEART:
			continue
		}
	}

	client.Close()
	logger.Info("Client:", client.UUID, "Closed")
}

func (client *ServerClient) End() {
	client.WriteClose()
	client.ReadClose()
	client.EndSelf()
}

func (client *ServerClient) EndSelf() {
	//delete(client.server.clients, client.UUID.String())
	client.heart.Stop()
	client.Slsn.Close()
	logger.Info("Client:", client.UUID, client.Sport, "Closed")
	client.Clsn.Close()
	logger.Info("Client:", client.UUID, client.Cport, "Closed")
	client.server = nil
}
