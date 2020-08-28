/*
Copyright 2020 The Flux CD contributors.

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

package v1alpha1

import (
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

const HelmChartKind = "HelmChart"

// HelmChartSpec defines the desired state of a Helm chart.
type HelmChartSpec struct {
	// The name or path the Helm chart is available at in the SourceRef.
	// +required
	Chart string `json:"chart"`

	// The chart version semver expression, ignored for charts from GitRepository
	// sources. Defaults to latest when omitted.
	// +optional
	Version string `json:"version,omitempty"`

	// The reference to the Source the chart is available at.
	// +required
	SourceRef LocalHelmChartSourceReference `json:"sourceRef"`

	// The interval at which to check the Source for updates.
	// +required
	Interval metav1.Duration `json:"interval"`
}

// LocalHelmChartSourceReference contains enough information to let you locate the
// typed referenced object at namespace level.
type LocalHelmChartSourceReference struct {
	// APIVersion of the referent.
	// +optional
	APIVersion string `json:"apiVersion,omitempty"`

	// Kind of the referent, valid values are ('HelmRepository', 'GitRepository').
	// +kubebuilder:validation:Enum=HelmRepository;GitRepository
	// +required
	Kind string `json:"kind"`

	// Name of the referent.
	// +required
	Name string `json:"name"`
}

// HelmChartStatus defines the observed state of the HelmChart.
type HelmChartStatus struct {
	// +optional
	Conditions []SourceCondition `json:"conditions,omitempty"`

	// URL is the download link for the last chart pulled.
	// +optional
	URL string `json:"url,omitempty"`

	// Artifact represents the output of the last successful chart sync.
	// +optional
	Artifact *Artifact `json:"artifact,omitempty"`
}

const (
	// ChartPullFailedReason represents the fact that the pull of the
	// Helm chart failed.
	ChartPullFailedReason string = "ChartPullFailed"

	// ChartPullSucceededReason represents the fact that the pull of
	// the Helm chart succeeded.
	ChartPullSucceededReason string = "ChartPullSucceeded"

	// ChartPackageFailedReason represent the fact that the package of
	// the Helm chart failed.
	ChartPackageFailedReason string = "ChartPackageFailed"

	// ChartPackageSucceededReason represents the fact that the package of
	// the Helm chart succeeded.
	ChartPackageSucceededReason string = "ChartPackageSucceeded"
)

// HelmReleaseProgressing resets any failures and registers progress toward reconciling the given HelmRelease
// by setting the ReadyCondition to ConditionUnknown for ProgressingReason.
func HelmChartProgressing(chart HelmChart) HelmChart {
	chart.Status.URL = ""
	chart.Status.Artifact = nil
	chart.Status.Conditions = []SourceCondition{}
	SetHelmChartCondition(&chart, ReadyCondition, corev1.ConditionUnknown, ProgressingReason, "reconciliation in progress")
	return chart
}

// SetHelmChartCondition sets the given condition with the given status, reason and message
// on the HelmChart.
func SetHelmChartCondition(chart *HelmChart, condition string, status corev1.ConditionStatus, reason, message string) {
	chart.Status.Conditions = filterOutSourceCondition(chart.Status.Conditions, condition)
	chart.Status.Conditions = append(chart.Status.Conditions, SourceCondition{
		Type:               condition,
		Status:             status,
		LastTransitionTime: metav1.Now(),
		Reason:             reason,
		Message:            message,
	})
}

// HelmChartReady sets the given artifact and url on the HelmChart
// and sets the ReadyCondition to True, with the given reason and
// message. It returns the modified HelmChart.
func HelmChartReady(chart HelmChart, artifact Artifact, url, reason, message string) HelmChart {
	chart.Status.Artifact = &artifact
	chart.Status.URL = url
	SetHelmChartCondition(&chart, ReadyCondition, corev1.ConditionTrue, reason, message)
	return chart
}

// HelmChartNotReady sets the ReadyCondition on the given HelmChart
// to False, with the given reason and message. It returns the modified
// HelmChart.
func HelmChartNotReady(chart HelmChart, reason, message string) HelmChart {
	SetHelmChartCondition(&chart, ReadyCondition, corev1.ConditionFalse, reason, message)
	return chart
}

// HelmChartReadyMessage returns the message of the ReadyCondition
// with status True, or an empty string.
func HelmChartReadyMessage(chart HelmChart) string {
	for _, condition := range chart.Status.Conditions {
		if condition.Type == ReadyCondition && condition.Status == corev1.ConditionTrue {
			return condition.Message
		}
	}
	return ""
}

// GetArtifact returns the latest artifact from the source
// if present in the status sub-resource.
func (in *HelmChart) GetArtifact() *Artifact {
	return in.Status.Artifact
}

// GetInterval returns the interval at which the source is updated.
func (in *HelmChart) GetInterval() metav1.Duration {
	return in.Spec.Interval
}

// +genclient
// +genclient:Namespaced
// +kubebuilder:object:root=true
// +kubebuilder:subresource:status
// +kubebuilder:printcolumn:name="Chart",type=string,JSONPath=`.spec.chart`
// +kubebuilder:printcolumn:name="Version",type=string,JSONPath=`.spec.version`
// +kubebuilder:printcolumn:name="Source Kind",type=string,JSONPath=`.spec.sourceRef.kind`
// +kubebuilder:printcolumn:name="Source Name",type=string,JSONPath=`.spec.sourceRef.name`
// +kubebuilder:printcolumn:name="Ready",type="string",JSONPath=".status.conditions[?(@.type==\"Ready\")].status",description=""
// +kubebuilder:printcolumn:name="Status",type="string",JSONPath=".status.conditions[?(@.type==\"Ready\")].message",description=""
// +kubebuilder:printcolumn:name="Age",type="date",JSONPath=".metadata.creationTimestamp",description=""

// HelmChart is the Schema for the helmcharts API
type HelmChart struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   HelmChartSpec   `json:"spec,omitempty"`
	Status HelmChartStatus `json:"status,omitempty"`
}

// +kubebuilder:object:root=true

// HelmChartList contains a list of HelmChart
type HelmChartList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []HelmChart `json:"items"`
}

func init() {
	SchemeBuilder.Register(&HelmChart{}, &HelmChartList{})
}
