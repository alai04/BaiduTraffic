package amap

import (
	"html/template"
	"log"
	"net/http"

	"github.com/gorilla/schema"
)

var (
	tplA    *template.Template
	decoder = schema.NewDecoder()
)

func init() {
	tplA = template.Must(template.New("amap.html.template").Parse(tplAText))
}

// Amap implement http.Handler, use to return amap
type Amap struct {
}

type mapParams struct {
	Lng  float64 `schema:"lng"`
	Lat  float64 `schema:"lat"`
	Zoom int     `schema:"z"`
	Base bool    `schema:"b"`
}

func (a Amap) getName() string {
	return "amap"
}

// ServeHTTP implement http.Handler
func (a Amap) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	log.Printf("got: %v", *r)
	mp := mapParams{
		Lng:  121.484,
		Lat:  31.216,
		Zoom: 16,
		Base: false,
	}

	err := decoder.Decode(&mp, r.URL.Query())
	if err != nil {
		log.Println("Error in GET parameters : ", err)
	} else {
		log.Println("GET parameters : ", mp)
	}

	w.Header().Set("Content-Type", "text/html")
	err = tplA.Execute(w, mp)
	if err != nil {
		log.Printf("Template execute error: %v", err)
	}
}

const tplAText = `<!DOCTYPE html>
<html>
  <head>
    <meta charset="utf-8" />
    <meta http-equiv="X-UA-Compatible" content="IE=edge" />
    <meta
      name="viewport"
      content="initial-scale=1.0, user-scalable=no, width=device-width"
    />
    <title>实时路况图层</title>
    <link
      rel="stylesheet"
      href="https://a.amap.com/jsapi_demos/static/demo-center/css/demo-center.css"
    />
    <style>
      html,
      body,
      #container {
        margin: 0;
        padding: 0;
        width: 100%;
        height: 100%;
      }
    </style>
  </head>
  <body>
    <div id="container"></div>
    <div id="tips" class="info">地图正在加载</div>
    <script src="https://webapi.amap.com/maps?v=1.4.15&key=1eae5203e9e7140c6ece93c745e422b5"></script>
    <script>
      var stdLayer = new AMap.TileLayer({
        visible: {{ .Base }},
        zIndex: 0,
      })

      //实时路况图层
      var trafficLayer = new AMap.TileLayer.Traffic({
        zIndex: 10,
      });

      var map = new AMap.Map("container", {
        resizeEnable: false,
        layers: [stdLayer, trafficLayer],
        center: [{{ .Lng }}, {{ .Lat }}],
        zoom: {{ .Zoom }},
      });

      map.on('complete', function() {
				setTimeout(function() {
					document.getElementById('tips').style.display = "none";
				}, 1000);
      });
    </script>
  </body>
</html>
`
