package config

import (
	"io/ioutil"
	"os"

	"gopkg.in/yaml.v2"
)

type Config struct {
	Port int
	Host string

	// if true, the server will attempt to obtain a
	// TLS cert from Let's Encrypt. The server will run
	// without TLS otherwise.
	AcmeTLS bool

	// if true, clients will need to authenticate by going through
	// oauth with Google. If false, the server will allow anyone
	// impersonate any user.
	VerifyUser bool

	// if true, every request made to the server will be logged
	// to stdout
	RequestLogs bool

	// if true, will verify the origin of incoming websocket connections
	CheckOrigin bool
}

// defaultConfig defines a config suitable for local development
// Production instances can customize it by creating a config.yaml
// and setting the environment variable UTTT_CONFIG_FILE
var defaultConfig Config = Config{
	Port:        8080,
	Host:        "localhost",
	AcmeTLS:     false,
	VerifyUser:  false,
	RequestLogs: true,
	CheckOrigin: false,
}

// Load returns the configuration for the server to
// use. If the environment does not specify a config file (or if
// said config file does not specify all fields), default
// values are used.
func Load() (*Config, error) {
	cfg := defaultConfig

	filename := os.Getenv("UTTT_CONFIG_FILE")
	if filename == "" {
		return &cfg, nil
	}

	data, err := ioutil.ReadFile(filename)
	if err != nil {
		return &cfg, err
	}

	err = yaml.Unmarshal(data, &cfg)
	return &cfg, err
}
