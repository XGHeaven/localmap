package conn

import (
  "net"
  "bufio"
)

type Connect struct {
  *bufio.ReadWriter
}

func NewConnect(conn net.Conn) *Connect {
  reader := bufio.NewReader(conn)
  writer := bufio.NewWriter(conn)
  return &Connect{ReadWriter: bufio.NewReadWriter(reader, writer)}
}

func (c *Connect) ReadWholeLine() ([]byte, error) {
  var (
    result = []byte{}
    isPrefex = true
    err error
    line []byte
  )

  for isPrefex {
    line, isPrefex, err = c.ReadLine()
    result = append(result, line...)
    if err != nil {
      return result, err
    }
  }
  return result, nil
}
