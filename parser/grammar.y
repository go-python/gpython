%{

package parser

// Grammar for Python

import (
	"fmt"
	"github.com/ncw/gpython/py"
	"github.com/ncw/gpython/ast"
)

// NB can put code blocks in not just at the end

// Returns a Tuple if > 1 items or a trailing comma, otherwise returns
// the first item in elts
func tupleOrExpr(pos ast.Pos, elts []ast.Expr, optional_comma bool) ast.Expr {
	if optional_comma || len(elts) > 1 {
		return &ast.Tuple{ExprBase: ast.ExprBase{Pos: pos}, Elts: elts, Ctx: ast.Load}
	} else {
		return  elts[0]
	}
}

// Apply trailers (if any) to expr
//
// trailers are half made Call, Subscript or Attribute
func applyTrailers(expr ast.Expr, trailers []ast.Expr) ast.Expr {
	//trailers := $1
	for _, trailer := range trailers {
		switch x := trailer.(type) {
		case *ast.Call:
			x.Func, expr = expr, x
		case *ast.Subscript:
			x.Value, expr = expr, x
		case *ast.Attribute:
			x.Value, expr = expr, x
		default:
			panic(fmt.Sprintf("Unknown trailer type: %T", expr))
		}
	}
	return expr
}

// Set the context for expr
func setCtx(yylex yyLexer, expr ast.Expr, ctx ast.ExprContext) {
	setctxer, ok := expr.(ast.SetCtxer)
	if !ok {
		expr_name := ""
		switch expr.(type) {
		case *ast.Lambda:
			expr_name = "lambda"
		case *ast.Call:
			expr_name = "function call"
		case *ast.BoolOp, *ast.BinOp, *ast.UnaryOp:
			expr_name = "operator"
		case *ast.GeneratorExp:
			expr_name = "generator expression"
		case *ast.Yield, *ast.YieldFrom:
			expr_name = "yield expression"
		case *ast.ListComp:
			expr_name = "list comprehension"
		case *ast.SetComp:
			expr_name = "set comprehension"
		case *ast.DictComp:
			expr_name = "dict comprehension"
		case *ast.Dict, *ast.Set, *ast.Num, *ast.Str, *ast.Bytes:
			expr_name = "literal"
		case *ast.NameConstant:
			expr_name = "keyword"
		case *ast.Ellipsis:
			expr_name = "Ellipsis"
		case *ast.Compare:
			expr_name = "comparison"
		case *ast.IfExp:
			expr_name = "conditional expression"
		default:
			expr_name = fmt.Sprintf("unexpected %T", expr)
		}
		action := "assign to"
		if ctx == ast.Del {
			action = "delete"
		}
		yylex.(*yyLex).SyntaxErrorf("can't %s %s", action, expr_name)
		return
	}
	setctxer.SetCtx(ctx)
}

// Set the context for all the items in exprs
func setCtxs(yylex yyLexer, exprs []ast.Expr, ctx ast.ExprContext) {
	for i := range exprs {
		setCtx(yylex, exprs[i], ctx)
	}
}

%}

%union {
	pos		ast.Pos		// kept up to date by the lexer
	str		string
	obj		py.Object
	mod		ast.Mod
	stmt		ast.Stmt
	stmts		[]ast.Stmt
	expr		ast.Expr
	exprs		[]ast.Expr
	op		ast.OperatorNumber
	cmpop		ast.CmpOp
	comma		bool
	comprehensions	[]ast.Comprehension
	isExpr		bool
	slice		ast.Slicer
	call		*ast.Call
	level		int
	alias		*ast.Alias
	aliases		[]*ast.Alias
	identifiers	[]ast.Identifier
	ifstmt		*ast.If
	lastif		*ast.If
	exchandlers	[]*ast.ExceptHandler
	withitem	*ast.WithItem
	withitems	[]*ast.WithItem
	arg		*ast.Arg
	args		[]*ast.Arg
	arguments	*ast.Arguments
}

%type <obj> strings
%type <mod> inputs file_input single_input eval_input
%type <stmts> simple_stmt stmt nl_or_stmt small_stmts stmts suite optional_else
%type <stmt> compound_stmt small_stmt expr_stmt del_stmt pass_stmt flow_stmt import_stmt global_stmt nonlocal_stmt assert_stmt break_stmt continue_stmt return_stmt raise_stmt yield_stmt import_name import_from while_stmt if_stmt for_stmt try_stmt with_stmt funcdef classdef classdef_or_funcdef decorated
%type <op> augassign
%type <expr> expr_or_star_expr expr star_expr xor_expr and_expr shift_expr arith_expr term factor power trailer atom test_or_star_expr test not_test lambdef test_nocond lambdef_nocond or_test and_test comparison testlist testlist_star_expr yield_expr_or_testlist yield_expr yield_expr_or_testlist_star_expr dictorsetmaker sliceop except_clause optional_return_type decorator
%type <exprs> exprlist testlistraw comp_if comp_iter expr_or_star_exprs test_or_star_exprs tests test_colon_tests trailers equals_yield_expr_or_testlist_star_expr decorators
%type <cmpop> comp_op
%type <comma> optional_comma
%type <comprehensions> comp_for
%type <slice> subscript subscriptlist subscripts
%type <call> argument arguments optional_arguments arguments2 arglist optional_arglist_call optional_arglist
%type <level> dot dots
%type <str> dotted_name from_arg
%type <identifiers> names
%type <alias> dotted_as_name import_as_name
%type <aliases> dotted_as_names import_as_names import_from_arg
%type <ifstmt> elifs
%type <exchandlers> except_clauses
%type <withitem> with_item
%type <withitems> with_items
%type <arg> vfpdeftest vfpdef optional_vfpdef tfpdeftest tfpdef optional_tfpdef
%type <args> vfpdeftests vfpdeftests1 tfpdeftests tfpdeftests1
%type <arguments> varargslist parameters optional_typedargslist typedargslist

%token NEWLINE
%token ENDMARKER
%token <str> NAME
%token INDENT
%token DEDENT
%token <obj> STRING
%token <obj> NUMBER

%token PLINGEQ // !=
%token PERCEQ // %=
%token ANDEQ // &=
%token STARSTAR // **
%token STARSTAREQ // **=
%token STAREQ // *=
%token PLUSEQ // +=
%token MINUSEQ // -=
%token MINUSGT // ->
%token ELIPSIS // ...
%token DIVDIV // //
%token DIVDIVEQ // //=
%token DIVEQ // /=
%token LTLT // <<
%token LTLTEQ // <<=
%token LTEQ // <=
%token LTGT // <>
%token EQEQ // ==
%token GTEQ // >=
%token GTGT // >>
%token GTGTEQ // >>=
%token HATEQ // ^=
%token PIPEEQ // |=

