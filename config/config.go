package config

import (
	"log"
	"os"

	"github.com/caarlos0/env/v10"
	"gopkg.in/yaml.v2"
)

type Config struct {
	Addr   string `yaml:"addr" env:"ADDR"`
	Secret string `yaml:"secret" env:"SECRET"`
	DB     string `yaml:"db" env:"DB"`
}

func New() *Config {
	cfg := Config{}
	cfg.readFile()
	cfg.readEnv()
	return &cfg
}

func (cfg *Config) readFile() {
	f, err := os.Open("config/config.yml")
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	decoder := yaml.NewDecoder(f)
	err = decoder.Decode(cfg)
	if err != nil {
		log.Fatal(err)
	}
}

func (cfg *Config) readEnv() {
	if err := env.Parse(cfg); err != nil {
		log.Printf("%+v\n", err)
	}
}
