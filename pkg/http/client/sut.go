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
	metadataFactory "github.com/project-alvarium/go-sdk/pkg/annotation/metadata/factory"
	assessMetadataFactory "github.com/project-alvarium/go-sdk/pkg/annotator/assess/metadata/factory"
	pkiMetadataFactory "github.com/project-alvarium/go-sdk/pkg/annotator/pki/metadata/factory"
	publishMetadataFactory "github.com/project-alvarium/go-sdk/pkg/annotator/publish/metadata/factory"
	identityFactory "github.com/project-alvarium/go-sdk/pkg/identity/factory"
)

// newSUT returns a new system under test.
func newSUTWithFactories(requestor Requestor, metadataFactories []metadataFactory.Contract) *instance {
	mFactory := metadataFactory.New(metadataFactories)
	iFactory := identityFactory.New()
	return New(requestor, mFactory, iFactory)
}

// newSUT returns a new system under test.
func newSUT(requestor Requestor) *instance {
	return newSUTWithFactories(
		requestor,
		[]metadataFactory.Contract{
			assessMetadataFactory.NewDefault(),
			pkiMetadataFactory.NewDefault(),
			publishMetadataFactory.NewDefault(),
		},
	)
}
