package amap

import (
	"context"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/httptest"
	"time"

	cdp "github.com/chromedp/chromedp"
)

type mapProvider interface {
	http.Handler
	getName() string
}

var mapProviders = []mapProvider{Amap{}, Bmap{}}

// PrepareForShot prepares the http server to return html for amap
// Prepare 1 time, shot anytime
func PrepareForShot(mapProvider string) (url string, close func()) {
	for _, provider := range mapProviders {
		if mapProvider == provider.getName() {
			ts := httptest.NewServer(provider)
			return ts.URL, ts.Close
		}
	}
	return "", nil
}

// Shot captures the screenshot from amap
func Shot(url string, w int64, h int64, fn string) error {
	var buf []byte
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	ctx, cancel = cdp.NewContext(ctx, cdp.WithLogf(log.Printf))
	defer cancel()

	return cdp.Run(ctx, cdp.Tasks{
		cdp.EmulateViewport(w, h),
		cdp.Navigate(url),
		cdp.WaitNotVisible(`#tcBtn`),
		cdp.CaptureScreenshot(&buf),
		screenshotSave(fn, &buf),
	})
}

func screenshotSave(fileName string, buf *[]byte) cdp.ActionFunc {
	return func(ctx context.Context) error {
		log.Printf("Write %v", fileName)
		return ioutil.WriteFile(fileName, *buf, 0644)
	}
}
