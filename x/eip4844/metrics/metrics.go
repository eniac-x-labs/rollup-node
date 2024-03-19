package metrics

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/collectors"

	"github.com/eniac-x-labs/rollup-node/metrics"
	txmetrics "github.com/eniac-x-labs/rollup-node/x/eip4844/txmgr/metrics"
)

const Namespace = "eip-4844-rollup"

type Metricer interface {
	txmetrics.TxMetricer
}

type Metrics struct {
	ns       string
	Registry *prometheus.Registry
	Factory  metrics.Factory
	txmetrics.TxMetrics
}

var _ Metricer = (*Metrics)(nil)

func NewMetrics(procName string) *Metrics {
	if procName == "" {
		procName = "default"
	}
	ns := Namespace + "_" + procName

	registry := NewRegistry()
	factory := metrics.With(registry)

	return &Metrics{
		ns:       ns,
		Registry: registry,
		Factory:  factory,

		TxMetrics: txmetrics.MakeTxMetrics(ns, factory),
	}
}

type RegistryMetricer interface {
	Registry() *prometheus.Registry
}

func NewRegistry() *prometheus.Registry {
	registry := prometheus.NewRegistry()
	registry.MustRegister(collectors.NewProcessCollector(collectors.ProcessCollectorOpts{}))
	registry.MustRegister(collectors.NewGoCollector())
	return registry
}
