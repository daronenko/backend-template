package metrics

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

const (
	Namespace = "backendtemplate"
)

var (
	JustCounter = promauto.NewCounterVec(prometheus.CounterOpts{
		Name: prometheus.BuildFQName(Namespace, "http", "just_counter"),
		Help: "Just counter",
	}, []string{"msg"})
)
