package timeHist

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

var (
	histogramByGroup = promauto.NewHistogramVec(prometheus.HistogramOpts{
		Namespace: "statistics",
		Subsystem: "counter",
		Name:      "timeHist",
		Buckets: []float64{
			0.5, 0.9, 0.99,
		},
	}, []string{
		"code", "handle",
	})
)
