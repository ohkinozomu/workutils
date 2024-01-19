package workutils

import (
	"k8s.io/apimachinery/pkg/runtime"
	workapiv1 "open-cluster-management.io/api/work/v1"
)

// Add adds the given object to the ManifestWork
func Add(work workapiv1.ManifestWork, obj runtime.Object) (workapiv1.ManifestWork, error) {
	manifests := work.Spec.Workload.Manifests
	rawExtension, err := objToRawExtension(obj)
	if err != nil {
		return work, err
	}
	manifests = append(manifests, workapiv1.Manifest{
		RawExtension: rawExtension,
	})
	work.Spec.Workload.Manifests = manifests
	return work, nil
}
