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

type Tweet struct {
	Id            int64    `json:"id"`
	UserId        int64    `json:"user_id"`
	FullText      string   `json:"full_text"`
	RetweetCount  int64    `json:"retweet_count"`
	FavoriteCount int64    `json:"favorite_count"`
	ReplyCount    int64    `json:"reply_count"`
	QuoteCount    int64    `json:"quote_count"`
	Geo           string   `json:"geo"`
	Coodinates    string   `json:"coordinates"`
	Place         string   `json:"place"`
	Lang          string   `json:"lang"`
	Source        string   `json:"source"`
	CreatedAt     RubyDate `json:"created_at"`
}

type User struct {
	Id              int64    `json:"id"`
	Name            string   `json:"name"`
	ScreenName      string   `json:"screen_name"`
	Location        string   `json:"location"`
	Description     string   `json:"description"`
	Url             string   `json:"url"`
	FollowersCount  int64    `json:"followers_count"`
	FriendsCount    int64    `json:"friends_count"`
	ListedCount     int64    `json:"listed_count"`
	FavouritesCount int64    `json:"favourites_count"`
	StatusesCount   int64    `json:"statuses_count"`
	MediaCount      int64    `json:"media_count"`
	Verified        bool     `json:"verified"`
	CreatedAt       RubyDate `json:"created_at"`
}

type GlobalObjects struct {
	Tweets map[string]Tweet `json:"tweets"`
	Users  map[string]User  `json:"users"`
}

type Error struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

type AdaptiveJson struct {
	GlobalObjects GlobalObjects `json:"globalObjects"`
	Errors        []Error       `json:"errors"`
}
