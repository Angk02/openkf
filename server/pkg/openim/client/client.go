// Copyright © 2023 OpenIM open source community. All rights reserved.
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

package client

import (
	"encoding/json"
	"fmt"
	"net/url"

	"github.com/go-resty/resty/v2"
)

// Client client.
type Client interface {
	GET(operationID, token string, params interface{}) (map[string]interface{}, error)
	POST(operationID, token string, params interface{}) (map[string]interface{}, error)
}

// httpClient http client.
type httpClient struct {
	url    string
	client *resty.Client
}

// NewClient new client.
func NewClient(url string) Client {
	return &httpClient{
		url:    url,
		client: resty.New(),
	}
}

// GET get.
func (c *httpClient) GET(operationID, token string, params interface{}) (map[string]interface{}, error) {
	req := c.client.R().
		SetHeader("Content-Type", "application/json").
		SetHeader("operationID", operationID)
	if token != "" {
		req.SetHeader("token", token)
	}
	switch v := params.(type) {
	case nil:
		// noop
	case map[string]string:
		req.SetQueryParams(v)
	case url.Values:
		req.SetQueryParamsFromValues(v)
	default:
		return nil, fmt.Errorf("unsupported query params type: %T", params)
	}

	resp, err := req.Get(c.url)
	if err != nil {
		return nil, err
	}

	var responseData map[string]interface{}
	if err := json.Unmarshal(resp.Body(), &responseData); err != nil {
		return nil, err
	}
	return responseData, nil
}

// POST post.
func (c *httpClient) POST(operationID, token string, params interface{}) (map[string]interface{}, error) {
	req := c.client.R().
		SetHeader("Content-Type", "application/json").
		SetHeader("operationID", operationID).
		SetBody(params)
	if token != "" {
		req.SetHeader("token", token)
	}

	resp, err := req.Post(c.url)
	if err != nil {
		return nil, err
	}

	var responseData map[string]interface{}
	err = json.Unmarshal(resp.Body(), &responseData)
	if err != nil {
		return nil, err
	}

	return responseData, nil
}
