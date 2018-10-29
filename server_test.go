package tspool

import (
	"net"
	"testing"
	"time"
)

func TestListenAndServe(t *testing.T) {

	var err error
	addr := "0.0.0.0:8080"
	handler := func(net.Conn) {
	}
	errHandlerFunc := func(net.Conn, string) {
	}
	rejectHandler := func(net.Conn, string) {
	}
	srv := &Server{}
	go func() {
		err = ListenAndServe(srv)
		if err == nil {
			t.Errorf("ListenAndServe falied. Got nil, expected %s\n", errAddr)
		}

	}()
	go func() {
		srv.Addr = "wrong addr string"
		err = ListenAndServe(srv)
		if err == nil {
			t.Errorf("ListenAndServe falied. Got nil, expected %s\n", errAddr)
		}

	}()
	go func() {
		srv.Addr = addr
		err = ListenAndServe(srv)
		if err == nil {
			t.Errorf("ListenAndServe falied. Got nil, expected %s\n", errHandler)
		}
	}()

	go func() {
		srv.Handler = handler
		err = ListenAndServe(srv)
		if err == nil {
			t.Errorf("ListenAndServe falied. Got nil, expected %s\n", errErrHandler)
		}
	}()

	go func() {
		srv.ErrHandler = errHandlerFunc
		err = ListenAndServe(srv)
		if err == nil {
			t.Errorf("ListenAndServe falied. Got nil, expected %s\n", errRejectHandler)
		}

	}()

	go func() {
		srv.RejectHandler = rejectHandler
		err = ListenAndServe(srv)
		if err == nil {
			t.Error(err)
		}
	}()

	select {
	case <-time.After(5 * time.Second):
		return
	}
}
