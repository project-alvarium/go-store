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
	"context"
	"sync"

	"github.com/project-alvarium/go-store/internal/pkg/interrupt"
	"github.com/project-alvarium/go-store/internal/pkg/routable"
	"github.com/project-alvarium/go-store/internal/pkg/server"

	"github.com/gorilla/mux"
)

// Run is the internal main entry point.
func Run(
	ctx context.Context,
	cancel context.CancelFunc,
	muxRouter *mux.Router,
	routables []routable.Contract,
	serverAddress *string) {

	for key := range routables {
		routables[key](muxRouter)
	}

	if serverAddress != nil {
		var wg sync.WaitGroup
		interrupt.TranslateToCancel(ctx, cancel, &wg)
		server.Serve(ctx, muxRouter, &wg, *serverAddress)
		wg.Wait()
	}
}
