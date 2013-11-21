// Cell object
//
// In the Go implementation this is just a pointer to an Object which
// can be nil

package py

// A python Cell object
type Cell struct {
	obj *Object
}

var CellType = NewType("cell", "cell object")

// Type of this object
func (o *Cell) Type() *Type {
	return CellType
}

// Define a new cell
func NewCell(obj Object) *Cell {
	return &Cell{&obj}
}

// Fetch the contents of the Cell or nil if not set
func (c *Cell) Get() Object {
	if c.obj == nil {
		return nil
	}
	return *c.obj
}

// Set the contents of the Cell
func (c *Cell) Set(obj Object) {
	c.obj = &obj
}
