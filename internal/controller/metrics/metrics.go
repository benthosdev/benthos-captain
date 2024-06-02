package metrics

import (
	"github.com/prometheus/client_golang/prometheus"
	"sigs.k8s.io/controller-runtime/pkg/metrics"
)

var (
	PipelineReconciles = prometheus.NewCounterVec(prometheus.CounterOpts{
		Namespace: "benthos_captain",
		Name:      "reconciles",
		Help:      "reconciles per Benthos pipeline",
	}, []string{"pipeline"})

	PipelineFailedReconciles = prometheus.NewCounterVec(prometheus.CounterOpts{
		Namespace: "benthos_captain",
		Name:      "failed_reconciles",
		Help:      "failed reconciles per Benthos pipeline",
	}, []string{"pipeline"})
)

func init() {
	metrics.Registry.MustRegister(PipelineReconciles, PipelineFailedReconciles)
}
