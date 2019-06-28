// Copyright © 2018 NAME HERE <EMAIL ADDRESS>
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package cmd

import (
	"fmt"

	lb "github.com/aeternity/aepp-contracts-library/aecl"
	"github.com/skratchdot/open-golang/open"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// configCmd represents the config command
var configCmd = &cobra.Command{
	Use:   "config",
	Short: "Print the configuration of the client",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Sprintf("%#v", lb.Config)
		viper.SetConfigType("yaml")
		err := viper.WriteConfigAs("examples/contracts_library.yaml")
		fmt.Println(err)
		// aeternity.Pp(
		// 	"Epoch URL", aeternity.Config.P.Epoch.URL,
		// 	"Epoch Internal URL", aeternity.Config.P.Epoch.InternalURL,
		// 	"Epoch Websocket URL", aeternity.Config.P.Epoch.WebsocketURL,
		// )
	},
}

var editCmd = &cobra.Command{
	Use:   "edit",
	Short: "Open the config file for editing",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Sprintf("%#v", lb.Config.ConfigPath)
		open.Run(lb.Config.ConfigPath)
	},
}

func init() {
	RootCmd.AddCommand(configCmd)
	configCmd.AddCommand(editCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// configCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// configCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")

}