%token FALSE // False
%token NONE // None
%token TRUE // True
%token AND // and
%token AS // as
%token ASSERT // assert
%token BREAK // break
%token CLASS // class
%token CONTINUE // continue
%token DEF // def
%token DEL // del
%token ELIF // elif
%token ELSE // else
%token EXCEPT // except
%token FINALLY // finally
%token FOR // for
%token FROM // from
%token GLOBAL // global
%token IF // if
%token IMPORT // import
%token IN // in
%token IS // is
%token LAMBDA // lambda
%token NONLOCAL // nonlocal
%token NOT // not
%token OR // or
%token PASS // pass
%token RAISE // raise
%token RETURN // return
%token TRY // try
%token WHILE // while
%token WITH // with
%token YIELD // yield

%token '(' ')' '[' ']' ':' ',' ';' '+' '-' '*' '/' '|' '&' '<' '>' '=' '.' '%' '{' '}' '^' '~' '@'

%token SINGLE_INPUT FILE_INPUT EVAL_INPUT

// Note:  Changing the grammar specified in this file will most likely
//        require corresponding changes in the parser module
//        (../Modules/parsermodule.c).  If you can't make the changes to
//        that module yourself, please co-ordinate the required changes
//        with someone who can; ask around on python-dev for help.  Fred
//        Drake <fdrake@acm.org> will probably be listening there.

// NOTE WELL: You should also follow all the steps listed in PEP 306,
// "How to Change Python's Grammar"

// Start symbols for the grammar:
//       single_input is a single interactive statement;
//       file_input is a module or sequence of commands read from an input file;
//       eval_input is the input for the eval() functions.
// NB: compound_stmt in single_input is followed by extra NEWLINE!

%%

// Start of grammar. This has 3 pseudo tokens which say which
// direction through the rest of the grammar we take.
inputs:
	SINGLE_INPUT single_input
	{
		yylex.(*yyLex).mod = $2
		return 0
	}
|	FILE_INPUT file_input
	{
		yylex.(*yyLex).mod = $2
		return 0
	}
|	EVAL_INPUT eval_input
	{
		yylex.(*yyLex).mod = $2
		return 0
	}

single_input:
/*	NEWLINE
	{
		// This is in the python grammar, but the interpreter
		// just gives "unexpected EOF while parsing" when you
		// give it a \n
		$$ = &ast.Interactive{ModBase: ast.ModBase{Pos: $<pos>$}}
	}
|*/	simple_stmt
	{
		$$ = &ast.Interactive{ModBase: ast.ModBase{Pos: $<pos>$}, Body: $1}
	}
|	compound_stmt NEWLINE
	{
		//  NB: compound_stmt in single_input is followed by extra NEWLINE!
		$$ = &ast.Interactive{ModBase: ast.ModBase{Pos: $<pos>$}, Body: []ast.Stmt{$1}}
	}

//file_input: (NEWLINE | stmt)* ENDMARKER
file_input:
	nl_or_stmt ENDMARKER
	{
		$$ = &ast.Module{ModBase: ast.ModBase{Pos: $<pos>$}, Body: $1}
	}

// (NEWLINE | stmt)*
nl_or_stmt:
	{
		$$ = nil
	}
|	nl_or_stmt NEWLINE
	{
	}
|	nl_or_stmt stmt
	{
		$$ = append($$, $2...)
	}

//eval_input: testlist NEWLINE* ENDMARKER
eval_input:
	testlist nls ENDMARKER
	{
		$$ = &ast.Expression{ModBase: ast.ModBase{Pos: $<pos>$}, Body: $1}
	}

// NEWLINE*
nls:
|	nls NEWLINE

optional_arglist:
	{
		$$ = &ast.Call{ExprBase: ast.ExprBase{Pos: $<pos>$}}
	}
|	arglist
	{
		$$ = $1
	}

optional_arglist_call:
	{
		$$ = nil
	}
|	'(' optional_arglist ')'
	{
		$$ = $2
	}

decorator:
	'@' dotted_name optional_arglist_call NEWLINE
	{
		fn := &ast.Name{ExprBase: ast.ExprBase{Pos: $<pos>$}, Id: ast.Identifier($2), Ctx: ast.Load}
		if $3 == nil {
			$$ = fn
		} else {
			call := *$3
			call.Func = fn
			$$ = &call
		}
	}

decorators:
	decorator
	{
		$$ = nil
		$$ = append($$, $1)
	}
|	decorators decorator
	{
		$$ = append($$, $2)
	}

classdef_or_funcdef:
	classdef
	{
		$$ = $1
	}
|	funcdef
	{
		$$ = $1
	}

decorated:
	decorators classdef_or_funcdef
	{
		switch x := ($2).(type) {
		case *ast.ClassDef:
			x.DecoratorList = $1
			$$ = x
		case *ast.FunctionDef:
			x.DecoratorList = $1
			$$ = x
		default:
			panic("bad type for decorated")
		}
	}

optional_return_type:
	{
		$$ = nil
	}
|	MINUSGT test
	{
		$$ = $2
	}

funcdef:
	DEF NAME parameters optional_return_type ':' suite
	{
		$$ = &ast.FunctionDef{StmtBase: ast.StmtBase{Pos: $<pos>$}, Name: ast.Identifier($2), Args: $3, Body: $6, Returns: $4}
	}

parameters:
	'(' optional_typedargslist ')'
	{
		$$ = $2
	}

optional_typedargslist:
	{
		$$ = &ast.Arguments{Pos: $<pos>$}
	}
|	typedargslist
	{
		$$ = $1
	}

// (',' tfpdef ['=' test])*
tfpdeftest:
	tfpdef
	{
		$$ = $1
		$<expr>$ = nil
	}
|	tfpdef '=' test
	{
		$$ = $1
		$<expr>$ = $3
	}

tfpdeftests:
	{
		$$ = nil
		$<exprs>$ = nil
	}
|	tfpdeftests ',' tfpdeftest
	{
		$$ = append($$, $3)
		if $<expr>3 != nil {
			$<exprs>$ = append($<exprs>$, $<expr>3)
		}
	}

tfpdeftests1:
	tfpdeftest
	{
		$$ = nil
		$$ = append($$, $1)
		$<exprs>$ = nil
		if $<expr>1 != nil {
			$<exprs>$ = append($<exprs>$, $<expr>1)
		}
	}
|	tfpdeftests1 ',' tfpdeftest
	{
		$$ = append($$, $3)
		if $<expr>3 != nil {
			$<exprs>$ = append($<exprs>$, $<expr>3)
		}
	}

optional_tfpdef:
	{
		$$ = nil
	}
|	tfpdef
	{
		$$ = $1
	}

// FIXME this isn't checking all the python rules for args before kwargs etc
typedargslist: 
	tfpdeftests1 optional_comma
	{
		$$ = &ast.Arguments{Pos: $<pos>$, Args: $1, Defaults: $<exprs>1}
	}
