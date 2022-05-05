// Copyright 2018 The go-python Authors.  All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Type objects - these make objects

// FIXME should be caching the expensive lookups in the superclasses
// and using the cache clearing machinery to clear the caches when the
// heirachy changes

// FIXME should make Mro and Bases be []*Type

package py

import (
	"fmt"
	"log"
)

// Type flags (tp_flags)
//
// These flags are used to extend the type structure in a backwards-compatible
// fashion. Extensions can use the flags to indicate (and test) when a given
// type structure contains a new feature. The Python core will use these when
// introducing new functionality between major revisions (to avoid mid-version
// changes in the PYTHON_API_VERSION).
//
// Arbitration of the flag bit positions will need to be coordinated among
// all extension writers who publically release their extensions (this will
// be fewer than you might expect!)..
//
// Most flags were removed as of Python 3.0 to make room for new flags.  (Some
// flags are not for backwards compatibility but to indicate the presence of an
// optional feature; these flags remain of course.)
//
// Type definitions should use TPFLAGS_DEFAULT for their tp_flags value.
//
// Code can use PyType_HasFeature(type_ob, flag_value) to test whether the
// given type object has a specified feature.

const (
	// Set if the type object is dynamically allocated
	TPFLAGS_HEAPTYPE uint = 1 << 9

	// Set if the type allows subclassing
	TPFLAGS_BASETYPE uint = 1 << 10

	// Set if the type is 'ready' -- fully initialized
	TPFLAGS_READY uint = 1 << 12

	// Set while the type is being 'readied', to prevent recursive ready calls
	TPFLAGS_READYING uint = 1 << 13

	// Objects support garbage collection (see objimp.h)
	TPFLAGS_HAVE_GC uint = 1 << 14

	// Objects support type attribute cache
	TPFLAGS_HAVE_VERSION_TAG  uint = 1 << 18
	TPFLAGS_VALID_VERSION_TAG uint = 1 << 19

	// Type is abstract and cannot be instantiated
	TPFLAGS_IS_ABSTRACT uint = 1 << 20

	// These flags are used to determine if a type is a subclass.
	TPFLAGS_INT_SUBCLASS      uint = 1 << 23
	TPFLAGS_LONG_SUBCLASS     uint = 1 << 24
	TPFLAGS_LIST_SUBCLASS     uint = 1 << 25
	TPFLAGS_TUPLE_SUBCLASS    uint = 1 << 26
	TPFLAGS_BYTES_SUBCLASS    uint = 1 << 27
	TPFLAGS_UNICODE_SUBCLASS  uint = 1 << 28
	TPFLAGS_DICT_SUBCLASS     uint = 1 << 29
	TPFLAGS_BASE_EXC_SUBCLASS uint = 1 << 30
	TPFLAGS_TYPE_SUBCLASS     uint = 1 << 31

	TPFLAGS_DEFAULT = TPFLAGS_HAVE_VERSION_TAG
)

type NewFunc func(metatype *Type, args Tuple, kwargs StringDict) (Object, error)

type InitFunc func(self Object, args Tuple, kwargs StringDict) error

type Type struct {
	ObjectType *Type  // Type of this object -- FIXME this is redundant in Base?
	Name       string // For printing, in format "<module>.<name>"
	Doc        string // Documentation string
	//	Methods    StringDict // *PyMethodDef
	//	Members    StringDict // *PyMemberDef
	//	Getset     *PyGetSetDef
	Base *Type
	Dict StringDict
	//	Dictoffset int
	Bases Tuple
	Mro   Tuple // method resolution order
	//	Cache      Object
	//	Subclasses Tuple
	//	Weaklist   Tuple
	New      NewFunc
	Init     InitFunc
	Flags    uint // Flags to define presence of optional/expanded features
	Qualname string

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

var TypeType *Type = &Type{
	Name: "type",
	Doc:  "type(object) -> the object's type\ntype(name, bases, dict) -> a new type",
	Dict: StringDict{},
}

var ObjectType = &Type{
	Name:  "object",
	Doc:   "The most base type",
	Flags: TPFLAGS_BASETYPE,
	Dict:  StringDict{},
}

func init() {
	// Initialised like this to avoid initialisation loops
	TypeType.New = TypeNew
	TypeType.Init = TypeInit
	TypeType.ObjectType = TypeType
	ObjectType.New = ObjectNew
	ObjectType.Init = ObjectInit
	ObjectType.ObjectType = TypeType
	err := TypeType.Ready()
	if err != nil {
		log.Fatal(err)
	}
	err = ObjectType.Ready()
	if err != nil {
		log.Fatal(err)
	}
}

// Type of this object
func (t *Type) Type() *Type {
	return t.ObjectType
}

// Satistfy error interface
func (t *Type) Error() string {
	return t.Name
}

// Get the Dict
func (t *Type) GetDict() StringDict {
	return t.Dict
}

// delayedReady holds types waiting to be intialised
var delayedReady = []*Type{}

// TypeDelayReady stores the list of types to initialise
//
// Call MakeReady when all initialised
func TypeDelayReady(t *Type) {
	delayedReady = append(delayedReady, t)
}

// TypeMakeReady readies all the types
func TypeMakeReady() (err error) {
	for _, t := range delayedReady {
		err = t.Ready()
		if err != nil {
			return fmt.Errorf("Error initialising go type %s: %v", t.Name, err)
		}
	}
	delayedReady = nil
	return nil
}

func init() {
	err := TypeMakeReady()
	if err != nil {
		log.Fatal(err)
	}
}

// Make a new type from a name
//
// For making Go types
func NewType(Name string, Doc string) *Type {
	t := &Type{
		ObjectType: TypeType,
		Name:       Name,
		Doc:        Doc,
		Dict:       StringDict{},
	}
	TypeDelayReady(t)
	return t
}

// Make a new type with constructors
//
// For making Go types
func NewTypeX(Name string, Doc string, New NewFunc, Init InitFunc) *Type {
	t := &Type{
		ObjectType: TypeType,
		Name:       Name,
		Doc:        Doc,
		New:        New,
		Init:       Init,
		Dict:       StringDict{},
	}
	TypeDelayReady(t)
	return t
}

// Make a subclass of a type
//
// For making Go types
func (t *Type) NewTypeFlags(Name string, Doc string, New NewFunc, Init InitFunc, Flags uint) *Type {
	// inherit constructors
	if New == nil {
		New = t.New
	}
	if Init == nil {
		Init = t.Init
	}
	// FIXME inherit more stuff
	tt := &Type{
		ObjectType: t,
		Name:       Name,
		Doc:        Doc,
		New:        New,
		Init:       Init,
		Flags:      Flags,
		Dict:       StringDict{},
		Bases:      Tuple{t},
	}
	TypeDelayReady(t)
	return tt
}

// Make a subclass of a type
//
// For making Go types
func (t *Type) NewType(Name string, Doc string, New NewFunc, Init InitFunc) *Type {
	// Inherit flags from superclass
	// FIXME not sure this is correct!
	return t.NewTypeFlags(Name, Doc, New, Init, t.Flags)
}

// Determine the most derived metatype.
func (metatype *Type) CalculateMetaclass(bases Tuple) (*Type, error) {
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
		return nil, ExceptionNewf(TypeError, "metaclass conflict: the metaclass of a derived class must be a (non-strict) subclass of the metaclasses of all its bases")
	}
	return winner, nil
}

