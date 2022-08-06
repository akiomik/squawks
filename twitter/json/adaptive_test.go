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
)

func TestFindCursorWhenReplaceEntryExists(t *testing.T) {
	j := Adaptive{
		Timeline: Timeline{
			Instructions: []Instruction{
				Instruction{
					ReplaceEntry: ReplaceEntry{
						EntryIdToReplace: "sq-cursor-top",
						Entry: Entry{
							EntryId: "sq-cursor-top",
							Content: Content{
								Operation: Operation{
									Cursor: Cursor{
										Value: "refresh:foobar",
									},
								},
							},
						},
					},
				},
				Instruction{
					ReplaceEntry: ReplaceEntry{
						EntryIdToReplace: "sq-cursor-bottom",
						Entry: Entry{
							EntryId: "sq-cursor-bottom",
							Content: Content{
								Operation: Operation{
									Cursor: Cursor{
										Value: "scroll:foobar",
									},
								},
							},
						},
					},
				},
			},
		},
	}

	expected := "scroll:foobar"
	actual, err := j.FindCursor()
	if err != nil {
		t.Errorf("Expect nil, got %v", err)
		return
	}

	if actual != expected {
		t.Errorf(`Expect "%s", got "%s"`, expected, actual)
		return
	}
}

func TestFindCursorWhenAddEntriesExist(t *testing.T) {
	j := Adaptive{
		Timeline: Timeline{
			Instructions: []Instruction{
				Instruction{
					AddEntries: AddEntries{
						Entries: []Entry{
							Entry{
								EntryId: "sq-cursor-top",
								Content: Content{
									Operation: Operation{
										Cursor: Cursor{
											Value: "refresh:foobar",
										},
									},
								},
							},
							Entry{
								EntryId: "sq-cursor-bottom",
								Content: Content{
									Operation: Operation{
										Cursor: Cursor{
											Value: "scroll:foobar",
										},
									},
								},
							},
						},
					},
				},
			},
		},
	}

	expected := "scroll:foobar"
	actual, err := j.FindCursor()
	if err != nil {
		t.Errorf("Expect nil, got %v", err)
		return
	}

	if actual != expected {
		t.Errorf(`Expect "%s", got "%s"`, expected, actual)
		return
	}
}

func TestFindCursorWhenNoCursorFound(t *testing.T) {
	j := Adaptive{}

	_, err := j.FindCursor()
	if err == nil {
		t.Errorf("Expect error object, got nil")
		return
	}
}
