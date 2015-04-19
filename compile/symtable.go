// Make a symbol table to identify which variables are glocal, closure or local

// FIXME tests for Varkeywords etc

// FIXME need to set locations for panics, eg
// PyErr_SyntaxLocationObject(st.ste_table.st_filename,
// 	st.ste_opt_lineno,
// 	st.ste_opt_col_offset)

package compile

import (
	"strings"

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

// Accumulate Scopes
type Scopes map[string]Scope

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

	// scopeGlobalExplicit and scopeGlobalImplicit are used internally by the symbol
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

// OptType for SymTable
type OptType uint8

// The following flag names are used for the Unoptimized bit field
const (
	optImportStar OptType = 1 << iota
	optTopLevel           // top-level names, including eval and exec
)

// Info about a symbol
type Symbol struct {
	Scope Scope
	Flags DefUse
}

type Symbols map[string]Symbol

func NewSymbols() Symbols {
	return make(Symbols)
}

type SymTable struct {
	Type              BlockType // 'class', 'module', and 'function'
	Name              string    // name of the class if the table is for a class, the name of the function if the table is for a function, or 'top' if the table is global (get_type() returns 'module').
	Lineno            int       // number of the first line in the block this table represents.
	Unoptimized       OptType   // false if namespace is optimized
	Nested            bool      // true if the block is a nested class or function.
	Free              bool      // true if block has free variables
	ChildFree         bool      // true if a child block has free vars, including free refs to globals
	Generator         bool      // true if namespace is a generator
	Varargs           bool      // true if block has varargs
	Varkeywords       bool      // true if block has varkeywords
	Returns_value     bool      // true if namespace uses return with an argument
	NeedsClassClosure bool      // for class scopes, true if a closure over __class__ should be created
	col_offset        int       // offset of first line of block
	opt_lineno        int       // lineno of last exec or import *
	opt_col_offset    int       // offset of last exec or import *
	tmpname           int       // counter for listcomp temp vars

	Symbols  Symbols
	Global   *SymTable // symbol table entry for module
	Parent   *SymTable
	Varnames []string             // list of function parameters
	Children map[string]*SymTable // Child SymTables keyed by symbol name
}

// Make a new top symbol table from the ast supplied
func NewSymTable(Ast ast.Ast) *SymTable {
	st := newSymTable(ModuleBlock, "top", nil)
	st.Unoptimized = optTopLevel
	// Parse into the symbol table
	st.Parse(Ast)
	// Analyze the symbolt table
	st.Analyze()
	return st
}

// Make a new symbol table from the ast supplied of the given type
func newSymTable(Type BlockType, Name string, parent *SymTable) *SymTable {
	st := &SymTable{
		Type:     Type,
		Name:     Name,
		Parent:   parent,
		Symbols:  NewSymbols(),
		Children: make(map[string]*SymTable),
	}
	if parent == nil {
		st.Global = st
	} else {
		st.Global = parent.Global
		st.Nested = parent.Nested || (parent.Type == FunctionBlock)
	}
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
			// Special-case super: it counts as a use of __class__
			if node.Ctx == ast.Load && st.Type == FunctionBlock && node.Id == "super" {
				st.AddDef(ast.Identifier("__class__"), defUse)
			}
		case *ast.FunctionDef:
			// Add the function name to the SymTable
			st.AddDef(node.Name, defLocal)
			name := string(node.Name)

			// Make a new symtable and add it to parent
			stNew := newSymTable(FunctionBlock, name, st)
			st.Children[name] = stNew
			// FIXME set stNew.Lineno

			// Walk the Decorators and Returns in this Symtable
			for _, expr := range node.DecoratorList {
				st.Parse(expr)
			}
			st.Parse(node.Returns)

			// Walk the Args and Body in the new symtable
			if node.Args != nil {
				stNew.Parse(node.Args)
			}
			for _, stmt := range node.Body {
				stNew.Parse(stmt)
			}

			// return false to stop the parse
			return false
		case *ast.ClassDef:
			st.AddDef(node.Name, defLocal)
		case *ast.Lambda:

			// Comprehensions
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

		case *ast.Arguments:
			// skip default arguments inside function block
			// XXX should ast be different?
			for _, arg := range node.Args {
				st.AddDef(arg.Arg, defParam)
			}
			for _, arg := range node.Kwonlyargs {
				st.AddDef(arg.Arg, defParam)
			}
			if node.Vararg != nil {
				st.AddDef(node.Vararg.Arg, defParam)
				st.Varargs = true
			}
			if node.Kwarg != nil {
				st.AddDef(node.Kwarg.Arg, defParam)
				st.Varkeywords = true
			}

		case *ast.ExceptHandler:
			if node.Name != "" {
				st.AddDef(node.Name, defLocal)
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
				store_name = name[dot+1:]
			}
			if name != "*" {
				st.AddDef(store_name, defImport)
			} else {
				if st.Type != ModuleBlock {
					panic(py.ExceptionNewf(py.SyntaxError, "import * only allowed at module level"))
				}
				st.Unoptimized |= optImportStar
			}
		}
		return true
	})
}

func (st *SymTable) parseComprehension(Ast ast.Ast, scope_name ast.Identifier, generators []ast.Comprehension, elt ast.Expr, value ast.Expr) {
	/* FIXME
	_, is_generator := Ast.(*ast.GeneratorExp)
	needs_tmp := !is_generator
	outermost := generators[0]
	// Outermost iterator is evaluated in current scope
	st.Parse(outermost.Iter)
	// Create comprehension scope for the rest
	if scope_name == "" || !symtable_enter_block(st, scope_name, FunctionBlock, e, e.lineno, e.col_offset) {
		return 0
	}
	st.st_cur.ste_generator = is_generator
	// Outermost iter is received as an argument
	id := ast.Identifier(fmt.Sprintf(".%d", pos))
	st.AddDef(id, defParam)
	// Allocate temporary name if needed
	if needs_tmp {
		symtable_new_tmpname(st)
	}
	VISIT(st, expr, outermost.target)
	parseSeq(st, expr, outermost.ifs)
	parseSeq_tail(st, comprehension, generators, 1)
	if value {
		VISIT(st, expr, value)
	}
	VISIT(st, expr, elt)
	return symtable_exit_block(st, e)
	*/
}

