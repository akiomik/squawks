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

package api

import (
	"fmt"

	"github.com/akiomik/squawks/api/json"
)

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
