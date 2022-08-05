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
	"os"

	"github.com/go-resty/resty/v2"

	"github.com/akiomik/get-old-tweets/config"
)

type Client struct {
	Client    *resty.Client
	UserAgent string
	AuthToken string
}

func NewClient() *Client {
	client := Client{}
	client.Client = resty.New()
	client.UserAgent = "get-old-tweets/" + config.Version
	client.AuthToken = "AAAAAAAAAAAAAAAAAAAAANRILgAAAAAAnNwIzUejRCOuH5E6I8xnZz4puTs%3D1Zv7ttfk8LF81IUq16cHjhLTvJu4FA33AGWWjCpTnA"

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

func (c *Client) Search(q Query, guestToken string, cursor string) (*Adaptive, error) {
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
		SetResult(Adaptive{}).
		SetError(Adaptive{}).
		SetHeader("x-guest-token", guestToken).
		SetQueryParams(params).
		Get("https://twitter.com/i/api/2/search/adaptive.json")

	if err != nil {
		return nil, err
	}

	return res.Result().(*Adaptive), nil
}

func (c *Client) SearchAll(q Query) <-chan *Adaptive {
	ch := make(chan *Adaptive)

	go func() {
		defer close(ch)

		cursor := ""
		guestToken := ""
		for {
			res, err := c.Search(q, guestToken, cursor)
			if err != nil {
				fmt.Fprintf(os.Stderr, "Error: %v\n", err)
				panic(err)
			}

			ch <- res
			if len(res.GlobalObjects.Tweets) == 0 {
				break
			}

			cursor, err = res.FindCursor()
			if err != nil {
				fmt.Fprintf(os.Stderr, "Error: %v\n", err)
				panic(err)
			}
		}
	}()

	return ch
}
