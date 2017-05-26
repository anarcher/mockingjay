package log

import (
	"github.com/go-kit/kit/log"

	stdlog "log"
	"os"
)

var Logger log.Logger

func init() {
	logger := log.NewLogfmtLogger(log.NewSyncWriter(os.Stderr))
	logger = log.With(logger, "ts", log.DefaultTimestampUTC)
	logger = log.With(logger, "caller", log.DefaultCaller)

	Logger = logger

	stdlog.SetFlags(0) // flags are handled by Go kit's logger
	stdlog.SetOutput(log.NewStdlibAdapter(Logger))

}
