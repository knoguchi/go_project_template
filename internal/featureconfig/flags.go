package featureconfig

import (
	"github.com/urfave/cli/v2"
)

var (
	devModeFlag = &cli.BoolFlag{
		Name:  "dev",
		Usage: "Enable experimental features still in development. These features may not be stable.",
	}
	kafkaBootstrapServersFlag = &cli.StringFlag{
		Name:  "kafka-url",
		Usage: "Stream attestations and blocks to specified kafka servers. This field is used for bootstrap.servers kafka config field.",
	}
)

// devModeFlags holds list of flags that are set when development mode is on.
var devModeFlags = []cli.Flag{
}

// MyCmdFlags contains a list of all the feature flags that apply to the beacon-chain client.
var MyCmdFlags = []cli.Flag{
	devModeFlag,
	kafkaBootstrapServersFlag,
}
