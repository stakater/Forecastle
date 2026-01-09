package wrappers

import (
	"reflect"
	"testing"

	"github.com/stakater/Forecastle/v1/pkg/annotations"
	"github.com/stakater/Forecastle/v1/pkg/testutil"
	v1 "k8s.io/api/networking/v1"
)

func TestNewIngressWrapper(t *testing.T) {
	type args struct {
		ingress *v1.Ingress
	}
	tests := []struct {
		name string
		args args
		want *IngressWrapper
	}{
		{
			name: "TestNewIngressWrapper",
			args: args{
				ingress: testutil.CreateIngress("test-ingress"),
			},
			want: &IngressWrapper{
				ingress: testutil.CreateIngress("test-ingress"),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewIngressWrapper(tt.args.ingress); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewIngressWrapper() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestIngressWrapper_GetAnnotationValue(t *testing.T) {
	type fields struct {
		ingress *v1.Ingress
	}
	type args struct {
		annotationKey string
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   string
	}{
		{
			name: "TestGetAnnotationValueWithExistingAnnotation",
			fields: fields{
				ingress: testutil.AddAnnotationToIngress(testutil.CreateIngress("test-ingress"), "someannotation", "hello"),
			},
			args: args{
				annotationKey: "someannotation",
			},
			want: "hello",
		},
		{
			name: "TestGetAnnotationValueWithInvalidAnnotation",
			fields: fields{
				ingress: testutil.CreateIngress("test-ingress"),
			},
			args: args{
				annotationKey: "someannotation",
			},
			want: "",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			iw := &IngressWrapper{
				ingress: tt.fields.ingress,
			}
			if got := iw.GetAnnotationValue(tt.args.annotationKey); got != tt.want {
				t.Errorf("IngressWrapper.GetAnnotationValue() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestIngressWrapper_GetName(t *testing.T) {
	type fields struct {
		ingress *v1.Ingress
	}
	tests := []struct {
		name   string
		fields fields
		want   string
	}{
		{
			name: "TestGetName",
			fields: fields{
				ingress: testutil.CreateIngress("test-ingress"),
			},
			want: "test-ingress",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			iw := &IngressWrapper{
				ingress: tt.fields.ingress,
			}
			if got := iw.GetName(); got != tt.want {
				t.Errorf("IngressWrapper.GetName() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestIngressWrapper_GetNamespace(t *testing.T) {
	type fields struct {
		ingress *v1.Ingress
	}
	tests := []struct {
		name   string
		fields fields
		want   string
	}{
		{
			name: "TestGetNamespace",
			fields: fields{
				ingress: testutil.CreateIngressWithNamespace("test-ingress", "test"),
			},
			want: "test",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			iw := &IngressWrapper{
				ingress: tt.fields.ingress,
			}
			if got := iw.GetNamespace(); got != tt.want {
				t.Errorf("IngressWrapper.GetNamespace() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestIngressWrapper_GetURL(t *testing.T) {
	type fields struct {
		ingress *v1.Ingress
	}
	tests := []struct {
		name   string
		fields fields
		want   string
	}{
		{
			name: "IngressWithNoURL",
			fields: fields{
				ingress: testutil.CreateIngress("someIngress"),
			},
			want: "",
		},
		{
			name: "IngressWithValidHost",
			fields: fields{
				ingress: testutil.CreateIngressWithHost("someIngress1", "google.com"),
			},
			want: "http://google.com",
		},
		{
			name: "IngressWithTLSHostButNoHost",
			fields: fields{
				ingress: testutil.CreateIngressWithTLSHost("someIngress2", "google.com"),
			},
			want: "https://google.com",
		},
		{
			name: "IngressWithTLSHostAndNormalHost",
			fields: fields{
				ingress: testutil.CreateIngressWithHostAndTLSHost("someIngress2", "google.com", "google.com"),
			},
			want: "https://google.com",
		},
		{
			name: "IngressWithValidHostWithOverridenURLWithoutScheme",
			fields: fields{
				ingress: testutil.AddAnnotationToIngress(testutil.CreateIngressWithHost("someIngress1", "google.com"), annotations.ForecastleURLAnnotation, "someotherurl.com"),
			},
			want: "",
		},
		{
			name: "IngressWithValidHostWithOverridenURLWithScheme",
			fields: fields{
				ingress: testutil.AddAnnotationToIngress(testutil.CreateIngressWithHost("someIngress1", "google.com"), annotations.ForecastleURLAnnotation, "https://someotherurl.com"),
			},
			want: "https://someotherurl.com",
		},
		{
			name: "IngressWithValidHostWithOverridenInvalidURL",
			fields: fields{
				ingress: testutil.AddAnnotationToIngress(testutil.CreateIngressWithHost("someIngress1", "google.com"), annotations.ForecastleURLAnnotation, "someotherurl42"),
			},
			want: "",
		},
		{
			name: "IngressWithStatusHostnameHostButNoHostNorTLSHost",
			fields: fields{
				ingress: testutil.CreateIngressWithStatusHostnameHost("someIngress2", "google.com"),
			},
			want: "http://google.com",
		},
		{
			name: "IngressWithStatusIPHostButNoHostNorTLSHost",
			fields: fields{
				ingress: testutil.CreateIngressWithStatusIPHost("someIngress2", "1.1.1.1"),
			},
			want: "http://1.1.1.1",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			iw := &IngressWrapper{
				ingress: tt.fields.ingress,
			}
			if got := iw.GetURL(); got != tt.want {
				t.Errorf("IngressWrapper.GetURL() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestIngressWrapper_rulesExist(t *testing.T) {
	type fields struct {
		ingress *v1.Ingress
	}
	tests := []struct {
		name   string
		fields fields
		want   bool
	}{
		{
			name: "IngressWithNoRules",
			fields: fields{
				ingress: testutil.CreateIngress("someIngress"),
			},
			want: false,
		},
		{
			name: "IngressWithRule",
			fields: fields{
				ingress: testutil.CreateIngressWithHost("someIngress", "http://google.com"),
			},
			want: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			iw := &IngressWrapper{
				ingress: tt.fields.ingress,
			}
			if got := iw.rulesExist(); got != tt.want {
				t.Errorf("IngressWrapper.rulesExist() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestIngressWrapper_tryGetTLSHost(t *testing.T) {
	type fields struct {
		ingress *v1.Ingress
	}
	tests := []struct {
		name   string
		fields fields
		want   string
		want1  bool
	}{
		{
			name: "IngressWithoutTLSHost",
			fields: fields{
				ingress: testutil.CreateIngressWithHost("someIngress", "google.com"),
			},
			want:  "",
			want1: false,
		},
		{
			name: "IngressWithTLSHost",
			fields: fields{
				ingress: testutil.CreateIngressWithHostAndTLSHost("someIngress", "google.com", "google.com"),
			},
			want:  "google.com",
			want1: true,
		},
		{
			name: "IngressWithTLSAndNoHosts",
			fields: fields{
				ingress: testutil.CreateIngressWithHostAndEmptyTLSHost("someIngress", "google.com"),
			},
			want:  "",
			want1: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			iw := &IngressWrapper{
				ingress: tt.fields.ingress,
			}
			got, got1 := iw.tryGetTLSHost()
			if got != tt.want {
				t.Errorf("IngressWrapper.tryGetTLSHost() got = %v, want %v", got, tt.want)
			}
			if got1 != tt.want1 {
				t.Errorf("IngressWrapper.tryGetTLSHost() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}
}

func TestIngressWrapper_supportsTLS(t *testing.T) {
	type fields struct {
		ingress *v1.Ingress
	}
	tests := []struct {
		name   string
		fields fields
		want   bool
	}{
		{
			name: "IngressWithoutTLS",
			fields: fields{
				ingress: testutil.CreateIngressWithHost("someIngress", "google.com"),
			},
			want: false,
		},
		{
			name: "IngressWithTLS",
			fields: fields{
				ingress: testutil.CreateIngressWithHostAndTLSHost("someIngress", "google.com", "google.com"),
			},
			want: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			iw := &IngressWrapper{
				ingress: tt.fields.ingress,
			}
			if got := iw.supportsTLS(); got != tt.want {
				t.Errorf("IngressWrapper.supportsTLS() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestIngressWrapper_tryGetRuleHost(t *testing.T) {
	type fields struct {
		ingress *v1.Ingress
	}
	tests := []struct {
		name   string
		fields fields
		want   string
		want1  bool
	}{
		{
			name: "IngressWithEmptyHost",
			fields: fields{
				ingress: testutil.CreateIngressWithHost("someIngress", ""),
			},
			want:  "",
			want1: false,
		},
		{
			name: "IngressWithCorrectHost",
			fields: fields{
				ingress: testutil.CreateIngressWithHost("someIngress", "google.com"),
			},
			want:  "google.com",
			want1: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			iw := &IngressWrapper{
				ingress: tt.fields.ingress,
			}
			got, got1 := iw.tryGetRuleHost()
			if got != tt.want {
				t.Errorf("IngressWrapper.tryGetRuleHost() got = %v, want %v", got, tt.want)
			}
			if got1 != tt.want1 {
				t.Errorf("IngressWrapper.tryGetRuleHost() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}
}

func TestIngressWrapper_tryGetStatusHost(t *testing.T) {
	type fields struct {
		ingress *v1.Ingress
	}
	tests := []struct {
		name   string
		fields fields
		want   string
		want1  bool
	}{
		{
			name: "IngressWithEmptyStatusHost",
			fields: fields{
				ingress: testutil.CreateIngressWithHost("someIngress", ""),
			},
			want:  "",
			want1: false,
		},
		{
			name: "IngressWithCorrectHostnameStatusHost",
			fields: fields{
				ingress: testutil.CreateIngressWithStatusHostnameHost("someIngress", "google.com"),
			},
			want:  "google.com",
			want1: true,
		},
		{
			name: "IngressWithCorrectIPStatusHost",
			fields: fields{
				ingress: testutil.CreateIngressWithStatusIPHost("someIngress", "1.1.1.1"),
			},
			want:  "1.1.1.1",
			want1: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			iw := &IngressWrapper{
				ingress: tt.fields.ingress,
			}
			got, got1 := iw.tryGetStatusHost()
			if got != tt.want {
				t.Errorf("IngressWrapper.tryGetStatusHost() got = %v, want %v", got, tt.want)
			}
			if got1 != tt.want1 {
				t.Errorf("IngressWrapper.tryGetStatusHost() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}
}

func TestIngressWrapper_getIngressSubPath(t *testing.T) {
	type fields struct {
		ingress *v1.Ingress
	}
	tests := []struct {
		name   string
		fields fields
		want   string
	}{
		{
			name: "IngressWithoutSubPath",
			fields: fields{
				ingress: testutil.CreateIngressWithHost("someIngress", "google.com"),
			},
			want: "",
		},
		{
			name: "IngressWithSubPath",
			fields: fields{
				ingress: testutil.CreateIngressWithHostAndSubPath("someIngress", "google.com", "/test", ""),
			},
			want: "/test",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			iw := &IngressWrapper{
				ingress: tt.fields.ingress,
			}
			if got := iw.getIngressSubPath(); got != tt.want {
				t.Errorf("IngressWrapper.getIngressSubPath() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestIngressWrapper_GetGroup(t *testing.T) {
	type fields struct {
		ingress *v1.Ingress
	}
	tests := []struct {
		name   string
		fields fields
		want   string
	}{
		{
			name: "IngressWithoutGroup",
			fields: fields{
				ingress: testutil.CreateIngressWithNamespace("someIngress", "test"),
			},
			want: "test", // Namespace is normalized to lowercase
		},
		{
			name: "IngressWithGroup",
			fields: fields{
				ingress: testutil.AddAnnotationToIngress(testutil.CreateIngressWithNamespace("someIngress", "test"), annotations.ForecastleGroupAnnotation, "My Group"),
			},
			want: "my group", // Group is normalized to lowercase for case-insensitive grouping
		},
		{
			name: "IngressWithMixedCaseNamespace",
			fields: fields{
				ingress: testutil.CreateIngressWithNamespace("someIngress", "TestNamespace"),
			},
			want: "testnamespace", // Namespace is normalized to lowercase
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			iw := &IngressWrapper{
				ingress: tt.fields.ingress,
			}
			if got := iw.GetGroup(); got != tt.want {
				t.Errorf("IngressWrapper.GetGroup() = %v, want %v", got, tt.want)
			}
		})
	}
}