|	tfpdeftests1 ',' '*' optional_tfpdef tfpdeftests
	{
		$$ = &ast.Arguments{Pos: $<pos>$, Args: $1, Defaults: $<exprs>1, Vararg: $4, Kwonlyargs: $5, KwDefaults: $<exprs>5}
	}
|	tfpdeftests1 ',' '*' optional_tfpdef tfpdeftests ',' STARSTAR tfpdef
	{
		$$ = &ast.Arguments{Pos: $<pos>$, Args: $1, Defaults: $<exprs>1, Vararg: $4, Kwonlyargs: $5, KwDefaults: $<exprs>5, Kwarg: $8}
	}
|	tfpdeftests1 ',' STARSTAR tfpdef
	{
		$$ = &ast.Arguments{Pos: $<pos>$, Args: $1, Defaults: $<exprs>1, Kwarg: $4}
	}
|	'*' optional_tfpdef tfpdeftests
	{
		$$ = &ast.Arguments{Pos: $<pos>$, Vararg: $2, Kwonlyargs: $3, KwDefaults: $<exprs>3}
	}
|	'*' optional_tfpdef tfpdeftests ',' STARSTAR tfpdef
	{
		$$ = &ast.Arguments{Pos: $<pos>$, Vararg: $2, Kwonlyargs: $3, KwDefaults: $<exprs>3, Kwarg: $6}
	}
|	STARSTAR tfpdef
	{
		$$ = &ast.Arguments{Pos: $<pos>$, Kwarg: $2}
	}

tfpdef:
	NAME
	{
		$$ = &ast.Arg{Pos: $<pos>$, Arg: ast.Identifier($1)}
	}
|	NAME ':' test
	{
		$$ = &ast.Arg{Pos: $<pos>$, Arg: ast.Identifier($1), Annotation: $3}
	}

vfpdeftest:
	vfpdef
	{
		$$ = $1
		$<expr>$ = nil
	}
|	vfpdef '=' test
	{
		$$ = $1
		$<expr>$ = $3
	}

vfpdeftests:
	{
		$$ = nil
		$<exprs>$ = nil
	}
|	vfpdeftests ',' vfpdeftest
	{
		$$ = append($$, $3)
		if $<expr>3 != nil {
			$<exprs>$ = append($<exprs>$, $<expr>3)
		}
	}

vfpdeftests1:
	vfpdeftest
	{
		$$ = nil
		$$ = append($$, $1)
		$<exprs>$ = nil
		if $<expr>1 != nil {
			$<exprs>$ = append($<exprs>$, $<expr>1)
		}
	}
|	vfpdeftests1 ',' vfpdeftest
	{
		$$ = append($$, $3)
		if $<expr>3 != nil {
			$<exprs>$ = append($<exprs>$, $<expr>3)
		}
	}

optional_vfpdef:
	{
		$$ = nil
	}
|	vfpdef
	{
		$$ = $1
	}

// FIXME this isn't checking all the python rules for args before kwargs etc
varargslist:
	vfpdeftests1 optional_comma
	{
		$$ = &ast.Arguments{Pos: $<pos>$, Args: $1, Defaults: $<exprs>1}
	}
|	vfpdeftests1 ',' '*' optional_vfpdef vfpdeftests
	{
		$$ = &ast.Arguments{Pos: $<pos>$, Args: $1, Defaults: $<exprs>1, Vararg: $4, Kwonlyargs: $5, KwDefaults: $<exprs>5}
	}
|	vfpdeftests1 ',' '*' optional_vfpdef vfpdeftests ',' STARSTAR vfpdef
	{
		$$ = &ast.Arguments{Pos: $<pos>$, Args: $1, Defaults: $<exprs>1, Vararg: $4, Kwonlyargs: $5, KwDefaults: $<exprs>5, Kwarg: $8}
	}
|	vfpdeftests1 ',' STARSTAR vfpdef
	{
		$$ = &ast.Arguments{Pos: $<pos>$, Args: $1, Defaults: $<exprs>1, Kwarg: $4}
	}
|	'*' optional_vfpdef vfpdeftests
	{
		$$ = &ast.Arguments{Pos: $<pos>$, Vararg: $2, Kwonlyargs: $3, KwDefaults: $<exprs>3}
	}
|	'*' optional_vfpdef vfpdeftests ',' STARSTAR vfpdef
	{
		$$ = &ast.Arguments{Pos: $<pos>$, Vararg: $2, Kwonlyargs: $3, KwDefaults: $<exprs>3, Kwarg: $6}
	}
|	STARSTAR vfpdef
	{
		$$ = &ast.Arguments{Pos: $<pos>$, Kwarg: $2}
	}

vfpdef:
	NAME
	{
		$$ = &ast.Arg{Pos: $<pos>$, Arg: ast.Identifier($1)}
	}

stmt:
	simple_stmt
	{
		$$ = $1
	}
|	compound_stmt
	{
		$$ = []ast.Stmt{$1}
	}

optional_semicolon: | ';'

small_stmts:
	small_stmt
	{
		$$ = nil
		$$ = append($$, $1)
	}
|	small_stmts ';' small_stmt
	{
		$$ = append($$, $3)
	}

simple_stmt:
	small_stmts optional_semicolon NEWLINE
	{
		$$ = $1
	}

small_stmt:
	expr_stmt
	{
		$$ = $1
	}
|	del_stmt
	{
		$$ = $1
	}
|	pass_stmt
	{
		$$ = $1
	}
|	flow_stmt
	{
		$$ = $1
	}
|	import_stmt
	{
		$$ = $1
	}
|	global_stmt
	{
		$$ = $1
	}
|	nonlocal_stmt
	{
		$$ = $1
	}
|	assert_stmt
	{
		$$ = $1
	}

/*
expr_stmt: testlist_star_expr (augassign (yield_expr|testlist) |
                     ('=' (yield_expr|testlist_star_expr))*)

expr_stmt:
testlist_star_expr (
    augassign (
        yield_expr|testlist
    ) | (
        '=' (
            yield_expr|testlist_star_expr
        )
    )*
)


expr_stmt: testlist_star_expr augassign yield_expr
expr_stmt: testlist_star_expr augassign testlist
expr_stmt: testlist_star_expr ('=' (yield_expr|testlist_star_expr))*
*/

expr_stmt:
	testlist_star_expr augassign yield_expr_or_testlist
	{
		target := $1
		setCtx(yylex, target, ast.Store)
		$$ = &ast.AugAssign{StmtBase: ast.StmtBase{Pos: $<pos>$}, Target: target, Op: $2, Value: $3}
	}
|	testlist_star_expr equals_yield_expr_or_testlist_star_expr
	{
		targets := []ast.Expr{$1}
		targets = append(targets, $2...)
		value := targets[len(targets)-1]
		targets = targets[:len(targets)-1]
		setCtxs(yylex, targets, ast.Store)
		$$ = &ast.Assign{StmtBase: ast.StmtBase{Pos: $<pos>$}, Targets: targets, Value: value}
	}
