package workutils

import (
	"bytes"
	"fmt"

	"k8s.io/apimachinery/pkg/runtime"
	workapiv1 "open-cluster-management.io/api/work/v1"
)

func (client *WorkUtilsClient) Diff(old workapiv1.ManifestWork, new workapiv1.ManifestWork) (added []runtime.Object, removed []runtime.Object, updated []runtime.Object, err error) {
	oldMap := make(map[string]runtime.Object)
	for _, manifest := range old.Spec.Workload.Manifests {
		obj, _, err := decode(manifest.Raw, client.Scheme)
		if err != nil {
			return nil, nil, nil, fmt.Errorf("failed to decode manifest in old: %v", err)
		}
		key, err := manifestKey(obj)
		if err != nil {
			return nil, nil, nil, fmt.Errorf("failed to get manifest key: %v", err)
		}
		oldMap[key] = obj
	}

	newMap := make(map[string]runtime.Object)
	for _, manifest := range new.Spec.Workload.Manifests {
		obj, _, err := decode(manifest.Raw, client.Scheme)
		if err != nil {
			return nil, nil, nil, fmt.Errorf("failed to decode manifest in new: %v", err)
		}
		key, err := manifestKey(obj)
		if err != nil {
			return nil, nil, nil, fmt.Errorf("failed to get manifest key: %v", err)
		}
		newMap[key] = obj
	}

	for key, obj1 := range oldMap {
		if obj2, found := newMap[key]; found {
			e, err := equal(client.Scheme, obj1, obj2)
			if err != nil {
				return nil, nil, nil, fmt.Errorf("failed to compare objects: %v", err)
			}
			if !e {
				updated = append(updated, obj2)
			}
			delete(newMap, key)
		} else {
			removed = append(removed, obj1)
		}
	}

	for _, obj := range newMap {
		added = append(added, obj)
	}

	return added, removed, updated, nil
}

func manifestKey(obj runtime.Object) (string, error) {
	r, err := getResourceFromObject(obj)
	if err != nil {
		return "", err
	}
	key := fmt.Sprintf("%s/%s/%s/%s/%s", r.Group, r.Version, r.Kind, r.Namespace, r.Name)
	return key, nil
}

func equal(scheme *runtime.Scheme, obj1, obj2 runtime.Object) (bool, error) {
	raw1, err := objToRawExtension(obj1, scheme)
	if err != nil {
		return false, err
	}
	raw2, err := objToRawExtension(obj2, scheme)
	if err != nil {
		return false, err
	}
	return bytes.Equal(raw1.Raw, raw2.Raw), nil
}
