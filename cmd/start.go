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
	"log"
	"net/http"

	"github.com/aeternity/aepp-contracts-library/aecl"
	"github.com/spf13/cobra"
)

func init() {
	RootCmd.AddCommand(startCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// chainCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// chainCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")

}

var startCmd = &cobra.Command{
	Use:   "start",
	Short: "Contract library builder",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("     ______    ___      _____      _   ")
		fmt.Println("   .' ___  | .'   `.   |_   _|    (_)  ")
		fmt.Println("  / .'   \\_|/  .-.  \\    | |      __   ")
		fmt.Println("  | |       | |   | |    | |   _ [  |  ")
		fmt.Println("  \\ `.___.'\\   `-'  /_  _| |__/ | | |  ")
		fmt.Println("   `.____ .' `.___.'(_)|________|[___] v", RootCmd.Version)
		fmt.Println()
		log.Println("Contract library for aeternity started", aecl.Config.ListenAddress)
		log.Println("Listening on address", aecl.Config.ListenAddress)
		log.Println("Available compilers: ", len(aecl.Config.Compilers))
		aecl.StartProxy()
		// start server
		http.HandleFunc("/", aecl.HandleRequestAndRedirect)
		if err := http.ListenAndServe(aecl.Config.ListenAddress, nil); err != nil {
			panic(err)
		}
	},
}
