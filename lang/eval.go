package lang

import (
	"fmt"

	"github.com/hashicorp/hcl2/ext/dynblock"
	"github.com/hashicorp/hcl2/hcl"
	"github.com/hashicorp/hcl2/hcldec"
	"github.com/hashicorp/terraform/config/configschema"
	"github.com/hashicorp/terraform/tfdiags"
	"github.com/zclconf/go-cty/cty"
	"github.com/zclconf/go-cty/cty/convert"
)

// ExpandBlock expands any "dynamic" blocks present in the given body. The
// result is a body with those blocks expanded, ready to be evaluated with
// EvalBlock.
//
// If the returned diagnostics contains errors then the result may be
// incomplete or invalid.
func (ctx *Context) ExpandBlock(body hcl.Body, schema *configschema.Block) (hcl.Body, tfdiags.Diagnostics) {
	spec := schema.DecoderSpec()
	traversals := dynblock.ForEachVariablesHCLDec(body, spec)
	vals, diags := ctx.Scope.RefValues(traversals)
	funcs := ctx.Functions()
	evalCtx := &hcl.EvalContext{
		Variables: vals,
		Functions: funcs,
	}
	return dynblock.Expand(body, evalCtx), diags
}

// EvalBlock evaluates the given body using the given block schema and returns
// a cty object value representing its contents. The type of the result is
// the implied object type of the given schema.
//
// This function does not automatically expand "dynamic" blocks within the
// body. If that is desired, first call the ExpandBlock method to obtain
// an expanded body to pass to this method.
//
// If the returned diagnostics contains errors then the result may be
// incomplete or invalid.
func (ctx *Context) EvalBlock(body hcl.Body, schema *configschema.Block) (cty.Value, tfdiags.Diagnostics) {
	spec := schema.DecoderSpec()
	traversals := hcldec.Variables(body, spec)
	vals, diags := ctx.Scope.RefValues(traversals)
	funcs := ctx.Functions()
	evalCtx := &hcl.EvalContext{
		Variables: vals,
		Functions: funcs,
	}
	val, evalDiags := hcldec.Decode(body, spec, evalCtx)
	diags = diags.Append(evalDiags)
	return val, diags
}

// EvalExpr evaluates a single expression in the receiving context and returns
// the resulting value. The value will be converted to the given type before
// it is returned if possible, or else an error diagnostic will be produced
// describing the conversion error.
//
// Pass an expected type of cty.DynamicPseudoType to skip automatic conversion
// and just obtain the returned value directly.
//
// If the returned diagnostics contains errors then the result may be
// incomplete, but will always be of the requested type.
func (ctx *Context) EvalExpr(expr hcl.Expression, wantType cty.Type) (cty.Value, tfdiags.Diagnostics) {
	traversals := expr.Variables()
	vals, diags := ctx.Scope.RefValues(traversals)
	funcs := ctx.Functions()
	evalCtx := &hcl.EvalContext{
		Variables: vals,
		Functions: funcs,
	}
	val, evalDiags := expr.Value(evalCtx)
	diags = diags.Append(evalDiags)

	var convErr error
	val, convErr = convert.Convert(val, wantType)
	if convErr != nil {
		val = cty.UnknownVal(wantType)
		diags = diags.Append(&hcl.Diagnostic{
			Severity: hcl.DiagError,
			Summary:  "Incorrect value type",
			Detail:   fmt.Sprintf("Invalid expression value: %s.", tfdiags.FormatError(convErr)),
			Subject:  expr.Range().Ptr(),
		})
	}

	return val, diags
}
