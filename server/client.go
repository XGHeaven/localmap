package server

import (
	"io"
	"net"

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
}

func NewServerClient(server *Server, clientConn *connect.TCPConnect) *ServerClient {
	client := &ServerClient{}
	client.TCPConnect = clientConn
	client.UUID, _ = uuid.NewUUID()
	client.server = server
	return client
}

func (client *ServerClient) Start() {
	_, err := client.ReadHello()
	if err != nil {
		client.Close()
		return
	}

	sLsn, sPort, sErr := client.NewRandomListener(25000, 30000)

	if sErr != nil {
		logger.Error("Listen Server Error")
		client.Close()
		return
	}

	cLsn, cPort, cErr := client.NewRandomListener(30000, 35000)

	if cErr != nil {
		logger.Error("Listen Client Error")
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

	go func() {
		for {
			sConn := <-client.Schan
			client.WriteRequestConnect()
			cConn := <-client.Cchan
			go util.LinkConnect(sConn, cConn)
			go util.LinkConnect(cConn, sConn)
		}
	}()

loop:
	for {
		block, err := client.ReadBlock()
		if err == io.EOF {
			logger.Error("Client:", client.UUID, "EOF")
			client.EndSelf()
			break
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
	delete(client.server.clients, client.UUID.String())
	client.Slsn.Close()
	logger.Info("Client:", client.UUID, client.Sport, "Closed")
	client.Clsn.Close()
	logger.Info("Client:", client.UUID, client.Cport, "Closed")
	client.server = nil
}
