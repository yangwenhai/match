package registry

import (
	"bytes"
	"encoding/binary"
	"errors"
	"kitfw/commom/amf"
	"time"

	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/sd/zk"
	stdzk "github.com/samuel/go-zookeeper/zk"
)

var (
	acl            = stdzk.WorldACL(stdzk.PermAll)
	connectTimeout = 3 * time.Second
	sessionTimeout = 10 * time.Second
)

//battle的zookeeper注册结构如下
/*
array(4) {
  ["host"]=>
  string(14) "192.168.10.193"
  ["port"]=>
  int(1234)
  ["weight"]=>
  int(30000)
  ["max_connect"]=>
  int(20)
}
*/
//和battle的结构保持一致，以便phpproxy访问
type ZkRegistry struct {
	Host       string `amf:"host"`
	Port       uint16 `amf:"port"`
	Weight     uint32 `amf:"weight"`
	MaxConnect uint32 `amf:"max_connect"`
}

type zkService struct {
	client zk.Client
}

func NewZkService(hosts []string) (*zkService, error) {
	c, err := zk.NewClient(
		hosts,
		log.NewNopLogger(),
		zk.ACL(acl),
		zk.ConnectTimeout(connectTimeout),
		zk.SessionTimeout(sessionTimeout),
	)
	if err != nil {
		return nil, err
	}
	if c == nil {
		return nil, errors.New("zookeeper client nil")
	}
	return &zkService{client: c}, nil
}

func (z *zkService) Register(path string, name string, host string, port uint16) error {

	re := &ZkRegistry{host, port, 30000, 20}
	buffer := bytes.NewBuffer(nil)
	err := amf.Encode(buffer, binary.BigEndian, re)
	if err != nil {
		return err
	}

	s := &zk.Service{
		Path: path,
		Name: name,
		Data: buffer.Bytes(),
	}

	if err := z.client.Register(s); err != nil {
		return err
	}

	//thers is a bug in go-kit-0.4.0 ,when we call Register,thers's no need to do "CreateParentNodes" and CreateProtectedEphemeralSequential will create a special path
	if err := z.client.Deregister(s); err != zk.ErrNodeNotFound && err != zk.ErrNotRegistered {
		return err
	}
	return nil
}
