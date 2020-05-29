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
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"

	urlIdentity "github.com/michaelestrin/go-store/internal/pkg/identity/url"

	"github.com/project-alvarium/go-sdk/pkg/annotation/store"
	"github.com/project-alvarium/go-sdk/pkg/identity"
	"github.com/project-alvarium/go-sdk/pkg/status"

	"github.com/gorilla/mux"
)

const (
	identityParam        = "identity"
	Method               = http.MethodGet
	CodeIdentityNotFound = http.StatusBadRequest
	codeMarshalFailed    = http.StatusBadRequest
	CodeSuccess          = http.StatusOK
)

// Route creates a url.
func Route(id string) string {
	return fmt.Sprintf("/findByIdentity/%s", id)
}

// EscapedRoute creates a url for client.
func EscapedRoute(id identity.Contract) string {
	return Route(url.PathEscape(id.Printable()))
}

// instance is a receiver that encapsulates required dependencies.
type instance struct {
	store store.Contract
}

// New is a factory function that returns instance.
func New(store store.Contract) *instance {
	return &instance{
		store: store,
	}
}

// Init adds package's route to muxRouter.
func (i *instance) Init(muxRouter *mux.Router) {
	muxRouter.HandleFunc(Route("{"+identityParam+"}"), i.handle).Methods(Method)
}

// handle implements package's functionality.
func (i *instance) handle(w http.ResponseWriter, r *http.Request) {
	value, result := i.store.FindByIdentity(urlIdentity.New(mux.Vars(r)[identityParam]))
	if result != status.Success {
		w.WriteHeader(CodeIdentityNotFound)
		return
	}

	body, err := json.Marshal(value)
	if err != nil {
		w.WriteHeader(codeMarshalFailed)
		return
	}

	w.WriteHeader(CodeSuccess)
	_, _ = w.Write(body)
}
