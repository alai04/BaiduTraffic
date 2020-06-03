package main

import (
	"flag"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	ginprometheus "github.com/zsais/go-gin-prometheus"
)

const (
	dateTimeFormat = "20060102150405"
)

var (
	fileDir   string
	urlTempls = []string{"/traffic_url", "/traffic_map"}
)

func main() {
	r := gin.Default()
	r.Use(cors.Default())

	p := ginprometheus.NewPrometheus("filesrv")
	p.ReqCntURLLabelMappingFn = func(c *gin.Context) string {
		url := c.Request.URL.Path
		for _, urlTempl := range urlTempls {
			if strings.HasPrefix(url, urlTempl) {
				return urlTempl
			}
		}
		return url
	}
	p.Use(r)

	r.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})
	r.GET(urlTempls[0], handleTraffic)
	r.Static(urlTempls[1], fileDir)

	r.Run() // listen and serve on 0.0.0.0:8080
}

func handleTraffic(c *gin.Context) {
	levelParam := c.Query("level")
	level, err := strconv.Atoi(levelParam)
	if err != nil || level < 9 || level > 12 {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"msg":     fmt.Sprintf("error deviation: %q", levelParam),
		})
		return
	}

	timeParam := c.DefaultQuery("time", time.Now().Format(dateTimeFormat))
	var tm time.Time
	tm, err = time.Parse(dateTimeFormat, timeParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"msg":     err.Error(),
		})
		return
	}

	deviationParam := c.DefaultQuery("deviation", "300")
	deviation, err := strconv.Atoi(deviationParam)
	if err != nil || deviation < 0 || deviation > 900 {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"msg":     fmt.Sprintf("error deviation: %q", deviationParam),
		})
		return
	}

	ff := fileFinder{level, tm, deviation}
	obj, err := ff.findFile()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"msg":     err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"msg":     "",
		"obj":     obj,
	})
}

func init() {
	flag.StringVar(&fileDir, "d", "./", "directory for serve")
	flag.Parse()
}
