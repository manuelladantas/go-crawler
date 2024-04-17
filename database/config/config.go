package config

import (
	"fmt"

	configuration "github.com/golobby/config/v3"
	"github.com/golobby/config/v3/pkg/feeder"
)

type Database struct {
	Name string
	Uri  string
	Host string
	Port string
	User string
	Pass string
}

type Config struct {
	Database Database
}

func NewConfig() Config {
	cfg := Config{}
	jsonFeed := feeder.Json{Path: "database/config/server.json"}
	c := configuration.New()

	c.AddFeeder(jsonFeed).AddStruct(&cfg)
	if err := c.Feed(); err != nil {
		fmt.Println("[Error] fail to load config file", err.Error())
	}

	return cfg

}
