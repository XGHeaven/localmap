package main

import (
	"flag"
	"io"
	"net"
	"errors"
	"./server"
	"./client"
	"./logger"
)

func getMessage(conn net.Conn) {
	io.Copy(conn, conn)
}

var (
	isServer bool
	isClient bool
	sPort int
	cPort int
	sAddr string
)

func init() {
	flag.BoolVar(&isServer, "server", false, "-server")
	flag.BoolVar(&isClient, "client", false, "-client")
	flag.IntVar(&sPort, "sport", 8000, "-sport=8000")
	flag.IntVar(&cPort, "cport", 8080, "-cport=8080")
	flag.StringVar(&sAddr, "addr", "127.0.0.1", "-addr=127.0.0.1")
	flag.Parse()
	if !isClient {
		isServer = true
		isClient = false
		err := checkServerOption()
		if err != nil {
			logger.Fatal(err)
		}
		logger.Info("start as server")
	} else {
		isServer = false
		isClient = true
		err := checkClientOption()
		if err != nil {
			logger.Fatal(err)
		}
		logger.Info("start as client")
	}
}

func checkServerOption() error {
	return nil
}

func checkClientOption() error {
	if sAddr == "" {
		return errors.New("empty server addr")
	}
	return nil
}

func main() {
	defer func() {
		err := recover()
		if err != nil {
			logger.Error(err)
		}
	}()
	if isServer {
		option := server.ServerOption{Cport: sPort}
		server.Start(option)
	}
	if isClient {
		option := client.ClientOption{Addr:sAddr, Cport: cPort, Sport: sPort}
		client.Start(option)
	}
}
