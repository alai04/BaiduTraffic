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
// adjust_consuming_total_seconds  总耗时，标签:level
// adjust_consuming_average_seconds  平均耗时，标签:level
// adjust_consuming_min_seconds  最小耗时，标签:level
// adjust_consuming_max_seconds  最大耗时，标签:level
// adjust_maps_average_bytes  图片平均大小，标签:level
// adjust_maps_min_bytes  图片最小大小，标签:level
// adjust_maps_max_bytes  图片最大大小，标签:level
// adjust_maps_failure  调整失败的图片数，标签:level,stage
// adjust_succ_precent  调整成功的比例，标签:level

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
	}, []string{"level"})
	adjustConsumingTotalSeconds = promauto.NewSummaryVec(prometheus.SummaryOpts{
		Subsystem: subsystem,
		Name:      "consuming_total_seconds",
		Help:      "调整图片总耗时",
	}, []string{"level"})
)

func metricsServe(port string) {
	http.Handle("/metrics", promhttp.Handler())
	log.Printf("Starting metrics server at http://localhost:%s/metrics", port)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%s", port), nil))
}
