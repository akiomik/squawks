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
	"net/http"
	"reflect"
	"testing"

	"github.com/jarcoal/httpmock"

	"github.com/akiomik/get-old-tweets/config"
)

func TestNewClient(t *testing.T) {
	c := NewClient()

	expectedUserAgent := "get-old-tweets/" + config.Version
	if c.UserAgent != expectedUserAgent {
		t.Errorf(`Expect "%s", got "%s"`, expectedUserAgent, c.UserAgent)
		return
	}

	if len(c.AuthToken) == 0 {
		t.Errorf(`Expect not "", got "%s"`, c.AuthToken)
		return
	}
}

func TestRequest(t *testing.T) {
	c := NewClient()

	expectedUserAgent := "custom-user-agent"
	expectedAuthToken := "my-auth-token"

	c.UserAgent = expectedUserAgent
	c.AuthToken = expectedAuthToken
	client := c.Request()

	if client.Header.Get("User-Agent") != expectedUserAgent {
		t.Errorf(`Expect "%v", got "%v"`, expectedUserAgent, client.Header.Get("User-Agent"))
	}

	if client.Token != expectedAuthToken {
		t.Errorf(`Expect "%v", got "%v"`, expectedAuthToken, client.Token)
	}
}

func TestSearchWhenWithoutCursor(t *testing.T) {
	c := NewClient()

	httpmock.ActivateNonDefault(c.Client.GetClient())
	defer httpmock.DeactivateAndReset()

	url := "https://twitter.com/i/api/2/search/adaptive.json?count=40&include_quote_count=true&include_reply_count=1&q=foo&query_source=typed_query&tweet_mode=extended&tweet_search_mode=live"
	res := `{ "globalObjects": { "tweets": {}, "users": {} } }`
	httpmock.RegisterResponder("GET", url, func(req *http.Request) (*http.Response, error) {
		res := httpmock.NewStringResponse(200, res)
		res.Header.Add("Content-Type", "application/json")
		return res, nil
	})

	q := Query{Text: "foo"}
	actual, err := c.Search(q, "", "")
	if err != nil {
		t.Errorf("Expect Client#Search not to return error objects, but got %v", err)
		return
	}

	expected := Adaptive{GlobalObjects: GlobalObjects{Tweets: map[string]Tweet{}, Users: map[string]User{}}}
	if !reflect.DeepEqual(*actual, expected) {
		t.Errorf("Expect Client#Search to return %+v, but got %+v", expected, *actual)
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
	c := NewClient()

	httpmock.ActivateNonDefault(c.Client.GetClient())
	defer httpmock.DeactivateAndReset()

	url := "https://twitter.com/i/api/2/search/adaptive.json?count=40&cursor=scroll%3Adeadbeef&include_quote_count=true&include_reply_count=1&q=foo&query_source=typed_query&tweet_mode=extended&tweet_search_mode=live"
	res := `{ "globalObjects": { "tweets": {}, "users": {} } }`
	httpmock.RegisterResponder("GET", url, func(req *http.Request) (*http.Response, error) {
		res := httpmock.NewStringResponse(200, res)
		res.Header.Add("Content-Type", "application/json")
		return res, nil
	})

	q := Query{Text: "foo"}
	actual, err := c.Search(q, "", "scroll:deadbeef")
	if err != nil {
		t.Errorf("Expect Client#Search not to return error objects, but got %v", err)
		return
	}

	expected := Adaptive{GlobalObjects: GlobalObjects{Tweets: map[string]Tweet{}, Users: map[string]User{}}}
	if !reflect.DeepEqual(*actual, expected) {
		t.Errorf("Expect Client#Search to return %+v, but got %+v", expected, *actual)
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
	c := NewClient()

	httpmock.ActivateNonDefault(c.Client.GetClient())
	defer httpmock.DeactivateAndReset()

	url := "https://twitter.com/i/api/2/search/adaptive.json?count=40&include_quote_count=true&include_reply_count=1&q=foo&query_source=typed_query&tweet_mode=extended&tweet_search_mode=live"
	res := `{}`
	httpmock.RegisterResponder("GET", url, httpmock.NewStringResponder(200, res))

	q := Query{Text: "foo"}
	ch := c.SearchAll(q)

	expected1 := Adaptive{}
	actual1 := <-ch
	if !reflect.DeepEqual(*actual1, expected1) {
		t.Errorf("Expect %+v first, got %+v", expected1, *actual1)
		return
	}

	actual2 := <-ch
	if actual2 != nil {
		t.Errorf("Expect nil second, got %+v", actual2)
		return
	}

	httpmock.GetTotalCallCount()
	info := httpmock.GetCallCountInfo()

	if info["GET "+url] != 1 {
		t.Errorf("The request GET %s was expected to execute once, but it executed %d time(s)", url, info["GET "+url])
		return
	}
}

func TestSearchAllWhenRestTweetsExist(t *testing.T) {
	c := NewClient()

	httpmock.ActivateNonDefault(c.Client.GetClient())
	defer httpmock.DeactivateAndReset()

	url1 := "https://twitter.com/i/api/2/search/adaptive.json?count=40&include_quote_count=true&include_reply_count=1&q=foo&query_source=typed_query&tweet_mode=extended&tweet_search_mode=live"
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
	httpmock.RegisterResponder("GET", url1, func(req *http.Request) (*http.Response, error) {
		res := httpmock.NewStringResponse(200, res1)
		res.Header.Add("Content-Type", "application/json")
		return res, nil
	})

	url2 := "https://twitter.com/i/api/2/search/adaptive.json?count=40&cursor=scroll%3Adeadbeef&include_quote_count=true&include_reply_count=1&q=foo&query_source=typed_query&tweet_mode=extended&tweet_search_mode=live"
	res2 := `{}`
	httpmock.RegisterResponder("GET", url2, func(req *http.Request) (*http.Response, error) {
		res := httpmock.NewStringResponse(200, res2)
		res.Header.Add("Content-Type", "application/json")
		return res, nil
	})

	q := Query{Text: "foo"}
	ch := c.SearchAll(q)

	expected1 := Adaptive{
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
	actual1 := <-ch
	if !reflect.DeepEqual(*actual1, expected1) {
		t.Errorf("Expect %+v first, got %+v", expected1, *actual1)
		return
	}

	expected2 := Adaptive{}
	actual2 := <-ch
	if !reflect.DeepEqual(*actual2, expected2) {
		t.Errorf("Expect %+v second, got %+v", expected2, *actual2)
		return
	}

	actual3 := <-ch
	if actual3 != nil {
		t.Errorf("Expect nil third, got %+v", actual3)
		return
	}

	httpmock.GetTotalCallCount()
	info := httpmock.GetCallCountInfo()

	if info["GET "+url1] != 1 {
		t.Errorf("The request GET %s was expected to execute once, but it executed %d time(s)", url1, info["GET "+url1])
		return
	}

	if info["GET "+url2] != 1 {
		t.Errorf("The request GET %s was expected to execute once, but it executed %d time(s)", url2, info["GET "+url2])
		return
	}
}
