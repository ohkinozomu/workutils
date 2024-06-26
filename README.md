# workutils

Utilities for Open Cluster Management ManifestWork

```go
func (client *WorkUtilsClient) Get(work workapiv1.ManifestWork, resource Resource) (runtime.Object, error)

func (client *WorkUtilsClient) Add(work workapiv1.ManifestWork, obj runtime.Object) (workapiv1.ManifestWork, error)

func (client *WorkUtilsClient) Update(work workapiv1.ManifestWork, obj runtime.Object) (workapiv1.ManifestWork, error)

func (client *WorkUtilsClient) Remove(work workapiv1.ManifestWork, resource Resource) (workapiv1.ManifestWork, error)
```