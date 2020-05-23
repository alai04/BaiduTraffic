# 设计方案
## 总体方案

## 前端取路况图接口
### 获取最接近路况的链接
HTTP GET
http://{host}:{port}/traffic_url?level={level}&time={YYYYMMDDhhmmss}&deviation={seconds}

#### 入参说明
- `level`     必选，取值范围9-12，level每提高1级，分辨率提高2x2倍，level9的分辨率为每像素512米
- `time`      可选，获取指定时间的路况图，14位数字表示年月日时分秒，缺省为当前时间
- `deviation` 可选，可接受的时间误差，秒为单位，缺省为60，表示time前后1分钟均可

#### 出参说明 JSON
```
{
  'success': true/false,
  'msg': '成功/失败原因',
  'obj': {
    'time': 'YYYYMMDDhhmmss',
    'filename': '图片文件名',
  }
}
```

### 获取路况图片
HTTP GET
http://{host}:{port}/traffic_map/{filename}

#### 入参说明
无

#### 出参说明
PNG图片数据
