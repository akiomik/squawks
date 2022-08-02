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

	"github.com/akiomik/get-old-tweets/config"
)

type Client struct {
	Client    *http.Client
	UserAgent string
}

func NewClient() *Client {
	client := Client{}
	client.Client = http.DefaultClient
	client.UserAgent = "get-old-tweets/v" + config.Version

	return &client
}

func (c *Client) get(url *url.URL) (*http.Response, error) {
	req, err := http.NewRequest("GET", url.String(), nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("User-Agent", c.UserAgent)
	req.Header.Set("Authorization", "Bearer AAAAAAAAAAAAAAAAAAAAANRILgAAAAAAnNwIzUejRCOuH5E6I8xnZz4puTs%3D1Zv7ttfk8LF81IUq16cHjhLTvJu4FA33AGWWjCpTnA")
	req.Header.Set("x-guest-token", "1554370317552091136")

	res, err := c.Client.Do(req)
	if err != nil {
		return nil, err
	}

	return res, nil
}

func (c *Client) Search(q Query, cursor string) (*AdaptiveJson, error) {
	urlString :=
		"https://twitter.com/i/api/2/search/adaptive.json" +
			"?q=" + q.Encode() +
			"&include_quote_count=true" +
			"&include_reply_count=1" +
			"&tweet_mode=extended" +
			"&count=40" +
			"&query_source=typed_query" +
			"&tweet_search_mode=live"

	if len(cursor) != 0 {
		urlString += "&cursor=" + url.QueryEscape(cursor)
	}

	url, _ := url.Parse(urlString)
	res, err := c.get(url)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	blob, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	adaptiveJson := new(AdaptiveJson)
	err = json.Unmarshal(blob, &adaptiveJson)
	if err != nil {
		return nil, err
	}

	if len(adaptiveJson.Errors) != 0 {
		message := ""
		for _, err := range adaptiveJson.Errors {
			message += fmt.Sprintf("[%d] %s\n", err.Code, err.Message)
		}

		return nil, fmt.Errorf("Error: %s", message)
	}

	return adaptiveJson, nil
}

func (c *Client) SearchAll(q Query) <-chan *AdaptiveJson {
	ch := make(chan *AdaptiveJson)

	go func() {
		defer close(ch)

		cursor := ""
		for {
			adaptiveJson, err := c.Search(q, cursor)
			if err != nil {
				panic(err)
			}

			ch <- adaptiveJson
			if len(adaptiveJson.GlobalObjects.Tweets) == 0 {
				break
			}

			cursor, err = adaptiveJson.FindCursor()
			if err != nil {
				panic(err)
			}
		}
	}()

	return ch
}