|	testlist_star_expr
	{
		$$ = &ast.ExprStmt{StmtBase: ast.StmtBase{Pos: $<pos>$}, Value: $1}
	}

yield_expr_or_testlist:
	yield_expr
	{
		$$ = $1
	}
|	testlist
	{
		$$ = $1
	}

yield_expr_or_testlist_star_expr:
	yield_expr
	{
		$$ = $1
	}
|	testlist_star_expr
	{
		$$ = $1
	}

equals_yield_expr_or_testlist_star_expr:
	'=' yield_expr_or_testlist_star_expr
	{
		$$ = nil
		$$ = append($$, $2)
	}
|	equals_yield_expr_or_testlist_star_expr '=' yield_expr_or_testlist_star_expr
	{
		$$ = append($$, $3)
	}

test_or_star_exprs:
	test_or_star_expr
	{
		$$ = nil
		$$ = append($$, $1)
	}
|	test_or_star_exprs ',' test_or_star_expr
	{
		$$ = append($$, $3)
	}

test_or_star_expr:
	test
	{
		$$ = $1
	}
|	star_expr
	{
		$$ = $1
	}

optional_comma:
	{
		$$ = false
	}
|	','
	{
		$$ = true
	}

testlist_star_expr:
	test_or_star_exprs optional_comma
	{
		$$ = tupleOrExpr($<pos>$, $1, $2)
	}

augassign:
	PLUSEQ
	{
		$$ = ast.Add
	}
|	MINUSEQ
	{
		$$ = ast.Sub
	}
|	STAREQ
	{
		$$ = ast.Mult
	}
|	DIVEQ
	{
		$$ = ast.Div
	}
|	PERCEQ
	{
		$$ = ast.Modulo
	}
|	ANDEQ
	{
		$$ = ast.BitAnd
	}
|	PIPEEQ
	{
		$$ = ast.BitOr
	}
|	HATEQ
	{
		$$ = ast.BitXor
	}
|	LTLTEQ
	{
		$$ = ast.LShift
	}
|	GTGTEQ
	{
		$$ = ast.RShift
	}
|	STARSTAREQ
	{
		$$ = ast.Pow
	}
|	DIVDIVEQ
	{
		$$ = ast.FloorDiv
	}

// For normal assignments, additional restrictions enforced by the interpreter
del_stmt:
	DEL exprlist
	{
		setCtxs(yylex, $2, ast.Del)
		$$ = &ast.Delete{StmtBase: ast.StmtBase{Pos: $<pos>$}, Targets: $2}
	}

pass_stmt:
	PASS
	{
		$$ = &ast.Pass{StmtBase: ast.StmtBase{Pos: $<pos>$}}
	}

flow_stmt:
	break_stmt
	{
		$$ = $1
	}
|	continue_stmt
	{
		$$ = $1
	}
|	return_stmt
	{
		$$ = $1
	}
|	raise_stmt
	{
		$$ = $1
	}
|	yield_stmt
	{
		$$ = $1
	}

break_stmt:
	BREAK
	{
		$$ = &ast.Break{StmtBase: ast.StmtBase{Pos: $<pos>$}}
	}

continue_stmt:
	CONTINUE
	{
		$$ = &ast.Continue{StmtBase: ast.StmtBase{Pos: $<pos>$}}
	}

return_stmt:
	RETURN
	{
		$$ = &ast.Return{StmtBase: ast.StmtBase{Pos: $<pos>$}}
	}
|	RETURN testlist
	{
		$$ = &ast.Return{StmtBase: ast.StmtBase{Pos: $<pos>$}, Value: $2}
	}

yield_stmt:
	yield_expr
	{
		$$ = &ast.ExprStmt{StmtBase: ast.StmtBase{Pos: $<pos>$}, Value: $1}
	}

raise_stmt:
	RAISE
	{
		$$ = &ast.Raise{StmtBase: ast.StmtBase{Pos: $<pos>$}}
	}
|	RAISE test
	{
		$$ = &ast.Raise{StmtBase: ast.StmtBase{Pos: $<pos>$}, Exc: $2}
	}
|	RAISE test FROM test
	{
		$$ = &ast.Raise{StmtBase: ast.StmtBase{Pos: $<pos>$}, Exc: $2, Cause: $4}
	}

import_stmt:
	import_name
	{
		$$ = $1
	}
|	import_from
	{
		$$ = $1
	}

import_name:
	IMPORT dotted_as_names
	{
		$$ = &ast.Import{StmtBase: ast.StmtBase{Pos: $<pos>$}, Names: $2}
	}

// note below: the '.' | ELIPSIS is necessary because '...' is tokenized as ELIPSIS
dot:
	'.'
	{
		$$ = 1
	}
|	ELIPSIS
	{
		$$ = 3
	}

dots:
	dot
	{
		$$ = $1
	}
|	dots dot
	{
		$$ += $2
	}

from_arg:
	dotted_name
	{
		$<level>$ = 0
		$$ = $1
	}
|	dots dotted_name
	{
		$<level>$ = $1
		$$ = $2
	}
|	dots
	{
		$<level>$ = $1
		$$ = ""
	}

import_from_arg:
	'*'
	{
		$$ = []*ast.Alias{&ast.Alias{Pos: $<pos>$, Name: ast.Identifier("*")}}
	}
|	'(' import_as_names optional_comma ')'
	{
		$$ = $2
	}
|	import_as_names optional_comma
	{
		$$ = $1
	}

import_from:
	FROM from_arg IMPORT import_from_arg
	{
		$$ = &ast.ImportFrom{StmtBase: ast.StmtBase{Pos: $<pos>$}, Module: ast.Identifier($2), Names: $4, Level: $<level>2}
	}

import_as_name:
	NAME
	{
		$$ = &ast.Alias{Pos: $<pos>$, Name: ast.Identifier($1)}
	}
|	NAME AS NAME
	{
		$$ = &ast.Alias{Pos: $<pos>$, Name: ast.Identifier($1), AsName: ast.Identifier($3)}
	}

dotted_as_name:
	dotted_name
	{
		$$ = &ast.Alias{Pos: $<pos>$, Name: ast.Identifier($1)}
	}
|	dotted_name AS NAME
	{
		$$ = &ast.Alias{Pos: $<pos>$, Name: ast.Identifier($1), AsName: ast.Identifier($3)}
	}

import_as_names:
	import_as_name
	{
		$$ = nil
		$$ = append($$, $1)
	}
|	import_as_names ',' import_as_name
	{
		$$ = append($$, $3)
	}

dotted_as_names:
	dotted_as_name
	{
		$$ = nil
		$$ = append($$, $1)
	}
|	dotted_as_names ',' dotted_as_name
	{
		$$ = append($$, $3)
	}

dotted_name:
	NAME
	{
		$$ = $1
	}
|	dotted_name '.' NAME
	{
		$$ += "." + $3
	}

names:
	NAME
	{
		$$ = nil
		$$ = append($$, ast.Identifier($1))
	}
