package main

import (
	"fmt"
	"image"
	"image/draw"
	"image/png"
	"log"
	"math"
	"net/http"
	"os"
	"path/filepath"
	"sync"
	"time"
)

const (
	urlPrefix    = "http://its.map.baidu.com:8002/traffic/TrafficTileService"
	xBeginLevel9 = 97
	xEndLevel9   = 102
	yBeginLevel9 = 25
	yEndLevel9   = 32
	xSmallSize   = 256
	ySmallSize   = 256
)

// TrafficMap specify level of map, num of threads to get data from baidu.com
type TrafficMap struct {
	level   int
	workers int
	fileDir string
	wg      *sync.WaitGroup
	result  draw.Image
	nErr    int
}

type mapRequest struct {
	x   int
	y   int
	url string
}

type mapResponse struct {
	x   int
	y   int
	img image.Image
	err error
}

// NewTrafficMap return a struct TrafficMap with level & workers
func NewTrafficMap(level int, workers int, fileDir string) *TrafficMap {
	return &TrafficMap{
		level:   level,
		workers: workers,
		fileDir: fileDir,
		wg:      new(sync.WaitGroup),
	}
}

// GetMap return the filename of whole traffic map
func (m TrafficMap) GetMap() (mapFilename string, err error) {
	scale := int(math.Pow(2, float64(m.level-9)))
	xBegin, xEnd := xBeginLevel9*scale, xEndLevel9*scale
	yBegin, yEnd := yBeginLevel9*scale, yEndLevel9*scale

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
			url := fmt.Sprintf("%v?time=%d&level=%d&x=%d&y=%d", urlPrefix, ts, m.level, x, y)
			n++
			log.Printf("%d %v", n, url)
			chRequest <- mapRequest{
				x:   x - xBegin,
				y:   yEnd - 1 - y,
				url: url,
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
			res := mapResponse{req.x, req.y, nil, nil}
			var resp *http.Response
			resp, res.err = http.Get(req.url)
			if res.err == nil {
				// body, _ := ioutil.ReadAll(resp.Body)
				res.img, res.err = png.Decode(resp.Body)
				resp.Body.Close()
			}
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
			dp := image.Pt(res.x*xSmallSize, res.y*ySmallSize)
			draw.Draw(m.result, image.Rectangle{dp, dp.Add(res.img.Bounds().Size())},
				res.img, image.Pt(0, 0), draw.Src)
			log.Printf("Paste small image @ %d,%d", res.x, res.y)
		} else {
			m.nErr++
			log.Printf("Get small image @ %d,%d error: %v", res.x, res.y, res.err)
		}
	}
	close(done)
}
