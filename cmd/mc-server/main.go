package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/candbright/go-log/log"
	"github.com/candbright/go-log/options"
	server "github.com/candbright/go-server/internal/mc-server"
	"github.com/candbright/go-server/pkg/config"
	"github.com/candbright/go-server/pkg/rest"
	"github.com/sirupsen/logrus"
)

var _BUILD_ = ""

var (
	help       bool
	version    bool
	configFile string
)

func init() {
	flag.BoolVar(&help, "h", false, "print help")
	flag.BoolVar(&version, "v", false, "print version")
	flag.StringVar(&configFile, "c", "conf/config.yaml", "configuration file path")
}

func main() {
	flag.Parse()
	if help {
		flag.Usage()
		return
	}
	if version {
		fmt.Println(_BUILD_)
		return
	}

	if err := loadConfig(configFile); err != nil {
		logAndExit("Load config failed! error: %v", err)
	}
	if err := initLogger(); err != nil {
		logAndExit("Logger init failed! error: %v", err)
	}

	if config.Global.GetBool("log.debug") {
		rest.SetErrorHandler(func(err error) {
			log.Errorf("%+v", err)
		})
	}

	s := server.NewServer()
	s.Serve()
}

func loadConfig(path string) error {
	return config.InitFromFile(path)
}

func initLogger() error {
	return log.Init(
		options.Path(config.Global.Get("log.path")),
		options.Level(func() logrus.Level {
			level, err := logrus.ParseLevel(config.Global.Get("log.level"))
			if err != nil {
				level = logrus.InfoLevel
			}
			return level
		}),
		options.Format(&logrus.JSONFormatter{}),
		options.GlobalField("app_name", config.Global.Get("server.name")),
	)
}

func logAndExit(format string, args ...interface{}) {
	fmt.Printf(format+"\n", args...)
	os.Exit(1)
}
