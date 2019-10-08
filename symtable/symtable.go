// Copyright 2018 The go-python Authors.  All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Make a symbol table to identify which variables are glocal, closure or local

// FIXME stop panics escaping this package

// FIXME tests for Varkeywords etc

// FIXME need to set locations for panics, eg
// PyErr_SyntaxLocationObject(st.ste_table.st_filename,
// 	st.ste_opt_lineno,
// 	st.ste_opt_col_offset)

package symtable

import (
	"fmt"
	"sort"
	"strings"

	"github.com/go-python/gpython/ast"
	"github.com/go-python/gpython/py"
)

//go:generate stringer -type=Scope,BlockType -output stringer.go

// Scope
type Scope uint8

// Scope definitions
const (
	ScopeInvalid Scope = iota
	ScopeLocal
	ScopeGlobalExplicit
	ScopeGlobalImplicit
	ScopeFree
	ScopeCell
)

// Accumulate Scopes
type Scopes map[string]Scope

// Def-use flag information
type DefUseFlags uint8

// Flags for def-use information
const (
	DefGlobal    DefUseFlags = 1 << iota // global stmt
	DefLocal                             // assignment in code block
	DefParam                             // formal parameter
	DefNonlocal                          // nonlocal stmt
	DefUse                               // name is used
	DefFree                              // name used but not defined in nested block
	DefFreeClass                         // free variable from class's method
	DefImport                            // assignment occurred via import

	DefBound = (DefLocal | DefParam | DefImport)

	// ScopeGlobalExplicit and ScopeGlobalImplicit are used internally by the symbol
	// table.  GLOBAL is returned from PyST_GetScope() for either of them.
	// It is stored in ste_symbols at bits 12-15.
	DefScopeMask = (DefGlobal | DefLocal | DefParam | DefNonlocal)
)

// BlockType for SymTable
type BlockType uint8

// BlockTypes
const (
	FunctionBlock BlockType = iota
	ClassBlock
	ModuleBlock
)

// OptType for SymTable
type OptType uint8

// The following flag names are used for the Unoptimized bit field
const (
	optImportStar OptType = 1 << iota
	optTopLevel           // top-level names, including eval and exec
)

// Info about a symbol
type Symbol struct {
	Scope     Scope
	Flags     DefUseFlags
	Lineno    int
	ColOffset int
}

type Symbols map[string]Symbol

type Children []*SymTable

type LookupChild map[ast.Ast]*SymTable

type SymTable struct {
	Type              BlockType // 'class', 'module', and 'function'
	Name              string    // name of the class if the table is for a class, the name of the function if the table is for a function, or 'top' if the table is global (get_type() returns 'module').
	Filename          string    // filename that this is being parsed from
	Lineno            int       // number of the first line in the block this table represents.
	Unoptimized       OptType   // false if namespace is optimized
	Nested            bool      // true if the block is a nested class or function.
	Free              bool      // true if block has free variables
	ChildFree         bool      // true if a child block has free vars, including free refs to globals
	Generator         bool      // true if namespace is a generator
	Varargs           bool      // true if block has varargs
	Varkeywords       bool      // true if block has varkeywords
	ReturnsValue      bool      // true if namespace uses return with an argument
	NeedsClassClosure bool      // for class scopes, true if a closure over __class__ should be created
	// col_offset        int       // offset of first line of block
	// opt_lineno        int       // lineno of last exec or import *
	// opt_col_offset    int       // offset of last exec or import *
	TmpName int    // counter for listcomp temp vars
	Private string // name of current class or ""

	Symbols     Symbols
	Global      *SymTable // symbol table entry for module
	Parent      *SymTable
	Varnames    []string    // list of function parameters
	Children    Children    // Child SymTables
	LookupChild LookupChild // Child symtables keyed by ast
}

