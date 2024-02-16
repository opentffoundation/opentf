// Copyright (c) The OpenTofu Authors
// SPDX-License-Identifier: MPL-2.0
// Copyright (c) 2023 HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package config

import (
	"github.com/hashicorp/hcl/v2"
	"github.com/opentofu/opentofu/internal/encryption/keyprovider"
	"github.com/opentofu/opentofu/internal/encryption/method"
)

// Config describes the terraform.encryption HCL block you can use to configure the state and plan encryption.
// The individual fields of this struct match the HCL structure directly.
type Config struct {
	KeyProviderConfigs []KeyProviderConfig `hcl:"key_provider,block"`
	MethodConfigs      []MethodConfig      `hcl:"method,block"`

	Backend   *EnforcableTargetConfig `hcl:"backend,block"`
	StateFile *EnforcableTargetConfig `hcl:"statefile,block"`
	PlanFile  *EnforcableTargetConfig `hcl:"planfile,block"`
	Remote    *RemoteConfig           `hcl:"remote_data_source,block"`
}

// Merge returns a merged configuration with  the current config and the specified override combined, the override
// taking precedence.
func (c *Config) Merge(override *Config) *Config {
	return MergeConfigs(c, override)
}

// KeyProviderConfig describes the terraform.encryption.key_provider.* block you can use to declare a key provider for
// encryption. The Body field will contain the remaining undeclared fields the key provider can consume.
type KeyProviderConfig struct {
	Type string   `hcl:"type,label"`
	Name string   `hcl:"name,label"`
	Body hcl.Body `hcl:",remain"`
}

// Addr returns a keyprovider.Addr from the current configuration.
func (k KeyProviderConfig) Addr() (keyprovider.Addr, hcl.Diagnostics) {
	return keyprovider.NewAddr(k.Type, k.Name)
}

// MethodConfig describes the terraform.encryption.method.* block you can use to declare the encryption method. The Body
// field will contain the remaining undeclared fields the method can consume.
type MethodConfig struct {
	Type string   `hcl:"type,label"`
	Name string   `hcl:"name,label"`
	Body hcl.Body `hcl:",remain"`
}

func (m MethodConfig) Addr() (method.Addr, hcl.Diagnostics) {
	return method.NewAddr(m.Type, m.Name)
}

// RemoteConfig describes the terraform.encryption.remote block you can use to declare encryption for remote state data
// sources.
type RemoteConfig struct {
	Default *TargetConfig       `hcl:"default,block"`
	Targets []NamedTargetConfig `hcl:"remote_data_source,block"`
}

// TargetConfig describes the target.encryption.state, target.encryption.plan, etc blocks.
type TargetConfig struct {
	Method   hcl.Expression `hcl:"method,optional"`
	Fallback *TargetConfig  `hcl:"fallback,block"`
}

// EnforcableTargetConfig is an extension of the TargetConfig that supports the enforced form.
//
// Note: This struct is copied because gohcl does not support embedding.
type EnforcableTargetConfig struct {
	Enforced bool           `hcl:"enforced,optional"`
	Method   hcl.Expression `hcl:"method,optional"`
	Fallback *TargetConfig  `hcl:"fallback,block"`
}

// AsTargetConfig converts the struct into its parent TargetConfig.
func (e EnforcableTargetConfig) AsTargetConfig() *TargetConfig {
	return &TargetConfig{
		Method:   e.Method,
		Fallback: e.Fallback,
	}
}

// NamedTargetConfig is an extension of the TargetConfig that describes a
// terraform.encryption.remote.remote_state_data.* block.
//
// Note: This struct is copied because gohcl does not support embedding.
type NamedTargetConfig struct {
	Name     string         `hcl:"name,label"`
	Method   hcl.Expression `hcl:"method,optional"`
	Fallback *TargetConfig  `hcl:"fallback,block"`
}

// AsTargetConfig converts the struct into its parent TargetConfig.
func (n NamedTargetConfig) AsTargetConfig() *TargetConfig {
	return &TargetConfig{
		Method:   n.Method,
		Fallback: n.Fallback,
	}
}
