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

package export

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFilter(t *testing.T) {
	examples := map[string]struct {
		input    []int
		f        func(n int) bool
		expected []int
	}{
		"even": {
			input:    []int{1, 2, 3, 4, 5},
			f:        func(n int) bool { return n%2 == 0 },
			expected: []int{2, 4},
		},
		"odd": {
			input:    []int{1, 2, 3, 4, 5},
			f:        func(n int) bool { return n%2 == 1 },
			expected: []int{1, 3, 5},
		},
		"empty": {
			input:    []int{},
			f:        func(n int) bool { return n%2 == 0 },
			expected: []int{},
		},
		"noop": {
			input:    []int{1, 2, 3, 4, 5},
			f:        func(n int) bool { return true },
			expected: []int{1, 2, 3, 4, 5},
		},
	}

	for name, e := range examples {
		t.Run(name, func(t *testing.T) {
			actual := Filter(e.input, e.f)
			assert.Equal(t, e.expected, actual)
		})
	}
}

func TestMap(t *testing.T) {
	examples := map[string]struct {
		input    []int
		f        func(n int) int
		expected []int
	}{
		"double": {
			input:    []int{1, 2, 3, 4, 5},
			f:        func(n int) int { return n * 2 },
			expected: []int{2, 4, 6, 8, 10},
		},
		"empty": {
			input:    []int{},
			f:        func(n int) int { return n * 2 },
			expected: []int{},
		},
		"id": {
			input:    []int{1, 2, 3, 4, 5},
			f:        func(n int) int { return n },
			expected: []int{1, 2, 3, 4, 5},
		},
	}

	for name, e := range examples {
		t.Run(name, func(t *testing.T) {
			actual := Map(e.input, e.f)
			assert.Equal(t, e.expected, actual)
		})
	}
}
