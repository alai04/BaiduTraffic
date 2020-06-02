package main

import (
	"fmt"
	"image"
	"image/color"
	"image/draw"
	"image/png"
	"log"
	"os"
	"time"

	"github.com/nfnt/resize"
)

var (
	timeConsuming time.Duration
	offset        = image.Pt(195, 214)
)

func timeCost() func() {
	start := time.Now()
	return func() {
		timeConsuming = time.Since(start)
	}
}

func adjust(src, dst string) (stage string, err error) {
	defer timeCost()()
	stage = "OPEN"
	bdFile, err := os.Open(src)
	if err != nil {
		log.Printf("Open %s error: %v", src, err)
		return
	}
	defer bdFile.Close()

	stage = "DECODE"
	bdImage, err := png.Decode(bdFile)
	if err != nil {
		log.Printf("Decode %s error: %v", src, err)
		return
	}
	log.Printf("bdImage type: %T", bdImage)
	hideRects(bdImage.(*image.NRGBA))

	stage = "ADJUST"
	segImages, err := splitVert(bdImage, adParam.vertical)
	if err != nil {
		log.Printf("Adjust %s error: %v", src, err)
		return
	}
	tmpImage := spliceVert(segImages)

	resImage := tmpImage.(*image.NRGBA)

	stage = "CREATE"
	resFile, err := os.Create(dst)
	if err != nil {
		log.Printf("Create %s error: %v", dst, err)
		return
	}
	defer resFile.Close()

	stage = "WRITE"
	err = png.Encode(resFile, resImage)
	return
}

func splitVert(img image.Image, adSegs []adjustSeg) (segImages []image.Image, err error) {
	tmpImage, ok := img.(*image.NRGBA)
	if !ok {
		err = fmt.Errorf("input image not image.NRGBA")
		return
	}
	w, h := tmpImage.Bounds().Dx(), tmpImage.Bounds().Dy()
	cur := 0
	for _, adSeg := range adSegs {
		b := int(adSeg.b * float64(h))
		if cur > b {
			err = fmt.Errorf("error parameter: %v", adSeg)
			return
		} else if cur < b {
			segImages = append(segImages, tmpImage.SubImage(image.Rect(0, cur, w, b)))
		}
		e := int(adSeg.e * float64(h))
		segImages = append(segImages,
			resize.Resize(uint(w), uint(float64(e-b)*adSeg.factor),
				tmpImage.SubImage(image.Rect(0, b, w, e)), resize.NearestNeighbor))
		cur = e
	}
	if cur < h {
		segImages = append(segImages, tmpImage.SubImage(image.Rect(0, cur, w, h)))
	}
	return
}

func spliceVert(images []image.Image) image.Image {
	var w, h int
	for _, img := range images {
		h += img.Bounds().Dy()
		if w < img.Bounds().Dx() {
			w = img.Bounds().Dx()
		}
	}
	resImage := image.NewNRGBA(image.Rect(0, 0, w, h))
	log.Printf("resImage size: %d x %d", w, h)
	var cur int
	for _, img := range images {
		draw.Draw(resImage, img.Bounds().Add(image.Pt(0, cur)), img, image.ZP, draw.Over)
		cur += img.Bounds().Dy()
	}
	return resImage
}

func hideRects(img *image.NRGBA) {
	rects := [][4]float64{
		{0.6837, 0.6438, 0.0260, 0.0510},
		{0.6927, 0.6938, 0.0330, 0.0586},
	}
	hide := color.NRGBA{0, 0, 0, 0}
	var x1, y1, w, h int
	wl, hl := img.Bounds().Dx(), img.Bounds().Dy()
	for _, rect := range rects {
		x1, y1 = int(rect[0]*float64(wl)), int(rect[1]*float64(hl))
		w, h = int(rect[2]*float64(wl)), int(rect[3]*float64(hl))
		x2, y2 := x1+w, y1+h
		draw.Draw(img, image.Rect(x1, y1, x2, y2), &image.Uniform{hide}, image.ZP, draw.Src)
	}
}

func combineToFW(img image.Image) image.Image {
	fwFilename := "IMG_6117.PNG"
	fwFile, err := os.Open(fwFilename)
	if err != nil {
		log.Fatalf("Open %s error: %v", fwFilename, err)
	}
	defer fwFile.Close()

	fwImage, err := png.Decode(fwFile)
	if err != nil {
		log.Fatalf("Decode %s error: %v", fwFilename, err)
	}

	tmpImage := resize.Resize(2140, 2550, img, resize.Bicubic)
	resImage := image.NewNRGBA(fwImage.Bounds())
	draw.Draw(resImage, resImage.Bounds(), fwImage, image.ZP, draw.Src)
	draw.Draw(resImage, tmpImage.Bounds(), tmpImage, offset, draw.Over)
	return resImage
}
