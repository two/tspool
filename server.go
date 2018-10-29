package tspool

import (
	"errors"
	"net"
)

var (
	errAddr          = errors.New("Addr must be given like 0.0.0.0:8080")
	errHandler       = errors.New("Handler must be given")
	errErrHandler    = errors.New("ErrHandler must be given")
	errRejectHandler = errors.New("RejectHandler must be given")
)

// Server have all functions that you should define to deal connect
type Server struct {
	Addr          string
	WorkerPool    WorkerPool
	Handler       func(net.Conn)
	ErrHandler    func(net.Conn, string)
	RejectHandler func(net.Conn, string)
}

// ListenAndServe will listen the port and serve your request
func ListenAndServe(srv *Server) error {
	if srv.Addr == "" {
		return errAddr
	}
	if srv.Handler == nil {
		return errHandler
	}
	if srv.ErrHandler == nil {
		return errErrHandler
	}
	if srv.RejectHandler == nil {
		return errRejectHandler
	}
	if srv.WorkerPool == nil {
		wp, err := DefaultWorkerPool()
		if err != nil {
			return err
		}
		srv.WorkerPool = wp
	}
	return srv.listenAndServe()
}

func (srv *Server) listenAndServe() error {
	l, err := net.Listen("tcp", srv.Addr)
	if err != nil {
		return err
	}
	for {
		conn, err := l.Accept()
		if err != nil {
			srv.ErrHandler(conn, err.Error())
			continue
		}
		worker, err := srv.WorkerPool.Get()
		if err != nil {
			srv.RejectHandler(conn, err.Error())
			continue
		}
		worker.Work(srv, conn)
	}
}
