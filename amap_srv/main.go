package main

import (
	"flag"
	"net/http/httptest"
)

var (
	flagServe bool
	flagPort  int
)

func init() {
	flag.BoolVar(&flagServe, "s", false, "Run as web server")
	flag.IntVar(&flagPort, "p", 8080, "Port for listening on")
	flag.Parse()
}

func main() {
	if flagServe {
		serve()
		return
	}

	ts := httptest.NewServer(amap{})
	defer ts.Close()
	shot(ts.URL)
}
