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

	"github.com/michaelestrin/go-store/internal/pkg/identity/url"
	"github.com/michaelestrin/go-store/internal/pkg/routes/find"
	testInternal "github.com/michaelestrin/go-store/internal/pkg/test"
	"github.com/michaelestrin/go-store/pkg/http/stub"

	"github.com/project-alvarium/go-sdk/pkg/annotation"
	metadataFactory "github.com/project-alvarium/go-sdk/pkg/annotation/metadata/factory"
	metadataStub "github.com/project-alvarium/go-sdk/pkg/annotation/metadata/stub"
	metadataStubFactory "github.com/project-alvarium/go-sdk/pkg/annotation/metadata/stub/factory"
	"github.com/project-alvarium/go-sdk/pkg/annotation/uniqueprovider/ulid"
	"github.com/project-alvarium/go-sdk/pkg/identity/hash"
	"github.com/project-alvarium/go-sdk/pkg/test"

	"github.com/stretchr/testify/assert"
)

// TestInstance_FindByIdentity tests FindByIdentity client method.
func TestInstance_FindByIdentity(t *testing.T) {
	type testCase struct {
		name string
		test func(t *testing.T)
	}

	cases := []testCase{
		{
			name: "Requestor failure",
			test: func(t *testing.T) {
				requestor := stub.New(nil, errors.New(""))
				sut := newSUT(requestor.Request)

				value, result := sut.FindByIdentity(url.New(test.FactoryRandomString()))

				assert.Nil(t, value)
				assert.Equal(t, findRequestorFailure, result)
			},
		},
		{
			name: "Unmarshal failure",
			test: func(t *testing.T) {
				requestor := stub.New(nil, nil)
				sut := newSUT(requestor.Request)

				value, result := sut.FindByIdentity(url.New(test.FactoryRandomString()))

				assert.Nil(t, value)
				assert.Equal(t, findUnmarshalFailure, result)
			},
		},
		{
			name: "Success",
			test: func(t *testing.T) {
				s := metadataStub.NewNullObject()
				response := []interface{}{
					annotation.New(ulid.New().Get(), hash.New(test.FactoryRandomByteSlice()), nil, s),
				}
				requestor := stub.New(testInternal.Marshal(t, response), nil)
				sut := newSUTWithFactories(
					requestor.Request,
					[]metadataFactory.Contract{
						metadataStubFactory.New(s),
					},
				)
				id := test.FactoryRandomString()
				idContract := url.New(id)

				value, result := sut.FindByIdentity(idContract)

				assert.Equal(t, find.Method, requestor.RequestMethod)
				assert.Equal(t, find.EscapedRoute(idContract), requestor.RequestURL)
				assert.Nil(t, requestor.RequestBody)
				assert.Equal(t, testInternal.Marshal(t, response), testInternal.Marshal(t, value))
				assert.Equal(t, findSuccess, result)
			},
		},
	}

	for i := range cases {
		t.Run(cases[i].name, cases[i].test)
	}
}
