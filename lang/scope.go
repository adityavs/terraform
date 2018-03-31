package lang

import (
	"github.com/hashicorp/hcl2/hcl"
	"github.com/hashicorp/terraform/tfdiags"
	"github.com/zclconf/go-cty/cty"
)

// A Scope is an object that can provide values for references that appear
// within expressions.
type Scope interface {
	// RefValues receives a set of traversals and must produce a map that
	// includes a value for each of these traversals. The result must ensure
	// that all given traversals are possible, but unknown values may be
	// inserted if a given traversal is erroneous.
	//
	// The returned map may contain additional values not directly requested,
	// for example if only a prefix of a given traversal is used to resolve
	// a complex object value. In that case, two different request traversals
	// may refer to different parts of such an object.
	//
	// An implementation of RefValues must behave as a pure function for the
	// lifetime of the object that is implementing it.
	RefValues([]hcl.Traversal) (map[string]cty.Value, tfdiags.Diagnostics)
}