// Make a new top symbol table from the ast supplied
func NewSymTable(Ast ast.Ast, filename string) (st *SymTable, err error) {
	defer func() {
		if r := recover(); r != nil {
			err = py.MakeException(r)
		}
	}()
	st = newSymTable(ModuleBlock, "top", nil)
	st.Unoptimized = optTopLevel
	st.Filename = filename
	// Parse into the symbol table
	st.Parse(Ast)
	// Analyze the symbolt table
	st.Analyze()
	return st, nil
}

// Make a new symbol table from the ast supplied of the given type
func newSymTable(Type BlockType, Name string, parent *SymTable) *SymTable {
	st := &SymTable{
		Type:        Type,
		Name:        Name,
		Parent:      parent,
		Symbols:     make(Symbols),
		Children:    make(Children, 0),
		LookupChild: make(LookupChild),
	}
	if parent == nil {
		st.Global = st
	} else {
		st.Global = parent.Global
		st.Nested = parent.Nested || (parent.Type == FunctionBlock)
		st.Filename = parent.Filename
	}
	return st
}

// Make a new symtable and add it to parent
func newSymTableBlock(Ast ast.Ast, Type BlockType, Name string, parent *SymTable) *SymTable {
	stNew := newSymTable(Type, Name, parent)
	parent.Children = append(parent.Children, stNew)
	parent.LookupChild[Ast] = stNew
	// FIXME set stNew.Lineno
	return stNew
}

// Panics abount a syntax error at this line and col
func (st *SymTable) panicSyntaxErrorLinenof(lineno, offset int, format string, a ...interface{}) {
	err := py.ExceptionNewf(py.SyntaxError, format, a...)
	err = py.MakeSyntaxError(err, st.Filename, lineno, offset, "")
	panic(err)
}

// Panics abount a syntax error on this ast node
func (st *SymTable) panicSyntaxErrorf(Ast ast.Ast, format string, a ...interface{}) {
	err := py.ExceptionNewf(py.SyntaxError, format, a...)
	err = py.MakeSyntaxError(err, st.Filename, Ast.GetLineno(), Ast.GetColOffset(), "")
	panic(err)
}

// FindChild finds SymTable attached to Ast - returns nil if not found
func (st *SymTable) FindChild(Ast ast.Ast) *SymTable {
	return st.LookupChild[Ast]
}

// GetScope finds the scope for the name, returns ScopeInvalid if not found
func (st *SymTable) GetScope(name string) Scope {
	symbol, ok := st.Symbols[name]
	if !ok {
		return ScopeInvalid
	}
	return symbol.Scope
}

// Add arguments to the symbol table
func (st *SymTable) addArgumentsToSymbolTable(node *ast.Arguments) {
	// skip default arguments inside function block
	// XXX should ast be different?
	for _, arg := range node.Args {
		st.AddDef(node, arg.Arg, DefParam)
	}
	for _, arg := range node.Kwonlyargs {
		st.AddDef(node, arg.Arg, DefParam)
	}
	if node.Vararg != nil {
		st.AddDef(node, node.Vararg.Arg, DefParam)
		st.Varargs = true
	}
	if node.Kwarg != nil {
		st.AddDef(node, node.Kwarg.Arg, DefParam)
		st.Varkeywords = true
	}
}

