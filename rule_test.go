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

func TestRule_Eval(t *testing.T) {
	// TODO: test Rule.Eval
}
