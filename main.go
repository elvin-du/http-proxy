package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"net/http"
)

var (
	dstUrl = flag.String("destination", "", "destination url. e.g. http://127.0.0.1:9999")
)

func main() {
	flag.Parse()
	if "" == *dstUrl {
		panic("destination must be set")
	}

	http.ListenAndServe(":8888", new(proxy))
}

type proxy struct{}

func (*proxy) ServeHTTP(rw http.ResponseWriter, req *http.Request) {
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
		fmt.Fprintln(rw, err.Error())
		return
	}
	defer resp.Body.Close()

	//proxy response
	rw.WriteHeader(resp.StatusCode)
	bin, err := ioutil.ReadAll(resp.Body)
	if nil != err {
		panic(err)
	}

	_, err = rw.Write(bin)
	if nil != err {
		fmt.Println("write to client, err", err.Error())
	}
}
