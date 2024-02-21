package configs

import (
	"github.com/hashicorp/hcl/v2"
	"github.com/opentofu/opentofu/internal/configs/hcl2shim"
	"github.com/zclconf/go-cty/cty"
)

// These were all moved to the hcl2shim package, but still have uses referenced from this package
func MergeBodies(base, override hcl.Body) hcl.Body {
	return hcl2shim.MergeBodies(base, override)
}

func exprIsNativeQuotedString(expr hcl.Expression) bool {
	return hcl2shim.ExprIsNativeQuotedString(expr)
}

func schemaForOverrides(schema *hcl.BodySchema) *hcl.BodySchema {
	return hcl2shim.SchemaForOverrides(schema)
}

func SynthBody(filename string, values map[string]cty.Value) hcl.Body {
	return hcl2shim.SynthBody(filename, values)
}
