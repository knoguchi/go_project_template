package configsvc

import (
	"encoding/json"
	"github.com/knoguchi/go_project_template/services"
	"testing"
)

type DummyConfig struct {
	services.ServiceConfig
	MyField string `json:"my_field"`
	MyValue float64 `json:"my_value"`
}

func TestConfig(t *testing.T) {
	cfg := Config{}
	j, err := json.Marshal(cfg)
	if err != nil {
		t.Errorf("can't marshal %v", err)
	}
	if string(j) != `{"reload":false,"save_on_exit":false}` {
		t.Errorf("Unexpected output %s", j)
	}
}

func TestConfigSvc_AddServiceConfig(t *testing.T) {
	cfg := Config{
		Reload: true,
		SaveOnExit: false,
	}
	dummy := &DummyConfig{
		MyField: "foo",
		MyValue: 1.23,
	}
	cfg.Services = append(cfg.Services, dummy)

	j, err := json.Marshal(cfg)
	if err != nil {
		t.Errorf("can't marshal %v", err)
	}
	if string(j) != `{"reload":true,"save_on_exit":false,"services":[{"name":"","enabled":false,"verbose":false,"my_field":"foo","my_value":1.23}]}` {
		t.Errorf("Unexpected output %s", j)
	}
}