// Parse the ast into the symbol table
func (st *SymTable) Parse(Ast ast.Ast) {
	ast.Walk(Ast, func(Ast ast.Ast) bool {
		switch node := Ast.(type) {
		case *ast.Nonlocal:
			for _, name := range node.Names {
				cur, ok := st.Symbols[string(name)]
				if ok {
					if (cur.Flags & DefLocal) != 0 {
						st.panicSyntaxErrorf(node, "name '%s' is assigned to before nonlocal declaration", name)
					}
					if (cur.Flags & DefUse) != 0 {
						st.panicSyntaxErrorf(node, "name '%s' is used prior to nonlocal declaration", name)
					}
				}
				st.AddDef(node, name, DefNonlocal)
			}
		case *ast.Global:
			for _, name := range node.Names {
				cur, ok := st.Symbols[string(name)]
				if ok {
					if (cur.Flags & DefLocal) != 0 {
						st.panicSyntaxErrorf(node, "name '%s' is assigned to before global declaration", name)

					}
					if (cur.Flags & DefUse) != 0 {
						st.panicSyntaxErrorf(node, "name '%s' is used prior to global declaration", name)
					}
				}
				st.AddDef(node, name, DefGlobal)
			}
		case *ast.Name:
			if node.Ctx == ast.Load {
				st.AddDef(node, node.Id, DefUse)
			} else {
				st.AddDef(node, node.Id, DefLocal)
			}
			// Special-case super: it counts as a use of __class__
			if node.Ctx == ast.Load && st.Type == FunctionBlock && node.Id == "super" {
				st.AddDef(node, ast.Identifier("__class__"), DefUse)
			}
		case *ast.FunctionDef:
			// Add the function name to the SymTable
			st.AddDef(node, node.Name, DefLocal)
			name := string(node.Name)

			// Walk these things in original symbol table
			if node.Args != nil {
				st.Parse(node.Args)
			}
			for _, expr := range node.DecoratorList {
				st.Parse(expr)
			}
			st.Parse(node.Returns)

			// Make a new symtable
			stNew := newSymTableBlock(Ast, FunctionBlock, name, st)

			// Add the arguments to the new symbol table
			stNew.addArgumentsToSymbolTable(node.Args)

			// Walk the Body in the new symtable
			for _, stmt := range node.Body {
				stNew.Parse(stmt)
			}

			// return false to stop the parse
			return false
		case *ast.ClassDef:
			st.AddDef(node, node.Name, DefLocal)
			name := string(node.Name)
			// Parse in the original symtable
			for _, expr := range node.Bases {
				st.Parse(expr)
			}
			for _, keyword := range node.Keywords {
				st.Parse(keyword)
			}
			if node.Starargs != nil {
				st.Parse(node.Starargs)
			}
			if node.Kwargs != nil {
				st.Parse(node.Kwargs)
			}
			for _, expr := range node.DecoratorList {
				st.Parse(expr)
			}
			// Make a new symtable
			stNew := newSymTableBlock(Ast, ClassBlock, name, st)
			stNew.Private = name // set name of class
			// Parse body in new symtable
			for _, stmt := range node.Body {
				stNew.Parse(stmt)
			}
			// return false to stop the parse
			return false
		case *ast.Lambda:
			// Parse in the original symtable
			if node.Args != nil {
				st.Parse(node.Args)
			}

			// Make a new symtable
			stNew := newSymTableBlock(Ast, FunctionBlock, "lambda", st)

			// Add the arguments to the new symbol table
			stNew.addArgumentsToSymbolTable(node.Args)

			// Walk the Body in the new symtable
			stNew.Parse(node.Body)

			// return false to stop the parse
			return false
		case *ast.ListComp:
			st.parseComprehension(Ast, "listcomp", node.Generators, node.Elt, nil)
			return false
		case *ast.SetComp:
			st.parseComprehension(Ast, "setcomp", node.Generators, node.Elt, nil)
			return false
		case *ast.DictComp:
			st.parseComprehension(Ast, "dictcomp", node.Generators, node.Key, node.Value)
			return false
		case *ast.GeneratorExp:
			st.parseComprehension(Ast, "genexpr", node.Generators, node.Elt, nil)
			return false
		case *ast.ExceptHandler:
			if node.Name != "" {
				st.AddDef(node, node.Name, DefLocal)
			}
		case *ast.Alias:
			// Compute store_name, the name actually bound by the import
			// operation.  It is different than node.name when node.name is a
			// dotted package name (e.g. spam.eggs)
			name := node.Name
			if node.AsName != "" {
				name = node.AsName
			}
			dot := strings.LastIndex(string(name), ".")
			store_name := name
			if dot >= 0 {
				store_name = name[:dot]
			}
			if name != "*" {
				st.AddDef(node, store_name, DefImport)
			} else {
				if st.Type != ModuleBlock {
					st.panicSyntaxErrorf(node, "import * only allowed at module level")
				}
				st.Unoptimized |= optImportStar
			}
		case *ast.Return:
			if node.Value != nil {
				st.ReturnsValue = true
			}
		case *ast.Yield, *ast.YieldFrom:
			st.Generator = true
		}
		return true
	})
}

