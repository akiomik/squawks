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

package api

import (
	"errors"
	"reflect"
	"testing"

	"github.com/jarcoal/httpmock"

	"github.com/akiomik/squawks/api/json"
)

func TestSearch(t *testing.T) {
	examples := map[string]struct {
		opts       *SearchOptions
		url        string
		statusCode int
		response   string
		expected   *json.Adaptive
		errored    bool
	}{
		"without-cursor": {
			opts:       &SearchOptions{Query: Query{Text: "foo"}},
			url:        "https://twitter.com/i/api/2/search/adaptive.json?count=40&include_quote_count=true&include_reply_count=1&q=foo&query_source=typed_query&tweet_mode=extended&tweet_search_mode=live",
			statusCode: 200,
			response:   `{ "globalObjects": { "tweets": {}, "users": {} } }`,
			expected:   &json.Adaptive{GlobalObjects: json.GlobalObjects{Tweets: map[string]json.Tweet{}, Users: map[string]json.User{}}},
			errored:    false,
		},
		"with-cursor": {
			opts:       &SearchOptions{Query: Query{Text: "foo"}, Cursor: "scroll:deadbeef"},
			url:        "https://twitter.com/i/api/2/search/adaptive.json?count=40&cursor=scroll%3Adeadbeef&include_quote_count=true&include_reply_count=1&q=foo&query_source=typed_query&tweet_mode=extended&tweet_search_mode=live",
			statusCode: 200,
			response:   `{ "globalObjects": { "tweets": {}, "users": {} } }`,
			expected:   &json.Adaptive{GlobalObjects: json.GlobalObjects{Tweets: map[string]json.Tweet{}, Users: map[string]json.User{}}},
			errored:    false,
		},
		"failure": {
			opts:       &SearchOptions{Query: Query{Text: "foo"}},
			url:        "https://twitter.com/i/api/2/search/adaptive.json?count=40&include_quote_count=true&include_reply_count=1&q=foo&query_source=typed_query&tweet_mode=extended&tweet_search_mode=live",
			statusCode: 403,
			response:   `{ "errors": [{ "code": 200, "message": "forbidden" }] }`,
			expected:   nil,
			errored:    true,
		},
		"top": {
			opts:       &SearchOptions{Query: Query{Text: "foo"}, Top: true},
			url:        "https://twitter.com/i/api/2/search/adaptive.json?count=40&include_quote_count=true&include_reply_count=1&q=foo&query_source=typed_query&tweet_mode=extended",
			statusCode: 200,
			response:   `{ "globalObjects": { "tweets": {}, "users": {} } }`,
			expected:   &json.Adaptive{GlobalObjects: json.GlobalObjects{Tweets: map[string]json.Tweet{}, Users: map[string]json.User{}}},
			errored:    false,
		},
	}

	for name, e := range examples {
		t.Run(name, func(t *testing.T) {
			c := NewClient()

			httpmock.ActivateNonDefault(c.Client.GetClient())
			defer httpmock.DeactivateAndReset()

			httpmock.RegisterResponder("GET", e.url, NewJsonResponse(e.statusCode, e.response))

			actual, err := c.Search(e.opts)
			if e.errored != (err != nil) {
				t.Errorf("Expect error %v, got %v", e.errored, err)
				return
			}

			if !reflect.DeepEqual(actual, e.expected) {
				t.Errorf("Expect %+v, got %+v", e.expected, actual)
				return
			}

			httpmock.GetTotalCallCount()
			info := httpmock.GetCallCountInfo()

			if info["GET "+e.url] != 1 {
				t.Errorf("The request GET %s was expected to execute once, but it executed %d time(s)", e.url, info["GET "+e.url])
				return
			}
		})
	}
}

