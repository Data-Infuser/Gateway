// config: Logger, Gateway 서버, 인증서버에 대한 설정 정의
package config

import (
	"fmt"
	"io/ioutil"
	"os"
	"strconv"

	"github.com/sirupsen/logrus"
	"gitlab.com/promptech1/infuser-gateway/constant"
	"gopkg.in/yaml.v2"
)

type Config struct {
	Logger   *logrus.Entry
	Author   Author   `yaml:"author"`
	Server   Server   `yaml:"server"`
	Executor Executor `yaml:"executor"`
}

type Author struct {
	Host   string `yaml:"host"`
	Port   int    `yaml:"port"`
	Tls    bool   `yaml:"tls"`
	CaFile string `yaml:"caFile"`
}

type Executor struct {
	Host   string `yaml:"host"`
	Port   int    `yaml:"port"`
	Tls    bool   `yaml:"tls"`
	CaFile string `yaml:"caFile"`
}

type Server struct {
	Host string `yaml:"host"`
	Port string `yaml:"port"`
}

// 구동 환경(Dev, Stage, Prod)에 따른 설정 정보 정의
func (ctx *Config) getConfEnv() {
	var authorConfig *Author
	var executorConfig *Executor
	var serverConfig *Server

	authorConfig = new(Author)
	executorConfig = new(Executor)
	serverConfig = new(Server)

	authorConfig.Host = os.Getenv("GATEWAY_AUTHOR_CONFIG_HOST")
	authorConfig.Port, _ = strconv.Atoi(os.Getenv("GATEWAY_AUTHOR_CONFIG_PORT"))
	authorConfig.Tls, _ = strconv.ParseBool(os.Getenv("GATEWAY_AUTHOR_CONFIG_TLS"))
	authorConfig.CaFile = os.Getenv("GATEWAY_AUTHOR_CONFIG_CA_FILE")

	serverConfig.Host = os.Getenv("GATEWAY_EXECUTOR_CONFIG_HOST")
	serverConfig.Port = os.Getenv("GATEWAY_EXECUTOR_CONFIG_PORT")
	authorConfig.Tls, _ = strconv.ParseBool(os.Getenv("GATEWAY_EXECUTOR_CONFIG_TLS"))
	authorConfig.CaFile = os.Getenv("GATEWAY_EXECUTOR_CONFIG_CA_FILE")

	serverConfig.Host = os.Getenv("GATEWAY_SERVER_CONFIG_HOST")
	serverConfig.Port = os.Getenv("GATEWAY_SERVER_CONFIG_PORT")

	ctx.Author = *authorConfig
	ctx.Server = *serverConfig
	ctx.Executor = *executorConfig
}

// Configration 정의: Log Level, 환경에 따른 Conf 설정 파일 Load
func (ctx *Config) InitConf() error {
	var fileName string
	env := os.Getenv("GATEWAY_ENV")

	logger := logrus.New()

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

	if _, err := os.Stat(fileName); os.IsNotExist(err) {
		ctx.getConfEnv()
	} else {
		var file []byte
		var err error

		if file, err = ioutil.ReadFile(fileName); err != nil {
			return err
		}
		if err = yaml.Unmarshal(file, ctx); err != nil {
			return err
		}
	}

	ctx.Logger.Info(fmt.Sprintf("Init configuration for '%s' env successfully =============", env))

	return nil
}
