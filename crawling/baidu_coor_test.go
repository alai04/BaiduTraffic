package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestLatLng2TileXY(t *testing.T) {
	bdc := NewBaiduCoor()
	x, y := bdc.LatLng2TileXY(31.1936, 121.4595, 16)
	assert.Equal(t, x, 13204, "LatLng2TileXY")
	assert.Equal(t, y, 3550, "LatLng2TileXY")
	x, y = bdc.LatLng2TileXY(31.2503, 121.5200, 16)
	assert.Equal(t, x, 13210, "LatLng2TileXY")
	assert.Equal(t, y, 3557, "LatLng2TileXY")
}
