package util

import (
  "net"
  "io"
  "../logger"
)

const (
  buffSize = 1024
)

func LinkConnect(dst, src net.Conn) {
  data := make([]byte, buffSize)
  for {
    n, err := src.Read(data)
    if n > 0 {
      dst.Write(data[:n])
    }
    if err == io.EOF {
      src.Close()
      dst.Close()
      logger.Debug(src.RemoteAddr().String() + " EOF")
      return
    }
  }
}
