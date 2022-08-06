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
	"github.com/akiomik/get-old-tweets/twitter/json"
)

func TestNewClient(t *testing.T) {
	c := NewClient()

	expectedUserAgent := "get-old-tweets/" + config.Version
	if c.UserAgent != expectedUserAgent {
		t.Errorf(`Expect "%s", got "%s"`, expectedUserAgent, c.UserAgent)
		return
	}

	if len(c.AuthToken) == 0 {
		t.Errorf(`Expect not "", got ""`)
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
		t.Errorf("Expect %v, got %v", expectedUserAgent, client.Header.Get("User-Agent"))
	}

	if client.Token != expectedAuthToken {
		t.Errorf("Expect %v, got %v", expectedAuthToken, client.Token)
	}
}

func TestGetGuestTokenSuccess(t *testing.T) {
	c := NewClient()

	httpmock.ActivateNonDefault(c.Client.GetClient())
	defer httpmock.DeactivateAndReset()

	url := "https://api.twitter.com/1.1/guest/activate.json"
	res := `{ "guest_token": "deadbeef" }`
	httpmock.RegisterResponder("POST", url, func(req *http.Request) (*http.Response, error) {
		res := httpmock.NewStringResponse(200, res)
		res.Header.Add("Content-Type", "application/json")
		return res, nil
	})

	actual, err := c.GetGuestToken()
	if err != nil {
		t.Errorf("Expect nil, got %v", err)
		return
	}

	expected := "deadbeef"
	if actual != expected {
		t.Errorf(`Expect "%s", got "%s"`, expected, actual)
		return
	}
}

func TestGetGuestTokenFailure(t *testing.T) {
	c := NewClient()

	httpmock.ActivateNonDefault(c.Client.GetClient())
	defer httpmock.DeactivateAndReset()

	url := "https://api.twitter.com/1.1/guest/activate.json"
	res := `{ "errors": [{ "code": 200, "message": "forbidden" }] }`
	httpmock.RegisterResponder("POST", url, func(req *http.Request) (*http.Response, error) {
		res := httpmock.NewStringResponse(403, res)
		res.Header.Add("Content-Type", "application/json")
		return res, nil
	})

	actual, err := c.GetGuestToken()
	if err == nil {
		t.Errorf("Expect error objects, got nil")
		return
	}

	if actual != "" {
		t.Errorf(`Expect "", got "%s"`, actual)
		return
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
		t.Errorf("Expect nil, got %v", err)
		return
	}

	expected := json.Adaptive{GlobalObjects: json.GlobalObjects{Tweets: map[string]json.Tweet{}, Users: map[string]json.User{}}}
	if !reflect.DeepEqual(*actual, expected) {
		t.Errorf("Expect %+v, got %+v", expected, *actual)
		return
	}

	httpmock.GetTotalCallCount()
	info := httpmock.GetCallCountInfo()

	if info["GET "+url] != 1 {
		t.Errorf("The request GET %s was expected to execute once, but it executed %d time(s)", url, info["GET "+url])
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
		t.Errorf("Expect nil, got %v", err)
		return
	}

	expected := json.Adaptive{GlobalObjects: json.GlobalObjects{Tweets: map[string]json.Tweet{}, Users: map[string]json.User{}}}
	if !reflect.DeepEqual(*actual, expected) {
		t.Errorf("Expect %+v, got %+v", expected, *actual)
		return
	}

	httpmock.GetTotalCallCount()
	info := httpmock.GetCallCountInfo()

	if info["GET "+url] != 1 {
		t.Errorf("The request GET %s was expected to execute once, but it executed %d time(s)", url, info["GET "+url])
		return
	}
}

