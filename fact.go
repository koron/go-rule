package rule

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