func TestSearchAll(t *testing.T) {
	examples := map[string]struct {
		maxRetryAttempts      uint
		activateStatusCode    int
		activateResponse      string
		adaptiveStatsCode     int
		adaptiveResponse      string
		expectedResults       []*SearchResult
		expectedActivateCount int
		expectedAdaptiveCount int
	}{
		"empty-tweets": {
			maxRetryAttempts:   uint(3),
			activateStatusCode: 200,
			activateResponse:   `{ "guest_token": "deadbeef" }`,
			adaptiveStatsCode:  200,
			adaptiveResponse:   `{}`,
			expectedResults: []*SearchResult{
				&SearchResult{&json.Adaptive{}, nil},
				nil,
			},
			expectedActivateCount: 1,
			expectedAdaptiveCount: 1,
		},
		"failed-get-guest-token": {
			maxRetryAttempts:   uint(3),
			activateStatusCode: 403,
			activateResponse:   `{ "errors": [{ "code": 200, "message": "forbidden" }] }`,
			adaptiveStatsCode:  200,
			adaptiveResponse:   `{}`,
			expectedResults: []*SearchResult{
				&SearchResult{nil, errors.New("failed to get guest token: 200: forbidden")},
				nil,
			},
			expectedActivateCount: 1,
			expectedAdaptiveCount: 0,
		},
		"retry-limit-exceeded": {
			maxRetryAttempts:   uint(3),
			activateStatusCode: 200,
			activateResponse:   `{}`,
			adaptiveStatsCode:  403,
			adaptiveResponse:   `{ "errors": [{ "code": 200, "message": "forbidden" }] }`,
			expectedResults: []*SearchResult{
				&SearchResult{nil, errors.New("retry limit exceeded: 200: forbidden")},
				nil,
			},
			expectedActivateCount: 4,
			expectedAdaptiveCount: 4,
		},
		"no-retries": {
			maxRetryAttempts:   uint(0),
			activateStatusCode: 200,
			activateResponse:   `{}`,
			adaptiveStatsCode:  403,
			adaptiveResponse:   `{ "errors": [{ "code": 200, "message": "forbidden" }] }`,
			expectedResults: []*SearchResult{
				&SearchResult{nil, errors.New("failed to search: 200: forbidden")},
				nil,
			},
			expectedActivateCount: 1,
			expectedAdaptiveCount: 1,
		},
	}

	for name, e := range examples {
		t.Run(name, func(t *testing.T) {
			c := NewClient()
			c.MaxRetryAttempts = e.maxRetryAttempts

			httpmock.ActivateNonDefault(c.Client.GetClient())
			defer httpmock.DeactivateAndReset()

			url1 := "https://api.twitter.com/1.1/guest/activate.json"
			httpmock.RegisterResponder("POST", url1, NewJsonResponse(e.activateStatusCode, e.activateResponse))

			url2 := "https://twitter.com/i/api/2/search/adaptive.json?count=40&include_quote_count=true&include_reply_count=1&q=foo&query_source=typed_query&tweet_mode=extended&tweet_search_mode=live"
			httpmock.RegisterResponder("GET", url2, NewJsonResponse(e.adaptiveStatsCode, e.adaptiveResponse))

			q := Query{Text: "foo"}
			opts := SearchOptions{Query: q}
			ch := c.SearchAll(opts)

			for _, expected := range e.expectedResults {
				actual := <-ch

				// TODO: Use reflect.DeepEqual
				if expected == nil {
					if actual != nil {
						t.Errorf("Expect nil, got %v", actual)
						return
					}
				} else {
					if expected.Error == nil {
						if actual.Error != nil {
							t.Errorf(`Expect nil, got "%v"`, actual.Error)
							return
						}
					} else {
						if actual.Error.Error() != expected.Error.Error() {
							t.Errorf(`Expect "%s", got "%s"`, expected.Error.Error(), actual.Error.Error())
							return
						}
					}

					if !reflect.DeepEqual(actual.Adaptive, expected.Adaptive) {
						t.Errorf("Expect %+v, got %+v", expected.Adaptive, actual.Adaptive)
						return
					}
				}
			}

			httpmock.GetTotalCallCount()
			info := httpmock.GetCallCountInfo()

			if info["POST "+url1] != e.expectedActivateCount {
				t.Errorf("The request POST %s was expected to execute %d time(s), but it executed %d time(s)", url1, e.expectedActivateCount, info["POST "+url1])
				return
			}

			if info["GET "+url2] != e.expectedAdaptiveCount {
				t.Errorf("The request GET %s was expected to execute %d time(s), but it executed %d time(s)", url2, e.expectedAdaptiveCount, info["GET "+url2])
				return
			}
		})
	}
}

