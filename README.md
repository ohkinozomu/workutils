# workutils

Utilities for Open Cluster Management ManifestWork

```go
func Get(work workapiv1.ManifestWork, resource Resource) (runtime.Object, error)

func Add(work workapiv1.ManifestWork, obj runtime.Object) (workapiv1.ManifestWork, error)

func Update(work workapiv1.ManifestWork, obj runtime.Object) (workapiv1.ManifestWork, error)

func Remove(work workapiv1.ManifestWork, resource Resource) (workapiv1.ManifestWork, error)
```