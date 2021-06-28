package config

import (
	"os"
	"strconv"
)

type ESProxyConfig struct {
	RedisHost     string
	RedisPort     int
	RedisPassword string
	RedisTTL      int
	EsHost        string
	EsPort        int
}

func GetConfig() *ESProxyConfig {
	conf := &ESProxyConfig{}

	if len(os.Getenv("REDIS_HOST")) > 0 {
		conf.RedisHost = os.Getenv("REDIS_HOST")
	} else {
		conf.RedisHost = "localhost"
	}

	if len(os.Getenv("REDIS_PORT")) > 0 {
		i, err := strconv.Atoi(os.Getenv("REDIS_PORT"))
		if err == nil {
			conf.RedisPort = i
		} else {
			conf.RedisPort = 6379
		}
	} else {
		conf.RedisPort = 6379
	}

	if len(os.Getenv("REDIS_PASSWORD")) > 0 {
		conf.RedisPassword = os.Getenv("REDIS_PASSWORD")
	} else {
		conf.RedisPassword = ""
	}

	if len(os.Getenv("REDIS_TTL")) > 0 {
		i, err := strconv.Atoi(os.Getenv("REDIS_TTL"))
		if err == nil {
			conf.RedisTTL = i
		} else {
			conf.RedisTTL = 30
		}
	} else {
		conf.RedisTTL = 30
	}

	return (conf)
}
