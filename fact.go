package rule

import (
	"fmt"
	"strings"
)

// Fact is evaluation target context.
type Fact map[string]interface{}

// Clone clones a fact except array members.
func (f Fact) Clone() Fact {
	dst := Fact{}
	for k, v := range f {
		switch w := v.(type) {
		case Fact:
			dst[k] = w.Clone()
		default:
			dst[k] = v
		}
	}
	return dst
}

func (f Fact) put(name string, value interface{}) error {
	n := strings.Index(name, ".")
	if n < 0 {
		f[name] = value
		return nil
	}
	k1 := name[:n]
	v, ok := f[k1]
	if !ok {
		return fmt.Errorf("not found key: %s", k1)
	}
	switch w := v.(type) {
	case map[string]interface{}:
		return Fact.put(w, name[n+1:], value)
	case Fact:
		return Fact.put(w, name[n+1:], value)
	}
	return fmt.Errorf("value for %q is not a map", k1)
}
