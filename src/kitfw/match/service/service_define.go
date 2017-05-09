package service

import (
	"bytes"
	"context"
	"encoding/binary"
	"kitfw/commom/amf"
	protocol "kitfw/match/protocol"
)

type BaseHandler interface {
	Process(context.Context, []byte) ([]byte, error)
	doProcess(ctx context.Context)
}

type NewHandlerFunc func() BaseHandler

var HandlerMap = map[int32]NewHandlerFunc{
	protocol.PROTOCOL_SUM_REQUEST:    func() BaseHandler { return NewSumHandler() },
	protocol.PROTOCOL_CONCAT_REQUEST: func() BaseHandler { return NewConcatHandler() },
	protocol.PROTOCOL_MATCH_REQUEST:  func() BaseHandler { return NewMatchHandler() },
}

func GetHandler(ProtoId int32) BaseHandler {
	if HandlerMap[ProtoId] != nil {
		return HandlerMap[ProtoId]()
	}
	return nil
}

type SumHandler struct {
	request *protocol.SumRequest
	reply   *protocol.SumReply
}

func NewSumHandler() *SumHandler {
	request := &protocol.SumRequest{}
	reply := &protocol.SumReply{}
	return &SumHandler{request, reply}
}

func (handler *SumHandler) Process(ctx context.Context, payload []byte) (ret []byte, err error) {

	encodetype := ctx.Value("encodetype").(string)
	if encodetype == "capnp" {
		if err = protocol.Decode(handler.request, payload); err != nil {
			return nil, err
		}
	} else if encodetype == "amf" {
		buffer := bytes.NewBuffer(payload)
		if err := amf.Decode(buffer, binary.BigEndian, handler.request); err != nil {
			return nil, err
		}
	}

	handler.doProcess(ctx)

	if encodetype == "capnp" {
		ret, err = protocol.Encode(handler.reply)
		if err != nil {
			return nil, err
		}
	} else if encodetype == "amf" {
		buffer := bytes.NewBuffer(nil)
		err = amf.Encode(buffer, binary.BigEndian, &handler.reply)
		if err != nil {
			return nil, err
		}
		ret = buffer.Bytes()
	}
	return
}

type ConcatHandler struct {
	request *protocol.ConcatRequest
	reply   *protocol.ConcatReply
}

func NewConcatHandler() *ConcatHandler {
	request := &protocol.ConcatRequest{}
	reply := &protocol.ConcatReply{}
	return &ConcatHandler{request, reply}
}

func (handler *ConcatHandler) Process(ctx context.Context, payload []byte) (ret []byte, err error) {

	encodetype := ctx.Value("encodetype").(string)
	if encodetype == "capnp" {
		if err = protocol.Decode(handler.request, payload); err != nil {
			return nil, err
		}
	} else if encodetype == "amf" {
		buffer := bytes.NewBuffer(payload)
		if err := amf.Decode(buffer, binary.BigEndian, handler.request); err != nil {
			return nil, err
		}
	}

	handler.doProcess(ctx)

	if encodetype == "capnp" {
		ret, err = protocol.Encode(handler.reply)
		if err != nil {
			return nil, err
		}
	} else if encodetype == "amf" {
		buffer := bytes.NewBuffer(nil)
		err = amf.Encode(buffer, binary.BigEndian, &handler.reply)
		if err != nil {
			return nil, err
		}
		ret = buffer.Bytes()
	}

	return
}

type MatchHandler struct {
	request *protocol.MatchRequest
	reply   *protocol.MatchReply
}

func NewMatchHandler() *MatchHandler {
	request := &protocol.MatchRequest{}
	reply := &protocol.MatchReply{}
	return &MatchHandler{request, reply}
}

func (handler *MatchHandler) Process(ctx context.Context, payload []byte) (ret []byte, err error) {

	encodetype := ctx.Value("encodetype").(string)
	if encodetype == "capnp" {
		if err = protocol.Decode(handler.request, payload); err != nil {
			return nil, err
		}
	} else if encodetype == "amf" {
		buffer := bytes.NewBuffer(payload)
		if err := amf.Decode(buffer, binary.BigEndian, handler.request); err != nil {
			return nil, err
		}
	}

	handler.doProcess(ctx)

	if encodetype == "capnp" {
		ret, err = protocol.Encode(handler.reply)
		if err != nil {
			return nil, err
		}
	} else if encodetype == "amf" {
		buffer := bytes.NewBuffer(nil)
		err = amf.Encode(buffer, binary.BigEndian, &handler.reply)
		if err != nil {
			return nil, err
		}
		ret = buffer.Bytes()
	}

	return
}
