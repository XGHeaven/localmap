package client

import (
	"io"
	"net"
	"syscall"
	"time"

	"github.com/xgheaven/localmap/connect"
	"github.com/xgheaven/localmap/logger"
	"github.com/xgheaven/localmap/util"
)

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

	logger.Info("connect to server success, server version", helloReplyBlock.Version)
	logger.Infof("please use %s\n", client.DSAddr.String())

	go client.WaitEnd()

	for {
		client.SetReadDeadline(time.Now().Add(time.Minute))
		block, err := client.ReadBlock()
		if err == io.EOF {
			logger.Error("server EOF")
			client.ForceEnd()
			break
		}

		if e, ok := err.(net.Error); ok {
			if e.Timeout() {
				logger.Error(e)
				break
			}
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
			break
		case connect.HEART:
			client.WriteHeart()
			continue
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
