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

package export

import (
	"encoding/csv"
	"fmt"
	"os"
	"strconv"

	"github.com/akiomik/get-old-tweets/twitter"
)

func ExportCsv(f *os.File, ch <-chan *twitter.Adaptive) <-chan struct{} {
	done := make(chan struct{})

	go func() {
		defer close(done)

		w := csv.NewWriter(f)
		err := w.Write([]string{"id", "username", "created_at", "full_text", "retweet_count", "favorite_count", "reply_count", "quote_count", "geo", "coordinates", "lang", "source"})
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			panic(err)
		}

		for j := range ch {
			for _, k := range ReversedKeysOf(j.GlobalObjects.Tweets) {
				t := j.GlobalObjects.Tweets[k]
				u := j.GlobalObjects.Users[strconv.FormatUint(t.UserId, 10)]
				r := []string{
					strconv.FormatUint(t.Id, 10),
					u.ScreenName,
					t.CreatedAt.Iso8601(),
					t.FullText,
					strconv.FormatUint(t.RetweetCount, 10),
					strconv.FormatUint(t.FavoriteCount, 10),
					strconv.FormatUint(t.ReplyCount, 10),
					strconv.FormatUint(t.QuoteCount, 10),
					t.Geo,
					t.Coodinates,
					t.Lang,
					t.Source,
				}

				err = w.Write(r)
				if err != nil {
					fmt.Fprintf(os.Stderr, "Error: %v\n", err)
					panic(err)
				}
			}

			w.Flush()
		}

		done <- struct{}{}
	}()

	return done
}
