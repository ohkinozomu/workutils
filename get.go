package workutils

import (
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/client-go/kubernetes/scheme"
	workapiv1 "open-cluster-management.io/api/work/v1"
)

// Get returns the resource with the given APIVersion, Kind, Name, and Namespace
func (client *WorkUtilsClient) Get(work workapiv1.ManifestWork, resource Resource) (runtime.Object, error) {
	err := validate(resource)
	if err != nil {
		return nil, err
	}

	manifests := work.Spec.Workload.Manifests
	for _, manifest := range manifests {
		obj, gvk, err := decode(manifest.Raw, scheme.Scheme)
		if err != nil {
			return nil, err
		}
		if gvk.Group != resource.Group {
			continue
		}
		if gvk.Version != resource.Version {
			continue
		}
		if gvk.Kind != resource.Kind {
			continue
		}

		name, err := getNameFromObject(obj)
		if err != nil {
			return nil, err
		}

		if name == resource.Name {
			if resource.Namespace == "" {
				return obj, nil
			}
			namespace, err := getNamespaceFromObject(obj)
			if err != nil {
				return nil, err
			}
			if namespace == resource.Namespace {
				return obj, nil
			}
		}
	}
	return nil, nil
}
