package server

import (
  "net"
  "math/rand"
  "time"
  "strconv"
  // "log"
  "errors"
  "../logger"
)

const (
  maxTry = 10
)

var (
  random = rand.New(rand.NewSource(time.Now().UnixNano()))
)

func NewRandomListener(min, max int) (net.Listener, int, error) {
  var (
    port int
    listener net.Listener
    err error
    try = 0
  )

  for ;try < maxTry; try++ {
    port = random.Intn(max-min) + min
    listener, err = net.Listen("tcp", ":" + strconv.Itoa(port))
    if err != nil {
      continue
    }
    logger.Infof("generator port at %d and listen at this", port)
    return listener, port, nil
  }

  return nil, 0, errors.New("try out")
}

func Push2Queue(listener net.Listener, queue chan net.Conn) {
  for {
    conn, err := listener.Accept()
    if err != nil {
      logger.Error(err)
      if opErr, ok := err.(*net.OpError); ok && opErr.Op == "accept" {
        logger.Debug("client close")
        return
      }
      continue
    }
    logger.Debug("connected: ", conn.RemoteAddr(), conn.LocalAddr())
    queue <- conn
  }
}
