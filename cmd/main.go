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

package main

import (
	"context"
	"flag"

	"github.com/michaelestrin/go-store/internal/pkg"
	"github.com/michaelestrin/go-store/internal/pkg/routable"
	"github.com/michaelestrin/go-store/internal/pkg/routes/append"
	"github.com/michaelestrin/go-store/internal/pkg/routes/create"
	"github.com/michaelestrin/go-store/internal/pkg/routes/find"

	metadataFactory "github.com/project-alvarium/go-sdk/pkg/annotation/metadata/factory"
	"github.com/project-alvarium/go-sdk/pkg/annotation/store/memory"
	assessMetadataFactory "github.com/project-alvarium/go-sdk/pkg/annotator/assess/metadata/factory"
	pkiMetadataFactory "github.com/project-alvarium/go-sdk/pkg/annotator/pki/metadata/factory"
	publishMetadataFactory "github.com/project-alvarium/go-sdk/pkg/annotator/publish/metadata/factory"
	identityFactory "github.com/project-alvarium/go-sdk/pkg/identity/factory"

	"github.com/gorilla/mux"
)

// main is the service's entry point.
func main() {
	var serverAddress string
	flag.StringVar(&serverAddress, "server", "localhost:8080", "Server address (localhost:8080)")
	flag.Parse()

	s := memory.New()
	mFactory := metadataFactory.New(
		[]metadataFactory.Contract{
			assessMetadataFactory.NewDefault(),
			pkiMetadataFactory.NewDefault(),
			publishMetadataFactory.NewDefault(),
		},
	)
	iFactory := identityFactory.New()
	ctx, cancel := context.WithCancel(context.Background())
	pkg.Run(
		ctx,
		cancel,
		mux.NewRouter().UseEncodedPath(),
		[]routable.Contract{
			find.New(s).Init,
			create.New(s, mFactory, iFactory).Init,
			append.New(s, mFactory, iFactory).Init,
		},
		&serverAddress,
	)
}
