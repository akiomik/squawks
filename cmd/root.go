// Copyright 2022 Akiomi Kamakura
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
	"os"

	"github.com/akiomik/get-old-tweets/config"
	"github.com/akiomik/get-old-tweets/twitter"
	"github.com/spf13/cobra"
)

var (
	text      string
	since     string
	until     string
	from      string
	to        string
	userAgent string
)

var rootCmd = &cobra.Command{
	Use:     "get-old-tweets",
	Short:   "get-old-tweets v" + config.Version,
	Version: config.Version,
	Run: func(cmd *cobra.Command, args []string) {
		query := twitter.Query{
			Text:  text,
			Since: since,
			Until: until,
			From:  from,
			To:    to,
		}
		if query.IsEmpty() {
			fmt.Fprintln(os.Stderr, "Error: One or more queries are required")
			os.Exit(1)
		}

		client := twitter.NewClient()
		if len(userAgent) > 0 {
			client.UserAgent = userAgent
		}

		ch := client.SearchAll(query)

		for json := range ch {
			for _, tweet := range json.GlobalObjects.Tweets {
				fmt.Printf("%+v\n", tweet)
			}
		}
	},
}

func init() {
	rootCmd.Flags().StringVarP(&text, "query", "q", "", "query text to be matched")
	rootCmd.Flags().StringVarP(&since, "since", "", "", "lower bound date to restrict search")
	rootCmd.Flags().StringVarP(&until, "until", "", "", "upper bound date to restrict search")
	rootCmd.Flags().StringVarP(&from, "from", "", "", "username from a twitter account")
	rootCmd.Flags().StringVarP(&to, "to", "", "", "username to a twitter account")
	rootCmd.Flags().StringVarP(&userAgent, "user-agent", "", "", "user-agent for request")
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err.Error())
		os.Exit(1)
	}
}
