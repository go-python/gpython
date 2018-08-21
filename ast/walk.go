// Copyright 2018 The go-python Authors.  All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package ast

import "fmt"

// Walk calls Visit on every node in the ast.  The parent is visited
// first, then the children
//
// Visit returns true to continue the walk
func Walk(ast Ast, Visit func(Ast) bool) {
	if ast == nil {
		return
	}
	if !Visit(ast) {
		return
	}

	// walk a single ast
	walk := func(ast Ast) {
		Walk(ast, Visit)
	}

	// walkStmts walks all the statements in the slice passed in
	walkStmts := func(stmts []Stmt) {
		for _, stmt := range stmts {
			walk(stmt)
		}
	}

	// walkExprs walks all the exprs in the slice passed in
	walkExprs := func(exprs []Expr) {
		for _, expr := range exprs {
			walk(expr)
		}
	}

	// walkComprehensions walks all the comprehensions in the slice passed in
	walkComprehensions := func(comprehensions []Comprehension) {
		for _, comprehension := range comprehensions {
			// Target Expr
			// Iter   Expr
			// Ifs    []Expr
			walk(comprehension.Target)
			walk(comprehension.Iter)
			walkExprs(comprehension.Ifs)
		}
	}

	switch node := ast.(type) {

	// Module nodes

	case *Module:
		// Body []Stmt
		walkStmts(node.Body)

	case *Interactive:
		// Body []Stmt
		walkStmts(node.Body)

	case *Expression:
		// Body Expr
		walk(node.Body)

	case *Suite:
		// Body []Stmt
		walkStmts(node.Body)

	// Statememt nodes

	case *FunctionDef:
		// Name          Identifier
		// Args          *Arguments
		// Body          []Stmt
		// DecoratorList []Expr
		// Returns       Expr
		if node.Args != nil {
			walk(node.Args)
		}
		walkStmts(node.Body)
		walkExprs(node.DecoratorList)
		walk(node.Returns)

	case *ClassDef:
		// Name          Identifier
		// Bases         []Expr
		// Keywords      []*Keyword
		// Starargs      Expr
		// Kwargs        Expr
		// Body          []Stmt
		// DecoratorList []Expr
		walkExprs(node.Bases)
		for _, k := range node.Keywords {
			walk(k)
		}
		walk(node.Starargs)
		walk(node.Kwargs)
		walkStmts(node.Body)
		walkExprs(node.DecoratorList)

	case *Return:
		// Value Expr
		walk(node.Value)

	case *Delete:
		// Targets []Expr
		walkExprs(node.Targets)

	case *Assign:
		// Targets []Expr
		// Value   Expr
		walkExprs(node.Targets)
		walk(node.Value)

	case *AugAssign:
		// Target Expr
		// Op     OperatorNumber
		// Value  Expr
		walk(node.Target)
		walk(node.Value)

	case *For:
		// Target Expr
		// Iter   Expr
		// Body   []Stmt
		// Orelse []Stmt
		walk(node.Target)
		walk(node.Iter)
		walkStmts(node.Body)
		walkStmts(node.Orelse)

	case *While:
		// Test   Expr
		// Body   []Stmt
		// Orelse []Stmt
		walk(node.Test)
		walkStmts(node.Body)
		walkStmts(node.Orelse)

	case *If:
		// Test   Expr
		// Body   []Stmt
		// Orelse []Stmt
		walk(node.Test)
		walkStmts(node.Body)
		walkStmts(node.Orelse)

	case *With:
		// Items []*WithItem
		// Body  []Stmt
		for _, wi := range node.Items {
			walk(wi)
		}
		walkStmts(node.Body)

	case *Raise:
		// Exc   Expr
		// Cause Expr
		walk(node.Exc)
		walk(node.Cause)

	case *Try:
		// Body      []Stmt
		// Handlers  []*ExceptHandler
		// Orelse    []Stmt
		// Finalbody []Stmt
		walkStmts(node.Body)
		for _, h := range node.Handlers {
			walk(h)
		}
		walkStmts(node.Orelse)
		walkStmts(node.Finalbody)

	case *Assert:
		// Test Expr
		// Msg  Expr
		walk(node.Test)
		walk(node.Msg)

	case *Import:
		// Names []*Alias
		for _, n := range node.Names {
			walk(n)
		}

	case *ImportFrom:
		// Module Identifier
		// Names  []*Alias
		// Level  int
		for _, n := range node.Names {
			walk(n)
		}

	case *Global:
		// Names []Identifier

	case *Nonlocal:
		// Names []Identifier

	case *ExprStmt:
		// Value Expr
		walk(node.Value)

	case *Pass:

	case *Break:

	case *Continue:

	// Expr nodes

	case *BoolOp:
		// Op     BoolOpNumber
		// Values []Expr
		walkExprs(node.Values)

	case *BinOp:
		// Left  Expr
		// Op    OperatorNumber
		// Right Expr
		walk(node.Left)
		walk(node.Right)

	case *UnaryOp:
		// Op      UnaryOpNumber
		// Operand Expr
		walk(node.Operand)

	case *Lambda:
		// Args *Arguments
		// Body Expr
		if node.Args != nil {
			walk(node.Args)
		}
		walk(node.Body)

	case *IfExp:
		// Test   Expr
		// Body   Expr
		// Orelse Expr
		walk(node.Test)
		walk(node.Body)
		walk(node.Orelse)

	case *Dict:
		// Keys   []Expr
		// Values []Expr
		walkExprs(node.Keys)
		walkExprs(node.Values)

	case *Set:
		// Elts []Expr
		walkExprs(node.Elts)

	case *ListComp:
		// Elt        Expr
		// Generators []Comprehension
		walk(node.Elt)
		walkComprehensions(node.Generators)

	case *SetComp:
		// Elt        Expr
		// Generators []Comprehension
		walk(node.Elt)
		walkComprehensions(node.Generators)

	case *DictComp:
		// Key        Expr
		// Value      Expr
		// Generators []Comprehension
		walk(node.Key)
		walk(node.Value)
		walkComprehensions(node.Generators)

	case *GeneratorExp:
		// Elt        Expr
		// Generators []Comprehension
		walk(node.Elt)
		walkComprehensions(node.Generators)

	case *Yield:
		// Value Expr
		walk(node.Value)

	case *YieldFrom:
		// Value Expr
		walk(node.Value)

	case *Compare:
		// Left        Expr
		// Ops         []CmpOp
		// Comparators []Expr
		walk(node.Left)
		walkExprs(node.Comparators)

	case *Call:
		// Func     Expr
		// Args     []Expr
		// Keywords []*Keyword
		// Starargs Expr
		// Kwargs   Expr
		walk(node.Func)
		walkExprs(node.Args)
		for _, k := range node.Keywords {
			walk(k)
		}
		walk(node.Starargs)
		walk(node.Kwargs)

	case *Num:
		// N Object

	case *Str:
		// S py.String

	case *Bytes:
		// S py.Bytes

	case *NameConstant:
		// Value Singleton

	case *Ellipsis:

	case *Attribute:
		// Value Expr
		// Attr  Identifier
		// Ctx   ExprContext
		walk(node.Value)

	case *Subscript:
		// Value Expr
		// Slice Slicer
		// Ctx   ExprContext
		walk(node.Value)
		walk(node.Slice)

	case *Starred:
		// Value Expr
		// Ctx   ExprContext
		walk(node.Value)

	case *Name:
		// Id  Identifier
		// Ctx ExprContext

	case *List:
		// Elts []Expr
		// Ctx  ExprContext
		walkExprs(node.Elts)

	case *Tuple:
		// Elts []Expr
		// Ctx  ExprContext
		walkExprs(node.Elts)

	// Slicer nodes

	case *Slice:
		// Lower Expr
		// Upper Expr
		// Step  Expr
		walk(node.Lower)
		walk(node.Upper)
		walk(node.Step)

	case *ExtSlice:
		// Dims []Slicer
		for _, s := range node.Dims {
			walk(s)
		}

	case *Index:
		// Value Expr
		walk(node.Value)

	// Misc nodes

	case *ExceptHandler:
		// ExprType Expr
		// Name     Identifier
		// Body     []Stmt
		walk(node.ExprType)
		walkStmts(node.Body)

	case *Arguments:
		// Args       []*Arg
		// Vararg     *Arg
		// Kwonlyargs []*Arg
		// KwDefaults []Expr
		// Kwarg      *Arg
		// Defaults   []Expr
		for _, arg := range node.Args {
			walk(arg)
		}
		if node.Vararg != nil {
			walk(node.Vararg)
		}
		for _, arg := range node.Kwonlyargs {
			walk(arg)
		}
		walkExprs(node.KwDefaults)
		if node.Kwarg != nil {
			walk(node.Kwarg)
		}
		walkExprs(node.Defaults)

	case *Arg:
		// Arg        Identifier
		// Annotation Expr
		if node.Annotation != nil {
			walk(node.Annotation)
		}

	case *Keyword:
		// Arg   Identifier
		// Value Expr
		walk(node.Value)

	case *Alias:
		// Name   Identifier
		// AsName Identifier

	case *WithItem:
		// ContextExpr  Expr
		// OptionalVars Expr
		walk(node.ContextExpr)
		walk(node.OptionalVars)

	default:
		panic(fmt.Sprintf("Unknown ast node %T, %#v", node, node))
	}
}
