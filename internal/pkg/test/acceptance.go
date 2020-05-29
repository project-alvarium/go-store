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

package test

import (
	"context"
	"sync"

	"github.com/michaelestrin/go-store/internal/pkg/routable"

	"github.com/gorilla/mux"
)

// MainFunc defines the signature of the function used to start a service.
type RunFunc func(
	ctx context.Context,
	cancel context.CancelFunc,
	muxRouter *mux.Router,
	routables []routable.Contract,
	serverAddress *string,
)

// NewSUT returns a new system under test for acceptance testing.
func NewSUT(runFunc RunFunc, routables []routable.Contract) (context.CancelFunc, *sync.WaitGroup, *mux.Router) {
	var wg sync.WaitGroup
	ctx, cancel := context.WithCancel(context.Background())
	muxRouter := mux.NewRouter()

	runFunc(ctx, cancel, muxRouter, routables, nil)

	return cancel, &wg, muxRouter
}
