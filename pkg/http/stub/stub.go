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

package stub

// instance is a receiver that encapsulates required dependencies.
type instance struct {
	RequestMethod string
	RequestURL    string
	RequestBody   []byte
	responseBody  []byte
	err           error
}

func New(responseBody []byte, err error) *instance {
	return &instance{
		responseBody: responseBody,
		err:          err,
	}
}

// Request encapsulates an http request of method to url with body.
func (i *instance) Request(method, url string, body []byte) (responseBody []byte, err error) {
	i.RequestMethod = method
	i.RequestURL = url
	i.RequestBody = body
	return i.responseBody, i.err
}