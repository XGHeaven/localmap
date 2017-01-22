package util

import (
  "bufio"
)

func ReadWholeLine(conn *bufio.Reader) ([]byte, error) {
  var (
    result = []byte{}
    isPrefex = true
    err error
    line []byte
  )

  for isPrefex {
    line, isPrefex, err = conn.ReadLine()
    result = append(result, line...)
    if err != nil {
      return result, err
    }
  }

  return result, nil
}
