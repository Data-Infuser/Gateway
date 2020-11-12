package config

import (
	"io/ioutil"
	"log"
	"os"
	"strconv"

	"gitlab.com/promptech1/infuser-gateway/constant"
	"gopkg.in/yaml.v2"
)

type Config struct {
	Author Author `yaml:"author"`
	Server Server `yaml:"server"`
}

type Author struct {
	Host   string `yaml:"host"`
	Port   int    `yaml:"port"`
	Tls    bool   `yaml:"tls"`
	CaFile string `yaml:"caFile"`
}

type Server struct {
	Host string `yaml:"host"`
	Port string `yaml:"Port"`
}

func (ctx *Config) getConfEnv() {
	var authorConfig *Author
	var serverConfig *Server

	authorConfig = new(Author)
	serverConfig = new(Server)

	authorConfig.Host = os.Getenv("GATEWAY_AUTHOR_CONFIG_HOST")
	authorConfig.Port, _ = strconv.Atoi(os.Getenv("GATEWAY_AUTHOR_CONFIG_PORT"))
	authorConfig.Tls, _ = strconv.ParseBool(os.Getenv("GATEWAY_AUTHOR_CONFIG_TLS"))
	authorConfig.CaFile = os.Getenv("GATEWAY_AUTHOR_CONFIG_CA_FILE")

	serverConfig.Host = os.Getenv("GATEWAY_SERVER_CONFIG_HOST")
	serverConfig.Port = os.Getenv("GATEWAY_SERVER_CONFIG_PORT")

	ctx.Author = *authorConfig
	ctx.Server = *serverConfig
}

func (ctx *Config) InitConf() error {
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

	return nil
}
