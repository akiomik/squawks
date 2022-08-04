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
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"

	"github.com/akiomik/get-old-tweets/config"
)

type Client struct {
	Client     *http.Client
	UserAgent  string
	AuthToken  string
	GuestToken string
}

func NewClient() *Client {
	client := Client{}
	client.Client = http.DefaultClient
	client.UserAgent = "get-old-tweets/" + config.Version
	client.AuthToken = "Bearer AAAAAAAAAAAAAAAAAAAAANRILgAAAAAAnNwIzUejRCOuH5E6I8xnZz4puTs%3D1Zv7ttfk8LF81IUq16cHjhLTvJu4FA33AGWWjCpTnA"

	return &client
}

func (c *Client) Request(m string, u string) (*http.Response, error) {
	req, err := http.NewRequest(m, u, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("User-Agent", c.UserAgent)
	req.Header.Set("Authorization", c.AuthToken)
	if len(c.GuestToken) != 0 {
		req.Header.Set("x-guest-token", c.GuestToken)
	}

	res, err := c.Client.Do(req)
	if err != nil {
		return nil, err
	}

	return res, nil
}

func (c *Client) JsonRequest(m string, u string, v any) error {
	res, err := c.Request(m, u)
	if err != nil {
		return err
	}
	defer res.Body.Close()

	blob, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return err
	}

	err = json.Unmarshal(blob, v)
	if err != nil {
		return err
	}

	return nil
}

func (c *Client) Search(q Query, cursor string) (*AdaptiveJson, error) {
	u :=
		"https://twitter.com/i/api/2/search/adaptive.json" +
			"?q=" + q.Encode() +
			"&include_quote_count=true" +
			"&include_reply_count=1" +
			"&tweet_mode=extended" +
			"&count=40" +
			"&query_source=typed_query" +
			"&tweet_search_mode=live"

	if len(cursor) != 0 {
		u += "&cursor=" + url.QueryEscape(cursor)
	}

	res := new(AdaptiveJson)
	err := c.JsonRequest("GET", u, res)
	if err != nil {
		return nil, err
	}

	if len(res.Errors) != 0 {
		message := ""
		for _, err := range res.Errors {
			message += fmt.Sprintf("[%d] %s\n", err.Code, err.Message)
		}

		return res, fmt.Errorf("Error: %s", message)
	}

	return res, nil
}

func (c *Client) SearchAll(q Query) <-chan *AdaptiveJson {
	ch := make(chan *AdaptiveJson)

	go func() {
		defer close(ch)

		cursor := ""
		for {
			res, err := c.Search(q, cursor)
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
