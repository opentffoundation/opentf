package encryptionconfig

import (
	"errors"
	"testing"
)

func TestRegisterKeyProviderConfigValidationFunction(t *testing.T) {
	testCases := []struct {
		testcase    string
		name        KeyProviderName
		validator   KeyProviderValidator
		expectedErr error
	}{
		{
			testcase:    "duplicate", // already done in init()
			name:        KeyProviderPassphrase,
			validator:   validateKPPassphraseConfig,
			expectedErr: errors.New("duplicate registration for key provider \"passphrase\""),
		},
		{
			testcase:    "no_validator",
			name:        "some_fancy_kms",
			validator:   nil,
			expectedErr: errors.New("missing validator during registration for key provider \"some_fancy_kms\": nil"),
		},
		// success case covered via init()
	}
	for _, tc := range testCases {
		t.Run(tc.testcase, func(t *testing.T) {
			err := RegisterKeyProviderValidator(tc.name, tc.validator)
			expectErr(t, err, tc.expectedErr)
		})
	}
}

func TestRegisterEncryptionMethodConfigValidationFunction(t *testing.T) {
	testCases := []struct {
		testcase    string
		name        MethodName
		validator   MethodValidator
		expectedErr error
	}{
		{
			testcase:    "duplicate", // already done in init()
			name:        MethodFull,
			validator:   validateEMFullConfig,
			expectedErr: errors.New("duplicate registration for encryption method \"full\""),
		},
		{
			testcase:    "no_validator",
			name:        "base64encrypt",
			validator:   nil,
			expectedErr: errors.New("missing validator during registration for encryption method \"base64encrypt\": nil"),
		},
		// success case covered via init()
	}
	for _, tc := range testCases {
		t.Run(tc.testcase, func(t *testing.T) {
			err := RegisterMethodValidator(tc.name, tc.validator)
			expectErr(t, err, tc.expectedErr)
		})
	}
}