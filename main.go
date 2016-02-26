package main

import (
	"net"
	"fmt"
	"os"
	"io"
	"strconv"
	"strings"
	_ "net/http/pprof"
	"log"
	"net/http"
)

var args struct{
	
}
func main() {

	go func() {
		log.Println(http.ListenAndServe("localhost:6060", nil))
	}()

	target := "127.0.0.1:8080"

	ipPort := strings.Split(target, ":")
	if len(ipPort) != 2 {
		fmt.Errorf("ip port not valid:%s", target)
		os.Exit(1)
	}

	targetPort, err := strconv.Atoi(ipPort[1])

	if err != nil {
		fmt.Errorf("port not valid:%s", ipPort[1])
	}

	targetIp := net.ParseIP(ipPort[0])

	listener, err := net.Listen("tcp", ":5000")
	if err != nil {
		fmt.Errorf(err.Error())
		os.Exit(1)
	}

	for {
		conn, err := listener.Accept()

		if err != nil {
			fmt.Errorf(err.Error())
		}

		go proxyConn(conn, targetIp, targetPort)
	}
}

func proxyConn(cliConn net.Conn, ip net.IP, port int) {
	serConn, err := net.Dial("tcp", ":8080")
	if err != nil {
		fmt.Errorf(err.Error())
		cliConn.Write([]byte(err.Error()))
	}

	finished := make(chan bool)
	go func(ch chan bool) {
		copyBuffer(cliConn, serConn, nil)
		fmt.Println("request finished.")
		ch <- true
	}(finished)

	go func(ch chan bool) {
		copyBuffer(serConn, cliConn, nil)
		fmt.Println("response finished.")
		ch <- true
	}(finished)

	<-finished

	serConn.Close()
	cliConn.Close()
}

func copyBuffer(dst io.Writer, src io.Reader, buf []byte) (written int64, err error) {
	if buf == nil {
		buf = make([]byte, 32 * 1024)
	}
	for {
		nr, er := src.Read(buf)
		fmt.Println("read:", nr)
		if nr > 0 {
			nw, ew := dst.Write(buf[0:nr])
			if nw > 0 {
				fmt.Println(nw)
				written += int64(nw)
			}
			if ew != nil {
				err = ew
				break
			}
			if nr != nw {
				err = io.ErrShortWrite
				break
			}
		}
		if er == io.EOF {
			break
		}
		if er != nil {
			err = er
			break
		}
	}
	return written, err
}


