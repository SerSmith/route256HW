package callCounter

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

var (
	HandlerCounter = promauto.NewCounterVec(prometheus.CounterOpts{
		Namespace: "statistics",
		Subsystem: "counter",
		Name:      "handlerCals",
	}, []string{
		"handle"})
)
