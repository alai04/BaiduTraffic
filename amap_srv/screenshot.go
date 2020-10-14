package main

import (
	"context"
	"fmt"
	"io/ioutil"
	"log"
	"time"

	"github.com/chromedp/chromedp"
	cdp "github.com/chromedp/chromedp"
)

func shot(url string) {
	var buf []byte
	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()
	ctx, cancel = chromedp.NewContext(ctx, chromedp.WithLogf(log.Printf))
	defer cancel()

	done := make(chan struct{})
	go promptWaiting(done)
	url += `?lng=121.484&lat=31.216&z=16&b=1`
	err := cdp.Run(ctx, cdp.Tasks{
		chromedp.EmulateViewport(2600, 3200),
		cdp.Navigate(url),
		cdp.WaitNotVisible(`#tip`),
		cdp.CaptureScreenshot(&buf),
		screenshotSave("foo.png", &buf),
	})
	close(done)
	if err != nil {
		log.Fatal(err)
	}
}

func screenshotSave(fileName string, buf *[]byte) cdp.ActionFunc {
	return func(ctx context.Context) error {
		log.Printf("Write %v", fileName)
		return ioutil.WriteFile(fileName, *buf, 0644)
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
