package main

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"
	"path/filepath"

	"github.com/gorilla/schema"
)

var (
	tpl     *template.Template
	decoder = schema.NewDecoder()
)

func init() {
	ex, _ := os.Executable()
	exPath := filepath.Dir(ex)
	realPath, _ := filepath.EvalSymlinks(exPath)
	tplFile := filepath.Join(realPath, "amap.html.template")
	log.Printf("Using template file %q", tplFile)
	tpl = template.Must(template.New("amap.html.template").ParseFiles(tplFile))
}

type amap struct {
}

func serve() {
	// example: http://localhost:8080/amap?lng=121.484&lat=31.216&z=16
	mux := http.NewServeMux()
	mux.Handle("/amap", amap{})
	log.Printf("Listen on port %d", flagPort)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", flagPort), mux))
}

type mapParams struct {
	// Width  int     `schema:"w"`
	// Height int     `schema:"h"`
	Lng  float64 `schema:"lng"`
	Lat  float64 `schema:"lat"`
	Zoom int     `schema:"z"`
	Base bool    `schema:"b"`
}

func (a amap) ServeHTTP(w http.ResponseWriter, r *http.Request) {
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
