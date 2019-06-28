package aecl

import (
	"fmt"
	"os"

	"github.com/spf13/viper"
)

const (
	// ConfigFilename default configuration file name
	ConfigFilename = "contracts_library"
)

//TuningSchema define fine tuning options schema
type TuningSchema struct {
	RequestMaxBodySize int64  `json:"max_body_size" yaml:"max_body_size" mapstructure:"max_body_size"`
	DbMaxOpenConns     int    `json:"max_open_connections" yaml:"max_open_connections" mapstructure:"max_open_connections"`
	DbMaxIdleConns     int    `json:"max_idle_connections" yaml:"max_idle_connections" mapstructure:"max_idle_connections"`
	VersionHeaderName  string `json:"version_header_name" yaml:"version_header_name" mapstructure:"version_header_name"`
}

// ConfigSchema define the configuration object
type ConfigSchema struct {
	ConfigPath    string           `json:"-" yaml:"-" mapstructure:"-"`
	Compilers     []CompilerSchema `json:"compilers" yaml:"compilers" mapstructure:"compilers"`
	DatabaseURL   string           `json:"db_url" yaml:"db_url" mapstructure:"db_url"`
	ListenAddress string           `json:"aecl_address" yaml:"aecl_address" mapstructure:"aecl_address"`
	Tuning        TuningSchema     `json:"tuning" yaml:"tuning" mapstructure:"tuning"`
}

// CompilerSchema is a configuration for the list of compilers
type CompilerSchema struct {
	URL       string `json:"url" yaml:"url" mapstructure:"url"`
	Version   string `json:"version" yaml:"version" mapstructure:"version"`
	IsDefault bool   `json:"is_default" yaml:"is_default" mapstructure:"is_default"`
}

//Defaults generate configuration defaults
func Defaults() {

	viper.SetDefault("DatabaseURL", "postgres://aecl:aecl@localhost/contracts_library?sslmode=disable")
	viper.SetDefault("ListenAddress", ":1905")
	viper.SetDefault("Tuning", map[string]string{
		"RequestMaxBodySize": "2000000",
		"DbMaxOpenConns":     "5",
		"DbMaxIdleConns":     "1",
		"VersionHeaderName":  "Sophia-Compiler-Version",
	})
	viper.SetDefault("Compilers", []map[string]string{
		map[string]string{
			"URL":       "http://localhost:3080",
			"Version":   "*",
			"IsDefault": "true",
		},
	})
}

//Validate configuration
func (c *ConfigSchema) Validate() {
	valid := true

	if !valid {
		fmt.Println("Invalid configuration")
		os.Exit(1)
	}
}

// Config system configuration
var Config ConfigSchema

// GenerateDefaultConfig generate a default configuration
func GenerateDefaultConfig(outFile, version string) {
	viper.Unmarshal(&Config)
}

// // Save save the configuration to disk
// func (c *ConfigSchema) Save() {
// 	b, _ := yaml.Marshal(c)
// 	data := strings.Join([]string{
// 		"#",
// 		"#\n# Configuration for aepp-contracts-library \n#\n",
// 		fmt.Sprintf("#\n# update on %s \n#\n", time.Now().Format(time.RFC3339)),
// 		fmt.Sprintf("%s", b),
// 		"#\n# Config end\n#",
// 	}, "\n")
// 	if err := os.MkdirAll(filepath.Dir(c.ConfigPath), os.ModePerm); err != nil {
// 		fmt.Println("Cannot create config file path", err)
// 		os.Exit(1)
// 	}
// 	err := ioutil.WriteFile(c.ConfigPath, []byte(data), 0600)
// 	if err != nil {
// 		fmt.Println("Cannot create config file ", err)
// 		os.Exit(1)
// 	}
// 	fmt.Println("config file written to", c.ConfigPath)
// }
