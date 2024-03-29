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
	ConfigPath    string             `json:"-" yaml:"-" mapstructure:"-"`
	Compilers     []CompilerSchema   `json:"compilers" yaml:"compilers" mapstructure:"compilers"`
	DatabaseURL   string             `json:"db_url" yaml:"db_url" mapstructure:"db_url"`
	ListenAddress string             `json:"aecl_address" yaml:"aecl_address" mapstructure:"aecl_address"`
	Tuning        TuningSchema       `json:"tuning" yaml:"tuning" mapstructure:"tuning"`
	Web           WebResourcesSchema `json:"web" yaml:"web" mapstructure:"web"`
}

//WebResourcesSchema define the configuration for web
type WebResourcesSchema struct {
	HomeTemplatePath string `json:"home_template_path" yaml:"home_template_path" mapstructure:"home_template_path"`
	TosTemplatePath  string `json:"tos_template_path" yaml:"tos_template_path" mapstructure:"tos_template_path"`
	AssetsFolderPath string `json:"assets_folder_path" yaml:"assets_folder_path" mapstructure:"assets_folder_path"`
	AssetsWebPath    string `json:"assets_web_path" yaml:"assets_web_path" mapstructure:"assets_web_path"`
}

// CompilerSchema is a configuration for the list of compilers
type CompilerSchema struct {
	URL       string `json:"-" yaml:"url" mapstructure:"url"`
	Version   string `json:"version" yaml:"version" mapstructure:"version"`
	IsDefault bool   `json:"is_default" yaml:"is_default" mapstructure:"is_default"`
}

//Defaults generate configuration defaults
func Defaults() {

	viper.SetDefault("db_url", "postgres://aecl:aecl@localhost/contracts_library?sslmode=disable")
	viper.SetDefault("aecl_address", ":1905")
	viper.SetDefault("tuning", map[string]interface{}{
		"max_body_size":        2000000,
		"max_open_connections": 5,
		"max_idle_connections": 1,
		"version_header_name":  "Sophia-Compiler-Version",
	})
	viper.SetDefault("web", map[string]interface{}{
		"home_template_path": "templates/home.html",
		"tos_template_path":  "templates/tos.html",
		"assets_folder_path": "templates/assets",
		"assets_web_path":    "/assets/*",
	})
	viper.SetDefault("compilers", []map[string]interface{}{
		map[string]interface{}{
			"url":        "http://localhost:3080",
			"version":    "*",
			"is_default": true,
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
