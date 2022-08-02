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
	"github.com/akiomik/get-old-tweets/export"
	"github.com/akiomik/get-old-tweets/twitter"
	"github.com/spf13/cobra"
)

var (
	out       string
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
		q := twitter.Query{
			Text:  text,
			Since: since,
			Until: until,
			From:  from,
			To:    to,
		}
		if q.IsEmpty() {
			fmt.Fprintln(os.Stderr, "Error: One or more queries are required")
			os.Exit(1)
		}

		f, err := os.OpenFile(out, os.O_WRONLY|os.O_CREATE|os.O_EXCL, 0644)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(1)
		}
		defer f.Close()

		client := twitter.NewClient()
		if len(userAgent) > 0 {
			client.UserAgent = userAgent
		}

		input := client.SearchAll(q)
		done := export.ExportCsv(f, input)

		<-done
	},
}

func init() {
	rootCmd.Flags().StringVarP(&out, "out", "o", "", "output csv filename")
	rootCmd.Flags().StringVarP(&text, "query", "q", "", "query text to be matched")
	rootCmd.Flags().StringVarP(&since, "since", "", "", "lower bound date to restrict search")
	rootCmd.Flags().StringVarP(&until, "until", "", "", "upper bound date to restrict search")
	rootCmd.Flags().StringVarP(&from, "from", "", "", "username from a twitter account")
	rootCmd.Flags().StringVarP(&to, "to", "", "", "username to a twitter account")
	rootCmd.Flags().StringVarP(&userAgent, "user-agent", "", "", "user-agent for request")
	rootCmd.MarkFlagRequired("out")
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err.Error())
		os.Exit(1)
	}
}
