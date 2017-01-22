package util

import (
	"io"
	"net"

	"../logger"
)

const (
	buffSize = 1024
)

func LinkConnect(dst, src *net.TCPConn) {
	data := make([]byte, buffSize)
	for {
		n, err := src.Read(data)
		if n > 0 {
			dst.Write(data[:n])
		}
		if err == io.EOF {
			src.CloseRead()
			dst.CloseWrite()
			logger.Debug(src.RemoteAddr().String() + " EOF")
			return
		}
	}
}
