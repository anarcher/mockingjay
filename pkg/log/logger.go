package log

import (
	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/log/level"

	stdlog "log"
	"os"
)

var Logger log.Logger

func init() {
	logger := log.NewLogfmtLogger(log.NewSyncWriter(os.Stderr))
	logger = level.NewFilter(logger, level.AllowInfo())
	logger = log.With(logger, "ts", log.DefaultTimestampUTC)
	logger = log.With(logger, "caller", log.DefaultCaller)

	Logger = logger

	stdlog.SetFlags(0) // flags are handled by Go kit's logger
	stdlog.SetOutput(log.NewStdlibAdapter(Logger))

}
