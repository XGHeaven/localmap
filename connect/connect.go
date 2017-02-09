package connect

import (
	"encoding/binary"
	"errors"
	"net"

	"github.com/xgheaven/localmap/logger"
)

/*
 block
 |type:4, none: 5|other...|
 hello block 1b
 |0000xxxx|
 hello reply block 5b
 |0001xxxx|server port||client port||
 request new connect 1b
 |0010xxxx|
 request new reply connect
*/

type TCPConnect net.TCPConn

func (conn *TCPConnect) ReadHello() (*HelloBlock, error) {
	block, err := conn.ReadBlock()
	if err != nil {
		return nil, err
	}
	if block.Type != HEL {
		return nil, errors.New("not hello block")
	}
	helloBlock, _ := NewHelloBlock(block)
	return helloBlock, nil
}

func (conn *TCPConnect) WriteHello() {
	block := &Block{BlockHeader: &BlockHeader{Type: HEL}}
	conn.WriteBlock(block)
}

func (conn *TCPConnect) ReadHelloReply() (*HelloReplyBlcok, error) {
	block, err := conn.ReadBlock()
	if err != nil {
		return nil, err
	}
	if block.Type != HELRLY {
		return nil, errors.New("not hello reply block")
	}
	helloReplyBlcok, _ := NewHelloReplyBlock(block)
	return helloReplyBlcok, nil
}

func (conn *TCPConnect) WriteHelloReply(sPort, cPort uint16) {
	data := make([]byte, 4)
	binary.LittleEndian.PutUint16(data[:2], sPort)
	binary.LittleEndian.PutUint16(data[2:], cPort)
	block := &Block{BlockHeader: &BlockHeader{Type: HELRLY}}
	block.Data = data
	conn.WriteBlock(block)
}

func (conn *TCPConnect) ReadRequestConnect() (*ReqConnBlock, error) {
	block, err := conn.ReadBlock()
	if err != nil {
		return nil, err
	}
	if block.Type != REQCON {
		return nil, errors.New("not request connect")
	}
	reqConnBlock, _ := NewReqConnBlock(block)
	return reqConnBlock, nil
}

func (conn *TCPConnect) WriteRequestConnect() {
	block := &Block{BlockHeader: &BlockHeader{Type: REQCON}}
	conn.WriteBlock(block)
}

func (conn *TCPConnect) ReadClose() (*CloseBlock, error) {
	block, err := conn.ReadBlock()
	if err != nil {
		return nil, err
	}
	if block.Type != CLOSE {
		return nil, errors.New("not close block")
	}
	closeBlock, _ := NewCloseBlock(block)
	return closeBlock, nil
}

func (conn *TCPConnect) WriteClose() {
	block := &Block{BlockHeader: &BlockHeader{Type: CLOSE}}
	conn.WriteBlock(block)
}

func (conn *TCPConnect) ReadBlock() (*Block, error) {
	return NewBlock(conn)
}

func (conn *TCPConnect) WriteBlock(block *Block) {
	if block.Data != nil {
		block.Len = uint16(len(block.Data))
	} else {
		block.Len = 0
	}
	data := make([]byte, 4, block.Len+4)
	data[0] = block.Type
	data[1] = block.Flag
	binary.LittleEndian.PutUint16(data[2:4], block.Len)
	data = append(data, block.Data...)
	logger.Debug("Block: WRITE", data)
	conn.Write(data)
}
