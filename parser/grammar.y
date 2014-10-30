%{

package parser

// Grammar for Python

import (
	"github.com/ncw/gpython/py"
)

%}

%union {
    str string
    obj py.Object
}

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

// FIXME figure out how to tell the parser to start from a given node
// inputs: single_input | file_input | eval_input
// In the mean time just do file_input
// inputs: single_input | file_input | eval_input
inputs: file_input

single_input: NEWLINE | simple_stmt | compound_stmt NEWLINE

// (NEWLINE | stmt)*
nl_or_stmt: | nl_or_stmt NEWLINE | nl_or_stmt stmt

//file_input: (NEWLINE | stmt)* ENDMARKER
file_input: nl_or_stmt ENDMARKER

// NEWLINE*
nls: | nls NEWLINE

//eval_input: testlist NEWLINE* ENDMARKER
eval_input: testlist nls ENDMARKER

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
     | tfpdeftest tfpdeftests ','
     | tfpdeftest tfpdeftests ',' '*' optional_tfpdef tfpdeftests
     | tfpdeftest tfpdeftests ',' '*' optional_tfpdef tfpdeftests ',' STARSTAR tfpdef
     | tfpdeftest tfpdeftests ',' STARSTAR tfpdef
     | '*' optional_tfpdef tfpdeftests
     | '*' optional_tfpdef tfpdeftests ',' STARSTAR tfpdef
     | STARSTAR tfpdef

tfpdef: NAME
      | NAME ':' test

vfpdeftest: vfpdef | vfpdef '=' test
vfpdeftests: | vfpdeftests ',' vfpdeftest
optional_vfpdef: | vfpdef
varargslist: vfpdeftest vfpdeftests
           | vfpdeftest vfpdeftests ','
           | vfpdeftest vfpdeftests ',' '*' optional_vfpdef vfpdeftests
           | vfpdeftest vfpdeftests ',' '*' optional_vfpdef vfpdeftests ',' STARSTAR vfpdef
           | vfpdeftest vfpdeftests ',' STARSTAR vfpdef
           | '*' optional_vfpdef vfpdeftests
           | '*' optional_vfpdef vfpdeftests ',' STARSTAR vfpdef
           | STARSTAR vfpdef

vfpdef: NAME

stmt: simple_stmt | compound_stmt
optional_semicolon: | ';'
small_stmts: small_stmt | small_stmts ';' small_stmt
simple_stmt: small_stmts optional_semicolon NEWLINE
small_stmt: expr_stmt | del_stmt | pass_stmt | flow_stmt |
             import_stmt | global_stmt | nonlocal_stmt | assert_stmt
yield_expr_or_testlist: yield_expr|testlist
yield_expr_or_testlist_star_expr: yield_expr|testlist_star_expr
equals_yield_expr_or_testlist_star_expr: | equals_yield_expr_or_testlist_star_expr '=' yield_expr_or_testlist_star_expr
expr_stmt: testlist_star_expr augassign yield_expr_or_testlist |
                     testlist_star_expr equals_yield_expr_or_testlist_star_expr
// FIXME this is making a reduce/reduce error for some reason :-(
test_or_star_expr: test | star_expr
test_or_star_exprs: test_or_star_expr | test_or_star_exprs ',' test_or_star_expr
optional_comma: | ','
testlist_star_expr: test_or_star_exprs optional_comma
augassign: PLUSEQ | MINUSEQ | STAREQ | DIVEQ | PERCEQ | ANDEQ | PIPEEQ | HATEQ |
            LTLTEQ | GTGTEQ | STARSTAREQ | DIVDIVEQ
// For normal assignments, additional restrictions enforced by the interpreter
del_stmt: DEL exprlist
pass_stmt: PASS
flow_stmt: break_stmt | continue_stmt | return_stmt | raise_stmt | yield_stmt
break_stmt: BREAK
continue_stmt: CONTINUE
return_stmt: RETURN | RETURN testlist
yield_stmt: yield_expr
raise_stmt: RAISE | RAISE test | RAISE test FROM test
import_stmt: import_name | import_from
import_name: IMPORT dotted_as_names
// note below: the ('.' | '...') is necessary because '...' is tokenized as ELLIPSIS
import_dots: | import_dots '.' | import_dots ELIPSIS
import_dots_plus: '.' | ELIPSIS | import_dots
from_arg: import_dots dotted_name | import_dots_plus
import_from_arg: '*' | '(' import_as_names ')' | import_as_names
import_from: FROM from_arg IMPORT import_from_arg
import_as_name: NAME | NAME AS NAME
dotted_as_name: dotted_name | dotted_name AS NAME
import_as_names: import_as_name optional_comma | import_as_name ',' import_as_names
dotted_as_names: dotted_as_name | dotted_as_names ',' dotted_as_name
dotted_name: NAME | dotted_name '.' NAME
names: NAME | names ',' NAME
global_stmt: GLOBAL names
nonlocal_stmt: NONLOCAL names
tests: test | tests ',' test
assert_stmt: ASSERT tests

