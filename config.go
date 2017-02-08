package mservice

import (
	"github.com/ilisin/configuration"
)

type Config struct {
	Host string `conf:"mservice.host,default(:8080)"`
}

func LoadAConfig() (*Config, error) {
	c := &Config{}
	err := configuration.Var(c)
	return c, err
}
