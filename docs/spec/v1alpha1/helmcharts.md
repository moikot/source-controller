# Helm Charts

The `HelmChart` API defines a source for Helm chart artifacts coming
from [`HelmRepository` sources](helmrepositories.md). The resource
exposes the latest pulled chart for the defined version as an artifact.

## Specification

Helm chart:

```go
// HelmChartSpec defines the desired state of a Helm chart source.
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
```

### Reference types

```go
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
```

### Status

```go
// HelmChartStatus defines the observed state of the HelmChart.
type HelmChartStatus struct {
	// +optional
	Conditions []SourceCondition `json:"conditions,omitempty"`

	// URL is the download link for the last chart fetched.
	// +optional
	URL string `json:"url,omitempty"`

	// Artifact represents the output of the last successful chart sync.
	// +optional
	Artifact *Artifact `json:"artifact,omitempty"`
}
```

### Condition reasons

```go
const (
	// ChartPullFailedReason represents the fact that the pull of the
	// given Helm chart failed.
	ChartPullFailedReason string = "ChartPullFailed"

	// ChartPullSucceededReason represents the fact that the pull of
	// the given Helm chart succeeded.
	ChartPullSucceededReason string = "ChartPullSucceeded"

	// ChartPackageFailedReason represent the fact that the package of
	// the Helm chart failed.
	ChartPackageFailedReason string = "ChartPackageFailed"

	// ChartPackageSucceededReason represents the fact that the package of
	// the Helm chart succeeded.
	ChartPackageSucceededReason string = "ChartPackageSucceeded"
)
```

## Spec examples

Pull a specific chart version every five minutes:

```yaml
apiVersion: source.toolkit.fluxcd.io/v1alpha1
kind: HelmChart
metadata:
  name: redis
spec:
  name: redis
  version: 10.5.7
  helmRepositoryRef:
    name: stable
  interval: 5m
```

Pull the latest chart version that matches the sermver range every ten minutes:

```yaml
apiVersion: source.toolkit.fluxcd.io/v1alpha1
kind: HelmChart
metadata:
  name: redis
spec:
  name: redis
  version: ^10.0.0
  helmRepositoryRef:
    name: stable
  interval: 10m
```

## Status examples

Successful chart pull:

```yaml
status:
  url: http://<host>/helmchart/default/redis/redis-10.5.7.tgz
  conditions:
    - lastTransitionTime: "2020-04-10T09:34:45Z"
      message: Helm chart is available at /data/helmchart/default/redis/redis-10.5.7.tgz
      reason: ChartPullSucceeded
      status: "True"
      type: Ready
```

Failed chart pull:

```yaml
status:
  conditions:
    - lastTransitionTime: "2020-04-10T09:34:45Z"
      message: 'invalid chart URL format'
      reason: ChartPullFailed
      status: "False"
      type: Ready
```

Wait for ready condition:

```bash
kubectl wait helmchart/redis --for=condition=ready --timeout=1m
```
