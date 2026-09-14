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

package model

import (
	"github.com/hashicorp/terraform-plugin-framework/types"

	sonatyperepo "github.com/sonatype-nexus-community/nexus-repo-api-client-go/v3"
)

type RepositoryModel struct {
	Name   types.String `tfsdk:"name"`
	Format types.String `tfsdk:"format"`
	Type   types.String `tfsdk:"type"`
	Url    types.String `tfsdk:"url"`
}

type RepositoriesModel struct {
	Repositories []RepositoryModel `tfsdk:"repositories"`
}

type BasicRepositoryModel struct {
	Name        types.String            `tfsdk:"name"`
	Online      types.Bool              `tfsdk:"online"`
	Url         types.String            `tfsdk:"url"`
	Cleanup     *RepositoryCleanupModel `tfsdk:"cleanup"`
	LastUpdated types.String            `tfsdk:"last_updated"`
}

type RepositoryCleanupModel struct {
	PolicyNames []types.String `tfsdk:"policy_names"`
}

func (m *RepositoryCleanupModel) MapFromApi(api *sonatyperepo.CleanupPolicyAttributes) RepositoryCleanupModel {
	// This used to assign `m = nil` and then `return *m`, panicking for every repository
	// without a cleanup policy. Write through the receiver and return a copy instead.
	m.PolicyNames = make([]types.String, 0)
	if api != nil {
		mapCleanupFromApi(api, m)
	}
	return *m
}

// cleanupFromApi builds the `cleanup` model from an API response, treating "no cleanup"
// and "cleanup with an empty policy_names" as the same thing.
//
// Why: Nexus cannot store an empty cleanup block. Sending `{"policyNames": []}` and
// reading back yields no cleanup at all. The previous code nil-ed the state in that case,
// so a configuration declaring `cleanup { policy_names = [] }` never converged: read
// gives null, null differs from the configuration, update, read gives null again...
// Measured on a live cluster on 2026-09-14 while driven by Crossplane: the provider
// rewrote the same repository 600 times per minute, each write stopping and starting a
// repository on a production registry.
//
// Keeping an existing empty state when the API returns no policies makes the read
// idempotent and the loop disappears. Behaviour with real policies is unchanged.
func cleanupFromApi(api *sonatyperepo.CleanupPolicyAttributes, current *RepositoryCleanupModel) *RepositoryCleanupModel {
	if api != nil && len(api.PolicyNames) > 0 {
		m := NewRepositoryCleanupModel()
		mapCleanupFromApi(api, m)
		return m
	}
	if current != nil && len(current.PolicyNames) == 0 {
		return current
	}
	return nil
}

func NewRepositoryCleanupModel() *RepositoryCleanupModel {
	return &RepositoryCleanupModel{
		PolicyNames: make([]types.String, 0),
	}
}

func mapCleanupFromApi(api *sonatyperepo.CleanupPolicyAttributes, m *RepositoryCleanupModel) {
	for _, p := range api.GetPolicyNames() {
		m.PolicyNames = append(m.PolicyNames, types.StringValue(p))
	}
}

func mapCleanupToApi(m *RepositoryCleanupModel, api *sonatyperepo.CleanupPolicyAttributes) {
	if m != nil {
		for _, p := range m.PolicyNames {
			api.PolicyNames = append(api.PolicyNames, p.ValueString())
		}
	}
}

// repositoryStorageModel
// ----------------------------------------
type repositoryStorageModel struct {
	BlobStoreName               types.String `tfsdk:"blob_store_name"`
	StrictContentTypeValidation types.Bool   `tfsdk:"strict_content_type_validation"`
}

func (m *repositoryStorageModel) MapFromApi(api *sonatyperepo.StorageAttributes) {
	m.BlobStoreName = types.StringValue(api.BlobStoreName)
	m.StrictContentTypeValidation = types.BoolValue(api.StrictContentTypeValidation)

}

func (m *repositoryStorageModel) MapToApi(api *sonatyperepo.StorageAttributes) {
	api.BlobStoreName = m.BlobStoreName.ValueString()
	api.StrictContentTypeValidation = m.StrictContentTypeValidation.ValueBool()
}
