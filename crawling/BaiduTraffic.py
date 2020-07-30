#!/usr/bin/env python3
# -*-coding:utf-8-*-

import time
import math
from urllib.request import urlretrieve
import PIL
import PIL.Image as Image

array1 = [75, 60, 45, 30, 15, 0]
array3 = [12890594.86, 8362377.87, 5591021, 3481989.83, 1678043.12, 0]
array2 = (
    [-0.0015702102444, 111320.7020616939, 1704480524535203,
     -10338987376042340, 26112667856603880, -35149669176653700,
     26595700718403920, -10725012454188240, 1800819912950474, 82.5],
    [0.0008277824516172526, 111320.7020463578, 647795574.6671607,
     -4082003173.641316, 10774905663.51142, -15171875531.51559,
     12053065338.62167, -5124939663.577472, 913311935.9512032, 67.5],
    [0.00337398766765, 111320.7020202162, 4481351.045890365,
     -23393751.19931662, 79682215.47186455, -115964993.2797253,
     97236711.15602145, -43661946.33752821, 8477230.501135234, 52.5],
    [0.00220636496208, 111320.7020209128, 51751.86112841131,
     3796837.749470245, 992013.7397791013, -1221952.21711287,
     1340652.697009075, -620943.6990984312, 144416.9293806241, 37.5],
    [-0.0003441963504368392, 111320.7020576856, 278.2353980772752,
     2485758.690035394, 6070.750963243378, 54821.18345352118,
     9540.606633304236, -2710.55326746645, 1405.483844121726, 22.5],
    [-0.0003218135878613132, 111320.7020701615, 0.00369383431289,
     823725.6402795718, 0.46104986909093, 2351.343141331292,
     1.58060784298199, 8.77738589078284, 0.37238884252424, 7.45])
array4 = (
    [1.410526172116255e-8, 0.00000898305509648872, -1.9939833816331,
     200.9824383106796, -187.2403703815547, 91.6087516669843,
     -23.38765649603339, 2.57121317296198, -0.03801003308653, 17337981.2],
    [-7.435856389565537e-9, 0.000008983055097726239, -0.78625201886289,
     96.32687599759846, -1.85204757529826, -59.36935905485877,
     47.40033549296737, -16.50741931063887, 2.28786674699375, 10260144.86],
    [-3.030883460898826e-8, 0.00000898305509983578, 0.30071316287616,
     59.74293618442277, 7.357984074871, -25.38371002664745,
     13.45380521110908, -3.29883767235584, 0.32710905363475, 6856817.37],
    [-1.981981304930552e-8, 0.000008983055099779535, 0.03278182852591,
     40.31678527705744, 0.65659298677277, -4.44255534477492,
     0.85341911805263, 0.12923347998204, -0.04625736007561, 4482777.06],
    [3.09191371068437e-9, 0.000008983055096812155, 0.00006995724062,
     23.10934304144901, -0.00023663490511, -0.6321817810242,
     -0.00663494467273, 0.03430082397953, -0.00466043876332, 2555164.4],
    [2.890871144776878e-9, 0.000008983055095805407, -3.068298e-8,
     7.47137025468032, -0.00000353937994, -0.02145144861037,
     -0.00001234426596, 0.00010322952773, -0.00000323890364, 826088.5])


# 百度坐标转墨卡托(纬度,经度)-->(横向x,纵向y)


def LatLng2Mercator(pLat, pLng):
    arr = None
    n_lat = pLat
    if(pLat > 74):
        n_lat = 74
    if(pLat < -74):
        n_lat = -74
    for i in range(0, len(array1)):
        if (pLat >= array1[i]):
            arr = array2[i]
            break
    if (arr == None):
        for i in range(-1, -len(array1)):
            if (pLat <= -array1[i]):
                arr = array2[i]
                break
    res = Convertor(pLng, n_lat, arr)
    res[0] = math.floor(res[0])
    res[1] = math.floor(res[1])
    return [res[0], res[1]]
