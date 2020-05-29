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

package create

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"

	urlIdentity "github.com/project-alvarium/go-store/internal/pkg/identity/url"

	"github.com/project-alvarium/go-sdk/pkg/annotation"
	metadataFactory "github.com/project-alvarium/go-sdk/pkg/annotation/metadata/factory"
	"github.com/project-alvarium/go-sdk/pkg/annotation/store"
	"github.com/project-alvarium/go-sdk/pkg/identity"
	identityFactory "github.com/project-alvarium/go-sdk/pkg/identity/factory"

	"github.com/gorilla/mux"
)

const (
	identityParam      = "identity"
	Method             = http.MethodPut
	codeBodyReadFailed = http.StatusBadRequest
	codeMarshalFailed  = http.StatusBadRequest
	CodeSuccess        = http.StatusOK
)

// Route creates a url.
func Route(id string) string {
	return fmt.Sprintf("/create/%s", id)
}

// EscapedRoute creates a url for client.
func EscapedRoute(id identity.Contract) string {
	return Route(url.PathEscape(id.Printable()))
}

// instance is a receiver that encapsulates required dependencies.
type instance struct {
	store    store.Contract
	mFactory metadataFactory.Contract
	iFactory identityFactory.Contract
}

// New is a factory function that returns instance.
func New(store store.Contract, mFactory metadataFactory.Contract, iFactory identityFactory.Contract) *instance {
	return &instance{
		store:    store,
		mFactory: mFactory,
		iFactory: iFactory,
	}
}

// Init adds package's route to muxRouter.
func (i *instance) Init(muxRouter *mux.Router) {
	muxRouter.HandleFunc(Route("{"+identityParam+"}"), i.handle).Methods(Method)
}

// handle implements package's functionality.
func (i *instance) handle(w http.ResponseWriter, r *http.Request) {
	id := urlIdentity.New(mux.Vars(r)[identityParam])

	body := make([]byte, r.ContentLength)
	if _, err := r.Body.Read(body); err != nil && err != io.EOF {
		w.WriteHeader(codeBodyReadFailed)
		return
	}

	var value annotation.Instance
	value.SetMetadataFactory(i.mFactory)
	value.SetIdentityFactory(i.iFactory)
	if err := json.Unmarshal(body, &value); err != nil {
		w.WriteHeader(codeMarshalFailed)
		return
	}

	resultInBytes, err := json.Marshal(i.store.Create(id, &value))
	if err != nil {
		w.WriteHeader(codeMarshalFailed)
		return
	}

	w.WriteHeader(CodeSuccess)
	_, _ = w.Write(resultInBytes)
}
