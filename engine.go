package rule

import "github.com/Knetic/govaluate"

// Func represents a function which can be used by actions.
type Func func(*Context, ...interface{}) (interface{}, error)

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
		rules: make([]*Rule, 0, 99),
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
func (eng *Engine) Eval(fact Fact, m Monitor) error {
	ctx := eng.newContext(fact, m)
	curr := make([]*Rule, len(eng.rules), len(eng.rules)+99)
	copy(curr, eng.rules)
	for {
		applied, pended, err := eval(curr, ctx, fact)
		if err != nil {
			return err
		}
		curr = curr[:]
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
	return nil
}

func eval(rules []*Rule, ctx *Context, fact Fact) ([]*Rule, []*Rule, error) {
	buf := make([]*Rule, len(rules))
	applied := 0
	pended := len(rules)
	for _, r := range rules {
		ok, _, err := r.Eval(ctx, fact)
		if err != nil {
			return nil, nil, err
		}
		if !ok {
			pended--
			buf[pended] = r
		}
		buf[applied] = r
		applied++
	}
	return buf[0:applied], buf[pended:], nil
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
	for k, v := range eng.Funcs {
		funcs[k] = func(args ...interface{}) (interface{}, error) {
			return v(ctx, args...)
		}
	}
	return ctx
}