# 墨卡托坐标转百度(横向像素x,纵向像素y)


def Mercator2LatLng(pX, pY):
    arr = None
    pX = math.floor(pX)
    pY = math.floor(pY)
    for i in range(0, len(array3)):
        if (abs(pY) >= array3[i]):
            arr = array4[i]
            break
    res = Convertor(abs(pX), abs(pY), arr)
    return [res[0], res[1]]

# 墨卡托坐标转换函数


def Convertor(x, y, param):
    T = param[0] + param[1] * abs(x)
    cC = abs(y) / param[9]
    cF = param[2] + param[3] * cC + param[4] * cC * cC + param[5] * cC * cC * cC + param[6] * \
        cC * cC * cC * cC + param[7] * cC * cC * cC * \
        cC * cC + param[8] * cC * cC * cC * cC * cC * cC
    if(x < 0):
        T = T*(-1)
    if(y < 0):
        cF = cF * (-1)
    return [T, cF]

# 获取瓦片编号
# 墨卡托坐标后，将结果除以地图分辨率(对2开(18-zoom)次方)即可得到平面像素坐标，然后将像素坐标除以256分别得到瓦片的行列号。


def Mercator2TileXY(pMercX, pMercY, zoom):
    resolution = pow(2, 18-zoom)
    x = int(pMercX/resolution)/256
    y = int(pMercY/resolution)/256
    x = math.floor(x)
    y = math.floor(y)
    return [x, y]


def get_traffic_map(level):

    api_url = 'http://its.map.baidu.com:8002/traffic/TrafficTileService'
    tm = time.time()
    ts = int(tm * 1000)
    tm_str = time.strftime('%Y%m%d_%H%M%S', time.localtime(tm))
    left, right = 114.4000, 120.0000
    bottom, top = 28.6000, 35.2000
    zoom = level
    pMerc1 = LatLng2Mercator(bottom, left)
    pMerc2 = LatLng2Mercator(top, right)
    bottom_left_tile = Mercator2TileXY(pMerc1[0], pMerc1[1], zoom)
    top_right_tile = Mercator2TileXY(pMerc2[0], pMerc2[1], zoom)
    print("左下:" + str(bottom_left_tile))
    print("右上:" + str(top_right_tile))
    return "****"

    # x_begin_level9, x_end_level9 = 97, 102
    # y_begin_level9, y_end_level9 = 25, 32
    # scale = 2**(level - 9)
    # x_begin, x_end = x_begin_level9 * scale, x_end_level9 * scale
    # y_begin, y_end = y_begin_level9 * scale, y_end_level9 * scale

    x_begin, y_begin = bottom_left_tile
    x_end, y_end = top_right_tile
    SMALL_SIZE_X = 256
    SMALL_SIZE_Y = 256
    large_image = Image.new(
        'RGBA', (SMALL_SIZE_X*(x_end-x_begin), SMALL_SIZE_Y*(y_end-y_begin)))

    n = 0
    n_err = 0
    for x in range(x_begin, x_end):
        for y in range(y_begin, y_end):
            n = n+1
            url = '{0}?time={1}&level={2}&x={3}&y={4}'.format(
                api_url, ts, level, x, y)
            print(n, url)
            png_fn = 'map/small/l{2}_x{0}_y{1}.png'.format(x, y, level)
            urlretrieve(url, png_fn)
            try:
                small_image = Image.open(png_fn)
                large_image.paste(
                    small_image, ((x-x_begin)*SMALL_SIZE_X, (y_end-1-y)*SMALL_SIZE_Y))
            except PIL.UnidentifiedImageError as err:
                print(err)
                n_err = n_err + 1
                # if n_err >= 5:
                #   print("Too many errors, exit...")
                #   exit()

    fn_result = 'map/L{0}_{1}.png'.format(level, tm_str)
    large_image.save(fn_result)
    print('Small images total: {0}, error: {1}'.format(n, n_err))
    return fn_result


if __name__ == '__main__':
    print('Traffic map save to:', get_traffic_map(9))
