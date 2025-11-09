package tree

type Class struct {
	Name    string
	Methods map[string]LoxFunction
}

func NewClass(name string, methods map[string]LoxFunction) Class {
	return Class{
		Name:    name,
		Methods: methods,
	}
}

func (c *Class) findMethod(name string) (LoxFunction, bool) {
	if method, ok := c.Methods[name]; ok {
		return method, true
	}
	return LoxFunction{}, false
}

func (c Class) String() string {
	return c.Name
}

func (c Class) Call(interpreter *Interpreter, arguments []any) any {
	instance := NewInstance(c)
	initializer, ok := c.findMethod("init")
	if ok {
		initializer.bind(instance).Call(interpreter, arguments)
	}
	return instance
}

func (c Class) Arity() int {
	initializer, ok := c.findMethod("init")
	if !ok {
		return 0
	}
	return initializer.Arity()
}
