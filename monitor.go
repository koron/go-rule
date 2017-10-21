package rule

// Monitor provides methods to monitoring evaluation.
type Monitor interface {
	ConditionError(*Context, *Rule, error)
	ConditionResult(*Context, *Rule, bool)
	ActionIgnore(*Context, *Rule)
	ActionCompileError(*Context, *Rule, error)
	ActionError(*Context, *Rule, error)
	ActionResult(*Context, *Rule, interface{})
}

type dummyMonitor struct {}
func (*dummyMonitor) ConditionError(*Context, *Rule, error) {}
func (*dummyMonitor) ConditionResult(*Context, *Rule, bool) {}
func (*dummyMonitor) ActionIgnore(*Context, *Rule) {}
func (*dummyMonitor) ActionCompileError(*Context, *Rule, error) {}
func (*dummyMonitor) ActionError(*Context, *Rule, error) {}
func (*dummyMonitor) ActionResult(*Context, *Rule, interface{}) {}
