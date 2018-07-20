package testutil

import (
	"k8s.io/api/extensions/v1beta1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/intstr"
)

func CreateIngress(name string) *v1beta1.Ingress {
	return &v1beta1.Ingress{
		ObjectMeta: metav1.ObjectMeta{
			Name: name,
		},
	}
}

func CreateIngressWithNamespace(name string, namespace string) *v1beta1.Ingress {
	ingress := CreateIngress(name)
	ingress.ObjectMeta.Namespace = namespace

	return ingress
}

func CreateIngressWithHost(name string, url string) *v1beta1.Ingress {
	ingress := CreateIngress(name)
	ingress.Spec.Rules = []v1beta1.IngressRule{
		{
			Host: url,
		},
	}

	return ingress
}

func AddAnnotationToIngress(ingress *v1beta1.Ingress, annotationKey string, annotationValue string) *v1beta1.Ingress {
	if ingress.Annotations == nil {
		ingress.Annotations = make(map[string]string)
	}

	ingress.Annotations[annotationKey] = annotationValue

	return ingress
}

func CreateIngressWithHostAndSubPath(name string, url string, subpath string, port string) *v1beta1.Ingress {
	ingress := CreateIngressWithHost(name, url)
	ingress.Spec.Rules[0].HTTP = &v1beta1.HTTPIngressRuleValue{
		Paths: []v1beta1.HTTPIngressPath{
			{
				Backend: v1beta1.IngressBackend{
					ServicePort: intstr.FromString(port),
				},
				Path: subpath,
			},
		},
	}

	return ingress
}

func CreateIngressWithTLSHost(name string, tlsurl string) *v1beta1.Ingress {
	ingress := CreateIngress(name)
	ingress.Spec.TLS = []v1beta1.IngressTLS{
		{
			Hosts: []string{
				tlsurl,
			},
		},
	}

	return ingress
}

func CreateIngressWithHostAndTLSHost(name string, host string, tlsurl string) *v1beta1.Ingress {
	ingress := CreateIngressWithHost(name, host)
	ingress.Spec.TLS = []v1beta1.IngressTLS{
		{
			Hosts: []string{
				tlsurl,
			},
		},
	}

	return ingress
}
