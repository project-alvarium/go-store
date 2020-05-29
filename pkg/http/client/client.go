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
	identityFactory "github.com/project-alvarium/go-sdk/pkg/identity/factory"
)

// Requestor defines the contract used to delegate http requests.
type Requestor func(method, path string, body []byte) (responseBody []byte, err error)

// instance is a receiver that encapsulates required dependencies.
type instance struct {
	requestor Requestor
	mFactory  metadataFactory.Contract
	iFactory  identityFactory.Contract
}

// New is a factory function that returns instance.
func New(requestor Requestor, mFactory metadataFactory.Contract, iFactory identityFactory.Contract) *instance {
	return &instance{
		requestor: requestor,
		mFactory:  mFactory,
		iFactory:  iFactory,
	}
}
