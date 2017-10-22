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

func newContext(f Fact) *Context {
	return &Context{
		Fact:  f,
		funcs: map[string]govaluate.ExpressionFunction{},
		m:     &dummyMonitor{},
	}
}

func (ctx *Context) addFuncs(funcs map[string]Func) {
	for k, v := range funcs {
		f := v
		ctx.funcs[k] = func(args ...interface{}) (interface{}, error) {
			return f(ctx, args...)
		}
	}
}
