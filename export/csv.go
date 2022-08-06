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
)

func ExportCsv(f *os.File, ch <-chan []Record) <-chan struct{} {
	done := make(chan struct{})

	go func() {
		defer close(done)

		w := csv.NewWriter(f)
		err := w.Write([]string{"id", "username", "created_at", "full_text", "retweet_count", "favorite_count", "reply_count", "quote_count", "geo", "coordinates", "lang", "source"})
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			panic(err)
		}

		for records := range ch {
			for _, record := range records {
				row := []string{
					strconv.FormatUint(record.Id, 10),
					record.Username,
					record.CreatedAt.String(),
					record.FullText,
					strconv.FormatUint(record.RetweetCount, 10),
					strconv.FormatUint(record.FavoriteCount, 10),
					strconv.FormatUint(record.ReplyCount, 10),
					strconv.FormatUint(record.QuoteCount, 10),
					record.Geo,
					record.Coodinates,
					record.Lang,
					record.Source,
				}

				err = w.Write(row)
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
