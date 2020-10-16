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
      #tips {
        padding: 0.75rem 1.25rem;
        margin-bottom: 1rem;
        border-radius: 0.25rem;
        position: fixed;
        top: 1rem;
        background-color: white;
        width: auto;
        min-width: 22rem;
        border-width: 0;
        right: 1rem;
        box-shadow: 0 2px 6px 0 rgba(114, 124, 245, 0.5);
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
    <div id="tips" class="info">地图正在加载</div>
  </body>
  <script type="text/javascript">
    var map = new BMap.Map("map");
    var traffic = new BMap.TrafficLayer();
    map.addTileLayer(traffic);
    map.addEventListener("tilesloaded", function () {
      setTimeout(function () {
        document.getElementById("tips").style.display = "none";
      }, 500);
    });
    map.centerAndZoom(new BMap.Point(121.484, 31.216), 16);
  </script>
</html>
`
