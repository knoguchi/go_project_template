package main

import (
	"context"
	"errors"
	"fmt"
	"github.com/knoguchi/go_project_template/cmd"
	"github.com/knoguchi/go_project_template/cmd/mycmd/flags"
	"github.com/knoguchi/go_project_template/internal/debug"
	"github.com/knoguchi/go_project_template/internal/featureconfig"
	"github.com/knoguchi/go_project_template/internal/logutil"
	"github.com/knoguchi/go_project_template/myapp"
	"github.com/knoguchi/go_project_template/version"
	"github.com/sirupsen/logrus"
	"github.com/urfave/cli/v2"
	prefixed "github.com/x-cray/logrus-prefixed-formatter"
	"golang.org/x/sync/errgroup"
	"os"
	"os/signal"
	"runtime"
	runtimeDebug "runtime/debug"
	"syscall"
)

var log = logrus.WithField("prefix", "mycmd")

// command line options are added in the order
var appFlags = []cli.Flag{
	flags.MyGlobalFlag,
	cmd.VerbosityFlag,
	cmd.LogFormat,
	cmd.LogFileName,
}

func init() {
	// wrap flags so that they can be loaded from alternative sources
	appFlags = cmd.WrapFlags(append(appFlags, featureconfig.MyCmdFlags...))
}

func main() {
	app := cli.App{}
	app.Name = "mycmd"
	app.Usage = "this is a my command"
	app.Action = startMyCmd
	app.Version = version.Version
	app.Commands = []*cli.Command{
		// add commands
	}

	app.Flags = appFlags

	app.Before = func(ctx *cli.Context) error {
		// Load flags from config file, if specified.
		if err := cmd.LoadFlagsFromConfig(ctx, app.Flags); err != nil {
			return err
		}

		format := ctx.String(cmd.LogFormat.Name)
		switch format {
		case "text":
			formatter := new(prefixed.TextFormatter)
			formatter.TimestampFormat = "2006-01-02 15:04:05"
			formatter.FullTimestamp = true
			// If persistent log files are written - we disable the log messages coloring because
			// the colors are ANSI codes and seen as gibberish in the log files.
			formatter.DisableColors = ctx.String(cmd.LogFileName.Name) != ""
			logrus.SetFormatter(formatter)
		case "json":
			logrus.SetFormatter(&logrus.JSONFormatter{})
		default:
			return fmt.Errorf("unknown log format \"%s\"", format)
		}

		logFileName := ctx.String(cmd.LogFileName.Name)
		if logFileName != "" {
			if err := logutil.ConfigurePersistentLogging(logFileName); err != nil {
				log.WithError(err).Error("Failed to configuring logging to disk.")
			}
		}
		if ctx.IsSet(flags.SetGCPercent.Name) {
			runtimeDebug.SetGCPercent(ctx.Int(flags.SetGCPercent.Name))
		}
		runtime.GOMAXPROCS(runtime.NumCPU())
		return debug.Setup(ctx)
	}

	defer func() {
		if x := recover(); x != nil {
			log.Errorf("Runtime panic: %v\n%v", x, string(runtimeDebug.Stack()))
			panic(x)
		}
	}()

	if err := app.Run(os.Args); err != nil {
		log.Error(err.Error())
	}
}

func startMyCmd(cliCtx *cli.Context) error {
	// setup logging
	verbosity := cliCtx.String(cmd.VerbosityFlag.Name)
	level, err := logrus.ParseLevel(verbosity)
	if err != nil {
		return err
	}
	logrus.SetLevel(level)

	// bootstrap configs
	configPath := cliCtx.String(cmd.ConfigFileFlag.FilePath)


	app, err := myapp.New(configPath)
	if err != nil {
		return err
	}

	// do not pass around cliCtx after this point
	ctx, cancel := context.WithCancel(cliCtx.Context)
	// setup errgroup
	g, gctx := errgroup.WithContext(ctx)

	// signal handling
	g.Go(func() error {
		// setup signal
		sigCh := make(chan os.Signal, 1)
		sigs := []os.Signal{
			os.Interrupt,
			os.Kill,
			syscall.SIGTERM,
			syscall.SIGABRT,
		}
		signal.Notify(sigCh, sigs...)

		// blocks until context is done or signal arrives
		select {
		case <-gctx.Done():
			log.Info("sig: gctx done")
			return gctx.Err()
		case <-sigCh:
			cancel()
			log.Info("sig: signal arrived")
			return gctx.Err()
		}
		return nil
	})

	// start myapp
	app.Start(gctx)

	// wait for all errgroup goroutines
	err = g.Wait()
	if err != nil {
		if errors.Is(err, context.Canceled) {
			log.Info("context was canceled")
		} else {
			log.Error("received error: %v", err)
		}
	} else {
		log.Infoln("startup: finished clean")
	}

	return err
}
