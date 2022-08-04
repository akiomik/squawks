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
		expected string
	}{
		"none": {
			text:     "",
			since:    "",
			until:    "",
			from:     "",
			to:       "",
			expected: "",
		},
		"all": {
			text:     "foo bar",
			since:    "2020-09-06",
			until:    "2020-09-07",
			from:     "foo",
			to:       "bar",
			expected: "foo+bar+since%3A2020-09-06+until%3A2020-09-07+from%3Afoo+to%3Abar",
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
			}

			actual := q.Encode()

			if actual != e.expected {
				t.Errorf("Expect Query#Encode() = \"%s\", but got \"%s\"", e.expected, actual)
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
		expected bool
	}{
		"none": {
			text:     "",
			since:    "",
			until:    "",
			from:     "",
			to:       "",
			expected: true,
		},
		"all": {
			text:     "foo bar",
			since:    "2020-09-06",
			until:    "2020-09-07",
			from:     "foo",
			to:       "bar",
			expected: false,
		},
	}

	for name, e := range examples {
		t.Run(name, func(t *testing.T) {
			q := Query{
				Text:  e.text,
				Since: e.since,
				Until: e.until,
			}

			actual := q.IsEmpty()

			if actual != e.expected {
				t.Errorf("Expect Query#IsEmpty() = %v, but got %v", e.expected, actual)
				return
			}
		})
	}
}
