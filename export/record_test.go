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
	"reflect"
	"testing"
	"time"

	"github.com/akiomik/get-old-tweets/twitter/json"
)

func TestNewRecordsFromAdaptive(t *testing.T) {
	j := json.Adaptive{
		GlobalObjects: json.GlobalObjects{
			Tweets: map[string]json.Tweet{
				"1000": json.Tweet{
					Id:            1000,
					UserId:        2000,
					CreatedAt:     json.RubyDate(time.Date(2020, 9, 6, 0, 1, 2, 0, time.UTC)),
					FullText:      "To Sherlock Holmes she is always the woman.",
					RetweetCount:  3000,
					FavoriteCount: 4000,
					ReplyCount:    5000,
					QuoteCount:    6000,
					Lang:          "en",
				},
				"100": json.Tweet{
					Id:            100,
					UserId:        200,
					CreatedAt:     json.RubyDate(time.Date(2020, 9, 6, 0, 1, 2, 0, time.UTC)),
					FullText:      "To Sherlock Holmes she is always the woman.",
					RetweetCount:  300,
					FavoriteCount: 400,
					ReplyCount:    500,
					QuoteCount:    600,
					Lang:          "en",
				},
			},
			Users: map[string]json.User{
				"2000": json.User{
					Id:         2000,
					Name:       "Watson",
					ScreenName: "watson1",
				},
				"200": json.User{
					Id:         200,
					Name:       "Watson",
					ScreenName: "watson2",
				},
			},
		},
	}

	expected := []Record{
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
			Lang:          "en",
		},
	}

	actual := NewRecordsFromAdaptive(&j)
	if !reflect.DeepEqual(actual, expected) {
		t.Errorf("Expect %+v, got %+v", expected, actual)
		return
	}
}
