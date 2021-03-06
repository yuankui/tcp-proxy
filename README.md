# tcp-proxy

in some cases, `client` may not connect to `realServer`,in this case, a proxy is required

this is a simple proxy with which when 

	client -> proxy(ip:port)

it acts like 
	
	client -> realServer(ip:port)


the process is like the following

	request -> proxy(ip:port) -> realServer(ip:port)
	client <- proxy(ip:port) <- response

## requirements

- [golang](http://golang.org) 1.5+
- [godep](https://github.com/tools/godep)

## build

	go get github.com/yuankui/tcp-proxy
	cd $GOPATH/src/github.com/yuankui/tcp-proxy
    godep restore
    go install .

## usage:

```
$ tcp-proxy

Usage:
  tcp-proxy

Application Options:
  -p, --port=   localhost listen port
  -r, --remote= remote ip:port
  -d, --debug=  debug port (default: 6060)
  -h, --help
```

## example

### http proxy

if `192.168.6.88:8080` is a http server

	$ tcp-proxy -p 5000 -r 192.168.6.88:8080
	
after the proxy is startup	

the following command acts the same

	$ curl localhost:5000
	$ curl 192.168.6.88:8080
	
### generic tcp proxy

if `192.168.6.88:8081` is a http server

	$ tcp-proxy -p 5000 -r 192.168.6.88:8081
	
after the proxy is startup	

the following command acts the same

	$ telnet localhost:5000
	$ telnet 192.168.6.88:8080
