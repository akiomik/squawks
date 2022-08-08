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

package json

import (
	"fmt"
)

type GlobalObjects struct {
	Tweets map[string]Tweet `json:"tweets"`
	Users  map[string]User  `json:"users"`
}

type Cursor struct {
	Value      string `json:"value"`
	CursorType string `json:"cursorType"`
}

type Operation struct {
	Cursor Cursor `json:"cursor"`
}

type ContentTweet struct {
	Id          string `json":id"`
	DisplayType string `json:"displayType"`
}

type ItemContent struct {
	Tweet ContentTweet `json:"tweet"`
}

type Item struct {
	Content ItemContent `json:"content"`
}

type Content struct {
	Operation Operation `json:"operation"`
	Item      Item      `json:"item"`
}

type Entry struct {
	EntryId   string  `json:"entryId"`
	SortIndex string  `json:"sortIndex"`
	Content   Content `json:"content"`
}

type AddEntries struct {
	Entries []Entry `json:"entries"`
}

type ReplaceEntry struct {
	EntryIdToReplace string `json:"entryIdToReplace"`
	Entry            Entry  `json:"entry"`
}

type Instruction struct {
	AddEntries   AddEntries   `json:"addEntries"`
	ReplaceEntry ReplaceEntry `json:"replaceEntry"`
}

type Timeline struct {
	Id           string        `json:"id"`
	Instructions []Instruction `json:"instructions"`
}

type Adaptive struct {
	GlobalObjects GlobalObjects `json:"globalObjects"`
	Timeline      Timeline      `json:"timeline"`
}

func (j *Adaptive) FindCursor() (string, error) {
	for _, i := range j.Timeline.Instructions {
		if i.ReplaceEntry.EntryIdToReplace == "sq-cursor-bottom" {
			return i.ReplaceEntry.Entry.Content.Operation.Cursor.Value, nil
		}

		for _, e := range i.AddEntries.Entries {
			if e.EntryId == "sq-cursor-bottom" {
				return e.Content.Operation.Cursor.Value, nil
			}
		}
	}

	return "", fmt.Errorf("cursor not found")
}
