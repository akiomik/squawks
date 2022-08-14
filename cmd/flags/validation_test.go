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
	"strings"
	"testing"

	"github.com/spf13/pflag"
	"github.com/stretchr/testify/assert"
)

func TestStringSliceWithValidationVarP(t *testing.T) {
	examples := map[string]struct {
		input    []string
		expected []string
		msg      string
	}{
		"none": {
			input:    []string{},
			expected: []string{},
			msg:      "",
		},
		"empty": {
			input:    []string{"--arg="},
			expected: []string{},
			msg:      "",
		},
		"valid-single": {
			input:    []string{"--arg=foo"},
			expected: []string{"foo"},
			msg:      "",
		},
		"invalid-single": {
			input:    []string{"--arg=bar"},
			expected: []string{},
			msg:      `invalid argument "bar" for "--arg" flag: string starting with "foo" are supported`,
		},
		"valid-multiple-flags": {
			input:    []string{"--arg=foo", "--arg=foobar"},
			expected: []string{"foo|foobar"},
			msg:      "",
		},
		"invalid-multiple-flags": {
			input:    []string{"--arg=foo", "--arg=bar"},
			expected: []string{"foo"},
			msg:      `invalid argument "bar" for "--arg" flag: string starting with "foo" are supported`,
		},
		"valid-multiple-values": {
			input:    []string{"--arg=foo,foobar"},
			expected: []string{"foo|foobar"},
			msg:      "",
		},
		"invalid-multiple-values": {
			input:    []string{"--arg=foo,bar"},
			expected: []string{},
			msg:      `invalid argument "foo,bar" for "--arg" flag: string starting with "foo" are supported`,
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

			if len(e.msg) == 0 {
				assert.NoError(t, err)
			} else {
				assert.EqualError(t, err, e.msg)
			}

			actual, err := flags.GetStringSlice("arg")
			assert.NoError(t, err)
			assert.Equal(t, e.expected, actual)
		})
	}
}
