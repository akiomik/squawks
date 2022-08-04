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

package twitter

import (
	"reflect"
	"testing"

	"github.com/jarcoal/httpmock"

	"github.com/akiomik/get-old-tweets/config"
)

func TestNewClient(t *testing.T) {
	c := NewClient()

	expectedUserAgent := "get-old-tweets/" + config.Version
	if c.UserAgent != expectedUserAgent {
		t.Errorf(`Expect "%s", but got "%s"`, expectedUserAgent, c.UserAgent)
		return
	}

	expectedGuestToken := ""
	if c.GuestToken != expectedGuestToken {
		t.Errorf(`Expect "%s", but got "%s"`, expectedGuestToken, c.GuestToken)
		return
	}
}

type TestEntityJson struct {
	Id uint `json:"id"`
}

func TestJsonRequestWhenValidJson(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	url := "https://example.com/entities/1"
	httpmock.RegisterResponder("GET", url, httpmock.NewStringResponder(200, `{ "id": 1 }`))

	c := NewClient()
	actual := new(TestEntityJson)
	err := c.JsonRequest("GET", url, actual)
	if err != nil {
		t.Errorf("Expect error, got nil")
		return
	}

	expected := TestEntityJson{Id: 1}
	if !reflect.DeepEqual(*actual, expected) {
		t.Errorf("Expect %+v, but got %+v", expected, *actual)
	}

	httpmock.GetTotalCallCount()
	info := httpmock.GetCallCountInfo()

	count := info["GET "+url]
	if count != 1 {
		t.Errorf("The request GET %s was expected to execute once, but it executed %d time(s)", url, count)
		return
	}
}

func TestJsonRequestWhenUnexpectedJson(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	url := "https://example.com/entities/1"
	httpmock.RegisterResponder("GET", url, httpmock.NewStringResponder(200, `[]`))

	c := NewClient()
	res := new(TestEntityJson)
	err := c.JsonRequest("GET", url, res)
	if err == nil {
		t.Errorf("Expect error, got nil")
		return
	}

	httpmock.GetTotalCallCount()
	info := httpmock.GetCallCountInfo()

	count := info["GET "+url]
	if count != 1 {
		t.Errorf("The request GET %s was expected to execute once, but it called %d time(s)", url, count)
		return
	}
}

func TestSearchWhenUnexpectedJson(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	url := "https://twitter.com/i/api/2/search/adaptive.json?q=foo&include_quote_count=true&include_reply_count=1&tweet_mode=extended&count=40&query_source=typed_query&tweet_search_mode=live"
	res := `{ "globalObjects": [] }`
	httpmock.RegisterResponder("GET", url, httpmock.NewStringResponder(200, res))

	c := NewClient()
	q := Query{Text: "foo"}
	actualJson, err := c.Search(q, "")
	if err == nil {
		t.Errorf("Expect Client#Search to return error objects, but got nil")
		return
	}

	if actualJson != nil {
		t.Errorf("Expect Client#Search to nil, but got %+v", actualJson)
	}

	httpmock.GetTotalCallCount()
	info := httpmock.GetCallCountInfo()

	if info["GET "+url] != 1 {
		t.Errorf("Expect Client#Search to call %s once, but it called %d time(s)", url, info["GET "+url])
		return
	}

	if len(info) != 1 {
		t.Errorf("Expect Client#Search to call %s only, but it called other urls too: %v", url, info)
		return
	}
}

