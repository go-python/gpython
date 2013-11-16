// Module objects

package py

// A python Module object
type Module struct {
	name    string
	doc     string
	methods []*Method
	//	dict Dict
}

var ModuleType = NewType("module")

// Type of this object
func (o *Module) Type() *Type {
	return ModuleType
}

// Define a new module
func NewModule(name, doc string, methods []*Method) *Module {
	return &Module{
		name:    name,
		doc:     doc,
		methods: methods,
	}
}
