// Copyright (c) The OpenTofu Authors
// SPDX-License-Identifier: MPL-2.0
// Copyright (c) 2023 HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package pbkdf2_test

import (
	"fmt"

	"github.com/hashicorp/hcl/v2/gohcl"
	"github.com/opentofu/opentofu/internal/encryption/keyprovider/pbkdf2"

	"github.com/opentofu/opentofu/internal/encryption/config"
)

var configuration = `key_provider "pbkdf2" "foo" {
  passphrase = "Hello world!"
}
`

// This example is a bare-bones configuration for a static key provider.
// It is mainly intended to demonstrate how you can use parse configuration
// and construct a static key provider from in.
// And is not intended to be used as a real-world example.
func Example_decrypt() {
	configStruct := pbkdf2.New().ConfigStruct()

	// Parse the config:
	parsedConfig, diags := config.LoadConfigFromString("config.hcl", configuration)
	if diags.HasErrors() {
		panic(diags)
	}

	// Use gohcl to parse the hcl block from parsedConfig into the static configuration struct:
	if err := gohcl.DecodeBody(
		parsedConfig.KeyProviderConfigs[0].Body,
		nil,
		configStruct,
	); err != nil {
		panic(err)
	}

	// Create the actual key provider.
	keyProvider, keyMeta, err := configStruct.Build()
	if err != nil {
		panic(err)
	}

	// Fill in the metadata stored with the encrypted form:
	meta := keyMeta.(*pbkdf2.Metadata)
	meta.Salt = []byte{0x10, 0xec, 0x3d, 0x3f, 0xe0, 0x2a, 0xd2, 0xbe, 0xe6, 0xf1, 0xf5, 0x54, 0xf, 0x8e, 0x6b, 0xbe, 0x3b, 0x8b, 0x29, 0x44, 0x5c, 0xf5, 0x2, 0xd2, 0x7d, 0x47, 0xad, 0x55, 0x4a, 0xa8, 0x97, 0x1f}
	meta.Iterations = 600000
	meta.HashFunction = "sha512"

	// Get decryption key from the provider.
	keys, _, err := keyProvider.Provide(meta)
	if err != nil {
		panic(err)
	}

	fmt.Printf("%x", keys.DecryptionKey)
	// Output: 7919af5a183ed2eb8bef7ab7555f5e9e3381afb91dbbc315be438a79de5c5fbd
}
