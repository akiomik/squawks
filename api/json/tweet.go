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

package json

// https://developer.twitter.com/en/docs/twitter-api/v1/data-dictionary/object-model/tweet

type Tweet struct {
	Id            uint64       `json:"id"`
	UserId        uint64       `json:"user_id"`
	FullText      string       `json:"full_text"`
	RetweetCount  uint64       `json:"retweet_count"`
	FavoriteCount uint64       `json:"favorite_count"`
	ReplyCount    uint64       `json:"reply_count"`
	QuoteCount    uint64       `json:"quote_count"`
	Geo           *Geo         `json:"geo"` // deprecated
	Coordinates   *Coordinates `json:"coordinates"`
	Place         *Place       `json:"place"`
	Lang          string       `json:"lang"`
	Source        string       `json:"source"`
	CreatedAt     RubyDate     `json:"created_at"`
}
