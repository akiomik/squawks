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
	"encoding/csv"
	"strings"

	"github.com/spf13/pflag"
)

func StringSliceWithValidationVarP(flags *pflag.FlagSet, p *[]string, name string, shorthand string, defaults []string, usage string, validator func([]string) error) *pflag.Flag {
	*p = defaults
	v := &StringSliceValueWithValidation{Values: p, Validator: validator}
	return flags.VarPF(v, name, shorthand, usage)
}

type StringSliceValueWithValidation struct {
	Values    *[]string
	Validator func([]string) error
}

func (e *StringSliceValueWithValidation) Set(v string) error {
	if v == "" {
		return nil
	}

	r := csv.NewReader(strings.NewReader(v))
	vs, err := r.Read()
	if err != nil {
		return err
	}

	err = e.Validator(vs)
	if err != nil {
		return err
	}

	*e.Values = append(*e.Values, vs...)
	return nil
}

func (e *StringSliceValueWithValidation) String() string {
	vs := *e.Values
	return "[" + strings.Join(vs[:], "|") + "]"
}

func (e *StringSliceValueWithValidation) Type() string {
	return "stringSlice"
}
