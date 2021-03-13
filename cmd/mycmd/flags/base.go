package flags

import "github.com/urfave/cli/v2"

// define global cli flags
var (
	// SetGCPercent is the percentage of current live allocations at which the garbage collector is to run.
	SetGCPercent = &cli.IntFlag{
		Name:  "gc-percent",
		Usage: "The percentage of freshly allocated data to live data on which the gc will be run again.",
		Value: 100,
	}

	MyGlobalFlag = &cli.BoolFlag{
		Name:  "my-global-flag",
		Usage: "global flag",
		Value: true,
	}
)
