// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package moduletest

import (
	"github.com/opentofu/opentofu/internal/configs"
	"github.com/opentofu/opentofu/internal/tfdiags"
)

type File struct {
	Config *configs.TestFile

	Name   string
	Status Status

	Runs []*Run

	Diagnostics tfdiags.Diagnostics
}
