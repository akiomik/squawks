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

//go:build small
// +build small

package twitter

import (
	"testing"
)

func TestEncode(t *testing.T) {
	examples := map[string]struct {
		text     string
		since    string
		until    string
		from     string
		to       string
		lang     string
		expected string
	}{
		"none": {
			text:     "",
			since:    "",
			until:    "",
			from:     "",
			to:       "",
			lang:     "",
			expected: "",
		},
		"all": {
			text:     "foo bar",
			since:    "2020-09-06",
			until:    "2020-09-07",
			from:     "foo",
			to:       "bar",
			lang:     "ja",
			expected: "foo bar since:2020-09-06 until:2020-09-07 from:foo to:bar lang:ja",
		},
	}

	for name, e := range examples {
		t.Run(name, func(t *testing.T) {
			q := Query{
				Text:  e.text,
				Since: e.since,
				Until: e.until,
				From:  e.from,
				To:    e.to,
				Lang:  e.lang,
			}

			actual := q.Encode()

			if actual != e.expected {
				t.Errorf(`Expect "%s", got "%s"`, e.expected, actual)
				return
			}
		})
	}
}

func TestIsEmpty(t *testing.T) {
	examples := map[string]struct {
		text     string
		since    string
		until    string
		from     string
		to       string
		lang     string
		expected bool
	}{
		"none": {
			text:     "",
			since:    "",
			until:    "",
			from:     "",
			to:       "",
			lang:     "",
			expected: true,
		},
		"all": {
			text:     "foo bar",
			since:    "2020-09-06",
			until:    "2020-09-07",
			from:     "foo",
			to:       "bar",
			lang:     "ja",
			expected: false,
		},
	}

	for name, e := range examples {
		t.Run(name, func(t *testing.T) {
			q := Query{
				Text:  e.text,
				Since: e.since,
				Until: e.until,
				From:  e.from,
				To:    e.to,
				Lang:  e.lang,
			}

			actual := q.IsEmpty()

			if actual != e.expected {
				t.Errorf("Expect %v, got %v", e.expected, actual)
				return
			}
		})
	}
}
