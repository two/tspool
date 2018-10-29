package main

import (
	"log"
	"net"
	"sync"
	"time"
)

var wg = new(sync.WaitGroup)

func main() {
	var num int = 20
	wg.Add(num)
	for i := 0; i < num; i++ {
		go conn()
	}
	wg.Wait()
	log.Println("done")
}
func conn() {
	defer wg.Done()
	addr := "0.0.0.0:8088"
	d := net.Dialer{Timeout: 100 * time.Millisecond}
	conn, err := d.Dial("tcp", addr)
	if err != nil {
		log.Fatalln(err)
	}
	defer conn.Close()
	raddr := conn.LocalAddr()
	var r = make([]byte, 1024)
	_, err = conn.Read(r)
	if err != nil {
		log.Printf(raddr.String() + " read error: " + err.Error())
	}
	log.Printf(raddr.String() + " got " + string(r))
}
