package main

import (
	"log"
	"math"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestAMapLatLng2TileXYHR(t *testing.T) {
	c := NewAmapCoor()
	tx, ty, x, y := c.LatLng2TileXYHR(31.1936, 121.4595, 16)
	// tx, ty, x, y := c.LatLng2TileXYHR(1, 91, 2)
	assert.Equal(t, tx, 54879, "LatLng2TileXY")
	assert.Equal(t, ty, 26786, "LatLng2TileXY")
	assert.Equal(t, x, 6, "LatLng2TileXY")
	assert.Equal(t, y, 22, "LatLng2TileXY")
	tx, ty, x, y = c.LatLng2TileXYHR(31.2503, 121.5200, 16)
	assert.Equal(t, tx, 54890, "LatLng2TileXY")
	assert.Equal(t, ty, 26774, "LatLng2TileXY")
	assert.Equal(t, x, 10, "LatLng2TileXY")
	assert.Equal(t, y, 4, "LatLng2TileXY")
}

func TestAMapTileXYHR2LatLng(t *testing.T) {
	c := NewAmapCoor()
	lat, lng := c.TileXYHR2LatLng(54879, 26786, 6, 22, 16)
	assert.Less(t, math.Abs(lat-31.1936), 0.001, "TileXYHR2LatLng")
	assert.Less(t, math.Abs(lng-121.4595), 0.001, "TileXYHR2LatLng")
	lat, lng = c.TileXYHR2LatLng(54890, 26774, 10, 4, 16)
	assert.Less(t, math.Abs(lat-31.2503), 0.001, "TileXYHR2LatLng")
	assert.Less(t, math.Abs(lng-121.5200), 0.001, "TileXYHR2LatLng")
}

func TestAMapTileAdjacent(t *testing.T) {
	c := NewAmapCoor()
	lat1, lng1 := c.TileXYHR2LatLng(54879, 26786, 255, 255, 16)
	lat2, lng2 := c.TileXYHR2LatLng(54880, 26787, 0, 0, 16)
	log.Printf("%f, %f -- %f, %f", lat1, lng1, lat2, lng2)
	assert.Less(t, math.Abs(lat1-lat2), 0.001, "TestAMapTileAdjacent")
	assert.Less(t, math.Abs(lng1-lng2), 0.001, "TestAMapTileAdjacent")
}