func TestSearchAllWhenRestTweetsExist(t *testing.T) {
	c := NewClient()

	httpmock.ActivateNonDefault(c.Client.GetClient())
	defer httpmock.DeactivateAndReset()

	url1 := "https://api.twitter.com/1.1/guest/activate.json"
	httpmock.RegisterResponder("POST", url1, NewJsonResponse(200, `{ "guest_token": "1234" }`))

	url2 := "https://twitter.com/i/api/2/search/adaptive.json?count=40&include_quote_count=true&include_reply_count=1&q=foo&query_source=typed_query&tweet_mode=extended&tweet_search_mode=live"
	res2 := `{
    "globalObjects": {
      "tweets": {
        "1": {
          "id": 1,
          "full_text": "To Sherlock Holmes she is always the woman."
        }
      },
      "users": {}
    },
    "timeline": {
      "instructions": [{
        "addEntries": {
          "entries": [{
            "entryId": "sq-cursor-bottom",
            "content": {
              "operation": {
                "cursor": { "value": "scroll:deadbeef", "cursorType": "Bottom" }
              }
            }
          }]
        }
      }]
    }
  }`
	httpmock.RegisterResponder("GET", url2, NewJsonResponse(200, res2))

	url3 := "https://twitter.com/i/api/2/search/adaptive.json?count=40&cursor=scroll%3Adeadbeef&include_quote_count=true&include_reply_count=1&q=foo&query_source=typed_query&tweet_mode=extended&tweet_search_mode=live"
	httpmock.RegisterResponder("GET", url3, NewJsonResponse(200, `{}`))

	q := Query{Text: "foo"}
	opts := SearchOptions{Query: q}
	ch := c.SearchAll(opts)

	actual1 := <-ch
	if actual1.Error != nil {
		t.Errorf("Expect nil, got %+v", actual1.Error)
		return
	}

	expected1 := &json.Adaptive{
		GlobalObjects: json.GlobalObjects{
			Tweets: map[string]json.Tweet{
				"1": json.Tweet{
					Id:       1,
					FullText: "To Sherlock Holmes she is always the woman.",
				},
			},
			Users: map[string]json.User{},
		},
		Timeline: json.Timeline{
			Instructions: []json.Instruction{
				json.Instruction{
					AddEntries: json.AddEntries{
						Entries: []json.Entry{
							json.Entry{
								EntryId: "sq-cursor-bottom",
								Content: json.Content{
									Operation: json.Operation{
										Cursor: json.Cursor{Value: "scroll:deadbeef", CursorType: "Bottom"},
									},
								},
							},
						},
					},
				},
			},
		},
	}
	if !reflect.DeepEqual(actual1.Adaptive, expected1) {
		t.Errorf("Expect %+v first, got %+v", expected1, actual1.Adaptive)
		return
	}

	actual2 := <-ch
	if actual2.Error != nil {
		t.Errorf("Expect nil, got %+v", actual2.Error)
		return
	}

	expected2 := &json.Adaptive{}
	if !reflect.DeepEqual(actual2.Adaptive, expected2) {
		t.Errorf("Expect %+v second, got %+v", expected2, actual2.Adaptive)
		return
	}

	actual3 := <-ch
	if actual3 != nil {
		t.Errorf("Expect nil third, got %+v", actual3)
		return
	}

	httpmock.GetTotalCallCount()
	info := httpmock.GetCallCountInfo()

	if info["POST "+url1] != 1 {
		t.Errorf("The request POST %s was expected to execute once, but it executed %d time(s)", url1, info["POST "+url1])
		return
	}

	if info["GET "+url2] != 1 {
		t.Errorf("The request GET %s was expected to execute once, but it executed %d time(s)", url2, info["GET "+url2])
		return
	}

	if info["GET "+url3] != 1 {
		t.Errorf("The request GET %s was expected to execute once, but it executed %d time(s)", url3, info["GET "+url3])
		return
	}
}
