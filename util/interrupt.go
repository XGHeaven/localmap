package util

import (
	"os"
	"os/signal"
	"syscall"
)

func NewInterruptChan(sig os.Signal) (c chan os.Signal) {
	c = make(chan os.Signal, 2)
	signal.Notify(c, syscall.SIGTERM, os.Interrupt)
	return
}
