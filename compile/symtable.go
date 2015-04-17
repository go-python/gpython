// Make a symbol table to identify which variables are glocal, closure or local

package compile

import (
	"github.com/ncw/gpython/ast"
	"github.com/ncw/gpython/py"
)

//go:generate stringer -type=Scope,BlockType -output stringer.go

// Scope
type Scope uint8

// Scope definitions
const (
	scopeInvalid Scope = iota
	scopeLocal
	scopeGlobalExplicit
	scopeGlobalImplicit
	scopeFree
	scopeCell
)

// Def-use flag information
type DefUse uint8

// Flags for def-use information
const (
	defGlobal    DefUse = 1 << iota // global stmt
	defLocal                        // assignment in code block
	defParam                        // formal parameter
	defNonlocal                     // nonlocal stmt
	defUse                          // name is used
	defFree                         // name used but not defined in nested block
	defFreeClass                    // free variable from class's method
	defImport                       // assignment occurred via import

	defBound = (defLocal | defParam | defImport)

	// GLOBAL_EXPLICIT and GLOBAL_IMPLICIT are used internally by the symbol
	// table.  GLOBAL is returned from PyST_GetScope() for either of them.
	// It is stored in ste_symbols at bits 12-15.
	defScopeMask = (defGlobal | defLocal | defParam | defNonlocal)
)

// BlockType for SymTable
type BlockType uint8

// BlockTypes
const (
	FunctionBlock BlockType = iota
	ClassBlock
	ModuleBlock
)

type Symbol struct {
	Name      string
	Scope     Scope
	Flags     DefUse
	Namespace *SymTable
}

type Symbols map[string]Symbol

func NewSymbols() Symbols {
	return make(Symbols)
}

type SymTable struct {
	Type       BlockType // 'class', 'module', and 'function'
	Name       string    // name of the class if the table is for a class, the name of the function if the table is for a function, or 'top' if the table is global (get_type() returns 'module').
	Lineno     int       // number of the first line in the block this table represents.
	Optimized  bool      // True if the locals in this table can be optimized.
	Nested     bool      // True if the block is a nested class or function.
	Exec       bool      // True if the block uses exec.
	ImportStar bool      // Return True if the block uses a starred from-import.
	Symbols    Symbols
	Global     *SymTable
	Parent     *SymTable
	Varnames   []string // list of function parameters
}

// Make a new top symbol table from the ast supplied
func NewSymTable(Ast ast.Ast) *SymTable {
	return newSymTable(Ast, ModuleBlock, "top", nil)
}

// Make a new symbol table from the ast supplied of the given type
func newSymTable(Ast ast.Ast, Type BlockType, Name string, parent *SymTable) *SymTable {
	st := &SymTable{
		Type:    ModuleBlock,
		Name:    "top",
		Parent:  parent,
		Symbols: NewSymbols(),
	}
	if parent == nil {
		st.Global = st
	} else {
		st.Global = parent.Global
	}
	st.Parse(Ast)
	return st
}

// Parse the ast into the symbol table
func (st *SymTable) Parse(Ast ast.Ast) {
	ast.Walk(Ast, func(Ast ast.Ast) bool {
		// New symbol tables needed at these points
		// FunctionDef
		// ClassDef
		// Lambda
		// Comprehension (all types of comprehension in py3)

		switch node := Ast.(type) {
		case *ast.Nonlocal:
			for _, name := range node.Names {
				st.AddDef(name, defNonlocal)
			}
		case *ast.Global:
			for _, name := range node.Names {
				st.AddDef(name, defGlobal)
			}
		case *ast.Name:
			if node.Ctx == ast.Load {
				st.AddDef(node.Id, defUse)
			} else {
				st.AddDef(node.Id, defLocal)
			}
			// FIXME Special-case super: it counts as a use of __class__
			// if node.Name.Ctx == ast.Load && st.st_cur.ste_type == FunctionBlock && e.v.Name.id == "super" {
			// 	st.AddDef(ast.Identifier("__class__"), defUse)
			// }
		case *ast.FunctionDef:
			st.AddDef(node.Name, defLocal)
		case *ast.ClassDef:
			st.AddDef(node.Name, defLocal)
		case *ast.Lambda:

			// Comprehensions
		case *ast.ListComp:
		case *ast.SetComp:
		case *ast.DictComp:
		case *ast.GeneratorExp:

		}
		return true
	})
}

const duplicateArgument = "duplicate argument %q in function definition"

// Add a symbol into the symble table
func (st *SymTable) AddDef(name ast.Identifier, flags DefUse) {
	// FIXME mangled := _Py_Mangle(st.st_private, name)
	mangled := string(name)

	// Add or update the symbol in the Symbols
	if sym, ok := st.Symbols[mangled]; ok {
		if (flags&defParam) != 0 && (sym.Flags&defParam) != 0 {
			/* Is it better to use 'mangled' or 'name' here? */
			panic(py.ExceptionNewf(py.SyntaxError, duplicateArgument, name))
			// FIXME
			// PyErr_SyntaxLocationObject(st.st_filename,
			// 	st.st_cur.ste_lineno,
			// 	st.st_cur.ste_col_offset)
		}
		sym.Flags |= flags
		st.Symbols[mangled] = sym
	} else {
		st.Symbols[mangled] = Symbol{
			Name:  string(name),
			Scope: 0, // FIXME
			Flags: flags,
		}
	}

	if (flags & defParam) != 0 {
		st.Varnames = append(st.Varnames, mangled)
	} else if (flags & defGlobal) != 0 {
		// If it is a global definition then add it in the global Symbols
		//
		// XXX need to update DEF_GLOBAL for other flags too;
		// perhaps only DEF_FREE_GLOBAL
		if sym, ok := st.Global.Symbols[mangled]; ok {
			sym.Flags |= flags
			st.Global.Symbols[mangled] = sym
		} else {
			st.Global.Symbols[mangled] = Symbol{
				Name:  string(name),
				Scope: 0, // FIXME
				Flags: flags,
			}
		}
	}
}
