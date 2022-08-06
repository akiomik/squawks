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

package twitter

import (
	"fmt"

	"github.com/go-resty/resty/v2"

	"github.com/akiomik/get-old-tweets/config"
	"github.com/akiomik/get-old-tweets/twitter/json"
)

type Client struct {
	Client           *resty.Client
	UserAgent        string
	AuthToken        string
	MaxRetryAttempts uint
}

func NewClient() *Client {
	client := Client{}
	client.Client = resty.New()
	client.UserAgent = "get-old-tweets/" + config.Version
	client.AuthToken = "AAAAAAAAAAAAAAAAAAAAANRILgAAAAAAnNwIzUejRCOuH5E6I8xnZz4puTs%3D1Zv7ttfk8LF81IUq16cHjhLTvJu4FA33AGWWjCpTnA"
	client.MaxRetryAttempts = 3

	return &client
}

func (c *Client) Request() *resty.Request {
	client := c.Client.R().SetHeader("Accept", "application/json")

	if len(c.AuthToken) != 0 {
		client = client.SetAuthToken(c.AuthToken)
	}

	if len(c.UserAgent) != 0 {
		client = client.SetHeader("User-Agent", c.UserAgent)
	}

	return client
}

func (c *Client) GetGuestToken() (string, error) {
	res, err := c.Request().
		SetResult(json.Activate{}).
		SetError(json.ErrorResponse{}).
		Post("https://api.twitter.com/1.1/guest/activate.json")

	if err != nil {
		return "", err
	}

	if res.IsError() {
		return "", res.Error().(*json.ErrorResponse)
	}

	return res.Result().(*json.Activate).GuestToken, nil
}

func (c *Client) Search(q Query, guestToken string, cursor string) (*json.Adaptive, error) {
	params := map[string]string{
		"q":                   q.Encode(),
		"include_quote_count": "true",
		"include_reply_count": "1",
		"tweet_mode":          "extended",
		"count":               "40",
		"query_source":        "typed_query",
		"tweet_search_mode":   "live",
	}

	if len(cursor) != 0 {
		params["cursor"] = cursor
	}

	res, err := c.Request().
		SetResult(json.Adaptive{}).
		SetError(json.ErrorResponse{}).
		SetHeader("x-guest-token", guestToken).
		SetQueryParams(params).
		Get("https://twitter.com/i/api/2/search/adaptive.json")

	if err != nil {
		return nil, err
	}

	if res.IsError() {
		return nil, res.Error().(*json.ErrorResponse)
	}

	return res.Result().(*json.Adaptive), nil
}

type SearchResult struct {
	Adaptive *json.Adaptive
	Error    error
}

func (c *Client) SearchAll(q Query) <-chan *SearchResult {
	ch := make(chan *SearchResult)

	go func() {
		defer close(ch)

		cursor := ""
		guestToken := ""
		attempts := uint(0)

		for {
			if guestToken == "" {
				newGuestToken, err := c.GetGuestToken()
				if err != nil {
					ch <- &SearchResult{nil, fmt.Errorf("Failed to get guest token: %w", err)}
					break
				}

				guestToken = newGuestToken
			}

			res, err := c.Search(q, guestToken, cursor)

			if err != nil {
				// TODO: check error code
				_, ok := err.(*json.ErrorResponse)
				if ok && c.MaxRetryAttempts != 0 {
					if attempts >= c.MaxRetryAttempts {
						ch <- &SearchResult{nil, fmt.Errorf("Retry limit exceeded: %w", err)}
						break
					}

					guestToken = ""
					attempts++
					continue
				} else {
					ch <- &SearchResult{nil, fmt.Errorf("Failed to search: %w", err)}
					break
				}
			}

			ch <- &SearchResult{res, nil}
			if len(res.GlobalObjects.Tweets) == 0 {
				break
			}

			cursor, err = res.FindCursor()
			if err != nil {
				ch <- &SearchResult{nil, fmt.Errorf("Failed to find cursor: %w", err)}
				break
			}
		}
	}()

	return ch
}
