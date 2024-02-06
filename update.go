package workutils

import (
	"k8s.io/apimachinery/pkg/runtime"
	workapiv1 "open-cluster-management.io/api/work/v1"
)

// Update updates ManifestWork
func Update(work workapiv1.ManifestWork, obj runtime.Object) (workapiv1.ManifestWork, error) {
	manifests := work.Spec.Workload.Manifests
	group, version, kind, err := getGVKFromObject(obj)
	if err != nil {
		return work, err
	}
	objName, err := getNameFromObject(obj)
	if err != nil {
		return work, err
	}
	objNamespace, err := getNamespaceFromObject(obj)
	if err != nil {
		return work, err
	}

	for i, manifest := range manifests {
		o, gvk, err := decode(manifest.Raw)
		if err != nil {
			return work, err
		}
		if gvk.Group != "" {
			if gvk.Group != group {
				continue
			}
		} else {
			if gvk.Version != version {
				continue
			}
		}
		if gvk.Kind != kind {
			continue
		}

		name, err := getNameFromObject(o)
		if err != nil {
			return work, err
		}

		rawExtension, err := objToRawExtension(obj)
		if err != nil {
			return work, err
		}

		if name == objName {
			if objNamespace == "" {
				manifests[i] = workapiv1.Manifest{
					RawExtension: rawExtension,
				}
				work.Spec.Workload.Manifests = manifests
				return work, nil
			}
			namespace, err := getNamespaceFromObject(o)
			if err != nil {
				return work, nil
			}
			if namespace == objNamespace {
				manifests[i] = workapiv1.Manifest{
					RawExtension: rawExtension,
				}
				work.Spec.Workload.Manifests = manifests
				return work, nil
			}
		}
	}
	return work, nil
}
