package client

import (
	"fmt"
	"io"
	"net"
	"syscall"

	"github.com/xgheaven/localmap/conn"
	"github.com/xgheaven/localmap/connect"
	"github.com/xgheaven/localmap/logger"
	"github.com/xgheaven/localmap/util"
)

var Option ClientOption

type (
	Client struct {
		*ClientOption
		DSAddr net.TCPAddr
		DCAddr net.TCPAddr
		*connect.TCPConnect
	}
)

func NewClient(option *ClientOption) (client *Client) {
	client = &Client{}
	client.ClientOption = option
	return
}

func (client *Client) Start() {
	cConn, err := net.DialTCP("tcp", nil, &client.SAddr)
	if err != nil {
		logger.Fatal("connect server error")
	}

	client.TCPConnect = (*connect.TCPConnect)(cConn)

	client.WriteHello()
	helloReplyBlock, err := client.ReadHelloReply()
	if err != nil {
		logger.Fatal("regnize error")
	}

	client.DSAddr = client.SAddr
	client.DCAddr = client.SAddr

	client.DSAddr.Port = int(helloReplyBlock.Sport)
	client.DCAddr.Port = int(helloReplyBlock.Cport)

	logger.Infof("connect to server, please use %s\n", client.DSAddr.String())

	go client.WaitEnd()

loop:
	for {
		block, err := client.ReadBlock()
		if err == io.EOF {
			logger.Error("server EOF")
			client.ForceEnd()
			break
		}
		if err != nil {
			logger.Error(err)
			continue
		}
		switch block.Type {
		case connect.REQCON:
			go func() {
				sConn, sErr := net.DialTCP("tcp", nil, &client.DCAddr)

				if sErr != nil {
					logger.Error("connect to server error")
					return
				}

				logger.Debug("connect to server")

				cConn, cErr := net.DialTCP("tcp", nil, &client.CAddr)

				if cErr != nil {
					logger.Error("connect to client error")
					return
				}

				logger.Debug("connect to client")

				go util.LinkConnect(cConn, sConn)
				go util.LinkConnect(sConn, cConn)
			}()
		case connect.CLOSE:
			client.EndSelf()
			break loop
		}
	}

	client.Close()
	logger.Info("Bye!")
}

func (client *Client) End() {
	client.WriteClose()
	client.ReadClose()
	client.EndSelf()
}

func (client *Client) ForceEnd() {
	client.EndSelf()
}

func (client *Client) EndSelf() {
}

func (client *Client) WaitEnd() {
	<-util.NewInterruptChan(syscall.SIGTERM)
	logger.Debug("receive SIGTERM sign")
	client.WriteClose()
}

func Start() {
	var (
		sPort, cPort int
	)
	cConn, err := net.DialTCP("tcp", nil, &Option.SAddr)
	if err != nil {
		logger.Fatal("connect server error")
	}

	defer cConn.Close()

	connect := conn.NewConnect(cConn)
	connect.WriteString("hello\n")
	connect.Flush()
	data, err := connect.ReadWholeLine()
	if err != nil {
		return
	}
	fmt.Sscanf(string(data), "hello %d %d", &sPort, &cPort)
	sAddr := Option.SAddr
	cAddr := Option.SAddr
	sAddr.Port = sPort
	cAddr.Port = cPort
	logger.Infof("connect to server, please use %s\n", sAddr.String())

	go func() {
		<-util.NewInterruptChan(syscall.SIGTERM)
		connect.WriteString("bye\n")
		connect.Flush()
		logger.Info("waiting for server close port")
		data, err := connect.ReadWholeLine()
		if string(data) != "bye" || err != nil {
			logger.Error("disconnect server port error")
		} else {
			logger.Info("server close port success")
		}
		logger.Info("client close")
		cConn.Close()
	}()

	for {
		data, err := connect.ReadWholeLine()
		if err == io.EOF {
			break
		}
		if err != nil {
			continue
		}
		if string(data) == "request" {
			go func() {
				sConn, sErr := net.DialTCP("tcp", nil, &cAddr)

				if sErr != nil {
					logger.Error("connect to server error")
					return
				}

				logger.Debug("connect to server")

				cConn, cErr := net.DialTCP("tcp", nil, &Option.CAddr)

				if cErr != nil {
					logger.Error("connect to client error")
					return
				}

				logger.Debug("connect to client")

				go util.LinkConnect(cConn, sConn)
				go util.LinkConnect(sConn, cConn)
			}()
		}
	}
}
