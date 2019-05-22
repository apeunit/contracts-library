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
	"crypto/tls"
	"encoding/json"
	"fmt"
	lb "github.com/aeternity/tool-nodelb/nodelb"
	"github.com/allegro/bigcache"
	"net/http"
	"net/http/httptrace"
	"net/http/httputil"
	"net/url"
	"os"
	"regexp"
	"time"

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
		var nc = lb.SetupCache()

		peers := new(Peers)
		fmt.Println("Proxy starting")

		fmt.Println("Requesting peers")
		// create a peer request
		cli := http.Client{Timeout: 10 * time.Second}
		req, err := http.NewRequest("GET", lb.Config.PeersURL, nil)
		if err != nil {
			fmt.Println("Error create peers request", err)
			return
		}
		res, err := cli.Do(req)
		if err != nil {
			fmt.Println("Error request peers", err)
			return
		}
		err = json.NewDecoder(res.Body).Decode(peers)
		if err != nil {
			fmt.Println("Error decode peers", err)
			return
		}
		// retrieve peer list
		r := regexp.MustCompile(`\d+\.\d+.\d+.\d+(\:\d+?)`)
		for _, p := range peers.Peers {
			ip := r.FindString(p)
			fmt.Println("Peer ", p, "found (", ip, ")")
			height, ct, ttfb, err := lb.TimeGet(fmt.Sprint("http://", ip, "/v2/key-blocks/current/height"))
			if err != nil {
				fmt.Println("Error ", err)
			}
			// add cached height
			lb.RegisterNode(ip, height, ttfb)
			fmt.Printf("ip: %10s H:%10d ct:%20s, ttfb: %10s", ip, height, ct, ttfb)
		}

		// start server
		http.HandleFunc("/", lb.HandleRequestAndRedirect)
		if err := http.ListenAndServe(getListenAddress(), nil); err != nil {
			panic(err)
		}
	},
}

// Get env var or default
func getEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}

// Get the port to listen on
func getListenAddress() string {
	port := getEnv("PORT", "1901")
	return ":" + port
}

// Peers the struct for peers
type Peers struct {
	Blocked  []string `json:"blocked"`
	Inbound  []string `json:"inbound"`
	Outbound []string `json:"outbound"`
	Peers    []string `json:"peers"`
}
