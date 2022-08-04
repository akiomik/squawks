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

package twitter

import (
	"fmt"
)

type BoundingBox struct {
	Type       string        `json:"type"`
	Coodinates [][][]float64 `json:"coordinates"`
}

type Place struct {
	Id          string      `json:"id"`
	Url         string      `json:"url"`
	PlaceType   string      `json:"place_type"`
	Name        string      `json:"name"`
	FullName    string      `json:"full_name"`
	CountryCode string      `json:"country_code"`
	Country     string      `json:"country"`
	BoundingBox BoundingBox `json:"bounding_box"`
}

type Tweet struct {
	Id            uint64   `json:"id"`
	UserId        uint64   `json:"user_id"`
	FullText      string   `json:"full_text"`
	RetweetCount  uint64   `json:"retweet_count"`
	FavoriteCount uint64   `json:"favorite_count"`
	ReplyCount    uint64   `json:"reply_count"`
	QuoteCount    uint64   `json:"quote_count"`
	Geo           string   `json:"geo"`
	Coodinates    string   `json:"coordinates"`
	Place         Place    `json:"place"`
	Lang          string   `json:"lang"`
	Source        string   `json:"source"`
	CreatedAt     RubyDate `json:"created_at"`
}

type User struct {
	Id              uint64   `json:"id"`
	Name            string   `json:"name"`
	ScreenName      string   `json:"screen_name"`
	Location        string   `json:"location"`
	Description     string   `json:"description"`
	Url             string   `json:"url"`
	FollowersCount  uint64   `json:"followers_count"`
	FriendsCount    uint64   `json:"friends_count"`
	ListedCount     uint64   `json:"listed_count"`
	FavouritesCount uint64   `json:"favourites_count"`
	StatusesCount   uint64   `json:"statuses_count"`
	MediaCount      uint64   `json:"media_count"`
	Verified        bool     `json:"verified"`
	CreatedAt       RubyDate `json:"created_at"`
}

type GlobalObjects struct {
	Tweets map[string]Tweet `json:"tweets"`
	Users  map[string]User  `json:"users"`
}

type Cursor struct {
	Value      string `json:"value"`
	CursorType string `json:"cursorType"`
}

type Operation struct {
	Cursor Cursor `json:"cursor"`
}

type Content struct {
	Operation Operation `json:"operation"`
}

type Entry struct {
	EntryId   string  `json:"entryId"`
	SortIndex string  `json:"sortIndex"`
	Content   Content `json:"content"`
}

type AddEntries struct {
	Entries []Entry `json:"entries"`
}

type ReplaceEntry struct {
	EntryIdToReplace string `json:"entryIdToReplace"`
	Entry            Entry  `json:"entry"`
}

type Instruction struct {
	AddEntries   AddEntries   `json:"addEntries"`
	ReplaceEntry ReplaceEntry `json:"replaceEntry"`
}

type Timeline struct {
	Id           string        `json:"id"`
	Instructions []Instruction `json:"instructions"`
}

type Error struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

type AdaptiveJson struct {
	GlobalObjects GlobalObjects `json:"globalObjects"`
	Timeline      Timeline      `json:"timeline"`
	Errors        []Error       `json:"errors"`
}

func (j *AdaptiveJson) FindCursor() (string, error) {
	for _, i := range j.Timeline.Instructions {
		if i.ReplaceEntry.EntryIdToReplace == "sq-cursor-bottom" {
			return i.ReplaceEntry.Entry.Content.Operation.Cursor.Value, nil
		}

		for _, e := range i.AddEntries.Entries {
			if e.EntryId == "sq-cursor-bottom" {
				return e.Content.Operation.Cursor.Value, nil
			}
		}
	}

	return "", fmt.Errorf("cursor not found")
}
