package nodelb

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"time"

	utils "github.com/aeternity/tool-nodelb/utils"

	yaml "gopkg.in/yaml.v2"
)

const (
	// ConfigFilename default configuration file name
	ConfigFilename = "config"
)

// TuningConfig fine tuning of parameters of the client
type TuningConfig struct {
	MovingAverageWindowSize int `json:"moving_average_window_size" yaml:"moving_average_window_size" mapstructure:"moving_average_window_size"`
	NodeBalancerSize        int `json:"node_balancer_size" yaml:"node_balancer_size" mapstructure:"node_balancer_size"`
}

// ConfigSchema define the configuration object
type ConfigSchema struct {
	Backends                []string     `json:"backends_url" yaml:"backends_url" mapstructure:"backends_url"`
	PeersURL                string       `json:"peer_url" yaml:"peer_url" mapstructure:"peer_url"`
	BackendPollIntervanlSec uint         `json:"backend_poll_interval_seconds" yaml:"backend_poll_interval_seconds" mapstructure:"backend_poll_interval_seconds"`
	ClientRetentionSec      uint         `json:"client_retention_seconds" yaml:"client_retention_seconds" mapstructure:"client_retention_seconds"`
	Tuning                  TuningConfig `json:"tuning" yaml:"tuning" mapstructure:"tuning"`
	ConfigPath              string       `json:"-" yaml:"-" mapstructure:"-"`
}

//Defaults generate configuration defaults
func (c *ConfigSchema) Defaults() *ConfigSchema {
	// for server
	c.Backends = []string{}
	utils.DefaultIfEmptyStr(&c.PeersURL, "https://sdk-mainnet.aepps.com/v2/debug/peers")
	utils.DefaultIfEmptyUint(&c.BackendPollIntervanlSec, 5)
	utils.DefaultIfEmptyUint(&c.ClientRetentionSec, 30)
	// for tuning
	utils.DefaultIfEmptyInt(&c.Tuning.MovingAverageWindowSize, 10)
	utils.DefaultIfEmptyInt(&c.Tuning.NodeBalancerSize, 20)
	return c
}

//Validate configuration
func (c *ConfigSchema) Validate() {
	valid := true

	if !valid {
		fmt.Println("Invalid configuration")
		os.Exit(1)
	}
}

// Config sytem configuration
var Config ConfigSchema

// GenerateDefaultConfig generate a default configuration
func GenerateDefaultConfig(outFile, version string) {
	Config = ConfigSchema{}
	Config.Defaults()
}

// Save save the configuration to disk
func (c *ConfigSchema) Save() {
	b, _ := yaml.Marshal(c)
	data := strings.Join([]string{
		"#",
		"#\n# Configuration for aepp-sdk-go \n#\n",
		fmt.Sprintf("#\n# update on %s \n#\n", time.Now().Format(time.RFC3339)),
		fmt.Sprintf("%s", b),
		"#\n# Config end\n#",
	}, "\n")
	if err := os.MkdirAll(filepath.Dir(c.ConfigPath), os.ModePerm); err != nil {
		fmt.Println("Cannot create config file path", err)
		os.Exit(1)
	}
	err := ioutil.WriteFile(c.ConfigPath, []byte(data), 0600)
	if err != nil {
		fmt.Println("Cannot create config file ", err)
		os.Exit(1)
	}
	fmt.Println("config file written to", c.ConfigPath)
}
