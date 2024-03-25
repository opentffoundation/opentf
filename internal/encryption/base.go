// Copyright (c) The OpenTofu Authors
// SPDX-License-Identifier: MPL-2.0
// Copyright (c) 2023 HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package encryption

import (
	"encoding/json"
	"fmt"

	"github.com/opentofu/opentofu/internal/encryption/config"
	"github.com/opentofu/opentofu/internal/encryption/keyprovider"
	"github.com/opentofu/opentofu/internal/encryption/method"

	"github.com/hashicorp/hcl/v2"
)

const (
	encryptionVersion = "v0"
)

type baseEncryption struct {
	enc        *encryption
	target     *config.TargetConfig
	enforced   bool
	name       string
	encMethods []method.Method
	encMeta    map[keyprovider.Addr][]byte
}

func newBaseEncryption(enc *encryption, target *config.TargetConfig, enforced bool, name string) (*baseEncryption, hcl.Diagnostics) {
	base := &baseEncryption{
		enc:      enc,
		target:   target,
		enforced: enforced,
		name:     name,
		encMeta:  make(map[keyprovider.Addr][]byte),
	}
	// Setup the encryptor
	//
	//     Instead of creating new encryption key data for each call to encrypt, we use the same encryptor for the given application (statefile or planfile).
	//
	// Why do we do this?
	//
	//   This allows us to always be in a state where we can encrypt data, which is particularly important when dealing with crashes. If the network is severed
	//   mid-apply, we still need to be able to write an encrypted errored.tfstate or dump to stdout. Additionally it reduces the overhead of encryption in
	//   general, as well as reducing cloud key provider costs.
	//
	// What are the security implications?
	//
	//   Plan file flow is fairly simple and is not impacted by this change. It only ever calls encrypt once at the end of plan generation.
	//
	//   State file is a bit more complex. The encrypted local state file (terraform.tfstate, .terraform.tfstate) will be written with the same
	//   keys as any remote state. These files should be identical, which will make debugging easier.
	//
	//   The major concern with this is that many of the encryption methods used have a limit to how much data a key can safely encrypt. Pbkdf2 for example
	//   has a limit of around 64GB before exhaustion is approached. Writing to the two local and one remote location specified above could not
	//   approach that limit. However the cached state file (.terraform/terraform.tfstate) is persisted every 30 seconds during long applies. For an
	//   extremely large state file (100MB) it would take an apply of over 5 hours to come close to the 64GB limit of pbkdf2 with some malicious actor recording
	//   every single change to the filesystem (or inspecting deleted blocks).
	//
	// What other benfits does this provide?
	//
	//   This performs a e2e validation run of the config -> methods flow. It serves as a validation step and allows us to return detailed
	//   diagnostics here and simple errors in the decrypt function below.
	//
	methods, diags := base.buildTargetMethods(base.encMeta)
	base.encMethods = methods
	return base, diags
}

type basedata struct {
	Meta    map[keyprovider.Addr][]byte `json:"meta"`
	Data    []byte                      `json:"encrypted_data"`
	Version string                      `json:"encryption_version"` // This is both a sigil for a valid encrypted payload and a future compatability field
}

func IsEncryptionPayload(data []byte) (bool, error) {
	es := basedata{}
	err := json.Unmarshal(data, &es)
	if err != nil {
		return false, err
	}

	// This could be extended with full version checking later on
	return es.Version != "", nil
}

func (s *baseEncryption) encrypt(data []byte) ([]byte, error) {
	// No configuration provided, don't do anything
	if s.target == nil {
		return data, nil
	}

	var encryptor method.Method = nil
	if len(s.encMethods) != 0 {
		// Use the pre-configured encryption method
		encryptor = s.encMethods[0]
	}

	if encryptor == nil {
		// ensure that the method is defined when Enforced is true
		if s.enforced {
			return nil, fmt.Errorf("encryption of %q is enforced, and therefore requires a method to be provided", s.name)
		}
		return data, nil
	}

	encd, err := encryptor.Encrypt(data)
	if err != nil {
		return nil, fmt.Errorf("encryption failed for %s: %w", s.name, err)
	}

	es := basedata{
		Version: encryptionVersion,
		Meta:    s.encMeta,
		Data:    encd,
	}
	jsond, err := json.Marshal(es)
	if err != nil {
		return nil, fmt.Errorf("unable to encode encrypted data as json: %w", err)
	}

	return jsond, nil
}

// TODO Find a way to make these errors actionable / clear
func (s *baseEncryption) decrypt(data []byte, validator func([]byte) error) ([]byte, error) {
	if s.target == nil {
		return data, nil
	}

	es := basedata{}
	err := json.Unmarshal(data, &es)

	if len(es.Version) == 0 || err != nil {
		// Not a valid payload, might be already decrypted
		verr := validator(data)
		if verr != nil {
			// Nope, just bad input

			// Return the outer json error if we have one
			if err != nil {
				return nil, fmt.Errorf("invalid data format for decryption: %w, %w", err, verr)
			}

			// Must have been invalid json payload
			return nil, fmt.Errorf("unable to determine data structure during decryption: %w", verr)
		}

		methods, diags := s.buildTargetMethods(make(map[keyprovider.Addr][]byte))
		if diags.HasErrors() {
			// This cast to error here is safe as we know that at least one error exists
			// This is also quite unlikely to happen as the constructor already has checked this code path
			return nil, diags
		}
		// Yep, it's already decrypted
		for _, method := range methods {
			if method == nil {
				// fallback allowed
				return data, nil
			}
		}
		return data, fmt.Errorf("decrypted payload provided without fallback specified")
	}

	if es.Version != encryptionVersion {
		return nil, fmt.Errorf("invalid encrypted payload version: %s != %s", es.Version, encryptionVersion)
	}

	// TODO Discuss if we should potentially cache this based on a json-encoded version of es.Meta and reduce overhead dramatically
	methods, diags := s.buildTargetMethods(es.Meta)
	if diags.HasErrors() {
		// This cast to error here is safe as we know that at least one error exists
		// This is also quite unlikely to happen as the constructor already has checked this code path
		return nil, diags
	}

	if len(methods) == 0 {
		err = validator(data)
		if err != nil {
			// TODO improve this error message
			return nil, err
		}
		// No methods/fallbacks specified and data is valid payload
		return data, nil
	}

	errs := make([]error, 0)
	for _, method := range methods {
		if method == nil {
			// No method specified for this target
			err = validator(data)
			if err == nil {
				return data, nil
			}
			// TODO improve this error message
			errs = append(errs, fmt.Errorf("payload is not already decrypted: %w", err))
			continue
		}
		uncd, err := method.Decrypt(es.Data)
		if err == nil {
			// Success
			return uncd, nil
		}
		// Record the failure
		errs = append(errs, fmt.Errorf("attempted decryption failed for %s: %w", s.name, err))
	}

	// This is good enough for now until we have better/distinct errors
	errMessage := "decryption failed for all provided methods: "
	sep := ""
	for _, err := range errs {
		errMessage += err.Error() + sep
		sep = "\n"
	}
	return nil, fmt.Errorf(errMessage)
}
