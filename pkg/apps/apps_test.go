package apps

import (
	"reflect"
	"testing"

	"k8s.io/api/extensions/v1beta1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/kubernetes/fake"
)

func TestNewList(t *testing.T) {
	kubeClient := fake.NewSimpleClientset()
	type args struct {
		kubeClient kubernetes.Interface
	}
	tests := []struct {
		name string
		args args
		want *List
	}{
		{
			name: "TestNewListWithNokubeClient",
			args: args{
				kubeClient: nil,
			},
			want: &List{
				kubeClient: nil,
			},
		},
		{
			name: "TestNewListWithDefaultkubeClient",
			args: args{
				kubeClient: kubeClient,
			},
			want: &List{
				kubeClient: kubeClient,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewList(tt.args.kubeClient); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewList() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestList_Get(t *testing.T) {
	kubeClient := fake.NewSimpleClientset()
	type fields struct {
		kubeClient kubernetes.Interface
		items      []ForecastleApp
		err        error
	}
	tests := []struct {
		name    string
		fields  fields
		want    []ForecastleApp
		wantErr bool
	}{
		{
			name: "TestGetForecastleApps",
			fields: fields{
				kubeClient: kubeClient,
				items: []ForecastleApp{
					ForecastleApp{
						Name:      "app",
						Icon:      "https://google.com/icon.png",
						Namespace: "test",
						URL:       "https://google.com",
					},
				},
			},
			want: []ForecastleApp{
				ForecastleApp{
					Name:      "app",
					Icon:      "https://google.com/icon.png",
					Namespace: "test",
					URL:       "https://google.com",
				},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			al := &List{
				kubeClient: tt.fields.kubeClient,
				items:      tt.fields.items,
				err:        tt.fields.err,
			}
			got, err := al.Get()
			if (err != nil) != tt.wantErr {
				t.Errorf("List.Get() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("List.Get() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_convertIngressesToForecastleApps(t *testing.T) {
	type args struct {
		ingresses []v1beta1.Ingress
	}
	tests := []struct {
		name     string
		args     args
		wantApps []ForecastleApp
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if gotApps := convertIngressesToForecastleApps(tt.args.ingresses); !reflect.DeepEqual(gotApps, tt.wantApps) {
				t.Errorf("convertIngressesToForecastleApps() = %v, want %v", gotApps, tt.wantApps)
			}
		})
	}
}
