package configsvc

import (
	"context"
	"github.com/knoguchi/go_project_template/services"
	"io"
	"reflect"
	"testing"
)

func TestConfigSvc_AddService(t *testing.T) {
	type fields struct {
		Service services.Service
		config  Config
	}
	type args struct {
		svc services.IService
	}
	tests := []struct {
		name   string
		fields fields
		args   args
	}{
		{
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &ConfigSvc{
				Service: tt.fields.Service,
				config:  tt.fields.config,
			}
			if c.Config != nil {
				t.Error("fail")
			}
		})
	}
}

func TestConfigSvc_CheckConfig(t *testing.T) {
	type fields struct {
		Service services.Service
		config  Config
	}
	tests := []struct {
		name    string
		fields  fields
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &ConfigSvc{
				Service: tt.fields.Service,
				config:  tt.fields.config,
			}
			if err := c.CheckConfig(); (err != nil) != tt.wantErr {
				t.Errorf("CheckConfig() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestConfigSvc_LoadConfig(t *testing.T) {
	type fields struct {
		Service services.Service
		config  Config
	}
	type args struct {
		configPath string
		dryrun     bool
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &ConfigSvc{
				Service: tt.fields.Service,
				config:  tt.fields.config,
			}
			if err := c.LoadConfig(tt.args.configPath, tt.args.dryrun); (err != nil) != tt.wantErr {
				t.Errorf("LoadConfig() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestConfigSvc_ReadConfigFromFile(t *testing.T) {
	type fields struct {
		Service services.Service
		config  Config
	}
	type args struct {
		configPath string
		dryrun     bool
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &ConfigSvc{
				Service: tt.fields.Service,
				config:  tt.fields.config,
			}
			if err := c.ReadConfigFromFile(tt.args.configPath, tt.args.dryrun); (err != nil) != tt.wantErr {
				t.Errorf("ReadConfigFromFile() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestConfigSvc_Start(t *testing.T) {
	type fields struct {
		Service services.Service
		config  Config
	}
	type args struct {
		ctx context.Context
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &ConfigSvc{
				Service: tt.fields.Service,
				config:  tt.fields.config,
			}
			if err := c.Start(tt.args.ctx); (err != nil) != tt.wantErr {
				t.Errorf("Start() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestConfigSvc_Status(t *testing.T) {
	type fields struct {
		Service services.Service
		config  Config
	}
	tests := []struct {
		name    string
		fields  fields
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &ConfigSvc{
				Service: tt.fields.Service,
				config:  tt.fields.config,
			}
			if err := c.Status(); (err != nil) != tt.wantErr {
				t.Errorf("Status() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestConfigSvc_Stop(t *testing.T) {
	type fields struct {
		Service services.Service
		config  Config
	}
	tests := []struct {
		name    string
		fields  fields
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &ConfigSvc{
				Service: tt.fields.Service,
				config:  tt.fields.config,
			}
			if err := c.Stop(); (err != nil) != tt.wantErr {
				t.Errorf("Stop() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestNew(t *testing.T) {
	tests := []struct {
		name string
		want *ConfigSvc
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := New(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("New() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestReadConfig(t *testing.T) {
	type args struct {
		configReader io.Reader
	}
	tests := []struct {
		name    string
		args    args
		want    *Config
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ReadConfig(tt.args.configReader)
			if (err != nil) != tt.wantErr {
				t.Errorf("ReadConfig() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ReadConfig() got = %v, want %v", got, tt.want)
			}
		})
	}
}
