package config

import (
	"errors"
	"fmt"
	"os"
	"time"

	"gopkg.in/yaml.v3"
)

type Query struct {
	Name       string        `yaml:"name"`
	Query      string        `yaml:"query"`
	Expression string        `yaml:"expression"`
	Range      time.Duration `yaml:"range"`
	Step       time.Duration `yaml:"step"`
	BoolValue  bool          `yaml:"bool"`
	Units      string        `yaml:"units"`
}

type Service struct {
	Name  string `yaml:"name"`
	Query Query  `yaml:"query"`
	Group string `yaml:"group"`
	// Extras []Query `yaml:"extras"`
}

type UI struct {
	PageTitle       string        `yaml:"title"`
	RefreshInterval time.Duration `yaml:"refresh"`
	Graphs          Graphs        `yaml:"graphs"`
}

type Metrics struct {
	Enabled bool `yaml:"enabled"`
	Port    int  `yaml:"port"`
}

type Graphs struct {
	Points int `yaml:"points"`
}

type Config struct {
	Port       int           `yaml:"port"`
	Metrics    Metrics       `yaml:"metrics"`
	Services   []Service     `yaml:"services"`
	Prometheus string        `yaml:"prometheus"`
	Refresh    time.Duration `yaml:"refresh"`

	UI UI `yaml:"ui"`
}

func Load(path string) (*Config, error) {
	if _, err := os.Stat(path); err != nil {
		return nil, fmt.Errorf("config file stat: %v", err)
	}

	fb, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("config file read: %v", err)
	}

	conf := &Config{}
	if err := yaml.Unmarshal(fb, conf); err != nil {
		return nil, fmt.Errorf("config file parse: %v", err)
	}

	setDefaults(conf)

	if err := conf.Validate(); err != nil {
		return nil, err
	}

	return conf, nil
}

func setDefaults(conf *Config) {
	if conf.Port == 0 {
		conf.Port = 3000
	}
	if conf.Metrics.Port == 0 {
		conf.Metrics.Port = 9743
	}

	for i, svc := range conf.Services {
		svc.Query.Name = "main"
		if svc.Query.Range == 0 {
			svc.Query.Range = time.Hour * 24
		}
		if svc.Query.Step == 0 {
			svc.Query.Step = time.Minute * 5
		}
		if svc.Group == "" {
			svc.Group = "default"
		}
		conf.Services[i] = svc
	}

	if conf.Refresh == 0 {
		conf.Refresh = time.Second * 30
	}

	if conf.UI.PageTitle == "" {
		conf.UI.PageTitle = "PromPage"
	}
	if conf.UI.RefreshInterval == 0 {
		conf.UI.RefreshInterval = time.Second * 30
	}
	if conf.UI.Graphs.Points <= 0 {
		conf.UI.Graphs.Points = 200
	}
}

func (c *Config) Validate() error {
	if c.Prometheus == "" {
		return errors.New("prometheus cannot be empty")
	}

	return nil
}
