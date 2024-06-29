# workutils

Utilities for Open Cluster Management ManifestWork

## Motivation

OCM ManifestWork is an array of Kubernetes resources.

I found the code for performing get, add, update, and remove operations on resources within ManifestWork a bit complex, so I thought about creating utilities for these tasks.

## Usage

You first need to create a WorkUtilsClient. It requires passing `*runtime.Scheme` as an argument.

```go
import(
    "github.com/ohkinozomu/workutils"
    "k8s.io/client-go/kubernetes/scheme"
)

    client := workutils.NewWorkUtilsClient(scheme.Scheme)
```

If resources not included in `scheme.Scheme` are present in ManifestWork, `AddToScheme` for those resources is necessary.

For example, in the case of CRD:

```go
import (
	apiextensionsv1 "k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1"
	"k8s.io/apimachinery/pkg/runtime"
)

	s = scheme.Scheme
	apiextensionsv1.AddToScheme(s)
	client := workutils.NewWorkUtilsClient(s)
```

After creating the client, you can perform Get, Add, Update, and Remove operations.

```go
func (client *WorkUtilsClient) Get(work workapiv1.ManifestWork, resource Resource) (runtime.Object, error)

func (client *WorkUtilsClient) Add(work workapiv1.ManifestWork, obj runtime.Object) (workapiv1.ManifestWork, error)

func (client *WorkUtilsClient) Update(work workapiv1.ManifestWork, obj runtime.Object) (workapiv1.ManifestWork, error)

func (client *WorkUtilsClient) Remove(work workapiv1.ManifestWork, resource Resource) (workapiv1.ManifestWork, error)
```