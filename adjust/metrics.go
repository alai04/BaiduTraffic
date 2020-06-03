package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

// metrics:
// adjust_maps_total  调整的图片总数，标签:level
// adjust_maps_failure  调整失败的图片数，标签:level,stage
// adjust_consuming_total_seconds  总耗时，标签:level

const subsystem = "adjust"

var (
	adjustMapsTotal = promauto.NewCounterVec(prometheus.CounterOpts{
		Subsystem: subsystem,
		Name:      "maps_total",
		Help:      "调整成功的图片总数",
	}, []string{"level"})
	adjustMapsFailure = promauto.NewCounterVec(prometheus.CounterOpts{
		Subsystem: subsystem,
		Name:      "maps_failure",
		Help:      "调整失败的图片数",
	}, []string{"level", "stage"})
	adjustConsumingTotalSeconds = promauto.NewSummaryVec(prometheus.SummaryOpts{
		Subsystem: subsystem,
		Name:      "consuming_total_seconds",
		Help:      "调整图片总耗时",
	}, []string{"level"})
)

func init() {
	prometheus.MustRegister(adjustMapsTotal)
	prometheus.MustRegister(adjustMapsFailure)
	prometheus.MustRegister(adjustConsumingTotalSeconds)
}

func metricsServe(port string) {
	// Initial some counter
	adjustMapsFailure.With(prometheus.Labels{"level": "9", "stage": "OPEN"}).Add(0)

	http.Handle("/metrics", promhttp.Handler())
	log.Printf("Starting metrics server at http://localhost:%s/metrics", port)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%s", port), nil))
}