// make a new temporary name
func (st *SymTable) newTmpName(node ast.Ast) {
	st.TmpName++
	id := ast.Identifier(fmt.Sprintf("_[%d]", st.TmpName))
	st.AddDef(node, id, DefLocal)
}

func (st *SymTable) parseComprehension(Ast ast.Ast, scopeName string, generators []ast.Comprehension, elt ast.Expr, value ast.Expr) {
	_, isGenerator := Ast.(*ast.GeneratorExp)
	needsTmp := !isGenerator
	outermost := generators[0]
	// Outermost iterator is evaluated in current scope
	st.Parse(outermost.Iter)
	// Create comprehension scope for the rest
	stNew := newSymTableBlock(Ast, FunctionBlock, scopeName, st)
	stNew.Generator = isGenerator
	// Outermost iter is received as an argument
	id := ast.Identifier(fmt.Sprintf(".%d", 0))
	stNew.AddDef(Ast, id, DefParam)
	// Allocate temporary name if needed
	if needsTmp {
		stNew.newTmpName(Ast)
	}
	stNew.Parse(outermost.Target)
	for _, expr := range outermost.Ifs {
		stNew.Parse(expr)
	}
	for _, comprehension := range generators[1:] {
		stNew.Parse(comprehension.Target)
		stNew.Parse(comprehension.Iter)
		for _, expr := range comprehension.Ifs {
			stNew.Parse(expr)
		}
	}
	if value != nil {
		stNew.Parse(value)
	}
	stNew.Parse(elt)
}

// Add a symbol into the symble table
func (st *SymTable) AddDef(node ast.Ast, name ast.Identifier, flags DefUseFlags) {
	// FIXME mangled := _Py_Mangle(st.Private, name)
	mangled := string(name)

	// Add or update the symbol in the Symbols
	if sym, ok := st.Symbols[mangled]; ok {
		if (flags&DefParam) != 0 && (sym.Flags&DefParam) != 0 {
			// Is it better to use 'mangled' or 'name' here?
			st.panicSyntaxErrorf(node, "duplicate argument '%s' in function definition", name)
			// FIXME
			// PyErr_SyntaxLocationObject(st.st_filename,
			// 	st.st_cur.ste_lineno,
			// 	st.st_cur.ste_col_offset)
		}
		sym.Flags |= flags
		st.Symbols[mangled] = sym
	} else {
		st.Symbols[mangled] = Symbol{
			Scope:     0, // FIXME
			Flags:     flags,
			Lineno:    node.GetLineno(),
			ColOffset: node.GetColOffset(),
		}
	}

	if (flags & DefParam) != 0 {
		st.Varnames = append(st.Varnames, mangled)
	} else if (flags & DefGlobal) != 0 {
		// If it is a global definition then add it in the global Symbols
		//
		// XXX need to update DefGlobal for other flags too;
		// perhaps only DefFreeClass
		if sym, ok := st.Global.Symbols[mangled]; ok {
			sym.Flags |= flags
			st.Global.Symbols[mangled] = sym
		} else {
			st.Global.Symbols[mangled] = Symbol{
				Scope:     0, // FIXME
				Flags:     flags,
				Lineno:    node.GetLineno(),
				ColOffset: node.GetColOffset(),
			}
		}
	}
}

