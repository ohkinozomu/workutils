package workutils

import (
	apiextensionsv1 "k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1"
	apiextensionsv1beta1 "k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1beta1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/client-go/kubernetes/scheme"
)

var wellknownScheme *runtime.Scheme

func init() {
	wellknownScheme = scheme.Scheme
	apiextensionsv1.AddToScheme(wellknownScheme)
	apiextensionsv1beta1.AddToScheme(wellknownScheme)
}
