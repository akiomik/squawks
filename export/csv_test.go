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

	"github.com/akiomik/get-old-tweets/twitter/json"
)

func TestExportCsvEmpty(t *testing.T) {
	f, err := os.CreateTemp(os.TempDir(), "get-old-tweets-test-export-csv-empty-")
	if err != nil {
		t.Errorf("Failed to create tempfile: %v", err)
		return
	}
	defer f.Close()

	ch := make(chan []Record)
	close(ch)

	done := ExportCsv(f, ch)
	<-done
	f.Seek(0, 0)

	reader := csv.NewReader(f)
	_, err = reader.Read()
	if err != io.EOF {
		t.Errorf("Expect to create an empty csv file, got %v", err)
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

	ch := make(chan []Record)
	go func() {
		defer close(ch)

		ch <- []Record{
			Record{
				Id:            1000,
				Username:      "watson1",
				CreatedAt:     Iso8601Date(time.Date(2020, 9, 6, 0, 1, 2, 0, time.UTC)),
				FullText:      "To Sherlock Holmes she is always the woman.",
				RetweetCount:  3000,
				FavoriteCount: 4000,
				ReplyCount:    5000,
				QuoteCount:    6000,
				Lang:          "en",
			},
			Record{
				Id:            100,
				Username:      "watson2",
				CreatedAt:     Iso8601Date(time.Date(2020, 9, 6, 0, 1, 2, 0, time.UTC)),
				FullText:      "To Sherlock Holmes she is always the woman.",
				RetweetCount:  300,
				FavoriteCount: 400,
				ReplyCount:    500,
				QuoteCount:    600,
				Coordinates:   &json.LongLat{-73.9998279, 40.74118764},
				Lang:          "en",
			},
		}

		ch <- []Record{
			Record{
				Id:            10,
				Username:      "watson3",
				CreatedAt:     Iso8601Date(time.Date(2020, 9, 6, 0, 1, 2, 0, time.UTC)),
				FullText:      "To Sherlock Holmes she is always the woman.",
				RetweetCount:  30,
				FavoriteCount: 40,
				ReplyCount:    50,
				QuoteCount:    60,
				Lang:          "en",
			},
			Record{
				Id:            1,
				Username:      "watson4",
				CreatedAt:     Iso8601Date(time.Date(2020, 9, 6, 0, 1, 2, 0, time.UTC)),
				FullText:      "To Sherlock Holmes she is always the woman.",
				RetweetCount:  3,
				FavoriteCount: 4,
				ReplyCount:    5,
				QuoteCount:    6,
				Lang:          "en",
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

	expectedHeader := []string{"id", "username", "created_at", "full_text", "retweet_count", "favorite_count", "reply_count", "quote_count", "coordinates", "lang", "source"}
	if !reflect.DeepEqual(actualHeader, expectedHeader) {
		t.Errorf("Expect to write %v as a header, got %v", expectedHeader, actualHeader)
		return
	}

	expectedRecords := [][]string{
		[]string{"1000", "watson1", "2020-09-06T00:01:02+00:00", "To Sherlock Holmes she is always the woman.", "3000", "4000", "5000", "6000", "", "en", ""},
		[]string{"100", "watson2", "2020-09-06T00:01:02+00:00", "To Sherlock Holmes she is always the woman.", "300", "400", "500", "600", "-73.9998279,40.74118764", "en", ""},
		[]string{"10", "watson3", "2020-09-06T00:01:02+00:00", "To Sherlock Holmes she is always the woman.", "30", "40", "50", "60", "", "en", ""},
		[]string{"1", "watson4", "2020-09-06T00:01:02+00:00", "To Sherlock Holmes she is always the woman.", "3", "4", "5", "6", "", "en", ""},
	}

	for i, expectedRecord := range expectedRecords {
		actualRecord, err := reader.Read()
		if err != nil {
			t.Errorf("Failed to read csv: %v", err)
			return
		}

		if !reflect.DeepEqual(actualRecord, expectedRecord) {
			t.Errorf("Expect to write %v as a record #%d, got %v", expectedRecord, i, actualRecord)
			return
		}
	}

	_, err = reader.Read()
	if err != io.EOF {
		t.Errorf("Expect to reach EOF, got not EOF %v", err)
		return
	}
}
