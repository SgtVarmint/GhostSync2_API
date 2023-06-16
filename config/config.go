package config

import (
	"fmt"
	"io/ioutil"
	"gopkg.in/yaml.v3"
)

type Config struct {
	AccessCode			string `yaml:"access_code"`
	VideoRoot			string `yaml:"video_root"`
	AuthenticatedCode	string `yaml:"authenticated_code"`
}

func ParseConfig() (*Config, error) {
	buf, err := ioutil.ReadFile("env.yml")
	if err != nil {
		return nil, err
	}

	config := &Config{}
	err = yaml.Unmarshal(buf, &config)
	if err != nil {
		return nil, fmt.Errorf("Error: ", err)
	}

	return config, nil
}