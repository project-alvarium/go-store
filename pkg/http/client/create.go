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

package client

import (
	"encoding/json"

	"github.com/michaelestrin/go-store/internal/pkg/routes/create"

	"github.com/project-alvarium/go-sdk/pkg/annotation"
	"github.com/project-alvarium/go-sdk/pkg/identity"
	"github.com/project-alvarium/go-sdk/pkg/status"
)

const (
	createMarshalFailure   = status.Unknown
	createRequestorFailure = status.Unknown
	createUnmarshalFailure = status.Unknown
	createSuccess          = status.Success
)

// Create stores annotation corresponding to a new identity and returns status.
func (i *instance) Create(id identity.Contract, m *annotation.Instance) (result status.Value) {
	var body, response []byte
	var err error

	if body, err = json.Marshal(m); err != nil {
		return createMarshalFailure
	}

	if response, err = i.requestor(create.Method, create.EscapedRoute(id), body); err != nil {
		return createRequestorFailure
	}

	if err := json.Unmarshal(response, &result); err != nil {
		return createUnmarshalFailure
	}

	return
}
