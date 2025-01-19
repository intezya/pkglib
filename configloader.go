package pkglib

import (
	"flag"
	"github.com/ilyakaznacheev/cleanenv"
	"os"
)

type configLoader struct{}

//type TConfig = st{}

func (c configLoader) fetchConfigPath() string {
	var result string

	flag.StringVar(&result, "config", "", "path to config file")
	flag.Parse()

	if result == "" {
		result = os.Getenv("CONFIG_PATH")
	}

	return result
}

func (c configLoader) loadConfigFile(cfg interface{}, extraArgs ...string) {
	_ = cleanenv.ReadEnv(cfg)
	if extraArgs != nil {
		if err := cleanenv.ReadConfig(extraArgs[0], cfg); err != nil {
			panic("failed to read config: " + err.Error())
		}
	}

	path := c.fetchConfigPath()

	if path == "" {
		panic("config path is empty")
	}
	if _, err := os.Stat(path); os.IsNotExist(err) {
		panic("config file does not exist: " + path)
	}

	if err := cleanenv.ReadConfig(path, cfg); err != nil {
		panic("failed to read config: " + err.Error())
	}
}

func (c configLoader) New(cfg interface{}) {
	c.loadConfigFile(cfg)
}

func (c configLoader) NewWithExtraPath(cfg interface{}, path string) {
	c.loadConfigFile(cfg, path)
}

var ConfigLoader configLoader
