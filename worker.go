package tspool

import (
	"net"
)

// Worker deal the request
type Worker interface {
	Work(*Server, net.Conn)
}

type defaultWorker struct {
	conn    chan net.Conn
	pos     uint
	server  *Server
	handler func(net.Conn)
}

// Work handler the request from client
func (w *defaultWorker) Work(srv *Server, conn net.Conn) {
	w.server = srv
	w.handler = srv.Handler
	w.conn <- conn
}

func (w *defaultWorker) run() {
	go func() {
		for {
			select {
			case c := <-w.conn:
				w.handler(c)
				w.server.WorkerPool.Put(w)
			}
		}
	}()
}

func newWorker(i uint) *defaultWorker {
	return &defaultWorker{
		conn: make(chan net.Conn),
		pos:  i,
	}
}
