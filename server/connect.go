package server

import (
  "net"
  // "encoding/base64"
  // "log"
  "../conn"
  "math/rand"
  "strconv"
  // "io"
  "../util"
  "../logger"
)

func HandleConnect(c net.Conn, random *rand.Rand) {
  var (
    // port int
    // sListener net.Listener
    // err error
    // trys int
  )

  connect := conn.NewConnect(c)

  data, err := connect.ReadWholeLine()
  // conn.SetReadDeadline(time.Now().Add(3e9))
  // n, err := conn.Read(data)
  // log.Println(err, n, data)
  if string(data) != "hello" || err != nil {
    c.Write([]byte("error"))
    c.Close()
    logger.Error("transform conn error")
    return
  }

  // for ; trys < 10; trys++ {
  //   port = random.Intn(100) + 30000
  //   sListener, err = net.Listen("tcp", ":" + strconv.Itoa(port))
  //   if err != nil {
  //     continue
  //   }
  //   break
  // }
  // if trys > 10 {
  //   log.Fatal("generator port error")
  // }
  // log.Printf("generator port at %d and listen at this\n", port)

  sListener, sPort, sErr := NewRandomListener(30000, 35000)

  if sErr != nil {
    connect.WriteString("listen server port error")
    connect.Flush()
    c.Close()
    logger.Error("can't listen server port")
    return
  }

  cListener, cPort, cErr := NewRandomListener(25000, 30000)

  if cErr != nil {
    connect.WriteString("listen client port error")
    connect.Flush()
    c.Close()
    logger.Error("can't listen client port")
    return
  }

  defer sListener.Close()
  defer cListener.Close()

  go func() {
    data, err := connect.ReadWholeLine()
    if (string(data) == "bye" || err != nil) {
      sListener.Close()
      cListener.Close()
      connect.WriteString("bye\n")
      connect.Flush()
    }
  }()

  sQueue := make(chan net.Conn, 10)
  cQueue := make(chan net.Conn, 10)

  go Push2Queue(cListener, cQueue)
  go Push2Queue(sListener, sQueue)

  connect.WriteString("hello " + strconv.Itoa(sPort) + " " + strconv.Itoa(cPort) + "\n")
  connect.Flush()

  for {
    sConn := <- sQueue
    connect.WriteString("request\n")
    connect.Flush()
    cConn := <- cQueue
    // go io.Copy(sConn, cConn)
    // go io.Copy(cConn, sConn)
    go util.LinkConnect(sConn, cConn)
    go util.LinkConnect(cConn, sConn)
  }
}
