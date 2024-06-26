package workutils

import (
	"slices"

	"k8s.io/client-go/kubernetes/scheme"
	workapiv1 "open-cluster-management.io/api/work/v1"
)

// Remove removes the resource with the given APIVersion, Kind, Name, and Namespace
func (client *WorkUtilsClient) Remove(work workapiv1.ManifestWork, resource Resource) (workapiv1.ManifestWork, error) {
	err := validate(resource)
	if err != nil {
		return work, err
	}

	manifests := work.Spec.Workload.Manifests

	for i, manifest := range manifests {
		o, gvk, err := decode(manifest.Raw, scheme.Scheme)
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
				work.Spec.Workload.Manifests = slices.Delete(manifests, i, i+1)
				return work, nil
			}
			namespace, err := getNamespaceFromObject(o)
			if err != nil {
				return work, nil
			}
			if namespace == resource.Namespace {
				work.Spec.Workload.Manifests = slices.Delete(manifests, i, i+1)
				return work, nil
			}
		}
	}
	return work, nil
}
