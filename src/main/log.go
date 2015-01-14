package main
import (
	"os"
	"github.com/op/go-logging"
)
var log = logging.MustGetLogger("server")
var logFormat = logging.MustStringFormatter(
	"%{color}%{time:15:04:05.000000} %{shortfunc} ▶ %{message}%{color:reset}")
var logBackend = logging.NewLogBackend(os.Stderr, "", 0)
var logFormattedBackend = logging.NewBackendFormatter(logBackend, logFormat)
var logLeveledBackend = logging.AddModuleLevel(logFormattedBackend)
func init () {
	// logLeveledBackend.SetLevel(logging.WARNING, "")
	logging.SetBackend(logLeveledBackend)
}