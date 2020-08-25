package config

import (
	"io/ioutil"
	"log"
	"os"

	"gitlab.com/promptech1/infuser-gateway/constant"
	"gopkg.in/yaml.v2"
)

type Config struct {
	Author Author `yaml:"author"`
}

type Author struct {
	Host   string `yaml:"host"`
	Port   int    `yaml:"port"`
	Tls    bool   `yaml:"tls"`
	CaFile string `yaml:"caFile"`
}

func (ctx *Config) InitConf() error {
	var file []byte
	var err error

	ballast := make([]byte, 10<<30)
	_ = ballast

	var fileName string
	env := os.Getenv("GATEWAY_ENV")
	log.Printf("Init config with '%s' environment", env)

	if len(env) > 0 && env == constant.ServiceProd {
		fileName = "config/config-prod.yaml"
	} else if len(env) > 0 && env == constant.ServiceStage {
		fileName = "config/config-stage.yaml"
	} else {
		fileName = "config/config-dev.yaml"
	}
	log.Printf("Load '%s' file", fileName)

	if file, err = ioutil.ReadFile(fileName); err != nil {
		return err
	}
	if err = yaml.Unmarshal(file, ctx); err != nil {
		return err
	}

	return nil
}
