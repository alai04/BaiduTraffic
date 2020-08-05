package main

import (
	"fmt"
	"image"
	"image/png"
	"log"
	"math"
	"net/http"
)

// AmapCoor 高德坐标系
type AmapCoor struct{}

// NewAmapCoor 返回高德坐标系
func NewAmapCoor() AmapCoor {
	return AmapCoor{}
}

// YDirectionN 返回Y轴方向是否向N
func (amc AmapCoor) YDirectionN() bool {
	return false
}

// ScalingFactor 高德的瓦片放大系数为4，一次可以读取1024*1024
func (amc AmapCoor) ScalingFactor() int {
	return 4
}

// LatLng2TileXY 高德坐标转瓦片编号
func (amc AmapCoor) LatLng2TileXY(pLat, pLng float64, zoom int) (tx int, ty int) {
	tx, ty, _, _ = amc.LatLng2TileXYHR(pLat, pLng, zoom)
	return
}

// LatLng2TileXYHR 高德坐标转瓦片编号及像素点位置
func (amc AmapCoor) LatLng2TileXYHR(lat, lng float64, zoom int) (tx, ty, x, y int) {
	fx := (lng + 180) / 360 * math.Pow(2, float64(zoom))
	tx = int(math.Floor(fx))
	x = int(fx*256) % 256
	latRad := lat * math.Pi / 180
	fy := (1 - math.Log(math.Tan(latRad)+1/math.Cos(latRad))/math.Pi) / 2 * math.Pow(2, float64(zoom))
	ty = int(math.Floor(fy))
	y = int(fy*256) % 256
	return
}

// TileXYHR2LatLng 瓦片编号及像素点位置转高德坐标
func (amc AmapCoor) TileXYHR2LatLng(tx, ty, x, y, zoom int) (lat, lng float64) {
	pixelXToTileAddition := float64(x) / 256.0
	lng = (float64(tx)+pixelXToTileAddition)/math.Pow(2, float64(zoom))*360 - 180
	pixelYToTileAddition := float64(y) / 256.0
	lat = math.Atan(sinh(math.Pi*(1-2*(float64(ty)+pixelYToTileAddition)/math.Pow(2, float64(zoom))))) * 180.0 / math.Pi
	return
}

// GetTile 获取指定瓦片编号的路况图片
func (amc AmapCoor) GetTile(req mapRequest) (img image.Image, err error) {
	lat, lng := amc.TileXYHR2LatLng(req.x, req.y, 0, 0, req.z)
	const urlPrefix = "https://restapi.amap.com/v3/staticmap?&size=1024*1024&key=5c3ff1a9ba96bd86123c21fd377a2e8c&traffic=1"
	url := fmt.Sprintf("%v&zoom=%d&location=%f,%f", urlPrefix, req.z-1, lng, lat)
	log.Printf("%d,%d %v", req.xPaste, req.yPaste, url)
	var resp *http.Response
	resp, err = http.Get(url)
	if err == nil {
		img, err = png.Decode(resp.Body)
		resp.Body.Close()
	}
	return
}

func sinh(x float64) float64 {
	return (math.Exp(x) - math.Exp(-x)) / 2
}
