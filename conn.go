package cpool

import (
	"net"
	"time"
)

// Conn _
type Conn interface {
	net.Conn
}

type cPoolConn struct {
	c        net.Conn
	createAt time.Time
}
