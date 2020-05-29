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

	"github.com/project-alvarium/go-store/internal/pkg/routes/find"

	"github.com/project-alvarium/go-sdk/pkg/annotation"
	"github.com/project-alvarium/go-sdk/pkg/identity"
	"github.com/project-alvarium/go-sdk/pkg/status"
)

const (
	findRequestorFailure = status.Unknown
	findUnmarshalFailure = status.Unknown
	findSuccess          = status.Success
)

// FindByIdentity returns annotations and status corresponding to identity.
func (i *instance) FindByIdentity(id identity.Contract) ([]*annotation.Instance, status.Value) {
	var response []byte
	var err error

	if response, err = i.requestor(find.Method, find.EscapedRoute(id), nil); err != nil {
		return nil, findRequestorFailure
	}

	var values []json.RawMessage
	if err := json.Unmarshal(response, &values); err != nil {
		return nil, findUnmarshalFailure
	}

	results := make([]*annotation.Instance, len(values))
	for valueIndex := range values {
		var value annotation.Instance
		value.SetMetadataFactory(i.mFactory)
		value.SetIdentityFactory(i.iFactory)
		if err := json.Unmarshal(values[valueIndex], &value); err != nil {
			return nil, findUnmarshalFailure
		}

		results[valueIndex] = &value
	}

	return results, findSuccess
}
