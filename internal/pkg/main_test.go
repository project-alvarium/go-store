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

package pkg

import (
	"testing"

	"github.com/project-alvarium/go-store/internal/pkg/identity/url"
	"github.com/project-alvarium/go-store/internal/pkg/routable"
	"github.com/project-alvarium/go-store/internal/pkg/routes/append"
	"github.com/project-alvarium/go-store/internal/pkg/routes/create"
	"github.com/project-alvarium/go-store/internal/pkg/routes/find"
	testInternal "github.com/project-alvarium/go-store/internal/pkg/test"

	"github.com/project-alvarium/go-sdk/pkg/annotation"
	metadataFactory "github.com/project-alvarium/go-sdk/pkg/annotation/metadata/factory"
	metadataStub "github.com/project-alvarium/go-sdk/pkg/annotation/metadata/stub"
	"github.com/project-alvarium/go-sdk/pkg/annotation/store"
	"github.com/project-alvarium/go-sdk/pkg/annotation/store/memory"
	"github.com/project-alvarium/go-sdk/pkg/annotation/uniqueprovider/ulid"
	assessMetadataFactory "github.com/project-alvarium/go-sdk/pkg/annotator/assess/metadata/factory"
	pkiMetadataFactory "github.com/project-alvarium/go-sdk/pkg/annotator/pki/metadata/factory"
	publishMetadataFactory "github.com/project-alvarium/go-sdk/pkg/annotator/publish/metadata/factory"
	identityFactory "github.com/project-alvarium/go-sdk/pkg/identity/factory"
	"github.com/project-alvarium/go-sdk/pkg/identity/hash"
	"github.com/project-alvarium/go-sdk/pkg/status"
	"github.com/project-alvarium/go-sdk/pkg/test"

	"github.com/gorilla/mux"
	"github.com/stretchr/testify/assert"
)

