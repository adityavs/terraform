package lang

import (
	"sync"

	"github.com/zclconf/go-cty/cty/function"
)

// Context is the main type in this package, allowing dynamic evaluation of
// blocks and expressions based on some contextual information that informs
// which variables and functions will be available.
type Context struct {
	// Scope is used to resolve any variables referenced in expressions.
	Scope Scope

	// BaseDir is the base directory used by any interpolation functions that
	// accept filesystem paths as arguments.
	BaseDir string

	// PureOnly can be set to true to request that any non-pure functions
	// produce unknown value results rather than actually executing. This is
	// important during a plan phase to avoid generating results that could
	// then differ during apply.
	PureOnly bool

	funcs     map[string]function.Function
	funcsLock sync.Mutex
}
