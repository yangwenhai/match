package phpproxy

import (
	"bytes"
	"encoding/binary"
	"kitfw/commom/amf"
	"time"
)

/*
/rpcfw/lib/PHPProxy.class.php :  __call
phpproxy ->> other serverice ( battle、lcserver)
array ('method' => $method, 'args' => $arrArg,
				'type' => $this->requestType, 'time' => time (), 'token' => $this->token,
				'callback' => $method, 'return' => $this->dummyReturn );
*/
type PHPProxy2Service struct {
	Method   string        `amf:"method"`
	Args     []interface{} `amf:"args"`
	Type     int           `amf:"type"`
	Time     int64         `amf:"time"` //amf的int只有29位，这里得用64位的数来代替
	Token    string        `amf:"token"`
	Callback string        `amf:"callback"`
	Return   bool          `amf:"return"`
}

/*
/rpcfw/lib/PHPProxy.class.php :  __call
php ->> phpproxy
$arrRequest = array ('token' => $this->token, 'method' => 'proxy', 'group' => $group,
				'db' => $db, 'args' => array ($this->module, $request ) );
*/
type Msg2PHPProxy struct {
	Token  string   `amf:"token"`
	Method string   `amf:"method"`
	Group  string   `amf:"group"`
	Db     string   `amf:"db"`
	Args   []string `amf:"args"`
}

/*
ServerProxy.class.php:asyncExecuteRequest
 ->> lcserver ->> php

array('method' => $method, 'args' => $arrArg, 'token' => $token,
'callback' => array('callbackName' => $callback),'serverId' => $this->serverId);
*/
type CallBack struct {
	CallbackName string `amf:"callbackName"`
}
type ExecuteTaskMsg struct {
	Method   string   `amf:"method"`
	Args     []string `amf:"args"`
	Token    string   `amf:"token"`
	Callback CallBack `amf:"callback"`
	ServerId int64    `amf:"serverId"`
}

/*
msg ->> phpproxy ->> lcserver ->> php
*/
func CreatExecuteTaskMsg(userid int64, group string, db string, serverid int64, token string, method string, args []string) ([]byte, error) {

	et := ExecuteTaskMsg{
		Method:   method,
		Args:     args, //[]string
		Token:    token,
		Callback: CallBack{CallbackName: "dummy"},
		ServerId: serverid,
	}

	ps := &PHPProxy2Service{
		Method:   "asyncExecuteRequest",
		Args:     []interface{}{userid, et, false},
		Type:     1, // RequestType::RELEASE
		Time:     int64(time.Now().Unix()),
		Token:    token,
		Callback: "asyncExecuteRequest",
		Return:   true, //dummyReturn
	}

	psBuffer := bytes.NewBuffer(nil)
	err := amf.Encode(psBuffer, binary.BigEndian, ps)
	if err != nil {
		return nil, err
	}
	paylaod := psBuffer.Bytes()

	//write message length
	headerlen := make([]byte, 4)
	binary.BigEndian.PutUint32(headerlen, uint32(len(paylaod)))

	//write message flage
	headerflage := make([]byte, 4)
	binary.BigEndian.PutUint32(headerflage, 0)

	buf := new(bytes.Buffer)
	buf.Write(headerlen)
	buf.Write(headerflage)
	buf.Write(paylaod)
	encodemsg := buf.Bytes()

	m := &Msg2PHPProxy{
		Token:  token,
		Method: "proxy",
		Group:  group,
		Db:     db,
		Args:   []string{"lcserver", string(encodemsg[:])},
	}

	buffer := bytes.NewBuffer(nil)
	err = amf.Encode(buffer, binary.BigEndian, m)
	if err != nil {
		return nil, err
	}

	return buffer.Bytes(), nil
}

/*
	php method:  RPCContext::getInstance()->executeTask(....)
*/
func ExecuteTask(userid int64, group string, db string, serverid int64, token string, method string, args []string) error {

	msg, err := CreatExecuteTaskMsg(
		userid,
		group,
		db,
		serverid,
		token,
		method,
		args,
	)

	if err != nil {
		return err
	}

	//chose a connection from phpprox pool and call phpproxy
	ret, err := Call(msg)
	if err != nil {
		return err
	}
	_ = ret
	return err
}