func TestSearchWhenError(t *testing.T) {
	c := NewClient()

	httpmock.ActivateNonDefault(c.Client.GetClient())
	defer httpmock.DeactivateAndReset()

	url := "https://twitter.com/i/api/2/search/adaptive.json?count=40&include_quote_count=true&include_reply_count=1&q=foo&query_source=typed_query&tweet_mode=extended&tweet_search_mode=live"
	res := `{ "errors": [{ "code": 200, "message": "forbidden" }] }`
	httpmock.RegisterResponder("GET", url, func(req *http.Request) (*http.Response, error) {
		res := httpmock.NewStringResponse(403, res)
		res.Header.Add("Content-Type", "application/json")
		return res, nil
	})

	q := Query{Text: "foo"}
	actualAdaptive, err := c.Search(q, "", "")

	expectedError := json.ErrorResponse{Errors: []json.Error{json.Error{Code: 200, Message: "forbidden"}}}
	actualError, ok := err.(*json.ErrorResponse)
	if !ok || !reflect.DeepEqual(*actualError, expectedError) {
		t.Errorf("Expect %+v, got %+v", expectedError, *actualError)
		return
	}

	if actualAdaptive != nil {
		t.Errorf("Expect nil, got %+v", actualAdaptive)
		return
	}

	httpmock.GetTotalCallCount()
	info := httpmock.GetCallCountInfo()

	if info["GET "+url] != 1 {
		t.Errorf("The request GET %s was expected to execute once, but it executed %d time(s)", url, info["GET "+url])
		return
	}
}

func TestSearchAllWhenRestTweetDoNotExist(t *testing.T) {
	c := NewClient()

	httpmock.ActivateNonDefault(c.Client.GetClient())
	defer httpmock.DeactivateAndReset()

	url1 := "https://api.twitter.com/1.1/guest/activate.json"
	httpmock.RegisterResponder("POST", url1, func(req *http.Request) (*http.Response, error) {
		res := httpmock.NewStringResponse(200, `{ "guest_token": "1234" }`)
		res.Header.Add("Content-Type", "application/json")
		return res, nil
	})

	url2 := "https://twitter.com/i/api/2/search/adaptive.json?count=40&include_quote_count=true&include_reply_count=1&q=foo&query_source=typed_query&tweet_mode=extended&tweet_search_mode=live"
	httpmock.RegisterResponder("GET", url2, func(req *http.Request) (*http.Response, error) {
		res := httpmock.NewStringResponse(200, `{}`)
		res.Header.Add("Content-Type", "application/json")
		return res, nil
	})

	q := Query{Text: "foo"}
	ch := c.SearchAll(q)

	actual1 := <-ch
	if actual1.Error != nil {
		t.Errorf("Expect nil, got %+v", actual1.Error)
		return
	}

	expected1 := &json.Adaptive{}
	if !reflect.DeepEqual(actual1.Adaptive, expected1) {
		t.Errorf("Expect %+v first, got %+v", expected1, actual1.Adaptive)
		return
	}

	actual2 := <-ch
	if actual2 != nil {
		t.Errorf("Expect nil second, got %+v", actual2)
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
}

func TestSearchAllWhenRestTweetsExist(t *testing.T) {
	c := NewClient()

	httpmock.ActivateNonDefault(c.Client.GetClient())
	defer httpmock.DeactivateAndReset()

	url1 := "https://api.twitter.com/1.1/guest/activate.json"
	httpmock.RegisterResponder("POST", url1, func(req *http.Request) (*http.Response, error) {
		res := httpmock.NewStringResponse(200, `{ "guest_token": "1234" }`)
		res.Header.Add("Content-Type", "application/json")
		return res, nil
	})

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
	httpmock.RegisterResponder("GET", url2, func(req *http.Request) (*http.Response, error) {
		res := httpmock.NewStringResponse(200, res2)
		res.Header.Add("Content-Type", "application/json")
		return res, nil
	})

	url3 := "https://twitter.com/i/api/2/search/adaptive.json?count=40&cursor=scroll%3Adeadbeef&include_quote_count=true&include_reply_count=1&q=foo&query_source=typed_query&tweet_mode=extended&tweet_search_mode=live"
	httpmock.RegisterResponder("GET", url3, func(req *http.Request) (*http.Response, error) {
		res := httpmock.NewStringResponse(200, `{}`)
		res.Header.Add("Content-Type", "application/json")
		return res, nil
	})

	q := Query{Text: "foo"}
	ch := c.SearchAll(q)

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
