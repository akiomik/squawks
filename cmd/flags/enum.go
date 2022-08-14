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

package flags

import (
	"fmt"
	"strings"

	"github.com/spf13/pflag"
)

func StringSliceEnumVarP(flags *pflag.FlagSet, p *[]string, name string, shorthand string, defaults []string, usage string, options []string) *pflag.Flag {
	formattedOptions := "[" + strings.Join(options[:], "|") + "]"
	validator := func(values []string) error {
		if All(values, func(value string) bool { return Includes(options, value) }) {
			return nil
		}

		return fmt.Errorf(`valid values are %s`, formattedOptions)
	}

	return StringSliceWithValidationVarP(flags, p, name, shorthand, defaults, usage+" "+formattedOptions, validator)
}
