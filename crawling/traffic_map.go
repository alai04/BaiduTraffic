package main

import (
	"fmt"
	"image"
	"image/draw"
	"image/png"
	"log"
	"os"
	"path/filepath"
	"sync"
	"time"
)

const (
	// urlPrefix  = "http://its.map.baidu.com:8002/traffic/TrafficTileService"
	xSmallSize = 256
	ySmallSize = 256
)

// TrafficMap specify level of map, num of threads to get data from baidu.com
type TrafficMap struct {
	config
	wg     *sync.WaitGroup
	result draw.Image
	nErr   int
}

type mapRequest struct {
	x      int
	y      int
	z      int
	xPaste int
	yPaste int
	coor   CoorSys
	ts     int64
}

type mapResponse struct {
	xPaste int
	yPaste int
	img    image.Image
	err    error
}

// NewTrafficMap return a struct TrafficMap with level & workers
func NewTrafficMap(cfg config) *TrafficMap {
	return &TrafficMap{
		config: cfg,
		wg:     new(sync.WaitGroup),
	}
}

// GetMap return the filename of whole traffic map
func (m TrafficMap) GetMap() (mapFilename string, err error) {
	var coor CoorSys = NewBaiduCoor()
	if m.amap {
		coor = NewAmapCoor()
	}
	xBegin, yBegin := coor.LatLng2TileXY(m.rect.b, m.rect.l, m.level)
	xEnd, yEnd := coor.LatLng2TileXY(m.rect.t, m.rect.r, m.level)
	if xBegin > xEnd {
		xBegin, xEnd = xEnd, xBegin
	}
	if yBegin > yEnd {
		yBegin, yEnd = yEnd, yBegin
	}
	xEnd++
	yEnd++

	tm := time.Now()
	tmStr := tm.Format("20060102_150405")
	ts := tm.UnixNano() / 1000000

	n := 0
	m.result = image.NewRGBA(image.Rect(0, 0, xSmallSize*(xEnd-xBegin), ySmallSize*(yEnd-yBegin)))

	chRequest := make(chan mapRequest, m.workers)
	chResponse := make(chan mapResponse)
	m.prepareWorkers(chRequest, chResponse)

	done := make(chan struct{})
	go m.pasteImages(done, chResponse)

	for x := xBegin; x < xEnd; x++ {
		for y := yBegin; y < yEnd; y++ {
			n++
			chRequest <- mapRequest{
				x:      x,
				y:      y,
				z:      m.level,
				xPaste: x - xBegin,
				yPaste: yEnd - 1 - y,
				coor:   coor,
				ts:     ts,
			}
		}
	}
	close(chRequest)
	log.Println("close(chRequest)")
	m.wg.Wait()
	log.Println("m.wg.Wait()")
	close(chResponse)
	log.Println("close(chResponse)")
	<-done
	log.Printf("Small images request: %d, error: %d", n, m.nErr)
	if m.nErr > n/2 {
		return "", fmt.Errorf("Too many error, skip the map")
	}

	mapFilename = filepath.Join(m.fileDir, fmt.Sprintf("L%d_%v.png", m.level, tmStr))
	tFilename := filepath.Join(m.fileDir, tmpFilename)
	out, _ := os.Create(tFilename)
	png.Encode(out, m.result)
	out.Close()
	return mapFilename, os.Rename(tFilename, mapFilename)
}

func (m TrafficMap) prepareWorkers(in <-chan mapRequest, out chan<- mapResponse) {
	worker := func(out chan<- mapResponse) {
		for req := range in {
			res := mapResponse{req.xPaste, req.yPaste, nil, nil}
			res.img, res.err = req.coor.GetTile(req)
			out <- res
		}
		m.wg.Done()
	}
	m.wg.Add(m.workers)
	for i := 0; i < m.workers; i++ {
		go worker(out)
	}
}

func (m *TrafficMap) pasteImages(done chan<- struct{}, ch <-chan mapResponse) {
	for res := range ch {
		if res.err == nil {
			dp := image.Pt(res.xPaste*xSmallSize, res.yPaste*ySmallSize)
			draw.Draw(m.result, image.Rectangle{dp, dp.Add(res.img.Bounds().Size())},
				res.img, image.Pt(0, 0), draw.Src)
			log.Printf("Paste small image @ %d,%d", res.xPaste, res.yPaste)
		} else {
			m.nErr++
			log.Printf("Get small image @ %d,%d error: %v", res.xPaste, res.yPaste, res.err)
		}
	}
	close(done)
}