// type test with subclassing support
// reads a IsSubtype of b
func (a *Type) IsSubtype(b *Type) bool {
	mro := a.Mro
	if len(mro) != 0 {
		// Deal with multiple inheritance without recursion
		// by walking the MRO tuple
		for _, baseObj := range mro {
			base := baseObj.(*Type)
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
		return b == ObjectType
	}
}

// Call type()
func (t *Type) M__call__(args Tuple, kwargs StringDict) (Object, error) {
	if t.New == nil {
		return nil, ExceptionNewf(TypeError, "cannot create '%s' instances", t.Name)
	}

	obj, err := t.New(t, args, kwargs)
	if err != nil {
		return nil, err
	}
	// Ugly exception: when the call was type(something),
	// don't call tp_init on the result.
	if t == TypeType && len(args) == 1 && len(kwargs) == 0 {
		return obj, nil
	}
	// If the returned object is not an instance of type,
	// it won't be initialized.
	if !obj.Type().IsSubtype(t) {
		return obj, nil
	}
	objType := obj.Type()
	if objType.Init != nil {
		err = objType.Init(obj, args, kwargs)
		if err != nil {
			return nil, err
		}
	}
	return obj, nil
}

// Internal API to look for a name through the MRO.
// This returns a borrowed reference, and doesn't set an exception,
// returning nil instead
func (t *Type) Lookup(name string) Object {
	// Py_ssize_t i, n;
	// PyObject *mro, *res, *base, *dict;
	// unsigned int h;

	// FIXME caching
	// if (MCACHE_CACHEABLE_NAME(name) &&
	//     PyType_HasFeature(type, TPFLAGS_VALID_VERSION_TAG)) {
	//     // fast path
	//     h = MCACHE_HASH_METHOD(type, name);
	//     if (method_cache[h].version == type->tp_version_tag &&
	//         method_cache[h].name == name)
	//         return method_cache[h].value;
	// }

	// Look in tp_dict of types in MRO
	mro := t.Mro

	// If mro is nil, the type is either not yet initialized
	// by PyType_Ready(), or already cleared by type_clear().
	// Either way the safest thing to do is to return nil.
	if mro == nil {
		return nil
	}

	var res Object
	// keep a strong reference to mro because type->tp_mro can be replaced
	// during PyDict_GetItem(dict, name)
	for _, baseObj := range mro {
		base := baseObj.(*Type)
		var ok bool
		res, ok = base.Dict[name]
		if ok {
			break
		}
	}

	// FIXME caching
	// if (MCACHE_CACHEABLE_NAME(name) && assign_version_tag(type)) {
	//     h = MCACHE_HASH_METHOD(type, name);
	//     method_cache[h].version = type->tp_version_tag;
	//     method_cache[h].value = res;  /* borrowed */
	//     Py_INCREF(name);
	//     Py_DECREF(method_cache[h].name);
	//     method_cache[h].name = name;
	// }

	return res
}

// Get an attribute from the type of a go type
//
// Doesn't call __getattr__ etc
//
// # Returns nil if not found
//
// Doesn't look in the instance dictionary
//
// FIXME this isn't totally correct!
// as we are ignoring getattribute etc
// See _PyObject_GenericGetAttrWithDict in object.c
func (t *Type) NativeGetAttrOrNil(name string) Object {
	// Look in type Dict
	if res, ok := t.Dict[name]; ok {
		return res
	}
	// Now look through base classes etc
	return t.Lookup(name)
}

// Get an attribute from the type
//
// Doesn't call __getattr__ etc
//
// # Returns nil if not found
//
// FIXME this isn't totally correct!
// as we are ignoring getattribute etc
// See _PyObject_GenericGetAttrWithDict in object.c
func (t *Type) GetAttrOrNil(name string) Object {
	// Look in instance dictionary first
	if res, ok := t.Dict[name]; ok {
		return res
	}
	// Then look in type Dict
	if res, ok := t.Type().Dict[name]; ok {
		return res
	}
	// Now look through base classes etc
	return t.Lookup(name)
}

// Calls method on name
//
// If method not found returns (nil, false, nil)
//
// If method found returns (object, true, err)
//
// May raise exceptions if calling the method failed
func (t *Type) CallMethod(name string, args Tuple, kwargs StringDict) (Object, bool, error) {
	fn := t.GetAttrOrNil(name) // FIXME this should use py.GetAttrOrNil?
	if fn == nil {
		return nil, false, nil
	}
	res, err := Call(fn, args, kwargs)
	return res, true, err
}

// Calls a type method on obj
//
// If obj isnt a *Type or the method isn't found on it returns (nil, false, nil)
//
// Otherwise returns (object, true, err)
//
// May raise exceptions if calling the method fails
func TypeCall(self Object, name string, args Tuple, kwargs StringDict) (Object, bool, error) {
	t, ok := self.(*Type)
	if !ok {
		return nil, false, nil
	}
	return t.CallMethod(name, args, kwargs)
}

// Calls TypeCall with 0 arguments
func TypeCall0(self Object, name string) (Object, bool, error) {
	return TypeCall(self, name, Tuple{self}, nil)
}

// Calls TypeCall with 1 argument
func TypeCall1(self Object, name string, arg Object) (Object, bool, error) {
	return TypeCall(self, name, Tuple{self, arg}, nil)
}

// Calls TypeCall with 2 arguments
func TypeCall2(self Object, name string, arg1, arg2 Object) (Object, bool, error) {
	return TypeCall(self, name, Tuple{self, arg1, arg2}, nil)
}

// Internal routines to do a method lookup in the type
// without looking in the instance dictionary
// (so we can't use PyObject_GetAttr) but still binding
// it to the instance.  The arguments are the object,
// the method name as a C string, and the address of a
// static variable used to cache the interned Python string.
//
// Two variants:
//
// - lookup_maybe() returns nil without raising an exception
//
//	when the _PyType_Lookup() call fails;
//
// - lookup_method() always raises an exception upon errors.
func lookup_maybe(self Object, attr string) Object {
	res := self.Type().Lookup(attr)
	// FIXME descriptor lookup
	// if (res != nil) {
	// 	descrgetfunc f;
	// 	if ((f = Py_TYPE(res)->tp_descr_get) == nil) {
	// 		Py_INCREF(res);
	// 	}else{
	// 		res = f(res, self, (PyObject *)(Py_TYPE(self)));
	// 	}
	// }
	return res
}

// func lookup_method(self Object, attr string) Object {
// 	res := lookup_maybe(self, attr)
// 	if res == nil {
// 		// FIXME PyErr_SetObject(PyExc_AttributeError, attrid->object);
// 		return ExceptionNewf(AttributeError, "'%s' object has no attribute '%s'", self.Type().Name, attr)
// 	}
// 	return res
// }

// Method resolution order algorithm C3 described in
// "A Monotonic Superclass Linearization for Dylan",
// by Kim Barrett, Bob Cassel, Paul Haahr,
// David A. Moon, Keith Playford, and P. Tucker Withington.
// (OOPSLA 1996)
//
// Some notes about the rules implied by C3:
//
// No duplicate bases.
// It isn't legal to repeat a class in a list of base classes.
//
// The next three properties are the 3 constraints in "C3".
//
// Local precendece order.
// If A precedes B in C's MRO, then A will precede B in the MRO of all
// subclasses of C.
//
// Monotonicity.
// The MRO of a class must be an extension without reordering of the
// MRO of each of its superclasses.
//
// Extended Precedence Graph (EPG).
// Linearization is consistent if there is a path in the EPG from
// each class to all its successors in the linearization.  See
// the paper for definition of EPG.

func tail_contains(list *List, whence int, o Object) bool {
	for j := whence + 1; j < len(list.Items); j++ {
		if list.Items[j] == o {
			return true
		}
	}
	return false
}

func class_name(cls Object) string {
	name := ObjectGetAttr(cls, "__name__")
	if name == nil {
		name = ObjectRepr(cls)
	}
	nameString, ok := name.(String)
	if !ok {
		return ""
	}
	return string(nameString)
}

func check_duplicates(list *List) error {
	// Let's use a quadratic time algorithm,
	// assuming that the bases lists is short.
	for i := range list.Items {
		o := list.Items[i]
		for j := i + 1; j < len(list.Items); j++ {
			if list.Items[j] == o {
				return ExceptionNewf(TypeError, "duplicate base class %s", class_name(o))
			}
		}
	}
	return nil
}

// Raise a TypeError for an MRO order disagreement.
//
// It's hard to produce a good error message.  In the absence of better
// insight into error reporting, report the classes that were candidates
// to be put next into the MRO.  There is some conflict between the
// order in which they should be put in the MRO, but it's hard to
// diagnose what constraint can't be satisfied.
func set_mro_error(to_merge *List, remain []int) error {
	return ExceptionNewf(TypeError, "mro is wonky")
	/* FIXME implement this!
	       Py_ssize_t i, n, off, to_merge_size;
	       char buf[1000];
	       PyObject *k, *v;
	       PyObject *set = PyDict_New();
	       if (!set) return;

	       to_merge_size = PyList_GET_SIZE(to_merge);
	       for (i = 0; i < to_merge_size; i++) {
	           PyObject *L = PyList_GET_ITEM(to_merge, i);
	           if (remain[i] < PyList_GET_SIZE(L)) {
	               PyObject *c = PyList_GET_ITEM(L, remain[i]);
	               if (PyDict_SetItem(set, c, Py_None) < 0) {
	                   Py_DECREF(set);
	                   return;
	               }
	           }
	       }
	       n = PyDict_Size(set);

	       off = PyOS_snprintf(buf, sizeof(buf), "Cannot create a \
	   consistent method resolution\norder (MRO) for bases");
	       i = 0;
	       while (PyDict_Next(set, &i, &k, &v) && (size_t)off < sizeof(buf)) {
	           PyObject *name = class_name(k);
	           char *name_str;
	           if (name != nil) {
	               name_str = _PyUnicode_AsString(name);
	               if (name_str == nil)
	                   name_str = "?";
	           } else
	               name_str = "?";
	           off += PyOS_snprintf(buf + off, sizeof(buf) - off, " %s", name_str);
	           Py_XDECREF(name);
	           if (--n && (size_t)(off+1) < sizeof(buf)) {
	               buf[off++] = ',';
	               buf[off] = '\0';
	           }
	       }
	       PyErr_SetString(PyExc_TypeError, buf);
	       Py_DECREF(set);
	*/
}

func pmerge(acc, to_merge *List) error {
	// Py_ssize_t i, j, to_merge_size, empty_cnt;
	// int *remain;
	// int ok;

	to_merge_size := len(to_merge.Items)

	// remain stores an index into each sublist of to_merge.
	// remain[i] is the index of the next base in to_merge[i]
	// that is not included in acc.
	remain := make([]int, to_merge_size)

again:
	empty_cnt := 0
	for i := 0; i < to_merge_size; i++ {
		cur_list := to_merge.Items[i].(*List)

		if remain[i] >= len(cur_list.Items) {
			empty_cnt++
			continue
		}

		// Choose next candidate for MRO.
		//
		// The input sequences alone can determine the choice.
		// If not, choose the class which appears in the MRO
		// of the earliest direct superclass of the new class.

		candidate := cur_list.Items[remain[i]]
		for j := 0; j < to_merge_size; j++ {
			j_lst := to_merge.Items[j].(*List)
			if tail_contains(j_lst, remain[j], candidate) {
				goto skip // continue outer loop
			}
		}
		acc.Append(candidate)
		for j := 0; j < to_merge_size; j++ {
			j_lst := to_merge.Items[j].(*List)
			if remain[j] < len(j_lst.Items) && j_lst.Items[remain[j]] == candidate {
				remain[j]++
			}
		}
		goto again
	skip:
	}

	if empty_cnt == to_merge_size {
		return nil
	}
	return set_mro_error(to_merge, remain)
}

func (t *Type) mro_implementation() (Object, error) {
	// Py_ssize_t i, n;
	// int ok;
	// PyObject *bases, *result;
	// PyObject *to_merge, *bases_aslist;
	var err error

	if t.Dict == nil {
		err = t.Ready()
		if err != nil {
			return nil, err
		}
	}

	// Find a superclass linearization that honors the constraints
	// of the explicit lists of bases and the constraints implied by
	// each base class.
	//
	// to_merge is a list of lists, where each list is a superclass
	// linearization implied by a base class.  The last element of
	// to_merge is the declared list of bases.

	bases := t.Bases
	n := len(bases)
	to_merge := NewListSized(n + 1)

	for i := range bases {
		base := bases[i].(*Type)
		parentMRO, err := SequenceList(base.Mro)
		if err != nil {
			return nil, err
		}
		to_merge.Items[i] = parentMRO
	}

	bases_aslist, err := SequenceList(bases)
	if err != nil {
		return nil, err
	}

	// This is just a basic sanity check.
	err = check_duplicates(bases_aslist)
	if err != nil {
		return nil, err
	}

	to_merge.Items[n] = bases_aslist

	result := NewListFromItems([]Object{t})

	err = pmerge(result, to_merge)
	if err != nil {
		return nil, err
	}

	return result, nil
}

func (t *Type) mro_internal() (err error) {
	// PyObject *mro, *result, *tuple;
	var result Object
	checkit := false

	if t == TypeType {
		result, err = t.mro_implementation()
		if err != nil {
			return err
		}
	} else {
		checkit = true
		// FIXME this is what it was originally
		// but we haven't put mro in slots or anything
		// mro := lookup_method(t, "mro")
		mro := lookup_maybe(t, "mro")
		if mro == nil {
			// Default to internal implementation
			result, err = t.mro_implementation()
			if err != nil {
				return err
			}
		} else {
			result, err = Call(mro, nil, nil)
			if err != nil {
				return err
			}
		}
	}
	tuple, err := SequenceTuple(result)
	if err != nil {
		return err
	}
	if checkit {
		// Py_ssize_t i, len;
		// PyObject *cls;
		// PyTypeObject *solid;

		solid := t.solid_base()

		for i := range tuple {
			cls := tuple[i]
			t, ok := cls.(*Type)
			if !ok {
				return ExceptionNewf(TypeError, "mro() returned a non-class ('%s')", cls.Type().Name)
			}
			if !solid.IsSubtype(t.solid_base()) {
				return ExceptionNewf(TypeError, "mro() returned base with unsuitable layout ('%.500s')", cls.Type().Name)
			}
		}
	}
	t.Mro = tuple

	// FIXME t.type_mro_modified(t.Mro)
	// corner case: the super class might have been hidden
	// from the custom MRO
	// FIXME t.type_mro_modified(t.Bases)

	// FIXME t.Modified()
	return nil
}

func (t *Type) inherit_special(base *Type) {
	//     /* Copying basicsize is connected to the GC flags */
	//     if (!(type->tp_flags & TPFLAGS_HAVE_GC) &&
	//         (base->tp_flags & TPFLAGS_HAVE_GC) &&
	//         (!type->tp_traverse && !type->tp_clear)) {
	//         type->tp_flags |= TPFLAGS_HAVE_GC;
	//         if (type->tp_traverse == nil)
	//             type->tp_traverse = base->tp_traverse;
	//         if (type->tp_clear == nil)
	//             type->tp_clear = base->tp_clear;
	//     }
	//     {
	//         /* The condition below could use some explanation.
	//            It appears that tp_new is not inherited for static types
	//            whose base class is 'object'; this seems to be a precaution
	//            so that old extension types don't suddenly become
	//            callable (object.__new__ wouldn't insure the invariants
	//            that the extension type's own factory function ensures).
	//            Heap types, of course, are under our control, so they do
	//            inherit tp_new; static extension types that specify some
	//            other built-in type as the default also
	//            inherit object.__new__. */
	//         if (base != &PyBaseObject_Type ||
	//             (type->tp_flags & TPFLAGS_HEAPTYPE)) {
	//             if (type->tp_new == nil)
	//                 type->tp_new = base->tp_new;
	//         }
	//     }
	//     if (type->tp_basicsize == 0)
	//         type->tp_basicsize = base->tp_basicsize;

	//     /* Copy other non-function slots */

	// #undef COPYVAL
	// #define COPYVAL(SLOT) \
	//     if (type->SLOT == 0) type->SLOT = base->SLOT

	//     COPYVAL(tp_itemsize);
	//     COPYVAL(tp_weaklistoffset);
	//     COPYVAL(tp_dictoffset);

	// Setup fast subclass flags
	switch {
	case base.IsSubtype(BaseException):
		t.Flags |= TPFLAGS_BASE_EXC_SUBCLASS
	case base.IsSubtype(TypeType):
		t.Flags |= TPFLAGS_TYPE_SUBCLASS
	case base.IsSubtype(IntType):
		t.Flags |= TPFLAGS_LONG_SUBCLASS
	case base.IsSubtype(BigIntType):
		t.Flags |= TPFLAGS_LONG_SUBCLASS
	case base.IsSubtype(BytesType):
		t.Flags |= TPFLAGS_BYTES_SUBCLASS
	case base.IsSubtype(StringType):
		t.Flags |= TPFLAGS_UNICODE_SUBCLASS
	case base.IsSubtype(TupleType):
		t.Flags |= TPFLAGS_TUPLE_SUBCLASS
	case base.IsSubtype(ListType):
		t.Flags |= TPFLAGS_LIST_SUBCLASS
	case base.IsSubtype(DictType):
		t.Flags |= TPFLAGS_DICT_SUBCLASS
	}

}

func add_subclass(base, t *Type) {
	// Py_ssize_t i;
	// int result;
	// PyObject *list, *ref, *newobj;

	// list = base->tp_subclasses;
	// if (list == nil) {
	//     base->tp_subclasses = list = PyList_New(0);
	//     if (list == nil)
	//         return -1;
	// }
	// assert(PyList_Check(list));
	// newobj = PyWeakref_NewRef((PyObject *)type, nil);
	// i = PyList_GET_SIZE(list);
	// while (--i >= 0) {
	//     ref = PyList_GET_ITEM(list, i);
	//     assert(PyWeakref_CheckRef(ref));
	//     if (PyWeakref_GET_OBJECT(ref) == Py_None)
	//         return PyList_SetItem(list, i, newobj);
	// }
	// result = PyList_Append(list, newobj);
	// Py_DECREF(newobj);
	// return result;
}

// func remove_subclass(base, t *Type) {
// 	// Py_ssize_t i;
// 	// PyObject *list, *ref;
//
// 	// list = base->tp_subclasses;
// 	// if (list == nil) {
// 	//     return;
// 	// }
// 	// assert(PyList_Check(list));
// 	// i = PyList_GET_SIZE(list);
// 	// while (--i >= 0) {
// 	//     ref = PyList_GET_ITEM(list, i);
// 	//     assert(PyWeakref_CheckRef(ref));
// 	//     if (PyWeakref_GET_OBJECT(ref) == (PyObject*)type) {
// 	//         /* this can't fail, right? */
// 	//         PySequence_DelItem(list, i);
// 	//         return;
// 	//     }
// 	// }
// }

// Ready the type for use
//
// Returns an error on problems
func (t *Type) Ready() error {
	// PyObject *dict, *bases;
	// PyTypeObject *base;
	// Py_ssize_t i, n;
	var err error

	if t.Flags&TPFLAGS_READY != 0 {
		if t.Dict == nil {
			return ExceptionNewf(SystemError, "Type.Ready is Ready but Dict is nil")
		}
		return nil
	}
	if t.Flags&TPFLAGS_READYING != 0 {
		return ExceptionNewf(SystemError, "Type.Ready already readying")
	}

	t.Flags |= TPFLAGS_READYING

	// Initialize tp_base (defaults to BaseObject unless that's us)
	base := t.Base
	if base == nil && t != ObjectType {
		base = ObjectType
		t.Base = base
	}

	// Now the only way base can still be nil is if type is
	// ObjectType.

	// Initialize the base class
	if base != nil && base.Dict == nil {
		err = base.Ready()
		if err != nil {
			return err
		}
	}

	// Initialize ob_type if nil.      This means extensions that want to be
	// compilable separately on Windows can call PyType_Ready() instead of
	// initializing the ob_type field of their type objects.
	// The test for base != nil is really unnecessary, since base is only
	// nil when type is ObjectType, and we know its ob_type is
	// not nil (it's initialized to &PyType_Type).      But coverity doesn't
	// know that.

	// FIXME - this can't work with the current Type scheme
	// if t.Type() == nil && base != nil {
	// 	Py_TYPE(t) = Py_TYPE(base)
	// }

	// Initialize tp_bases
	bases := t.Bases
	if bases == nil {
		if base == nil {
			bases = Tuple{}
		} else {
			bases = Tuple{base}
		}
		t.Bases = bases
	}

	// Initialize tp_dict
	dict := t.Dict
	if dict == nil {
		dict = NewStringDict()
		t.Dict = dict
	}

	// Add type-specific descriptors to tp_dict
	// FIXME not doing this
	// if add_operators(t) < 0 {
	// 	goto error
	// }
	// if t.Methods != nil {
	// 	if add_methods(t, t.Methods) < 0 {
	// 		goto error
	// 	}
	// }
	// if t.Members != nil {
	// 	if add_members(t, t.Members) < 0 {
	// 		goto error
	// 	}
	// }
	// if t.Getset != nil {
	// 	if add_getset(t, t.Getset) < 0 {
	// 		goto error
	// 	}
	// }

	// Calculate method resolution order
	err = t.mro_internal()
	if err != nil {
		return err
	}

	// Inherit special flags from dominant base
	if t.Base != nil {
		t.inherit_special(t.Base)
	}

	// Initialize tp_dict properly
	bases = t.Mro
	if bases == nil {
		panic("Type.Ready: bases is nil")
	}
	// Ignore slots
	// for i := 1; i < len(bases); i++ {
	// 	b, ok := bases[i].(*Type)
	// 	if ok {
	// 		inherit_slots(t, b)
	// 	}
	// }

	// if the type dictionary doesn't contain a __doc__, set it from
	// the tp_doc slot.
	if _, ok := t.Dict["__doc__"]; ok {
		if t.Doc != "" {
			t.Dict["__doc__"] = String(t.Doc)
		} else {
			t.Dict["__doc__"] = None
		}
	}

	// Link into each base class's list of subclasses
	bases = t.Bases
	for i := range bases {
		b, ok := bases[i].(*Type)
		if ok {
			add_subclass(b, t)
		}
	}

	// All done -- set the ready flag
	if t.Dict == nil {
		panic("Type.Ready Dict is nil")
	}
	t.Flags = (t.Flags &^ TPFLAGS_READYING) | TPFLAGS_READY
	return nil
}

func (t *Type) extra_ivars(base *Type) bool {
	return false
	/* FIXME implement this
	   	t_size := t.Basicsize;
	   	b_size := base.Basicsize;

	       assert(t_size >= b_size); // Else type smaller than base!
	       if (t.Itemsize || base.Itemsize) {
	           // If itemsize is involved, stricter rules
	           return t_size != b_size ||
	               t.Itemsize != base.Itemsize;
	       }
	       if (t.Weaklistoffset && base.Weaklistoffset == 0 &&
	           t.Weaklistoffset + sizeof(PyObject *) == t_size &&
	           t.Flags & TPFLAGS_HEAPTYPE)
	           t_size -= sizeof(PyObject *);
	       if (t.Dictoffset && base.Dictoffset == 0 &&
	           t.Dictoffset + sizeof(PyObject *) == t_size &&
	           t.Flags & TPFLAGS_HEAPTYPE)
	           t_size -= sizeof(PyObject *);

	       return t_size != b_size;
	*/
}

func (t *Type) solid_base() *Type {
	var base *Type

	if t.Base != nil {
		base = t.Base.solid_base()
	} else {
		base = ObjectType
	}
	if t.extra_ivars(base) {
		return t
	} else {
		return base
	}
}

// Calculate the best base amongst multiple base classes.
// This is the first one that's on the path to the "solid base".
func best_base(bases Tuple) (*Type, error) {
	// Py_ssize_t i, n;
	// PyTypeObject *base, *winner, *candidate, *base_i;
	// PyObject *base_proto;
	var err error

	if len(bases) == 0 {
		panic("best_base: no bases supplied")
	}
	var base *Type
	var winner *Type
	for i := range bases {
		base_i, ok := bases[i].(*Type)
		if !ok {
			return nil, ExceptionNewf(TypeError, "bases must be types")
		}
		if base_i.Dict == nil {
			err = base_i.Ready()
			if err != nil {
				return nil, err
			}
		}
		candidate := base_i.solid_base()
		if winner == nil {
			winner = candidate
			base = base_i
		} else if winner.IsSubtype(candidate) {
		} else if candidate.IsSubtype(winner) {
			winner = candidate
			base = base_i
		} else {
			return nil, ExceptionNewf(TypeError, "multiple bases have instance lay-out conflict")
		}
	}
	if base == nil {
		return nil, ExceptionNewf(SystemError, "best_base: none found")
	}
	return base, nil
}

// Generic object allocator
func (t *Type) Alloc() *Type {
	// Set the type of the new object to this type
	obj := &Type{
		ObjectType: t,
		Base:       t,
		Dict:       StringDict{},
	}
	return obj
}

// Create a new type
func TypeNew(metatype *Type, args Tuple, kwargs StringDict) (Object, error) {
	// fmt.Printf("TypeNew(type=%q, args=%v, kwargs=%v\n", metatype.Name, args, kwargs)
	var nameObj, basesObj, orig_dictObj Object
	var new_type, base, winner *Type
	// PyHeapTypeObject et;
	// PyMemberDef mp;
	// Py_ssize_t i, nbases, nslots, slotoffset, add_dict, add_weak;
	// _Py_IDENTIFIER(__qualname__);
	// _Py_IDENTIFIER(__slots__);

	// Special case: type(x) should return x.ob_type
	if metatype != nil && len(args) == 1 && len(kwargs) == 0 {
		return args[0].Type(), nil
	}

	// SF bug 475327 -- if that didn't trigger, we need 3
	// arguments. but PyArg_ParseTupleAndKeywords below may give
	// a msg saying type() needs exactly 3.
	if len(args)+len(kwargs) != 3 {
		return nil, ExceptionNewf(TypeError, "type() takes 1 or 3 arguments")
	}

	// Check arguments: (name, bases, dict)
	err := ParseTupleAndKeywords(args, kwargs, "UOO:type", []string{"name", "bases", "dict"},
		&nameObj,
		&basesObj,
		&orig_dictObj)
	if err != nil {
		return nil, err
	}
	name := nameObj.(String)
	bases := basesObj.(Tuple)
	orig_dict := orig_dictObj.(StringDict)

	// Determine the proper metatype to deal with this:
	winner, err = metatype.CalculateMetaclass(bases)
	if err != nil {
		return nil, err
	}

	if winner != metatype {
		//if winner.New != TypeNew { // Pass it to the winner
		// FIXME Nasty hack since you can't compare function pointers in Go
		if fmt.Sprintf("%p", winner.New) != fmt.Sprintf("%p", TypeNew) { // Pass it to the winner
			return winner.New(winner, args, kwargs)
		}
		metatype = winner
	}

	// Adjust for empty tuple bases
	if len(bases) == 0 {
		bases = Tuple{Object(ObjectType)}
	}

	// Calculate best base, and check that all bases are type objects
	base, err = best_base(bases)
	if err != nil {
		return nil, err
	}
	if base.Flags&TPFLAGS_BASETYPE == 0 {
		return nil, ExceptionNewf(TypeError, "type '%s' is not an acceptable base type", base.Name)
	}

	dict := orig_dict.Copy()

	// Check for a __slots__ sequence variable in dict, and count it
	slots, haveSlots := dict["__slots__"]
	nslots := 0
	// add_dict := 0
	// add_weak := 0
	// may_add_dict = base.tp_dictoffset == 0
	// may_add_weak = base.tp_weaklistoffset == 0 && base.tp_itemsize == 0
	if !haveSlots {
		// if may_add_dict {
		// 	add_dict++
		// }
		// if may_add_weak {
		// 	add_weak++
		// }
	} else {
		_ = slots
		return nil, ExceptionNewf(SystemError, "Can't do __slots__ yet")
		/* FIXME ignore slots for the moment
		// Have slots

		// Make it into a tuple
		if PyUnicode_Check(slots) {
			slots = PyTuple_Pack(1, slots)
		} else {
			slots = PySequence_Tuple(slots)
		}
		if slots == nil {
			goto error
		}
		assert(PyTuple_Check(slots))

		// Are slots allowed?
		nslots = PyTuple_GET_SIZE(slots)
		if nslots > 0 && base.tp_itemsize != 0 {
			PyErr_Format(PyExc_TypeError,
				"nonempty __slots__ not supported for subtype of '%s'",
				base.tp_name)
			goto error
		}

		// Check for valid slot names and two special cases
		for i = 0; i < nslots; i++ {
			PyObject * tmp = PyTuple_GET_ITEM(slots, i)
			if !valid_identifier(tmp) {
				goto error
			}
			assert(PyUnicode_Check(tmp))
			if _PyUnicode_CompareWithId(tmp, &PyId___dict__) == 0 {
				if !may_add_dict || add_dict {
					PyErr_SetString(PyExc_TypeError,
						"__dict__ slot disallowed: we already got one")
					goto error
				}
				add_dict++
			}
			if PyUnicode_CompareWithASCIIString(tmp, "__weakref__") == 0 {
				if !may_add_weak || add_weak {
					PyErr_SetString(PyExc_TypeError,
						"__weakref__ slot disallowed: either we already got one, or __itemsize__ != 0")
					goto error
				}
				add_weak++
			}
		}

		// Copy slots into a list, mangle names and sort them.
		// Sorted names are needed for __class__ assignment.
		// Convert them back to tuple at the end.

		newslots = PyList_New(nslots - add_dict - add_weak)
		for i, j = 0, 0; i < nslots; i++ {
			tmp = PyTuple_GET_ITEM(slots, i)
			if (add_dict &&
				_PyUnicode_CompareWithId(tmp, &PyId___dict__) == 0) ||
				(add_weak &&
					PyUnicode_CompareWithASCIIString(tmp, "__weakref__") == 0) {
				continue
			}
			tmp = _Py_Mangle(name, tmp)
			if !tmp {
				goto error
			}
			PyList_SET_ITEM(newslots, j, tmp)
			if PyDict_GetItem(dict, tmp) {
				PyErr_Format(PyExc_ValueError,
					"%R in __slots__ conflicts with class variable",
					tmp)
				goto error
			}
			j++
		}
		assert(j == nslots-add_dict-add_weak)
		nslots = j
		Py_CLEAR(slots)
		if PyList_Sort(newslots) == -1 {
			goto error
		}
		slots = PyList_AsTuple(newslots)

		// Secondary bases may provide weakrefs or dict
		if nbases > 1 &&
			((may_add_dict && !add_dict) ||
				(may_add_weak && !add_weak)) {
			for i = 0; i < nbases; i++ {
				tmp = PyTuple_GET_ITEM(bases, i)
				if tmp == base {
					continue // Skip primary base
				}
				// assert(PyType_Check(tmp));
				tmptype = tmp.(*Type)
				if may_add_dict && !add_dict &&
					tmptype.tp_dictoffset != 0 {
					add_dict++
				}
				if may_add_weak && !add_weak &&
					tmptype.tp_weaklistoffset != 0 {
					add_weak++
				}
				if may_add_dict && !add_dict {
					continue
				}
				if may_add_weak && !add_weak {
					continue
				}
				// Nothing more to check
				break
			}
		}
		*/
	}

	// Allocate the type object
	_ = nslots // FIXME
	new_type = metatype.Alloc()
	new_type.New = ObjectNew   // FIXME metatype.New // FIXME?
	new_type.Init = ObjectInit // FIXME metatype.New // FIXME?

	// Keep name and slots alive in the extended type object
	et := new_type
	et.Name = string(name)
	// FIXME et.Slots = slots
	slots = nil

	// Initialize tp_flags
	new_type.Flags = TPFLAGS_DEFAULT | TPFLAGS_HEAPTYPE | TPFLAGS_BASETYPE

	// Set tp_base and tp_bases
	new_type.Bases = bases
	bases = nil
	new_type.Base = base

	// Initialize tp_dict from passed-in dict
	new_type.Dict = dict
	// fmt.Printf("New type dict is %v\n", dict)

	// Set __module__ in the dict
	if _, ok := dict["__module__"]; !ok {
		fmt.Printf("*** FIXME need to get the current vm globals somehow\n")
		// tmp = PyEval_GetGlobals()
		// if tmp != nil {
		// 	tmp, ok := tmp["__name__"]
		// 	if ok {
		// 		dict["__module__"] = tmp
		// 	}
		// }
	}

	// Set ht_qualname to dict['__qualname__'] if available, else to
	// __name__.  The __qualname__ accessor will look for ht_qualname.
	if qualname, ok := dict["__qualname__"]; ok {
		if Qualname, ok := qualname.(String); !ok {
			return nil, ExceptionNewf(TypeError, "type __qualname__ must be a str, not %s", qualname.Type().Name)
		} else {
			et.Qualname = string(Qualname)
		}
		delete(dict, "__qualname__")
	} else {
		et.Qualname = et.Name
	}

	// Set tp_doc to a copy of dict['__doc__'], if the latter is there
	// and is a string.  The __doc__ accessor will first look for tp_doc;
	// if that fails, it will still look into __dict__.
	if doc, ok := dict["__doc__"]; ok {
		if Doc, ok := doc.(String); ok {
			new_type.Doc = string(Doc)
		}
	}

	// Special-case __new__: if it's a plain function,
	// make it a static function
	// FIXME
	// tmp = dict["__new__"]
	// if tmp != nil && PyFunction_Check(tmp) {
	// 	tmp = PyStaticMethod_New(tmp)
	// 	if _PyDict_SetItemId(dict, &PyId___new__, tmp) < 0 {
	// 		goto error
	// 	}
	// }

	/*
		// Add descriptors for custom slots from __slots__, or for __dict__
		mp = PyHeapType_GET_MEMBERS(et)
		slotoffset = base.tp_basicsize
		if et.ht_slots != nil {
			for i = 0; i < nslots; i++ {
				mp.name = _PyUnicode_AsString(
					PyTuple_GET_ITEM(et.ht_slots, i))
				mp.new_type = T_OBJECT_EX
				mp.offset = slotoffset

				// __dict__ and __weakref__ are already filtered out
				assert(strcmp(mp.name, "__dict__") != 0)
				assert(strcmp(mp.name, "__weakref__") != 0)

				slotoffset += 1 // FIXME sizeof(PyObject *);
				mp++
			}
		}
		if add_dict {
			// if (base.tp_itemsize)
			//     new_type.tp_dictoffset = -sizeof(PyObject *);
			// else
			//     new_type.tp_dictoffset = slotoffset;
			slotoffset += 1 // sizeof(PyObject *);
		}
		if new_type.tp_dictoffset {
			et.ht_cached_keys = _PyDict_NewKeysForClass()
		}
		if add_weak {
			assert(!base.tp_itemsize)
			new_type.tp_weaklistoffset = slotoffset
			slotoffset += 1 // FIXME sizeof(PyObject *);
		}
		new_type.tp_basicsize = slotoffset
		new_type.tp_itemsize = base.tp_itemsize
		new_type.tp_members = PyHeapType_GET_MEMBERS(et)

		if new_type.tp_weaklistoffset && new_type.tp_dictoffset {
			new_type.tp_getset = subtype_getsets_full
		} else if new_type.tp_weaklistoffset && !new_type.tp_dictoffset {
			new_type.tp_getset = subtype_getsets_weakref_only
		} else if !new_type.tp_weaklistoffset && new_type.tp_dictoffset {
			new_type.tp_getset = subtype_getsets_dict_only
		} else {
			new_type.tp_getset = nil
		}

		// Special case some slots
		if new_type.tp_dictoffset != 0 || nslots > 0 {
			if base.tp_getattr == nil && base.tp_getattro == nil {
				new_type.tp_getattro = PyObject_GenericGetAttr
			}
			if base.tp_setattr == nil && base.tp_setattro == nil {
				new_type.tp_setattro = PyObject_GenericSetAttr
			}
		}
		new_type.tp_dealloc = subtype_dealloc

		// Enable GC unless there are really no instance variables possible
		if !(new_type.tp_basicsize == sizeof(PyObject) &&
			new_type.tp_itemsize == 0) {
			new_type.tp_flags |= TPFLAGS_HAVE_GC
		}

		// Always override allocation strategy to use regular heap
		new_type.tp_alloc = PyType_GenericAlloc
		if new_type.tp_flags & TPFLAGS_HAVE_GC {
			new_type.tp_free = PyObject_GC_Del
			new_type.tp_traverse = subtype_traverse
			new_type.tp_clear = subtype_clear
		} else {
			new_type.tp_free = PyObject_Del
		}

	*/
	// Initialize the rest
	err = new_type.Ready()
	if err != nil {
		return nil, err
	}

	// Put the proper slots in place
	// fixup_slot_dispatchers(new_type)

	return new_type, nil
}

func TypeInit(cls Object, args Tuple, kwargs StringDict) error {
	if len(kwargs) != 0 {
		return ExceptionNewf(TypeError, "type.__init__() takes no keyword arguments")
	}

	if len(args) != 1 && len(args) != 3 {
		return ExceptionNewf(TypeError, "type.__init__() takes 1 or 3 arguments")
	}

	// Call object.__init__(self) now.
	// XXX Could call super(type, cls).__init__() but what's the point?
	return ObjectInit(cls, nil, nil)
}

// The base type of all types (eventually)... except itself.

// You may wonder why object.__new__() only complains about arguments
// when object.__init__() is not overridden, and vice versa.
//
// Consider the use cases:
//
// 1. When neither is overridden, we want to hear complaints about
//    excess (i.e., any) arguments, since their presence could
//    indicate there's a bug.
//
// 2. When defining an Immutable type, we are likely to override only
//    __new__(), since __init__() is called too late to initialize an
//    Immutable object.  Since __new__() defines the signature for the
//    type, it would be a pain to have to override __init__() just to
//    stop it from complaining about excess arguments.
//
// 3. When defining a Mutable type, we are likely to override only
//    __init__().  So here the converse reasoning applies: we don't
//    want to have to override __new__() just to stop it from
//    complaining.
//
// 4. When __init__() is overridden, and the subclass __init__() calls
//    object.__init__(), the latter should complain about excess
//    arguments; ditto for __new__().
//
// Use cases 2 and 3 make it unattractive to unconditionally check for
// excess arguments.  The best solution that addresses all four use
// cases is as follows: __init__() complains about excess arguments
// unless __new__() is overridden and __init__() is not overridden
// (IOW, if __init__() is overridden or __new__() is not overridden);
// symmetrically, __new__() complains about excess arguments unless
// __init__() is overridden and __new__() is not overridden
// (IOW, if __new__() is overridden or __init__() is not overridden).
//
// However, for backwards compatibility, this breaks too much code.
// Therefore, in 2.6, we'll *warn* about excess arguments when both
// methods are overridden; for all other cases we'll use the above
// rules.

// Return true if any arguments supplied
func excess_args(args Tuple, kwargs StringDict) bool {
	return len(args) != 0 || len(kwargs) != 0
}

func ObjectInit(self Object, args Tuple, kwargs StringDict) error {
	t := self.Type()
	// FIXME bodge to compare function pointers
	// if excess_args(args, kwargs) && (fmt.Sprintf("%p", t.New) == fmt.Sprintf("%p", ObjectNew) || fmt.Sprintf("%p", t.Init) != fmt.Sprintf("%p", ObjectInit)) {
	// 	return ExceptionNewf(TypeError, "object.__init__() takes no parameters")
	// }

	// FIXME this isn't correct probably
	// Check args for object()
	if t == ObjectType && excess_args(args, kwargs) {
		return ExceptionNewf(TypeError, "object.__init__() takes no parameters")
	}

	// Call the __init__ method if it exists
	// FIXME this isn't the way cpython does it - it adjusts the function pointers
	// Only do this for non built in types
	if _, ok := self.(*Type); ok {
		init := t.GetAttrOrNil("__init__")
		// fmt.Printf("init = %v\n", init)
		if init != nil {
			newArgs := make(Tuple, len(args)+1)
			newArgs[0] = self
			copy(newArgs[1:], args)
			_, err := Call(init, newArgs, kwargs)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func ObjectNew(t *Type, args Tuple, kwargs StringDict) (Object, error) {
	// FIXME bodge to compare function pointers
	// if excess_args(args, kwargs) && (fmt.Sprintf("%p", t.Init) == fmt.Sprintf("%p", ObjectInit) || fmt.Sprintf("%p", t.New) != fmt.Sprintf("%p", ObjectNew)) {
	// 	return ExceptionNewf(TypeError, "object() takes no parameters")
	// }

	// FIXME this isn't correct probably
	// Check arguments to new only for object
	if t == ObjectType && excess_args(args, kwargs) {
		return nil, ExceptionNewf(TypeError, "object() takes no parameters")
	}

	// FIXME abstrac ty pes
	// if (type->tp_flags & TPFLAGS_IS_ABSTRACT) {
	// 	PyObject *abstract_methods = NULL;
	// 	PyObject *builtins;
	// 	PyObject *sorted;
	// 	PyObject *sorted_methods = NULL;
	// 	PyObject *joined = NULL;
	// 	PyObject *comma;
	// 	_Py_static_string(comma_id, ", ");
	// 	_Py_IDENTIFIER(sorted);

	// 	// Compute ", ".join(sorted(type.__abstractmethods__))
	// 	// into joined.
	// 	abstract_methods = type_abstractmethods(type, NULL);
	// 	if (abstract_methods == NULL) {
	// 		goto error;
	// 	}
	// 	builtins = PyEval_GetBuiltins();
	// 	if (builtins == NULL) {
	// 		goto error;
	// 	}
	// 	sorted = _PyDict_GetItemId(builtins, &PyId_sorted);
	// 	if (sorted == NULL) {
	// 		goto error;
	// 	}
	// 	sorted_methods = PyObject_CallFunctionObjArgs(sorted,
	// 		abstract_methods,
	// 		NULL);
	// 	if (sorted_methods == NULL) {
	// 		goto error;
	// 	}
	// 	comma = _PyUnicode_FromId(&comma_id);
	// 	if (comma == NULL) {
	// 		goto error;
	// 	}
	// 	joined = PyUnicode_Join(comma, sorted_methods);
	// 	if (joined == NULL) {
	// 		goto error;
	// 	}

	// 	PyErr_Format(PyExc_TypeError,
	// 		"Can't instantiate abstract class %s "
	// 		"with abstract methods %U",
	// 		type->tp_name,
	// 		joined);
	// error:
	// 	Py_XDECREF(joined);
	// 	Py_XDECREF(sorted_methods);
	// 	Py_XDECREF(abstract_methods);
	// 	return NULL;
	// }
	return t.Alloc(), nil
}

// FIXME this should be the default?
func (ty *Type) M__eq__(other Object) (Object, error) {
	if otherTy, ok := other.(*Type); ok && ty == otherTy {
		return True, nil
	}
	return False, nil
}

// FIXME this should be the default?
func (ty *Type) M__ne__(other Object) (Object, error) {
	if otherTy, ok := other.(*Type); ok && ty == otherTy {
		return False, nil
	}
	return True, nil
}

func (ty *Type) M__str__() (Object, error) {
	if res, ok, err := ty.CallMethod("__str__", Tuple{ty}, nil); ok {
		return res, err
	}
	return ty.M__repr__()
}

func (ty *Type) M__repr__() (Object, error) {
	if res, ok, err := ty.CallMethod("__repr__", Tuple{ty}, nil); ok {
		return res, err
	}
	if ty.Name == "" {
		// FIXME not a good way to tell objects from classes!
		return String(fmt.Sprintf("<%s object at %p>", ty.Type().Name, ty)), nil
	}
	return String(fmt.Sprintf("<class '%s'>", ty.Name)), nil

}

// Make sure it satisfies the interface
var _ Object = (*Type)(nil)
var _ I__call__ = (*Type)(nil)
var _ IGetDict = (*Type)(nil)
var _ I__repr__ = (*Type)(nil)
var _ I__str__ = (*Type)(nil)
