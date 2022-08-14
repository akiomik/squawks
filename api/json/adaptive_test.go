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

package json

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

var (
	sqCursorTopEntry = Entry{
		EntryId: "sq-cursor-top",
		Content: Content{
			Operation: Operation{
				Cursor: Cursor{
					Value: "refresh:foobar",
				},
			},
		},
	}

	sqCursorBottomEntry = Entry{
		EntryId: "sq-cursor-bottom",
		Content: Content{
			Operation: Operation{
				Cursor: Cursor{
					Value: "scroll:foobar",
				},
			},
		},
	}
)

func TestFindCursorWhenReplaceEntryExists(t *testing.T) {
	j := Adaptive{
		Timeline: Timeline{
			Instructions: []Instruction{
				Instruction{
					ReplaceEntry: ReplaceEntry{
						EntryIdToReplace: "sq-cursor-top",
						Entry:            sqCursorTopEntry,
					},
				},
				Instruction{
					ReplaceEntry: ReplaceEntry{
						EntryIdToReplace: "sq-cursor-bottom",
						Entry:            sqCursorBottomEntry,
					},
				},
			},
		},
	}

	expected := "scroll:foobar"
	actual, err := j.FindCursor()
	assert.NoError(t, err)
	assert.Equal(t, expected, actual)
}

func TestFindCursorWhenAddEntriesExist(t *testing.T) {
	j := Adaptive{
		Timeline: Timeline{
			Instructions: []Instruction{
				Instruction{
					AddEntries: AddEntries{
						Entries: []Entry{
							sqCursorTopEntry,
							sqCursorBottomEntry,
						},
					},
				},
			},
		},
	}

	expected := "scroll:foobar"
	actual, err := j.FindCursor()
	assert.NoError(t, err)
	assert.Equal(t, expected, actual)
}

func TestFindCursorWhenNoCursorFound(t *testing.T) {
	j := Adaptive{}

	_, err := j.FindCursor()
	assert.Error(t, err)
}
