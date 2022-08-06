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
	"strconv"
	"time"

	"github.com/akiomik/get-old-tweets/twitter/json"
)

type Record struct {
	Id            uint64
	Username      string
	CreatedAt     Iso8601Date
	FullText      string
	RetweetCount  uint64
	FavoriteCount uint64
	ReplyCount    uint64
	QuoteCount    uint64
	Geo           string
	Coodinates    string
	Lang          string
	Source        string
}

func NewRecordsFromAdaptive(j *json.Adaptive) []Record {
	l := make([]Record, len(j.GlobalObjects.Tweets))

	for i, k := range ReversedKeysOf(j.GlobalObjects.Tweets) {
		t := j.GlobalObjects.Tweets[k]
		u := j.GlobalObjects.Users[strconv.FormatUint(t.UserId, 10)]
		l[i] = Record{
			Id:            t.Id,
			Username:      u.ScreenName,
			CreatedAt:     Iso8601Date(time.Time(t.CreatedAt)),
			FullText:      t.FullText,
			RetweetCount:  t.RetweetCount,
			FavoriteCount: t.FavoriteCount,
			ReplyCount:    t.ReplyCount,
			QuoteCount:    t.QuoteCount,
			Geo:           t.Geo,
			Coodinates:    t.Coodinates,
			Lang:          t.Lang,
			Source:        t.Source,
		}
	}

	return l
}