compound_stmt: if_stmt | while_stmt | for_stmt | try_stmt | with_stmt | funcdef | classdef | decorated
elifs: | elifs ELIF test ':' suite
optional_else: | ELSE ':' suite
if_stmt: IF test ':' suite elifs optional_else
while_stmt: WHILE test ':' suite optional_else
for_stmt: FOR exprlist IN testlist ':' suite optional_else
except_clauses: | except_clauses except_clause ':' suite
try_stmt: TRY ':' suite except_clauses optional_else
 | TRY ':' suite except_clauses optional_else FINALLY ':' suite
 | TRY ':' suite FINALLY ':' suite
with_items: with_item | with_items ',' with_item
with_stmt: WITH with_items  ':' suite
with_item: test | test AS expr
// NB compile.c makes sure that the default except clause is last
except_clause: EXCEPT | EXCEPT test | EXCEPT test AS NAME
stmts: stmt | stmts stmt
suite: simple_stmt | NEWLINE INDENT stmts DEDENT

test: or_test | or_test IF or_test ELSE test | lambdef
test_nocond: or_test | lambdef_nocond
lambdef: LAMBDA ':' test | LAMBDA varargslist ':' test
lambdef_nocond: LAMBDA ':' test_nocond | LAMBDA varargslist ':' test_nocond
or_test: and_test | or_test OR and_test
and_test: not_test | and_test AND not_test
not_test: NOT not_test | comparison
comparison: expr | comparison comp_op expr
// <> LTGT isn't actually a valid comparison operator in Python. It's here for the
// sake of a __future__ import described in PEP 401
comp_op: '<'|'>'|EQEQ|GTEQ|LTEQ|LTGT|PLINGEQ|IN|NOT IN|IS|IS NOT
star_expr: '*' expr
expr: xor_expr | expr '|' xor_expr
xor_expr: and_expr | xor_expr '^' and_expr
and_expr: shift_expr | and_expr '&' shift_expr
shift_expr: arith_expr | shift_expr LTLT arith_expr| shift_expr GTGT arith_expr
arith_expr: term | arith_expr '+' term | arith_expr '-' term
term: factor | term '*' factor| term '/' factor| term '%' factor| term DIVDIV factor
factor: '+' factor | '-' factor | '~' factor | power
trailers: | trailers trailer
power: atom trailers | atom trailers STARSTAR factor
strings: STRING | strings STRING
atom: '(' ')' |
      '(' yield_expr ')' |
      '(' testlist_comp ')' |
       '[' ']' |
       '[' testlist_comp ']' |
       '{' '}' |
       '{' dictorsetmaker '}' |
       NAME | NUMBER | strings | ELIPSIS | NONE | TRUE | FALSE
testlist_comp: test_or_star_expr comp_for | test_or_star_expr test_or_star_exprs optional_comma
trailer: '(' ')' | '(' arglist ')' | '[' subscriptlist ']' | '.' NAME
subscripts: subscript | subscripts ',' subscript
subscriptlist: subscripts optional_comma
subscript: test
| ':'
| ':' sliceop
| ':' test
| ':' test sliceop
| test ':'
| test ':' sliceop
| test ':' test
| test ':' test sliceop
sliceop: ':' | ':' test
expr_or_star_expr: expr|star_expr
expr_or_star_exprs: expr_or_star_expr | expr_or_star_exprs ',' expr_or_star_expr
exprlist: expr_or_star_exprs optional_comma
testlist: tests optional_comma
// (',' test ':' test)*
test_colon_tests: test ':' test | test_colon_tests ',' test ':' test
dictorsetmaker: test_colon_tests optional_comma
                | test ':' test comp_for
                | test testlist
                | test comp_for

classdef: CLASS NAME optional_arglist_call ':' suite

arguments: argument | arguments ',' argument
optional_arguments: | arguments ','
arguments2: | arguments2 ',' argument
arglist: arguments optional_comma
       | optional_arguments '*' test arguments2
       | optional_arguments '*' test arguments2 ',' STARSTAR test
       | optional_arguments STARSTAR test
// The reason that keywords are test nodes instead of NAME is that using NAME
// results in an ambiguity. ast.c makes sure it's a NAME.
argument: test
        | test comp_for
        | test '=' test  // Really [keyword '='] test
comp_iter: comp_for | comp_if
comp_for: FOR exprlist IN or_test
        | FOR exprlist IN or_test comp_iter
comp_if: IF test_nocond
       | IF test_nocond comp_iter

// not used in grammar, but may appear in "node" passed from Parser to Compiler
// encoding_decl: NAME

yield_expr: YIELD
          | YIELD yield_arg
yield_arg: FROM test | testlist
