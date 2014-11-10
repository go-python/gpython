%{

package parser

// Grammar for Python

import (
	"github.com/ncw/gpython/py"
	"github.com/ncw/gpython/ast"
)

%}

%union {
	str	string
	obj	py.Object
	ast	ast.Ast
	mod	ast.Mod
	stmt	ast.Stmt
	stmts	[]ast.Stmt
	stmts1	[]ast.Stmt	// nl_or_stmt accumulator
	stmts2	[]ast.Stmt	// small_stmts accumulator
	stmts3	[]ast.Stmt	// stmts accumulator
	pos	ast.Pos		// kept up to date by the lexer
	op	ast.OperatorNumber
	cmpop	ast.CmpOp
	expr	ast.Expr
	exprs	[]ast.Expr
	exprs1	[]ast.Expr	// expr_or_star_exprs accumulator
	exprs2	[]ast.Expr	// test_or_star_exprs accumulator
	trailers []ast.Expr	// list of trailer expressions
	cmp	*ast.Compare
	boolop1 *ast.BoolOp	// or_test
	boolop2 *ast.BoolOp	// and_test
}

%type <str> strings
%type <mod> inputs file_input single_input eval_input
%type <stmts> simple_stmt stmt 
%type <stmts1> nl_or_stmt 
%type <stmts2> small_stmts
%type <stmts3> stmts
%type <stmt> compound_stmt small_stmt expr_stmt del_stmt pass_stmt flow_stmt import_stmt global_stmt nonlocal_stmt assert_stmt
%type <op> augassign
%type <expr> expr_or_star_expr expr star_expr xor_expr and_expr shift_expr arith_expr term factor power trailer atom test_or_star_expr test not_test lambdef test_nocond lambdef_nocond
%type <exprs> exprlist
%type <exprs1> expr_or_star_exprs
%type <exprs2> test_or_star_exprs
%type <trailers> trailers
%type <cmpop> comp_op
%type <cmp> comparison
%type <boolop1> or_test
%type <boolop2> and_test 

%token NEWLINE
%token ENDMARKER
%token <str> NAME
%token INDENT
%token DEDENT
%token <str> STRING
%token <str> NUMBER

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
	NEWLINE
	{
		$$ = &ast.Interactive{ModBase: ast.ModBase{$<pos>$}}
	}
|	simple_stmt
	{
		$$ = &ast.Interactive{ModBase: ast.ModBase{$<pos>$}, Body: $1}
	}
|	compound_stmt NEWLINE
	{
		$$ = &ast.Interactive{ModBase: ast.ModBase{$<pos>$}, Body: []ast.Stmt{$1}}
	}

//file_input: (NEWLINE | stmt)* ENDMARKER
file_input:
	nl_or_stmt ENDMARKER
	{
		$$ = &ast.Module{ModBase: ast.ModBase{$<pos>$}, Body: $1}
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
	}

// NEWLINE*
nls: | nls NEWLINE

optional_arglist: | arglist

optional_arglist_call: | '(' optional_arglist ')'

decorator: '@' dotted_name optional_arglist_call NEWLINE

decorators: decorator | decorators decorator

classdef_or_funcdef: classdef | funcdef

decorated: decorators classdef_or_funcdef

optional_return_type: | MINUSGT test

funcdef: DEF NAME parameters optional_return_type ':' suite

parameters: '(' optional_typedargslist ')'

optional_typedargslist: | typedargslist

// (',' tfpdef ['=' test])*
tfpdeftest: tfpdef | tfpdef '=' test

tfpdeftests: | tfpdeftests ',' tfpdeftest

optional_tfpdef: | tfpdef

typedargslist: 
	tfpdeftest tfpdeftests
|	tfpdeftest tfpdeftests ','
|	tfpdeftest tfpdeftests ',' '*' optional_tfpdef tfpdeftests
|	tfpdeftest tfpdeftests ',' '*' optional_tfpdef tfpdeftests ',' STARSTAR tfpdef
|	tfpdeftest tfpdeftests ',' STARSTAR tfpdef
|	'*' optional_tfpdef tfpdeftests
|	'*' optional_tfpdef tfpdeftests ',' STARSTAR tfpdef
|	STARSTAR tfpdef

tfpdef:
	NAME
|	NAME ':' test

vfpdeftest:
	vfpdef
|	vfpdef '=' test

vfpdeftests:
|	vfpdeftests ',' vfpdeftest

optional_vfpdef:
|	vfpdef

varargslist:
	vfpdeftest vfpdeftests
