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

package url

const name = "url"

// identity is a receiver that encapsulates required dependencies.
type identity struct {
	id string
}

// New is a factory function that returns an initialized identity.
func New(id string) *identity {
	return &identity{
		id: id,
	}
}

// Binary returns a unique key based on identity used within the SDK.
func (i *identity) Binary() []byte {
	return []byte(i.id)
}

// Printable returns a unique key based on identity used within the SDK.
func (i *identity) Printable() string {
	return i.id
}

// Kind returns the type of concrete implementation.
func (*identity) Kind() string {
	return name
}
