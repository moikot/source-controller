# Helm Charts

The `HelmChart` API defines a source for Helm chart artifacts coming
from [`HelmRepository` sources](helmrepositories.md). The resource
exposes the latest pulled chart for the defined version as an artifact.

## Specification

Helm chart:

```go
// HelmChartSpec defines the desired state of a Helm chart source.
type HelmChartSpec struct {
	// The name of the Helm chart, as made available by the referenced
	// Helm repository.
	// +required
	Name string `json:"name"`

	// The chart version semver expression, defaults to latest when
	// omitted.
	// +optional
	Version string `json:"version,omitempty"`

	// The name of the HelmRepository the chart is available at.
	// +required
	HelmRepositoryRef v1.LocalObjectReference `json:"helmRepositoryRef"`

	// The interval at which to check the referenced HelmRepository index
	// for updates.
	// +required
	Interval metav1.Duration `json:"interval"`
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