|	names ',' NAME
	{
		$$ = append($$, ast.Identifier($3))
	}

global_stmt:
	GLOBAL names
	{
		$$ = &ast.Global{StmtBase: ast.StmtBase{Pos: $<pos>$}, Names: $2}
	}

nonlocal_stmt:
	NONLOCAL names
	{
		$$ = &ast.Nonlocal{StmtBase: ast.StmtBase{Pos: $<pos>$}, Names: $2}
	}

tests:
	test
	{
		$$ = nil
		$$ = append($$, $1)
	}
|	tests ',' test
	{
		$$ = append($$, $3)
	}

assert_stmt:
	ASSERT test
	{
		$$ = &ast.Assert{StmtBase: ast.StmtBase{Pos: $<pos>$}, Test: $2}
	}
|	ASSERT test ',' test
	{
		$$ = &ast.Assert{StmtBase: ast.StmtBase{Pos: $<pos>$}, Test: $2, Msg: $4}
	}

compound_stmt:
	if_stmt
	{
		$$ = $1
	}
|	while_stmt
	{
		$$ = $1
	}
|	for_stmt
	{
		$$ = $1
	}
|	try_stmt
	{
		$$ = $1
	}
|	with_stmt
	{
		$$ = $1
	}
|	funcdef
	{
		$$ = $1
	}
|	classdef
	{
		$$ = $1
	}
|	decorated
	{
		$$ = $1
	}

elifs:
	{
		$$ = nil
		$<lastif>$ = nil
	}
|	elifs ELIF test ':' suite
	{
		elifs := $$
		newif := &ast.If{StmtBase: ast.StmtBase{Pos: $<pos>$}, Test: $3, Body: $5}
		if elifs == nil {
			$$ = newif
		} else {
			$<lastif>$.Orelse = []ast.Stmt{newif}
		}
		$<lastif>$ = newif
	}

optional_else:
	{
		$$ = nil
	}
|	ELSE ':' suite
	{
		$$ = $3
	}

if_stmt:
	IF test ':' suite elifs optional_else
	{
		newif := &ast.If{StmtBase: ast.StmtBase{Pos: $<pos>$}, Test: $2, Body: $4}
		$$ = newif
		elifs := $5
		optional_else := $6
		if len(optional_else) != 0 {
			if elifs != nil {
				$<lastif>5.Orelse = optional_else
				newif.Orelse = []ast.Stmt{elifs}
			} else {
				newif.Orelse = optional_else
			}
		} else {
			if elifs != nil {
				newif.Orelse = []ast.Stmt{elifs}
			}
		}
	}

while_stmt:
	WHILE test ':' suite optional_else
	{
		$$ = &ast.While{StmtBase: ast.StmtBase{Pos: $<pos>$}, Test: $2, Body: $4, Orelse: $5}
	}

for_stmt:
	FOR exprlist IN testlist ':' suite optional_else
	{
		target := tupleOrExpr($<pos>$, $2, false)
		setCtx(yylex, target, ast.Store)
		$$ = &ast.For{StmtBase: ast.StmtBase{Pos: $<pos>$}, Target: target, Iter: $4, Body: $6, Orelse: $7}
	}

except_clauses:
	{
		$$ = nil
	}
|	except_clauses except_clause ':' suite
	{
		exc := &ast.ExceptHandler{Pos: $<pos>$, ExprType: $2, Name: ast.Identifier($<str>2), Body: $4}
		$$ = append($$, exc)
	}

try_stmt:
	TRY ':' suite except_clauses
	{
		$$ = &ast.Try{StmtBase: ast.StmtBase{Pos: $<pos>$}, Body: $3, Handlers: $4}
	}
|	TRY ':' suite except_clauses ELSE ':' suite
	{
		$$ = &ast.Try{StmtBase: ast.StmtBase{Pos: $<pos>$}, Body: $3, Handlers: $4, Orelse: $7}
	}
|	TRY ':' suite except_clauses FINALLY ':' suite
	{
		$$ = &ast.Try{StmtBase: ast.StmtBase{Pos: $<pos>$}, Body: $3, Handlers: $4, Finalbody: $7}
	}
|	TRY ':' suite except_clauses ELSE ':' suite FINALLY ':' suite
	{
		$$ = &ast.Try{StmtBase: ast.StmtBase{Pos: $<pos>$}, Body: $3, Handlers: $4, Orelse: $7, Finalbody: $10}
	}

with_items:
	with_item
	{
		$$ = nil
		$$ = append($$, $1)
	}
|	with_items ',' with_item
	{
		$$ = append($$, $3)
	}

with_stmt:
	WITH with_items ':' suite
	{
		$$ = &ast.With{StmtBase: ast.StmtBase{Pos: $<pos>$}, Items: $2, Body: $4}
	}

with_item:
	test
	{
		$$ = &ast.WithItem{Pos: $<pos>$, ContextExpr: $1}
	}
|	test AS expr
	{
		v := $3
		setCtx(yylex, v, ast.Store)
		$$ = &ast.WithItem{Pos: $<pos>$, ContextExpr: $1, OptionalVars: v}
	}

// NB compile.c makes sure that the default except clause is last
except_clause:
	EXCEPT
	{
		$$ = nil
		$<str>$ = ""
	}
|	EXCEPT test
	{
		$$ = $2
		$<str>$ = ""
	}
|	EXCEPT test AS NAME
	{
		$$ = $2
		$<str>$ = $4
	}

stmts:
	stmt
	{
		$$ = nil
		$$ = append($$, $1...)
	}
|	stmts stmt
	{
		$$ = append($$, $2...)
	}

suite:
	simple_stmt
	{
		$$ = $1
	}
|	NEWLINE INDENT stmts DEDENT
	{
		$$ = $3
	}

test:
	or_test
	{
		$$ = $1
	}
|	or_test IF or_test ELSE test
	{
		$$ = &ast.IfExp{ExprBase: ast.ExprBase{Pos: $<pos>$}, Test:$3, Body: $1, Orelse: $5}
	}
|	lambdef
	{
		$$ = $1
	}

test_nocond:
	or_test
	{
		$$ = $1
	}
|	lambdef_nocond
	{
		$$ = $1
	}

lambdef:
	LAMBDA ':' test
	{
		args := &ast.Arguments{Pos: $<pos>$}
		$$ = &ast.Lambda{ExprBase: ast.ExprBase{Pos: $<pos>$}, Args: args, Body: $3}
	}
|	LAMBDA varargslist ':' test
	{
		$$ = &ast.Lambda{ExprBase: ast.ExprBase{Pos: $<pos>$}, Args: $2, Body: $4}
	}

lambdef_nocond:
	LAMBDA ':' test_nocond
	{
		args := &ast.Arguments{Pos: $<pos>$}
		$$ = &ast.Lambda{ExprBase: ast.ExprBase{Pos: $<pos>$}, Args: args, Body: $3}
	}
