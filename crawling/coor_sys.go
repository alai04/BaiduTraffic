package main

import (
	"image"
)

// CoorSys 表示各个不同坐标系
type CoorSys interface {
	YDirectionN() bool
	LatLng2TileXY(pLat, pLng float64, zoom int) (tx, ty int)
	LatLng2TileXYHR(pLat, pLng float64, zoom int) (tx, ty, x, y int)
	TileXYHR2LatLng(tx, ty, x, y, zoom int) (lat, lng float64)
	GetTile(req mapRequest) (img image.Image, err error)
}
