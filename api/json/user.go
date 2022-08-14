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
