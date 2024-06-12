package configs

import (
	"github.com/hashicorp/hcl/v2"
	"github.com/hashicorp/hcl/v2/gohcl"
	"github.com/hashicorp/hcl/v2/hcldec"
	"github.com/opentofu/opentofu/internal/addrs"
	"github.com/opentofu/opentofu/internal/lang"
	"github.com/zclconf/go-cty/cty"
)

// StaticIdentifier holds a Referencable item and where it was declared
type StaticIdentifier struct {
	Module    addrs.Module
	Subject   addrs.Referenceable
	DeclRange hcl.Range
}

func (ref StaticIdentifier) String() string {
	val := ref.Subject.String()
	if len(ref.Module) != 0 {
		val = ref.Module.String() + "." + val
	}
	return val
}

type StaticModuleVariables func(v *Variable) (cty.Value, hcl.Diagnostics)

// StaticModuleCall contains the information required to call a given module
type StaticModuleCall struct {
	Addr addrs.Module

	Variables StaticModuleVariables

	RootPath string
}

func (c StaticModuleCall) Child(name string, vars StaticModuleVariables) StaticModuleCall {
	addr := make([]string, len(c.Addr)+1)
	copy(addr, c.Addr)
	addr[len(addr)-1] = name

	return StaticModuleCall{
		Addr:      addr,
		Variables: vars,
		RootPath:  c.RootPath,
	}
}

type StaticContext struct {
	Call StaticModuleCall
	cfg  *Module
}

// Creates a static context based from the given module and module call
func NewStaticContext(mod *Module, call StaticModuleCall) *StaticContext {
	return &StaticContext{
		Call: call,
		cfg:  mod,
	}
}

func (s *StaticContext) scope(ident StaticIdentifier) *lang.Scope {
	return newStaticScope(s, ident, nil)
}

func (s StaticContext) Evaluate(expr hcl.Expression, ident StaticIdentifier) (cty.Value, hcl.Diagnostics) {
	val, diags := s.scope(ident).EvalExpr(expr, cty.DynamicPseudoType)
	return val, diags.ToHCL()
}

func (s StaticContext) DecodeExpression(expr hcl.Expression, ident StaticIdentifier, val any) hcl.Diagnostics {
	var diags hcl.Diagnostics

	refs, refsDiags := lang.ReferencesInExpr(addrs.ParseRef, expr)
	diags = append(diags, refsDiags.ToHCL()...)
	if diags.HasErrors() {
		return diags
	}

	ctx, ctxDiags := s.scope(ident).EvalContext(refs)
	diags = append(diags, ctxDiags.ToHCL()...)
	if diags.HasErrors() {
		return diags
	}

	return gohcl.DecodeExpression(expr, ctx, val)
}

func (s StaticContext) DecodeBlock(body hcl.Body, spec hcldec.Spec, ident StaticIdentifier) (cty.Value, hcl.Diagnostics) {
	var diags hcl.Diagnostics

	refs, refsDiags := lang.References(addrs.ParseRef, hcldec.Variables(body, spec))
	diags = append(diags, refsDiags.ToHCL()...)
	if diags.HasErrors() {
		return cty.DynamicVal, diags
	}

	ctx, ctxDiags := s.scope(ident).EvalContext(refs)
	diags = append(diags, ctxDiags.ToHCL()...)
	if diags.HasErrors() {
		return cty.DynamicVal, diags
	}

	val, valDiags := hcldec.Decode(body, spec, ctx)
	diags = append(diags, valDiags...)
	return val, diags
}
