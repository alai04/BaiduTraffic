package main

import (
	"image"
	"image/color"
	"image/draw"
	"image/png"
	"log"
	"os"

	"github.com/nfnt/resize"
)

var offset = image.Pt(195, 214)

func adjust(src, dst string) (err error) {
	bdFile, err := os.Open(src)
	if err != nil {
		log.Printf("Open %s error: %v", src, err)
		return
	}
	defer bdFile.Close()

	bdImage, err := png.Decode(bdFile)
	if err != nil {
		log.Printf("Decode %s error: %v", src, err)
		return
	}
	log.Printf("bdImage type: %T", bdImage)
	hideRects(bdImage.(*image.NRGBA))

	segImages := splitVert(bdImage, adParam.vertical)
	tmpImage := spliceVert(segImages)

	resImage := tmpImage.(*image.NRGBA)

	resFile, err := os.Create(dst)
	if err != nil {
		log.Printf("Create %s error: %v", dst, err)
		return
	}
	defer resFile.Close()
	return png.Encode(resFile, resImage)
}

func splitVert(img image.Image, adSegs []adjustSeg) []image.Image {
	tmpImage, ok := img.(*image.NRGBA)
	if !ok {
		log.Fatalln("input image not image.NRGBA")
	}
	w, h := tmpImage.Bounds().Dx(), tmpImage.Bounds().Dy()
	var segImages []image.Image
	cur := 0
	for _, adSeg := range adSegs {
		b := int(adSeg.b * float64(h))
		if cur > b {
			log.Fatalf("error parameter: %v", adSeg)
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
	return segImages
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
