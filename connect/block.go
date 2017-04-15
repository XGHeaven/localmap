package connect

import (
	"encoding/binary"
	"errors"
	"io"

	"github.com/xgheaven/localmap/logger"
)

const (
	HEL = iota
	HELRLY
	REQCON
	REQRLY
	CLOSE
)

const (
	BLOCKREADSIZE = 1024
)

type (
	BlockHeader struct {
		Type uint8
		Flag uint8
		Len  uint16
	}

	Block struct {
		*BlockHeader
		Data []byte
	}

	baseBlock struct {
		Len  uint32
		Raw  []byte
		Type byte
	}

	HelloBlock struct {
		*Block
	}

	HelloReplyBlcok struct {
		*Block
		Sport uint16
		Cport uint16
	}

	ReqConnBlock struct {
		*Block
	}

	CloseBlock struct {
		*Block
	}
)

func NewBlock(reader io.Reader) (*Block, error) {
	rawHeader := make([]byte, 4)
	n, err := reader.Read(rawHeader)
	logger.Debug(n, err == nil)
	if err != nil {
		return nil, err
	}
	if n != 4 {
		return nil, errors.New("wrong block header")
	}

	length := binary.LittleEndian.Uint16(rawHeader[2:4])
	typ := rawHeader[0]
	flag := rawHeader[1]
	base := &BlockHeader{Len: length, Flag: flag, Type: typ}

	if err != nil && length > 0 {
		return nil, err
	}

	data := make([]byte, length, length)

	for length > 0 {
		n, err := reader.Read(data[uint16(len(data))-length:])
		if uint16(n) > length {
			return nil, errors.New("block length error")
		}
		length -= uint16(n)

		if err != nil {
			return nil, err
		}
	}

	logger.Debug("Block: HEADER ", rawHeader, "BODY", data)
	return &Block{BlockHeader: base, Data: data}, nil
}

func NewHelloBlock(block *Block) (helloBlock *HelloBlock, err error) {
	helloBlock = &HelloBlock{Block: block}
	return
}

func NewHelloReplyBlock(block *Block) (helloReplyBlcok *HelloReplyBlcok, err error) {
	helloReplyBlcok = &HelloReplyBlcok{Block: block}
	helloReplyBlcok.Sport = binary.LittleEndian.Uint16(helloReplyBlcok.Data[:2])
	helloReplyBlcok.Cport = binary.LittleEndian.Uint16(helloReplyBlcok.Data[2:])
	return
}

func NewReqConnBlock(block *Block) (reqConnBlock *ReqConnBlock, err error) {
	reqConnBlock = &ReqConnBlock{Block: block}
	return
}

func NewCloseBlock(block *Block) (closeBlock *CloseBlock, err error) {
	closeBlock = &CloseBlock{Block: block}
	return
}
