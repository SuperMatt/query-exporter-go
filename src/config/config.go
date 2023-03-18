package config

import (
	"os"

	//yaml
	"gopkg.in/yaml.v3"

	// mergo
	"github.com/imdario/mergo"
)

type Config struct {
	Endpoints []Endpoint `yaml:"endpoints"`
	// The port to listen on
	Server  Server  `yaml:"server"`
	Metrics Metrics `yaml:"metrics"`
	Debug   bool    `yaml:"debug"`
}

type Server struct {

	// The port to listen on
	Port            int    `yaml:"port"`
	MetricsEndpoint string `yaml:"metrics_endpoint"`
}

type Metrics struct {
	Prefix string `yaml:"prefix"`
}

type Endpoint struct {
	Name         string   `yaml:"name"`
	Address      string   `yaml:"address"`
	Headers      []Header `yaml:"headers"`
	Query        string   `yaml:"query"`
	QueryOffsets []string `yaml:"query_offsets"`
}

type Header struct {
	Name  string `yaml:"name"`
	Value string `yaml:"value"`
}

// A function to create a new config object
func NewConfig() *Config {
	cfg := &Config{
		Server: Server{
			Port:            8080,
			MetricsEndpoint: "/metrics",
		},
		Metrics: Metrics{
			Prefix: "query_exporter",
		},
		Debug: false,
	}
	return cfg
}

// load the config from a file.
// the file should either be the one provided, or config.yaml or config.yml
func LoadConfig(file string) (*Config, error) {
	cfg := NewConfig()
	defaultCfg := NewConfig()

	// check if the config file is provided
	if file == "" {
		// check if config.yaml or config.yml exists
		// if it does, use that
		// if it doesn't, use the default config
		for _, f := range []string{"config.yaml", "config.yml"} {
			_, err := os.Stat(f)
			if err == nil {
				file = f
				break
			}
		}

		// if the file is still empty, use the default config
		if file == "" {
			return defaultCfg, nil
		}
	}

	// open the contents of the file
	f, err := os.Open(file)
	if err != nil {
		return defaultCfg, err
	}
	defer f.Close()

	// decode the yaml
	decoder := yaml.NewDecoder(f)
	err = decoder.Decode(cfg)

	// if there was an error, return the default config
	if err != nil {
		return defaultCfg, err
	}

	// use mergo to deep merge the default config with the loaded config
	// this will ensure that all fields are present
	err = mergo.Merge(cfg, defaultCfg)
	if err != nil {
		return defaultCfg, err
	}

	return cfg, nil
}
