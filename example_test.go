package rule_test

import (
	"fmt"

	rule "github.com/koron/go-rule"
)

func ExampleEngine() {
	eng := rule.NewEngine()
	eng.AddFuncs(map[string]rule.Func{
		"PRINT": func(ctx *rule.Context, args ...interface{}) (interface{}, error) {
			fmt.Println(args...)
			return true, nil
		},
	})
	eng.AddRule("too cold", `temp < 15`, `PRINT("COLD!")`)
	eng.AddRule("too hot", `temp >= 25`, `PRINT("HOT!")`)
	eng.AddRule("OK", `temp >= 15 && temp < 25`, `PRINT("OK")`)

	// Output:
	// COLD!
	// HOT!
	// OK
	eng.Eval(rule.Fact{"temp": 10.0}, nil)
	eng.Eval(rule.Fact{"temp": 30.0}, nil)
	eng.Eval(rule.Fact{"temp": 20.0}, nil)
}
