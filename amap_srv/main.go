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
	flagServe       bool
	flagPort        int
	flagCenter      string
	flagWidth       int64
	flagHeight      int64
	flagZoom        int
	flagBase        bool
	flagMapDir      string
	flagMapProvider string
)

func init() {
	flag.BoolVar(&flagServe, "s", false, "以Web服务方式运行")
	flag.IntVar(&flagPort, "p", 8080, "监听端口")
	flag.BoolVar(&flagBase, "b", false, "是否带底图(仅对高德地图有效)")
	flag.StringVar(&flagMapProvider, "m", "amap", "地图提供方，缺省为高德地图")
	flag.Int64Var(&flagWidth, "w", 256, "地图图片宽度")
	flag.Int64Var(&flagHeight, "h", 256, "地图图片高度")
	flag.IntVar(&flagZoom, "z", 10, "地图缩放级别(3-20)")
	flag.StringVar(&flagCenter, "c", "121.4833305,31.216065", "地图中心点坐标")
	flag.StringVar(&flagMapDir, "d", "map", "地图图片输出目录")
	flag.Parse()
}

func main() {
	if flagServe {
		// example: http://localhost:8080/amap?lng=121.484&lat=31.216&z=16
		mux := http.NewServeMux()
		mux.Handle("/amap", amap.Amap{})
		mux.Handle("/bmap", amap.Bmap{})
		log.Printf("Listen on port %d", flagPort)
		log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", flagPort), mux))
		return
	}

	url, closeFunc := amap.PrepareForShot(flagMapProvider)
	if url == "" || closeFunc == nil {
		log.Fatal(fmt.Errorf("MapProvider %q not found", flagMapProvider))
	}
	defer closeFunc()

	var lng, lat float64
	n, err := fmt.Sscanf(flagCenter, "%f,%f", &lng, &lat)
	if err != nil {
		log.Fatalf("中心点坐标不正确: %q, %v", flagCenter, err)
	}
	if n != 2 {
		log.Fatalf("中心点坐标不正确: %q", flagCenter)
	}
	url = fmt.Sprintf("%s?lng=%f&lat=%f&z=%d", url, lng, lat, flagZoom)
	if flagBase {
		url += "&b=1"
	}
	filename := fmt.Sprintf("%s/L%02d_%s.png", flagMapDir, flagZoom,
		time.Now().Format("20060102_150405"))

	done := make(chan struct{})
	go promptWaiting(done)
	err = amap.Shot(url, flagWidth, flagHeight, filename)
	close(done)
	if err != nil {
		log.Fatal(err)
	} else {
		log.Printf("Traffic map save to: %v", filename)
	}
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
