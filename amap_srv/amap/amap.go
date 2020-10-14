package amap

import (
	"context"
	"html/template"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/httptest"
	"time"

	cdp "github.com/chromedp/chromedp"
	"github.com/gorilla/schema"
)

var (
	tpl     *template.Template
	decoder = schema.NewDecoder()
)

func init() {
	tpl = template.Must(template.New("amap.html.template").Parse(tplText))
}

// Amap implement http.Handler, use to return amap
type Amap struct {
}

type mapParams struct {
	// Width  int     `schema:"w"`
	// Height int     `schema:"h"`
	Lng  float64 `schema:"lng"`
	Lat  float64 `schema:"lat"`
	Zoom int     `schema:"z"`
	Base bool    `schema:"b"`
}

func (a Amap) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	log.Printf("got: %v", *r)
	mp := mapParams{
		// Width:  2800,
		// Height: 3500,
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
	err = tpl.Execute(w, mp)
	if err != nil {
		log.Printf("Template execute error: %v", err)
	}
}

// PrepareForShot prepares the http server to return html for amap
// Prepare 1 time, shot anytime
func PrepareForShot() (url string, close func()) {
	ts := httptest.NewServer(Amap{})
	return ts.URL, ts.Close
}

// Shot captures the screenshot from amap
func Shot(url string) {
	var buf []byte
	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()
	ctx, cancel = cdp.NewContext(ctx, cdp.WithLogf(log.Printf))
	defer cancel()

	url += `?lng=121.484&lat=31.216&z=16&b=1`
	err := cdp.Run(ctx, cdp.Tasks{
		cdp.EmulateViewport(2600, 3200),
		cdp.Navigate(url),
		cdp.WaitNotVisible(`#tip`),
		cdp.CaptureScreenshot(&buf),
		screenshotSave("foo.png", &buf),
	})
	if err != nil {
		log.Fatal(err)
	}
}

func screenshotSave(fileName string, buf *[]byte) cdp.ActionFunc {
	return func(ctx context.Context) error {
		log.Printf("Write %v", fileName)
		return ioutil.WriteFile(fileName, *buf, 0644)
	}
}

const tplText = `<!DOCTYPE html>
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
    <div id="tip" class="info">地图正在加载</div>
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
          document.getElementById('tip').style.display = "none";
      });
    </script>
  </body>
</html>
`
