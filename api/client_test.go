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
	"net/http"
	"testing"

	"github.com/jarcoal/httpmock"
	"github.com/stretchr/testify/assert"

	"github.com/akiomik/squawks/config"
)

func NewJsonResponse(code int, body string) func(req *http.Request) (*http.Response, error) {
	return func(req *http.Request) (*http.Response, error) {
		res := httpmock.NewStringResponse(code, body)
		res.Header.Add("Content-Type", "application/json")
		return res, nil
	}
}

func TestNewClient(t *testing.T) {
	c := NewClient()

	assert.Equal(t, "squawks/"+config.Version, c.UserAgent)
	assert.NotEmpty(t, c.AuthToken)
}

func TestRequest(t *testing.T) {
	c := NewClient()

	expectedUserAgent := "custom-user-agent"
	expectedAuthToken := "my-auth-token"

	c.UserAgent = expectedUserAgent
	c.AuthToken = expectedAuthToken
	client := c.Request()

	assert.Equal(t, expectedUserAgent, client.Header.Get("User-Agent"))
	assert.Equal(t, expectedAuthToken, client.Token)
}

func TestGetGuestToken(t *testing.T) {
	examples := map[string]struct {
		statusCode  int
		response    string
		expected    string
		expectError bool
	}{
		"success": {
			statusCode:  200,
			response:    `{ "guest_token": "deadbeef" }`,
			expected:    "deadbeef",
			expectError: false,
		},
		"failure": {
			statusCode:  403,
			response:    `{ "errors": [{ "code": 200, "message": "forbidden" }] }`,
			expected:    "",
			expectError: true,
		},
	}

	for name, e := range examples {
		t.Run(name, func(t *testing.T) {
			c := NewClient()

			httpmock.ActivateNonDefault(c.Client.GetClient())
			defer httpmock.DeactivateAndReset()

			url := "https://api.twitter.com/1.1/guest/activate.json"
			httpmock.RegisterResponder("POST", url, NewJsonResponse(e.statusCode, e.response))

			actual, err := c.GetGuestToken()
			if e.expectError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}

			assert.Equal(t, e.expected, actual)

			httpmock.GetTotalCallCount()
			info := httpmock.GetCallCountInfo()

			assert.Equal(t, 1, info["POST "+url])
		})
	}
}
