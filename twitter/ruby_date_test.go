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
	"encoding/json"
	"testing"
	"time"
)

type TestSchema struct {
	CreatedAt RubyDate `json:"created_at"`
}

func TestString(t *testing.T) {
	d := RubyDate(time.Date(2013, 8, 19, 2, 4, 28, 0, time.UTC))
	expected := "Mon Aug 19 02:04:28 +0000 2013"
	actual := d.String()
	if actual != expected {
		t.Errorf("Expect RubyDate#String() = \"%s\", but got \"%s\"", expected, actual)
		return
	}
}

func TestIso8601(t *testing.T) {
	d := RubyDate(time.Date(2013, 8, 19, 2, 4, 28, 0, time.UTC))
	expected := "2013-08-19T02:04:28+00:00"
	actual := d.Iso8601()
	if actual != expected {
		t.Errorf("Expect RubyDate#Iso8601() = \"%s\", but got \"%s\"", expected, actual)
		return
	}
}

func TestEqual(t *testing.T) {
	examples := map[string]struct {
		this     RubyDate
		that     RubyDate
		expected bool
	}{
		"true": {
			this:     RubyDate(time.Date(2013, 8, 19, 2, 4, 28, 0, time.UTC)),
			that:     RubyDate(time.Date(2013, 8, 19, 2, 4, 28, 0, time.UTC)),
			expected: true,
		},
		"false": {
			this:     RubyDate(time.Date(2013, 8, 19, 2, 4, 28, 0, time.UTC)),
			that:     RubyDate(time.Date(2013, 8, 19, 2, 4, 28, 1, time.UTC)),
			expected: false,
		},
	}

	for name, e := range examples {
		t.Run(name, func(t *testing.T) {
			actual := e.this.Equal(e.that)
			if actual != e.expected {
				t.Errorf("Expect RubyDate#Equal() = %v, but got %v", e.expected, actual)
				return
			}
		})
	}
}

func TestUnmarshall(t *testing.T) {
	jsonString := `{
    "created_at": "Mon Aug 19 02:04:28 +0000 2013"
  }`

	schema := new(TestSchema)
	err := json.Unmarshal([]byte(jsonString), &schema)
	if err != nil {
		t.Errorf("Expect Query#Encode() not to return error object, but got \"%v\"", err)
		return
	}

	expected := RubyDate(time.Date(2013, 8, 19, 2, 4, 28, 0, time.UTC))
	actual := schema.CreatedAt
	if !actual.Equal(expected) {
		t.Errorf("Expect RubyDate#Unmarshal() = \"%v\", but got \"%v\"", expected, actual)
		return
	}
}