// StringSet for storing strings
type StringSet map[string]struct{}

// Add adds elem to the set, returning the set
//
// If the element already exists then it has no effect
func (s StringSet) Add(elem string) {
	s[elem] = struct{}{}
}

// Update adds all the elements from the other set to this set.
func (s StringSet) Update(other StringSet) {
	for elem := range other {
		s[elem] = struct{}{}
	}
}

// Copy makes a shallow copy of s
func (s StringSet) Copy() StringSet {
	copy := make(StringSet, len(s))
	copy.Update(s)
	return copy
}

// Discard ensures k is not in the set returning true if it was found
func (s StringSet) Discard(k string) bool {
	if _, ok := s[k]; !ok {
		return false
	}
	delete(s, k)
	return true
}

// Contains returns true if k is in s
func (s StringSet) Contains(k string) bool {
	_, ok := s[k]
	return ok
}

/* Analyze raw symbol information to determine scope of each name.

   The next several functions are helpers for symtable_analyze(),
   which determines whether a name is local, global, or free.  In addition,
   it determines which local variables are cell variables; they provide
   bindings that are used for free variables in enclosed blocks.

   There are also two kinds of global variables, implicit and explicit.  An
   explicit global is declared with the global statement.  An implicit
   global is a free variable for which the compiler has found no binding
   in an enclosing function scope.  The implicit global is either a global
   or a builtin.  Python's module and class blocks use the xxx_NAME opcodes
   to handle these names to implement slightly odd semantics.  In such a
   block, the name is treated as global until it is assigned to; then it
   is treated as a local.

   The symbol table requires two passes to determine the scope of each name.
   The first pass collects raw facts from the AST via the symtable_visit_*
   functions: the name is a parameter here, the name is used but not defined
   here, etc.  The second pass analyzes these facts during a pass over the
   PySTEntryObjects created during pass 1.

   When a function is entered during the second pass, the parent passes
   the set of all name bindings visible to its children.  These bindings
   are used to determine if non-local variables are free or implicit globals.
   Names which are explicitly declared nonlocal must exist in this set of
   visible names - if they do not, a syntax error is raised. After doing
   the local analysis, it analyzes each of its child blocks using an
   updated set of name bindings.

   The children update the free variable set.  If a local variable is added to
   the free variable set by the child, the variable is marked as a cell.  The
   function object being defined must provide runtime storage for the variable
   that may outlive the function's frame.  Cell variables are removed from the
   free set before the analyze function returns to its parent.

   During analysis, the names are:
      symbols: dict mapping from symbol names to flag values (including offset scope values)
      scopes: dict mapping from symbol names to scope values (no offset)
      local: set of all symbol names local to the current scope
      bound: set of all symbol names local to a containing function scope
      free: set of all symbol names referenced but not bound in child scopes
      global: set of all symbol names explicitly declared as global
*/

