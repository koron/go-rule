package rule

import (
	"math"

	"github.com/Knetic/govaluate"
)

// Rule represents a rule.
type Rule struct {
	name  string // rule name
	rev   bool   // revivable flag: re-evaluate when other rules are applied.
	act   string // action statement
	acErr error  // error when compiling the action

	// compiled expression for condition
	cond *govaluate.EvaluableExpression

	// childRules holds depended rules. Depended rules are evaluated when the
	// parent is applied.
	childRules []*Rule
}

// Name returns name of rule.
func (r *Rule) Name() string {
	return r.name
}

// WithName changes name of rule.
func (r *Rule) WithName(s string) *Rule {
	r.name = s
	return r
}

// Revivable returns revivable flag. Revivable flag means whether evaluate
// again when other rules are applied or not.
func (r *Rule) Revivable() bool {
	return r.rev
}

// WithRevivable changes revivable flag.
func (r *Rule) WithRevivable(v bool) *Rule {
	r.rev = v
	return r
}

func compileRule(name, cond, act string) (*Rule, error) {
	c, err := govaluate.NewEvaluableExpression(cond)
	if err != nil {
		return nil, err
	}
	r := &Rule{
		name: name,
		cond: c,
		act:  act,
	}
	return r, nil
}

func isTrue(v interface{}) bool {
	switch w := v.(type) {
	case bool:
		return w
	case float64:
		if math.IsNaN(w) {
			return false
		}
		return w != 0.0
	case string:
		return len(w) > 0
	default:
		return false
	}
}

// Eval evaluates a fact.  If condition is true then evaluate action.
func (r *Rule) eval(ctx *Context, fact Fact) (bool, interface{}) {
	res, err := r.cond.Evaluate(fact)
	if err != nil {
		ctx.m.ConditionError(ctx, r, err)
		return false, nil
	}
	resTrue := isTrue(res)
	ctx.m.ConditionResult(ctx, r, resTrue)
	if !resTrue {
		return false, nil
	}
	res2, err := r.doAct(ctx, fact)
	if err != nil {
		return true, nil
	}
	return true, res2
}

func (r *Rule) doAct(ctx *Context, fact Fact) (interface{}, error) {
	if r.act == "" || r.acErr != nil {
		ctx.m.ActionIgnore(ctx, r)
		return nil, nil
	}
	ex, err := govaluate.NewEvaluableExpressionWithFunctions(r.act, ctx.funcs)
	if err != nil {
		r.acErr = err
		ctx.m.ActionCompileError(ctx, r, err)
		return nil, nil
	}
	res, err := ex.Evaluate(fact)
	if err != nil {
		ctx.m.ActionError(ctx, r, err)
		return nil, err
	}
	ctx.m.ActionResult(ctx, r, res)
	return res, nil
}

// AddRule compiles a rule and puts into child rules.
func (r *Rule) AddRule(name, condition, action string) (*Rule, error) {
	cr, err := compileRule(name, condition, action)
	if err != nil {
		return nil, err
	}
	r.childRules = append(r.childRules, cr)
	return cr, nil
}
