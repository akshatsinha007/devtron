/*
 * Copyright (c) 2020-2024. Devtron Inc.
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package util

import (
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strconv"
	"time"
)

func ReadFromUrlWithRetry(url string) ([]byte, error) {
	var (
		err      error
		response *http.Response
		retries  = 3
	)

	for retries > 0 {
		response, err = http.Get(url)
		if err != nil {
			retries -= 1
			time.Sleep(1 * time.Second)
		} else {
			break
		}
	}
	if response != nil {
		defer response.Body.Close()
		statusCode := response.StatusCode
		if statusCode != http.StatusOK {
			return nil, errors.New(fmt.Sprintf("Error in getting content from url - %s. Status code : %s", url, strconv.Itoa(statusCode)))
		}
		body, err := ioutil.ReadAll(response.Body)
		if err != nil {
			return nil, err
		}
		return body, nil
	}
	return nil, err
}

func GetHost(urlStr string) (string, error) {
	u, err := url.Parse(urlStr)
	if err == nil {
		return u.Host, nil
	}
	u, err = url.Parse("//" + urlStr)
	if err != nil {
		return "", fmt.Errorf("invalid url: %w", err)
	}
	return u.Host, nil
}
