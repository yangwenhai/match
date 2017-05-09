package tcp

import (
	"bytes"
	"context"
	"encoding/binary"
	"fmt"
	"kitfw/commom/amf"
	logger "kitfw/commom/log"
	"kitfw/commom/pb"
	"kitfw/match/envconf"
	"kitfw/match/protocol"
	"net"

	"github.com/go-kit/kit/endpoint"
)

type AmfRequest struct {
	Method   string   `amf:"method"`
	Args     []string `amf:"args"`
	Type     int      `amf:"type"`
	Time     int64    `amf:"time"` //amf的int只有29位，这里得用64位的数来代替
	Token    string   `amf:"token"`
	Callback string   `amf:"callback"`
	Return   bool     `amf:"return"`
}

type AmfResponse struct {
	Err string `amf:"err"`
	Ret string `amf:"ret"`
}

type client struct {
	conn     net.Conn
	endpoint endpoint.Endpoint
	request  AmfRequest
}

func RunServer(listener net.Listener, endpoint endpoint.Endpoint) error {
	for {
		conn, err := listener.Accept()
		if err != nil {
			return err
		}
		c := &client{conn: conn, endpoint: endpoint, request: AmfRequest{}}
		go c.Run()
	}
}

func (c *client) Run() {
	defer c.conn.Close()
	for {
		if err := c.recvMsg(); err != nil {
			logger.Error("err", err)
			break
		}
		if err := c.processMsg(); err != nil {
			logger.Error("err", err)
			break
		}
		c.clearRequest()
	}
}
func (c *client) recvMsg() error {

	//message length
	var length uint32
	err := binary.Read(c.conn, binary.BigEndian, &length)
	if nil != err {
		return fmt.Errorf("recvMsg err!failed when read request length, %v", err)
	}

	//flage
	var flage uint32
	err = binary.Read(c.conn, binary.BigEndian, &flage)
	if nil != err {
		return fmt.Errorf("recvMsg err!failed when read flage, %v", err)
	}

	//read message body
	data := make([]byte, length)
	err = binary.Read(c.conn, binary.BigEndian, data)
	if nil != err {
		return fmt.Errorf("recvMsg err!failed when read request data, %v", err)
	}

	//decode amf
	var re AmfRequest
	buffer := bytes.NewBuffer(data)
	if err := amf.Decode(buffer, binary.BigEndian, &re); err != nil {
		return fmt.Errorf("recvMsg err!decoad amf error, %v", err)
	}
	if re.Args == nil || len(re.Args) == 0 {
		return fmt.Errorf("recvMsg err!decoad amf error, args nil!method:%s token:%s", re.Method, re.Token)
	}
	c.request = re

	return nil
}

func (c *client) processMsg() error {

	//get protodid from method name
	protoid, ok := protocol.METHOD_PROTOCOL_MAP[c.request.Method]
	if !ok {
		return fmt.Errorf("processMsg err!invalid method%s token:%s", c.request.Method, c.request.Token)
	}
	logger.Info("encodetype", "amf", "logid", c.request.Token, "method", c.request.Method)

	//create logger
	rlogger := logger.NewLogger()
	rlogger.SetLogLevel(envconf.EnvCfg.LOG_LEVEL)
	rlogger.With("type", "tcpamf", "logid", c.request.Token)

	//create context
	ctx := context.Background()
	ctx = context.WithValue(ctx, "logger", rlogger)
	ctx = context.WithValue(ctx, "logid", c.request.Token)
	ctx = context.WithValue(ctx, "encodetype", "amf")

	//create response
	req := &pb.KitfwRequest{
		Protoid: protoid,
		Payload: []byte(c.request.Args[0]),
	}

	//call method
	res, err := c.endpoint(ctx, req)
	if err != nil {
		return err
	}
	if res == nil {
		return fmt.Errorf("processMsg err!response nil!method:%s token:%s", c.request.Method, c.request.Token)
	}

	//encode amf response
	response := &AmfResponse{Err: "ok", Ret: string(res.(*pb.KitfwReply).Payload[:])}
	buffer := bytes.NewBuffer(nil)
	err = amf.Encode(buffer, binary.BigEndian, response)
	if err != nil {
		return fmt.Errorf("processMsg err!encode amf response err! method:%s token:%s,%v", c.request.Method, c.request.Token, err)
	}
	msg := buffer.Bytes()

	//write message length
	var length uint32 = uint32(len(msg))
	err = binary.Write(c.conn, binary.BigEndian, length)
	if nil != err {
		return fmt.Errorf("processMsg err!write response length err!method:%s token:%s,%v", c.request.Method, c.request.Token, err)
	}

	//write message flage
	var flage uint32 = 0
	err = binary.Write(c.conn, binary.BigEndian, flage)
	if nil != err {
		return fmt.Errorf("processMsg err!write response flage err!method:%s token:%s,%v", c.request.Method, c.request.Token, err)
	}

	//write response
	err = binary.Write(c.conn, binary.BigEndian, msg)
	if err != nil {
		return fmt.Errorf("processMsg err!send response err!method:%s token:%s,%v", c.request.Method, c.request.Token, err)
	}

	return nil
}

func (c *client) clearRequest() {
	c.request = AmfRequest{}

}
