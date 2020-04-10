/*
Copyright 2019 The Kubernetes Authors.

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

// For an example usage see: https://github.com/kubernetes-sigs/cluster-api/blob/8206c9252cb1a8fdcd44ad405ae08a81c72c6cfb/controllers/cluster_controller.go#L177

package metrics

import (
	"github.com/prometheus/client_golang/prometheus"
	"sigs.k8s.io/controller-runtime/pkg/metrics"
)

var (
	// AzureMachineFailed is a metric that is set to 1 if the machine becomes ready
	AzureMachineReady = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "capz_machine_ready",
			Help: "Azure machine is ready if set to 1 and not if 0.",
		},
		[]string{"machine", "namespace", "cluster"},
	)

	// AzureMachineFailed is a metric that is set to 1 if the machine has failed
	AzureMachineFailed = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "capz_machine_failed",
			Help: "Azure machine is failed if set to 1 and not if 0.",
		},
		[]string{"machine", "namespace", "cluster"},
	)
)

func init() {
	metrics.Registry.MustRegister(
		AzureMachineReady,
		AzureMachineFailed,
	)
}