package rule

import "github.com/Knetic/govaluate"

// Context is an evaluation context.
type Context struct {
	// Fact hols current Fact in evaluation. Note that any modifications on
	// this affects later evaluation.
	Fact Fact

	funcs map[string]govaluate.ExpressionFunction
	m     Monitor
}
