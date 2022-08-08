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

package export

func Filter[A any](xs []A, f func(A) bool) []A {
	acc := []A{}

	for _, x := range xs {
		if f(x) {
			acc = append(acc, x)
		}
	}

	return acc
}

func Map[A any, B any](xs []A, f func(A) B) []B {
	newXs := make([]B, len(xs))
	for i, x := range xs {
		newXs[i] = f(x)
	}

	return newXs
}
