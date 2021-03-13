package flags

import (
	log "github.com/sirupsen/logrus"
	"github.com/urfave/cli/v2"
)

// GlobalFlags specifies all the global flags for the
// beacon node.
type GlobalFlags struct {
	MyGlobalFlag bool
}

var globalConfig *GlobalFlags

// Get retrieves the global config.
func Get() *GlobalFlags {
	if globalConfig == nil {
		return &GlobalFlags{}
	}
	return globalConfig
}

// Init sets the global config equal to the config that is passed in.
func Init(c *GlobalFlags) {
	globalConfig = c
}

// ConfigureGlobalFlags initializes the global config.
// based on the provided cli context.
func ConfigureGlobalFlags(ctx *cli.Context) {
	cfg := &GlobalFlags{}
	if ctx.Bool(MyGlobalFlag.Name) {
		log.Warn("Using Head Sync flag, it starts syncing from last saved head.")
		cfg.MyGlobalFlag = true
	}
	// do some config work

	Init(cfg)
}

