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

package find

import (
	"testing"

	"github.com/project-alvarium/go-store/internal/pkg"
	"github.com/project-alvarium/go-store/internal/pkg/identity/url"
	"github.com/project-alvarium/go-store/internal/pkg/routable"
	testInternal "github.com/project-alvarium/go-store/internal/pkg/test"

	"github.com/project-alvarium/go-sdk/pkg/annotation"
	metadataStub "github.com/project-alvarium/go-sdk/pkg/annotation/metadata/stub"
	"github.com/project-alvarium/go-sdk/pkg/annotation/store"
	"github.com/project-alvarium/go-sdk/pkg/annotation/store/memory"
	"github.com/project-alvarium/go-sdk/pkg/annotation/uniqueprovider/ulid"
	"github.com/project-alvarium/go-sdk/pkg/identity/hash"
	"github.com/project-alvarium/go-sdk/pkg/test"

	"github.com/gorilla/mux"
	"github.com/stretchr/testify/assert"
)

// TestFind tests find route.
func TestFind(t *testing.T) {
	type testCase struct {
		name string
		test func(t *testing.T, muxRouter *mux.Router, store store.Contract)
	}

	cases := []testCase{
		{
			name: "Identity not found",
			test: func(t *testing.T, muxRouter *mux.Router, _ store.Contract) {
				response := testInternal.SendRequestWithoutBody(t, muxRouter, Method, Route(test.FactoryRandomString()))

				assert.Equal(t, CodeIdentityNotFound, response.Code)
				assert.Nil(t, response.Body.Bytes())
			},
		},
		{
			name: "Success",
			test: func(t *testing.T, muxRouter *mux.Router, store store.Contract) {
				id := hash.New(test.FactoryRandomByteSlice())
				idContract := url.New(id.Printable())
				value := annotation.New(ulid.New().Get(), id, nil, metadataStub.NewNullObject())
				store.Create(idContract, value)

				response := testInternal.SendRequestWithoutBody(t, muxRouter, Method, EscapedRoute(idContract))

				assert.Equal(t, CodeSuccess, response.Code)
				assert.Equal(t, testInternal.Marshal(t, []*annotation.Instance{value}), response.Body.Bytes())
			},
		},
	}

	for i := range cases {
		s := memory.New()
		cancel, wg, muxRouter := testInternal.NewSUT(pkg.Run, []routable.Contract{New(s).Init})
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
