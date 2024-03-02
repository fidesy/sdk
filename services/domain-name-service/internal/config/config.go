package config

import (
	"fmt"
	"gopkg.in/yaml.v3"
	"os"
	"strings"
)

const (
	RedisHost     = "redis-host"
	RedisPassword = "redis-password"
)

var (
	ErrConfigNotFoundByKey = func(key string) error {
		return fmt.Errorf("config not found by key = %q", key)
	}
)

var conf *Config

type Config struct {
	RedisHost     string `yaml:"redis-host"`
	RedisPassword string `yaml:"redis-password"`
}

func Init() error {
	ENV := os.Getenv("ENV")

	body, err := os.ReadFile(fmt.Sprintf("./configs/values_%s.yaml", strings.ToLower(ENV)))
	if err != nil {
		return err
	}

	err = yaml.Unmarshal(body, &conf)
	return err
}

func Get(key string) interface{} {
	switch key {
	case RedisHost:
		return conf.RedisHost
	case RedisPassword:
		return conf.RedisPassword
	default:
		panic(ErrConfigNotFoundByKey(key))
	}
}