|	LAMBDA varargslist ':' test_nocond
	{
		$$ = &ast.Lambda{ExprBase: ast.ExprBase{Pos: $<pos>$}, Args: $2, Body: $4}
	}

or_test:
	and_test
	{
		$$ = $1
		$<isExpr>$ = true
	}
|	or_test OR and_test
	{
		if !$<isExpr>1 {
			boolop := $$.(*ast.BoolOp)
			boolop.Values = append(boolop.Values, $3)
		} else {
			$$ = &ast.BoolOp{ExprBase: ast.ExprBase{Pos: $<pos>$}, Op: ast.Or, Values: []ast.Expr{$$, $3}}
		}
		$<isExpr>$ = false
	}

and_test:
	not_test
	{
		$$ = $1
		$<isExpr>$ = true
	}
|	and_test AND not_test
	{
		if !$<isExpr>1 {
			boolop := $$.(*ast.BoolOp)
			boolop.Values = append(boolop.Values, $3)
		} else {
			$$ = &ast.BoolOp{ExprBase: ast.ExprBase{Pos: $<pos>$}, Op: ast.And, Values: []ast.Expr{$$, $3}}
		}
		$<isExpr>$ = false
	}

not_test:
	NOT not_test
	{
		$$ = &ast.UnaryOp{ExprBase: ast.ExprBase{Pos: $<pos>$}, Op: ast.Not, Operand: $2}
	}
|	comparison
	{
		$$ = $1
	}

comparison:
	expr
	{
		$$ = $1
		$<isExpr>$ = true
	}
|	comparison comp_op expr
	{
		if !$<isExpr>1 {
			comp := $$.(*ast.Compare)
			comp.Ops = append(comp.Ops, $2)
			comp.Comparators = append(comp.Comparators, $3)
		} else{
			$$ = &ast.Compare{ExprBase: ast.ExprBase{Pos: $<pos>$}, Left: $$, Ops: []ast.CmpOp{$2}, Comparators: []ast.Expr{$3}}
		}
		$<isExpr>$ = false
	}

// <> LTGT isn't actually a valid comparison operator in Python. It's here for the
// sake of a __future__ import described in PEP 401
comp_op:
	'<'
	{
		$$ = ast.Lt
	}
|	'>'
	{
		$$ = ast.Gt
	}
|	EQEQ
	{
		$$ = ast.Eq
	}
|	GTEQ
	{
		$$ = ast.GtE
	}
|	LTEQ
	{
		$$ = ast.LtE
	}
|	LTGT
	{
		yylex.(*yyLex).SyntaxError("invalid syntax")
	}
|	PLINGEQ
	{
		$$ = ast.NotEq
	}
|	IN
	{
		$$ = ast.In
	}
|	NOT IN
	{
		$$ = ast.NotIn
	}
|	IS
	{
		$$ = ast.Is
	}
|	IS NOT
	{
		$$ = ast.IsNot
	}

star_expr:
	'*' expr
	{
		$$ = &ast.Starred{ExprBase: ast.ExprBase{Pos: $<pos>$}, Value: $2, Ctx: ast.Load}
	}

expr:
	xor_expr
	{
		$$ = $1
	}
|	expr '|' xor_expr
	{
		$$ = &ast.BinOp{ExprBase: ast.ExprBase{Pos: $<pos>$}, Left: $1, Op: ast.BitOr, Right: $3}
	}

xor_expr:
	and_expr
	{
		$$ = $1
	}
|	xor_expr '^' and_expr
	{
		$$ = &ast.BinOp{ExprBase: ast.ExprBase{Pos: $<pos>$}, Left: $1, Op: ast.BitXor, Right: $3}
	}

and_expr:
	shift_expr
	{
		$$ = $1
	}
|	and_expr '&' shift_expr
	{
		$$ = &ast.BinOp{ExprBase: ast.ExprBase{Pos: $<pos>$}, Left: $1, Op: ast.BitAnd, Right: $3}
	}

shift_expr:
	arith_expr
	{
		$$ = $1
	}
|	shift_expr LTLT arith_expr
	{
		$$ = &ast.BinOp{ExprBase: ast.ExprBase{Pos: $<pos>$}, Left: $1, Op: ast.LShift, Right: $3}
	}
|	shift_expr GTGT arith_expr
	{
		$$ = &ast.BinOp{ExprBase: ast.ExprBase{Pos: $<pos>$}, Left: $1, Op: ast.RShift, Right: $3}
	}

arith_expr:
	term
	{
		$$ = $1
	}
|	arith_expr '+' term
	{
		$$ = &ast.BinOp{ExprBase: ast.ExprBase{Pos: $<pos>$}, Left: $1, Op: ast.Add, Right: $3}
	}
|	arith_expr '-' term
	{
		$$ = &ast.BinOp{ExprBase: ast.ExprBase{Pos: $<pos>$}, Left: $1, Op: ast.Sub, Right: $3}
	}

term:
	factor
	{
		$$ = $1
	}
|	term '*' factor
	{
		$$ = &ast.BinOp{ExprBase: ast.ExprBase{Pos: $<pos>$}, Left: $1, Op: ast.Mult, Right: $3}
	}
|	term '/' factor
	{
		$$ = &ast.BinOp{ExprBase: ast.ExprBase{Pos: $<pos>$}, Left: $1, Op: ast.Div, Right: $3}
	}
|	term '%' factor
	{
		$$ = &ast.BinOp{ExprBase: ast.ExprBase{Pos: $<pos>$}, Left: $1, Op: ast.Modulo, Right: $3}
	}
|	term DIVDIV factor
	{
		$$ = &ast.BinOp{ExprBase: ast.ExprBase{Pos: $<pos>$}, Left: $1, Op: ast.FloorDiv, Right: $3}
	}

factor:
	'+' factor
	{
		$$ = &ast.UnaryOp{ExprBase: ast.ExprBase{Pos: $<pos>$}, Op: ast.UAdd, Operand: $2}
	}
|	'-' factor
	{
		$$ = &ast.UnaryOp{ExprBase: ast.ExprBase{Pos: $<pos>$}, Op: ast.USub, Operand: $2}
	}
|	'~' factor
	{
		$$ = &ast.UnaryOp{ExprBase: ast.ExprBase{Pos: $<pos>$}, Op: ast.Invert, Operand: $2}
	}
|	power
	{
		$$ = $1
	}

power:
	atom trailers
	{
		$$ = applyTrailers($1, $2)
	}
|	atom trailers STARSTAR factor
	{
		$$ = &ast.BinOp{ExprBase: ast.ExprBase{Pos: $<pos>$}, Left: applyTrailers($1, $2), Op: ast.Pow, Right: $4}
	}

// Trailers are half made Call, Attribute or Subscript
trailers:
	{
		$$ = nil
	}
|	trailers trailer
	{
		$$ = append($$, $2)
	}

strings:
	STRING
	{
		$$ = $1
	}
