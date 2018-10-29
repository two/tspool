# tspool
A TCP Library use worker pool to improve performance and protect your server.

[![Build Status](https://travis-ci.org/two/tspool.svg?branch=master)](https://travis-ci.org/two/tspool)
[![Code Coverage](https://codecov.io/gh/two/tspool/branch/master/graph/badge.svg)](https://codecov.io/gh/two/tspool)
[![License](https://img.shields.io/github/license/mashape/apistatus.svg)](https://opensource.org/licenses/MIT)
[![GoDoc](https://godoc.org/github.com/two/tspool?status.svg)](https://godoc.org/github.com/two/tspool)


## Install

```
go get github.com/two/tspool
```

## Usage

Build your server with tspool
### server
```go
package main

import (
	"github.com/two/tspool"
	"log"
	"net"
)

func main() {
	wp, err := tspool.DefaultWorkerPool(100, 200)
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
```
Build your tspool client
```go

package main

import (
	"log"
	"net"
	"sync"
	"time"
)

var wg = new(sync.WaitGroup)

func main() {
	var num int = 2000
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
```
## Define your own worker pool and worker

## Example
[example](https://github.com/two/tspool/tree/master/examples)
