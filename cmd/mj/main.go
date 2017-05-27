package main

import (
	"github.com/anarcher/mockingjay/pkg/config"
	"github.com/anarcher/mockingjay/pkg/handler"
	"github.com/anarcher/mockingjay/pkg/log"

	"github.com/elazarl/goproxy"

	"flag"
	"fmt"
	"net/http"
	"os"
)

var (
	configFile string
	verbose    bool
	addr       string
)

var Usage = func() {
	fmt.Fprintf(os.Stderr, "Usage of %s:\n", os.Args[0])
	flag.PrintDefaults()
}

func init() {
	flag.StringVar(&configFile, "config", "mj.yaml", "")
	flag.BoolVar(&verbose, "v", false, "verbose: print additional output")
	flag.StringVar(&addr, "addr", ":8080", "")
}

func main() {
	logger := log.Logger
	logger.Log("mj", "start")

	flag.Parse()
	if configFile == "" {
		Usage()
		os.Exit(1)
	}
	if addr == "" {
		Usage()
		os.Exit(1)
	}

	cfg, err := config.ReadConfigFile(configFile)
	if err != nil {
		logger.Log("err", err)
		os.Exit(1)
	}

	proxy := goproxy.NewProxyHttpServer()
	proxy.Verbose = verbose

	for i, proxyCfg := range cfg.Proxies {
		ph, err := handler.NewProxyHandler(proxyCfg)
		if err != nil {
			logger.Log("err", err)
			os.Exit(1)
		}

		proxy.OnRequest(ph).Do(ph)
		logger.Log("handler", "add", "idx", i)
	}

	{
		err := http.ListenAndServe(addr, proxy)
		logger.Log("err", err)
	}
	logger.Log("mj", "end")
}
