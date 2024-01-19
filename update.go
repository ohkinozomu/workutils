package workutils

import (
	"k8s.io/apimachinery/pkg/runtime"
	workapiv1 "open-cluster-management.io/api/work/v1"
)

// Update updates ManifestWork
func Update(work workapiv1.ManifestWork, resource Resource, obj runtime.Object) (workapiv1.ManifestWork, error) {
	err := validate(resource)
	if err != nil {
		return work, err
	}

	manifests := work.Spec.Workload.Manifests

	for i, manifest := range manifests {
		o, gvk, err := decode(manifest.Raw)
		if err != nil {
			return work, err
		}
		if gvk.Group != "" {
			if gvk.Group != resource.Group {
				continue
			}
		} else {
			if gvk.Version != resource.Version {
				continue
			}
		}
		if gvk.Kind != resource.Kind {
			continue
		}

		name, err := getNameFromObject(o)
		if err != nil {
			return work, err
		}

		if name == resource.Name {
			if resource.Namespace == "" {
				manifests[i] = workapiv1.Manifest{
					RawExtension: runtime.RawExtension{
						Raw:    manifest.Raw,
						Object: o,
					},
				}
				work.Spec.Workload.Manifests = manifests
				return work, nil
			}
			namespace, err := getNamespaceFromObject(o)
			if err != nil {
				return work, nil
			}
			if namespace == resource.Namespace {
				manifests[i] = workapiv1.Manifest{
					RawExtension: runtime.RawExtension{
						Raw:    manifest.Raw,
						Object: o,
					},
				}
				work.Spec.Workload.Manifests = manifests
				return work, nil
			}
		}
	}
	return work, nil
}