/* Decide on scope of name, given flags.

   The namespace dictionaries may be modified to record information
   about the new name.  For example, a new global will add an entry to
   global.  A name that was global can be changed to local.
*/
func (st *SymTable) AnalyzeName(scopes Scopes, name string, symbol Symbol, bound, local, free, global StringSet) {
	flags := symbol.Flags
	if (flags & DefGlobal) != 0 {
		if (flags & DefParam) != 0 {
			st.panicSyntaxErrorLinenof(symbol.Lineno, symbol.ColOffset, "name '%s' is parameter and global", name)
		}
		if (flags & DefNonlocal) != 0 {
			st.panicSyntaxErrorLinenof(symbol.Lineno, symbol.ColOffset, "name '%s' is nonlocal and global", name)
		}
		scopes[name] = ScopeGlobalExplicit
		global.Add(name)
		if bound != nil {
			bound.Discard(name)
		}
		return
	}
	if (flags & DefNonlocal) != 0 {
		if (flags & DefParam) != 0 {
			st.panicSyntaxErrorLinenof(symbol.Lineno, symbol.ColOffset, "name '%s' is parameter and nonlocal", name)
		}
		if bound == nil {
			st.panicSyntaxErrorLinenof(symbol.Lineno, symbol.ColOffset, "nonlocal declaration not allowed at module level")
		}
		if !bound.Contains(name) {
			st.panicSyntaxErrorLinenof(symbol.Lineno, symbol.ColOffset, "no binding for nonlocal '%s' found", name)
		}
		scopes[name] = ScopeFree
		st.Free = true
		free.Add(name)
		return
	}
	if (flags & DefBound) != 0 {
		scopes[name] = ScopeLocal
		local.Add(name)
		global.Discard(name)
		return
	}
	/* If an enclosing block has a binding for this name, it
	   is a free variable rather than a global variable.
	   Note that having a non-NULL bound implies that the block
	   is nested.
	*/
	if bound != nil && bound.Contains(name) {
		scopes[name] = ScopeFree
		st.Free = true
		free.Add(name)
		return
	}
	/* If a parent has a global statement, then call it global
	   explicit?  It could also be global implicit.
	*/
	if global != nil && global.Contains(name) {
		scopes[name] = ScopeGlobalImplicit
		return
	}
	if st.Nested {
		st.Free = true
	}
	scopes[name] = ScopeGlobalImplicit
}

/* If a name is defined in free and also in locals, then this block
   provides the binding for the free variable.  The name should be
   marked CELL in this block and removed from the free list.

   Note that the current block's free variables are included in free.
   That's safe because no name can be free and local in the same scope.
*/
func AnalyzeCells(scopes Scopes, free StringSet) {
	for name, scope := range scopes {
		if scope != ScopeLocal {
			continue
		}
		if !free.Contains(name) {
			continue
		}
		/* Replace LOCAL with CELL for this name, and remove
		   from free. It is safe to replace the value of name
		   in the dict, because it will not cause a resize.
		*/
		scopes[name] = ScopeCell
		free.Discard(name)
	}
}

func (st *SymTable) DropClassFree(free StringSet) {
	res := free.Discard("__class__")
	if res {
		st.NeedsClassClosure = true
	}
}

/* Enter the final scope information into the st.Symbols dict.
 *
 * All arguments are dicts.  Modifies symbols, others are read-only.
 */
func (symbols Symbols) Update(scopes Scopes, bound, free StringSet, classflag bool) {
	/* Update scope information for all symbols in this scope */
	for name, symbol := range symbols {
		symbol.Scope = scopes[name]
		symbols[name] = symbol
	}

	/* Record not yet resolved free variables from children (if any) */
	for name := range free {
		// FIXME haven't managed to find a test case for this code
		// suspect a problem!

		/* Handle symbol that already exists in this scope */
		if symbol, ok := symbols[name]; ok {
			/* Handle a free variable in a method of
			   the class that has the same name as a local
			   or global in the class scope.
			*/
			if classflag && (symbol.Flags&(DefBound|DefGlobal)) != 0 {
				symbol.Flags |= DefFreeClass
				symbols[name] = symbol
			}
			/* It's a cell, or already free in this scope */
			continue
		}
		/* Handle global symbol */
		if !bound.Contains(name) {
			continue /* it's a global */
		}
		/* Propagate new free symbol up the lexical stack */
		symbols[name] = Symbol{
			Scope: ScopeFree,
			// FIXME Lineno: node.GetLineno(),
			// FIXME ColOffset: node.GetColOffset(),
		}
	}
}

