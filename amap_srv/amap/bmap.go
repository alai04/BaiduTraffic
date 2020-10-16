package amap

import (
	"html/template"
	"log"
	"net/http"
)

var (
	tplB *template.Template
)

func init() {
	tplB = template.Must(template.New("amap.html.template").Parse(tplBText))
}

// Bmap implement http.Handler, use to return bmap
type Bmap struct {
}

type bmapParams struct {
	Lng  float64 `schema:"lng"`
	Lat  float64 `schema:"lat"`
	Zoom int     `schema:"z"`
}

func (b Bmap) getName() string {
	return "baidu"
}

// ServeHTTP implement http.Handler
func (b Bmap) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	log.Printf("got: %v", *r)
	mp := bmapParams{
		Lng:  121.484,
		Lat:  31.216,
		Zoom: 16,
	}

	err := decoder.Decode(&mp, r.URL.Query())
	if err != nil {
		log.Println("Error in GET parameters : ", err)
	} else {
		log.Println("GET parameters : ", mp)
	}

	w.Header().Set("Content-Type", "text/html")
	err = tplB.Execute(w, mp)
	if err != nil {
		log.Printf("Template execute error: %v", err)
	}
}

const tplBText = `<!DOCTYPE html>
<html xmlns="http://www.w3.org/1999/xhtml">
  <head>
    <meta name="viewport" content="initial-scale=1.0, user-scalable=no" />
    <meta http-equiv="Content-Type" content="text/html; charset=utf-8" />
    <title>百度地图实时路况</title>
    <style type="text/css">
      html {
        height: 100%;
      }
      body {
        height: 100%;
        margin: 0px;
        padding: 0px;
      }
      #map {
        height: 100%;
        width: 100%;
      }
    </style>
    <!--引用百度地图API-->
    <script
      type="text/javascript"
      src="http://api.map.baidu.com/api?v=2.0&ak=viWC4tz2prZbv6ju00qxXAAdubAoxZ8F"
    ></script>
    <link
      href="http://api.map.baidu.com/library/TrafficControl/1.4/src/TrafficControl_min.css"
      rel="stylesheet"
      type="text/css"
    />
    <script
      type="text/javascript"
      src="http://api.map.baidu.com/library/TrafficControl/1.4/src/TrafficControl_min.js"
    ></script>
  </head>

  <body>
    <!--百度地图容器-->
    <div id="map"></div>
  </body>
  <script type="text/javascript">
    //创建和初始化地图函数：
    function initMap() {
      createMap(); //创建地图
      setMapEvent(); //设置地图事件
      addMapControl(); //向地图添加控件
      addMapOverlay(); //向地图添加覆盖物
    }
    function createMap() {
      map = new BMap.Map("map");
      map.centerAndZoom(new BMap.Point({{ .Lng }}, {{ .Lat }}), {{ .Zoom }});
    }
    function setMapEvent() {}
    function addMapOverlay() {}
    //向地图添加控件
    function addMapControl() {
      var ctrl = new BMapLib.TrafficControl({
        showPanel: true, //是否显示路况提示面板
      });
      map.addControl(ctrl);
    }
    var map;
    initMap();
    document.getElementById("tcBtn").click();
    document.getElementById("tcWrap").style.display = "none";
    setTimeout(function () {
      document.getElementById("tcBtn").style.display = "none";
    }, 2000);
  </script>
</html>
`
