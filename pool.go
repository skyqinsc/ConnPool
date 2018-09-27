package cpool

import (
	"container/list"
	"errors"
	"net"
	"sync"
	"time"
)

// CPool Err
var (
	ErrOverMaxConn = errors.New("cPool: over maximum of conn")
)

// ConnPool _
type ConnPool interface {
	SetIdleTime(time.Duration)
	SetMaxConns(uint32)
	Get() (Conn, error)
}

// an simple net pool
type cPool struct {
	// constructor of conn
	dial func() (net.Conn, error)
	// the container of idle conn
	idle *list.List
	// the number of current idel conn
	count uint32
	// maximum of conn
	maxConn uint32
	// idle time
	idleTime time.Duration

	locker sync.Mutex
}

// SetIdelTime _
func (c *cPool) SetIdleTime(timeout time.Duration) {
	c.idleTime = timeout
}

// SetMaxConns _
func (c *cPool) SetMaxConns(num uint32) {
	c.maxConn = num
}

func (c *cPool) Get() (Conn, error) {
	c.locker.Lock()
	if c.idle.Len() > 0 {
		ct := c.idle.Front()
		c.idle.Remove(ct)

		cn := ct.Value.(*conn)
		if now().Sub(cn.rTime) < c.idleTime {
			c.locker.Unlock()
			return cn, nil
		}
		cn.Close()
		c.count--
	}

	if c.count >= c.maxConn {
		return nil, ErrOverMaxConn
	}
	c.count++
	c.locker.Unlock()

	cn, err := c.dial()
	if err != nil {
		c.locker.Lock()
		c.count--
		c.locker.Unlock()
		return nil, err
	}

	ret := &conn{
		Conn:  cn,
		rTime: now(),
		pool:  c,
	}

	return ret, nil
}

// NewCPool create a net pool
func NewCPool(dial func() (net.Conn, error), maxConn uint32, timeout time.Duration) ConnPool {
	return &cPool{
		dial:     dial,
		idle:     list.New(),
		maxConn:  maxConn,
		idleTime: timeout,
	}
}
