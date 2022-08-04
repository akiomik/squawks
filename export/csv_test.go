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

//go:build medium
// +build medium

package export

import (
	"encoding/csv"
	"io"
	"os"
	"reflect"
	"testing"
	"time"

	"github.com/akiomik/get-old-tweets/twitter"
)

func TestExportCsvEmpty(t *testing.T) {
	f, err := os.CreateTemp(os.TempDir(), "get-old-tweets-test-export-csv-empty-")
	if err != nil {
		t.Errorf("Failed to create tempfile: %v", err)
		return
	}
	defer f.Close()

	ch := make(chan *twitter.AdaptiveJson)
	close(ch)

	done := ExportCsv(f, ch)
	<-done
	f.Seek(0, 0)

	reader := csv.NewReader(f)
	_, err = reader.Read()
	if err != io.EOF {
		t.Errorf("Expect ExportCsv() to create an empty csv file, but not empty: %v", err)
		return
	}
}

func TestExportCsvNonEmpty(t *testing.T) {
	f, err := os.CreateTemp(os.TempDir(), "get-old-tweets-test-export-csv-non-empty-")
	if err != nil {
		t.Errorf("Failed to create tempfile: %v", err)
		return
	}
	defer f.Close()

	ch := make(chan *twitter.AdaptiveJson)
	go func() {
		defer close(ch)

		ch <- &twitter.AdaptiveJson{
			GlobalObjects: twitter.GlobalObjects{
				Tweets: map[string]twitter.Tweet{
					"1000": twitter.Tweet{
						Id:            1000,
						UserId:        2000,
						CreatedAt:     twitter.RubyDate(time.Date(2020, 9, 6, 0, 1, 2, 0, time.UTC)),
						FullText:      "To Sherlock Holmes she is always the woman.",
						RetweetCount:  3000,
						FavoriteCount: 4000,
						ReplyCount:    5000,
						QuoteCount:    6000,
						Lang:          "en",
					},
					"100": twitter.Tweet{
						Id:            100,
						UserId:        200,
						CreatedAt:     twitter.RubyDate(time.Date(2020, 9, 6, 0, 1, 2, 0, time.UTC)),
						FullText:      "To Sherlock Holmes she is always the woman.",
						RetweetCount:  300,
						FavoriteCount: 400,
						ReplyCount:    500,
						QuoteCount:    600,
						Lang:          "en",
					},
				},
				Users: map[string]twitter.User{
					"2000": twitter.User{
						Id:         2000,
						Name:       "Watson",
						ScreenName: "watson",
					},
					"200": twitter.User{
						Id:         200,
						Name:       "Watson",
						ScreenName: "watson",
					},
				},
			},
		}

		ch <- &twitter.AdaptiveJson{
			GlobalObjects: twitter.GlobalObjects{
				Tweets: map[string]twitter.Tweet{
					"10": twitter.Tweet{
						Id:            10,
						UserId:        20,
						CreatedAt:     twitter.RubyDate(time.Date(2020, 9, 6, 0, 1, 2, 0, time.UTC)),
						FullText:      "To Sherlock Holmes she is always the woman.",
						RetweetCount:  30,
						FavoriteCount: 40,
						ReplyCount:    50,
						QuoteCount:    60,
						Lang:          "en",
					},
					"1": twitter.Tweet{
						Id:            1,
						UserId:        2,
						CreatedAt:     twitter.RubyDate(time.Date(2020, 9, 6, 0, 1, 2, 0, time.UTC)),
						FullText:      "To Sherlock Holmes she is always the woman.",
						RetweetCount:  3,
						FavoriteCount: 4,
						ReplyCount:    5,
						QuoteCount:    6,
						Lang:          "en",
					},
				},
				Users: map[string]twitter.User{
					"2": twitter.User{
						Id:         2,
						Name:       "Watson",
						ScreenName: "watson",
					},
					"20": twitter.User{
						Id:         20,
						Name:       "Watson",
						ScreenName: "watson",
					},
				},
			},
		}
	}()

	done := ExportCsv(f, ch)
	<-done
	f.Seek(0, 0)

	reader := csv.NewReader(f)
	actualHeader, err := reader.Read()
	if err != nil {
		t.Errorf("Failed to read csv: %v", err)
		return
	}

	expectedHeader := []string{"id", "username", "created_at", "full_text", "retweet_count", "favorite_count", "reply_count", "quote_count", "geo", "coordinates", "lang", "source"}
	if !reflect.DeepEqual(actualHeader, expectedHeader) {
		t.Errorf("Expect ExportCsv() to write %v as a header, but got %v", expectedHeader, expectedHeader)
		return
	}

	expectedRecords := [][]string{
		[]string{"1000", "watson", "2020-09-06T00:01:02+00:00", "To Sherlock Holmes she is always the woman.", "3000", "4000", "5000", "6000", "", "", "en", ""},
		[]string{"100", "watson", "2020-09-06T00:01:02+00:00", "To Sherlock Holmes she is always the woman.", "300", "400", "500", "600", "", "", "en", ""},
		[]string{"10", "watson", "2020-09-06T00:01:02+00:00", "To Sherlock Holmes she is always the woman.", "30", "40", "50", "60", "", "", "en", ""},
		[]string{"1", "watson", "2020-09-06T00:01:02+00:00", "To Sherlock Holmes she is always the woman.", "3", "4", "5", "6", "", "", "en", ""},
	}

	for i, expectedRecord := range expectedRecords {
		actualRecord, err := reader.Read()
		if err != nil {
			t.Errorf("Failed to read csv: %v", err)
			return
		}

		if !reflect.DeepEqual(actualRecord, expectedRecord) {
			t.Errorf("Expect ExportCsv() to write %v as a record #%d, but got %v", expectedRecord, i, actualRecord)
			return
		}
	}

	_, err = reader.Read()
	if err != io.EOF {
		t.Errorf("Expect ExportCsv() to reach EOF, but not EOF: %v", err)
		return
	}
}
