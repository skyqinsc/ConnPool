package cpool

import (
	"log"
	"net"
	"time"
)

// Conn _
type Conn interface {
	net.Conn
	Put()
}

type conn struct {
	net.Conn
	rTime time.Time
	pool  *cPool
}

func (c *conn) Put() {
	c.pool.locker.Lock()
	defer c.pool.locker.Unlock()

	c.rTime = now()
	log.Println(c.pool.idle.Len())
	c.pool.idle.PushBack(c)
}
