package rule

// Func represents a function which can be used by actions.
type Func func(*Context, ...interface{}) (interface{}, error)

var defaultFuncs = map[string]Func{
	"PUT_FACT": putFact,
}

func putFact(ctx *Context, args ...interface{}) (interface{}, error) {
	if len(args) != 2 {
		return false, &ArgsCountError{Expected: 2, Given: len(args)}
	}
	name, ok := args[0].(string)
	if !ok {
		return false, &ArgTypeError{Pos: 0, Expected: "string", Given: args[0]}
	}
	err := ctx.Fact.put(name, args[1])
	if err != nil {
		return false, &NativeFuncError{Name: "Fact.put", Err: err}
	}
	return true, nil
}
