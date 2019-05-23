package aecl

import (
	"fmt"
	"os"

	utils "github.com/aeternity/aepp-contracts-library/utils"
)

const (
	// ConfigFilename default configuration file name
	ConfigFilename = "config"
)

type TuningSchema struct {
	MaxBodySize int64 `json:"max_body_size" yaml:"max_body_size" mapstructure:"max_body_size"`
}

// ConfigSchema define the configuration object
type ConfigSchema struct {
	ConfigPath    string       `json:"-" yaml:"-" mapstructure:"-"`
	CompilerURL   string       `json:"compiler_url" yaml:"compiler_url" mapstructure:"compiler_url"`
	DatabaseURL   string       `json:"db_url" yaml:"db_url" mapstructure:"db_url"`
	ListenAddress string       `json:"aecl_address" yaml:"aecl_address" mapstructure:"aecl_address"`
	Tuning        TuningSchema `json:"tuning" yaml:"tuning" mapstructure:"tuning"`
}

//Defaults generate configuration defaults
func (c *ConfigSchema) Defaults() *ConfigSchema {
	// for server
	c.CompilerURL = utils.GetEnv("COMPILER_URL", "https://compiler.aepps.com")
	c.DatabaseURL = utils.GetEnv("DATABASE_URL", "postgres://middleware:middleware@35.228.174.89:5432/contracts_library?sslmode=disable")
	c.ListenAddress = utils.GetEnv("AECL_ADDRESS", ":1905")
	c.Tuning = TuningSchema{
		MaxBodySize: utils.GetEnvInt64("MAX_BODY_SIZE", 2000000),
	}
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

// Config system configuration
var Config ConfigSchema

// GenerateDefaultConfig generate a default configuration
func GenerateDefaultConfig(outFile, version string) {
	Config = ConfigSchema{}
	Config.Defaults()
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
