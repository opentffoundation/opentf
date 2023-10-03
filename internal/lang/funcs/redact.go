// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package funcs

import (
	"fmt"

	"github.com/zclconf/go-cty/cty"

	"github.com/opentofu/opentofu/internal/lang/marks"
)

func redactIfSensitive(value interface{}, markses ...cty.ValueMarks) string {
	if marks.Has(cty.DynamicVal.WithMarks(markses...), marks.Sensitive) {
		return "(sensitive value)"
	}
	switch v := value.(type) {
	case string:
		return fmt.Sprintf("%q", v)
	default:
		return fmt.Sprintf("%v", v)
	}
}
