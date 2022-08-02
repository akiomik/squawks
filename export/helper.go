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

package export

import (
	"sort"

	"github.com/akiomik/get-old-tweets/twitter"
)

func Contains(ss []string, needle string) bool {
	for _, s := range ss {
		if s == needle {
			return true
		}
	}

	return false
}

func KeysOf(m map[string]twitter.Tweet) []string {
	keys := make([]string, len(m))

	i := 0
	for k := range m {
		keys[i] = k
		i++
	}

	return keys
}

func ReversedKeysOf(m map[string]twitter.Tweet) []string {
	ks := KeysOf(m)
	sort.Sort(sort.Reverse(sort.StringSlice(ks)))
	return ks
}
