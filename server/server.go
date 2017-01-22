package server

import (
  "net"
  "math/rand"
  "time"
  "strconv"
  "../logger"
)

type ServerOption struct {
  Cport int
}

func Start(option ServerOption) {
  listener, err := net.Listen("tcp", ":" + strconv.Itoa(option.Cport))

  if err != nil {
    logger.Fatal(err)
  }

  logger.Info("listen at " + strconv.Itoa(option.Cport))

  random := rand.New(rand.NewSource(time.Now().UnixNano()))
  // ports := make(map[int][net.Conn])
  for {
    c, err := listener.Accept()
    if err != nil {
      logger.Error(err)
      continue
    }
    logger.Debug("accept new client")
    go HandleConnect(c, random)
  }
}
