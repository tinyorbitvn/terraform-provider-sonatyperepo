/*
 * Copyright (c) 2019-present Sonatype, Inc.
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *      http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

// Package xpprovider exposes the provider constructor for embedding.
//
// Crossplane providers generated with Upjet run Terraform providers
// in-process and therefore need to import the provider constructor. The
// constructor lives in internal/provider, which Go does not allow other
// modules to import; this package re-exports it. It has no other purpose
// and no runtime cost for the Terraform binary.
package xpprovider

import (
	fwprovider "github.com/hashicorp/terraform-plugin-framework/provider"

	"github.com/tinyorbitvn/terraform-provider-sonatyperepo/internal/provider"
)

// New returns a configured-later terraform-plugin-framework provider,
// identical to what main.go serves.
func New(version string) fwprovider.Provider {
	return provider.New(version)()
}
