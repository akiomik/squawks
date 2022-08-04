// copyright 2022 akiomi kamakura
//
// licensed under the apache license, version 2.0 (the "license");
// you may not use this file except in compliance with the license.
// you may obtain a copy of the license at
//
//     http://www.apache.org/licenses/license-2.0
//
// unless required by applicable law or agreed to in writing, software
// distributed under the license is distributed on an "as is" basis,
// without warranties or conditions of any kind, either express or implied.
// see the license for the specific language governing permissions and
// limitations under the license.

//go:build small
// +build small

package export

import (
	"reflect"
	"testing"

	"github.com/akiomik/get-old-tweets/twitter"
)

func TestContains(t *testing.T) {
	ss := []string{"foo", "bar", "baz"}

	examples := map[string]struct {
		value    string
		expected bool
	}{
		"foo": {
			value:    "foo",
			expected: true,
		},
		"bar": {
			value:    "bar",
			expected: true,
		},
		"baz": {
			value:    "baz",
			expected: true,
		},
		"qux": {
			value:    "qux",
			expected: false,
		},
	}

	for name, e := range examples {
		t.Run(name, func(t *testing.T) {
			actual := Contains(ss, e.value)
			if actual != e.expected {
				t.Errorf("Expect Contains(%v, %s) = %v, but got %v", ss, e.value, e.expected, actual)
				return
			}
		})
	}
}

func TestKeysOf(t *testing.T) {
	m := map[string]twitter.Tweet{
		"z":  twitter.Tweet{},
		"zz": twitter.Tweet{},
		"a":  twitter.Tweet{},
		"c":  twitter.Tweet{},
		"c1": twitter.Tweet{},
		"c0": twitter.Tweet{},
	}

	ks := KeysOf(m)

	if len(ks) != len(m) {
		t.Errorf("Expect len(KeysOf()) to be %d, but got %d", len(m), len(ks))
		return
	}

	for k, _ := range m {
		if !Contains(ks, k) {
			t.Errorf("Expect KeysOf() to include \"%s\", but none", k)
			return
		}
	}
}

func TestReversedKeysOf(t *testing.T) {
	m := map[string]twitter.Tweet{
		"z":  twitter.Tweet{},
		"zz": twitter.Tweet{},
		"a":  twitter.Tweet{},
		"c":  twitter.Tweet{},
		"c1": twitter.Tweet{},
		"c0": twitter.Tweet{},
	}

	expected := []string{
		"zz",
		"z",
		"c1",
		"c0",
		"c",
		"a",
	}
	actual := ReversedKeysOf(m)

	if !reflect.DeepEqual(actual, expected) {
		t.Errorf("Expect ReversedKeysOf() = %v, but got %v", expected, actual)
		return
	}
}
