package main

import (
	"fmt"
	"github.com/jessevdk/go-flags"
	"io"
	"net"
	"net/http"
	_ "net/http/pprof"
	"os"
)

var args struct {
	LocalPort      int    `short:"p" long:"port" description:"localhost listen port" required:"true"`
	RemoteHostPort string `short:"r" long:"remote" description:"remote ip:port" required:"true"`
	DebugPort      int    `short:"d" long:"debug" description:"debug port" default:"6060"`
	Help           bool   `short:"h" long:"help" descrition:"the help message"`
}

var usage string
var parser *flags.Parser

func parseArgs(args interface{}) (e error) {
	_, err := parser.ParseArgs(os.Args)

	if err != nil {
		e = err
	}
	return
}

func printUsage() {
	parser.WriteHelp(os.Stderr)
}

func init() {
	parser = flags.NewParser(&args, flags.None)
	usage = parser.Usage
}

func printError(format string, a ...interface{}) {
	fmt.Fprintf(os.Stderr, format, a...)
	fmt.Printf("\n")
}

func main() {

	var err error
	err = parseArgs(&args)

	if err != nil {
		printUsage()
		return
	}

	go func() {
		err := http.ListenAndServe(fmt.Sprintf(":%d", args.DebugPort), nil)
		if err != nil {
			printError(err.Error())
			printUsage()
			os.Exit(1)
		}
	}()

	localBinding := fmt.Sprintf(":%d", args.LocalPort)
	listener, err := net.Listen("tcp", localBinding)
	if err != nil {
		printError(err.Error())
		printUsage()
		return
	}

	for {
		conn, err := listener.Accept()

		if err != nil {
			printError(err.Error())
		}

		go proxyConn(conn, args.RemoteHostPort)
	}
}

func proxyConn(cliConn net.Conn, addr string) {
	serConn, err := net.Dial("tcp", addr)
	if err != nil {
		printError(err.Error())
		cliConn.Write([]byte(err.Error()))
	}

	finished := make(chan bool)
	go func(ch chan bool) {
		copyBuffer(cliConn, serConn, nil)
		ch <- true
	}(finished)

	go func(ch chan bool) {
		copyBuffer(serConn, cliConn, nil)
		ch <- true
	}(finished)

	<-finished

	serConn.Close()
	cliConn.Close()
}

func copyBuffer(dst io.Writer, src io.Reader, buf []byte) (written int64, err error) {
	if buf == nil {
		buf = make([]byte, 32*1024)
	}
	for {
		nr, er := src.Read(buf)
		if nr > 0 {
			nw, ew := dst.Write(buf[0:nr])
			if nw > 0 {
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
