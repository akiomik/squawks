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

type SearchOptions struct {
	GuestToken string
	Cursor     string
	Query      Query
	Top        bool
}

func (c *Client) Search(opts *SearchOptions) (*json.Adaptive, error) {
	params := map[string]string{
		"q":                   opts.Query.Encode(),
		"include_quote_count": "true",
		"include_reply_count": "1",
		"tweet_mode":          "extended",
		"count":               "40",
		"query_source":        "typed_query",
	}

	if len(opts.Cursor) != 0 {
		params["cursor"] = opts.Cursor
	}

	if !opts.Top {
		params["tweet_search_mode"] = "live"
	}

	res, err := c.Request().
		SetResult(json.Adaptive{}).
		SetError(json.ErrorResponse{}).
		SetHeader("x-guest-token", opts.GuestToken).
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

func (c *Client) SearchAll(opts SearchOptions) <-chan *SearchResult {
	ch := make(chan *SearchResult)

	go func() {
		defer close(ch)

		cursor := opts.Cursor
		guestToken := opts.GuestToken
		attempts := uint(0)

		for {
			if guestToken == "" {
				newGuestToken, err := c.GetGuestToken()
				if err != nil {
					ch <- &SearchResult{nil, fmt.Errorf("failed to get guest token: %w", err)}
					break
				}

				guestToken = newGuestToken
			}

			opts.GuestToken = guestToken
			opts.Cursor = cursor
			res, err := c.Search(&opts)

			if err != nil {
				// TODO: check error code
				_, ok := err.(*json.ErrorResponse)
				if ok && c.MaxRetryAttempts != 0 {
					if attempts >= c.MaxRetryAttempts {
						ch <- &SearchResult{nil, fmt.Errorf("retry limit exceeded: %w", err)}
						break
					}

					guestToken = ""
					attempts++
					continue
				} else {
					ch <- &SearchResult{nil, fmt.Errorf("failed to search: %w", err)}
					break
				}
			}

			ch <- &SearchResult{res, nil}
			if len(res.GlobalObjects.Tweets) == 0 {
				break
			}

			cursor, err = res.FindCursor()
			if err != nil {
				ch <- &SearchResult{nil, fmt.Errorf("failed to find cursor: %w", err)}
				break
			}
		}
	}()

	return ch
}
