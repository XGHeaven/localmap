package client

import (
  "net"
  "strconv"
  "../conn"
  "../util"
  // "log"
  "fmt"
  "io"
  "syscall"
  "../logger"
)

type ClientOption struct {
  Addr string
  Sport int
  Cport int
}

func Start(option ClientOption) {
  var (
    sPort, cPort int
  )
  cConn, err := net.Dial("tcp", option.Addr + ":" + strconv.Itoa(option.Sport))
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
  logger.Infof("connect to server, please use %s:%d\n", option.Addr, sPort)

  go func() {
    <-util.NewInterruptChan(syscall.SIGTERM)
    connect.WriteString("bye\n")
    connect.Flush()
    logger.Info("waiting for server close port")
    data, err := connect.ReadWholeLine()
    if (string(data) != "bye" || err != nil) {
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
      break;
    }
    if err != nil {
      continue
    }
    if string(data) == "request" {
      go func() {
        sConn, sErr := net.Dial("tcp", option.Addr + ":" + strconv.Itoa(cPort))

        if sErr != nil {
          logger.Error("connect to server error")
          return
        }

        cConn, cErr := net.Dial("tcp", option.Addr + ":" + strconv.Itoa(option.Cport))

        if cErr != nil {
          logger.Error("connect to client error")
          return
        }

        go util.LinkConnect(cConn, sConn)
        go util.LinkConnect(sConn, cConn)
      }()
    }
  }
}
