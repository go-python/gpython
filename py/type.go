// Type objects - these make objects

package py

import (
	"fmt"
)

type Type struct {
	Name    string     // For printing, in format "<module>.<name>"
	Doc     string     // Documentation string
	Methods StringDict // *PyMethodDef
	Members StringDict // *PyMemberDef
	//	Getset     *PyGetSetDef
	Base       *Type
	Dict       Object
	Dictoffset int
	Bases      Tuple
	Mro        Tuple // method resolution order
	//	Cache      Object
	Subclasses Tuple
	Weaklist   Tuple

	/*
	   Py_ssize_t tp_basicsize, tp_itemsize; // For allocation

	   // Methods to implement standard operations

	   destructor tp_dealloc;
	   printfunc tp_print;
	   getattrfunc tp_getattr;
	   setattrfunc tp_setattr;
	   void *tp_reserved; // formerly known as tp_compare
	   reprfunc tp_repr;

	   // Method suites for standard classes

	   PyNumberMethods *tp_as_number;
	   PySequenceMethods *tp_as_sequence;
	   PyMappingMethods *tp_as_mapping;

	   // More standard operations (here for binary compatibility)

	   hashfunc tp_hash;
	   ternaryfunc tp_call;
	   reprfunc tp_str;
	   getattrofunc tp_getattro;
	   setattrofunc tp_setattro;

	   // Functions to access object as input/output buffer
	   PyBufferProcs *tp_as_buffer;

	   // Flags to define presence of optional/expanded features
	   unsigned long tp_flags;

	   const char *tp_doc; // Documentation string

	   // Assigned meaning in release 2.0
	   // call function for all accessible objects
	   traverseproc tp_traverse;

	   // delete references to contained objects
	   inquiry tp_clear;

	   // Assigned meaning in release 2.1
	   // rich comparisons
	   richcmpfunc tp_richcompare;

	   // weak reference enabler
	   Py_ssize_t tp_weaklistoffset;

	   // Iterators
	   getiterfunc tp_iter;
	   iternextfunc tp_iternext;

	   // Attribute descriptor and subclassing stuff
	   struct PyMethodDef *tp_methods;
	   struct PyMemberDef *tp_members;
	   struct PyGetSetDef *tp_getset;
	   struct _typeobject *tp_base;
	   PyObject *tp_dict;
	   descrgetfunc tp_descr_get;
	   descrsetfunc tp_descr_set;
	   Py_ssize_t tp_dictoffset;
	   initproc tp_init;
	   allocfunc tp_alloc;
	   newfunc tp_new;
	   freefunc tp_free; // Low-level free-memory routine
	   inquiry tp_is_gc; // For PyObject_IS_GC
	   PyObject *tp_bases;
	   PyObject *tp_mro; // method resolution order
	   PyObject *tp_cache;
	   PyObject *tp_subclasses;
	   PyObject *tp_weaklist;
	   destructor tp_del;

	   // Type attribute cache version tag. Added in version 2.6
	   unsigned int tp_version_tag;

	   destructor tp_finalize;
	*/
}

var TypeType = NewType("type", "type(object) -> the object's type\ntype(name, bases, dict) -> a new type")
var BaseObjectType = NewType("object", "The most base type")

// Type of this object
func (o *Type) Type() *Type {
	return TypeType
}

// Make a new type from a name
func NewType(Name string, Doc string) *Type {
	return &Type{
		Name: Name,
		Doc:  Doc,
	}
}

// Determine the most derived metatype.
func (metatype *Type) CalculateMetaclass(bases Tuple) *Type {
	// Determine the proper metatype to deal with this,
	// and check for metatype conflicts while we're at it.
	// Note that if some other metatype wins to contract,
	// it's possible that its instances are not types. */

	winner := metatype
	for _, tmp := range bases {
		tmptype := tmp.Type()
		if winner.IsSubtype(tmptype) {
			continue
		}
		if tmptype.IsSubtype(winner) {
			winner = tmptype
			continue
		}
		// else:
		// FIXME TypeError
		panic(fmt.Sprintf("TypeError: metaclass conflict: the metaclass of a derived class must be a (non-strict) subclass of the metaclasses of all its bases"))
	}
	return winner
}

// type test with subclassing support
func (a *Type) IsSubtype(b *Type) bool {
	mro := a.Mro
	if mro != nil {
		// Deal with multiple inheritance without recursion
		// by walking the MRO tuple
		for _, base := range mro {
			if base == b {
				return true
			}
		}
		return false
	} else {
		// a is not completely initilized yet; follow tp_base
		for {
			if a == b {
				return true
			}
			a = a.Base
			if a == nil {
				break
			}
		}
		return b == BaseObjectType
	}
}

// call a type
func (t *Type) M__call__(args Tuple, kwargs StringDict) Object {
	fmt.Printf("Type __call__ FIXME\n")
	return None
}

// Make sure it satisfies the interface
var _ Object = (*Type)(nil)
var _ I__call__ = (*Type)(nil)
