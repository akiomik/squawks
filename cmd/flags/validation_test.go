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
	"fmt"
	"reflect"
	"strings"
	"testing"

	"github.com/spf13/pflag"
)

func TestStringSliceWithValidationVarP(t *testing.T) {
	examples := map[string]struct {
		input    []string
		expected []string
		error    string
	}{
		"none": {
			input:    []string{},
			expected: []string{},
			error:    "",
		},
		"empty": {
			input:    []string{"--arg="},
			expected: []string{},
			error:    "",
		},
		"valid-single": {
			input:    []string{"--arg=foo"},
			expected: []string{"foo"},
			error:    "",
		},
		"invalid-single": {
			input:    []string{"--arg=bar"},
			expected: []string{},
			error:    `invalid argument "bar" for "--arg" flag: string starting with "foo" are supported`,
		},
		"valid-multiple-flags": {
			input:    []string{"--arg=foo", "--arg=foobar"},
			expected: []string{"foo|foobar"},
			error:    "",
		},
		"invalid-multiple-flags": {
			input:    []string{"--arg=foo", "--arg=bar"},
			expected: []string{"foo"},
			error:    `invalid argument "bar" for "--arg" flag: string starting with "foo" are supported`,
		},
		"valid-multiple-values": {
			input:    []string{"--arg=foo,foobar"},
			expected: []string{"foo|foobar"},
			error:    "",
		},
		"invalid-multiple-values": {
			input:    []string{"--arg=foo,bar"},
			expected: []string{},
			error:    `invalid argument "foo,bar" for "--arg" flag: string starting with "foo" are supported`,
		},
	}

	for name, e := range examples {
		t.Run(name, func(t *testing.T) {
			var args []string

			flags := pflag.NewFlagSet("test", pflag.ContinueOnError)
			prefix := "foo"
			validator := func(values []string) error {
				if All(values, func(value string) bool { return strings.HasPrefix(value, prefix) }) {
					return nil
				}

				return fmt.Errorf(`string starting with "%s" are supported`, prefix)
			}

			StringSliceWithValidationVarP(flags, &args, "arg", "", []string{}, "args for testing", validator)
			err := flags.Parse(e.input)

			if err == nil {
				if e.error != "" {
					t.Errorf(`Expect error %v, got nil`, e.error)
					return
				}
			} else {
				if e.error != err.Error() {
					t.Errorf(`Expect "%s", got "%v"`, e.error, err)
					return
				}
			}

			actual, err := flags.GetStringSlice("arg")
			if err != nil {
				t.Errorf(`Expect no error, got "%v"`, err)
				return
			}

			if !reflect.DeepEqual(actual, e.expected) {
				t.Errorf("Expect %v, got %v", e.expected, actual)
				return
			}
		})
	}
}
