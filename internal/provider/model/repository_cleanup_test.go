/*
 * Copyright (c) 2019-present Sonatype, Inc.
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package model

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/types"
	sonatyperepo "github.com/sonatype-nexus-community/nexus-repo-api-client-go/v3"
	"github.com/stretchr/testify/assert"
)

// Nexus cannot store an empty cleanup block: sending empty policyNames and reading back
// yields no cleanup. If the read nils the state, a configuration declaring an empty cleanup
// never converges and Terraform/Crossplane keeps rewriting it (measured on a live cluster
// on 2026-09-14: 600 writes per minute).
func TestCleanupFromApiKeepsEmptyCleanup(t *testing.T) {
	current := NewRepositoryCleanupModel()

	got := cleanupFromApi(nil, current)

	assert.NotNil(t, got, "an empty cleanup must be preserved; nil-ing it is a write loop")
	assert.Same(t, current, got)
	assert.Empty(t, got.PolicyNames)
}

func TestCleanupFromApiStaysNilWhenNothingIsSet(t *testing.T) {
	assert.Nil(t, cleanupFromApi(nil, nil))
	assert.Nil(t, cleanupFromApi(&sonatyperepo.CleanupPolicyAttributes{}, nil))
}

func TestCleanupFromApiReadsPolicyNames(t *testing.T) {
	api := &sonatyperepo.CleanupPolicyAttributes{PolicyNames: []string{"proxy-cache-30d"}}

	got := cleanupFromApi(api, nil)

	assert.NotNil(t, got)
	assert.Equal(t, []types.String{types.StringValue("proxy-cache-30d")}, got.PolicyNames)
}

// A policy removed from the repository on the Nexus side IS real drift: the state must
// drop to nil so Terraform sees the change instead of keeping the stale names.
func TestCleanupFromApiReportsRemovedPolicyAsDrift(t *testing.T) {
	current := &RepositoryCleanupModel{PolicyNames: []types.String{types.StringValue("removed-upstream")}}

	assert.Nil(t, cleanupFromApi(nil, current))
}

// The previous code assigned `m = nil` and then `return *m`, panicking for every
// repository without a cleanup policy.
func TestRepositoryCleanupModelMapFromApiDoesNotPanic(t *testing.T) {
	m := NewRepositoryCleanupModel()

	assert.NotPanics(t, func() { m.MapFromApi(nil) })
	assert.Empty(t, m.PolicyNames)

	out := m.MapFromApi(&sonatyperepo.CleanupPolicyAttributes{PolicyNames: []string{"p1"}})
	assert.Equal(t, []types.String{types.StringValue("p1")}, out.PolicyNames)
}
