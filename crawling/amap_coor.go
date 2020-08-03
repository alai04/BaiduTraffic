package main

import (
	"image"
)

// AmapCoor 百度坐标系
type AmapCoor struct{}

// NewAmapCoor 返回百度坐标系
func NewAmapCoor() AmapCoor {
	return AmapCoor{}
}

// LatLng2TileXY 百度坐标转瓦片编号
func (amc AmapCoor) LatLng2TileXY(pLat, pLng float64, zoom int) (x int, y int) {
	x, y, _, _ = amc.LatLng2TileXYHR(pLat, pLng, zoom)
	return
}

// LatLng2TileXYHR 百度坐标转瓦片编号及像素点位置
func (amc AmapCoor) LatLng2TileXYHR(pLat, pLng float64, zoom int) (tx, ty, x, y int) {
	// mercX, mercY := amc.LatLng2Mercator(pLat, pLng)
	tx, ty = 0, 0 //amc.Mercator2TileXY(mercX, mercY, zoom)
	return
}

// TileXYHR2LatLng 瓦片编号及像素点位置转百度坐标
func (amc AmapCoor) TileXYHR2LatLng(tx, ty, x, y, zoom int) (lat, lng float64) {
	return
}

// GetTile 获取指定瓦片编号的路况图片
func (amc AmapCoor) GetTile(req mapRequest) (img image.Image, err error) {
	img = image.NewRGBA(image.Rect(0, 0, xSmallSize, ySmallSize))
	return
}
