package main

import (
	"flag"
	stdlog "log"
	"net/http"
	"os"
	"strings"

	"github.com/go-kit/log"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/prometheus/common/promlog"
    
)

var (
	ExporterName = "diy_exporter"

	endpoint      string
	listenAddress string
	withProfiling bool
	logLevel      string
)

func setFlags(addFlags ...func()) {
	flag.StringVar(&endpoint, "web.telemetry-path", "/metrics", "Path under which to expose metrics.")
	flag.StringVar(&listenAddress, "web.listen-address", ":9090", "Address on which to expose metrics and web interface.")
	flag.BoolVar(&withProfiling, "profiling", false, "Enable profiling via web interface host:port/debug/pprof/")
	flag.StringVar(&logLevel, "log.level", "info", "Only log messages with the given severity or above. Valid levels: [debug, info, warn, error, fatal]")

	for _, f := range addFlags {
		f()
	}

	flag.Parse()
}

func newLogger() log.Logger {
	switch strings.ToLower(logLevel) {
	case "debug":
		logLevel = "debug"
	case "info":
		logLevel = "info"
	case "warn", "warning":
		logLevel = "warn"
	case "err", "error":
		logLevel = "error"
	default:
		logLevel = "info"
	}

	pLoglevel := &promlog.AllowedLevel{}
	_ = pLoglevel.Set(logLevel)
	promLogCfg := &promlog.Config{
		Level: pLoglevel,
	}

	return promlog.New(promLogCfg)
}

func main() {
	setFlags(
	)
	logger := newLogger()

	reg := prometheus.NewRegistry()

	http.Handle(endpoint, promhttp.HandlerFor(reg, promhttp.HandlerOpts{
		ErrorLog:      stdlog.New(log.NewStdlibAdapter(logger), "", 0),
		ErrorHandling: promhttp.ContinueOnError,
	}))
	if err := http.ListenAndServe(listenAddress, nil); err != nil {
		_ = logger.Log("msg", "failed to start web server", "err", err)
		os.Exit(1)
	}
}
