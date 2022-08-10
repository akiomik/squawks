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
	"strings"
)

type Query struct {
	Text     string
	Since    string
	Until    string
	From     string
	To       string
	Lang     string
	Filters  []string
	Excludes []string
	Geocode  string
	Url      string
}

func (q *Query) Encode() string {
	var ss []string

	if len(q.Text) != 0 {
		ss = append(ss, q.Text)
	}

	if len(q.Since) != 0 {
		ss = append(ss, "since:"+q.Since)
	}

	if len(q.Until) != 0 {
		ss = append(ss, "until:"+q.Until)
	}

	if len(q.From) != 0 {
		ss = append(ss, "from:"+q.From)
	}

	if len(q.To) != 0 {
		ss = append(ss, "to:"+q.To)
	}

	if len(q.Lang) != 0 {
		ss = append(ss, "lang:"+q.Lang)
	}

	for _, filter := range q.Filters {
		ss = append(ss, "filter:"+filter)
	}

	for _, exclude := range q.Excludes {
		ss = append(ss, "exclude:"+exclude)
	}

	if len(q.Geocode) != 0 {
		ss = append(ss, "geocode:"+q.Geocode)
	}

	if len(q.Url) != 0 {
		ss = append(ss, "url:"+q.Url)
	}

	return strings.Join(ss[:], " ")
}

func (q *Query) IsEmpty() bool {
	return len(q.Encode()) == 0
}
