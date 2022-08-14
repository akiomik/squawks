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

package json

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestLatLongString(t *testing.T) {
	examples := map[string]struct {
		coordinate LatLong
		expected   string
	}{
		"normal": {
			coordinate: LatLong{40.74118764, -73.9998279},
			expected:   "40.74118764,-73.9998279",
		},
		"zero": {
			coordinate: LatLong{},
			expected:   "0,0",
		},
	}

	for name, e := range examples {
		t.Run(name, func(t *testing.T) {
			actual := e.coordinate.String()
			assert.Equal(t, e.expected, actual)
		})
	}
}

func TestLongLatString(t *testing.T) {
	examples := map[string]struct {
		coordinate LongLat
		expected   string
	}{
		"normal": {
			coordinate: LongLat{-73.9998279, 40.74118764},
			expected:   "-73.9998279,40.74118764",
		},
		"zero": {
			coordinate: LongLat{},
			expected:   "0,0",
		},
	}

	for name, e := range examples {
		t.Run(name, func(t *testing.T) {
			actual := e.coordinate.String()
			assert.Equal(t, e.expected, actual)
		})
	}
}