func TestSearchWhenWithoutCursor(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	url := "https://twitter.com/i/api/2/search/adaptive.json?q=foo&include_quote_count=true&include_reply_count=1&tweet_mode=extended&count=40&query_source=typed_query&tweet_search_mode=live"
	res := `{ "globalObjects": { "tweets": {}, "users": {} } }`
	httpmock.RegisterResponder("GET", url, httpmock.NewStringResponder(200, res))

	c := NewClient()
	q := Query{Text: "foo"}
	actualJson, err := c.Search(q, "")
	if err != nil {
		t.Errorf("Expect Client#Search not to return error objects, but got %v", err)
		return
	}

	expectedJson := AdaptiveJson{GlobalObjects: GlobalObjects{Tweets: map[string]Tweet{}, Users: map[string]User{}}}
	if !reflect.DeepEqual(*actualJson, expectedJson) {
		t.Errorf("Expect Client#Search to return %+v, but got %+v", expectedJson, actualJson)
		return
	}

	httpmock.GetTotalCallCount()
	info := httpmock.GetCallCountInfo()

	if info["GET "+url] != 1 {
		t.Errorf("Expect Client#Search to call %s once, but it called %d time(s)", url, info["GET "+url])
		return
	}

	if len(info) != 1 {
		t.Errorf("Expect Client#Search to call %s only, but it called other urls too: %v", url, info)
		return
	}
}

func TestSearchWhenWithCursor(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	url := "https://twitter.com/i/api/2/search/adaptive.json?q=foo&include_quote_count=true&include_reply_count=1&tweet_mode=extended&count=40&query_source=typed_query&tweet_search_mode=live&cursor=scroll%3Adeadbeef"
	res := `{ "globalObjects": { "tweets": {}, "users": {} } }`
	httpmock.RegisterResponder("GET", url, httpmock.NewStringResponder(200, res))

	c := NewClient()
	q := Query{Text: "foo"}
	actualJson, err := c.Search(q, "scroll:deadbeef")
	if err != nil {
		t.Errorf("Expect Client#Search not to return error objects, but got %v", err)
		return
	}

	expectedJson := AdaptiveJson{GlobalObjects: GlobalObjects{Tweets: map[string]Tweet{}, Users: map[string]User{}}}
	if !reflect.DeepEqual(*actualJson, expectedJson) {
		t.Errorf("Expect Client#Search to return %+v, but got %+v", expectedJson, actualJson)
		return
	}

	httpmock.GetTotalCallCount()
	info := httpmock.GetCallCountInfo()

	if info["GET "+url] != 1 {
		t.Errorf("Expect Client#Search to call %s once, but it called %d time(s)", url, info["GET "+url])
		return
	}

	if len(info) != 1 {
		t.Errorf("Expect Client#Search to call %s only, but it called other urls too: %v", url, info)
		return
	}
}

func TestSearchWhenErrorsExist(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	url := "https://twitter.com/i/api/2/search/adaptive.json?q=foo&include_quote_count=true&include_reply_count=1&tweet_mode=extended&count=40&query_source=typed_query&tweet_search_mode=live"
	res := `{ "errors": [{ "code": 200, "message": "Forbidden" }] }`
	httpmock.RegisterResponder("GET", url, httpmock.NewStringResponder(200, res))

	c := NewClient()
	q := Query{Text: "foo"}
	actualJson, err := c.Search(q, "")
	if err == nil {
		t.Errorf("Expect Client#Search to return error objects, but got nil")
		return
	}

	expectedJson := AdaptiveJson{Errors: []Error{Error{Code: 200, Message: "Forbidden"}}}
	if !reflect.DeepEqual(*actualJson, expectedJson) {
		t.Errorf("Expect Client#Search to return %+v, but got %+v", expectedJson, actualJson)
		return
	}

	httpmock.GetTotalCallCount()
	info := httpmock.GetCallCountInfo()

	if info["GET "+url] != 1 {
		t.Errorf("Expect Client#Search to call %s once, but it called %d time(s)", url, info["GET "+url])
		return
	}

	if len(info) != 1 {
		t.Errorf("Expect Client#Search to call %s only, but it called other urls too: %v", url, info)
		return
	}
}

