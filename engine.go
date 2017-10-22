package rule

import "github.com/Knetic/govaluate"

// Engine is a rule engine.
type Engine struct {
	// Funcs can be added functions which can be used by actions.
	Funcs map[string]Func

	rules []*Rule
}

// NewEngine creates a new rule engine.
func NewEngine() *Engine {
	return &Engine{
		Funcs: map[string]Func{},
		rules: make([]*Rule, 0, 64),
	}
}

// AddFuncs adds functions (set of named Func).
func (eng *Engine) AddFuncs(funcs map[string]Func) {
	for k, fn := range funcs {
		if fn == nil {
			continue
		}
		eng.Funcs[k] = fn
	}
}

// AddRule compiles a new rule and puts into root rules of Engine.
func (eng *Engine) AddRule(name, condition, action string) (*Rule, error) {
	r, err := compileRule(name, condition, action)
	if err != nil {
		return nil, err
	}
	eng.rules = append(eng.rules, r)
	return r, nil
}

// Eval evaluates rules with a fact. Evaluation can be monitored by Monitor.
// nil for m is allowed it make monitor disabled.
func (eng *Engine) Eval(fact Fact, m Monitor) {
	ctx := eng.newContext(fact, m)
	curr := make([]*Rule, len(eng.rules), len(eng.rules)+99)
	copy(curr, eng.rules)
	for {
		applied, pended := eval(curr, ctx, ctx.Fact)
		curr = curr[:0]
		for _, r := range applied {
			if len(r.childRules) > 0 {
				curr = append(curr, r.childRules...)
			}
		}
		if len(applied) == 0 && len(curr) == 0 {
			break
		}
		for _, r := range pended {
			if r.rev {
				curr = append(curr, r)
			}
		}
	}
}

func eval(rules []*Rule, ctx *Context, fact Fact) ([]*Rule, []*Rule) {
	buf := make([]*Rule, len(rules))
	applied := 0
	pended := len(rules)
	for _, r := range rules {
		ok, _ := r.eval(ctx, fact)
		if !ok {
			pended--
			buf[pended] = r
			continue
		}
		buf[applied] = r
		applied++
	}
	return buf[0:applied], buf[pended:]
}

func (eng *Engine) newContext(fact Fact, m Monitor) *Context {
	if m == nil {
		m = &dummyMonitor{}
	}
	funcs := map[string]govaluate.ExpressionFunction{}
	ctx := &Context{
		Fact:  fact.Clone(),
		funcs: funcs,
		m:     m,
	}
	ctx.addFuncs(eng.Funcs)
	ctx.addFuncs(defaultFuncs)
	return ctx
}
