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

package search

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"github.com/akiomik/squawks/api"
	"github.com/akiomik/squawks/cmd/flags"
	"github.com/akiomik/squawks/export"
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
	near      string
	within    string
	top       bool
	userAgent string
)

func NewTweetsCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "tweets --out FILENAME",
		Short: "Search for tweets",
		Run: func(cmd *cobra.Command, args []string) {
			q := api.Query{
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
				Near:     near,
				Within:   within,
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

			client := api.NewClient()
			if len(userAgent) > 0 {
				client.UserAgent = userAgent
			}

			ch := make(chan []export.Record)
			go func() {
				defer close(ch)

				opts := api.SearchOptions{Query: q, Top: top}
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

	flags.StringSliceEnumVarP(cmd.Flags(), &excludes, "exclude", "", []string{}, "exclude tweets by type of tweet", []string{"hashtags", "nativeretweets", "retweets", "replies"})
	flags.StringSliceEnumVarP(cmd.Flags(), &filters, "filter", "", []string{}, "find tweets by type of account or tweet", []string{"verified", "follows", "media", "images", "twimg", "videos", "periscope", "vine", "consumer_video", "pro_video", "native_video", "links", "hashtags", "nativeretweets", "retweets", "replies", "safe", "news"})
	cmd.Flags().StringVarP(&from, "from", "", "", "find tweets sent from a certain user")
	cmd.Flags().StringVarP(&geocode, "geocode", "", "", "find tweets sent from certain coordinates (e.g. 35.6851508,139.7526768,0.1km)")
	flags.StringSliceEnumVarP(cmd.Flags(), &includes, "include", "", []string{}, "include tweets by type of tweet", []string{"hashtags", "nativeretweets", "retweets", "replies"})
	cmd.Flags().StringVarP(&lang, "lang", "", "", "find tweets by a certain language (e.g. en, es, fr)")
	cmd.Flags().StringVarP(&near, "near", "", "", "find tweets nearby a certain location (e.g. tokyo)")
	cmd.Flags().StringVarP(&out, "out", "o", "", "output csv filename (required)")
	cmd.Flags().StringVarP(&text, "query", "q", "", "query text to search")
	cmd.Flags().StringVarP(&since, "since", "", "", "find tweets since a certain day (e.g. 2014-07-21)")
	cmd.Flags().StringVarP(&to, "to", "", "", "find tweets sent in reply to a certain user")
	cmd.Flags().BoolVarP(&top, "top", "", false, "find top tweets")
	cmd.Flags().StringVarP(&until, "until", "", "", "find tweets until a certain day (e.g. 2020-09-06)")
	cmd.Flags().StringVarP(&url, "url", "", "", "find tweets containing a certain url (e.g. www.example.com)")
	cmd.Flags().StringVarP(&userAgent, "user-agent", "", "", "set custom user-agent")
	cmd.Flags().StringVarP(&within, "within", "", "", "find tweets nearby a certain location (e.g. 1km)")
	cmd.MarkFlagRequired("out")

	return cmd
}
