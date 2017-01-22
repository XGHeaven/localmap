package client

import (
	"net"

	"../conn"
	"../util"
	// "log"
	"fmt"
	"io"
	"syscall"

	"../logger"
)

var Option ClientOption

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