|	strings STRING
	{
		switch a := $$.(type) {
		case py.String:
			switch b := $2.(type) {
			case py.String:
				$$ = a + b
			default:
				yylex.(*yyLex).SyntaxError("cannot mix string and nonstring literals")
			}
		case py.Bytes:
			switch b := $2.(type) {
			case py.Bytes:
				$$ = append(a, b...)
			default:
				yylex.(*yyLex).SyntaxError("cannot mix bytes and nonbytes literals")
			}
		}
	}

atom:
	'(' ')'
	{
		$$ = &ast.Tuple{ExprBase: ast.ExprBase{Pos: $<pos>$}, Ctx: ast.Load}
	}
|	'(' yield_expr ')'
	{
		$$ = $2
	}
|	'(' test_or_star_expr comp_for ')'
	{
		$$ = &ast.GeneratorExp{ExprBase: ast.ExprBase{Pos: $<pos>$}, Elt: $2, Generators: $3}
	}
|	'(' test_or_star_exprs optional_comma ')' 
	{
		$$ = tupleOrExpr($<pos>$, $2, $3)
	}
|	'[' ']'
	{
		$$ = &ast.List{ExprBase: ast.ExprBase{Pos: $<pos>$}, Ctx: ast.Load}
	}
|	'[' test_or_star_expr comp_for ']'
	{
		$$ = &ast.ListComp{ExprBase: ast.ExprBase{Pos: $<pos>$}, Elt: $2, Generators: $3}
	}
|	'[' test_or_star_exprs optional_comma ']'
	{
		$$ = &ast.List{ExprBase: ast.ExprBase{Pos: $<pos>$}, Elts: $2, Ctx: ast.Load}
	}
|	'{' '}'
	{
		$$ = &ast.Dict{ExprBase: ast.ExprBase{Pos: $<pos>$}}
	}
|	'{' dictorsetmaker '}'
	{
		$$ = $2
	}
|	NAME
	{
		$$ = &ast.Name{ExprBase: ast.ExprBase{Pos: $<pos>$}, Id: ast.Identifier($1), Ctx: ast.Load}
	}
|	NUMBER
	{
		$$ = &ast.Num{ExprBase: ast.ExprBase{Pos: $<pos>$}, N: $1}
	}
|	strings
	{
		switch s := $1.(type) {
		case py.String:
			$$ = &ast.Str{ExprBase: ast.ExprBase{Pos: $<pos>$}, S: s}
		case py.Bytes:
			$$ = &ast.Bytes{ExprBase: ast.ExprBase{Pos: $<pos>$}, S: s}
		default:
			panic("not Bytes or String in strings")
		}
	}
|	ELIPSIS
	{
		$$ = &ast.Ellipsis{ExprBase: ast.ExprBase{Pos: $<pos>$}}
	}
|	NONE
	{
		$$ = &ast.NameConstant{ExprBase: ast.ExprBase{Pos: $<pos>$}, Value: py.None}
	}
|	TRUE
	{
		$$ = &ast.NameConstant{ExprBase: ast.ExprBase{Pos: $<pos>$}, Value: py.True}
	}
|	FALSE
	{
		$$ = &ast.NameConstant{ExprBase: ast.ExprBase{Pos: $<pos>$}, Value: py.False}
	}

// Trailers are half made Call, Attribute or Subscript
trailer:
	'(' ')'
	{
		$$ = &ast.Call{ExprBase: ast.ExprBase{Pos: $<pos>$}}
	}
|	'(' arglist ')'
	{
		$$ = $2
	}
|	'[' subscriptlist ']'
	{
		slice := $2
		// If all items of a ExtSlice are just Index then return as tuple
		if extslice, ok := slice.(*ast.ExtSlice); ok {
			elts := make([]ast.Expr, len(extslice.Dims))
			for i, item := range(extslice.Dims) {
				if index, isIndex := item.(*ast.Index); isIndex {
					elts[i] = index.Value
				} else {
					goto notAllIndex
				}
			}
                slice = &ast.Index{SliceBase: extslice.SliceBase, Value: &ast.Tuple{ExprBase: ast.ExprBase{Pos: extslice.SliceBase.Pos}, Elts: elts, Ctx: ast.Load}}
		notAllIndex:
		}
		$$ = &ast.Subscript{ExprBase: ast.ExprBase{Pos: $<pos>$}, Slice: slice, Ctx: ast.Load}
	}
|	'.' NAME
	{
		$$ = &ast.Attribute{ExprBase: ast.ExprBase{Pos: $<pos>$}, Attr: ast.Identifier($2), Ctx: ast.Load}
	}

subscripts:
	subscript
	{
		$$ = $1
		$<isExpr>$ = true
	}
|	subscripts ',' subscript
	{
		if !$<isExpr>1 {
			extSlice := $$.(*ast.ExtSlice)
			extSlice.Dims = append(extSlice.Dims, $3)
		} else {
			$$ = &ast.ExtSlice{SliceBase: ast.SliceBase{Pos: $<pos>$}, Dims: []ast.Slicer{$1, $3}}
		}
		$<isExpr>$ = false
	}

subscriptlist:
	subscripts optional_comma
	{
		if $2 && $<isExpr>1 {
			$$ = &ast.ExtSlice{SliceBase: ast.SliceBase{Pos: $<pos>$}, Dims: []ast.Slicer{$1}}
		} else {
			$$ = $1
		}
	}

subscript:
	test
	{
		$$ = &ast.Index{SliceBase: ast.SliceBase{Pos: $<pos>$}, Value: $1}
	}
|	':'
	{
		$$ = &ast.Slice{SliceBase: ast.SliceBase{Pos: $<pos>$}, Lower: nil, Upper: nil, Step: nil}
	}
|	':' sliceop
	{
		$$ = &ast.Slice{SliceBase: ast.SliceBase{Pos: $<pos>$}, Lower: nil, Upper: nil, Step: $2}
	}
|	':' test
	{
		$$ = &ast.Slice{SliceBase: ast.SliceBase{Pos: $<pos>$}, Lower: nil, Upper: $2, Step: nil}
	}
|	':' test sliceop
	{
		$$ = &ast.Slice{SliceBase: ast.SliceBase{Pos: $<pos>$}, Lower: nil, Upper: $2, Step: $3}
	}
|	test ':'
	{
		$$ = &ast.Slice{SliceBase: ast.SliceBase{Pos: $<pos>$}, Lower: $1, Upper: nil, Step: nil}
	}
|	test ':' sliceop
	{
		$$ = &ast.Slice{SliceBase: ast.SliceBase{Pos: $<pos>$}, Lower: $1, Upper: nil, Step: $3}
	}
|	test ':' test
	{
		$$ = &ast.Slice{SliceBase: ast.SliceBase{Pos: $<pos>$}, Lower: $1, Upper: $3, Step: nil}
	}
