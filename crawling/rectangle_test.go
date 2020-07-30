package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRectangle(t *testing.T) {
	rect := rectangle{121.4595, 31.1936, 121.5200, 31.2503}
	strRect := "121.4595,31.1936,121.5200,31.2503"
	assert.Equal(t, rect.String(), strRect, "rectangle.String()")
	var rect2 rectangle
	rect2.Set(strRect)
	assert.Equal(t, rect2, rect, "rectangle.Set()")
	assert.Equal(t, rect2.String(), strRect, "rectangle.String()")
}

func TestRectangleInvalid(t *testing.T) {
	strRect := "121.4595,31.1936,121.5200"
	var rect rectangle
	assert.NotNil(t, rect.Set(strRect))
	strRect = "aa,121.4595,31.1936,121.5200"
	assert.NotNil(t, rect.Set(strRect))
}
