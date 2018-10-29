package main

import (
	"github.com/two/tspool"
	"log"
	"net"
)

func main() {
	wp, err := tspool.DefaultWorkerPool(1, 2)
	if err != nil {
		log.Fatalln(err)
	}
	server := &tspool.Server{
		Addr:          "0.0.0.0:8088",
		Handler:       handler,
		ErrHandler:    errHandler,
		WorkerPool:    wp,
		RejectHandler: rejectHandler,
	}
	err = tspool.ListenAndServe(server)
	if err != nil {
		log.Fatalln(err)
	}
}

func handler(c net.Conn) {
	addr := c.RemoteAddr()
	c.Write([]byte("hello"))
	c.Close()
	log.Println(addr.String())
}

func errHandler(c net.Conn, err string) {
	defer c.Close()
	log.Fatalln("run server error: " + err)
}

func rejectHandler(c net.Conn, err string) {
	defer c.Close()
	log.Println("reject connect error: " + err)
}
