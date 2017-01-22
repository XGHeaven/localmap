package util

import(
  "os/signal"
  "os"
  "syscall"
)

func  NewInterruptChan(sig os.Signal) (c chan os.Signal) {
  c = make(chan os.Signal, 5)
  signal.Notify(c, syscall.SIGTERM)
  return
}
