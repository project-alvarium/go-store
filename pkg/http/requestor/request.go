/*******************************************************************************
 * Copyright 2020 Dell Inc.
 *
 * Licensed under the Apache License, Version 2.0 (the "License"); you may not use this file except
 * in compliance with the License. You may obtain a copy of the License at
 *
 * http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software distributed under the License
 * is distributed on an "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express
 * or implied. See the License for the specific language governing permissions and limitations under
 * the License.
 *******************************************************************************/

package requestor

import (
	"bytes"
	"errors"
	"io"
	"io/ioutil"
	"net/http"
	"time"
)

// instance is a receiver that encapsulates required dependencies.
type instance struct {
	url string
}

// New is a factory function that returns instance.
func New(url string) *instance {
	return &instance{
		url: url,
	}
}

// Handler encapsulates making an http request of method to url with body.
func (i *instance) Handler(method, path string, body []byte) (responseBody []byte, err error) {
	var reader io.Reader
	var request *http.Request
	var response *http.Response

	if len(body) > 0 {
		reader = bytes.NewReader(body)
	}
	if request, err = http.NewRequest(method, i.url+path, reader); err != nil {
		return
	}

	client := &http.Client{
		Timeout: time.Second * time.Duration(30),
	}
	if response, err = client.Do(request); err != nil {
		return
	}

	defer func() {
		_ = response.Body.Close()
	}()

	if response.StatusCode != http.StatusOK {
		return nil, errors.New("response.StatusCode != http.StatusOK")
	}

	if responseBody, err = ioutil.ReadAll(response.Body); err != nil {
		return
	}

	return responseBody, nil
}
