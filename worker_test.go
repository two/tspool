package tspool

import (
	"log"
	"net"
	"testing"
	"time"
)

func server(t *testing.T) {
	wp, err := DefaultWorkerPool(1, 2)
	if err != nil {
		t.Error(err)
	}
	server := &Server{
		Addr:          "0.0.0.0:8088",
		Handler:       handler,
		ErrHandler:    errHandlerFunc,
		WorkerPool:    wp,
		RejectHandler: rejectHandler,
	}
	go func() {
		err = ListenAndServe(server)
		if err != nil {
			t.Error(err)
		}
	}()
	select {
	case <-time.After(5 * time.Second):
		return
	}

}
func handler(c net.Conn) {
	addr := c.RemoteAddr()
	c.Write([]byte("hello"))
	c.Close()
	log.Println(addr.String())
}

func errHandlerFunc(c net.Conn, err string) {
	defer c.Close()
	log.Fatalln("run server error: " + err)
}

func rejectHandler(c net.Conn, err string) {
	defer c.Close()
	log.Println("reject connect error: " + err)
}

func TestWork(t *testing.T) {
	server(t)
	addr := "0.0.0.0:8088"
	d := net.Dialer{Timeout: 100 * time.Millisecond}
	conn, err := d.Dial("tcp", addr)
	if err != nil {
		log.Fatalln(err)
	}
	defer conn.Close()
	var r = make([]byte, 1024)
	_, err = conn.Read(r)
	if err != nil {
		t.Error(err)
	}
}
