package rule

import (
	"fmt"
	"testing"
)

func TestPutFunc_failure(t *testing.T) {
	ng := func(args []interface{}, exp error) {
		ctx := &Context{Fact: Fact{}}
		_, err := putFact(ctx, args...)
		if err == nil {
			t.Fatalf(`putFact(%+v) should be failed but suceeded: %#v`, args, exp)
		}
		assertEquals(t, err, exp, `putFact(%+v) should be failed`, args)
	}
	ng(nil, &ArgsCountError{Expected: 2, Given: 0})
	ng([]interface{}{}, &ArgsCountError{Expected: 2, Given: 0})
	ng([]interface{}{"foo"}, &ArgsCountError{Expected: 2, Given: 1})
	ng([]interface{}{"foo", "bar", "baz"},
		&ArgsCountError{Expected: 2, Given: 3})
	ng([]interface{}{123, "foo"},
		&ArgTypeError{Pos: 0, Expected: "string", Given: 123})
	ng([]interface{}{"foo.bar", 123},
		&NativeFuncError{
			Name: "Fact.put",
			Err:  fmt.Errorf("not found key: %s", "foo"),
		})
}
