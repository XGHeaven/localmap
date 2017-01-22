package server

import (
	"net"
)

type ServerOption struct {
	Addr net.TCPAddr
}
