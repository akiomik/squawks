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

package flags

import (
	"strings"
	"testing"
)

func TestAll(t *testing.T) {
	examples := map[string]struct {
		xs       []string
		f        func(x string) bool
		expected bool
	}{
		"true": {
			xs:       []string{"foo", "bar"},
			f:        func(x string) bool { return len(x) == 3 },
			expected: true,
		},
		"false": {
			xs:       []string{"foo", "bar"},
			f:        func(x string) bool { return strings.HasPrefix(x, "fo") },
			expected: false,
		},
		"empty": {
			xs:       []string{},
			f:        func(x string) bool { return len(x) == 3 },
			expected: true,
		},
	}

	for name, e := range examples {
		t.Run(name, func(t *testing.T) {
			actual := All(e.xs, e.f)
			if actual != e.expected {
				t.Errorf("Expect %v, got %v", e.expected, actual)
				return
			}
		})
	}
}

func TestAny(t *testing.T) {
	examples := map[string]struct {
		xs       []string
		f        func(x string) bool
		expected bool
	}{
		"true": {
			xs:       []string{"foo", "bar"},
			f:        func(x string) bool { return strings.HasPrefix(x, "fo") },
			expected: true,
		},
		"false": {
			xs:       []string{"foo", "bar"},
			f:        func(x string) bool { return strings.HasPrefix(x, "oo") },
			expected: false,
		},
		"empty": {
			xs:       []string{},
			f:        func(x string) bool { return strings.HasPrefix(x, "fo") },
			expected: false,
		},
	}

	for name, e := range examples {
		t.Run(name, func(t *testing.T) {
			actual := Any(e.xs, e.f)
			if actual != e.expected {
				t.Errorf("Expect %v, got %v", e.expected, actual)
				return
			}
		})
	}
}

func TestIncludes(t *testing.T) {
	examples := map[string]struct {
		xs       []string
		n        string
		expected bool
	}{
		"true": {
			xs:       []string{"foo", "bar"},
			n:        "foo",
			expected: true,
		},
		"false": {
			xs:       []string{"foo", "bar"},
			n:        "baz",
			expected: false,
		},
		"empty": {
			xs:       []string{},
			n:        "foo",
			expected: false,
		},
	}

	for name, e := range examples {
		t.Run(name, func(t *testing.T) {
			actual := Includes(e.xs, e.n)
			if actual != e.expected {
				t.Errorf("Expect %v, got %v", e.expected, actual)
				return
			}
		})
	}
}