func TestSearchAllWhenRestTweetDoNotExist(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	url := "https://twitter.com/i/api/2/search/adaptive.json?q=foo&include_quote_count=true&include_reply_count=1&tweet_mode=extended&count=40&query_source=typed_query&tweet_search_mode=live"
	res := `{}`
	httpmock.RegisterResponder("GET", url, httpmock.NewStringResponder(200, res))

	c := NewClient()
	q := Query{Text: "foo"}
	ch := c.SearchAll(q)

	expectedJson1 := AdaptiveJson{}
	actualJson1 := <-ch
	if !reflect.DeepEqual(*actualJson1, expectedJson1) {
		t.Errorf("Expect Client#SearchAll to send %+v first, but got %+v", expectedJson1, *actualJson1)
		return
	}

	actualJson2 := <-ch
	if actualJson2 != nil {
		t.Errorf("Expect Client#SearchAll to send nil second, but got %+v", actualJson2)
		return
	}

	httpmock.GetTotalCallCount()
	info := httpmock.GetCallCountInfo()

	if info["GET "+url] != 1 {
		t.Errorf("Expect Client#SearchAll to call %s once, but it called %d time(s)", url, info["GET "+url])
		return
	}

	if len(info) != 1 {
		t.Errorf("Expect Client#SearchAll to call %s only, but it called other urls too: %v", url, info)
		return
	}
}

func TestSearchAllWhenRestTweetsExist(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	url1 := "https://twitter.com/i/api/2/search/adaptive.json?q=foo&include_quote_count=true&include_reply_count=1&tweet_mode=extended&count=40&query_source=typed_query&tweet_search_mode=live"
	res1 := `{
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
	httpmock.RegisterResponder("GET", url1, httpmock.NewStringResponder(200, res1))

	url2 := "https://twitter.com/i/api/2/search/adaptive.json?q=foo&include_quote_count=true&include_reply_count=1&tweet_mode=extended&count=40&query_source=typed_query&tweet_search_mode=live&cursor=scroll%3Adeadbeef"
	res2 := `{}`
	httpmock.RegisterResponder("GET", url2, httpmock.NewStringResponder(200, res2))

	c := NewClient()
	q := Query{Text: "foo"}
	ch := c.SearchAll(q)

	expectedJson1 := AdaptiveJson{
		GlobalObjects: GlobalObjects{
			Tweets: map[string]Tweet{
				"1": Tweet{
					Id:       1,
					FullText: "To Sherlock Holmes she is always the woman.",
				},
			},
			Users: map[string]User{},
		},
		Timeline: Timeline{
			Instructions: []Instruction{
				Instruction{
					AddEntries: AddEntries{
						Entries: []Entry{
							Entry{
								EntryId: "sq-cursor-bottom",
								Content: Content{
									Operation: Operation{
										Cursor: Cursor{Value: "scroll:deadbeef", CursorType: "Bottom"},
									},
								},
							},
						},
					},
				},
			},
		},
	}
	actualJson1 := <-ch
	if !reflect.DeepEqual(*actualJson1, expectedJson1) {
		t.Errorf("Expect Client#SearchAll to send %+v first, but got %+v", expectedJson1, *actualJson1)
		return
	}

	expectedJson2 := AdaptiveJson{}
	actualJson2 := <-ch
	if !reflect.DeepEqual(*actualJson2, expectedJson2) {
		t.Errorf("Expect Client#SearchAll to send %+v second, but got %+v", expectedJson2, *actualJson2)
		return
	}

	actualJson3 := <-ch
	if actualJson3 != nil {
		t.Errorf("Expect Client#SearchAll to send nil third, but got %+v", actualJson3)
		return
	}

	httpmock.GetTotalCallCount()
	info := httpmock.GetCallCountInfo()

	if info["GET "+url1] != 1 {
		t.Errorf("Expect Client#SearchAll to call %s once, but it called %d time(s)", url1, info["GET "+url1])
		return
	}

	if info["GET "+url2] != 1 {
		t.Errorf("Expect Client#SearchAll to call %s once, but it called %d time(s)", url2, info["GET "+url2])
		return
	}

	if len(info) != 2 {
		t.Errorf("Expect Client#SearchAll to call just 2 urls, but it called %d urls", len(info))
		return
	}
}