const duplicateArgument = "duplicate argument %q in function definition"

// Add a symbol into the symble table
func (st *SymTable) AddDef(name ast.Identifier, flags DefUse) {
	// FIXME mangled := _Py_Mangle(st.st_private, name)
	mangled := string(name)

	// Add or update the symbol in the Symbols
	if sym, ok := st.Symbols[mangled]; ok {
		if (flags&defParam) != 0 && (sym.Flags&defParam) != 0 {
			// Is it better to use 'mangled' or 'name' here?
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
			Scope: 0, // FIXME
			Flags: flags,
		}
	}

	if (flags & defParam) != 0 {
		st.Varnames = append(st.Varnames, mangled)
	} else if (flags & defGlobal) != 0 {
		// If it is a global definition then add it in the global Symbols
		//
		// XXX need to update defGlobal for other flags too;
		// perhaps only defFreeClass
		if sym, ok := st.Global.Symbols[mangled]; ok {
			sym.Flags |= flags
			st.Global.Symbols[mangled] = sym
		} else {
			st.Global.Symbols[mangled] = Symbol{
				Scope: 0, // FIXME
				Flags: flags,
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
func (st *SymTable) AnalyzeName(scopes Scopes, name string, flags DefUse, bound, local, free, global StringSet) {
	if (flags & defGlobal) != 0 {
		if (flags & defParam) != 0 {
			panic(py.ExceptionNewf(py.SyntaxError, "name '%s' is parameter and global", name))
		}
		if (flags & defNonlocal) != 0 {
			panic(py.ExceptionNewf(py.SyntaxError, "name '%s' is nonlocal and global", name))
		}
		scopes[name] = scopeGlobalExplicit
		global.Add(name)
		if bound != nil {
			bound.Discard(name)
		}
		return
	}
	if (flags & defNonlocal) != 0 {
		if (flags & defParam) != 0 {
			panic(py.ExceptionNewf(py.SyntaxError, "name '%s' is parameter and nonlocal", name))
		}
		if bound == nil {
			panic(py.ExceptionNewf(py.SyntaxError, "nonlocal declaration not allowed at module level"))
		}
		if bound.Contains(name) {
			panic(py.ExceptionNewf(py.SyntaxError, "no binding for nonlocal '%s' found", name))
		}
		scopes[name] = scopeFree
		st.Free = true
		free.Add(name)
		return
	}
	if (flags & defBound) != 0 {
		scopes[name] = scopeLocal
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
		scopes[name] = scopeFree
		st.Free = true
		free.Add(name)
		return
	}
	/* If a parent has a global statement, then call it global
	   explicit?  It could also be global implicit.
	*/
	if global != nil && global.Contains(name) {
		scopes[name] = scopeGlobalImplicit
		return
	}
	if st.Nested {
		st.Free = true
	}
	scopes[name] = scopeGlobalImplicit
}

/* If a name is defined in free and also in locals, then this block
   provides the binding for the free variable.  The name should be
   marked CELL in this block and removed from the free list.

   Note that the current block's free variables are included in free.
   That's safe because no name can be free and local in the same scope.
*/
func AnalyzeCells(scopes Scopes, free StringSet) {
	for name, scope := range scopes {
		if scope != scopeLocal {
			continue
		}
		if !free.Contains(name) {
			continue
		}
		/* Replace LOCAL with CELL for this name, and remove
		   from free. It is safe to replace the value of name
		   in the dict, because it will not cause a resize.
		*/
		scopes[name] = scopeCell
		free.Discard(name)
	}
}

func (st *SymTable) DropClassFree(free StringSet) {
	res := free.Discard("__class__")
	if res {
		st.NeedsClassClosure = true
	}
}

/* Check for illegal statements in unoptimized namespaces */
func (st *SymTable) CheckUnoptimized() {
	if st.Type != FunctionBlock || st.Unoptimized == 0 || !(st.Free || st.ChildFree) {
		return
	}

	trailer := "contains a nested function with free variables"
	if !st.ChildFree {
		trailer = "is a nested function"
	}

	switch st.Unoptimized {
	case optTopLevel: /* import * at top-level is fine */
		return
	case optImportStar:
		panic(py.ExceptionNewf(py.SyntaxError, "import * is not allowed in function '%s' because it %s", st.Name, trailer))
		break
	}
}

/* Enter the final scope information into the ste_symbols dict.
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
		/* Handle symbol that already exists in this scope */
		if symbol, ok := symbols[name]; ok {
			/* Handle a free variable in a method of
			   the class that has the same name as a local
			   or global in the class scope.
			*/
			if classflag && (symbol.Flags&(defBound|defGlobal)) != 0 {
				symbol.Flags |= defFreeClass
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
		symbols[name] = Symbol{Scope: scopeFree}
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
	// PyObject *name, *v, *local = nil, *scopes = nil, *newbound = nil;
	// PyObject *newglobal = nil, *newfree = nil, *allfree = nil;
	// PyObject *temp;
	// int i, success = 0;
	// Py_ssize_t pos = 0;

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
		st.AnalyzeName(scopes, name, v.Flags, bound, local, free, global)
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
	st.CheckUnoptimized()

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
