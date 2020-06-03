package main

import (
	"flag"
	"fmt"
	"io"
	"net/http"
)

var (
	dstUrl = flag.String("destination", "", "destination url. e.g. http://127.0.0.1:9999")
	addr   = flag.String("address", ":9500", "server address")
)

func main() {
	flag.Parse()
	if "" == *dstUrl {
		panic("destination must be set")
	}

	fmt.Println("Serve on:", *addr)
	if err := http.ListenAndServe(*addr, new(proxy)); nil != err {
		panic(err)
	}
}

type proxy struct{}

func (*proxy) ServeHTTP(rw http.ResponseWriter, req *http.Request) {
	fmt.Sprintf("Received request %s %s\n", req.Method, req.RequestURI)

	//proxy request
	newUrl := fmt.Sprintf("%s%s", *dstUrl, req.URL.RequestURI())
	r, err := http.NewRequest(req.Method, newUrl, req.Body)
	if nil != err {
		panic(err)
	}
	r.Header = req.Header

	//proxy go
	resp, err := http.DefaultClient.Do(r)
	if nil != err {
		rw.WriteHeader(http.StatusBadGateway)
		fmt.Fprintln(rw, err.Error())
		return
	}

	//proxy response
	for key, values := range resp.Header {
		for _, v := range values {
			rw.Header().Add(key, v)
		}
	}
	rw.WriteHeader(resp.StatusCode)
	_, err = io.Copy(rw, resp.Body)
	if nil != err {
		fmt.Println("write to client, err", err.Error())
	}

	resp.Body.Close()
}
