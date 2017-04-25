/*
Copyright 2016 The Kubernetes Authors.

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

package metrics

import (
	"time"

	"github.com/prometheus/client_golang/prometheus"
)

var (
	lastTimestamp = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Namespace: "cluster_autoscaler",
			Name:      "last_time_seconds",
			Help:      "Last time CA run some main loop fragment.",
		}, []string{"main"},
	)

	lastDuration = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Namespace: "cluster_autoscaler",
			Name:      "last_duration_microseconds",
			Help:      "Time spent in last main loop fragments in microseconds.",
		}, []string{"main"},
	)

	duration = prometheus.NewSummaryVec(
		prometheus.SummaryOpts{
			Namespace: "cluster_autoscaler",
			Name:      "duration_microseconds",
			Help:      "Time spent in main loop fragments in microseconds.",
		}, []string{"main"},
	)

	nodegroupminstate = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Namespace: "cluster_autoscaler",
			Name:      "node_group_min_spec",
			Help:      "Current minimum bound of the node group.",
		}, []string{
			"node_group",
		},
	)

	nodegroupmaxstate = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Namespace: "cluster_autoscaler",
			Name:      "node_group_max_spec",
			Help:      "Current maximum bound of the node group.",
		}, []string{
			"node_group",
		},
	)

	nodegroupstate = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Namespace: "cluster_autoscaler",
			Name:      "node_group_size",
			Help:      "Current size of the node group.",
		}, []string{
			"node_group",
		},
	)

	scalefailures = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Namespace: "cluster_autoscaler",
			Name:      "node_group_scaling_failures",
			Help:      "Current size of the node group.",
		}, []string{
			"node_group",
			"type",
		},
	)
)

func init() {
	prometheus.MustRegister(duration)
	prometheus.MustRegister(lastDuration)
	prometheus.MustRegister(lastTimestamp)
}

func durationToMicro(start time.Time) float64 {
	return float64(time.Now().Sub(start).Nanoseconds() / 1000)
}

// UpdateDuration records the duration of the step identified by the label
func UpdateDuration(label string, start time.Time) {
	duration.WithLabelValues(label).Observe(durationToMicro(start))
	lastDuration.WithLabelValues(label).Set(durationToMicro(start))
}

// UpdateLastTime records the time the step identified by the label was started
func UpdateLastTime(label string) {
	lastTimestamp.WithLabelValues(label).Set(float64(time.Now().Unix()))
}

// UpdateNodeGroupMinState records the current minimum size of a given node group
func UpdateNodeGroupMinState(nodegroup string, min int) {
	nodegroupminstate.WithLabelValues(nodegroup).Set(float64(min))
}

// UpdateNodeGroupMaxState records the current maximum size of a given node group
func UpdateNodeGroupMaxState(nodegroup string, max int) {
	nodegroupmaxstate.WithLabelValues(nodegroup).Set(float64(max))
}

// UpdateNodeGroupState records the current size of a given node group
func UpdateNodeGroupState(nodegroup string, current int) {
	nodegroupstate.WithLabelValues(nodegroup).Set(float64(current))
}

// UpdateNodeRemoved decriments the nodegroup size
func UpdateNodeRemoved(nodegroup string) {
	nodegroupstate.WithLabelValues(nodegroup).Dec()
}

// UpdateNodeAdded incriments the nodegroup size
func UpdateNodeAdded(nodegroup string) {
	nodegroupstate.WithLabelValues(nodegroup).Inc()
}

// UpdateScaleFailures inc the counter of failures of a nodegroup and type
func UpdateScaleFailures(nodegroup string, failuretype string) {
	scalefailures.WithLabelValues(nodegroup, failuretype).Inc()
}
