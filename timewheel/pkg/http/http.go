/*
 * MIT License
 *
 * Copyright (c) 2023 Runze Wu
 *
 * Permission is hereby granted, free of charge, to any person obtaining a copy
 * of this software and associated documentation files (the "Software"), to deal
 * in the Software without restriction, including without limitation the rights
 * to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
 * copies of the Software, and to permit persons to whom the Software is
 * furnished to do so, subject to the following conditions:
 *
 * The above copyright notice and this permission notice shall be included in all * copies or substantial portions of the Software.
 *
 * THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
 * IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
 * FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
 * AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
 * LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
 * OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
 * SOFTWARE.
 *
 */

package http

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
)

type Client struct {
	core *http.Client
}

func NewClient() *Client {
	return &Client{
		core: http.DefaultClient,
	}
}

func (c *Client) JSONGet(ctx context.Context, url string, header, params map[string]string, resp interface{}) error {
	return c.JSONDo(ctx, http.MethodGet, getCompleteURL(url, params), header, nil, resp)
}

func (c *Client) JSONPost(ctx context.Context, url string, header map[string]string, req, resp interface{}) error {
	return c.JSONDo(ctx, http.MethodPost, url, header, req, resp)
}

func (c *Client) JSONDo(ctx context.Context, method, url string, header map[string]string, req, resp interface{}) error {
	var reqReader io.Reader
	if req != nil {
		body, _ := json.Marshal(req)
		reqReader = bytes.NewReader(body)
	}
	request, err := http.NewRequestWithContext(ctx, method, url, reqReader)
	if err != nil {
		return err
	}
	if request.Header == nil {
		request.Header = make(http.Header)
	}
	for k, v := range header {
		request.Header.Add(k, v)
	}
	request.Header.Add("Content-Type", "application/json")
	response, err := c.core.Do(request)
	if err != nil {
		return err
	}

	defer response.Body.Close()
	if response.StatusCode != http.StatusOK {
		return fmt.Errorf("invalid status: %d", response.StatusCode)
	}
	if resp == nil {
		return nil
	}
	responseBody, err := io.ReadAll(response.Body)
	if err != nil {
		return nil
	}
	return json.Unmarshal(responseBody, resp)
}

func getCompleteURL(original string, params map[string]string) string {
	if len(params) == 0 {
		return original
	}
	values := url.Values{}
	for k, v := range params {
		values.Add(k, v)
	}
	query, _ := url.QueryUnescape(values.Encode())
	return fmt.Sprintf("%s?%s", original, query)
}