/* Make final symbol table decisions for block of ste.

   Arguments:
   st -- current symtable entry (input/output)
   bound -- set of variables bound in enclosing scopes (input).  bound
       is nil for module blocks.
   free -- set of free variables in enclosed scopes (output)
   globals -- set of declared global variables in enclosing scopes (input)

   The implementation uses two mutually recursive functions,
   analyze_block() and analyze_child_block().  analyze_block() is
   responsible for analyzing the individual names defined in a block.
   analyze_child_block() prepares temporary namespace dictionaries
   used to evaluated nested blocks.

   The two functions exist because a child block should see the name
   bindings of its enclosing blocks, but those bindings should not
   propagate back to a parent block.
*/
func (st *SymTable) AnalyzeBlock(bound, free, global StringSet) {
	local := make(StringSet) // collect new names bound in block
	scopes := make(Scopes)   // collect scopes defined for each name

	/* Allocate new global and bound variable dictionaries.  These
	   dictionaries hold the names visible in nested blocks.  For
	   ClassBlocks, the bound and global names are initialized
	   before analyzing names, because class bindings aren't
	   visible in methods.  For other blocks, they are initialized
	   after names are analyzed.
	*/

	/* TODO(jhylton): Package these dicts in a struct so that we
	   can write reasonable helper functions?
	*/
	newglobal := make(StringSet)
	newfree := make(StringSet)
	newbound := make(StringSet)

	/* Class namespace has no effect on names visible in
	   nested functions, so populate the global and bound
	   sets to be passed to child blocks before analyzing
	   this one.
	*/
	if st.Type == ClassBlock {
		/* Pass down known globals */
		newglobal.Update(global)
		/* Pass down previously bound symbols */
		newbound.Update(bound)
	}

	for name, v := range st.Symbols {
		st.AnalyzeName(scopes, name, v, bound, local, free, global)
	}

	/* Populate global and bound sets to be passed to children. */
	if st.Type != ClassBlock {
		/* Add function locals to bound set */
		if st.Type == FunctionBlock {
			newbound.Update(local)
		}
		/* Pass down previously bound symbols */
		newbound.Update(bound)
		/* Pass down known globals */
		newglobal.Update(global)
	} else {
		/* Special-case __class__ */
		newbound.Add("__class__")
	}

	/* Recursively call analyze_child_block() on each child block.

	   newbound, newglobal now contain the names visible in
	   nested blocks.  The free variables in the children will
	   be collected in allfree.
	*/
	allfree := make(StringSet)
	for _, entry := range st.Children {
		entry.AnalyzeChildBlock(newbound, newfree, newglobal, allfree)
		/* Check if any children have free variables */
		if entry.Free || entry.ChildFree {
			st.ChildFree = true
		}
	}

	newfree.Update(allfree)

	/* Check if any local variables must be converted to cell variables */
	if st.Type == FunctionBlock {
		AnalyzeCells(scopes, newfree)
	} else if st.Type == ClassBlock {
		st.DropClassFree(newfree)
	}
	/* Records the results of the analysis in the symbol table entry */
	st.Symbols.Update(scopes, bound, newfree, st.Type == ClassBlock)

	free.Update(newfree)
}

func (st *SymTable) AnalyzeChildBlock(bound, free, global, child_free StringSet) {
	/* Copy the bound and global dictionaries.

	   These dictionary are used by all blocks enclosed by the
	   current block.  The analyze_block() call modifies these
	   dictionaries.
	*/
	temp_bound := bound.Copy()
	temp_free := free.Copy()
	temp_global := global.Copy()

	st.AnalyzeBlock(temp_bound, temp_free, temp_global)
	child_free.Update(temp_free)
}

// Analyze the SymTable
func (st *SymTable) Analyze() {
	free := make(StringSet)
	global := make(StringSet)
	st.AnalyzeBlock(nil, free, global)
}

// Return a sorted list of symbol names if the scope of a name matches
// either scopeType or flag is set
func (st *SymTable) Find(scopeType Scope, flag DefUseFlags) (out []string) {
	for name, v := range st.Symbols {
		if v.Scope == scopeType || (v.Flags&flag) != 0 {
			out = append(out, name)
		}
	}
	sort.Strings(out)
	return out
}
