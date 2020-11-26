package config

import (
	"fmt"
	"io/ioutil"
	"os"

	"github.com/sirupsen/logrus"
	"gitlab.com/promptech1/infuser-gateway/constant"
	"gopkg.in/yaml.v2"
)

type Context struct {
	Logger *logrus.Entry
	Author Author `yaml:"author"`
}

type Author struct {
	Host   string `yaml:"host"`
	Port   int    `yaml:"port"`
	Tls    bool   `yaml:"tls"`
	CaFile string `yaml:"caFile"`
}

func (ctx *Context) InitContext() error {
	logger := logrus.New()

	var file []byte
	var err error

	var fileName string
	env := os.Getenv("GATEWAY_ENV")

	if len(env) > 0 && env == constant.ServiceProd {
		logger.SetLevel(logrus.InfoLevel)
		fileName = "config/config-prod.yaml"
	} else if len(env) > 0 && env == constant.ServiceStage {
		logger.SetLevel(logrus.InfoLevel)
		fileName = "config/config-stage.yaml"
	} else {
		logger.SetLevel(logrus.DebugLevel)
		fileName = "config/config-dev.yaml"
	}

	logger.Out = os.Stdout

	ctx.Logger = logger.WithFields(logrus.Fields{
		"tag": "gateway",
		"id":  os.Getpid(),
	})

	if file, err = ioutil.ReadFile(fileName); err != nil {
		return err
	}
	if err = yaml.Unmarshal(file, ctx); err != nil {
		return err
	}

	ctx.Logger.Info(fmt.Sprintf("Init configuration for '%s' env successfully =============", env))

	return nil
}
