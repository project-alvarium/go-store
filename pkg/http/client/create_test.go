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
	"errors"
	"testing"

	"github.com/project-alvarium/go-store/internal/pkg/identity/url"
	"github.com/project-alvarium/go-store/internal/pkg/routes/create"
	testInternal "github.com/project-alvarium/go-store/internal/pkg/test"
	"github.com/project-alvarium/go-store/pkg/http/stub"

	"github.com/project-alvarium/go-sdk/pkg/annotation"
	metadataStub "github.com/project-alvarium/go-sdk/pkg/annotation/metadata/stub"
	"github.com/project-alvarium/go-sdk/pkg/annotation/uniqueprovider/ulid"
	identityHash "github.com/project-alvarium/go-sdk/pkg/identity/hash"
	"github.com/project-alvarium/go-sdk/pkg/status"
	"github.com/project-alvarium/go-sdk/pkg/test"

	"github.com/stretchr/testify/assert"
)

// TestInstance_Create tests Create client method.
func TestInstance_Create(t *testing.T) {
	type testCase struct {
		name string
		test func(t *testing.T)
	}

	cases := []testCase{
		{
			name: "requestor failure",
			test: func(t *testing.T) {
				requestor := stub.New(nil, errors.New(""))
				sut := newSUT(requestor.Request)

				result := sut.Create(url.New(test.FactoryRandomString()), nil)

				assert.Equal(t, createRequestorFailure, result)
			},
		},
		{
			name: "unmarshal failure",
			test: func(t *testing.T) {
				requestor := stub.New(nil, nil)
				sut := newSUT(requestor.Request)

				result := sut.Create(url.New(test.FactoryRandomString()), nil)

				assert.Equal(t, createUnmarshalFailure, result)
			},
		},
		{
			name: "success",
			test: func(t *testing.T) {
				id := identityHash.New(test.FactoryRandomByteSlice())
				idContract := url.New(id.Printable())
				m := annotation.New(ulid.New().Get(), id, nil, metadataStub.NewNullObject())
				requestor := stub.New(testInternal.Marshal(t, status.Success), nil)
				sut := newSUT(requestor.Request)

				result := sut.Create(idContract, m)

				assert.Equal(t, create.Method, requestor.RequestMethod)
				assert.Equal(t, create.EscapedRoute(idContract), requestor.RequestURL)
				assert.Equal(t, testInternal.Marshal(t, m), requestor.RequestBody)
				assert.Equal(t, createSuccess, result)
			},
		},
	}

	for i := range cases {
		t.Run(cases[i].name, cases[i].test)
	}
}
