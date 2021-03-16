package logutil

import (
	"github.com/knoguchi/go_project_template/internal/config"
	"github.com/sirupsen/logrus"
	"io"
	"os"
)

func addLogWriter(w io.Writer) {
	mw := io.MultiWriter(logrus.StandardLogger().Out, w)
	logrus.SetOutput(mw)
}

// ConfigurePersistentLogging adds a log-to-file writer. File content is identical to stdout.
func ConfigurePersistentLogging(logFileName string) error {
	logrus.WithField("logFileName", logFileName).Info("Logs will be made persistent")
	f, err := os.OpenFile(logFileName, os.O_CREATE|os.O_WRONLY|os.O_APPEND, config.CurrentIoConfig().ReadWritePermissions)
	if err != nil {
		return err
	}

	addLogWriter(f)

	logrus.Info("File logging initialized")
	return nil
}
