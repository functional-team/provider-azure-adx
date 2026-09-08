/*
Copyright 2026 The provider-azure-adx Authors.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package kusto

import (
	"time"

	"github.com/prometheus/client_golang/prometheus"
)

var (
	commandsTotal = prometheus.NewCounterVec(prometheus.CounterOpts{
		Name: "adx_commands_total",
		Help: "Kusto management commands issued by the provider, by managed resource kind, operation and result class.",
	}, []string{"kind", "op", "class"})
	commandDuration = prometheus.NewHistogramVec(prometheus.HistogramOpts{
		Name:    "adx_command_duration_seconds",
		Help:    "Duration of Kusto management commands.",
		Buckets: []float64{0.05, 0.1, 0.25, 0.5, 1, 2.5, 5, 10, 30, 60},
	}, []string{"kind", "op"})
	commandsThrottled = prometheus.NewCounterVec(prometheus.CounterOpts{
		Name: "adx_commands_throttled_total",
		Help: "Kusto management commands rejected by the cluster because of throttling.",
	}, []string{"kind"})
)

// RegisterMetrics registers the client metrics with reg. Call once from main.
func RegisterMetrics(reg prometheus.Registerer) error {
	for _, c := range []prometheus.Collector{commandsTotal, commandDuration, commandsThrottled} {
		if err := reg.Register(c); err != nil {
			if _, ok := err.(prometheus.AlreadyRegisteredError); !ok { //nolint:errorlint // prometheus returns the struct by value
				return err
			}
		}
	}
	return nil
}

func observe(op Op, class string, d time.Duration) {
	commandsTotal.WithLabelValues(op.Kind, op.Op, class).Inc()
	commandDuration.WithLabelValues(op.Kind, op.Op).Observe(d.Seconds())
	if class == "Throttled" {
		commandsThrottled.WithLabelValues(op.Kind).Inc()
	}
}