// TestAcceptance implements multi-route acceptance tests.
func TestAcceptance(t *testing.T) {
	type testCase struct {
		name string
		test func(t *testing.T, muxRouter *mux.Router, store store.Contract)
	}

	cases := []testCase{
		{
			name: "Add new",
			test: func(t *testing.T, muxRouter *mux.Router, store store.Contract) {
				id := hash.New(test.FactoryRandomByteSlice())
				value := annotation.New(ulid.New().Get(), id, nil, metadataStub.NewNullObject())

				response := testInternal.SendRequestWithBody(
					t,
					muxRouter,
					create.Method,
					create.EscapedRoute(url.New(id.Printable())),
					testInternal.Marshal(t, value),
				)

				assert.Equal(t, create.CodeSuccess, response.Code)
				assert.Equal(t, testInternal.Marshal(t, status.Success), response.Body.Bytes())
			},
		},
		{
			name: "Add existing",
			test: func(t *testing.T, muxRouter *mux.Router, store store.Contract) {
				id := hash.New(test.FactoryRandomByteSlice())
				idContract := url.New(id.Printable())
				value := annotation.New(ulid.New().Get(), id, nil, metadataStub.NewNullObject())
				assert.Equal(t, status.Success, store.Create(idContract, value))

				response := testInternal.SendRequestWithBody(
					t,
					muxRouter,
					create.Method,
					create.EscapedRoute(idContract),
					testInternal.Marshal(t, value),
				)

				assert.Equal(t, create.CodeSuccess, response.Code)
				assert.Equal(t, testInternal.Marshal(t, status.Exists), response.Body.Bytes())
			},
		},
		{
			name: "Append new",
			test: func(t *testing.T, muxRouter *mux.Router, store store.Contract) {
				id := hash.New(test.FactoryRandomByteSlice())
				value := annotation.New(ulid.New().Get(), id, nil, metadataStub.NewNullObject())

				response := testInternal.SendRequestWithBody(
					t,
					muxRouter,
					append.Method,
					append.EscapedRoute(url.New(id.Printable())),
					testInternal.Marshal(t, value),
				)

				assert.Equal(t, append.CodeSuccess, response.Code)
				assert.Equal(t, testInternal.Marshal(t, status.NotFound), response.Body.Bytes())
			},
		},
		{
			name: "Append existing",
			test: func(t *testing.T, muxRouter *mux.Router, store store.Contract) {
				id := hash.New(test.FactoryRandomByteSlice())
				idContract := url.New(id.Printable())
				value := annotation.New(ulid.New().Get(), id, nil, metadataStub.NewNullObject())
				assert.Equal(t, status.Success, store.Create(idContract, value))

				response := testInternal.SendRequestWithBody(
					t,
					muxRouter,
					append.Method,
					append.EscapedRoute(idContract),
					testInternal.Marshal(t, value),
				)

				assert.Equal(t, append.CodeSuccess, response.Code)
				assert.Equal(t, testInternal.Marshal(t, status.Success), response.Body.Bytes())
			},
		},
		{
			name: "Find non-existent",
			test: func(t *testing.T, muxRouter *mux.Router, store store.Contract) {
				id := hash.New(test.FactoryRandomByteSlice())

				response := testInternal.SendRequestWithoutBody(
					t,
					muxRouter,
					find.Method,
					find.EscapedRoute(url.New(id.Printable())),
				)

				assert.Equal(t, find.CodeIdentityNotFound, response.Code)
			},
		},
		{
			name: "Find created",
			test: func(t *testing.T, muxRouter *mux.Router, store store.Contract) {
				id := hash.New(test.FactoryRandomByteSlice())
				idContract := url.New(id.Printable())
				value := annotation.New(ulid.New().Get(), id, nil, metadataStub.NewNullObject())
				assert.Equal(t, status.Success, store.Create(idContract, value))

				response := testInternal.SendRequestWithoutBody(
					t,
					muxRouter,
					find.Method,
					find.EscapedRoute(idContract),
				)

				assert.Equal(t, find.CodeSuccess, response.Code)
				assert.Equal(t, testInternal.Marshal(t, []interface{}{value}), response.Body.Bytes())
			},
		},
		{
			name: "Find appended",
			test: func(t *testing.T, muxRouter *mux.Router, store store.Contract) {
				id := hash.New(test.FactoryRandomByteSlice())
				idContract := url.New(id.Printable())
				createdValue := annotation.New(ulid.New().Get(), id, nil, metadataStub.NewNullObject())
				appendedValue := annotation.New(ulid.New().Get(), id, nil, metadataStub.NewNullObject())
				assert.Equal(t, status.Success, store.Create(idContract, createdValue))
				assert.Equal(t, status.Success, store.Append(idContract, appendedValue))

				response := testInternal.SendRequestWithoutBody(
					t,
					muxRouter,
					find.Method,
					find.EscapedRoute(idContract),
				)

				assert.Equal(t, find.CodeSuccess, response.Code)
				assert.Equal(
					t,
					testInternal.Marshal(t, []interface{}{createdValue, appendedValue}),
					response.Body.Bytes(),
				)
			},
		},
	}

	for i := range cases {
		s := memory.New()
		mFactory := metadataFactory.New(
			[]metadataFactory.Contract{
				assessMetadataFactory.NewDefault(),
				pkiMetadataFactory.NewDefault(),
				publishMetadataFactory.NewDefault(),
			},
		)
		iFactory := identityFactory.New()
		cancel, wg, muxRouter := testInternal.NewSUT(
			Run,
			[]routable.Contract{
				append.New(s, mFactory, iFactory).Init,
				create.New(s, mFactory, iFactory).Init,
				find.New(s).Init,
			},
		)
		t.Run(
			cases[i].name,
			func(t *testing.T) {
				cases[i].test(t, muxRouter, s)
				cancel()
				wg.Wait()
			},
		)
	}
}