|	test ':' test sliceop
	{
		$$ = &ast.Slice{SliceBase: ast.SliceBase{Pos: $<pos>$}, Lower: $1, Upper: $3, Step: $4}
	}

sliceop:
	':'
	{
		$$ = nil
	}
|	':' test
	{
		$$ = $2
	}

expr_or_star_expr:
	expr
	{
		$$ = $1
	}
|	star_expr
	{
		$$ = $1
	}

expr_or_star_exprs:
	expr_or_star_expr
	{
		$$ = nil
		$$ = append($$, $1)
	}
|	expr_or_star_exprs ',' expr_or_star_expr
	{
		$$ = append($$, $3)
	}

exprlist:
	expr_or_star_exprs optional_comma
	{
		$$ = $1
		$<comma>$ = $2
	}

testlist:
	tests optional_comma
	{
		elts := $1
		if $2 || len(elts) > 1 {
			$$ = &ast.Tuple{ExprBase: ast.ExprBase{Pos: $<pos>$}, Elts: elts, Ctx: ast.Load}
		} else {
			$$ = elts[0]
		}
	}

testlistraw:
	tests optional_comma
	{
		$$ = $1
	}

// (',' test ':' test)*
test_colon_tests:
	test ':' test
	{
		$$ = nil
		$$ = append($$, $1, $3)	// key, value order
	}
|	test_colon_tests ',' test ':' test
	{
		$$ = append($$, $3, $5)
	}

dictorsetmaker:
	test_colon_tests optional_comma
	{
		keyValues := $1
		d := &ast.Dict{ExprBase: ast.ExprBase{Pos: $<pos>$}, Keys: nil, Values: nil}
		for i := 0; i < len(keyValues)-1; i += 2 {
			d.Keys = append(d.Keys, keyValues[i])
			d.Values = append(d.Values, keyValues[i+1])
		}
		$$ = d
	}
|	test ':' test comp_for
	{
		$$ = &ast.DictComp{ExprBase: ast.ExprBase{Pos: $<pos>$}, Key: $1, Value: $3, Generators: $4}
	}
|	testlistraw
	{
		$$ = &ast.Set{ExprBase: ast.ExprBase{Pos: $<pos>$}, Elts: $1}
	}
|	test comp_for
	{
		$$ = &ast.SetComp{ExprBase: ast.ExprBase{Pos: $<pos>$}, Elt: $1, Generators: $2}
	}

classdef:
	CLASS NAME optional_arglist_call ':' suite
	{
		classDef := &ast.ClassDef{StmtBase: ast.StmtBase{Pos: $<pos>$}, Name: ast.Identifier($2), Body: $5}
		$$ = classDef
		args := $3
		if args != nil {
			classDef.Bases = args.Args
			classDef.Keywords = args.Keywords
			classDef.Starargs = args.Starargs
			classDef.Kwargs = args.Kwargs
		}
	}

arguments:
	argument
	{
		$$ = $1
	}
|	arguments ',' argument
	{
		$$.Args = append($$.Args, $3.Args...)
		$$.Keywords = append($$.Keywords, $3.Keywords...)
	}

optional_arguments:
	{
		$$ = &ast.Call{}
	}
|	arguments ','
	{
		$$ = $1
	}

arguments2:
	{
		$$ = &ast.Call{}
	}
|	arguments2 ',' argument
	{
		$$.Args = append($$.Args, $3.Args...)
		$$.Keywords = append($$.Keywords, $3.Keywords...)
	}

arglist:
	arguments optional_comma
	{
		$$ = $1
	}
|	optional_arguments '*' test arguments2
	{
		call := $1
		call.Starargs = $3
		if len($4.Args) != 0 {
			yylex.(*yyLex).SyntaxError("only named arguments may follow *expression")
		}
		call.Keywords = append(call.Keywords, $4.Keywords...)
		$$ = call
	}
|	optional_arguments '*' test arguments2 ',' STARSTAR test
	{
		call := $1
		call.Starargs = $3
		call.Kwargs = $7
		if len($4.Args) != 0 {
			yylex.(*yyLex).SyntaxError("only named arguments may follow *expression")
		}
		call.Keywords = append(call.Keywords, $4.Keywords...)
		$$ = call
	}
|	optional_arguments STARSTAR test
	{
		call := $1
		call.Kwargs = $3
		$$ = call
	}

// The reason that keywords are test nodes instead of NAME is that using NAME
// results in an ambiguity. ast.c makes sure it's a NAME.
argument:
	test
	{
		$$ = &ast.Call{}
		$$.Args = []ast.Expr{$1}
	}
|	test comp_for
	{
		$$ = &ast.Call{}
		$$.Args = []ast.Expr{
			&ast.GeneratorExp{ExprBase: ast.ExprBase{Pos: $<pos>$}, Elt: $1, Generators: $2},
		}
	}
|	test '=' test  // Really [keyword '='] test
	{
		$$ = &ast.Call{}
		test := $1
		if name, ok := test.(*ast.Name); ok {
			$$.Keywords = []*ast.Keyword{&ast.Keyword{Pos: name.Pos, Arg: name.Id, Value: $3}}
		} else {
			yylex.(*yyLex).SyntaxError("keyword can't be an expression")
		}
	}

comp_iter:
	comp_for
	{
		$<comprehensions>$ = $1
		$$ = nil
	}
|	comp_if
	{
		$<comprehensions>$ = $<comprehensions>1
		$$ = $1
	}

comp_for:
	FOR exprlist IN or_test
	{
		c := ast.Comprehension{
			Target: tupleOrExpr($<pos>$, $2, $<comma>2),
			Iter: $4,
		}
		setCtx(yylex, c.Target, ast.Store)
		$$ = []ast.Comprehension{c}
	}
|	FOR exprlist IN or_test comp_iter
	{
		c := ast.Comprehension{
			Target: tupleOrExpr($<pos>$, $2, $<comma>2),
			Iter: $4,
			Ifs: $5,
		}
		setCtx(yylex, c.Target, ast.Store)
		$$ = []ast.Comprehension{c}
		$$ = append($$, $<comprehensions>5...)
	}

comp_if:
	IF test_nocond
	{
		$$ = []ast.Expr{$2}
		$<comprehensions>$ = nil
	}
|	IF test_nocond comp_iter
	{
		$$ = []ast.Expr{$2}
		$$ = append($$, $3...)
		$<comprehensions>$ = $<comprehensions>3
	}

// not used in grammar, but may appear in "node" passed from Parser to Compiler
// encoding_decl: NAME

yield_expr:
	YIELD
	{
		$$ = &ast.Yield{ExprBase: ast.ExprBase{Pos: $<pos>$}}
	}
|	YIELD FROM test
	{
		$$= &ast.YieldFrom{ExprBase: ast.ExprBase{Pos: $<pos>$}, Value: $3}
	}
|	YIELD testlist
	{
		$$= &ast.Yield{ExprBase: ast.ExprBase{Pos: $<pos>$}, Value: $2}
	}
