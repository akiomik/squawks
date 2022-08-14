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
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/akiomik/squawks/api/json"
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
	Latitude      *float64
	Longitude     *float64
	Lang          string
	Source        string
}

func ReverseSortedTweetIds(j *json.Adaptive) []string {
	is := j.Timeline.Instructions
	if len(is) == 0 {
		return make([]string, 0)
	}

	es := Filter(is[0].AddEntries.Entries, func(e json.Entry) bool {
		return strings.HasPrefix(e.EntryId, "sq-I-t-") &&
			len(e.Content.Item.Content.Tweet.Id) != 0 &&
			e.Content.Item.Content.Tweet.DisplayType == "Tweet"
	})

	sort.Slice(es, func(i int, j int) bool {
		return es[i].SortIndex > es[j].SortIndex
	})

	return Map(es, func(e json.Entry) string {
		return e.Content.Item.Content.Tweet.Id
	})
}

func NewRecordsFromAdaptive(j *json.Adaptive) []Record {
	return Map(ReverseSortedTweetIds(j), func(id string) Record {
		t := j.GlobalObjects.Tweets[id]
		u := j.GlobalObjects.Users[strconv.FormatUint(t.UserId, 10)]

		var latitude *float64
		var longitude *float64
		if t.Coordinates != nil {
			lat := t.Coordinates.Coordinates.Latitude()
			long := t.Coordinates.Coordinates.Longitude()
			latitude = &lat
			longitude = &long
		}

		return Record{
			Id:            t.Id,
			Username:      u.ScreenName,
			CreatedAt:     Iso8601Date(time.Time(t.CreatedAt)),
			FullText:      t.FullText,
			RetweetCount:  t.RetweetCount,
			FavoriteCount: t.FavoriteCount,
			ReplyCount:    t.ReplyCount,
			QuoteCount:    t.QuoteCount,
			Latitude:      latitude,
			Longitude:     longitude,
			Lang:          t.Lang,
			Source:        t.Source,
		}
	})
}
