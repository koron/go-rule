package rule

import "testing"

func TestFact_Clone(t *testing.T) {
	a := Fact{"abc": 1, "foo": "foo"}
	b := a.Clone()
	b["abc"] = 2
	assertEquals(t, a, Fact{"abc": 1, "foo": "foo"}, "a keeps original")
	assertEquals(t, b, Fact{"abc": 2, "foo": "foo"}, "b should be modified")
}

func TestFact_put(t *testing.T) {
	// basic behavior
	a := Fact{
		"number": 1,
		"fact": Fact{
			"number": 2,
			"string": "hello",
		},
		"map": map[string]interface{}{
			"number": 3,
			"string": "world",
		},
	}
	modified := Fact{
		"number": 10,
		"append": "added1",
		"fact": Fact{
			"number": 20,
			"string": "olleh",
			"append": "added2",
		},
		"map": map[string]interface{}{
			"number": 30,
			"string": "dlrow",
			"append": "added3",
		},
	}
	a.put("number", 10)
	a.put("append", "added1")
	a.put("fact.number", 20)
	a.put("fact.string", "olleh")
	a.put("fact.append", "added2")
	a.put("map.number", 30)
	a.put("map.string", "dlrow")
	a.put("map.append", "added3")
	assertEquals(t, a, modified, "a should be modified")

	// failure
	fail := func(name, msg string) {
		err := a.put(name, nil)
		if err == nil {
			t.Errorf("put(%q) succeeded but expected failure: %s", name, msg)
			return
		}
		assertEquals(t, err.Error(), msg, "put() failed with unexpected error", name)
	}
	fail("foo.bar", "not found key: foo")
	fail("fact.abc.def", "not found key: abc")
	fail("fact.number.foo", `value for "number" is not a map`)
	fail("fact.string.foo", `value for "string" is not a map`)
	fail("map.number.foo", `value for "number" is not a map`)
	fail("map.string.foo", `value for "string" is not a map`)
}
