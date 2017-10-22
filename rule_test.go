package rule

import (
	"math"
	"testing"
)

func TestIsTrue(t *testing.T) {
	ok := func(v interface{}, exp bool) {
		act := isTrue(v)
		if act != exp {
			t.Errorf("isTrue(%+v) returns %t but expected %t", v, act, exp)
		}
	}
	ok(true, true)
	ok(false, false)
	ok("", false)
	ok("foo", true)
	ok("bar", true)
	ok("false", true)
	ok(0.0, false)
	ok(0.1, true)
	ok(-0.1, true)
	ok(math.NaN(), false)
}

type ruleDesc struct {
	name string
	cond string
	act  string
}

type ruleResult struct {
	cond bool
	act  interface{}
}

func TestRule_eval(t *testing.T) {
	fact := Fact{
		"country": "Japan",
		"sex":     "male",
		"age":     27.0,
	}
	ok := func(rd ruleDesc, rr ruleResult) {
		r, err := compileRule(rd.name, rd.cond, rd.act)
		if err != nil {
			t.Errorf(`compileRule(%v) failed: %s`, rd, err)
			return
		}
		f := fact.Clone()
		ctx := newContext(f)
		ctx.addFuncs(defaultFuncs)
		rc, ra := r.eval(ctx, f)
		assertEquals(t, ruleResult{cond: rc, act: ra}, rr,
			"evaluation failure: %+v", *r)
	}
	ok(ruleDesc{cond: "age < 30", act: "true"},
		ruleResult{cond: true, act: true})
	ok(ruleDesc{cond: "age > 30", act: "true"},
		ruleResult{cond: false, act: nil})
	ok(ruleDesc{cond: "age < 30", act: "10 * 99"},
		ruleResult{cond: true, act: 990.0})
	ok(ruleDesc{cond: "age < 30", act: "country"},
		ruleResult{cond: true, act: "Japan"})
	ok(ruleDesc{cond: "age < 30", act: `PUT_FACT("country.foo", 123)`},
		ruleResult{cond: true, act: nil})
}
