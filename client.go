package workutils

import "k8s.io/apimachinery/pkg/runtime"

type WorkUtilsClient struct {
	Scheme *runtime.Scheme
}

func NewWorkUtilsClient(scheme *runtime.Scheme) *WorkUtilsClient {
	return &WorkUtilsClient{Scheme: scheme}
}
