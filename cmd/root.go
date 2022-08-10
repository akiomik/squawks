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
	lang      string
	filters   []string
	includes  []string
	excludes  []string
	geocode   string
	url       string
	top       bool
	userAgent string
)

var rootCmd = &cobra.Command{
	Use:     "get-old-tweets --out FILENAME",
	Short:   "get-old-tweets v" + config.Version,
	Version: config.Version,
	Run: func(cmd *cobra.Command, args []string) {
		q := twitter.Query{
			Text:     text,
			Since:    since,
			Until:    until,
			From:     from,
			To:       to,
			Lang:     lang,
			Filters:  filters,
			Includes: includes,
			Excludes: excludes,
			Geocode:  geocode,
			Url:      url,
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

		ch := make(chan []export.Record)
		go func() {
			defer close(ch)

			opts := twitter.SearchOptions{Query: q, Top: top}
			for res := range client.SearchAll(opts) {
				if res.Error != nil {
					fmt.Fprintln(os.Stderr, "Error: %w", res.Error)
					os.Exit(1)
				}

				ch <- export.NewRecordsFromAdaptive(res.Adaptive)
			}
		}()

		<-export.ExportCsv(f, ch)
	},
}

func init() {
	rootCmd.Flags().StringSliceVarP(&excludes, "exclude", "", []string{}, "exclude tweets by type of tweet (e.g. hashtags, retweets, replies)")
	rootCmd.Flags().StringSliceVarP(&filters, "filter", "", []string{}, "find tweets by type of account or tweet (e.g. verified, follows, images, links)")
	rootCmd.Flags().StringVarP(&from, "from", "", "", "find tweets sent from a certain user")
	rootCmd.Flags().StringVarP(&geocode, "geocode", "", "", "find tweets sent from certain coordinates (e.g. 35.6851508,139.7526768,0.1km)")
	rootCmd.Flags().StringSliceVarP(&includes, "include", "", []string{}, "include tweets by type of tweet (e.g. hashtags, retweets, replies)")
	rootCmd.Flags().StringVarP(&lang, "lang", "", "", "find tweets by a certain language (e.g. en, es, fr)")
	rootCmd.Flags().StringVarP(&out, "out", "o", "", "output csv filename (required)")
	rootCmd.Flags().StringVarP(&text, "query", "q", "", "query text to search")
	rootCmd.Flags().StringVarP(&since, "since", "", "", "find tweets since a certain day (e.g. 2014-07-21)")
	rootCmd.Flags().StringVarP(&to, "to", "", "", "find tweets sent in reply to a certain user")
	rootCmd.Flags().BoolVarP(&top, "top", "", false, "find top tweets")
	rootCmd.Flags().StringVarP(&until, "until", "", "", "find tweets until a certain day (e.g. 2020-09-06)")
	rootCmd.Flags().StringVarP(&url, "url", "", "", "find tweets containing a certain url (e.g. www.example.com)")
	rootCmd.Flags().StringVarP(&userAgent, "user-agent", "", "", "set custom user-agent")
	rootCmd.MarkFlagRequired("out")
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err.Error())
		os.Exit(1)
	}
}
