package rule

import "testing"

func TestEngine_Eval(t *testing.T) {
	// types for the test.
	type callDesc struct {
		name string
		args []interface{}
	}
	type ruleAdder interface {
		AddRule(_, _, _ string) (*Rule, error)
	}

	eng := NewEngine()
	called := make([]callDesc, 0, 16)
	addRule := func(ra ruleAdder, name, cond, act string) *Rule {
		r, err := ra.AddRule(name, cond, act)
		if err != nil {
			t.Fatalf("failed to add rule: name=%q cond=%q act=%q", name, cond, act)
		}
		return r
	}
	ok := func(f Fact, exp []callDesc) {
		called = called[:0]
		eng.Eval(f, nil)
		assertEquals(t, called, exp, `unexpected actions: fact=%+v`, f)
	}

	// build rule engine to test.
	eng.AddFuncs(map[string]Func{
		"FOO": func(ctx *Context, args ...interface{}) (interface{}, error) {
			called = append(called, callDesc{name: "foo", args: args})
			return true, nil
		},
		"BAR": func(ctx *Context, args ...interface{}) (interface{}, error) {
			called = append(called, callDesc{name: "bar", args: args})
			return true, nil
		},
	})
	r1 := addRule(eng, "rule1", `sex == "male"`, `FOO("rule1 fired")`)
	addRule(r1, "rule1a", `age >= 30`, `FOO("rule1a FIRED")`)
	addRule(r1, "rule1b", `age < 30`, `FOO("rule1b FIRED")`)
	r2 := addRule(eng, "rule2", `sex == "female"`, `BAR("rule2 fired")`)
	addRule(r2, "rule2a", `age >= 30`, `BAR("rule2a FIRED")`)
	addRule(r2, "rule2b", `age < 30`, `BAR("rule2b FIRED")`)
	addRule(eng, "rule3", `marked`, `FOO("rule3 fired")`).WithRevivable(true)
	r4 := addRule(eng, "rule4", `country == "Japan"`, `BAR("rule4 fired")`)
	addRule(r4, "rule4a", `age < 30`, `PUT_FACT("marked", true)`)

	// test normal rules.
	ok(Fact{"sex": "male"}, []callDesc{
		{name: "foo", args: []interface{}{"rule1 fired"}},
	})
	ok(Fact{"sex": "female"}, []callDesc{
		{name: "bar", args: []interface{}{"rule2 fired"}},
	})

	// test child rules.
	ok(Fact{"sex": "male", "age": 32.0}, []callDesc{
		{name: "foo", args: []interface{}{"rule1 fired"}},
		{name: "foo", args: []interface{}{"rule1a FIRED"}},
	})
	ok(Fact{"sex": "male", "age": 27.0}, []callDesc{
		{name: "foo", args: []interface{}{"rule1 fired"}},
		{name: "foo", args: []interface{}{"rule1b FIRED"}},
	})
	ok(Fact{"sex": "female", "age": 32.0}, []callDesc{
		{name: "bar", args: []interface{}{"rule2 fired"}},
		{name: "bar", args: []interface{}{"rule2a FIRED"}},
	})
	ok(Fact{"sex": "female", "age": 27.0}, []callDesc{
		{name: "bar", args: []interface{}{"rule2 fired"}},
		{name: "bar", args: []interface{}{"rule2b FIRED"}},
	})

	// test revival rule.
	ok(Fact{"country":"Japan"}, []callDesc{
		{name: "bar", args: []interface{}{"rule4 fired"}},
	})
	ok(Fact{"country":"Japan", "age": 27.0}, []callDesc{
		{name: "bar", args: []interface{}{"rule4 fired"}},
		{name: "foo", args: []interface{}{"rule3 fired"}},
	})
	ok(Fact{"country":"Japan", "age": 30.0}, []callDesc{
		{name: "bar", args: []interface{}{"rule4 fired"}},
	})
}
