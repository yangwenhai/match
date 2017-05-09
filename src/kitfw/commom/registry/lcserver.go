package registry

import (
	"bytes"
	"encoding/binary"
	"errors"
	"fmt"
	"kitfw/commom/amf"
	logger "kitfw/commom/log"
	"net"
	"strings"
	"sync"

	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/sd/zk"
	stdzk "github.com/samuel/go-zookeeper/zk"
)

type watchService struct {
	path     string
	client   zk.Client
	conMap   map[string]net.Conn
	groupMap map[string]ZkLcserver
	quitc    chan struct{}
	sync.RWMutex
}

type ZkLcserver struct {
	Host       string `amf:"host"`
	Port       uint16 `amf:"port"`
	WanHost    string `amf:"wan_host"`
	WanPort    uint16 `amf:"wan_port"`
	Db         string `amf:"db"`
	Weight     uint32 `amf:"weight"`
	MaxConnect uint32 `amf:"max_connect"`
}
type ZkGroup struct {
	Group string `amf:"group"`
}

var wservice *watchService

func InitWatchLcserver(zkhosts []string, path string) error {
	c, err := zk.NewClient(
		zkhosts,
		log.NewNopLogger(),
		zk.ACL(acl),
		zk.ConnectTimeout(connectTimeout),
		zk.SessionTimeout(sessionTimeout),
	)
	if err != nil {
		return err
	}
	if c == nil {
		return errors.New("zookeeper client nil")
	}
	wservice = &watchService{
		path:     path,
		client:   c,
		conMap:   make(map[string]net.Conn),
		groupMap: make(map[string]ZkLcserver),
		quitc:    make(chan struct{}),
	}

	if err = wservice.WatchLcserver(); err != nil {
		return err
	}
	return nil
}

func CallLcserver(group string, msg []byte) error {

	lc, ok := wservice.groupMap[group]
	if !ok {
		return errors.New(fmt.Sprintf("error group:%s", group))
	}

	con, ok := wservice.conMap[group]
	if !ok {
		addr := fmt.Sprintf("%s:%d", lc.Host, lc.Port)
		conn, err := net.Dial("tcp", addr)
		if err != nil {
			return err
		}
		wservice.Lock()
		wservice.conMap[group] = conn
		con = conn
		wservice.Unlock()
	}

	length := len(msg)
	for length > 0 {
		n, err := con.Write(msg)
		if err != nil {
			return err
		}
		length = length - n
	}

	return nil
}

func (w *watchService) WatchLcserver() error {
	instances, eventc, err := w.client.GetEntries(w.path)
	if err != nil {
		return err
	}

	if err := w.mapInstances(instances); err != nil {
		return err
	}

	go w.loop(eventc)

	return nil

}

func (w *watchService) loop(eventc <-chan stdzk.Event) {
	var (
		instances []string
		err       error
	)
	for {
		select {
		case <-eventc:
			// We received a path update notification. Call GetEntries to
			// retrieve child node data, and set a new watch, as ZK watches are
			// one-time triggers.
			instances, eventc, err = w.client.GetEntries(w.path)
			if err != nil {
				logger.Error("path", w.path, "msg", "failed to retrieve entries", "err", err)
				continue
			}
			logger.Info("path", w.path, "instances", len(instances))

			if err := w.mapInstances(instances); err != nil {
				logger.Error("err", err)
			}

		case <-w.quitc:
			return
		}
	}
}

func (w *watchService) Stop() {
	close(w.quitc)
}

func (w *watchService) mapInstances(instances []string) error {

	// str1 := "CgsBCWhvc3QGGzE5Mi4xNjguODguNjgJcG9ydAS/EQE="
	// d1, err := base64.StdEncoding.DecodeString(str1)
	// if err != nil {
	// 	fmt.Println(err)
	// }
	// var re struct {
	// 	Host string `amf:"host"`
	// 	Port int    `amf:"port"`
	// }
	// b1 := bytes.NewBuffer(d1)
	// if err := amf.Decode(b1, binary.BigEndian, &re); err != nil {
	// 	fmt.Println(err)
	// }
	// fmt.Println(re)

	// str2 := "CgsBC2dyb3VwBhdnYW1lMTgwMDAwMQE="
	// d2, err := base64.StdEncoding.DecodeString(str2)
	// if err != nil {
	// 	fmt.Println(err)
	// }
	// var mp struct {
	// 	Group string `amf:"group"`
	// }
	// b2 := bytes.NewBuffer(d2)
	// if err := amf.Decode(b2, binary.BigEndian, &mp); err != nil {
	// 	fmt.Println(err)
	// }
	// fmt.Println(mp)

	// encodeString := "CgsBCWhvc3QGHTE5Mi4xNjguMTAuMTkzCXBvcnQEvkERd2FuX2hvc3QGHzExNS4xODIuMjUxLjE5MxF3YW5fcG9ydAS+QQVkYgYbcGlyYXRlMTgwMDAwMQ13ZWlnaHQEgeowF21heF9jb25uZWN0BAIB"
	// decodeBytes, err := base64.StdEncoding.DecodeString(encodeString)
	// if err != nil {
	// 	fmt.Println(err)
	// }

	// var lc ZkLcserver
	// buffer := bytes.NewBuffer(decodeBytes)
	// if err := amf.Decode(buffer, binary.BigEndian, &lc); err != nil {
	// 	fmt.Println(err)
	// }
	// fmt.Println("xxx", lc)

	for _, v := range instances {
		var lc ZkLcserver
		buffer := bytes.NewBuffer([]byte(v))
		if err := amf.Decode(buffer, binary.BigEndian, &lc); err != nil {
			return err
		}
		//go-kit库的GetEntries 应该返回一个map，好让goup信息和lcserver信息一一对应，但他们的接口返回的只是lcserver信息，没法和group对应起来
		//暂时简单处理一下，先不考虑和服问题，
		group := strings.Replace(lc.Db, "pirate", "group", 1)
		w.groupMap[group] = lc
	}
	fmt.Println(w.groupMap)
	return nil
}
