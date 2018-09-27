package cpool

import (
	"log"
	"net"
	"testing"
	"time"
)

func server() {
	tcpListener, err := net.Listen("tcp", "127.0.0.1:8123")
	if err != nil {
		panic("failed to listen")
	}
	defer tcpListener.Close()

	log.Println("Listen on :8123")
	for {
		conn, err := tcpListener.Accept()
		if err != nil {
			log.Panicln("failed to accept:", err)
		} else {
			log.Println(conn.RemoteAddr().String(), " tcp connect success")
		}
	}
}

func Test_Get(t *testing.T) {
	go server()

	dial := func() (net.Conn, error) {
		return net.DialTimeout("tcp", "127.0.0.1:8123", 1*time.Second)
	}
	p := NewCPool(dial, 4, 100*time.Millisecond)

	var clist []Conn
	for i := 0; i < 4; i++ {
		c, err := p.Get()
		if err != nil {
			if err != ErrOverMaxConn {
				t.Error(err)
			}
			t.Log(err)
		}
		clist = append(clist, c)
	}
}
