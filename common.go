package workutils

import (
	"bytes"

	"k8s.io/apimachinery/pkg/api/meta"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apimachinery/pkg/runtime/serializer"
	"k8s.io/apimachinery/pkg/runtime/serializer/json"
	"k8s.io/client-go/kubernetes/scheme"
)

func decode(data []byte) (runtime.Object, *schema.GroupVersionKind, error) {
	scheme := scheme.Scheme
	// https://github.com/kubernetes/apimachinery/issues/102#issue-713181306
	codecs := serializer.NewCodecFactory(scheme)
	deserializer := codecs.UniversalDeserializer()
	obj, gvk, err := deserializer.Decode(data, nil, nil)
	if err != nil {
		return nil, nil, err
	}
	return obj, gvk, nil
}

func getGVKFromObject(obj runtime.Object) (string, string, string, error) {
	gvk := obj.GetObjectKind().GroupVersionKind()
	return gvk.Group, gvk.Version, gvk.Kind, nil
}

func getNameFromObject(obj runtime.Object) (string, error) {
	accessor, err := meta.Accessor(obj)
	if err != nil {
		return "", err
	}
	return accessor.GetName(), nil
}

func getNamespaceFromObject(obj runtime.Object) (string, error) {
	accessor, err := meta.Accessor(obj)
	if err != nil {
		return "", err
	}
	return accessor.GetNamespace(), nil
}

func stringToRawExtension(manifest string) (runtime.RawExtension, error) {
	obj, _, err := decode([]byte(manifest))
	if err != nil {
		return runtime.RawExtension{}, err
	}

	raw := runtime.RawExtension{
		Raw:    []byte(manifest),
		Object: obj,
	}

	return raw, nil
}

func objToRawExtension(obj runtime.Object) (runtime.RawExtension, error) {
	serializer := json.NewYAMLSerializer(json.DefaultMetaFactory, scheme.Scheme, scheme.Scheme)
	var buffer bytes.Buffer

	err := serializer.Encode(obj, &buffer)
	if err != nil {
		return runtime.RawExtension{}, err
	}

	return runtime.RawExtension{
		Raw:    buffer.Bytes(),
		Object: obj,
	}, nil
}
