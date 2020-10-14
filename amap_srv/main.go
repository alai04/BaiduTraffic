package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/alai04/BaiduTraffic/amap_srv/amap"
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
		// example: http://localhost:8080/amap?lng=121.484&lat=31.216&z=16
		mux := http.NewServeMux()
		mux.Handle("/amap", amap.Amap{})
		log.Printf("Listen on port %d", flagPort)
		log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", flagPort), mux))
		return
	}

	url, closeFunc := amap.PrepareForShot()
	defer closeFunc()

	done := make(chan struct{})
	go promptWaiting(done)
	amap.Shot(url)
	close(done)
}

func promptWaiting(done <-chan struct{}) {
	for {
		select {
		case <-done:
			return
		case <-time.After(time.Second):
			fmt.Print(".")
		}
	}
}