|	vfpdeftest vfpdeftests ','
|	vfpdeftest vfpdeftests ',' '*' optional_vfpdef vfpdeftests
|	vfpdeftest vfpdeftests ',' '*' optional_vfpdef vfpdeftests ',' STARSTAR vfpdef
|	vfpdeftest vfpdeftests ',' STARSTAR vfpdef
|	'*' optional_vfpdef vfpdeftests
|	'*' optional_vfpdef vfpdeftests ',' STARSTAR vfpdef
|	STARSTAR vfpdef

vfpdef: NAME

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
		$$ = []ast.Stmt{$1}
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

expr_stmt:
	testlist_star_expr augassign yield_expr_or_testlist
	{
	}
|	testlist_star_expr equals_yield_expr_or_testlist_star_expr
	{
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
	{
	}
|	equals_yield_expr_or_testlist_star_expr '=' yield_expr_or_testlist_star_expr
	{
	}

test_or_star_exprs:
	test_or_star_expr
	{
		$$ = []ast.Expr{$1}
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

optional_comma: | ','

testlist_star_expr:
	test_or_star_exprs optional_comma
	{
		$$ = $1
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
		$$ = &ast.Delete{StmtBase: ast.StmtBase{$<pos>$}, Targets: $2}
	}

pass_stmt:
	PASS
	{
		$$ = &ast.Pass{StmtBase: ast.StmtBase{$<pos>$}}
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
	}

continue_stmt:
	CONTINUE
	{
	}

return_stmt:
	RETURN
	{
	}
|	RETURN testlist
	{
	}

yield_stmt:
	yield_expr
	{
	}

raise_stmt:
	RAISE
	{
	}
|	RAISE test
	{
	}
|	RAISE test FROM test
	{
	}

import_stmt:
	import_name
	{
	}
|	import_from
	{
	}

import_name:
	IMPORT dotted_as_names
	{
	}

// note below: the '.' | ELIPSIS is necessary because '...' is tokenized as ELIPSIS
dot: '.' | ELIPSIS

dots: dot | dots dot

from_arg: dotted_name | dots dotted_name | dots

import_from_arg: '*' | '(' import_as_names ')' | import_as_names

import_from: FROM from_arg IMPORT import_from_arg

import_as_name: NAME | NAME AS NAME

dotted_as_name: dotted_name | dotted_name AS NAME

import_as_names: import_as_name optional_comma | import_as_name ',' import_as_names

dotted_as_names: dotted_as_name | dotted_as_names ',' dotted_as_name

dotted_name: NAME | dotted_name '.' NAME

names: NAME | names ',' NAME

global_stmt:
	GLOBAL names
	{
	}

nonlocal_stmt:
	NONLOCAL names
	{
	}

tests: test | tests ',' test

assert_stmt:
	ASSERT tests
	{
	}

compound_stmt:
	if_stmt
	{
	}
|	while_stmt
	{
	}
|	for_stmt
	{
	}
|	try_stmt
	{
	}
|	with_stmt
	{
	}
|	funcdef
	{
	}
|	classdef
	{
	}
|	decorated
	{
	}

elifs:
|	elifs ELIF test ':' suite

optional_else:
|	ELSE ':' suite

if_stmt: IF test ':' suite elifs optional_else

while_stmt: WHILE test ':' suite optional_else

for_stmt: FOR exprlist IN testlist ':' suite optional_else

except_clauses:
|	except_clauses except_clause ':' suite

try_stmt: TRY ':' suite except_clauses
       
|	TRY ':' suite except_clauses ELSE ':' suite
       
|	TRY ':' suite except_clauses FINALLY ':' suite
       
|	TRY ':' suite except_clauses ELSE ':' suite FINALLY ':' suite

with_items: with_item
|	with_items ',' with_item

with_stmt: WITH with_items  ':' suite

with_item: test
|	test AS expr

// NB compile.c makes sure that the default except clause is last
except_clause: EXCEPT
|	EXCEPT test
|	EXCEPT test AS NAME

stmts:
	stmt
	{
		$$ = make([]ast.Stmt, len($1))
		copy($$, $1)
	}
|	stmts stmt
	{
		$$ = append($$, $2...)
	}

suite:
	simple_stmt
|	NEWLINE INDENT stmts DEDENT

test:
	or_test
	{
		$$ = $1
	}
|	or_test IF or_test ELSE test
	{
		$$ = &ast.IfExp{ExprBase: ast.ExprBase{$<pos>$}, Test:$1, Body: $3, Orelse: $5} // FIXME Ctx
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
	}
|	LAMBDA varargslist ':' test
	{
	}

lambdef_nocond:
	LAMBDA ':' test_nocond
	{
	}
|	LAMBDA varargslist ':' test_nocond
	{
	}

or_test:
	and_test
	{
		$$ = &ast.BoolOp{ExprBase: ast.ExprBase{$<pos>$}, Op: ast.Or} // FIXME Ctx
	}
|	or_test OR and_test
	{
		$$.Values = append($$.Values, $3)
	}

and_test:
	not_test
	{
		$$ = &ast.BoolOp{ExprBase: ast.ExprBase{$<pos>$}, Op: ast.And}
	}
|	and_test AND not_test
	{
		$$.Values = append($$.Values, $3)
	}

not_test:
	NOT not_test
	{
		$$ = &ast.UnaryOp{ExprBase: ast.ExprBase{$<pos>$}, Op: ast.Not, Operand: $2} // FIXME Ctx
	}
|	comparison
	{
		$$ = $1
	}

comparison:
	expr
	{
		$$ = &ast.Compare{ExprBase: ast.ExprBase{$<pos>$}, Left: $1} // FIXME Ctx
	}
|	comparison comp_op expr
	{
		$$.Ops = append($$.Ops, $2)
		$$.Comparators = append($$.Comparators, $3)
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
		yylex.Error("Invalid syntax")
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
		$$ = &ast.Starred{ExprBase: ast.ExprBase{$<pos>$}, Value: $2} // FIXME Ctx
	}

expr:
	xor_expr
	{
		$$ = $1
	}
|	expr '|' xor_expr
	{
		$$ = &ast.BinOp{ExprBase: ast.ExprBase{$<pos>$}, Left: $1, Op: ast.BitOr, Right: $3} // FIXME Ctx
	}

xor_expr:
	and_expr
	{
		$$ = $1
	}
|	xor_expr '^' and_expr
	{
		$$ = &ast.BinOp{ExprBase: ast.ExprBase{$<pos>$}, Left: $1, Op: ast.BitXor, Right: $3} // FIXME Ctx
	}

and_expr:
	shift_expr
	{
		$$ = $1
	}
|	and_expr '&' shift_expr
	{
		$$ = &ast.BinOp{ExprBase: ast.ExprBase{$<pos>$}, Left: $1, Op: ast.BitAnd, Right: $3} // FIXME Ctx
	}

shift_expr:
	arith_expr
	{
		$$ = $1
	}
|	shift_expr LTLT arith_expr
	{
		$$ = &ast.BinOp{ExprBase: ast.ExprBase{$<pos>$}, Left: $1, Op: ast.LShift, Right: $3} // FIXME Ctx
	}
|	shift_expr GTGT arith_expr
	{
		$$ = &ast.BinOp{ExprBase: ast.ExprBase{$<pos>$}, Left: $1, Op: ast.RShift, Right: $3} // FIXME Ctx
	}

arith_expr:
	term
	{
		$$ = $1
	}
|	arith_expr '+' term
	{
		$$ = &ast.BinOp{ExprBase: ast.ExprBase{$<pos>$}, Left: $1, Op: ast.Add, Right: $3} // FIXME Ctx
	}
|	arith_expr '-' term
	{
		$$ = &ast.BinOp{ExprBase: ast.ExprBase{$<pos>$}, Left: $1, Op: ast.Sub, Right: $3} // FIXME Ctx
	}

term:
	factor
	{
		$$ = $1
	}
|	term '*' factor
	{
		$$ = &ast.BinOp{ExprBase: ast.ExprBase{$<pos>$}, Left: $1, Op: ast.Mult, Right: $3} // FIXME Ctx
	}
|	term '/' factor
	{
		$$ = &ast.BinOp{ExprBase: ast.ExprBase{$<pos>$}, Left: $1, Op: ast.Div, Right: $3} // FIXME Ctx
	}
|	term '%' factor
	{
		$$ = &ast.BinOp{ExprBase: ast.ExprBase{$<pos>$}, Left: $1, Op: ast.Modulo, Right: $3} // FIXME Ctx
	}
|	term DIVDIV factor
	{
		$$ = &ast.BinOp{ExprBase: ast.ExprBase{$<pos>$}, Left: $1, Op: ast.FloorDiv, Right: $3} // FIXME Ctx
	}

factor:
	'+' factor
	{
		$$ = &ast.UnaryOp{ExprBase: ast.ExprBase{$<pos>$}, Op: ast.UAdd, Operand: $2} // FIXME Ctx
	}
|	'-' factor
	{
		$$ = &ast.UnaryOp{ExprBase: ast.ExprBase{$<pos>$}, Op: ast.USub, Operand: $2} // FIXME Ctx
	}
|	'~' factor
	{
		$$ = &ast.UnaryOp{ExprBase: ast.ExprBase{$<pos>$}, Op: ast.Invert, Operand: $2} // FIXME Ctx
	}
|	power
	{
		$$ = $1
	}

power:
	atom trailers
	{
		// FIXME apply trailers (if any) to atom
		$$ = $1
	}
|	atom trailers STARSTAR factor
	{
		// FIXME apply trailers (if any) to atom
		$$ = $1
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
|	strings STRING
	{
		$$ += $2
	}

atom:
	'(' ')'
	{
		$$ = &ast.Tuple{ExprBase: ast.ExprBase{$<pos>$}}
	}
|	'(' yield_expr ')'
	{
		// FIXME
		$$ = nil
	}
|	'(' testlist_comp ')'
	{
		// FIXME
		$$ = nil
	}
|	'[' ']'
	{
		$$ = &ast.List{ExprBase: ast.ExprBase{$<pos>$}}
	}
|	'[' testlist_comp ']'
	{
		// FIXME
		$$ = nil
	}
|	'{' '}'
	{
		$$ = &ast.Dict{ExprBase: ast.ExprBase{$<pos>$}}
	}
|	'{' dictorsetmaker '}'
	{
		// FIXME
		$$ = nil
	}
|	NAME
	{
		$$ = &ast.Name{ExprBase: ast.ExprBase{$<pos>$}, Id: ast.Identifier($1)}
	}
|	NUMBER
	{
		// FIXME
		$$ = nil
	}
|	strings
	{
		// FIXME
		$$ = nil
	}
|	ELIPSIS
	{
		// FIXME
		$$ = nil
	}
|	NONE
	{
		// FIXME
		$$ = nil
	}
|	TRUE
	{
		// FIXME
		$$ = nil
	}
|	FALSE
	{
		// FIXME
		$$ = nil
	}

testlist_comp:
	test_or_star_expr comp_for
|	test_or_star_exprs optional_comma

// Trailers are half made Call, Attribute or Subscript
trailer:
	'(' ')'
	{
	}
|	'(' arglist ')'
	{
	}
|	'[' subscriptlist ']'
	{
	}
|	'.' NAME
	{
		$$ = &ast.Attribute{ExprBase: ast.ExprBase{$<pos>$}, Attr: ast.Identifier($2)} // FIXME Ctx
	}

subscripts:
	subscript
	{
	}
|	subscripts ',' subscript
	{
	}

subscriptlist:
	subscripts optional_comma
	{
	}

subscript:
	test
|	':'
|	':' sliceop
|	':' test
|	':' test sliceop
|	test ':'
|	test ':' sliceop
|	test ':' test
|	test ':' test sliceop

sliceop:
	':'
|	':' test

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
		$$ = []ast.Expr{$1}
	}
|	expr_or_star_exprs ',' expr_or_star_expr
	{
		$$ = append($$, $3)
	}

exprlist:
	expr_or_star_exprs optional_comma
	{
		$$ = $1
	}

testlist: tests optional_comma

// (',' test ':' test)*
test_colon_tests:
	test ':' test
|	test_colon_tests ',' test ':' test

dictorsetmaker:
	test_colon_tests optional_comma
|	test ':' test comp_for
|	testlist
|	test comp_for

classdef:
	CLASS NAME optional_arglist_call ':' suite

arguments:
	argument
|	arguments ',' argument

optional_arguments:
|	arguments ','

arguments2:
|	arguments2 ',' argument

arglist:
	arguments optional_comma
|	optional_arguments '*' test arguments2
|	optional_arguments '*' test arguments2 ',' STARSTAR test
|	optional_arguments STARSTAR test

// The reason that keywords are test nodes instead of NAME is that using NAME
// results in an ambiguity. ast.c makes sure it's a NAME.
argument:
	test
|	test comp_for
|	test '=' test  // Really [keyword '='] test

comp_iter:
	comp_for
|	comp_if

comp_for:
	FOR exprlist IN or_test
|	FOR exprlist IN or_test comp_iter

comp_if:
	IF test_nocond
|	IF test_nocond comp_iter

// not used in grammar, but may appear in "node" passed from Parser to Compiler
// encoding_decl: NAME

yield_expr:
	YIELD
|	YIELD yield_arg

yield_arg:
	FROM test
|	testlist
