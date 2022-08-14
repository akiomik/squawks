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

func All[T any](xs []T, f func(x T) bool) bool {
	for _, x := range xs {
		if !f(x) {
			return false
		}
	}

	return true
}

func Any[T any](xs []T, f func(x T) bool) bool {
	for _, x := range xs {
		if f(x) {
			return true
		}
	}

	return false
}

func Includes[T comparable](xs []T, n T) bool {
	return Any(xs, func(x T) bool { return x == n })
}
