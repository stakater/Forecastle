package customapps

import (
	"reflect"
	"testing"

	"github.com/stakater/Forecastle/pkg/config"
	"github.com/stakater/Forecastle/pkg/forecastle"
)

func TestNewList(t *testing.T) {
	type args struct {
		appConfig config.Config
	}
	tests := []struct {
		name string
		args args
		want *List
	}{
		{
			name: "TestNewListWithAppConfig",
			args: args{
				appConfig: config.Config{Title: "Test"},
			},
			want: &List{
				appConfig: config.Config{Title: "Test"},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewList(tt.args.appConfig); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewList() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestList_Populate(t *testing.T) {
	type fields struct {
		appConfig config.Config
		err       error
		items     []forecastle.App
	}
	tests := []struct {
		name   string
		fields fields
		want   *List
	}{
		{
			name: "TestListPopulate",
			fields: fields{
				appConfig: config.Config{
					CustomApps: []config.CustomApp{
						config.CustomApp{
							Name:  "Test",
							URL:   "http://google.com",
							Icon:  "http://google.com",
							Group: "My Group",
						},
					},
				},
			},
			want: &List{
				appConfig: config.Config{
					CustomApps: []config.CustomApp{
						config.CustomApp{
							Name:  "Test",
							URL:   "http://google.com",
							Icon:  "http://google.com",
							Group: "My Group",
						},
					},
				},
				items: []forecastle.App{
					{
						Name:     "Test",
						URL:      "http://google.com",
						Icon:     "http://google.com",
						Group:    "My Group",
						IsCustom: true,
					},
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			al := &List{
				appConfig: tt.fields.appConfig,
				err:       tt.fields.err,
				items:     tt.fields.items,
			}
			if got := al.Populate(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("List.Populate() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestList_Get(t *testing.T) {
	type fields struct {
		appConfig config.Config
		err       error
		items     []forecastle.App
	}
	tests := []struct {
		name    string
		fields  fields
		want    []forecastle.App
		wantErr bool
	}{
		{
			name: "TestGetCustomForecastleApps",
			fields: fields{
				items: []forecastle.App{
					{
						Name:     "app",
						Icon:     "https://google.com/icon.png",
						Group:    "test",
						URL:      "https://google.com",
						IsCustom: true,
					},
				},
			},
			want: []forecastle.App{
				{
					Name:     "app",
					Icon:     "https://google.com/icon.png",
					Group:    "test",
					URL:      "https://google.com",
					IsCustom: true,
				},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			al := &List{
				appConfig: tt.fields.appConfig,
				err:       tt.fields.err,
				items:     tt.fields.items,
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

func Test_convertCustomAppsToForecastleApps(t *testing.T) {
	type args struct {
		customApps []config.CustomApp
	}
	tests := []struct {
		name     string
		args     args
		wantApps []forecastle.App
	}{
		{
			name: "TestConvertCustomAppsToForecastleAppsWithNoApps",
			args: args{
				customApps: []config.CustomApp{},
			},
			wantApps: nil,
		},
		{
			name: "TestConvertCustomAppsToForecastleApps",
			args: args{
				customApps: []config.CustomApp{
					config.CustomApp{
						Name:  "test",
						Icon:  "http://google.com/image.png",
						Group: "New",
						URL:   "http://google.com",
					},
				},
			},
			wantApps: []forecastle.App{
				forecastle.App{
					Name:     "test",
					Icon:     "http://google.com/image.png",
					Group:    "New",
					URL:      "http://google.com",
					IsCustom: true,
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if gotApps := convertCustomAppsToForecastleApps(tt.args.customApps); !reflect.DeepEqual(gotApps, tt.wantApps) {
				t.Errorf("convertCustomAppsToForecastleApps() = %v, want %v", gotApps, tt.wantApps)
			}
		})
	}
}
