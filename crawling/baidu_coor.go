package main

import (
	"fmt"
	"image"
	"image/png"
	"log"
	"math"
	"net/http"
)

// BaiduCoor 百度坐标系
type BaiduCoor struct{}

// NewBaiduCoor 返回百度坐标系
func NewBaiduCoor() BaiduCoor {
	return BaiduCoor{}
}

var (
	array1 = []float64{75, 60, 45, 30, 15, 0}
	array3 = []float64{12890594.86, 8362377.87, 5591021, 3481989.83, 1678043.12, 0}
	array2 = [][]float64{
		{-0.0015702102444, 111320.7020616939, 1704480524535203,
			-10338987376042340, 26112667856603880, -35149669176653700,
			26595700718403920, -10725012454188240, 1800819912950474, 82.5},
		{0.0008277824516172526, 111320.7020463578, 647795574.6671607,
			-4082003173.641316, 10774905663.51142, -15171875531.51559,
			12053065338.62167, -5124939663.577472, 913311935.9512032, 67.5},
		{0.00337398766765, 111320.7020202162, 4481351.045890365,
			-23393751.19931662, 79682215.47186455, -115964993.2797253,
			97236711.15602145, -43661946.33752821, 8477230.501135234, 52.5},
		{0.00220636496208, 111320.7020209128, 51751.86112841131,
			3796837.749470245, 992013.7397791013, -1221952.21711287,
			1340652.697009075, -620943.6990984312, 144416.9293806241, 37.5},
		{-0.0003441963504368392, 111320.7020576856, 278.2353980772752,
			2485758.690035394, 6070.750963243378, 54821.18345352118,
			9540.606633304236, -2710.55326746645, 1405.483844121726, 22.5},
		{-0.0003218135878613132, 111320.7020701615, 0.00369383431289,
			823725.6402795718, 0.46104986909093, 2351.343141331292,
			1.58060784298199, 8.77738589078284, 0.37238884252424, 7.45}}
	array4 = [][]float64{
		{1.410526172116255e-8, 0.00000898305509648872, -1.9939833816331,
			200.9824383106796, -187.2403703815547, 91.6087516669843,
			-23.38765649603339, 2.57121317296198, -0.03801003308653, 17337981.2},
		{-7.435856389565537e-9, 0.000008983055097726239, -0.78625201886289,
			96.32687599759846, -1.85204757529826, -59.36935905485877,
			47.40033549296737, -16.50741931063887, 2.28786674699375, 10260144.86},
		{-3.030883460898826e-8, 0.00000898305509983578, 0.30071316287616,
			59.74293618442277, 7.357984074871, -25.38371002664745,
			13.45380521110908, -3.29883767235584, 0.32710905363475, 6856817.37},
		{-1.981981304930552e-8, 0.000008983055099779535, 0.03278182852591,
			40.31678527705744, 0.65659298677277, -4.44255534477492,
			0.85341911805263, 0.12923347998204, -0.04625736007561, 4482777.06},
		{3.09191371068437e-9, 0.000008983055096812155, 0.00006995724062,
			23.10934304144901, -0.00023663490511, -0.6321817810242,
			-0.00663494467273, 0.03430082397953, -0.00466043876332, 2555164.4},
		{2.890871144776878e-9, 0.000008983055095805407, -3.068298e-8,
			7.47137025468032, -0.00000353937994, -0.02145144861037,
			-0.00001234426596, 0.00010322952773, -0.00000323890364, 826088.5}}
)

// LatLng2Mercator 百度坐标转墨卡托(纬度,经度)-->(横向x,纵向y)
func (bdc BaiduCoor) LatLng2Mercator(pLat, pLng float64) (x int, y int) {
	var arr []float64
	nLat := pLat
	if pLat > 74 {
		nLat = 74
	}
	if pLat < -74 {
		nLat = -74
	}
	for i := range array1 {
		if pLat >= array1[i] {
			arr = array2[i]
			break
		}
	}
	if arr == nil {
		for i := len(array1) - 1; i >= 0; i-- {
			if pLat <= -array1[i] {
				arr = array2[i]
				break
			}
		}
	}
	fx, fy := bdc.Convertor(pLng, nLat, arr)
	x, y = int(math.Floor(fx)), int(math.Floor(fy))
	return
}

// Mercator2LatLng 墨卡托坐标转百度(横向像素x,纵向像素y)
func (bdc BaiduCoor) Mercator2LatLng(pX, pY int) (lat, lon float64) {
	var arr []float64
	for i := range array3 {
		if math.Abs(float64(pY)) >= array3[i] {
			arr = array4[i]
			break
		}
	}
	return bdc.Convertor(math.Abs(float64(pX)), math.Abs(float64(pY)), arr)
}

// Convertor 墨卡托坐标转换函数
func (bdc BaiduCoor) Convertor(x, y float64, param []float64) (float64, float64) {
	T := param[0] + param[1]*math.Abs(x)
	cC := math.Abs(y) / param[9]
	cF := param[2] + param[3]*cC + param[4]*cC*cC + param[5]*cC*cC*cC +
		param[6]*cC*cC*cC*cC + param[7]*cC*cC*cC*cC*cC + param[8]*cC*cC*cC*cC*cC*cC
	if x < 0 {
		T = -T
	}
	if y < 0 {
		cF = -cF
	}
	return T, cF
}

// Mercator2TileXY 获取瓦片编号
func (bdc BaiduCoor) Mercator2TileXY(pMercX, pMercY, zoom int) (int, int) {
	resolution := math.Pow(2, float64(18-zoom))
	x := float64(pMercX) / resolution / 256
	y := float64(pMercY) / resolution / 256
	return int(math.Floor(x)), int(math.Floor(y))
}

// LatLng2TileXY 百度坐标转瓦片编号
func (bdc BaiduCoor) LatLng2TileXY(pLat, pLng float64, zoom int) (x int, y int) {
	x, y, _, _ = bdc.LatLng2TileXYHR(pLat, pLng, zoom)
	return
}

// LatLng2TileXYHR 百度坐标转瓦片编号及像素点位置
func (bdc BaiduCoor) LatLng2TileXYHR(pLat, pLng float64, zoom int) (tx, ty, x, y int) {
	mercX, mercY := bdc.LatLng2Mercator(pLat, pLng)
	tx, ty = bdc.Mercator2TileXY(mercX, mercY, zoom)
	return
}

// TileXYHR2LatLng 瓦片编号及像素点位置转百度坐标
func (bdc BaiduCoor) TileXYHR2LatLng(tx, ty, x, y, zoom int) (lat, lng float64) {
	return
}

// GetTile 获取指定瓦片编号的路况图片
func (bdc BaiduCoor) GetTile(req mapRequest) (img image.Image, err error) {
	const urlPrefix string = "http://its.map.baidu.com:8002/traffic/TrafficTileService"
	url := fmt.Sprintf("%v?time=%d&level=%d&x=%d&y=%d", urlPrefix, req.ts, req.z, req.x, req.y)
	log.Printf("%d,%d %v", req.xPaste, req.yPaste, url)
	var resp *http.Response
	resp, err = http.Get(url)
	if err == nil {
		img, err = png.Decode(resp.Body)
		resp.Body.Close()
	}
	return
}
