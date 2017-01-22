package client

import (
	"net"
)

type ClientOption struct {
	SAddr net.TCPAddr
	CAddr net.TCPAddr
}
