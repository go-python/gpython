//line grammar.y:2
package parser

import __yyfmt__ "fmt"

//line grammar.y:3
// Grammar for Python

import (
	"fmt"
	"github.com/go-python/gpython/ast"
	"github.com/go-python/gpython/py"
)

// NB can put code blocks in not just at the end

// Returns a Tuple if > 1 items or a trailing comma, otherwise returns
// the first item in elts
func tupleOrExpr(pos ast.Pos, elts []ast.Expr, optional_comma bool) ast.Expr {
	if optional_comma || len(elts) > 1 {
		return &ast.Tuple{ExprBase: ast.ExprBase{Pos: pos}, Elts: elts, Ctx: ast.Load}
	} else {
		return elts[0]
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

//line grammar.y:99
type yySymType struct {
	yys            int
	pos            ast.Pos // kept up to date by the lexer
	str            string
	obj            py.Object
	mod            ast.Mod
	stmt           ast.Stmt
	stmts          []ast.Stmt
	expr           ast.Expr
	exprs          []ast.Expr
	op             ast.OperatorNumber
	cmpop          ast.CmpOp
	comma          bool
	comprehensions []ast.Comprehension
	isExpr         bool
	slice          ast.Slicer
	call           *ast.Call
	level          int
	alias          *ast.Alias
	aliases        []*ast.Alias
	identifiers    []ast.Identifier
	ifstmt         *ast.If
	lastif         *ast.If
	exchandlers    []*ast.ExceptHandler
	withitem       *ast.WithItem
	withitems      []*ast.WithItem
	arg            *ast.Arg
	args           []*ast.Arg
	arguments      *ast.Arguments
}

const NEWLINE = 57346
const ENDMARKER = 57347
const NAME = 57348
const INDENT = 57349
const DEDENT = 57350
const STRING = 57351
const NUMBER = 57352
const PLINGEQ = 57353
const PERCEQ = 57354
const ANDEQ = 57355
const STARSTAR = 57356
const STARSTAREQ = 57357
const STAREQ = 57358
const PLUSEQ = 57359
const MINUSEQ = 57360
const MINUSGT = 57361
const ELIPSIS = 57362
const DIVDIV = 57363
const DIVDIVEQ = 57364
const DIVEQ = 57365
const LTLT = 57366
const LTLTEQ = 57367
const LTEQ = 57368
const LTGT = 57369
const EQEQ = 57370
const GTEQ = 57371
const GTGT = 57372
const GTGTEQ = 57373
const HATEQ = 57374
const PIPEEQ = 57375
const FALSE = 57376
const NONE = 57377
const TRUE = 57378
const AND = 57379
const AS = 57380
const ASSERT = 57381
const BREAK = 57382
const CLASS = 57383
const CONTINUE = 57384
const DEF = 57385
const DEL = 57386
const ELIF = 57387
const ELSE = 57388
const EXCEPT = 57389
const FINALLY = 57390
const FOR = 57391
const FROM = 57392
const GLOBAL = 57393
const IF = 57394
const IMPORT = 57395
const IN = 57396
const IS = 57397
const LAMBDA = 57398
const NONLOCAL = 57399
const NOT = 57400
const OR = 57401
const PASS = 57402
const RAISE = 57403
const RETURN = 57404
const TRY = 57405
const WHILE = 57406
const WITH = 57407
const YIELD = 57408
const SINGLE_INPUT = 57409
const FILE_INPUT = 57410
const EVAL_INPUT = 57411

var yyToknames = []string{
	"NEWLINE",
	"ENDMARKER",
	"NAME",
	"INDENT",
	"DEDENT",
	"STRING",
	"NUMBER",
	"PLINGEQ",
	"PERCEQ",
	"ANDEQ",
	"STARSTAR",
	"STARSTAREQ",
	"STAREQ",
	"PLUSEQ",
	"MINUSEQ",
	"MINUSGT",
	"ELIPSIS",
	"DIVDIV",
	"DIVDIVEQ",
	"DIVEQ",
	"LTLT",
	"LTLTEQ",
	"LTEQ",
	"LTGT",
	"EQEQ",
	"GTEQ",
	"GTGT",
	"GTGTEQ",
	"HATEQ",
	"PIPEEQ",
	"FALSE",
	"NONE",
	"TRUE",
	"AND",
	"AS",
	"ASSERT",
	"BREAK",
	"CLASS",
	"CONTINUE",
	"DEF",
	"DEL",
	"ELIF",
	"ELSE",
	"EXCEPT",
	"FINALLY",
	"FOR",
	"FROM",
	"GLOBAL",
	"IF",
	"IMPORT",
	"IN",
	"IS",
	"LAMBDA",
	"NONLOCAL",
	"NOT",
	"OR",
	"PASS",
	"RAISE",
	"RETURN",
	"TRY",
	"WHILE",
	"WITH",
	"YIELD",
	"'('",
	"')'",
	"'['",
	"']'",
	"':'",
	"','",
	"';'",
	"'+'",
	"'-'",
	"'*'",
	"'/'",
	"'|'",
	"'&'",
	"'<'",
	"'>'",
	"'='",
	"'.'",
	"'%'",
	"'{'",
	"'}'",
	"'^'",
	"'~'",
	"'@'",
	"SINGLE_INPUT",
	"FILE_INPUT",
	"EVAL_INPUT",
}
var yyStatenames = []string{}

const yyEofCode = 1
const yyErrCode = 2
const yyMaxDepth = 200

//line yacctab:1
var yyExca = []int{
	-1, 1,
	1, -1,
	-2, 0,
	-1, 232,
	68, 13,
	-2, 291,
	-1, 382,
	68, 91,
	-2, 292,
}

const yyNprod = 311
const yyPrivate = 57344

var yyTokenNames []string
var yyStates []string

const yyLast = 1441

var yyAct = []int{

	59, 468, 61, 314, 160, 97, 165, 164, 456, 421,
	401, 375, 321, 349, 361, 342, 141, 464, 224, 101,
	102, 6, 259, 111, 335, 223, 334, 103, 210, 69,
	318, 145, 457, 54, 60, 35, 237, 110, 95, 105,
	74, 71, 72, 66, 64, 146, 150, 73, 137, 57,
	106, 75, 231, 97, 143, 107, 70, 180, 24, 97,
	23, 17, 189, 288, 133, 247, 106, 277, 49, 124,
	125, 107, 130, 122, 120, 121, 2, 3, 4, 131,
	123, 284, 128, 96, 204, 139, 243, 232, 129, 127,
	126, 142, 378, 157, 138, 262, 236, 154, 181, 152,
	179, 315, 148, 243, 184, 185, 48, 227, 226, 420,
	243, 167, 211, 215, 387, 195, 166, 186, 187, 97,
	279, 222, 280, 336, 341, 188, 385, 166, 190, 191,
	192, 196, 199, 315, 166, 163, 281, 99, 474, 132,
	466, 312, 163, 453, 450, 390, 398, 395, 382, 373,
	197, 200, 234, 214, 251, 284, 289, 235, 252, 140,
	255, 216, 151, 257, 246, 175, 238, 206, 239, 260,
	261, 419, 241, 240, 221, 486, 473, 291, 258, 460,
	173, 174, 171, 172, 333, 403, 340, 413, 384, 412,
	244, 242, 480, 332, 411, 250, 249, 162, 263, 159,
	409, 253, 254, 311, 162, 405, 400, 379, 176, 178,
	370, 363, 177, 316, 285, 296, 256, 287, 219, 218,
	290, 97, 267, 293, 268, 271, 272, 111, 108, 283,
	269, 270, 286, 322, 169, 170, 266, 292, 273, 274,
	275, 276, 325, 397, 297, 298, 328, 357, 313, 356,
	452, 106, 396, 304, 381, 372, 107, 338, 305, 299,
	355, 300, 353, 343, 303, 339, 282, 232, 230, 337,
	238, 84, 239, 323, 91, 85, 284, 155, 329, 459,
	322, 350, 156, 156, 156, 87, 156, 265, 404, 264,
	358, 220, 359, 284, 248, 245, 459, 284, 461, 90,
	88, 89, 365, 367, 366, 407, 362, 134, 371, 362,
	346, 447, 354, 391, 106, 376, 377, 228, 158, 107,
	182, 211, 307, 67, 207, 302, 183, 374, 315, 344,
	34, 369, 81, 15, 82, 14, 383, 315, 392, 76,
	77, 166, 380, 166, 462, 483, 430, 260, 394, 467,
	83, 389, 402, 78, 136, 386, 114, 336, 315, 116,
	388, 117, 166, 393, 352, 399, 465, 330, 414, 139,
	433, 327, 324, 295, 294, 408, 135, 113, 112, 422,
	423, 326, 217, 322, 98, 425, 426, 211, 427, 410,
	212, 418, 406, 7, 100, 424, 417, 415, 213, 350,
	309, 436, 308, 432, 438, 428, 440, 439, 441, 431,
	229, 435, 434, 437, 310, 429, 161, 109, 301, 360,
	331, 144, 147, 376, 449, 443, 149, 317, 451, 320,
	319, 448, 348, 347, 168, 442, 25, 444, 445, 446,
	454, 119, 193, 203, 104, 458, 205, 455, 306, 364,
	202, 233, 68, 470, 62, 80, 278, 79, 463, 118,
	16, 432, 469, 115, 13, 12, 11, 322, 9, 475,
	10, 44, 43, 42, 478, 41, 481, 479, 484, 476,
	40, 39, 485, 469, 38, 33, 472, 487, 488, 469,
	209, 208, 84, 32, 31, 91, 85, 30, 29, 482,
	28, 27, 26, 368, 8, 93, 87, 94, 5, 92,
	1, 86, 0, 0, 0, 0, 0, 0, 0, 0,
	90, 88, 89, 0, 0, 47, 50, 24, 51, 23,
	36, 0, 0, 0, 0, 20, 56, 45, 18, 55,
	0, 0, 65, 46, 67, 0, 37, 53, 52, 21,
	19, 22, 58, 81, 84, 82, 416, 91, 85, 0,
	76, 77, 63, 0, 0, 0, 0, 0, 87, 0,
	0, 83, 0, 0, 78, 48, 0, 0, 0, 0,
	0, 0, 90, 88, 89, 0, 0, 47, 50, 24,
	51, 23, 36, 0, 0, 0, 0, 20, 56, 45,
	18, 55, 0, 0, 65, 46, 67, 0, 37, 53,
	52, 21, 19, 22, 58, 81, 84, 82, 0, 91,
	85, 0, 76, 77, 63, 0, 0, 0, 0, 0,
	87, 0, 0, 83, 0, 0, 78, 48, 0, 0,
	0, 0, 0, 0, 90, 88, 89, 0, 0, 47,
	50, 24, 51, 23, 36, 0, 0, 0, 0, 20,
	56, 45, 18, 55, 0, 0, 65, 46, 67, 0,
	37, 53, 52, 21, 19, 22, 58, 81, 0, 82,
	0, 0, 0, 0, 76, 77, 63, 225, 0, 84,
	0, 0, 91, 85, 0, 83, 0, 0, 78, 48,
	0, 0, 0, 87, 0, 0, 0, 0, 0, 0,
	0, 0, 0, 0, 0, 0, 0, 90, 88, 89,
	0, 0, 47, 50, 0, 51, 0, 36, 0, 0,
	0, 0, 0, 56, 45, 0, 55, 0, 0, 65,
	46, 67, 0, 37, 53, 52, 0, 0, 0, 58,
	81, 84, 82, 0, 91, 85, 0, 76, 77, 63,
	0, 0, 0, 0, 0, 87, 0, 0, 83, 0,
	0, 78, 0, 0, 0, 0, 0, 0, 0, 90,
	88, 89, 0, 0, 47, 50, 0, 51, 0, 36,
	0, 0, 0, 0, 0, 56, 45, 0, 55, 0,
	0, 65, 46, 67, 0, 37, 53, 52, 0, 0,
	0, 58, 81, 84, 82, 0, 91, 85, 0, 76,
	77, 63, 0, 0, 0, 0, 0, 87, 0, 0,
	83, 0, 0, 78, 0, 0, 0, 0, 0, 0,
	0, 90, 88, 89, 0, 0, 0, 0, 0, 0,
	84, 0, 0, 91, 85, 0, 0, 0, 0, 0,
	0, 0, 0, 65, 87, 67, 0, 0, 0, 0,
	0, 0, 0, 58, 81, 194, 82, 0, 90, 88,
	89, 76, 77, 63, 0, 0, 0, 84, 0, 0,
	91, 85, 83, 0, 0, 78, 0, 0, 0, 0,
	65, 87, 67, 0, 0, 0, 0, 0, 0, 0,
	58, 81, 0, 82, 0, 90, 88, 89, 76, 77,
	63, 0, 0, 0, 0, 0, 0, 0, 0, 83,
	84, 0, 78, 91, 85, 0, 0, 65, 477, 67,
	0, 0, 0, 0, 87, 0, 0, 0, 81, 0,
	82, 198, 0, 0, 0, 76, 77, 63, 90, 88,
	89, 0, 0, 0, 0, 0, 83, 84, 0, 78,
	91, 85, 0, 0, 0, 0, 0, 0, 0, 0,
	65, 87, 67, 0, 0, 0, 0, 0, 0, 0,
	0, 81, 0, 82, 0, 90, 88, 89, 76, 77,
	0, 0, 0, 0, 0, 0, 0, 0, 0, 83,
	84, 0, 78, 91, 85, 0, 0, 65, 0, 67,
	0, 0, 0, 0, 87, 0, 0, 0, 81, 0,
	82, 0, 403, 0, 0, 76, 77, 0, 90, 88,
	89, 0, 0, 0, 0, 0, 83, 0, 0, 78,
	0, 0, 84, 0, 0, 91, 85, 0, 0, 0,
	65, 0, 67, 0, 0, 0, 87, 0, 0, 0,
	0, 81, 0, 82, 0, 351, 0, 0, 76, 77,
	90, 88, 89, 0, 0, 0, 0, 0, 0, 83,
	0, 0, 78, 0, 84, 0, 0, 91, 85, 0,
	0, 0, 65, 0, 67, 0, 0, 0, 87, 0,
	0, 0, 0, 81, 345, 82, 0, 0, 0, 0,
	76, 77, 90, 88, 89, 0, 0, 0, 0, 0,
	0, 83, 0, 0, 78, 0, 0, 84, 0, 0,
	91, 85, 0, 0, 65, 0, 67, 0, 0, 0,
	0, 87, 0, 0, 0, 81, 0, 82, 0, 0,
	0, 0, 76, 77, 63, 90, 88, 89, 0, 0,
	0, 0, 0, 83, 84, 0, 78, 91, 85, 0,
	0, 0, 0, 0, 0, 0, 0, 65, 87, 67,
	0, 0, 0, 0, 0, 0, 0, 58, 81, 0,
	82, 0, 90, 88, 89, 76, 77, 0, 0, 0,
	0, 84, 0, 0, 91, 85, 83, 0, 0, 78,
	0, 0, 0, 0, 65, 87, 67, 0, 0, 0,
	0, 0, 0, 0, 0, 81, 0, 82, 0, 90,
	88, 89, 76, 77, 0, 0, 0, 0, 84, 0,
	0, 91, 85, 83, 201, 153, 78, 0, 0, 0,
	0, 65, 87, 67, 0, 0, 0, 0, 0, 0,
	0, 0, 81, 0, 82, 0, 90, 88, 89, 76,
	77, 0, 0, 0, 0, 84, 0, 0, 91, 85,
	83, 0, 0, 78, 0, 0, 0, 0, 471, 87,
	67, 0, 0, 0, 0, 0, 0, 0, 0, 81,
	0, 82, 0, 90, 88, 89, 76, 77, 0, 0,
	0, 0, 84, 0, 0, 91, 85, 83, 0, 0,
	78, 0, 0, 0, 0, 65, 87, 67, 0, 0,
	0, 0, 0, 0, 0, 0, 81, 0, 82, 0,
	90, 88, 89, 76, 77, 0, 0, 0, 84, 0,
	0, 91, 85, 0, 83, 0, 0, 78, 0, 0,
	0, 0, 87, 0, 0, 0, 0, 0, 0, 0,
	0, 0, 0, 81, 0, 82, 90, 88, 89, 0,
	76, 77, 63, 0, 0, 0, 0, 0, 0, 0,
	0, 83, 0, 0, 78, 0, 0, 0, 0, 0,
	0, 0, 0, 0, 0, 0, 0, 0, 0, 81,
	0, 82, 0, 0, 0, 0, 76, 77, 0, 0,
	0, 0, 0, 0, 0, 0, 0, 83, 0, 0,
	78,
}
var yyPact = []int{

	-14, -1000, 610, -1000, 1279, -1000, -1000, 380, 64, -1000,
	-1000, -1000, -1000, -1000, -1000, -1000, -1000, -1000, 1279, 1279,
	1316, 157, 1279, 372, 371, 17, -1000, -1000, -1000, -1000,
	-1000, -1000, -1000, -1000, -1000, 57, 1316, -1000, -1000, -1000,
	-1000, -1000, -1000, -1000, -1000, 370, 370, 1279, 363, 87,
	-1000, -1000, 1279, 1279, -1000, 363, 79, -1000, 1205, -1000,
	-1000, 225, -1000, 1352, 281, 128, -1000, 265, 154, 22,
	-30, 19, 296, 30, 41, -1000, 1352, 1352, 1352, -1000,
	-1000, 807, 881, 1168, -1000, -1000, 315, -1000, -1000, -1000,
	-1000, -1000, -1000, 486, -1000, -1000, 81, -1000, -1000, 745,
	378, 148, 147, 237, 102, -1000, 22, -1000, 683, 36,
	-1000, 279, 201, 200, -1000, -1000, -1000, -1000, 1131, 14,
	-1000, -1000, -1000, -1000, -1000, -1000, -1000, -1000, -1000, -1000,
	-1000, -1000, 844, -1000, 101, -1000, 101, 100, 20, -1000,
	1088, -1000, -1000, 245, 92, -1000, 27, 241, 3, 79,
	-1000, -1000, -1000, 1279, -1000, 265, 265, 22, 265, 1279,
	145, 91, 337, 337, -1000, 13, -1000, -1000, 1352, -1000,
	-1000, -1000, -1000, -1000, -1000, -1000, -1000, 235, 229, 1352,
	1352, 1352, 1352, 1352, 1352, 1352, 1352, 1352, 1352, 1352,
	-1000, -1000, -1000, 53, -1000, 198, 248, 87, -1000, 248,
	87, -1000, -23, 84, 106, -1000, 81, -1000, -1000, -1000,
	-1000, -1000, -1000, 369, 1279, -1000, -1000, -1000, 683, 683,
	1279, 1316, -1000, -1000, -1000, 318, 1279, 683, 1352, 303,
	127, 142, 1279, -1000, -1000, -1000, 844, -1000, -1000, -1000,
	366, 1279, 377, 365, -1000, 1279, 363, 361, 117, -1000,
	3, -1000, 223, 281, -1000, -1000, 1279, 110, -1000, -1000,
	-1000, -1000, 1279, 22, -1000, -1000, -30, 19, 296, 30,
	30, 41, 41, -1000, -1000, -1000, -1000, 1352, -1000, 1046,
	1004, 358, -1000, 194, 1316, 192, 179, 177, -1000, 1279,
	-1000, 1279, -1000, -1000, -1000, -1000, -1000, -1000, 263, 140,
	-1000, 256, 610, -1000, -1000, 22, 139, 1279, 187, -1000,
	77, 322, 322, -1000, 10, 136, 683, 186, -1000, 76,
	112, -1000, 32, -1000, -1000, -1000, -1000, -1000, -1000, -1000,
	-1000, -1000, -1000, 351, 73, -1000, 275, 1279, -1000, -1000,
	337, 337, 75, -1000, -1000, -1000, 184, 173, 74, -1000,
	135, 961, -1000, -1000, 234, -1000, -1000, -1000, 134, 248,
	260, -1000, 129, 683, 123, 118, 116, 1279, 548, -1000,
	683, -1000, -1000, 95, -1000, -1000, -1000, -1000, 1279, 1279,
	-1000, -1000, 1279, -1000, 1279, 1279, -1000, 1279, 73, -1000,
	351, 340, -1000, -1000, -1000, 356, -1000, -1000, 1004, -1000,
	961, -1000, 114, 1279, 265, 1279, -1000, 1279, -1000, 683,
	263, 683, 683, 683, 273, -1000, -1000, -1000, -1000, 322,
	322, 72, -1000, -1000, -1000, -1000, -1000, -1000, 182, -1000,
	-1000, 71, -1000, 337, -1000, -1000, 114, -1000, -1000, 227,
	-1000, 108, -1000, -1000, -1000, 250, -1000, 338, -1000, -1000,
	352, 68, -1000, 335, -1000, -1000, -1000, -1000, -1000, 1242,
	683, 105, -1000, 66, -1000, 322, 924, 337, 244, 224,
	-1000, 121, -1000, 683, 331, -1000, -1000, 1279, -1000, -1000,
	1242, 104, -1000, 322, -1000, -1000, 1242, -1000, -1000,
}
var yyPgo = []int{

	0, 511, 510, 509, 508, 507, 18, 28, 505, 504,
	503, 25, 14, 390, 61, 502, 501, 500, 498, 497,
	494, 493, 485, 484, 481, 480, 475, 473, 472, 471,
	470, 468, 466, 465, 464, 335, 333, 463, 460, 459,
	39, 29, 34, 56, 41, 42, 47, 40, 51, 457,
	456, 455, 49, 0, 43, 454, 1, 453, 2, 44,
	452, 38, 35, 451, 33, 36, 450, 10, 449, 448,
	330, 27, 446, 445, 8, 444, 68, 83, 443, 442,
	441, 436, 434, 16, 32, 13, 433, 432, 12, 430,
	429, 428, 30, 52, 427, 46, 426, 45, 422, 307,
	31, 24, 421, 26, 420, 419, 418, 37, 417, 7,
	6, 22, 17, 3, 11, 15, 416, 9, 414, 4,
	410, 402, 400, 398, 394,
}
var yyR1 = []int{

	0, 2, 2, 2, 4, 4, 3, 8, 8, 8,
	5, 123, 123, 94, 94, 93, 93, 70, 81, 81,
	37, 37, 38, 69, 69, 35, 120, 121, 121, 112,
	112, 117, 117, 118, 118, 114, 114, 122, 122, 122,
	122, 122, 122, 122, 113, 113, 109, 109, 115, 115,
	116, 116, 111, 111, 119, 119, 119, 119, 119, 119,
	119, 110, 7, 7, 124, 124, 9, 9, 6, 14,
	14, 14, 14, 14, 14, 14, 14, 15, 15, 15,
	63, 63, 65, 65, 80, 80, 76, 76, 52, 52,
	83, 83, 62, 39, 39, 39, 39, 39, 39, 39,
	39, 39, 39, 39, 39, 16, 17, 18, 18, 18,
	18, 18, 23, 24, 25, 25, 27, 26, 26, 26,
	19, 19, 28, 95, 95, 96, 96, 98, 98, 98,
	104, 104, 104, 29, 101, 101, 100, 100, 103, 103,
	102, 102, 97, 97, 99, 99, 20, 21, 77, 77,
	22, 22, 13, 13, 13, 13, 13, 13, 13, 13,
	105, 105, 12, 12, 31, 30, 32, 106, 106, 33,
	33, 33, 33, 108, 108, 34, 107, 107, 68, 68,
	68, 10, 10, 11, 11, 53, 53, 53, 56, 56,
	55, 55, 57, 57, 58, 58, 59, 59, 54, 54,
	60, 60, 82, 82, 82, 82, 82, 82, 82, 82,
	82, 82, 82, 42, 41, 41, 43, 43, 44, 44,
	45, 45, 45, 46, 46, 46, 47, 47, 47, 47,
	47, 48, 48, 48, 48, 49, 49, 79, 79, 1,
	1, 51, 51, 51, 51, 51, 51, 51, 51, 51,
	51, 51, 51, 51, 51, 51, 51, 50, 50, 50,
	50, 87, 87, 86, 85, 85, 85, 85, 85, 85,
	85, 85, 85, 67, 67, 40, 40, 75, 75, 71,
	61, 72, 78, 78, 66, 66, 66, 66, 36, 89,
	89, 90, 90, 91, 91, 92, 92, 92, 92, 88,
	88, 88, 74, 74, 84, 84, 73, 73, 64, 64,
	64,
}
var yyR2 = []int{

	0, 2, 2, 2, 1, 2, 2, 0, 2, 2,
	3, 0, 2, 0, 1, 0, 3, 4, 1, 2,
	1, 1, 2, 0, 2, 6, 3, 0, 1, 1,
	3, 0, 3, 1, 3, 0, 1, 2, 5, 8,
	4, 3, 6, 2, 1, 3, 1, 3, 0, 3,
	1, 3, 0, 1, 2, 5, 8, 4, 3, 6,
	2, 1, 1, 1, 0, 1, 1, 3, 3, 1,
	1, 1, 1, 1, 1, 1, 1, 3, 2, 1,
	1, 1, 1, 1, 2, 3, 1, 3, 1, 1,
	0, 1, 2, 1, 1, 1, 1, 1, 1, 1,
	1, 1, 1, 1, 1, 2, 1, 1, 1, 1,
	1, 1, 1, 1, 1, 2, 1, 1, 2, 4,
	1, 1, 2, 1, 1, 1, 2, 1, 2, 1,
	1, 4, 2, 4, 1, 3, 1, 3, 1, 3,
	1, 3, 1, 3, 1, 3, 2, 2, 1, 3,
	2, 4, 1, 1, 1, 1, 1, 1, 1, 1,
	0, 5, 0, 3, 6, 5, 7, 0, 4, 4,
	7, 7, 10, 1, 3, 4, 1, 3, 1, 2,
	4, 1, 2, 1, 4, 1, 5, 1, 1, 1,
	3, 4, 3, 4, 1, 3, 1, 3, 2, 1,
	1, 3, 1, 1, 1, 1, 1, 1, 1, 1,
	2, 1, 2, 2, 1, 3, 1, 3, 1, 3,
	1, 3, 3, 1, 3, 3, 1, 3, 3, 3,
	3, 2, 2, 2, 1, 2, 4, 0, 2, 1,
	2, 2, 3, 4, 4, 2, 4, 4, 2, 3,
	1, 1, 1, 1, 1, 1, 1, 2, 3, 3,
	2, 1, 3, 2, 1, 1, 2, 2, 3, 2,
	3, 3, 4, 1, 2, 1, 1, 1, 3, 2,
	2, 2, 3, 5, 2, 4, 1, 2, 5, 1,
	3, 0, 2, 0, 3, 2, 4, 7, 3, 1,
	2, 3, 1, 1, 4, 5, 2, 3, 1, 3,
	2,
}
var yyChk = []int{

	-1000, -2, 90, 91, 92, -4, -6, -13, -9, -31,
	-30, -32, -33, -34, -35, -36, -38, -14, 52, 64,
	49, 63, 65, 43, 41, -81, -15, -16, -17, -18,
	-19, -20, -21, -22, -70, -62, 44, 60, -23, -24,
	-25, -26, -27, -28, -29, 51, 57, 39, 89, -76,
	40, 42, 62, 61, -64, 53, 50, -52, 66, -53,
	-42, -58, -55, 76, -59, 56, -54, 58, -60, -41,
	-43, -44, -45, -46, -47, -48, 74, 75, 88, -49,
	-51, 67, 69, 85, 6, 10, -1, 20, 35, 36,
	34, 9, -3, -8, -5, -61, -77, -53, 4, 73,
	-124, -53, -53, -71, -75, -40, -41, -42, 71, -108,
	-107, -53, 6, 6, -70, -37, -36, -35, -39, -80,
	17, 18, 16, 23, 12, 13, 33, 32, 25, 31,
	15, 22, 82, -71, -99, 6, -99, -53, -97, 6,
	72, -83, -61, -53, -102, -100, -97, -98, -97, -96,
	-95, 83, 20, 50, -61, 52, 59, -41, 37, 71,
	-119, -116, 76, 14, -109, -110, 6, -54, -82, 80,
	81, 28, 29, 26, 27, 11, 54, 58, 55, 78,
	87, 79, 24, 30, 74, 75, 76, 77, 84, 21,
	-48, -48, -48, -79, 68, -64, -52, -76, 70, -52,
	-76, 86, -66, -78, -53, -72, -77, 9, 5, 4,
	-7, -6, -13, -123, 72, -83, -14, 4, 71, 71,
	54, 72, -83, -11, -6, 4, 72, 71, 38, -120,
	67, -93, 67, -63, -64, -61, 82, -65, -64, -62,
	72, 72, -93, 83, -52, 50, 72, 38, 53, -95,
	-97, -53, -58, -59, -54, -53, 71, 72, -83, -111,
	-110, -110, 82, -41, 54, 58, -43, -44, -45, -46,
	-46, -47, -47, -48, -48, -48, -48, 14, -50, 67,
	69, 83, 68, -84, 49, -83, -84, -83, 86, 72,
	-83, 71, -84, -83, 5, 4, -53, -11, -11, -61,
	-40, -106, 7, -107, -11, -41, -69, 19, -121, -122,
	-118, 76, 14, -112, -113, 6, 71, -94, -92, -89,
	-90, -88, -53, -65, 6, -53, 4, 6, -53, -100,
	6, -104, 76, 67, -103, -101, 6, 46, -53, -109,
	76, 14, -115, -53, -48, 68, -92, -86, -87, -85,
	-53, 71, 6, 68, -71, 68, 70, 70, -53, -53,
	-105, -12, 46, 71, -68, 46, 48, 47, -10, -7,
	71, -53, 68, 72, -83, -114, -113, -113, 82, 71,
	-11, 68, 72, -83, 76, 14, -84, 82, -103, -83,
	72, 38, -53, -111, -110, 72, 68, 70, 72, -83,
	71, -67, -53, 71, 54, 71, -84, 45, -12, 71,
	-11, 71, 71, 71, -53, -7, 8, -11, -112, 76,
	14, -117, -53, -53, -88, -53, -53, -53, -83, -101,
	6, -115, -109, 14, -85, -67, -53, -67, -53, -58,
	-53, -53, -11, -12, -11, -11, -11, 38, -114, -113,
	72, -91, 68, 72, -110, -67, -74, -84, -73, 52,
	71, 48, 6, -117, -112, 14, 72, 14, -56, -58,
	-57, 56, -11, 71, 72, -113, -88, 14, -110, -74,
	71, -119, -11, 14, -53, -56, 71, -113, -56,
}
var yyDef = []int{

	0, -2, 0, 7, 0, 1, 4, 0, 64, 152,
	153, 154, 155, 156, 157, 158, 159, 66, 0, 0,
	0, 0, 0, 0, 0, 0, 69, 70, 71, 72,
	73, 74, 75, 76, 18, 79, 0, 106, 107, 108,
	109, 110, 111, 120, 121, 0, 0, 0, 0, 90,
	112, 113, 114, 117, 116, 0, 0, 86, 308, 88,
	89, 185, 187, 0, 194, 0, 196, 0, 199, 200,
	214, 216, 218, 220, 223, 226, 0, 0, 0, 234,
	237, 0, 0, 0, 250, 251, 252, 253, 254, 255,
	256, 239, 2, 0, 3, 11, 90, 148, 5, 65,
	0, 0, 0, 0, 90, 277, 275, 276, 0, 0,
	173, 176, 0, 15, 19, 22, 20, 21, 0, 78,
	93, 94, 95, 96, 97, 98, 99, 100, 101, 102,
	103, 104, 0, 105, 146, 144, 147, 150, 15, 142,
	91, 92, 115, 118, 122, 140, 136, 0, 127, 129,
	125, 123, 124, 0, 310, 0, 0, 213, 0, 0,
	0, 90, 52, 0, 50, 46, 61, 198, 0, 202,
	203, 204, 205, 206, 207, 208, 209, 0, 211, 0,
	0, 0, 0, 0, 0, 0, 0, 0, 0, 0,
	231, 232, 233, 235, 241, 0, 86, 90, 245, 86,
	90, 248, 0, 90, 148, 286, 90, 240, 6, 8,
	9, 62, 63, 0, 91, 280, 67, 68, 0, 0,
	0, 91, 279, 167, 183, 0, 0, 0, 0, 23,
	27, 0, -2, 77, 80, 81, 0, 84, 82, 83,
	0, 0, 0, 0, 87, 0, 0, 0, 0, 126,
	128, 309, 0, 195, 197, 190, 0, 91, 54, 48,
	53, 60, 0, 201, 210, 212, 215, 217, 219, 221,
	222, 224, 225, 227, 228, 229, 230, 0, 238, 291,
	0, 0, 242, 0, 0, 0, 0, 0, 249, 91,
	284, 0, 287, 281, 10, 12, 149, 160, 162, 0,
	278, 169, 0, 174, 175, 177, 0, 0, 0, 28,
	90, 35, 0, 33, 29, 44, 0, 0, 14, 90,
	0, 289, 299, 85, 145, 151, 17, 143, 119, 141,
	137, 133, 130, 0, 90, 138, 134, 0, 191, 51,
	52, 0, 58, 47, 236, 257, 0, 0, 90, 261,
	264, 265, 260, 243, 0, 244, 246, 247, 0, 282,
	162, 165, 0, 0, 0, 0, 0, 178, 0, 181,
	0, 24, 26, 91, 37, 31, 36, 43, 0, 0,
	288, 16, -2, 295, 0, 0, 300, 0, 90, 132,
	91, 0, 186, 48, 57, 0, 258, 259, 91, 263,
	269, 266, 267, 273, 0, 0, 285, 0, 164, 0,
	162, 0, 0, 0, 179, 182, 184, 25, 34, 35,
	0, 41, 30, 45, 290, 293, 298, 301, 0, 139,
	135, 55, 49, 0, 262, 270, 271, 268, 274, 304,
	283, 0, 163, 166, 168, 170, 171, 0, 31, 40,
	0, 296, 131, 0, 59, 272, 305, 302, 303, 0,
	0, 0, 180, 38, 32, 0, 0, 0, 306, 188,
	189, 0, 161, 0, 0, 42, 294, 0, 56, 307,
	0, 0, 172, 0, 297, 192, 0, 39, 193,
}
var yyTok1 = []int{

	1, 3, 3, 3, 3, 3, 3, 3, 3, 3,
	3, 3, 3, 3, 3, 3, 3, 3, 3, 3,
	3, 3, 3, 3, 3, 3, 3, 3, 3, 3,
	3, 3, 3, 3, 3, 3, 3, 84, 79, 3,
	67, 68, 76, 74, 72, 75, 83, 77, 3, 3,
	3, 3, 3, 3, 3, 3, 3, 3, 71, 73,
	80, 82, 81, 3, 89, 3, 3, 3, 3, 3,
	3, 3, 3, 3, 3, 3, 3, 3, 3, 3,
	3, 3, 3, 3, 3, 3, 3, 3, 3, 3,
	3, 69, 3, 70, 87, 3, 3, 3, 3, 3,
	3, 3, 3, 3, 3, 3, 3, 3, 3, 3,
	3, 3, 3, 3, 3, 3, 3, 3, 3, 3,
	3, 3, 3, 85, 78, 86, 88,
}
var yyTok2 = []int{

	2, 3, 4, 5, 6, 7, 8, 9, 10, 11,
	12, 13, 14, 15, 16, 17, 18, 19, 20, 21,
	22, 23, 24, 25, 26, 27, 28, 29, 30, 31,
	32, 33, 34, 35, 36, 37, 38, 39, 40, 41,
	42, 43, 44, 45, 46, 47, 48, 49, 50, 51,
	52, 53, 54, 55, 56, 57, 58, 59, 60, 61,
	62, 63, 64, 65, 66, 90, 91, 92,
}
var yyTok3 = []int{
	0,
}

//line yaccpar:1

/*	parser for yacc output	*/

var yyDebug = 0

type yyLexer interface {
	Lex(lval *yySymType) int
	Error(s string)
}

const yyFlag = -1000

func yyTokname(c int) string {
	// 4 is TOKSTART above
	if c >= 4 && c-4 < len(yyToknames) {
		if yyToknames[c-4] != "" {
			return yyToknames[c-4]
		}
	}
	return __yyfmt__.Sprintf("tok-%v", c)
}

func yyStatname(s int) string {
	if s >= 0 && s < len(yyStatenames) {
		if yyStatenames[s] != "" {
			return yyStatenames[s]
		}
	}
	return __yyfmt__.Sprintf("state-%v", s)
}

func yylex1(lex yyLexer, lval *yySymType) int {
	c := 0
	char := lex.Lex(lval)
	if char <= 0 {
		c = yyTok1[0]
		goto out
	}
	if char < len(yyTok1) {
		c = yyTok1[char]
		goto out
	}
	if char >= yyPrivate {
		if char < yyPrivate+len(yyTok2) {
			c = yyTok2[char-yyPrivate]
			goto out
		}
	}
	for i := 0; i < len(yyTok3); i += 2 {
		c = yyTok3[i+0]
		if c == char {
			c = yyTok3[i+1]
			goto out
		}
	}

out:
	if c == 0 {
		c = yyTok2[1] /* unknown char */
	}
	if yyDebug >= 3 {
		__yyfmt__.Printf("lex %s(%d)\n", yyTokname(c), uint(char))
	}
	return c
}

func yyParse(yylex yyLexer) int {
	var yyn int
	var yylval yySymType
	var yyVAL yySymType
	yyS := make([]yySymType, yyMaxDepth)

	Nerrs := 0   /* number of errors */
	Errflag := 0 /* error recovery flag */
	yystate := 0
	yychar := -1
	yyp := -1
	goto yystack

ret0:
	return 0

ret1:
	return 1

yystack:
	/* put a state and value onto the stack */
	if yyDebug >= 4 {
		__yyfmt__.Printf("char %v in %v\n", yyTokname(yychar), yyStatname(yystate))
	}

	yyp++
	if yyp >= len(yyS) {
		nyys := make([]yySymType, len(yyS)*2)
		copy(nyys, yyS)
		yyS = nyys
	}
	yyS[yyp] = yyVAL
	yyS[yyp].yys = yystate

yynewstate:
	yyn = yyPact[yystate]
	if yyn <= yyFlag {
		goto yydefault /* simple state */
	}
	if yychar < 0 {
		yychar = yylex1(yylex, &yylval)
	}
	yyn += yychar
	if yyn < 0 || yyn >= yyLast {
		goto yydefault
	}
	yyn = yyAct[yyn]
	if yyChk[yyn] == yychar { /* valid shift */
		yychar = -1
		yyVAL = yylval
		yystate = yyn
		if Errflag > 0 {
			Errflag--
		}
		goto yystack
	}

yydefault:
	/* default state action */
	yyn = yyDef[yystate]
	if yyn == -2 {
		if yychar < 0 {
			yychar = yylex1(yylex, &yylval)
		}

		/* look through exception table */
		xi := 0
		for {
			if yyExca[xi+0] == -1 && yyExca[xi+1] == yystate {
				break
			}
			xi += 2
		}
		for xi += 2; ; xi += 2 {
			yyn = yyExca[xi+0]
			if yyn < 0 || yyn == yychar {
				break
			}
		}
		yyn = yyExca[xi+1]
		if yyn < 0 {
			goto ret0
		}
	}
	if yyn == 0 {
		/* error ... attempt to resume parsing */
		switch Errflag {
		case 0: /* brand new error */
			yylex.Error("syntax error")
			Nerrs++
			if yyDebug >= 1 {
				__yyfmt__.Printf("%s", yyStatname(yystate))
				__yyfmt__.Printf(" saw %s\n", yyTokname(yychar))
			}
			fallthrough

		case 1, 2: /* incompletely recovered error ... try again */
			Errflag = 3

			/* find a state where "error" is a legal shift action */
			for yyp >= 0 {
				yyn = yyPact[yyS[yyp].yys] + yyErrCode
				if yyn >= 0 && yyn < yyLast {
					yystate = yyAct[yyn] /* simulate a shift of "error" */
					if yyChk[yystate] == yyErrCode {
						goto yystack
					}
				}

				/* the current p has no shift on "error", pop stack */
				if yyDebug >= 2 {
					__yyfmt__.Printf("error recovery pops state %d\n", yyS[yyp].yys)
				}
				yyp--
			}
			/* there is no state on the stack with an error shift ... abort */
			goto ret1

		case 3: /* no shift yet; clobber input char */
			if yyDebug >= 2 {
				__yyfmt__.Printf("error recovery discards %s\n", yyTokname(yychar))
			}
			if yychar == yyEofCode {
				goto ret1
			}
			yychar = -1
			goto yynewstate /* try again in the same state */
		}
	}

	/* reduction by production yyn */
	if yyDebug >= 2 {
		__yyfmt__.Printf("reduce %v in:\n\t%v\n", yyn, yyStatname(yystate))
	}

	yynt := yyn
	yypt := yyp
	_ = yypt // guard against "declared and not used"

	yyp -= yyR2[yyn]
	yyVAL = yyS[yyp+1]

	/* consult goto table to find next state */
	yyn = yyR1[yyn]
	yyg := yyPgo[yyn]
	yyj := yyg + yyS[yyp].yys + 1

	if yyj >= yyLast {
		yystate = yyAct[yyg]
	} else {
		yystate = yyAct[yyj]
		if yyChk[yystate] != -yyn {
			yystate = yyAct[yyg]
		}
	}
	// dummy call; replaced with literal code
	switch yynt {

	case 1:
		//line grammar.y:246
		{
			yylex.(*yyLex).mod = yyS[yypt-0].mod
			return 0
		}
	case 2:
		//line grammar.y:251
		{
			yylex.(*yyLex).mod = yyS[yypt-0].mod
			return 0
		}
	case 3:
		//line grammar.y:256
		{
			yylex.(*yyLex).mod = yyS[yypt-0].mod
			return 0
		}
	case 4:
		//line grammar.y:270
		{
			yyVAL.mod = &ast.Interactive{ModBase: ast.ModBase{Pos: yyVAL.pos}, Body: yyS[yypt-0].stmts}
		}
	case 5:
		//line grammar.y:274
		{
			//  NB: compound_stmt in single_input is followed by extra NEWLINE!
			yyVAL.mod = &ast.Interactive{ModBase: ast.ModBase{Pos: yyVAL.pos}, Body: []ast.Stmt{yyS[yypt-1].stmt}}
		}
	case 6:
		//line grammar.y:282
		{
			yyVAL.mod = &ast.Module{ModBase: ast.ModBase{Pos: yyVAL.pos}, Body: yyS[yypt-1].stmts}
		}
	case 7:
		//line grammar.y:288
		{
			yyVAL.stmts = nil
		}
	case 8:
		//line grammar.y:292
		{
		}
	case 9:
		//line grammar.y:295
		{
			yyVAL.stmts = append(yyVAL.stmts, yyS[yypt-0].stmts...)
		}
	case 10:
		//line grammar.y:302
		{
			yyVAL.mod = &ast.Expression{ModBase: ast.ModBase{Pos: yyVAL.pos}, Body: yyS[yypt-2].expr}
		}
	case 13:
		//line grammar.y:311
		{
			yyVAL.call = &ast.Call{ExprBase: ast.ExprBase{Pos: yyVAL.pos}}
		}
	case 14:
		//line grammar.y:315
		{
			yyVAL.call = yyS[yypt-0].call
		}
	case 15:
		//line grammar.y:320
		{
			yyVAL.call = nil
		}
	case 16:
		//line grammar.y:324
		{
			yyVAL.call = yyS[yypt-1].call
		}
	case 17:
		//line grammar.y:330
		{
			fn := &ast.Name{ExprBase: ast.ExprBase{Pos: yyVAL.pos}, Id: ast.Identifier(yyS[yypt-2].str), Ctx: ast.Load}
			if yyS[yypt-1].call == nil {
				yyVAL.expr = fn
			} else {
				call := *yyS[yypt-1].call
				call.Func = fn
				yyVAL.expr = &call
			}
		}
	case 18:
		//line grammar.y:343
		{
			yyVAL.exprs = nil
			yyVAL.exprs = append(yyVAL.exprs, yyS[yypt-0].expr)
		}
	case 19:
		//line grammar.y:348
		{
			yyVAL.exprs = append(yyVAL.exprs, yyS[yypt-0].expr)
		}
	case 20:
		//line grammar.y:354
		{
			yyVAL.stmt = yyS[yypt-0].stmt
		}
	case 21:
		//line grammar.y:358
		{
			yyVAL.stmt = yyS[yypt-0].stmt
		}
	case 22:
		//line grammar.y:364
		{
			switch x := (yyS[yypt-0].stmt).(type) {
			case *ast.ClassDef:
				x.DecoratorList = yyS[yypt-1].exprs
				yyVAL.stmt = x
			case *ast.FunctionDef:
				x.DecoratorList = yyS[yypt-1].exprs
				yyVAL.stmt = x
			default:
				panic("bad type for decorated")
			}
		}
	case 23:
		//line grammar.y:378
		{
			yyVAL.expr = nil
		}
	case 24:
		//line grammar.y:382
		{
			yyVAL.expr = yyS[yypt-0].expr
		}
	case 25:
		//line grammar.y:388
		{
			yyVAL.stmt = &ast.FunctionDef{StmtBase: ast.StmtBase{Pos: yyVAL.pos}, Name: ast.Identifier(yyS[yypt-4].str), Args: yyS[yypt-3].arguments, Body: yyS[yypt-0].stmts, Returns: yyS[yypt-2].expr}
		}
	case 26:
		//line grammar.y:394
		{
			yyVAL.arguments = yyS[yypt-1].arguments
		}
	case 27:
		//line grammar.y:399
		{
			yyVAL.arguments = &ast.Arguments{Pos: yyVAL.pos}
		}
	case 28:
		//line grammar.y:403
		{
			yyVAL.arguments = yyS[yypt-0].arguments
		}
	case 29:
		//line grammar.y:410
		{
			yyVAL.arg = yyS[yypt-0].arg
			yyVAL.expr = nil
		}
	case 30:
		//line grammar.y:415
		{
			yyVAL.arg = yyS[yypt-2].arg
			yyVAL.expr = yyS[yypt-0].expr
		}
	case 31:
		//line grammar.y:421
		{
			yyVAL.args = nil
			yyVAL.exprs = nil
		}
	case 32:
		//line grammar.y:426
		{
			yyVAL.args = append(yyVAL.args, yyS[yypt-0].arg)
			if yyS[yypt-0].expr != nil {
				yyVAL.exprs = append(yyVAL.exprs, yyS[yypt-0].expr)
			}
		}
	case 33:
		//line grammar.y:435
		{
			yyVAL.args = nil
			yyVAL.args = append(yyVAL.args, yyS[yypt-0].arg)
			yyVAL.exprs = nil
			if yyS[yypt-0].expr != nil {
				yyVAL.exprs = append(yyVAL.exprs, yyS[yypt-0].expr)
			}
		}
	case 34:
		//line grammar.y:444
		{
			yyVAL.args = append(yyVAL.args, yyS[yypt-0].arg)
			if yyS[yypt-0].expr != nil {
				yyVAL.exprs = append(yyVAL.exprs, yyS[yypt-0].expr)
			}
		}
	case 35:
		//line grammar.y:452
		{
			yyVAL.arg = nil
		}
	case 36:
		//line grammar.y:456
		{
			yyVAL.arg = yyS[yypt-0].arg
		}
	case 37:
		//line grammar.y:463
		{
			yyVAL.arguments = &ast.Arguments{Pos: yyVAL.pos, Args: yyS[yypt-1].args, Defaults: yyS[yypt-1].exprs}
		}
	case 38:
		//line grammar.y:467
		{
			yyVAL.arguments = &ast.Arguments{Pos: yyVAL.pos, Args: yyS[yypt-4].args, Defaults: yyS[yypt-4].exprs, Vararg: yyS[yypt-1].arg, Kwonlyargs: yyS[yypt-0].args, KwDefaults: yyS[yypt-0].exprs}
		}
	case 39:
		//line grammar.y:471
		{
			yyVAL.arguments = &ast.Arguments{Pos: yyVAL.pos, Args: yyS[yypt-7].args, Defaults: yyS[yypt-7].exprs, Vararg: yyS[yypt-4].arg, Kwonlyargs: yyS[yypt-3].args, KwDefaults: yyS[yypt-3].exprs, Kwarg: yyS[yypt-0].arg}
		}
	case 40:
		//line grammar.y:475
		{
			yyVAL.arguments = &ast.Arguments{Pos: yyVAL.pos, Args: yyS[yypt-3].args, Defaults: yyS[yypt-3].exprs, Kwarg: yyS[yypt-0].arg}
		}
	case 41:
		//line grammar.y:479
		{
			yyVAL.arguments = &ast.Arguments{Pos: yyVAL.pos, Vararg: yyS[yypt-1].arg, Kwonlyargs: yyS[yypt-0].args, KwDefaults: yyS[yypt-0].exprs}
		}
	case 42:
		//line grammar.y:483
		{
			yyVAL.arguments = &ast.Arguments{Pos: yyVAL.pos, Vararg: yyS[yypt-4].arg, Kwonlyargs: yyS[yypt-3].args, KwDefaults: yyS[yypt-3].exprs, Kwarg: yyS[yypt-0].arg}
		}
	case 43:
		//line grammar.y:487
		{
			yyVAL.arguments = &ast.Arguments{Pos: yyVAL.pos, Kwarg: yyS[yypt-0].arg}
		}
	case 44:
		//line grammar.y:493
		{
			yyVAL.arg = &ast.Arg{Pos: yyVAL.pos, Arg: ast.Identifier(yyS[yypt-0].str)}
		}
	case 45:
		//line grammar.y:497
		{
			yyVAL.arg = &ast.Arg{Pos: yyVAL.pos, Arg: ast.Identifier(yyS[yypt-2].str), Annotation: yyS[yypt-0].expr}
		}
	case 46:
		//line grammar.y:503
		{
			yyVAL.arg = yyS[yypt-0].arg
			yyVAL.expr = nil
		}
	case 47:
		//line grammar.y:508
		{
			yyVAL.arg = yyS[yypt-2].arg
			yyVAL.expr = yyS[yypt-0].expr
		}
	case 48:
		//line grammar.y:514
		{
			yyVAL.args = nil
			yyVAL.exprs = nil
		}
	case 49:
		//line grammar.y:519
		{
			yyVAL.args = append(yyVAL.args, yyS[yypt-0].arg)
			if yyS[yypt-0].expr != nil {
				yyVAL.exprs = append(yyVAL.exprs, yyS[yypt-0].expr)
			}
		}
	case 50:
		//line grammar.y:528
		{
			yyVAL.args = nil
			yyVAL.args = append(yyVAL.args, yyS[yypt-0].arg)
			yyVAL.exprs = nil
			if yyS[yypt-0].expr != nil {
				yyVAL.exprs = append(yyVAL.exprs, yyS[yypt-0].expr)
			}
		}
	case 51:
		//line grammar.y:537
		{
			yyVAL.args = append(yyVAL.args, yyS[yypt-0].arg)
			if yyS[yypt-0].expr != nil {
				yyVAL.exprs = append(yyVAL.exprs, yyS[yypt-0].expr)
			}
		}
	case 52:
		//line grammar.y:545
		{
			yyVAL.arg = nil
		}
	case 53:
		//line grammar.y:549
		{
			yyVAL.arg = yyS[yypt-0].arg
		}
	case 54:
		//line grammar.y:556
		{
			yyVAL.arguments = &ast.Arguments{Pos: yyVAL.pos, Args: yyS[yypt-1].args, Defaults: yyS[yypt-1].exprs}
		}
	case 55:
		//line grammar.y:560
		{
			yyVAL.arguments = &ast.Arguments{Pos: yyVAL.pos, Args: yyS[yypt-4].args, Defaults: yyS[yypt-4].exprs, Vararg: yyS[yypt-1].arg, Kwonlyargs: yyS[yypt-0].args, KwDefaults: yyS[yypt-0].exprs}
		}
	case 56:
		//line grammar.y:564
		{
			yyVAL.arguments = &ast.Arguments{Pos: yyVAL.pos, Args: yyS[yypt-7].args, Defaults: yyS[yypt-7].exprs, Vararg: yyS[yypt-4].arg, Kwonlyargs: yyS[yypt-3].args, KwDefaults: yyS[yypt-3].exprs, Kwarg: yyS[yypt-0].arg}
		}
	case 57:
		//line grammar.y:568
		{
			yyVAL.arguments = &ast.Arguments{Pos: yyVAL.pos, Args: yyS[yypt-3].args, Defaults: yyS[yypt-3].exprs, Kwarg: yyS[yypt-0].arg}
		}
	case 58:
		//line grammar.y:572
		{
			yyVAL.arguments = &ast.Arguments{Pos: yyVAL.pos, Vararg: yyS[yypt-1].arg, Kwonlyargs: yyS[yypt-0].args, KwDefaults: yyS[yypt-0].exprs}
		}
	case 59:
		//line grammar.y:576
		{
			yyVAL.arguments = &ast.Arguments{Pos: yyVAL.pos, Vararg: yyS[yypt-4].arg, Kwonlyargs: yyS[yypt-3].args, KwDefaults: yyS[yypt-3].exprs, Kwarg: yyS[yypt-0].arg}
		}
	case 60:
		//line grammar.y:580
		{
			yyVAL.arguments = &ast.Arguments{Pos: yyVAL.pos, Kwarg: yyS[yypt-0].arg}
		}
	case 61:
		//line grammar.y:586
		{
			yyVAL.arg = &ast.Arg{Pos: yyVAL.pos, Arg: ast.Identifier(yyS[yypt-0].str)}
		}
	case 62:
		//line grammar.y:592
		{
			yyVAL.stmts = yyS[yypt-0].stmts
		}
	case 63:
		//line grammar.y:596
		{
			yyVAL.stmts = []ast.Stmt{yyS[yypt-0].stmt}
		}
	case 66:
		//line grammar.y:604
		{
			yyVAL.stmts = nil
			yyVAL.stmts = append(yyVAL.stmts, yyS[yypt-0].stmt)
		}
	case 67:
		//line grammar.y:609
		{
			yyVAL.stmts = append(yyVAL.stmts, yyS[yypt-0].stmt)
		}
	case 68:
		//line grammar.y:615
		{
			yyVAL.stmts = yyS[yypt-2].stmts
		}
	case 69:
		//line grammar.y:621
		{
			yyVAL.stmt = yyS[yypt-0].stmt
		}
	case 70:
		//line grammar.y:625
		{
			yyVAL.stmt = yyS[yypt-0].stmt
		}
	case 71:
		//line grammar.y:629
		{
			yyVAL.stmt = yyS[yypt-0].stmt
		}
	case 72:
		//line grammar.y:633
		{
			yyVAL.stmt = yyS[yypt-0].stmt
		}
	case 73:
		//line grammar.y:637
		{
			yyVAL.stmt = yyS[yypt-0].stmt
		}
	case 74:
		//line grammar.y:641
		{
			yyVAL.stmt = yyS[yypt-0].stmt
		}
	case 75:
		//line grammar.y:645
		{
			yyVAL.stmt = yyS[yypt-0].stmt
		}
	case 76:
		//line grammar.y:649
		{
			yyVAL.stmt = yyS[yypt-0].stmt
		}
	case 77:
		//line grammar.y:676
		{
			target := yyS[yypt-2].expr
			setCtx(yylex, target, ast.Store)
			yyVAL.stmt = &ast.AugAssign{StmtBase: ast.StmtBase{Pos: yyVAL.pos}, Target: target, Op: yyS[yypt-1].op, Value: yyS[yypt-0].expr}
		}
	case 78:
		//line grammar.y:682
		{
			targets := []ast.Expr{yyS[yypt-1].expr}
			targets = append(targets, yyS[yypt-0].exprs...)
			value := targets[len(targets)-1]
			targets = targets[:len(targets)-1]
			setCtxs(yylex, targets, ast.Store)
			yyVAL.stmt = &ast.Assign{StmtBase: ast.StmtBase{Pos: yyVAL.pos}, Targets: targets, Value: value}
		}
	case 79:
		//line grammar.y:691
		{
			yyVAL.stmt = &ast.ExprStmt{StmtBase: ast.StmtBase{Pos: yyVAL.pos}, Value: yyS[yypt-0].expr}
		}
	case 80:
		//line grammar.y:697
		{
			yyVAL.expr = yyS[yypt-0].expr
		}
	case 81:
		//line grammar.y:701
		{
			yyVAL.expr = yyS[yypt-0].expr
		}
	case 82:
		//line grammar.y:707
		{
			yyVAL.expr = yyS[yypt-0].expr
		}
	case 83:
		//line grammar.y:711
		{
			yyVAL.expr = yyS[yypt-0].expr
		}
	case 84:
		//line grammar.y:717
		{
			yyVAL.exprs = nil
			yyVAL.exprs = append(yyVAL.exprs, yyS[yypt-0].expr)
		}
	case 85:
		//line grammar.y:722
		{
			yyVAL.exprs = append(yyVAL.exprs, yyS[yypt-0].expr)
		}
	case 86:
		//line grammar.y:728
		{
			yyVAL.exprs = nil
			yyVAL.exprs = append(yyVAL.exprs, yyS[yypt-0].expr)
		}
	case 87:
		//line grammar.y:733
		{
			yyVAL.exprs = append(yyVAL.exprs, yyS[yypt-0].expr)
		}
	case 88:
		//line grammar.y:739
		{
			yyVAL.expr = yyS[yypt-0].expr
		}
	case 89:
		//line grammar.y:743
		{
			yyVAL.expr = yyS[yypt-0].expr
		}
	case 90:
		//line grammar.y:748
		{
			yyVAL.comma = false
		}
	case 91:
		//line grammar.y:752
		{
			yyVAL.comma = true
		}
	case 92:
		//line grammar.y:758
		{
			yyVAL.expr = tupleOrExpr(yyVAL.pos, yyS[yypt-1].exprs, yyS[yypt-0].comma)
		}
	case 93:
		//line grammar.y:764
		{
			yyVAL.op = ast.Add
		}
	case 94:
		//line grammar.y:768
		{
			yyVAL.op = ast.Sub
		}
	case 95:
		//line grammar.y:772
		{
			yyVAL.op = ast.Mult
		}
	case 96:
		//line grammar.y:776
		{
			yyVAL.op = ast.Div
		}
	case 97:
		//line grammar.y:780
		{
			yyVAL.op = ast.Modulo
		}
	case 98:
		//line grammar.y:784
		{
			yyVAL.op = ast.BitAnd
		}
	case 99:
		//line grammar.y:788
		{
			yyVAL.op = ast.BitOr
		}
	case 100:
		//line grammar.y:792
		{
			yyVAL.op = ast.BitXor
		}
	case 101:
		//line grammar.y:796
		{
			yyVAL.op = ast.LShift
		}
	case 102:
		//line grammar.y:800
		{
			yyVAL.op = ast.RShift
		}
	case 103:
		//line grammar.y:804
		{
			yyVAL.op = ast.Pow
		}
	case 104:
		//line grammar.y:808
		{
			yyVAL.op = ast.FloorDiv
		}
	case 105:
		//line grammar.y:815
		{
			setCtxs(yylex, yyS[yypt-0].exprs, ast.Del)
			yyVAL.stmt = &ast.Delete{StmtBase: ast.StmtBase{Pos: yyVAL.pos}, Targets: yyS[yypt-0].exprs}
		}
	case 106:
		//line grammar.y:822
		{
			yyVAL.stmt = &ast.Pass{StmtBase: ast.StmtBase{Pos: yyVAL.pos}}
		}
	case 107:
		//line grammar.y:828
		{
			yyVAL.stmt = yyS[yypt-0].stmt
		}
	case 108:
		//line grammar.y:832
		{
			yyVAL.stmt = yyS[yypt-0].stmt
		}
	case 109:
		//line grammar.y:836
		{
			yyVAL.stmt = yyS[yypt-0].stmt
		}
	case 110:
		//line grammar.y:840
		{
			yyVAL.stmt = yyS[yypt-0].stmt
		}
	case 111:
		//line grammar.y:844
		{
			yyVAL.stmt = yyS[yypt-0].stmt
		}
	case 112:
		//line grammar.y:850
		{
			yyVAL.stmt = &ast.Break{StmtBase: ast.StmtBase{Pos: yyVAL.pos}}
		}
	case 113:
		//line grammar.y:856
		{
			yyVAL.stmt = &ast.Continue{StmtBase: ast.StmtBase{Pos: yyVAL.pos}}
		}
	case 114:
		//line grammar.y:862
		{
			yyVAL.stmt = &ast.Return{StmtBase: ast.StmtBase{Pos: yyVAL.pos}}
		}
	case 115:
		//line grammar.y:866
		{
			yyVAL.stmt = &ast.Return{StmtBase: ast.StmtBase{Pos: yyVAL.pos}, Value: yyS[yypt-0].expr}
		}
	case 116:
		//line grammar.y:872
		{
			yyVAL.stmt = &ast.ExprStmt{StmtBase: ast.StmtBase{Pos: yyVAL.pos}, Value: yyS[yypt-0].expr}
		}
	case 117:
		//line grammar.y:878
		{
			yyVAL.stmt = &ast.Raise{StmtBase: ast.StmtBase{Pos: yyVAL.pos}}
		}
	case 118:
		//line grammar.y:882
		{
			yyVAL.stmt = &ast.Raise{StmtBase: ast.StmtBase{Pos: yyVAL.pos}, Exc: yyS[yypt-0].expr}
		}
	case 119:
		//line grammar.y:886
		{
			yyVAL.stmt = &ast.Raise{StmtBase: ast.StmtBase{Pos: yyVAL.pos}, Exc: yyS[yypt-2].expr, Cause: yyS[yypt-0].expr}
		}
	case 120:
		//line grammar.y:892
		{
			yyVAL.stmt = yyS[yypt-0].stmt
		}
	case 121:
		//line grammar.y:896
		{
			yyVAL.stmt = yyS[yypt-0].stmt
		}
	case 122:
		//line grammar.y:902
		{
			yyVAL.stmt = &ast.Import{StmtBase: ast.StmtBase{Pos: yyVAL.pos}, Names: yyS[yypt-0].aliases}
		}
	case 123:
		//line grammar.y:909
		{
			yyVAL.level = 1
		}
	case 124:
		//line grammar.y:913
		{
			yyVAL.level = 3
		}
	case 125:
		//line grammar.y:919
		{
			yyVAL.level = yyS[yypt-0].level
		}
	case 126:
		//line grammar.y:923
		{
			yyVAL.level += yyS[yypt-0].level
		}
	case 127:
		//line grammar.y:929
		{
			yyVAL.level = 0
			yyVAL.str = yyS[yypt-0].str
		}
	case 128:
		//line grammar.y:934
		{
			yyVAL.level = yyS[yypt-1].level
			yyVAL.str = yyS[yypt-0].str
		}
	case 129:
		//line grammar.y:939
		{
			yyVAL.level = yyS[yypt-0].level
			yyVAL.str = ""
		}
	case 130:
		//line grammar.y:946
		{
			yyVAL.aliases = []*ast.Alias{&ast.Alias{Pos: yyVAL.pos, Name: ast.Identifier("*")}}
		}
	case 131:
		//line grammar.y:950
		{
			yyVAL.aliases = yyS[yypt-2].aliases
		}
	case 132:
		//line grammar.y:954
		{
			yyVAL.aliases = yyS[yypt-1].aliases
		}
	case 133:
		//line grammar.y:960
		{
			yyVAL.stmt = &ast.ImportFrom{StmtBase: ast.StmtBase{Pos: yyVAL.pos}, Module: ast.Identifier(yyS[yypt-2].str), Names: yyS[yypt-0].aliases, Level: yyS[yypt-2].level}
		}
	case 134:
		//line grammar.y:966
		{
			yyVAL.alias = &ast.Alias{Pos: yyVAL.pos, Name: ast.Identifier(yyS[yypt-0].str)}
		}
	case 135:
		//line grammar.y:970
		{
			yyVAL.alias = &ast.Alias{Pos: yyVAL.pos, Name: ast.Identifier(yyS[yypt-2].str), AsName: ast.Identifier(yyS[yypt-0].str)}
		}
	case 136:
		//line grammar.y:976
		{
			yyVAL.alias = &ast.Alias{Pos: yyVAL.pos, Name: ast.Identifier(yyS[yypt-0].str)}
		}
	case 137:
		//line grammar.y:980
		{
			yyVAL.alias = &ast.Alias{Pos: yyVAL.pos, Name: ast.Identifier(yyS[yypt-2].str), AsName: ast.Identifier(yyS[yypt-0].str)}
		}
	case 138:
		//line grammar.y:986
		{
			yyVAL.aliases = nil
			yyVAL.aliases = append(yyVAL.aliases, yyS[yypt-0].alias)
		}
	case 139:
		//line grammar.y:991
		{
			yyVAL.aliases = append(yyVAL.aliases, yyS[yypt-0].alias)
		}
	case 140:
		//line grammar.y:997
		{
			yyVAL.aliases = nil
			yyVAL.aliases = append(yyVAL.aliases, yyS[yypt-0].alias)
		}
	case 141:
		//line grammar.y:1002
		{
			yyVAL.aliases = append(yyVAL.aliases, yyS[yypt-0].alias)
		}
	case 142:
		//line grammar.y:1008
		{
			yyVAL.str = yyS[yypt-0].str
		}
	case 143:
		//line grammar.y:1012
		{
			yyVAL.str += "." + yyS[yypt-0].str
		}
	case 144:
		//line grammar.y:1018
		{
			yyVAL.identifiers = nil
			yyVAL.identifiers = append(yyVAL.identifiers, ast.Identifier(yyS[yypt-0].str))
		}
	case 145:
		//line grammar.y:1023
		{
			yyVAL.identifiers = append(yyVAL.identifiers, ast.Identifier(yyS[yypt-0].str))
		}
	case 146:
		//line grammar.y:1029
		{
			yyVAL.stmt = &ast.Global{StmtBase: ast.StmtBase{Pos: yyVAL.pos}, Names: yyS[yypt-0].identifiers}
		}
	case 147:
		//line grammar.y:1035
		{
			yyVAL.stmt = &ast.Nonlocal{StmtBase: ast.StmtBase{Pos: yyVAL.pos}, Names: yyS[yypt-0].identifiers}
		}
	case 148:
		//line grammar.y:1041
		{
			yyVAL.exprs = nil
			yyVAL.exprs = append(yyVAL.exprs, yyS[yypt-0].expr)
		}
	case 149:
		//line grammar.y:1046
		{
			yyVAL.exprs = append(yyVAL.exprs, yyS[yypt-0].expr)
		}
	case 150:
		//line grammar.y:1052
		{
			yyVAL.stmt = &ast.Assert{StmtBase: ast.StmtBase{Pos: yyVAL.pos}, Test: yyS[yypt-0].expr}
		}
	case 151:
		//line grammar.y:1056
		{
			yyVAL.stmt = &ast.Assert{StmtBase: ast.StmtBase{Pos: yyVAL.pos}, Test: yyS[yypt-2].expr, Msg: yyS[yypt-0].expr}
		}
	case 152:
		//line grammar.y:1062
		{
			yyVAL.stmt = yyS[yypt-0].stmt
		}
	case 153:
		//line grammar.y:1066
		{
			yyVAL.stmt = yyS[yypt-0].stmt
		}
	case 154:
		//line grammar.y:1070
		{
			yyVAL.stmt = yyS[yypt-0].stmt
		}
	case 155:
		//line grammar.y:1074
		{
			yyVAL.stmt = yyS[yypt-0].stmt
		}
	case 156:
		//line grammar.y:1078
		{
			yyVAL.stmt = yyS[yypt-0].stmt
		}
	case 157:
		//line grammar.y:1082
		{
			yyVAL.stmt = yyS[yypt-0].stmt
		}
	case 158:
		//line grammar.y:1086
		{
			yyVAL.stmt = yyS[yypt-0].stmt
		}
	case 159:
		//line grammar.y:1090
		{
			yyVAL.stmt = yyS[yypt-0].stmt
		}
	case 160:
		//line grammar.y:1095
		{
			yyVAL.ifstmt = nil
			yyVAL.lastif = nil
		}
	case 161:
		//line grammar.y:1100
		{
			elifs := yyVAL.ifstmt
			newif := &ast.If{StmtBase: ast.StmtBase{Pos: yyVAL.pos}, Test: yyS[yypt-2].expr, Body: yyS[yypt-0].stmts}
			if elifs == nil {
				yyVAL.ifstmt = newif
			} else {
				yyVAL.lastif.Orelse = []ast.Stmt{newif}
			}
			yyVAL.lastif = newif
		}
	case 162:
		//line grammar.y:1112
		{
			yyVAL.stmts = nil
		}
	case 163:
		//line grammar.y:1116
		{
			yyVAL.stmts = yyS[yypt-0].stmts
		}
	case 164:
		//line grammar.y:1122
		{
			newif := &ast.If{StmtBase: ast.StmtBase{Pos: yyVAL.pos}, Test: yyS[yypt-4].expr, Body: yyS[yypt-2].stmts}
			yyVAL.stmt = newif
			elifs := yyS[yypt-1].ifstmt
			optional_else := yyS[yypt-0].stmts
			if len(optional_else) != 0 {
				if elifs != nil {
					yyS[yypt-1].lastif.Orelse = optional_else
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
	case 165:
		//line grammar.y:1143
		{
			yyVAL.stmt = &ast.While{StmtBase: ast.StmtBase{Pos: yyVAL.pos}, Test: yyS[yypt-3].expr, Body: yyS[yypt-1].stmts, Orelse: yyS[yypt-0].stmts}
		}
	case 166:
		//line grammar.y:1149
		{
			target := tupleOrExpr(yyVAL.pos, yyS[yypt-5].exprs, false)
			setCtx(yylex, target, ast.Store)
			yyVAL.stmt = &ast.For{StmtBase: ast.StmtBase{Pos: yyVAL.pos}, Target: target, Iter: yyS[yypt-3].expr, Body: yyS[yypt-1].stmts, Orelse: yyS[yypt-0].stmts}
		}
	case 167:
		//line grammar.y:1156
		{
			yyVAL.exchandlers = nil
		}
	case 168:
		//line grammar.y:1160
		{
			exc := &ast.ExceptHandler{Pos: yyVAL.pos, ExprType: yyS[yypt-2].expr, Name: ast.Identifier(yyS[yypt-2].str), Body: yyS[yypt-0].stmts}
			yyVAL.exchandlers = append(yyVAL.exchandlers, exc)
		}
	case 169:
		//line grammar.y:1167
		{
			yyVAL.stmt = &ast.Try{StmtBase: ast.StmtBase{Pos: yyVAL.pos}, Body: yyS[yypt-1].stmts, Handlers: yyS[yypt-0].exchandlers}
		}
	case 170:
		//line grammar.y:1171
		{
			yyVAL.stmt = &ast.Try{StmtBase: ast.StmtBase{Pos: yyVAL.pos}, Body: yyS[yypt-4].stmts, Handlers: yyS[yypt-3].exchandlers, Orelse: yyS[yypt-0].stmts}
		}
	case 171:
		//line grammar.y:1175
		{
			yyVAL.stmt = &ast.Try{StmtBase: ast.StmtBase{Pos: yyVAL.pos}, Body: yyS[yypt-4].stmts, Handlers: yyS[yypt-3].exchandlers, Finalbody: yyS[yypt-0].stmts}
		}
	case 172:
		//line grammar.y:1179
		{
			yyVAL.stmt = &ast.Try{StmtBase: ast.StmtBase{Pos: yyVAL.pos}, Body: yyS[yypt-7].stmts, Handlers: yyS[yypt-6].exchandlers, Orelse: yyS[yypt-3].stmts, Finalbody: yyS[yypt-0].stmts}
		}
	case 173:
		//line grammar.y:1185
		{
			yyVAL.withitems = nil
			yyVAL.withitems = append(yyVAL.withitems, yyS[yypt-0].withitem)
		}
	case 174:
		//line grammar.y:1190
		{
			yyVAL.withitems = append(yyVAL.withitems, yyS[yypt-0].withitem)
		}
	case 175:
		//line grammar.y:1196
		{
			yyVAL.stmt = &ast.With{StmtBase: ast.StmtBase{Pos: yyVAL.pos}, Items: yyS[yypt-2].withitems, Body: yyS[yypt-0].stmts}
		}
	case 176:
		//line grammar.y:1202
		{
			yyVAL.withitem = &ast.WithItem{Pos: yyVAL.pos, ContextExpr: yyS[yypt-0].expr}
		}
	case 177:
		//line grammar.y:1206
		{
			v := yyS[yypt-0].expr
			setCtx(yylex, v, ast.Store)
			yyVAL.withitem = &ast.WithItem{Pos: yyVAL.pos, ContextExpr: yyS[yypt-2].expr, OptionalVars: v}
		}
	case 178:
		//line grammar.y:1215
		{
			yyVAL.expr = nil
			yyVAL.str = ""
		}
	case 179:
		//line grammar.y:1220
		{
			yyVAL.expr = yyS[yypt-0].expr
			yyVAL.str = ""
		}
	case 180:
		//line grammar.y:1225
		{
			yyVAL.expr = yyS[yypt-2].expr
			yyVAL.str = yyS[yypt-0].str
		}
	case 181:
		//line grammar.y:1232
		{
			yyVAL.stmts = nil
			yyVAL.stmts = append(yyVAL.stmts, yyS[yypt-0].stmts...)
		}
	case 182:
		//line grammar.y:1237
		{
			yyVAL.stmts = append(yyVAL.stmts, yyS[yypt-0].stmts...)
		}
	case 183:
		//line grammar.y:1243
		{
			yyVAL.stmts = yyS[yypt-0].stmts
		}
	case 184:
		//line grammar.y:1247
		{
			yyVAL.stmts = yyS[yypt-1].stmts
		}
	case 185:
		//line grammar.y:1253
		{
			yyVAL.expr = yyS[yypt-0].expr
		}
	case 186:
		//line grammar.y:1257
		{
			yyVAL.expr = &ast.IfExp{ExprBase: ast.ExprBase{Pos: yyVAL.pos}, Test: yyS[yypt-2].expr, Body: yyS[yypt-4].expr, Orelse: yyS[yypt-0].expr}
		}
	case 187:
		//line grammar.y:1261
		{
			yyVAL.expr = yyS[yypt-0].expr
		}
	case 188:
		//line grammar.y:1267
		{
			yyVAL.expr = yyS[yypt-0].expr
		}
	case 189:
		//line grammar.y:1271
		{
			yyVAL.expr = yyS[yypt-0].expr
		}
	case 190:
		//line grammar.y:1277
		{
			args := &ast.Arguments{Pos: yyVAL.pos}
			yyVAL.expr = &ast.Lambda{ExprBase: ast.ExprBase{Pos: yyVAL.pos}, Args: args, Body: yyS[yypt-0].expr}
		}
	case 191:
		//line grammar.y:1282
		{
			yyVAL.expr = &ast.Lambda{ExprBase: ast.ExprBase{Pos: yyVAL.pos}, Args: yyS[yypt-2].arguments, Body: yyS[yypt-0].expr}
		}
	case 192:
		//line grammar.y:1288
		{
			args := &ast.Arguments{Pos: yyVAL.pos}
			yyVAL.expr = &ast.Lambda{ExprBase: ast.ExprBase{Pos: yyVAL.pos}, Args: args, Body: yyS[yypt-0].expr}
		}
	case 193:
		//line grammar.y:1293
		{
			yyVAL.expr = &ast.Lambda{ExprBase: ast.ExprBase{Pos: yyVAL.pos}, Args: yyS[yypt-2].arguments, Body: yyS[yypt-0].expr}
		}
	case 194:
		//line grammar.y:1299
		{
			yyVAL.expr = yyS[yypt-0].expr
			yyVAL.isExpr = true
		}
	case 195:
		//line grammar.y:1304
		{
			if !yyS[yypt-2].isExpr {
				boolop := yyVAL.expr.(*ast.BoolOp)
				boolop.Values = append(boolop.Values, yyS[yypt-0].expr)
			} else {
				yyVAL.expr = &ast.BoolOp{ExprBase: ast.ExprBase{Pos: yyVAL.pos}, Op: ast.Or, Values: []ast.Expr{yyVAL.expr, yyS[yypt-0].expr}}
			}
			yyVAL.isExpr = false
		}
	case 196:
		//line grammar.y:1316
		{
			yyVAL.expr = yyS[yypt-0].expr
			yyVAL.isExpr = true
		}
	case 197:
		//line grammar.y:1321
		{
			if !yyS[yypt-2].isExpr {
				boolop := yyVAL.expr.(*ast.BoolOp)
				boolop.Values = append(boolop.Values, yyS[yypt-0].expr)
			} else {
				yyVAL.expr = &ast.BoolOp{ExprBase: ast.ExprBase{Pos: yyVAL.pos}, Op: ast.And, Values: []ast.Expr{yyVAL.expr, yyS[yypt-0].expr}}
			}
			yyVAL.isExpr = false
		}
	case 198:
		//line grammar.y:1333
		{
			yyVAL.expr = &ast.UnaryOp{ExprBase: ast.ExprBase{Pos: yyVAL.pos}, Op: ast.Not, Operand: yyS[yypt-0].expr}
		}
	case 199:
		//line grammar.y:1337
		{
			yyVAL.expr = yyS[yypt-0].expr
		}
	case 200:
		//line grammar.y:1343
		{
			yyVAL.expr = yyS[yypt-0].expr
			yyVAL.isExpr = true
		}
	case 201:
		//line grammar.y:1348
		{
			if !yyS[yypt-2].isExpr {
				comp := yyVAL.expr.(*ast.Compare)
				comp.Ops = append(comp.Ops, yyS[yypt-1].cmpop)
				comp.Comparators = append(comp.Comparators, yyS[yypt-0].expr)
			} else {
				yyVAL.expr = &ast.Compare{ExprBase: ast.ExprBase{Pos: yyVAL.pos}, Left: yyVAL.expr, Ops: []ast.CmpOp{yyS[yypt-1].cmpop}, Comparators: []ast.Expr{yyS[yypt-0].expr}}
			}
			yyVAL.isExpr = false
		}
	case 202:
		//line grammar.y:1363
		{
			yyVAL.cmpop = ast.Lt
		}
	case 203:
		//line grammar.y:1367
		{
			yyVAL.cmpop = ast.Gt
		}
	case 204:
		//line grammar.y:1371
		{
			yyVAL.cmpop = ast.Eq
		}
	case 205:
		//line grammar.y:1375
		{
			yyVAL.cmpop = ast.GtE
		}
	case 206:
		//line grammar.y:1379
		{
			yyVAL.cmpop = ast.LtE
		}
	case 207:
		//line grammar.y:1383
		{
			yylex.(*yyLex).SyntaxError("invalid syntax")
		}
	case 208:
		//line grammar.y:1387
		{
			yyVAL.cmpop = ast.NotEq
		}
	case 209:
		//line grammar.y:1391
		{
			yyVAL.cmpop = ast.In
		}
	case 210:
		//line grammar.y:1395
		{
			yyVAL.cmpop = ast.NotIn
		}
	case 211:
		//line grammar.y:1399
		{
			yyVAL.cmpop = ast.Is
		}
	case 212:
		//line grammar.y:1403
		{
			yyVAL.cmpop = ast.IsNot
		}
	case 213:
		//line grammar.y:1409
		{
			yyVAL.expr = &ast.Starred{ExprBase: ast.ExprBase{Pos: yyVAL.pos}, Value: yyS[yypt-0].expr, Ctx: ast.Load}
		}
	case 214:
		//line grammar.y:1415
		{
			yyVAL.expr = yyS[yypt-0].expr
		}
	case 215:
		//line grammar.y:1419
		{
			yyVAL.expr = &ast.BinOp{ExprBase: ast.ExprBase{Pos: yyVAL.pos}, Left: yyS[yypt-2].expr, Op: ast.BitOr, Right: yyS[yypt-0].expr}
		}
	case 216:
		//line grammar.y:1425
		{
			yyVAL.expr = yyS[yypt-0].expr
		}
	case 217:
		//line grammar.y:1429
		{
			yyVAL.expr = &ast.BinOp{ExprBase: ast.ExprBase{Pos: yyVAL.pos}, Left: yyS[yypt-2].expr, Op: ast.BitXor, Right: yyS[yypt-0].expr}
		}
	case 218:
		//line grammar.y:1435
		{
			yyVAL.expr = yyS[yypt-0].expr
		}
	case 219:
		//line grammar.y:1439
		{
			yyVAL.expr = &ast.BinOp{ExprBase: ast.ExprBase{Pos: yyVAL.pos}, Left: yyS[yypt-2].expr, Op: ast.BitAnd, Right: yyS[yypt-0].expr}
		}
	case 220:
		//line grammar.y:1445
		{
			yyVAL.expr = yyS[yypt-0].expr
		}
	case 221:
		//line grammar.y:1449
		{
			yyVAL.expr = &ast.BinOp{ExprBase: ast.ExprBase{Pos: yyVAL.pos}, Left: yyS[yypt-2].expr, Op: ast.LShift, Right: yyS[yypt-0].expr}
		}
	case 222:
		//line grammar.y:1453
		{
			yyVAL.expr = &ast.BinOp{ExprBase: ast.ExprBase{Pos: yyVAL.pos}, Left: yyS[yypt-2].expr, Op: ast.RShift, Right: yyS[yypt-0].expr}
		}
	case 223:
		//line grammar.y:1459
		{
			yyVAL.expr = yyS[yypt-0].expr
		}
	case 224:
		//line grammar.y:1463
		{
			yyVAL.expr = &ast.BinOp{ExprBase: ast.ExprBase{Pos: yyVAL.pos}, Left: yyS[yypt-2].expr, Op: ast.Add, Right: yyS[yypt-0].expr}
		}
	case 225:
		//line grammar.y:1467
		{
			yyVAL.expr = &ast.BinOp{ExprBase: ast.ExprBase{Pos: yyVAL.pos}, Left: yyS[yypt-2].expr, Op: ast.Sub, Right: yyS[yypt-0].expr}
		}
	case 226:
		//line grammar.y:1473
		{
			yyVAL.expr = yyS[yypt-0].expr
		}
	case 227:
		//line grammar.y:1477
		{
			yyVAL.expr = &ast.BinOp{ExprBase: ast.ExprBase{Pos: yyVAL.pos}, Left: yyS[yypt-2].expr, Op: ast.Mult, Right: yyS[yypt-0].expr}
		}
	case 228:
		//line grammar.y:1481
		{
			yyVAL.expr = &ast.BinOp{ExprBase: ast.ExprBase{Pos: yyVAL.pos}, Left: yyS[yypt-2].expr, Op: ast.Div, Right: yyS[yypt-0].expr}
		}
	case 229:
		//line grammar.y:1485
		{
			yyVAL.expr = &ast.BinOp{ExprBase: ast.ExprBase{Pos: yyVAL.pos}, Left: yyS[yypt-2].expr, Op: ast.Modulo, Right: yyS[yypt-0].expr}
		}
	case 230:
		//line grammar.y:1489
		{
			yyVAL.expr = &ast.BinOp{ExprBase: ast.ExprBase{Pos: yyVAL.pos}, Left: yyS[yypt-2].expr, Op: ast.FloorDiv, Right: yyS[yypt-0].expr}
		}
	case 231:
		//line grammar.y:1495
		{
			yyVAL.expr = &ast.UnaryOp{ExprBase: ast.ExprBase{Pos: yyVAL.pos}, Op: ast.UAdd, Operand: yyS[yypt-0].expr}
		}
	case 232:
		//line grammar.y:1499
		{
			yyVAL.expr = &ast.UnaryOp{ExprBase: ast.ExprBase{Pos: yyVAL.pos}, Op: ast.USub, Operand: yyS[yypt-0].expr}
		}
	case 233:
		//line grammar.y:1503
		{
			yyVAL.expr = &ast.UnaryOp{ExprBase: ast.ExprBase{Pos: yyVAL.pos}, Op: ast.Invert, Operand: yyS[yypt-0].expr}
		}
	case 234:
		//line grammar.y:1507
		{
			yyVAL.expr = yyS[yypt-0].expr
		}
	case 235:
		//line grammar.y:1513
		{
			yyVAL.expr = applyTrailers(yyS[yypt-1].expr, yyS[yypt-0].exprs)
		}
	case 236:
		//line grammar.y:1517
		{
			yyVAL.expr = &ast.BinOp{ExprBase: ast.ExprBase{Pos: yyVAL.pos}, Left: applyTrailers(yyS[yypt-3].expr, yyS[yypt-2].exprs), Op: ast.Pow, Right: yyS[yypt-0].expr}
		}
	case 237:
		//line grammar.y:1523
		{
			yyVAL.exprs = nil
		}
	case 238:
		//line grammar.y:1527
		{
			yyVAL.exprs = append(yyVAL.exprs, yyS[yypt-0].expr)
		}
	case 239:
		//line grammar.y:1533
		{
			yyVAL.obj = yyS[yypt-0].obj
		}
	case 240:
		//line grammar.y:1537
		{
			switch a := yyVAL.obj.(type) {
			case py.String:
				switch b := yyS[yypt-0].obj.(type) {
				case py.String:
					yyVAL.obj = a + b
				default:
					yylex.(*yyLex).SyntaxError("cannot mix string and nonstring literals")
				}
			case py.Bytes:
				switch b := yyS[yypt-0].obj.(type) {
				case py.Bytes:
					yyVAL.obj = append(a, b...)
				default:
					yylex.(*yyLex).SyntaxError("cannot mix bytes and nonbytes literals")
				}
			}
		}
	case 241:
		//line grammar.y:1558
		{
			yyVAL.expr = &ast.Tuple{ExprBase: ast.ExprBase{Pos: yyVAL.pos}, Ctx: ast.Load}
		}
	case 242:
		//line grammar.y:1562
		{
			yyVAL.expr = yyS[yypt-1].expr
		}
	case 243:
		//line grammar.y:1566
		{
			yyVAL.expr = &ast.GeneratorExp{ExprBase: ast.ExprBase{Pos: yyVAL.pos}, Elt: yyS[yypt-2].expr, Generators: yyS[yypt-1].comprehensions}
		}
	case 244:
		//line grammar.y:1570
		{
			yyVAL.expr = tupleOrExpr(yyVAL.pos, yyS[yypt-2].exprs, yyS[yypt-1].comma)
		}
	case 245:
		//line grammar.y:1574
		{
			yyVAL.expr = &ast.List{ExprBase: ast.ExprBase{Pos: yyVAL.pos}, Ctx: ast.Load}
		}
	case 246:
		//line grammar.y:1578
		{
			yyVAL.expr = &ast.ListComp{ExprBase: ast.ExprBase{Pos: yyVAL.pos}, Elt: yyS[yypt-2].expr, Generators: yyS[yypt-1].comprehensions}
		}
	case 247:
		//line grammar.y:1582
		{
			yyVAL.expr = &ast.List{ExprBase: ast.ExprBase{Pos: yyVAL.pos}, Elts: yyS[yypt-2].exprs, Ctx: ast.Load}
		}
	case 248:
		//line grammar.y:1586
		{
			yyVAL.expr = &ast.Dict{ExprBase: ast.ExprBase{Pos: yyVAL.pos}}
		}
	case 249:
		//line grammar.y:1590
		{
			yyVAL.expr = yyS[yypt-1].expr
		}
	case 250:
		//line grammar.y:1594
		{
			yyVAL.expr = &ast.Name{ExprBase: ast.ExprBase{Pos: yyVAL.pos}, Id: ast.Identifier(yyS[yypt-0].str), Ctx: ast.Load}
		}
	case 251:
		//line grammar.y:1598
		{
			yyVAL.expr = &ast.Num{ExprBase: ast.ExprBase{Pos: yyVAL.pos}, N: yyS[yypt-0].obj}
		}
	case 252:
		//line grammar.y:1602
		{
			switch s := yyS[yypt-0].obj.(type) {
			case py.String:
				yyVAL.expr = &ast.Str{ExprBase: ast.ExprBase{Pos: yyVAL.pos}, S: s}
			case py.Bytes:
				yyVAL.expr = &ast.Bytes{ExprBase: ast.ExprBase{Pos: yyVAL.pos}, S: s}
			default:
				panic("not Bytes or String in strings")
			}
		}
	case 253:
		//line grammar.y:1613
		{
			yyVAL.expr = &ast.Ellipsis{ExprBase: ast.ExprBase{Pos: yyVAL.pos}}
		}
	case 254:
		//line grammar.y:1617
		{
			yyVAL.expr = &ast.NameConstant{ExprBase: ast.ExprBase{Pos: yyVAL.pos}, Value: py.None}
		}
	case 255:
		//line grammar.y:1621
		{
			yyVAL.expr = &ast.NameConstant{ExprBase: ast.ExprBase{Pos: yyVAL.pos}, Value: py.True}
		}
	case 256:
		//line grammar.y:1625
		{
			yyVAL.expr = &ast.NameConstant{ExprBase: ast.ExprBase{Pos: yyVAL.pos}, Value: py.False}
		}
	case 257:
		//line grammar.y:1632
		{
			yyVAL.expr = &ast.Call{ExprBase: ast.ExprBase{Pos: yyVAL.pos}}
		}
	case 258:
		//line grammar.y:1636
		{
			yyVAL.expr = yyS[yypt-1].call
		}
	case 259:
		//line grammar.y:1640
		{
			slice := yyS[yypt-1].slice
			// If all items of a ExtSlice are just Index then return as tuple
			if extslice, ok := slice.(*ast.ExtSlice); ok {
				elts := make([]ast.Expr, len(extslice.Dims))
				for i, item := range extslice.Dims {
					if index, isIndex := item.(*ast.Index); isIndex {
						elts[i] = index.Value
					} else {
						goto notAllIndex
					}
				}
				slice = &ast.Index{SliceBase: extslice.SliceBase, Value: &ast.Tuple{ExprBase: ast.ExprBase{Pos: extslice.SliceBase.Pos}, Elts: elts, Ctx: ast.Load}}
			notAllIndex:
			}
			yyVAL.expr = &ast.Subscript{ExprBase: ast.ExprBase{Pos: yyVAL.pos}, Slice: slice, Ctx: ast.Load}
		}
	case 260:
		//line grammar.y:1658
		{
			yyVAL.expr = &ast.Attribute{ExprBase: ast.ExprBase{Pos: yyVAL.pos}, Attr: ast.Identifier(yyS[yypt-0].str), Ctx: ast.Load}
		}
	case 261:
		//line grammar.y:1664
		{
			yyVAL.slice = yyS[yypt-0].slice
			yyVAL.isExpr = true
		}
	case 262:
		//line grammar.y:1669
		{
			if !yyS[yypt-2].isExpr {
				extSlice := yyVAL.slice.(*ast.ExtSlice)
				extSlice.Dims = append(extSlice.Dims, yyS[yypt-0].slice)
			} else {
				yyVAL.slice = &ast.ExtSlice{SliceBase: ast.SliceBase{Pos: yyVAL.pos}, Dims: []ast.Slicer{yyS[yypt-2].slice, yyS[yypt-0].slice}}
			}
			yyVAL.isExpr = false
		}
	case 263:
		//line grammar.y:1681
		{
			if yyS[yypt-0].comma && yyS[yypt-1].isExpr {
				yyVAL.slice = &ast.ExtSlice{SliceBase: ast.SliceBase{Pos: yyVAL.pos}, Dims: []ast.Slicer{yyS[yypt-1].slice}}
			} else {
				yyVAL.slice = yyS[yypt-1].slice
			}
		}
	case 264:
		//line grammar.y:1691
		{
			yyVAL.slice = &ast.Index{SliceBase: ast.SliceBase{Pos: yyVAL.pos}, Value: yyS[yypt-0].expr}
		}
	case 265:
		//line grammar.y:1695
		{
			yyVAL.slice = &ast.Slice{SliceBase: ast.SliceBase{Pos: yyVAL.pos}, Lower: nil, Upper: nil, Step: nil}
		}
	case 266:
		//line grammar.y:1699
		{
			yyVAL.slice = &ast.Slice{SliceBase: ast.SliceBase{Pos: yyVAL.pos}, Lower: nil, Upper: nil, Step: yyS[yypt-0].expr}
		}
	case 267:
		//line grammar.y:1703
		{
			yyVAL.slice = &ast.Slice{SliceBase: ast.SliceBase{Pos: yyVAL.pos}, Lower: nil, Upper: yyS[yypt-0].expr, Step: nil}
		}
	case 268:
		//line grammar.y:1707
		{
			yyVAL.slice = &ast.Slice{SliceBase: ast.SliceBase{Pos: yyVAL.pos}, Lower: nil, Upper: yyS[yypt-1].expr, Step: yyS[yypt-0].expr}
		}
	case 269:
		//line grammar.y:1711
		{
			yyVAL.slice = &ast.Slice{SliceBase: ast.SliceBase{Pos: yyVAL.pos}, Lower: yyS[yypt-1].expr, Upper: nil, Step: nil}
		}
	case 270:
		//line grammar.y:1715
		{
			yyVAL.slice = &ast.Slice{SliceBase: ast.SliceBase{Pos: yyVAL.pos}, Lower: yyS[yypt-2].expr, Upper: nil, Step: yyS[yypt-0].expr}
		}
	case 271:
		//line grammar.y:1719
		{
			yyVAL.slice = &ast.Slice{SliceBase: ast.SliceBase{Pos: yyVAL.pos}, Lower: yyS[yypt-2].expr, Upper: yyS[yypt-0].expr, Step: nil}
		}
	case 272:
		//line grammar.y:1723
		{
			yyVAL.slice = &ast.Slice{SliceBase: ast.SliceBase{Pos: yyVAL.pos}, Lower: yyS[yypt-3].expr, Upper: yyS[yypt-1].expr, Step: yyS[yypt-0].expr}
		}
	case 273:
		//line grammar.y:1729
		{
			yyVAL.expr = nil
		}
	case 274:
		//line grammar.y:1733
		{
			yyVAL.expr = yyS[yypt-0].expr
		}
	case 275:
		//line grammar.y:1739
		{
			yyVAL.expr = yyS[yypt-0].expr
		}
	case 276:
		//line grammar.y:1743
		{
			yyVAL.expr = yyS[yypt-0].expr
		}
	case 277:
		//line grammar.y:1749
		{
			yyVAL.exprs = nil
			yyVAL.exprs = append(yyVAL.exprs, yyS[yypt-0].expr)
		}
	case 278:
		//line grammar.y:1754
		{
			yyVAL.exprs = append(yyVAL.exprs, yyS[yypt-0].expr)
		}
	case 279:
		//line grammar.y:1760
		{
			yyVAL.exprs = yyS[yypt-1].exprs
			yyVAL.comma = yyS[yypt-0].comma
		}
	case 280:
		//line grammar.y:1767
		{
			elts := yyS[yypt-1].exprs
			if yyS[yypt-0].comma || len(elts) > 1 {
				yyVAL.expr = &ast.Tuple{ExprBase: ast.ExprBase{Pos: yyVAL.pos}, Elts: elts, Ctx: ast.Load}
			} else {
				yyVAL.expr = elts[0]
			}
		}
	case 281:
		//line grammar.y:1778
		{
			yyVAL.exprs = yyS[yypt-1].exprs
		}
	case 282:
		//line grammar.y:1785
		{
			yyVAL.exprs = nil
			yyVAL.exprs = append(yyVAL.exprs, yyS[yypt-2].expr, yyS[yypt-0].expr) // key, value order
		}
	case 283:
		//line grammar.y:1790
		{
			yyVAL.exprs = append(yyVAL.exprs, yyS[yypt-2].expr, yyS[yypt-0].expr)
		}
	case 284:
		//line grammar.y:1796
		{
			keyValues := yyS[yypt-1].exprs
			d := &ast.Dict{ExprBase: ast.ExprBase{Pos: yyVAL.pos}, Keys: nil, Values: nil}
			for i := 0; i < len(keyValues)-1; i += 2 {
				d.Keys = append(d.Keys, keyValues[i])
				d.Values = append(d.Values, keyValues[i+1])
			}
			yyVAL.expr = d
		}
	case 285:
		//line grammar.y:1806
		{
			yyVAL.expr = &ast.DictComp{ExprBase: ast.ExprBase{Pos: yyVAL.pos}, Key: yyS[yypt-3].expr, Value: yyS[yypt-1].expr, Generators: yyS[yypt-0].comprehensions}
		}
	case 286:
		//line grammar.y:1810
		{
			yyVAL.expr = &ast.Set{ExprBase: ast.ExprBase{Pos: yyVAL.pos}, Elts: yyS[yypt-0].exprs}
		}
	case 287:
		//line grammar.y:1814
		{
			yyVAL.expr = &ast.SetComp{ExprBase: ast.ExprBase{Pos: yyVAL.pos}, Elt: yyS[yypt-1].expr, Generators: yyS[yypt-0].comprehensions}
		}
	case 288:
		//line grammar.y:1820
		{
			classDef := &ast.ClassDef{StmtBase: ast.StmtBase{Pos: yyVAL.pos}, Name: ast.Identifier(yyS[yypt-3].str), Body: yyS[yypt-0].stmts}
			yyVAL.stmt = classDef
			args := yyS[yypt-2].call
			if args != nil {
				classDef.Bases = args.Args
				classDef.Keywords = args.Keywords
				classDef.Starargs = args.Starargs
				classDef.Kwargs = args.Kwargs
			}
		}
	case 289:
		//line grammar.y:1834
		{
			yyVAL.call = yyS[yypt-0].call
		}
	case 290:
		//line grammar.y:1838
		{
			yyVAL.call.Args = append(yyVAL.call.Args, yyS[yypt-0].call.Args...)
			yyVAL.call.Keywords = append(yyVAL.call.Keywords, yyS[yypt-0].call.Keywords...)
		}
	case 291:
		//line grammar.y:1844
		{
			yyVAL.call = &ast.Call{}
		}
	case 292:
		//line grammar.y:1848
		{
			yyVAL.call = yyS[yypt-1].call
		}
	case 293:
		//line grammar.y:1853
		{
			yyVAL.call = &ast.Call{}
		}
	case 294:
		//line grammar.y:1857
		{
			yyVAL.call.Args = append(yyVAL.call.Args, yyS[yypt-0].call.Args...)
			yyVAL.call.Keywords = append(yyVAL.call.Keywords, yyS[yypt-0].call.Keywords...)
		}
	case 295:
		//line grammar.y:1864
		{
			yyVAL.call = yyS[yypt-1].call
		}
	case 296:
		//line grammar.y:1868
		{
			call := yyS[yypt-3].call
			call.Starargs = yyS[yypt-1].expr
			if len(yyS[yypt-0].call.Args) != 0 {
				yylex.(*yyLex).SyntaxError("only named arguments may follow *expression")
			}
			call.Keywords = append(call.Keywords, yyS[yypt-0].call.Keywords...)
			yyVAL.call = call
		}
	case 297:
		//line grammar.y:1878
		{
			call := yyS[yypt-6].call
			call.Starargs = yyS[yypt-4].expr
			call.Kwargs = yyS[yypt-0].expr
			if len(yyS[yypt-3].call.Args) != 0 {
				yylex.(*yyLex).SyntaxError("only named arguments may follow *expression")
			}
			call.Keywords = append(call.Keywords, yyS[yypt-3].call.Keywords...)
			yyVAL.call = call
		}
	case 298:
		//line grammar.y:1889
		{
			call := yyS[yypt-2].call
			call.Kwargs = yyS[yypt-0].expr
			yyVAL.call = call
		}
	case 299:
		//line grammar.y:1899
		{
			yyVAL.call = &ast.Call{}
			yyVAL.call.Args = []ast.Expr{yyS[yypt-0].expr}
		}
	case 300:
		//line grammar.y:1904
		{
			yyVAL.call = &ast.Call{}
			yyVAL.call.Args = []ast.Expr{
				&ast.GeneratorExp{ExprBase: ast.ExprBase{Pos: yyVAL.pos}, Elt: yyS[yypt-1].expr, Generators: yyS[yypt-0].comprehensions},
			}
		}
	case 301:
		//line grammar.y:1911
		{
			yyVAL.call = &ast.Call{}
			test := yyS[yypt-2].expr
			if name, ok := test.(*ast.Name); ok {
				yyVAL.call.Keywords = []*ast.Keyword{&ast.Keyword{Pos: name.Pos, Arg: name.Id, Value: yyS[yypt-0].expr}}
			} else {
				yylex.(*yyLex).SyntaxError("keyword can't be an expression")
			}
		}
	case 302:
		//line grammar.y:1923
		{
			yyVAL.comprehensions = yyS[yypt-0].comprehensions
			yyVAL.exprs = nil
		}
	case 303:
		//line grammar.y:1928
		{
			yyVAL.comprehensions = yyS[yypt-0].comprehensions
			yyVAL.exprs = yyS[yypt-0].exprs
		}
	case 304:
		//line grammar.y:1935
		{
			c := ast.Comprehension{
				Target: tupleOrExpr(yyVAL.pos, yyS[yypt-2].exprs, yyS[yypt-2].comma),
				Iter:   yyS[yypt-0].expr,
			}
			setCtx(yylex, c.Target, ast.Store)
			yyVAL.comprehensions = []ast.Comprehension{c}
		}
	case 305:
		//line grammar.y:1944
		{
			c := ast.Comprehension{
				Target: tupleOrExpr(yyVAL.pos, yyS[yypt-3].exprs, yyS[yypt-3].comma),
				Iter:   yyS[yypt-1].expr,
				Ifs:    yyS[yypt-0].exprs,
			}
			setCtx(yylex, c.Target, ast.Store)
			yyVAL.comprehensions = []ast.Comprehension{c}
			yyVAL.comprehensions = append(yyVAL.comprehensions, yyS[yypt-0].comprehensions...)
		}
	case 306:
		//line grammar.y:1957
		{
			yyVAL.exprs = []ast.Expr{yyS[yypt-0].expr}
			yyVAL.comprehensions = nil
		}
	case 307:
		//line grammar.y:1962
		{
			yyVAL.exprs = []ast.Expr{yyS[yypt-1].expr}
			yyVAL.exprs = append(yyVAL.exprs, yyS[yypt-0].exprs...)
			yyVAL.comprehensions = yyS[yypt-0].comprehensions
		}
	case 308:
		//line grammar.y:1973
		{
			yyVAL.expr = &ast.Yield{ExprBase: ast.ExprBase{Pos: yyVAL.pos}}
		}
	case 309:
		//line grammar.y:1977
		{
			yyVAL.expr = &ast.YieldFrom{ExprBase: ast.ExprBase{Pos: yyVAL.pos}, Value: yyS[yypt-0].expr}
		}
	case 310:
		//line grammar.y:1981
		{
			yyVAL.expr = &ast.Yield{ExprBase: ast.ExprBase{Pos: yyVAL.pos}, Value: yyS[yypt-0].expr}
		}
	}
	goto yystack /* stack new state and value */
}
