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

import (
	"reflect"
	"testing"
	"time"

	"github.com/akiomik/get-old-tweets/twitter/json"
)

func TestReverseSortedTweetIdsWhenInstructionsAreEmpty(t *testing.T) {
	j := &json.Adaptive{
		Timeline: json.Timeline{
			Instructions: []json.Instruction{},
		},
	}

	expected := []string{}
	actual := ReverseSortedTweetIds(j)
	if !reflect.DeepEqual(actual, expected) {
		t.Errorf("Expect %v, got %v", expected, actual)
		return
	}
}

func TestReverseSortedTweetIdsWhenInstructionsAreNotEmpty(t *testing.T) {
	j := &json.Adaptive{
		Timeline: json.Timeline{
			Instructions: []json.Instruction{
				json.Instruction{
					AddEntries: json.AddEntries{
						Entries: []json.Entry{
							json.Entry{
								EntryId:   "sq-I-t-300",
								SortIndex: "999990",
								Content: json.Content{
									Item: json.Item{
										Content: json.ItemContent{
											Tweet: json.ContentTweet{
												Id:          "300",
												DisplayType: "Tweet",
											},
										},
									},
								},
							},
							json.Entry{
								EntryId:   "sq-I-t-100",
								SortIndex: "999980",
								Content: json.Content{
									Item: json.Item{
										Content: json.ItemContent{
											Tweet: json.ContentTweet{
												Id:          "100",
												DisplayType: "Tweet",
											},
										},
									},
								},
							},
							json.Entry{
								EntryId:   "sq-I-t-200",
								SortIndex: "999970",
								Content: json.Content{
									Item: json.Item{
										Content: json.ItemContent{
											Tweet: json.ContentTweet{
												Id:          "200",
												DisplayType: "Tweet",
											},
										},
									},
								},
							},
							json.Entry{
								EntryId:   "sq-cursor-top",
								SortIndex: "999999",
								Content: json.Content{
									Operation: json.Operation{
										Cursor: json.Cursor{
											Value: "refresh:foobar",
										},
									},
								},
							},
							json.Entry{
								EntryId:   "sq-cursor-bottom",
								SortIndex: "0",
								Content: json.Content{
									Operation: json.Operation{
										Cursor: json.Cursor{
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

	expected := []string{"300", "100", "200"}
	actual := ReverseSortedTweetIds(j)
	if !reflect.DeepEqual(actual, expected) {
		t.Errorf("Expect %v, got %v", expected, actual)
		return
	}
}

func TestNewRecordsFromAdaptive(t *testing.T) {
	j := json.Adaptive{
		GlobalObjects: json.GlobalObjects{
			Tweets: map[string]json.Tweet{
				"1000": json.Tweet{
					Id:            1000,
					UserId:        2000,
					CreatedAt:     json.RubyDate(time.Date(2020, 9, 6, 0, 1, 2, 0, time.UTC)),
					FullText:      "To Sherlock Holmes she is always the woman.",
					RetweetCount:  3000,
					FavoriteCount: 4000,
					ReplyCount:    5000,
					QuoteCount:    6000,
					Lang:          "en",
				},
				"100": json.Tweet{
					Id:            100,
					UserId:        200,
					CreatedAt:     json.RubyDate(time.Date(2020, 9, 6, 0, 1, 2, 0, time.UTC)),
					FullText:      "To Sherlock Holmes she is always the woman.",
					RetweetCount:  300,
					FavoriteCount: 400,
					ReplyCount:    500,
					QuoteCount:    600,
					Coordinates: &json.Coordinates{
						Type:        "Point",
						Coordinates: json.LongLat{-73.9998279, 40.74118764},
					},
					Lang: "en",
				},
			},
			Users: map[string]json.User{
				"2000": json.User{
					Id:         2000,
					Name:       "Watson",
					ScreenName: "watson1",
				},
				"200": json.User{
					Id:         200,
					Name:       "Watson",
					ScreenName: "watson2",
				},
			},
		},
		Timeline: json.Timeline{
			Instructions: []json.Instruction{
				json.Instruction{
					AddEntries: json.AddEntries{
						Entries: []json.Entry{
							json.Entry{
								EntryId:   "sq-I-t-1000",
								SortIndex: "999990",
								Content: json.Content{
									Item: json.Item{
										Content: json.ItemContent{
											Tweet: json.ContentTweet{
												Id:          "1000",
												DisplayType: "Tweet",
											},
										},
									},
								},
							},
							json.Entry{
								EntryId:   "sq-I-t-100",
								SortIndex: "999980",
								Content: json.Content{
									Item: json.Item{
										Content: json.ItemContent{
											Tweet: json.ContentTweet{
												Id:          "100",
												DisplayType: "Tweet",
											},
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

	latitude := 40.74118764
	longitude := -73.9998279
	expected := []Record{
		Record{
			Id:            1000,
			Username:      "watson1",
			CreatedAt:     Iso8601Date(time.Date(2020, 9, 6, 0, 1, 2, 0, time.UTC)),
			FullText:      "To Sherlock Holmes she is always the woman.",
			RetweetCount:  3000,
			FavoriteCount: 4000,
			ReplyCount:    5000,
			QuoteCount:    6000,
			Latitude:      nil,
			Longitude:     nil,
			Lang:          "en",
		},
		Record{
			Id:            100,
			Username:      "watson2",
			CreatedAt:     Iso8601Date(time.Date(2020, 9, 6, 0, 1, 2, 0, time.UTC)),
			FullText:      "To Sherlock Holmes she is always the woman.",
			RetweetCount:  300,
			FavoriteCount: 400,
			ReplyCount:    500,
			QuoteCount:    600,
			Latitude:      &latitude,
			Longitude:     &longitude,
			Lang:          "en",
		},
	}

	actual := NewRecordsFromAdaptive(&j)
	if !reflect.DeepEqual(actual, expected) {
		t.Errorf("Expect %+v, got %+v", expected, actual)
		return
	}
}
