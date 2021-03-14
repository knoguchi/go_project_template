package configsvc

import (
	"encoding/json"
	"fmt"
	"github.com/knoguchi/go_project_template/services"
	"github.com/knoguchi/go_project_template/services/kafka"
	"github.com/knoguchi/go_project_template/services/webservice"
)

// Constants here hold some messages
const (
	ErrExchangeNameEmpty                       = "exchange #%d name is empty"
	ErrExchangeAvailablePairsEmpty             = "exchange %s available pairs is empty"
	ErrExchangeEnabledPairsEmpty               = "exchange %s enabled pairs is empty"
	ErrExchangeBaseCurrenciesEmpty             = "exchange %s base currencies is empty"
	ErrExchangeNotFound                        = "exchange %s not found"
	ErrNoEnabledExchanges                      = "no exchanges enabled"
	ErrCryptocurrenciesEmpty                   = "cryptocurrencies variable is empty"
	ErrFailureOpeningConfig                    = "fatal error opening %s file. Error: %s"
	ErrCheckingConfigValues                    = "fatal error checking config values. Error: %s"
	ErrSavingConfigBytesMismatch               = "config file %q bytes comparison doesn't match, read %s expected %s"
	WarningWebserverCredentialValuesEmpty      = "webserver support disabled due to empty Username/Password values"
	WarningWebserverListenAddressInvalid       = "webserver support disabled due to invalid listen address"
	WarningExchangeAuthAPIDefaultOrEmptyValues = "exchange %s authenticated API support disabled due to default/empty APIKey/Secret/ClientID values"
	WarningPairsLastUpdatedThresholdExceeded   = "exchange %s last manual update of available currency pairs has exceeded %d days. Manual update required!"
)

type ConfigSvc struct {
	services.Service
	MainConfig MainConfig
}

// MainConfig is the overarching object that holds all the information
// Golang can marshal Services to JSON, but it can't unmarshal JSON to Services
// Hence -services.  See _MainConfig for unmarshal
type MainConfig struct {
	Reload     bool                               `json:"reload"`
	SaveOnExit bool                               `json:"save_on_exit"`
	Services   map[string]services.IServiceConfig `json:"-services,omitempty"` // can marshal, but can't unmarshal
}

// Rest of the code is for unmarshaling Services
type _MainConfig struct {
	Services   map[string]services.IServiceConfig `json:"services,omitempty"`
}
func (c *_MainConfig) UnmarshalJSON(data []byte) error {
	// store into generic object map
	var objmap map[string]*json.RawMessage
	if err := json.Unmarshal(data, &objmap); err != nil {
		return err
	}

	// services needs special handling because unmarshaller can't guess the type
	var svcs map[string]*json.RawMessage
	if err := json.Unmarshal(*objmap["services"], &svcs); err != nil {
		return err
	}
	services := map[string]services.IServiceConfig{}
	for k := range svcs {
		switch k {
		case "kafka":
			x := &kafka.KafkaServiceConfig{}
			json.Unmarshal(*svcs[k], x)
			services[k] = x
		case "webservice":
			y := &webservice.WebServiceConfig{}
			json.Unmarshal(*svcs[k], y)
			services[k] = y
		}
		fmt.Printf("LOOP %#v\n", k)
	}
	c.Services = services

	return nil
}
