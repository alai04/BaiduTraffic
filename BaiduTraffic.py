#!/usr/bin/env python3
#-*-coding:utf-8-*-

import time
from urllib.request import urlretrieve
import PIL
import PIL.Image as Image

def get_traffic_map(level):

  api_url = 'http://its.map.baidu.com:8002/traffic/TrafficTileService'
  tm = time.time()
  ts = int(tm * 1000)
  tm_str = time.strftime('%Y%m%d_%H%M%S', time.localtime(tm))
  x_begin_level9, x_end_level9 = 97, 102
  y_begin_level9, y_end_level9 = 25, 32
  scale = 2**(level - 9)
  x_begin, x_end = x_begin_level9 * scale, x_end_level9 * scale
  y_begin, y_end = y_begin_level9 * scale, y_end_level9 * scale
  SMALL_SIZE_X = 256
  SMALL_SIZE_Y = 256
  large_image = Image.new('RGBA', (SMALL_SIZE_X*(x_end-x_begin), SMALL_SIZE_Y*(y_end-y_begin)))

  n = 0
  n_err = 0
  for x in range(x_begin, x_end):
    for y in range(y_begin, y_end):
      n = n+1
      url = '{0}?time={1}&level={2}&x={3}&y={4}'.format(api_url, ts, level, x, y)
      print(n, url)
      png_fn = 'map/small/l{2}_x{0}_y{1}.png'.format(x, y, level)
      urlretrieve(url, png_fn)
      try:
        small_image = Image.open(png_fn)
        large_image.paste(small_image, ((x-x_begin)*SMALL_SIZE_X, (y_end-1-y)*SMALL_SIZE_Y))
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