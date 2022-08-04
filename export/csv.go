// copyright 2022 akiomi kamakura
//
// licensed under the apache license, version 2.0 (the "license");
// you may not use this file except in compliance with the license.
// you may obtain a copy of the license at
//
//     http://www.apache.org/licenses/license-2.0
//
// unless required by applicable law or agreed to in writing, software
// distributed under the license is distributed on an "as is" basis,
// without warranties or conditions of any kind, either express or implied.
// see the license for the specific language governing permissions and
// limitations under the license.

package export

import (
	"encoding/csv"
	"fmt"
	"os"
	"strconv"

	"github.com/akiomik/get-old-tweets/twitter"
)

func ExportCsv(f *os.File, ch <-chan *twitter.AdaptiveJson) <-chan struct{} {
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
