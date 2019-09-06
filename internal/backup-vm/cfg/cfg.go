package cfg

import (
	"flag"
	"fmt"
	"os"

	"go.uber.org/zap"
)

var (
	// Cfg stores loggers and port data
	Cfg config

	// Execution env type: prod or dev
	envType = "ENVIRONMENT"
	envProd bool
)

type config struct {
	Logger     *zap.SugaredLogger
	FastLogger *zap.Logger
	Port       int
}

func init() {
	// Setup logger(s)
	env := os.Getenv(envType)
	if env == "prod" || env == "production" {
		envProd = true
	}
	var loggerFunc func(opts ...zap.Option) (*zap.Logger, error)
	if envProd {
		loggerFunc = zap.NewProduction
	} else {
		loggerFunc = zap.NewDevelopment
	}
	logger, err := loggerFunc()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
	defer logger.Sync()
	sugar := logger.Sugar()

	Cfg.Logger = sugar
	Cfg.FastLogger = logger

	// Args from cmd
	port := flag.Int("port", 8488, "Port to listen")
	help := flag.Bool("help", false, "Show usage")
	flag.Parse()

	if *help {
		flag.PrintDefaults()
		os.Exit(2)
	}

	Cfg.Port = *port
}
