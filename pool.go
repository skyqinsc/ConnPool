package cpool

import (
	"container/list"
	"time"
)

// ConnPool _
type ConnPool interface {
	SetIdleTimeout()
	SetMaxOpenConns(int)
	SetMaxIdleConns(int)
	Get(string) (Conn, error)
}

type cPool struct {
	dial        func(string, time.Duration) Conn
	idel        list.List
	maxIdel     int
	maxOpen     int
	maxLifetime time.Duration

	openCh  chan struct{}
	clearCh chan struct{}
}
