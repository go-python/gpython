// Property object

package py

// A python Property object
type Property struct {
	Fget func(self Object) Object
	Fset func(self, value Object)
	Fdel func(self Object)
	Doc  string
}

var PropertyType = NewType("property", "property object")

// Type of this object
func (o *Property) Type() *Type {
	return PropertyType
}

func (p *Property) M__get__(instance, owner Object) Object {
	if p.Fget == nil {
		panic(ExceptionNewf(AttributeError, "can't get attribute"))
	}
	return p.Fget(instance)
}

func (p *Property) M__set__(instance, value Object) Object {
	if p.Fset == nil {
		panic(ExceptionNewf(AttributeError, "can't set attribute"))
	}
	p.Fset(instance, value)
	return None
}

func (p *Property) M__delete__(instance Object) Object {
	if p.Fdel == nil {
		panic(ExceptionNewf(AttributeError, "can't delete attribute"))
	}
	p.Fdel(instance)
	return None
}

// Interfaces
var _ I__get__ = (*Property)(nil)
var _ I__set__ = (*Property)(nil)
var _ I__delete__ = (*Property)(nil)
