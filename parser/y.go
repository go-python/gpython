//line grammar.y:2
package parser

import __yyfmt__ "fmt"

//line grammar.y:3
// Grammar for Python

import (
	"fmt"
	"github.com/ncw/gpython/ast"
	"github.com/ncw/gpython/py"
)

// NB can put code blocks in not just at the end

// Returns a Tuple if > 1 items or a trailing comma, otherwise returns
// the first item in elts
func tupleOrExpr(pos ast.Pos, elts []ast.Expr, optional_comma bool) ast.Expr {
	if optional_comma || len(elts) > 1 {
		return &ast.Tuple{ExprBase: ast.ExprBase{pos}, Elts: elts, Ctx: ast.Load}
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
	-1, 233,
	68, 14,
	-2, 292,
	-1, 383,
	68, 92,
	-2, 293,
}

const yyNprod = 312
const yyPrivate = 57344

var yyTokenNames []string
var yyStates []string

const yyLast = 1515

var yyAct = []int{

	60, 469, 62, 315, 161, 98, 166, 165, 457, 422,
	402, 376, 322, 350, 362, 343, 142, 465, 211, 336,
	102, 103, 260, 225, 112, 224, 7, 335, 61, 70,
	319, 146, 104, 55, 458, 36, 238, 111, 72, 75,
	73, 67, 65, 147, 151, 96, 106, 58, 71, 138,
	108, 107, 74, 232, 98, 144, 25, 18, 24, 181,
	98, 2, 3, 4, 289, 140, 108, 107, 278, 190,
	134, 97, 76, 233, 248, 244, 285, 379, 85, 153,
	50, 92, 86, 263, 237, 205, 182, 180, 487, 244,
	185, 186, 88, 139, 158, 228, 227, 474, 100, 143,
	386, 149, 475, 167, 49, 155, 91, 89, 90, 388,
	168, 164, 467, 337, 216, 454, 196, 451, 212, 244,
	98, 280, 223, 281, 187, 188, 167, 316, 167, 391,
	197, 200, 189, 399, 164, 421, 342, 282, 396, 82,
	383, 83, 152, 374, 285, 215, 77, 78, 64, 290,
	191, 192, 193, 235, 141, 252, 207, 84, 217, 253,
	79, 256, 385, 198, 201, 236, 292, 239, 481, 240,
	261, 262, 316, 163, 334, 258, 247, 242, 241, 259,
	313, 222, 461, 333, 404, 414, 413, 412, 410, 245,
	406, 160, 453, 243, 251, 250, 163, 420, 341, 264,
	254, 255, 401, 380, 371, 364, 317, 257, 220, 219,
	109, 398, 358, 357, 397, 286, 297, 382, 288, 373,
	268, 291, 98, 269, 294, 272, 273, 356, 112, 267,
	354, 283, 284, 157, 323, 287, 270, 271, 233, 231,
	293, 156, 312, 326, 338, 298, 299, 329, 157, 314,
	266, 108, 107, 405, 305, 265, 221, 157, 339, 306,
	274, 275, 276, 277, 344, 304, 340, 300, 249, 301,
	246, 239, 285, 240, 324, 460, 366, 368, 367, 330,
	285, 323, 351, 462, 125, 126, 135, 131, 123, 121,
	122, 359, 363, 360, 132, 124, 285, 129, 448, 460,
	408, 363, 392, 130, 128, 127, 157, 229, 159, 372,
	308, 347, 208, 183, 108, 107, 377, 378, 355, 184,
	35, 303, 370, 16, 316, 167, 316, 212, 375, 15,
	316, 167, 484, 468, 137, 463, 167, 384, 466, 393,
	431, 337, 353, 381, 434, 331, 140, 115, 261, 395,
	117, 345, 390, 403, 133, 328, 118, 325, 387, 296,
	295, 101, 389, 136, 394, 114, 400, 113, 327, 415,
	218, 99, 213, 214, 310, 8, 409, 309, 230, 311,
	423, 424, 162, 110, 323, 302, 426, 427, 416, 428,
	411, 361, 419, 212, 332, 407, 425, 418, 145, 148,
	351, 150, 437, 176, 433, 439, 429, 441, 440, 442,
	432, 430, 436, 435, 438, 318, 452, 321, 174, 175,
	172, 173, 320, 349, 377, 450, 444, 348, 169, 26,
	120, 194, 449, 204, 105, 459, 443, 206, 445, 446,
	447, 455, 307, 365, 203, 234, 177, 179, 456, 69,
	178, 471, 63, 81, 279, 80, 119, 17, 116, 464,
	14, 13, 433, 470, 12, 10, 11, 45, 323, 44,
	476, 43, 170, 171, 42, 479, 41, 482, 480, 485,
	477, 40, 39, 486, 470, 34, 33, 473, 488, 489,
	470, 210, 209, 85, 32, 31, 92, 86, 30, 29,
	483, 28, 27, 369, 9, 94, 95, 88, 5, 93,
	1, 87, 0, 0, 0, 0, 0, 0, 0, 0,
	0, 91, 89, 90, 0, 0, 48, 51, 25, 52,
	24, 37, 0, 0, 0, 0, 21, 57, 46, 19,
	56, 0, 0, 66, 47, 68, 0, 38, 54, 53,
	22, 20, 23, 59, 82, 85, 83, 417, 92, 86,
	0, 77, 78, 64, 0, 0, 0, 0, 0, 88,
	0, 0, 84, 0, 0, 79, 49, 0, 0, 0,
	0, 0, 0, 91, 89, 90, 0, 0, 48, 51,
	25, 52, 24, 37, 0, 0, 0, 0, 21, 57,
	46, 19, 56, 0, 0, 66, 47, 68, 0, 38,
	54, 53, 22, 20, 23, 59, 82, 0, 83, 0,
	0, 0, 0, 77, 78, 64, 6, 0, 85, 0,
	0, 92, 86, 0, 84, 0, 0, 79, 49, 0,
	0, 0, 88, 0, 0, 0, 0, 0, 0, 0,
	0, 0, 0, 0, 0, 0, 91, 89, 90, 0,
	0, 48, 51, 25, 52, 24, 37, 0, 0, 0,
	0, 21, 57, 46, 19, 56, 0, 0, 66, 47,
	68, 0, 38, 54, 53, 22, 20, 23, 59, 82,
	85, 83, 0, 92, 86, 0, 77, 78, 64, 0,
	0, 0, 0, 0, 88, 0, 0, 84, 0, 0,
	79, 49, 0, 0, 0, 0, 0, 0, 91, 89,
	90, 0, 0, 48, 51, 25, 52, 24, 37, 0,
	0, 0, 0, 21, 57, 46, 19, 56, 0, 0,
	66, 47, 68, 0, 38, 54, 53, 22, 20, 23,
	59, 82, 0, 83, 0, 0, 0, 0, 77, 78,
	64, 226, 0, 85, 0, 0, 92, 86, 0, 84,
	0, 0, 79, 49, 0, 0, 0, 88, 0, 0,
	0, 0, 0, 0, 0, 0, 0, 0, 0, 0,
	0, 91, 89, 90, 0, 0, 48, 51, 0, 52,
	0, 37, 0, 0, 0, 0, 0, 57, 46, 0,
	56, 0, 0, 66, 47, 68, 0, 38, 54, 53,
	0, 0, 0, 59, 82, 85, 83, 0, 92, 86,
	0, 77, 78, 64, 0, 0, 0, 0, 0, 88,
	0, 0, 84, 0, 0, 79, 0, 0, 0, 0,
	0, 0, 0, 91, 89, 90, 0, 0, 48, 51,
	0, 52, 0, 37, 0, 0, 0, 0, 0, 57,
	46, 0, 56, 0, 0, 66, 47, 68, 0, 38,
	54, 53, 0, 0, 0, 59, 82, 85, 83, 0,
	92, 86, 0, 77, 78, 64, 0, 0, 0, 0,
	0, 88, 0, 0, 84, 0, 0, 79, 0, 0,
	0, 0, 0, 0, 0, 91, 89, 90, 0, 0,
	0, 0, 0, 0, 85, 0, 0, 92, 86, 0,
	0, 0, 0, 0, 0, 0, 0, 66, 88, 68,
	0, 0, 0, 0, 0, 0, 0, 59, 82, 195,
	83, 0, 91, 89, 90, 77, 78, 64, 0, 0,
	0, 85, 0, 0, 92, 86, 84, 0, 0, 79,
	0, 0, 0, 0, 66, 88, 68, 0, 0, 0,
	0, 0, 0, 0, 59, 82, 0, 83, 0, 91,
	89, 90, 77, 78, 64, 0, 0, 0, 0, 0,
	0, 0, 0, 84, 85, 0, 79, 92, 86, 0,
	0, 66, 478, 68, 0, 0, 0, 0, 88, 0,
	0, 0, 82, 0, 83, 199, 0, 0, 0, 77,
	78, 64, 91, 89, 90, 0, 0, 0, 0, 0,
	84, 85, 0, 79, 92, 86, 0, 0, 0, 0,
	0, 0, 0, 0, 66, 88, 68, 0, 0, 0,
	0, 0, 0, 0, 0, 82, 0, 83, 0, 91,
	89, 90, 77, 78, 0, 0, 0, 0, 0, 0,
	0, 0, 0, 84, 85, 0, 79, 92, 86, 0,
	0, 66, 0, 68, 0, 0, 0, 0, 88, 0,
	0, 0, 82, 0, 83, 0, 404, 0, 0, 77,
	78, 0, 91, 89, 90, 0, 0, 0, 0, 0,
	84, 0, 0, 79, 0, 0, 85, 0, 0, 92,
	86, 0, 0, 0, 66, 0, 68, 0, 0, 0,
	88, 0, 0, 0, 0, 82, 0, 83, 0, 352,
	0, 0, 77, 78, 91, 89, 90, 0, 0, 0,
	0, 0, 0, 84, 0, 0, 79, 0, 85, 0,
	0, 92, 86, 0, 0, 0, 66, 0, 68, 0,
	0, 0, 88, 0, 0, 0, 0, 82, 346, 83,
	0, 0, 0, 0, 77, 78, 91, 89, 90, 0,
	0, 0, 0, 0, 0, 84, 0, 0, 79, 0,
	0, 85, 0, 0, 92, 86, 0, 0, 66, 0,
	68, 0, 0, 0, 0, 88, 0, 0, 0, 82,
	0, 83, 0, 0, 0, 0, 77, 78, 64, 91,
	89, 90, 0, 0, 0, 0, 0, 84, 85, 0,
	79, 92, 86, 0, 0, 0, 0, 0, 0, 0,
	0, 66, 88, 68, 0, 0, 0, 0, 0, 0,
	0, 59, 82, 0, 83, 0, 91, 89, 90, 77,
	78, 0, 0, 0, 0, 85, 0, 0, 92, 86,
	84, 0, 0, 79, 0, 0, 0, 0, 66, 88,
	68, 0, 0, 0, 0, 0, 0, 0, 0, 82,
	0, 83, 0, 91, 89, 90, 77, 78, 0, 0,
	0, 0, 85, 0, 0, 92, 86, 84, 202, 154,
	79, 0, 0, 0, 0, 66, 88, 68, 0, 0,
	0, 0, 0, 0, 0, 0, 82, 0, 83, 0,
	91, 89, 90, 77, 78, 0, 0, 0, 0, 85,
	0, 0, 92, 86, 84, 0, 0, 79, 0, 0,
	0, 0, 472, 88, 68, 0, 0, 0, 0, 0,
	0, 0, 0, 82, 0, 83, 0, 91, 89, 90,
	77, 78, 0, 0, 0, 0, 85, 0, 0, 92,
	86, 84, 0, 0, 79, 0, 0, 0, 0, 66,
	88, 68, 0, 0, 0, 0, 0, 0, 0, 0,
	82, 0, 83, 0, 91, 89, 90, 77, 78, 0,
	0, 0, 85, 0, 0, 92, 86, 0, 84, 0,
	0, 79, 0, 0, 0, 0, 88, 0, 68, 0,
	0, 0, 0, 0, 0, 0, 0, 82, 0, 83,
	91, 89, 90, 0, 77, 78, 0, 0, 0, 0,
	0, 0, 0, 0, 0, 84, 0, 0, 79, 0,
	0, 0, 0, 0, 0, 0, 0, 0, 0, 0,
	0, 0, 0, 82, 0, 83, 0, 0, 0, 0,
	77, 78, 0, 0, 0, 0, 0, 0, 0, 0,
	0, 84, 0, 0, 79,
}
var yyPact = []int{

	-29, -1000, 622, -1000, 1353, -1000, -1000, -1000, 367, 25,
	-1000, -1000, -1000, -1000, -1000, -1000, -1000, -1000, -1000, 1353,
	1353, 72, 139, 1353, 361, 359, 15, -1000, -1000, -1000,
	-1000, -1000, -1000, -1000, -1000, -1000, 272, 72, -1000, -1000,
	-1000, -1000, -1000, -1000, -1000, -1000, 357, 357, 1353, 340,
	82, -1000, -1000, 1353, 1353, -1000, 340, 59, -1000, 1279,
	-1000, -1000, 189, -1000, 1426, 271, 120, -1000, 1390, 392,
	9, -28, 7, 289, 16, 48, -1000, 1426, 1426, 1426,
	-1000, -1000, 881, 955, 1242, -1000, -1000, 303, -1000, -1000,
	-1000, -1000, -1000, -1000, 487, -1000, -1000, 73, -1000, -1000,
	819, 366, 138, 137, 202, 109, -1000, 9, -1000, 757,
	24, -1000, 269, 172, 171, -1000, -1000, -1000, -1000, 1205,
	2, -1000, -1000, -1000, -1000, -1000, -1000, -1000, -1000, -1000,
	-1000, -1000, -1000, 918, -1000, 106, -1000, 106, 105, 6,
	-1000, 1162, -1000, -1000, 220, 104, -1000, 36, 215, -8,
	59, -1000, -1000, -1000, 1353, -1000, 1390, 1390, 9, 1390,
	1353, 136, 103, 325, 325, -1000, 1, -1000, -1000, 1426,
	-1000, -1000, -1000, -1000, -1000, -1000, -1000, -1000, 201, 192,
	1426, 1426, 1426, 1426, 1426, 1426, 1426, 1426, 1426, 1426,
	1426, -1000, -1000, -1000, 54, -1000, 163, 231, 82, -1000,
	231, 82, -1000, -22, 77, 95, -1000, 73, -1000, -1000,
	-1000, -1000, -1000, -1000, 355, 1353, -1000, -1000, -1000, 757,
	757, 1353, 72, -1000, -1000, -1000, 314, 1353, 757, 1426,
	291, 166, 135, 1353, -1000, -1000, -1000, 918, -1000, -1000,
	-1000, 351, 1353, 364, 349, -1000, 1353, 340, 339, 107,
	-1000, -8, -1000, 198, 271, -1000, -1000, 1353, 122, -1000,
	-1000, -1000, -1000, 1353, 9, -1000, -1000, -28, 7, 289,
	16, 16, 48, 48, -1000, -1000, -1000, -1000, 1426, -1000,
	1120, 1078, 336, -1000, 162, 72, 159, 143, 142, -1000,
	1353, -1000, 1353, -1000, -1000, -1000, -1000, -1000, -1000, 246,
	134, -1000, 230, 684, -1000, -1000, 9, 133, 1353, 151,
	-1000, 71, 320, 320, -1000, -5, 132, 757, 149, -1000,
	68, 86, -1000, 27, -1000, -1000, -1000, -1000, -1000, -1000,
	-1000, -1000, -1000, -1000, 335, 57, -1000, 264, 1353, -1000,
	-1000, 325, 325, 66, -1000, -1000, -1000, 146, 141, 61,
	-1000, 131, 1035, -1000, -1000, 199, -1000, -1000, -1000, 119,
	231, 255, -1000, 117, 757, 116, 115, 114, 1353, 549,
	-1000, 757, -1000, -1000, 121, -1000, -1000, -1000, -1000, 1353,
	1353, -1000, -1000, 1353, -1000, 1353, 1353, -1000, 1353, 57,
	-1000, 335, 334, -1000, -1000, -1000, 330, -1000, -1000, 1078,
	-1000, 1035, -1000, 113, 1353, 1390, 1353, -1000, 1353, -1000,
	757, 246, 757, 757, 757, 260, -1000, -1000, -1000, -1000,
	320, 320, 45, -1000, -1000, -1000, -1000, -1000, -1000, 124,
	-1000, -1000, 43, -1000, 325, -1000, -1000, 113, -1000, -1000,
	247, -1000, 111, -1000, -1000, -1000, 235, -1000, 329, -1000,
	-1000, 324, 40, -1000, 319, -1000, -1000, -1000, -1000, -1000,
	1316, 757, 26, -1000, 30, -1000, 320, 998, 325, 223,
	174, -1000, 97, -1000, 757, 318, -1000, -1000, 1353, -1000,
	-1000, 1316, 17, -1000, 320, -1000, -1000, 1316, -1000, -1000,
}
var yyPgo = []int{

	0, 511, 510, 509, 508, 506, 23, 18, 505, 504,
	503, 25, 14, 372, 57, 502, 501, 499, 498, 495,
	494, 486, 485, 482, 481, 476, 474, 471, 469, 467,
	466, 465, 464, 461, 460, 329, 323, 458, 457, 456,
	46, 29, 28, 48, 38, 40, 52, 39, 72, 455,
	454, 453, 47, 0, 41, 452, 1, 451, 2, 42,
	449, 45, 35, 445, 33, 36, 444, 10, 443, 442,
	320, 32, 437, 435, 8, 434, 80, 71, 433, 431,
	430, 429, 428, 16, 34, 13, 427, 423, 12, 422,
	417, 416, 30, 53, 415, 44, 401, 43, 399, 286,
	31, 19, 398, 27, 394, 391, 385, 37, 383, 7,
	6, 22, 17, 3, 11, 15, 382, 9, 379, 4,
	378, 377, 374, 373, 361,
}
var yyR1 = []int{

	0, 2, 2, 2, 4, 4, 4, 3, 8, 8,
	8, 5, 123, 123, 94, 94, 93, 93, 70, 81,
	81, 37, 37, 38, 69, 69, 35, 120, 121, 121,
	112, 112, 117, 117, 118, 118, 114, 114, 122, 122,
	122, 122, 122, 122, 122, 113, 113, 109, 109, 115,
	115, 116, 116, 111, 111, 119, 119, 119, 119, 119,
	119, 119, 110, 7, 7, 124, 124, 9, 9, 6,
	14, 14, 14, 14, 14, 14, 14, 14, 15, 15,
	15, 63, 63, 65, 65, 80, 80, 76, 76, 52,
	52, 83, 83, 62, 39, 39, 39, 39, 39, 39,
	39, 39, 39, 39, 39, 39, 16, 17, 18, 18,
	18, 18, 18, 23, 24, 25, 25, 27, 26, 26,
	26, 19, 19, 28, 95, 95, 96, 96, 98, 98,
	98, 104, 104, 104, 29, 101, 101, 100, 100, 103,
	103, 102, 102, 97, 97, 99, 99, 20, 21, 77,
	77, 22, 22, 13, 13, 13, 13, 13, 13, 13,
	13, 105, 105, 12, 12, 31, 30, 32, 106, 106,
	33, 33, 33, 33, 108, 108, 34, 107, 107, 68,
	68, 68, 10, 10, 11, 11, 53, 53, 53, 56,
	56, 55, 55, 57, 57, 58, 58, 59, 59, 54,
	54, 60, 60, 82, 82, 82, 82, 82, 82, 82,
	82, 82, 82, 82, 42, 41, 41, 43, 43, 44,
	44, 45, 45, 45, 46, 46, 46, 47, 47, 47,
	47, 47, 48, 48, 48, 48, 49, 49, 79, 79,
	1, 1, 51, 51, 51, 51, 51, 51, 51, 51,
	51, 51, 51, 51, 51, 51, 51, 51, 50, 50,
	50, 50, 87, 87, 86, 85, 85, 85, 85, 85,
	85, 85, 85, 85, 67, 67, 40, 40, 75, 75,
	71, 61, 72, 78, 78, 66, 66, 66, 66, 36,
	89, 89, 90, 90, 91, 91, 92, 92, 92, 92,
	88, 88, 88, 74, 74, 84, 84, 73, 73, 64,
	64, 64,
}
var yyR2 = []int{

	0, 2, 2, 2, 1, 1, 2, 2, 0, 2,
	2, 3, 0, 2, 0, 1, 0, 3, 4, 1,
	2, 1, 1, 2, 0, 2, 6, 3, 0, 1,
	1, 3, 0, 3, 1, 3, 0, 1, 2, 5,
	8, 4, 3, 6, 2, 1, 3, 1, 3, 0,
	3, 1, 3, 0, 1, 2, 5, 8, 4, 3,
	6, 2, 1, 1, 1, 0, 1, 1, 3, 3,
	1, 1, 1, 1, 1, 1, 1, 1, 3, 2,
	1, 1, 1, 1, 1, 2, 3, 1, 3, 1,
	1, 0, 1, 2, 1, 1, 1, 1, 1, 1,
	1, 1, 1, 1, 1, 1, 2, 1, 1, 1,
	1, 1, 1, 1, 1, 1, 2, 1, 1, 2,
	4, 1, 1, 2, 1, 1, 1, 2, 1, 2,
	1, 1, 4, 2, 4, 1, 3, 1, 3, 1,
	3, 1, 3, 1, 3, 1, 3, 2, 2, 1,
	3, 2, 4, 1, 1, 1, 1, 1, 1, 1,
	1, 0, 5, 0, 3, 6, 5, 7, 0, 4,
	4, 7, 7, 10, 1, 3, 4, 1, 3, 1,
	2, 4, 1, 2, 1, 4, 1, 5, 1, 1,
	1, 3, 4, 3, 4, 1, 3, 1, 3, 2,
	1, 1, 3, 1, 1, 1, 1, 1, 1, 1,
	1, 2, 1, 2, 2, 1, 3, 1, 3, 1,
	3, 1, 3, 3, 1, 3, 3, 1, 3, 3,
	3, 3, 2, 2, 2, 1, 2, 4, 0, 2,
	1, 2, 2, 3, 4, 4, 2, 4, 4, 2,
	3, 1, 1, 1, 1, 1, 1, 1, 2, 3,
	3, 2, 1, 3, 2, 1, 1, 2, 2, 3,
	2, 3, 3, 4, 1, 2, 1, 1, 1, 3,
	2, 2, 2, 3, 5, 2, 4, 1, 2, 5,
	1, 3, 0, 2, 0, 3, 2, 4, 7, 3,
	1, 2, 3, 1, 1, 4, 5, 2, 3, 1,
	3, 2,
}
var yyChk = []int{

	-1000, -2, 90, 91, 92, -4, 4, -6, -13, -9,
	-31, -30, -32, -33, -34, -35, -36, -38, -14, 52,
	64, 49, 63, 65, 43, 41, -81, -15, -16, -17,
	-18, -19, -20, -21, -22, -70, -62, 44, 60, -23,
	-24, -25, -26, -27, -28, -29, 51, 57, 39, 89,
	-76, 40, 42, 62, 61, -64, 53, 50, -52, 66,
	-53, -42, -58, -55, 76, -59, 56, -54, 58, -60,
	-41, -43, -44, -45, -46, -47, -48, 74, 75, 88,
	-49, -51, 67, 69, 85, 6, 10, -1, 20, 35,
	36, 34, 9, -3, -8, -5, -61, -77, -53, 4,
	73, -124, -53, -53, -71, -75, -40, -41, -42, 71,
	-108, -107, -53, 6, 6, -70, -37, -36, -35, -39,
	-80, 17, 18, 16, 23, 12, 13, 33, 32, 25,
	31, 15, 22, 82, -71, -99, 6, -99, -53, -97,
	6, 72, -83, -61, -53, -102, -100, -97, -98, -97,
	-96, -95, 83, 20, 50, -61, 52, 59, -41, 37,
	71, -119, -116, 76, 14, -109, -110, 6, -54, -82,
	80, 81, 28, 29, 26, 27, 11, 54, 58, 55,
	78, 87, 79, 24, 30, 74, 75, 76, 77, 84,
	21, -48, -48, -48, -79, 68, -64, -52, -76, 70,
	-52, -76, 86, -66, -78, -53, -72, -77, 9, 5,
	4, -7, -6, -13, -123, 72, -83, -14, 4, 71,
	71, 54, 72, -83, -11, -6, 4, 72, 71, 38,
	-120, 67, -93, 67, -63, -64, -61, 82, -65, -64,
	-62, 72, 72, -93, 83, -52, 50, 72, 38, 53,
	-95, -97, -53, -58, -59, -54, -53, 71, 72, -83,
	-111, -110, -110, 82, -41, 54, 58, -43, -44, -45,
	-46, -46, -47, -47, -48, -48, -48, -48, 14, -50,
	67, 69, 83, 68, -84, 49, -83, -84, -83, 86,
	72, -83, 71, -84, -83, 5, 4, -53, -11, -11,
	-61, -40, -106, 7, -107, -11, -41, -69, 19, -121,
	-122, -118, 76, 14, -112, -113, 6, 71, -94, -92,
	-89, -90, -88, -53, -65, 6, -53, 4, 6, -53,
	-100, 6, -104, 76, 67, -103, -101, 6, 46, -53,
	-109, 76, 14, -115, -53, -48, 68, -92, -86, -87,
	-85, -53, 71, 6, 68, -71, 68, 70, 70, -53,
	-53, -105, -12, 46, 71, -68, 46, 48, 47, -10,
	-7, 71, -53, 68, 72, -83, -114, -113, -113, 82,
	71, -11, 68, 72, -83, 76, 14, -84, 82, -103,
	-83, 72, 38, -53, -111, -110, 72, 68, 70, 72,
	-83, 71, -67, -53, 71, 54, 71, -84, 45, -12,
	71, -11, 71, 71, 71, -53, -7, 8, -11, -112,
	76, 14, -117, -53, -53, -88, -53, -53, -53, -83,
	-101, 6, -115, -109, 14, -85, -67, -53, -67, -53,
	-58, -53, -53, -11, -12, -11, -11, -11, 38, -114,
	-113, 72, -91, 68, 72, -110, -67, -74, -84, -73,
	52, 71, 48, 6, -117, -112, 14, 72, 14, -56,
	-58, -57, 56, -11, 71, 72, -113, -88, 14, -110,
	-74, 71, -119, -11, 14, -53, -56, 71, -113, -56,
}
var yyDef = []int{

	0, -2, 0, 8, 0, 1, 4, 5, 0, 65,
	153, 154, 155, 156, 157, 158, 159, 160, 67, 0,
	0, 0, 0, 0, 0, 0, 0, 70, 71, 72,
	73, 74, 75, 76, 77, 19, 80, 0, 107, 108,
	109, 110, 111, 112, 121, 122, 0, 0, 0, 0,
	91, 113, 114, 115, 118, 117, 0, 0, 87, 309,
	89, 90, 186, 188, 0, 195, 0, 197, 0, 200,
	201, 215, 217, 219, 221, 224, 227, 0, 0, 0,
	235, 238, 0, 0, 0, 251, 252, 253, 254, 255,
	256, 257, 240, 2, 0, 3, 12, 91, 149, 6,
	66, 0, 0, 0, 0, 91, 278, 276, 277, 0,
	0, 174, 177, 0, 16, 20, 23, 21, 22, 0,
	79, 94, 95, 96, 97, 98, 99, 100, 101, 102,
	103, 104, 105, 0, 106, 147, 145, 148, 151, 16,
	143, 92, 93, 116, 119, 123, 141, 137, 0, 128,
	130, 126, 124, 125, 0, 311, 0, 0, 214, 0,
	0, 0, 91, 53, 0, 51, 47, 62, 199, 0,
	203, 204, 205, 206, 207, 208, 209, 210, 0, 212,
	0, 0, 0, 0, 0, 0, 0, 0, 0, 0,
	0, 232, 233, 234, 236, 242, 0, 87, 91, 246,
	87, 91, 249, 0, 91, 149, 287, 91, 241, 7,
	9, 10, 63, 64, 0, 92, 281, 68, 69, 0,
	0, 0, 92, 280, 168, 184, 0, 0, 0, 0,
	24, 28, 0, -2, 78, 81, 82, 0, 85, 83,
	84, 0, 0, 0, 0, 88, 0, 0, 0, 0,
	127, 129, 310, 0, 196, 198, 191, 0, 92, 55,
	49, 54, 61, 0, 202, 211, 213, 216, 218, 220,
	222, 223, 225, 226, 228, 229, 230, 231, 0, 239,
	292, 0, 0, 243, 0, 0, 0, 0, 0, 250,
	92, 285, 0, 288, 282, 11, 13, 150, 161, 163,
	0, 279, 170, 0, 175, 176, 178, 0, 0, 0,
	29, 91, 36, 0, 34, 30, 45, 0, 0, 15,
	91, 0, 290, 300, 86, 146, 152, 18, 144, 120,
	142, 138, 134, 131, 0, 91, 139, 135, 0, 192,
	52, 53, 0, 59, 48, 237, 258, 0, 0, 91,
	262, 265, 266, 261, 244, 0, 245, 247, 248, 0,
	283, 163, 166, 0, 0, 0, 0, 0, 179, 0,
	182, 0, 25, 27, 92, 38, 32, 37, 44, 0,
	0, 289, 17, -2, 296, 0, 0, 301, 0, 91,
	133, 92, 0, 187, 49, 58, 0, 259, 260, 92,
	264, 270, 267, 268, 274, 0, 0, 286, 0, 165,
	0, 163, 0, 0, 0, 180, 183, 185, 26, 35,
	36, 0, 42, 31, 46, 291, 294, 299, 302, 0,
	140, 136, 56, 50, 0, 263, 271, 272, 269, 275,
	305, 284, 0, 164, 167, 169, 171, 172, 0, 32,
	41, 0, 297, 132, 0, 60, 273, 306, 303, 304,
	0, 0, 0, 181, 39, 33, 0, 0, 0, 307,
	189, 190, 0, 162, 0, 0, 43, 295, 0, 57,
	308, 0, 0, 173, 0, 298, 193, 0, 40, 194,
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
		//line grammar.y:263
		{
			// panic("FIXME no coverage")
			yyVAL.mod = &ast.Interactive{ModBase: ast.ModBase{yyVAL.pos}}
		}
	case 5:
		//line grammar.y:268
		{
			yyVAL.mod = &ast.Interactive{ModBase: ast.ModBase{yyVAL.pos}, Body: yyS[yypt-0].stmts}
		}
	case 6:
		//line grammar.y:272
		{
			// panic("FIXME no coverage")
			yyVAL.mod = &ast.Interactive{ModBase: ast.ModBase{yyVAL.pos}, Body: []ast.Stmt{yyS[yypt-1].stmt}}
		}
	case 7:
		//line grammar.y:280
		{
			yyVAL.mod = &ast.Module{ModBase: ast.ModBase{yyVAL.pos}, Body: yyS[yypt-1].stmts}
		}
	case 8:
		//line grammar.y:286
		{
			yyVAL.stmts = nil
		}
	case 9:
		//line grammar.y:290
		{
		}
	case 10:
		//line grammar.y:293
		{
			yyVAL.stmts = append(yyVAL.stmts, yyS[yypt-0].stmts...)
		}
	case 11:
		//line grammar.y:300
		{
			yyVAL.mod = &ast.Expression{ModBase: ast.ModBase{yyVAL.pos}, Body: yyS[yypt-2].expr}
		}
	case 14:
		//line grammar.y:309
		{
			yyVAL.call = &ast.Call{ExprBase: ast.ExprBase{yyVAL.pos}}
		}
	case 15:
		//line grammar.y:313
		{
			yyVAL.call = yyS[yypt-0].call
		}
	case 16:
		//line grammar.y:318
		{
			yyVAL.call = nil
		}
	case 17:
		//line grammar.y:322
		{
			yyVAL.call = yyS[yypt-1].call
		}
	case 18:
		//line grammar.y:328
		{
			fn := &ast.Name{ExprBase: ast.ExprBase{yyVAL.pos}, Id: ast.Identifier(yyS[yypt-2].str), Ctx: ast.Load}
			if yyS[yypt-1].call == nil {
				yyVAL.expr = fn
			} else {
				call := *yyS[yypt-1].call
				call.Func = fn
				yyVAL.expr = &call
			}
		}
	case 19:
		//line grammar.y:341
		{
			yyVAL.exprs = nil
			yyVAL.exprs = append(yyVAL.exprs, yyS[yypt-0].expr)
		}
	case 20:
		//line grammar.y:346
		{
			yyVAL.exprs = append(yyVAL.exprs, yyS[yypt-0].expr)
		}
	case 21:
		//line grammar.y:352
		{
			yyVAL.stmt = yyS[yypt-0].stmt
		}
	case 22:
		//line grammar.y:356
		{
			yyVAL.stmt = yyS[yypt-0].stmt
		}
	case 23:
		//line grammar.y:362
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
	case 24:
		//line grammar.y:376
		{
			yyVAL.expr = nil
		}
	case 25:
		//line grammar.y:380
		{
			yyVAL.expr = yyS[yypt-0].expr
		}
	case 26:
		//line grammar.y:386
		{
			yyVAL.stmt = &ast.FunctionDef{StmtBase: ast.StmtBase{yyVAL.pos}, Name: ast.Identifier(yyS[yypt-4].str), Args: yyS[yypt-3].arguments, Body: yyS[yypt-0].stmts, Returns: yyS[yypt-2].expr}
		}
	case 27:
		//line grammar.y:392
		{
			yyVAL.arguments = yyS[yypt-1].arguments
		}
	case 28:
		//line grammar.y:397
		{
			yyVAL.arguments = &ast.Arguments{Pos: yyVAL.pos}
		}
	case 29:
		//line grammar.y:401
		{
			yyVAL.arguments = yyS[yypt-0].arguments
		}
	case 30:
		//line grammar.y:408
		{
			yyVAL.arg = yyS[yypt-0].arg
			yyVAL.expr = nil
		}
	case 31:
		//line grammar.y:413
		{
			yyVAL.arg = yyS[yypt-2].arg
			yyVAL.expr = yyS[yypt-0].expr
		}
	case 32:
		//line grammar.y:419
		{
			yyVAL.args = nil
			yyVAL.exprs = nil
		}
	case 33:
		//line grammar.y:424
		{
			yyVAL.args = append(yyVAL.args, yyS[yypt-0].arg)
			if yyS[yypt-0].expr != nil {
				yyVAL.exprs = append(yyVAL.exprs, yyS[yypt-0].expr)
			}
		}
	case 34:
		//line grammar.y:433
		{
			yyVAL.args = nil
			yyVAL.args = append(yyVAL.args, yyS[yypt-0].arg)
			yyVAL.exprs = nil
			if yyS[yypt-0].expr != nil {
				yyVAL.exprs = append(yyVAL.exprs, yyS[yypt-0].expr)
			}
		}
	case 35:
		//line grammar.y:442
		{
			yyVAL.args = append(yyVAL.args, yyS[yypt-0].arg)
			if yyS[yypt-0].expr != nil {
				yyVAL.exprs = append(yyVAL.exprs, yyS[yypt-0].expr)
			}
		}
	case 36:
		//line grammar.y:450
		{
			yyVAL.arg = nil
		}
	case 37:
		//line grammar.y:454
		{
			yyVAL.arg = yyS[yypt-0].arg
		}
	case 38:
		//line grammar.y:461
		{
			yyVAL.arguments = &ast.Arguments{Pos: yyVAL.pos, Args: yyS[yypt-1].args, Defaults: yyS[yypt-1].exprs}
		}
	case 39:
		//line grammar.y:465
		{
			yyVAL.arguments = &ast.Arguments{Pos: yyVAL.pos, Args: yyS[yypt-4].args, Defaults: yyS[yypt-4].exprs, Vararg: yyS[yypt-1].arg, Kwonlyargs: yyS[yypt-0].args, KwDefaults: yyS[yypt-0].exprs}
		}
	case 40:
		//line grammar.y:469
		{
			yyVAL.arguments = &ast.Arguments{Pos: yyVAL.pos, Args: yyS[yypt-7].args, Defaults: yyS[yypt-7].exprs, Vararg: yyS[yypt-4].arg, Kwonlyargs: yyS[yypt-3].args, KwDefaults: yyS[yypt-3].exprs, Kwarg: yyS[yypt-0].arg}
		}
	case 41:
		//line grammar.y:473
		{
			yyVAL.arguments = &ast.Arguments{Pos: yyVAL.pos, Args: yyS[yypt-3].args, Defaults: yyS[yypt-3].exprs, Kwarg: yyS[yypt-0].arg}
		}
	case 42:
		//line grammar.y:477
		{
			yyVAL.arguments = &ast.Arguments{Pos: yyVAL.pos, Vararg: yyS[yypt-1].arg, Kwonlyargs: yyS[yypt-0].args, KwDefaults: yyS[yypt-0].exprs}
		}
	case 43:
		//line grammar.y:481
		{
			yyVAL.arguments = &ast.Arguments{Pos: yyVAL.pos, Vararg: yyS[yypt-4].arg, Kwonlyargs: yyS[yypt-3].args, KwDefaults: yyS[yypt-3].exprs, Kwarg: yyS[yypt-0].arg}
		}
	case 44:
		//line grammar.y:485
		{
			yyVAL.arguments = &ast.Arguments{Pos: yyVAL.pos, Kwarg: yyS[yypt-0].arg}
		}
	case 45:
		//line grammar.y:491
		{
			yyVAL.arg = &ast.Arg{Pos: yyVAL.pos, Arg: ast.Identifier(yyS[yypt-0].str)}
		}
	case 46:
		//line grammar.y:495
		{
			yyVAL.arg = &ast.Arg{Pos: yyVAL.pos, Arg: ast.Identifier(yyS[yypt-2].str), Annotation: yyS[yypt-0].expr}
		}
	case 47:
		//line grammar.y:501
		{
			yyVAL.arg = yyS[yypt-0].arg
			yyVAL.expr = nil
		}
	case 48:
		//line grammar.y:506
		{
			yyVAL.arg = yyS[yypt-2].arg
			yyVAL.expr = yyS[yypt-0].expr
		}
	case 49:
		//line grammar.y:512
		{
			yyVAL.args = nil
			yyVAL.exprs = nil
		}
	case 50:
		//line grammar.y:517
		{
			yyVAL.args = append(yyVAL.args, yyS[yypt-0].arg)
			if yyS[yypt-0].expr != nil {
				yyVAL.exprs = append(yyVAL.exprs, yyS[yypt-0].expr)
			}
		}
	case 51:
		//line grammar.y:526
		{
			yyVAL.args = nil
			yyVAL.args = append(yyVAL.args, yyS[yypt-0].arg)
			yyVAL.exprs = nil
			if yyS[yypt-0].expr != nil {
				yyVAL.exprs = append(yyVAL.exprs, yyS[yypt-0].expr)
			}
		}
	case 52:
		//line grammar.y:535
		{
			yyVAL.args = append(yyVAL.args, yyS[yypt-0].arg)
			if yyS[yypt-0].expr != nil {
				yyVAL.exprs = append(yyVAL.exprs, yyS[yypt-0].expr)
			}
		}
	case 53:
		//line grammar.y:543
		{
			yyVAL.arg = nil
		}
	case 54:
		//line grammar.y:547
		{
			yyVAL.arg = yyS[yypt-0].arg
		}
	case 55:
		//line grammar.y:554
		{
			yyVAL.arguments = &ast.Arguments{Pos: yyVAL.pos, Args: yyS[yypt-1].args, Defaults: yyS[yypt-1].exprs}
		}
	case 56:
		//line grammar.y:558
		{
			yyVAL.arguments = &ast.Arguments{Pos: yyVAL.pos, Args: yyS[yypt-4].args, Defaults: yyS[yypt-4].exprs, Vararg: yyS[yypt-1].arg, Kwonlyargs: yyS[yypt-0].args, KwDefaults: yyS[yypt-0].exprs}
		}
	case 57:
		//line grammar.y:562
		{
			yyVAL.arguments = &ast.Arguments{Pos: yyVAL.pos, Args: yyS[yypt-7].args, Defaults: yyS[yypt-7].exprs, Vararg: yyS[yypt-4].arg, Kwonlyargs: yyS[yypt-3].args, KwDefaults: yyS[yypt-3].exprs, Kwarg: yyS[yypt-0].arg}
		}
	case 58:
		//line grammar.y:566
		{
			yyVAL.arguments = &ast.Arguments{Pos: yyVAL.pos, Args: yyS[yypt-3].args, Defaults: yyS[yypt-3].exprs, Kwarg: yyS[yypt-0].arg}
		}
	case 59:
		//line grammar.y:570
		{
			yyVAL.arguments = &ast.Arguments{Pos: yyVAL.pos, Vararg: yyS[yypt-1].arg, Kwonlyargs: yyS[yypt-0].args, KwDefaults: yyS[yypt-0].exprs}
		}
	case 60:
		//line grammar.y:574
		{
			yyVAL.arguments = &ast.Arguments{Pos: yyVAL.pos, Vararg: yyS[yypt-4].arg, Kwonlyargs: yyS[yypt-3].args, KwDefaults: yyS[yypt-3].exprs, Kwarg: yyS[yypt-0].arg}
		}
	case 61:
		//line grammar.y:578
		{
			yyVAL.arguments = &ast.Arguments{Pos: yyVAL.pos, Kwarg: yyS[yypt-0].arg}
		}
	case 62:
		//line grammar.y:584
		{
			yyVAL.arg = &ast.Arg{Pos: yyVAL.pos, Arg: ast.Identifier(yyS[yypt-0].str)}
		}
	case 63:
		//line grammar.y:590
		{
			yyVAL.stmts = yyS[yypt-0].stmts
		}
	case 64:
		//line grammar.y:594
		{
			yyVAL.stmts = []ast.Stmt{yyS[yypt-0].stmt}
		}
	case 67:
		//line grammar.y:602
		{
			yyVAL.stmts = nil
			yyVAL.stmts = append(yyVAL.stmts, yyS[yypt-0].stmt)
		}
	case 68:
		//line grammar.y:607
		{
			yyVAL.stmts = append(yyVAL.stmts, yyS[yypt-0].stmt)
		}
	case 69:
		//line grammar.y:613
		{
			yyVAL.stmts = yyS[yypt-2].stmts
		}
	case 70:
		//line grammar.y:619
		{
			yyVAL.stmt = yyS[yypt-0].stmt
		}
	case 71:
		//line grammar.y:623
		{
			yyVAL.stmt = yyS[yypt-0].stmt
		}
	case 72:
		//line grammar.y:627
		{
			yyVAL.stmt = yyS[yypt-0].stmt
		}
	case 73:
		//line grammar.y:631
		{
			yyVAL.stmt = yyS[yypt-0].stmt
		}
	case 74:
		//line grammar.y:635
		{
			yyVAL.stmt = yyS[yypt-0].stmt
		}
	case 75:
		//line grammar.y:639
		{
			yyVAL.stmt = yyS[yypt-0].stmt
		}
	case 76:
		//line grammar.y:643
		{
			yyVAL.stmt = yyS[yypt-0].stmt
		}
	case 77:
		//line grammar.y:647
		{
			yyVAL.stmt = yyS[yypt-0].stmt
		}
	case 78:
		//line grammar.y:674
		{
			target := yyS[yypt-2].expr
			setCtx(yylex, target, ast.Store)
			yyVAL.stmt = &ast.AugAssign{StmtBase: ast.StmtBase{yyVAL.pos}, Target: target, Op: yyS[yypt-1].op, Value: yyS[yypt-0].expr}
		}
	case 79:
		//line grammar.y:680
		{
			targets := []ast.Expr{yyS[yypt-1].expr}
			targets = append(targets, yyS[yypt-0].exprs...)
			value := targets[len(targets)-1]
			targets = targets[:len(targets)-1]
			setCtxs(yylex, targets, ast.Store)
			yyVAL.stmt = &ast.Assign{StmtBase: ast.StmtBase{yyVAL.pos}, Targets: targets, Value: value}
		}
	case 80:
		//line grammar.y:689
		{
			yyVAL.stmt = &ast.ExprStmt{StmtBase: ast.StmtBase{yyVAL.pos}, Value: yyS[yypt-0].expr}
		}
	case 81:
		//line grammar.y:695
		{
			yyVAL.expr = yyS[yypt-0].expr
		}
	case 82:
		//line grammar.y:699
		{
			yyVAL.expr = yyS[yypt-0].expr
		}
	case 83:
		//line grammar.y:705
		{
			yyVAL.expr = yyS[yypt-0].expr
		}
	case 84:
		//line grammar.y:709
		{
			yyVAL.expr = yyS[yypt-0].expr
		}
	case 85:
		//line grammar.y:715
		{
			yyVAL.exprs = nil
			yyVAL.exprs = append(yyVAL.exprs, yyS[yypt-0].expr)
		}
	case 86:
		//line grammar.y:720
		{
			yyVAL.exprs = append(yyVAL.exprs, yyS[yypt-0].expr)
		}
	case 87:
		//line grammar.y:726
		{
			yyVAL.exprs = nil
			yyVAL.exprs = append(yyVAL.exprs, yyS[yypt-0].expr)
		}
	case 88:
		//line grammar.y:731
		{
			yyVAL.exprs = append(yyVAL.exprs, yyS[yypt-0].expr)
		}
	case 89:
		//line grammar.y:737
		{
			yyVAL.expr = yyS[yypt-0].expr
		}
	case 90:
		//line grammar.y:741
		{
			yyVAL.expr = yyS[yypt-0].expr
		}
	case 91:
		//line grammar.y:746
		{
			yyVAL.comma = false
		}
	case 92:
		//line grammar.y:750
		{
			yyVAL.comma = true
		}
	case 93:
		//line grammar.y:756
		{
			yyVAL.expr = tupleOrExpr(yyVAL.pos, yyS[yypt-1].exprs, yyS[yypt-0].comma)
		}
	case 94:
		//line grammar.y:762
		{
			yyVAL.op = ast.Add
		}
	case 95:
		//line grammar.y:766
		{
			yyVAL.op = ast.Sub
		}
	case 96:
		//line grammar.y:770
		{
			yyVAL.op = ast.Mult
		}
	case 97:
		//line grammar.y:774
		{
			yyVAL.op = ast.Div
		}
	case 98:
		//line grammar.y:778
		{
			yyVAL.op = ast.Modulo
		}
	case 99:
		//line grammar.y:782
		{
			yyVAL.op = ast.BitAnd
		}
	case 100:
		//line grammar.y:786
		{
			yyVAL.op = ast.BitOr
		}
	case 101:
		//line grammar.y:790
		{
			yyVAL.op = ast.BitXor
		}
	case 102:
		//line grammar.y:794
		{
			yyVAL.op = ast.LShift
		}
	case 103:
		//line grammar.y:798
		{
			yyVAL.op = ast.RShift
		}
	case 104:
		//line grammar.y:802
		{
			yyVAL.op = ast.Pow
		}
	case 105:
		//line grammar.y:806
		{
			yyVAL.op = ast.FloorDiv
		}
	case 106:
		//line grammar.y:813
		{
			setCtxs(yylex, yyS[yypt-0].exprs, ast.Del)
			yyVAL.stmt = &ast.Delete{StmtBase: ast.StmtBase{yyVAL.pos}, Targets: yyS[yypt-0].exprs}
		}
	case 107:
		//line grammar.y:820
		{
			yyVAL.stmt = &ast.Pass{StmtBase: ast.StmtBase{yyVAL.pos}}
		}
	case 108:
		//line grammar.y:826
		{
			yyVAL.stmt = yyS[yypt-0].stmt
		}
	case 109:
		//line grammar.y:830
		{
			yyVAL.stmt = yyS[yypt-0].stmt
		}
	case 110:
		//line grammar.y:834
		{
			yyVAL.stmt = yyS[yypt-0].stmt
		}
	case 111:
		//line grammar.y:838
		{
			yyVAL.stmt = yyS[yypt-0].stmt
		}
	case 112:
		//line grammar.y:842
		{
			yyVAL.stmt = yyS[yypt-0].stmt
		}
	case 113:
		//line grammar.y:848
		{
			yyVAL.stmt = &ast.Break{StmtBase: ast.StmtBase{yyVAL.pos}}
		}
	case 114:
		//line grammar.y:854
		{
			yyVAL.stmt = &ast.Continue{StmtBase: ast.StmtBase{yyVAL.pos}}
		}
	case 115:
		//line grammar.y:860
		{
			yyVAL.stmt = &ast.Return{StmtBase: ast.StmtBase{yyVAL.pos}}
		}
	case 116:
		//line grammar.y:864
		{
			yyVAL.stmt = &ast.Return{StmtBase: ast.StmtBase{yyVAL.pos}, Value: yyS[yypt-0].expr}
		}
	case 117:
		//line grammar.y:870
		{
			yyVAL.stmt = &ast.ExprStmt{StmtBase: ast.StmtBase{yyVAL.pos}, Value: yyS[yypt-0].expr}
		}
	case 118:
		//line grammar.y:876
		{
			yyVAL.stmt = &ast.Raise{StmtBase: ast.StmtBase{yyVAL.pos}}
		}
	case 119:
		//line grammar.y:880
		{
			yyVAL.stmt = &ast.Raise{StmtBase: ast.StmtBase{yyVAL.pos}, Exc: yyS[yypt-0].expr}
		}
	case 120:
		//line grammar.y:884
		{
			yyVAL.stmt = &ast.Raise{StmtBase: ast.StmtBase{yyVAL.pos}, Exc: yyS[yypt-2].expr, Cause: yyS[yypt-0].expr}
		}
	case 121:
		//line grammar.y:890
		{
			yyVAL.stmt = yyS[yypt-0].stmt
		}
	case 122:
		//line grammar.y:894
		{
			yyVAL.stmt = yyS[yypt-0].stmt
		}
	case 123:
		//line grammar.y:900
		{
			yyVAL.stmt = &ast.Import{StmtBase: ast.StmtBase{yyVAL.pos}, Names: yyS[yypt-0].aliases}
		}
	case 124:
		//line grammar.y:907
		{
			yyVAL.level = 1
		}
	case 125:
		//line grammar.y:911
		{
			yyVAL.level = 3
		}
	case 126:
		//line grammar.y:917
		{
			yyVAL.level = yyS[yypt-0].level
		}
	case 127:
		//line grammar.y:921
		{
			yyVAL.level += yyS[yypt-0].level
		}
	case 128:
		//line grammar.y:927
		{
			yyVAL.level = 0
			yyVAL.str = yyS[yypt-0].str
		}
	case 129:
		//line grammar.y:932
		{
			yyVAL.level = yyS[yypt-1].level
			yyVAL.str = yyS[yypt-0].str
		}
	case 130:
		//line grammar.y:937
		{
			yyVAL.level = yyS[yypt-0].level
			yyVAL.str = ""
		}
	case 131:
		//line grammar.y:944
		{
			yyVAL.aliases = []*ast.Alias{&ast.Alias{Pos: yyVAL.pos, Name: ast.Identifier("*")}}
		}
	case 132:
		//line grammar.y:948
		{
			yyVAL.aliases = yyS[yypt-2].aliases
		}
	case 133:
		//line grammar.y:952
		{
			yyVAL.aliases = yyS[yypt-1].aliases
		}
	case 134:
		//line grammar.y:958
		{
			yyVAL.stmt = &ast.ImportFrom{StmtBase: ast.StmtBase{yyVAL.pos}, Module: ast.Identifier(yyS[yypt-2].str), Names: yyS[yypt-0].aliases, Level: yyS[yypt-2].level}
		}
	case 135:
		//line grammar.y:964
		{
			yyVAL.alias = &ast.Alias{Pos: yyVAL.pos, Name: ast.Identifier(yyS[yypt-0].str)}
		}
	case 136:
		//line grammar.y:968
		{
			yyVAL.alias = &ast.Alias{Pos: yyVAL.pos, Name: ast.Identifier(yyS[yypt-2].str), AsName: ast.Identifier(yyS[yypt-0].str)}
		}
	case 137:
		//line grammar.y:974
		{
			yyVAL.alias = &ast.Alias{Pos: yyVAL.pos, Name: ast.Identifier(yyS[yypt-0].str)}
		}
	case 138:
		//line grammar.y:978
		{
			yyVAL.alias = &ast.Alias{Pos: yyVAL.pos, Name: ast.Identifier(yyS[yypt-2].str), AsName: ast.Identifier(yyS[yypt-0].str)}
		}
	case 139:
		//line grammar.y:984
		{
			yyVAL.aliases = nil
			yyVAL.aliases = append(yyVAL.aliases, yyS[yypt-0].alias)
		}
	case 140:
		//line grammar.y:989
		{
			yyVAL.aliases = append(yyVAL.aliases, yyS[yypt-0].alias)
		}
	case 141:
		//line grammar.y:995
		{
			yyVAL.aliases = nil
			yyVAL.aliases = append(yyVAL.aliases, yyS[yypt-0].alias)
		}
	case 142:
		//line grammar.y:1000
		{
			yyVAL.aliases = append(yyVAL.aliases, yyS[yypt-0].alias)
		}
	case 143:
		//line grammar.y:1006
		{
			yyVAL.str = yyS[yypt-0].str
		}
	case 144:
		//line grammar.y:1010
		{
			yyVAL.str += "." + yyS[yypt-0].str
		}
	case 145:
		//line grammar.y:1016
		{
			yyVAL.identifiers = nil
			yyVAL.identifiers = append(yyVAL.identifiers, ast.Identifier(yyS[yypt-0].str))
		}
	case 146:
		//line grammar.y:1021
		{
			yyVAL.identifiers = append(yyVAL.identifiers, ast.Identifier(yyS[yypt-0].str))
		}
	case 147:
		//line grammar.y:1027
		{
			yyVAL.stmt = &ast.Global{StmtBase: ast.StmtBase{yyVAL.pos}, Names: yyS[yypt-0].identifiers}
		}
	case 148:
		//line grammar.y:1033
		{
			yyVAL.stmt = &ast.Nonlocal{StmtBase: ast.StmtBase{yyVAL.pos}, Names: yyS[yypt-0].identifiers}
		}
	case 149:
		//line grammar.y:1039
		{
			yyVAL.exprs = nil
			yyVAL.exprs = append(yyVAL.exprs, yyS[yypt-0].expr)
		}
	case 150:
		//line grammar.y:1044
		{
			yyVAL.exprs = append(yyVAL.exprs, yyS[yypt-0].expr)
		}
	case 151:
		//line grammar.y:1050
		{
			yyVAL.stmt = &ast.Assert{StmtBase: ast.StmtBase{yyVAL.pos}, Test: yyS[yypt-0].expr}
		}
	case 152:
		//line grammar.y:1054
		{
			yyVAL.stmt = &ast.Assert{StmtBase: ast.StmtBase{yyVAL.pos}, Test: yyS[yypt-2].expr, Msg: yyS[yypt-0].expr}
		}
	case 153:
		//line grammar.y:1060
		{
			yyVAL.stmt = yyS[yypt-0].stmt
		}
	case 154:
		//line grammar.y:1064
		{
			yyVAL.stmt = yyS[yypt-0].stmt
		}
	case 155:
		//line grammar.y:1068
		{
			yyVAL.stmt = yyS[yypt-0].stmt
		}
	case 156:
		//line grammar.y:1072
		{
			yyVAL.stmt = yyS[yypt-0].stmt
		}
	case 157:
		//line grammar.y:1076
		{
			yyVAL.stmt = yyS[yypt-0].stmt
		}
	case 158:
		//line grammar.y:1080
		{
			yyVAL.stmt = yyS[yypt-0].stmt
		}
	case 159:
		//line grammar.y:1084
		{
			yyVAL.stmt = yyS[yypt-0].stmt
		}
	case 160:
		//line grammar.y:1088
		{
			yyVAL.stmt = yyS[yypt-0].stmt
		}
	case 161:
		//line grammar.y:1093
		{
			yyVAL.ifstmt = nil
			yyVAL.lastif = nil
		}
	case 162:
		//line grammar.y:1098
		{
			elifs := yyVAL.ifstmt
			newif := &ast.If{StmtBase: ast.StmtBase{yyVAL.pos}, Test: yyS[yypt-2].expr, Body: yyS[yypt-0].stmts}
			if elifs == nil {
				yyVAL.ifstmt = newif
			} else {
				yyVAL.lastif.Orelse = []ast.Stmt{newif}
			}
			yyVAL.lastif = newif
		}
	case 163:
		//line grammar.y:1110
		{
			yyVAL.stmts = nil
		}
	case 164:
		//line grammar.y:1114
		{
			yyVAL.stmts = yyS[yypt-0].stmts
		}
	case 165:
		//line grammar.y:1120
		{
			newif := &ast.If{StmtBase: ast.StmtBase{yyVAL.pos}, Test: yyS[yypt-4].expr, Body: yyS[yypt-2].stmts}
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
	case 166:
		//line grammar.y:1141
		{
			yyVAL.stmt = &ast.While{StmtBase: ast.StmtBase{yyVAL.pos}, Test: yyS[yypt-3].expr, Body: yyS[yypt-1].stmts, Orelse: yyS[yypt-0].stmts}
		}
	case 167:
		//line grammar.y:1147
		{
			target := tupleOrExpr(yyVAL.pos, yyS[yypt-5].exprs, false)
			setCtx(yylex, target, ast.Store)
			yyVAL.stmt = &ast.For{StmtBase: ast.StmtBase{yyVAL.pos}, Target: target, Iter: yyS[yypt-3].expr, Body: yyS[yypt-1].stmts, Orelse: yyS[yypt-0].stmts}
		}
	case 168:
		//line grammar.y:1154
		{
			yyVAL.exchandlers = nil
		}
	case 169:
		//line grammar.y:1158
		{
			exc := &ast.ExceptHandler{Pos: yyVAL.pos, ExprType: yyS[yypt-2].expr, Name: ast.Identifier(yyS[yypt-2].str), Body: yyS[yypt-0].stmts}
			yyVAL.exchandlers = append(yyVAL.exchandlers, exc)
		}
	case 170:
		//line grammar.y:1165
		{
			yyVAL.stmt = &ast.Try{StmtBase: ast.StmtBase{yyVAL.pos}, Body: yyS[yypt-1].stmts, Handlers: yyS[yypt-0].exchandlers}
		}
	case 171:
		//line grammar.y:1169
		{
			yyVAL.stmt = &ast.Try{StmtBase: ast.StmtBase{yyVAL.pos}, Body: yyS[yypt-4].stmts, Handlers: yyS[yypt-3].exchandlers, Orelse: yyS[yypt-0].stmts}
		}
	case 172:
		//line grammar.y:1173
		{
			yyVAL.stmt = &ast.Try{StmtBase: ast.StmtBase{yyVAL.pos}, Body: yyS[yypt-4].stmts, Handlers: yyS[yypt-3].exchandlers, Finalbody: yyS[yypt-0].stmts}
		}
	case 173:
		//line grammar.y:1177
		{
			yyVAL.stmt = &ast.Try{StmtBase: ast.StmtBase{yyVAL.pos}, Body: yyS[yypt-7].stmts, Handlers: yyS[yypt-6].exchandlers, Orelse: yyS[yypt-3].stmts, Finalbody: yyS[yypt-0].stmts}
		}
	case 174:
		//line grammar.y:1183
		{
			yyVAL.withitems = nil
			yyVAL.withitems = append(yyVAL.withitems, yyS[yypt-0].withitem)
		}
	case 175:
		//line grammar.y:1188
		{
			yyVAL.withitems = append(yyVAL.withitems, yyS[yypt-0].withitem)
		}
	case 176:
		//line grammar.y:1194
		{
			yyVAL.stmt = &ast.With{StmtBase: ast.StmtBase{yyVAL.pos}, Items: yyS[yypt-2].withitems, Body: yyS[yypt-0].stmts}
		}
	case 177:
		//line grammar.y:1200
		{
			yyVAL.withitem = &ast.WithItem{Pos: yyVAL.pos, ContextExpr: yyS[yypt-0].expr}
		}
	case 178:
		//line grammar.y:1204
		{
			v := yyS[yypt-0].expr
			setCtx(yylex, v, ast.Store)
			yyVAL.withitem = &ast.WithItem{Pos: yyVAL.pos, ContextExpr: yyS[yypt-2].expr, OptionalVars: v}
		}
	case 179:
		//line grammar.y:1213
		{
			yyVAL.expr = nil
			yyVAL.str = ""
		}
	case 180:
		//line grammar.y:1218
		{
			yyVAL.expr = yyS[yypt-0].expr
			yyVAL.str = ""
		}
	case 181:
		//line grammar.y:1223
		{
			yyVAL.expr = yyS[yypt-2].expr
			yyVAL.str = yyS[yypt-0].str
		}
	case 182:
		//line grammar.y:1230
		{
			yyVAL.stmts = nil
			yyVAL.stmts = append(yyVAL.stmts, yyS[yypt-0].stmts...)
		}
	case 183:
		//line grammar.y:1235
		{
			yyVAL.stmts = append(yyVAL.stmts, yyS[yypt-0].stmts...)
		}
	case 184:
		//line grammar.y:1241
		{
			yyVAL.stmts = yyS[yypt-0].stmts
		}
	case 185:
		//line grammar.y:1245
		{
			yyVAL.stmts = yyS[yypt-1].stmts
		}
	case 186:
		//line grammar.y:1251
		{
			yyVAL.expr = yyS[yypt-0].expr
		}
	case 187:
		//line grammar.y:1255
		{
			yyVAL.expr = &ast.IfExp{ExprBase: ast.ExprBase{yyVAL.pos}, Test: yyS[yypt-2].expr, Body: yyS[yypt-4].expr, Orelse: yyS[yypt-0].expr}
		}
	case 188:
		//line grammar.y:1259
		{
			yyVAL.expr = yyS[yypt-0].expr
		}
	case 189:
		//line grammar.y:1265
		{
			yyVAL.expr = yyS[yypt-0].expr
		}
	case 190:
		//line grammar.y:1269
		{
			yyVAL.expr = yyS[yypt-0].expr
		}
	case 191:
		//line grammar.y:1275
		{
			args := &ast.Arguments{Pos: yyVAL.pos}
			yyVAL.expr = &ast.Lambda{ExprBase: ast.ExprBase{yyVAL.pos}, Args: args, Body: yyS[yypt-0].expr}
		}
	case 192:
		//line grammar.y:1280
		{
			yyVAL.expr = &ast.Lambda{ExprBase: ast.ExprBase{yyVAL.pos}, Args: yyS[yypt-2].arguments, Body: yyS[yypt-0].expr}
		}
	case 193:
		//line grammar.y:1286
		{
			args := &ast.Arguments{Pos: yyVAL.pos}
			yyVAL.expr = &ast.Lambda{ExprBase: ast.ExprBase{yyVAL.pos}, Args: args, Body: yyS[yypt-0].expr}
		}
	case 194:
		//line grammar.y:1291
		{
			yyVAL.expr = &ast.Lambda{ExprBase: ast.ExprBase{yyVAL.pos}, Args: yyS[yypt-2].arguments, Body: yyS[yypt-0].expr}
		}
	case 195:
		//line grammar.y:1297
		{
			yyVAL.expr = yyS[yypt-0].expr
			yyVAL.isExpr = true
		}
	case 196:
		//line grammar.y:1302
		{
			if !yyS[yypt-2].isExpr {
				boolop := yyVAL.expr.(*ast.BoolOp)
				boolop.Values = append(boolop.Values, yyS[yypt-0].expr)
			} else {
				yyVAL.expr = &ast.BoolOp{ExprBase: ast.ExprBase{yyVAL.pos}, Op: ast.Or, Values: []ast.Expr{yyVAL.expr, yyS[yypt-0].expr}}
			}
			yyVAL.isExpr = false
		}
	case 197:
		//line grammar.y:1314
		{
			yyVAL.expr = yyS[yypt-0].expr
			yyVAL.isExpr = true
		}
	case 198:
		//line grammar.y:1319
		{
			if !yyS[yypt-2].isExpr {
				boolop := yyVAL.expr.(*ast.BoolOp)
				boolop.Values = append(boolop.Values, yyS[yypt-0].expr)
			} else {
				yyVAL.expr = &ast.BoolOp{ExprBase: ast.ExprBase{yyVAL.pos}, Op: ast.And, Values: []ast.Expr{yyVAL.expr, yyS[yypt-0].expr}}
			}
			yyVAL.isExpr = false
		}
	case 199:
		//line grammar.y:1331
		{
			yyVAL.expr = &ast.UnaryOp{ExprBase: ast.ExprBase{yyVAL.pos}, Op: ast.Not, Operand: yyS[yypt-0].expr}
		}
	case 200:
		//line grammar.y:1335
		{
			yyVAL.expr = yyS[yypt-0].expr
		}
	case 201:
		//line grammar.y:1341
		{
			yyVAL.expr = yyS[yypt-0].expr
			yyVAL.isExpr = true
		}
	case 202:
		//line grammar.y:1346
		{
			if !yyS[yypt-2].isExpr {
				comp := yyVAL.expr.(*ast.Compare)
				comp.Ops = append(comp.Ops, yyS[yypt-1].cmpop)
				comp.Comparators = append(comp.Comparators, yyS[yypt-0].expr)
			} else {
				yyVAL.expr = &ast.Compare{ExprBase: ast.ExprBase{yyVAL.pos}, Left: yyVAL.expr, Ops: []ast.CmpOp{yyS[yypt-1].cmpop}, Comparators: []ast.Expr{yyS[yypt-0].expr}}
			}
			yyVAL.isExpr = false
		}
	case 203:
		//line grammar.y:1361
		{
			yyVAL.cmpop = ast.Lt
		}
	case 204:
		//line grammar.y:1365
		{
			yyVAL.cmpop = ast.Gt
		}
	case 205:
		//line grammar.y:1369
		{
			yyVAL.cmpop = ast.Eq
		}
	case 206:
		//line grammar.y:1373
		{
			yyVAL.cmpop = ast.GtE
		}
	case 207:
		//line grammar.y:1377
		{
			yyVAL.cmpop = ast.LtE
		}
	case 208:
		//line grammar.y:1381
		{
			yylex.(*yyLex).SyntaxError("invalid syntax")
		}
	case 209:
		//line grammar.y:1385
		{
			yyVAL.cmpop = ast.NotEq
		}
	case 210:
		//line grammar.y:1389
		{
			yyVAL.cmpop = ast.In
		}
	case 211:
		//line grammar.y:1393
		{
			yyVAL.cmpop = ast.NotIn
		}
	case 212:
		//line grammar.y:1397
		{
			yyVAL.cmpop = ast.Is
		}
	case 213:
		//line grammar.y:1401
		{
			yyVAL.cmpop = ast.IsNot
		}
	case 214:
		//line grammar.y:1407
		{
			yyVAL.expr = &ast.Starred{ExprBase: ast.ExprBase{yyVAL.pos}, Value: yyS[yypt-0].expr, Ctx: ast.Load}
		}
	case 215:
		//line grammar.y:1413
		{
			yyVAL.expr = yyS[yypt-0].expr
		}
	case 216:
		//line grammar.y:1417
		{
			yyVAL.expr = &ast.BinOp{ExprBase: ast.ExprBase{yyVAL.pos}, Left: yyS[yypt-2].expr, Op: ast.BitOr, Right: yyS[yypt-0].expr}
		}
	case 217:
		//line grammar.y:1423
		{
			yyVAL.expr = yyS[yypt-0].expr
		}
	case 218:
		//line grammar.y:1427
		{
			yyVAL.expr = &ast.BinOp{ExprBase: ast.ExprBase{yyVAL.pos}, Left: yyS[yypt-2].expr, Op: ast.BitXor, Right: yyS[yypt-0].expr}
		}
	case 219:
		//line grammar.y:1433
		{
			yyVAL.expr = yyS[yypt-0].expr
		}
	case 220:
		//line grammar.y:1437
		{
			yyVAL.expr = &ast.BinOp{ExprBase: ast.ExprBase{yyVAL.pos}, Left: yyS[yypt-2].expr, Op: ast.BitAnd, Right: yyS[yypt-0].expr}
		}
	case 221:
		//line grammar.y:1443
		{
			yyVAL.expr = yyS[yypt-0].expr
		}
	case 222:
		//line grammar.y:1447
		{
			yyVAL.expr = &ast.BinOp{ExprBase: ast.ExprBase{yyVAL.pos}, Left: yyS[yypt-2].expr, Op: ast.LShift, Right: yyS[yypt-0].expr}
		}
	case 223:
		//line grammar.y:1451
		{
			yyVAL.expr = &ast.BinOp{ExprBase: ast.ExprBase{yyVAL.pos}, Left: yyS[yypt-2].expr, Op: ast.RShift, Right: yyS[yypt-0].expr}
		}
	case 224:
		//line grammar.y:1457
		{
			yyVAL.expr = yyS[yypt-0].expr
		}
	case 225:
		//line grammar.y:1461
		{
			yyVAL.expr = &ast.BinOp{ExprBase: ast.ExprBase{yyVAL.pos}, Left: yyS[yypt-2].expr, Op: ast.Add, Right: yyS[yypt-0].expr}
		}
	case 226:
		//line grammar.y:1465
		{
			yyVAL.expr = &ast.BinOp{ExprBase: ast.ExprBase{yyVAL.pos}, Left: yyS[yypt-2].expr, Op: ast.Sub, Right: yyS[yypt-0].expr}
		}
	case 227:
		//line grammar.y:1471
		{
			yyVAL.expr = yyS[yypt-0].expr
		}
	case 228:
		//line grammar.y:1475
		{
			yyVAL.expr = &ast.BinOp{ExprBase: ast.ExprBase{yyVAL.pos}, Left: yyS[yypt-2].expr, Op: ast.Mult, Right: yyS[yypt-0].expr}
		}
	case 229:
		//line grammar.y:1479
		{
			yyVAL.expr = &ast.BinOp{ExprBase: ast.ExprBase{yyVAL.pos}, Left: yyS[yypt-2].expr, Op: ast.Div, Right: yyS[yypt-0].expr}
		}
	case 230:
		//line grammar.y:1483
		{
			yyVAL.expr = &ast.BinOp{ExprBase: ast.ExprBase{yyVAL.pos}, Left: yyS[yypt-2].expr, Op: ast.Modulo, Right: yyS[yypt-0].expr}
		}
	case 231:
		//line grammar.y:1487
		{
			yyVAL.expr = &ast.BinOp{ExprBase: ast.ExprBase{yyVAL.pos}, Left: yyS[yypt-2].expr, Op: ast.FloorDiv, Right: yyS[yypt-0].expr}
		}
	case 232:
		//line grammar.y:1493
		{
			yyVAL.expr = &ast.UnaryOp{ExprBase: ast.ExprBase{yyVAL.pos}, Op: ast.UAdd, Operand: yyS[yypt-0].expr}
		}
	case 233:
		//line grammar.y:1497
		{
			yyVAL.expr = &ast.UnaryOp{ExprBase: ast.ExprBase{yyVAL.pos}, Op: ast.USub, Operand: yyS[yypt-0].expr}
		}
	case 234:
		//line grammar.y:1501
		{
			yyVAL.expr = &ast.UnaryOp{ExprBase: ast.ExprBase{yyVAL.pos}, Op: ast.Invert, Operand: yyS[yypt-0].expr}
		}
	case 235:
		//line grammar.y:1505
		{
			yyVAL.expr = yyS[yypt-0].expr
		}
	case 236:
		//line grammar.y:1511
		{
			yyVAL.expr = applyTrailers(yyS[yypt-1].expr, yyS[yypt-0].exprs)
		}
	case 237:
		//line grammar.y:1515
		{
			yyVAL.expr = &ast.BinOp{ExprBase: ast.ExprBase{yyVAL.pos}, Left: applyTrailers(yyS[yypt-3].expr, yyS[yypt-2].exprs), Op: ast.Pow, Right: yyS[yypt-0].expr}
		}
	case 238:
		//line grammar.y:1521
		{
			yyVAL.exprs = nil
		}
	case 239:
		//line grammar.y:1525
		{
			yyVAL.exprs = append(yyVAL.exprs, yyS[yypt-0].expr)
		}
	case 240:
		//line grammar.y:1531
		{
			yyVAL.obj = yyS[yypt-0].obj
		}
	case 241:
		//line grammar.y:1535
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
	case 242:
		//line grammar.y:1556
		{
			yyVAL.expr = &ast.Tuple{ExprBase: ast.ExprBase{yyVAL.pos}, Ctx: ast.Load}
		}
	case 243:
		//line grammar.y:1560
		{
			yyVAL.expr = yyS[yypt-1].expr
		}
	case 244:
		//line grammar.y:1564
		{
			yyVAL.expr = &ast.GeneratorExp{ExprBase: ast.ExprBase{yyVAL.pos}, Elt: yyS[yypt-2].expr, Generators: yyS[yypt-1].comprehensions}
		}
	case 245:
		//line grammar.y:1568
		{
			yyVAL.expr = tupleOrExpr(yyVAL.pos, yyS[yypt-2].exprs, yyS[yypt-1].comma)
		}
	case 246:
		//line grammar.y:1572
		{
			yyVAL.expr = &ast.List{ExprBase: ast.ExprBase{yyVAL.pos}, Ctx: ast.Load}
		}
	case 247:
		//line grammar.y:1576
		{
			yyVAL.expr = &ast.ListComp{ExprBase: ast.ExprBase{yyVAL.pos}, Elt: yyS[yypt-2].expr, Generators: yyS[yypt-1].comprehensions}
		}
	case 248:
		//line grammar.y:1580
		{
			yyVAL.expr = &ast.List{ExprBase: ast.ExprBase{yyVAL.pos}, Elts: yyS[yypt-2].exprs, Ctx: ast.Load}
		}
	case 249:
		//line grammar.y:1584
		{
			yyVAL.expr = &ast.Dict{ExprBase: ast.ExprBase{yyVAL.pos}}
		}
	case 250:
		//line grammar.y:1588
		{
			yyVAL.expr = yyS[yypt-1].expr
		}
	case 251:
		//line grammar.y:1592
		{
			yyVAL.expr = &ast.Name{ExprBase: ast.ExprBase{yyVAL.pos}, Id: ast.Identifier(yyS[yypt-0].str), Ctx: ast.Load}
		}
	case 252:
		//line grammar.y:1596
		{
			yyVAL.expr = &ast.Num{ExprBase: ast.ExprBase{yyVAL.pos}, N: yyS[yypt-0].obj}
		}
	case 253:
		//line grammar.y:1600
		{
			switch s := yyS[yypt-0].obj.(type) {
			case py.String:
				yyVAL.expr = &ast.Str{ExprBase: ast.ExprBase{yyVAL.pos}, S: s}
			case py.Bytes:
				yyVAL.expr = &ast.Bytes{ExprBase: ast.ExprBase{yyVAL.pos}, S: s}
			default:
				panic("not Bytes or String in strings")
			}
		}
	case 254:
		//line grammar.y:1611
		{
			yyVAL.expr = &ast.Ellipsis{ExprBase: ast.ExprBase{yyVAL.pos}}
		}
	case 255:
		//line grammar.y:1615
		{
			yyVAL.expr = &ast.NameConstant{ExprBase: ast.ExprBase{yyVAL.pos}, Value: py.None}
		}
	case 256:
		//line grammar.y:1619
		{
			yyVAL.expr = &ast.NameConstant{ExprBase: ast.ExprBase{yyVAL.pos}, Value: py.True}
		}
	case 257:
		//line grammar.y:1623
		{
			yyVAL.expr = &ast.NameConstant{ExprBase: ast.ExprBase{yyVAL.pos}, Value: py.False}
		}
	case 258:
		//line grammar.y:1630
		{
			yyVAL.expr = &ast.Call{ExprBase: ast.ExprBase{yyVAL.pos}}
		}
	case 259:
		//line grammar.y:1634
		{
			yyVAL.expr = yyS[yypt-1].call
		}
	case 260:
		//line grammar.y:1638
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
				slice = &ast.Index{SliceBase: extslice.SliceBase, Value: &ast.Tuple{ExprBase: ast.ExprBase{extslice.SliceBase.Pos}, Elts: elts, Ctx: ast.Load}}
			notAllIndex:
			}
			yyVAL.expr = &ast.Subscript{ExprBase: ast.ExprBase{yyVAL.pos}, Slice: slice, Ctx: ast.Load}
		}
	case 261:
		//line grammar.y:1656
		{
			yyVAL.expr = &ast.Attribute{ExprBase: ast.ExprBase{yyVAL.pos}, Attr: ast.Identifier(yyS[yypt-0].str), Ctx: ast.Load}
		}
	case 262:
		//line grammar.y:1662
		{
			yyVAL.slice = yyS[yypt-0].slice
			yyVAL.isExpr = true
		}
	case 263:
		//line grammar.y:1667
		{
			if !yyS[yypt-2].isExpr {
				extSlice := yyVAL.slice.(*ast.ExtSlice)
				extSlice.Dims = append(extSlice.Dims, yyS[yypt-0].slice)
			} else {
				yyVAL.slice = &ast.ExtSlice{SliceBase: ast.SliceBase{yyVAL.pos}, Dims: []ast.Slicer{yyS[yypt-2].slice, yyS[yypt-0].slice}}
			}
			yyVAL.isExpr = false
		}
	case 264:
		//line grammar.y:1679
		{
			if yyS[yypt-0].comma && yyS[yypt-1].isExpr {
				yyVAL.slice = &ast.ExtSlice{SliceBase: ast.SliceBase{yyVAL.pos}, Dims: []ast.Slicer{yyS[yypt-1].slice}}
			} else {
				yyVAL.slice = yyS[yypt-1].slice
			}
		}
	case 265:
		//line grammar.y:1689
		{
			yyVAL.slice = &ast.Index{SliceBase: ast.SliceBase{yyVAL.pos}, Value: yyS[yypt-0].expr}
		}
	case 266:
		//line grammar.y:1693
		{
			yyVAL.slice = &ast.Slice{SliceBase: ast.SliceBase{yyVAL.pos}, Lower: nil, Upper: nil, Step: nil}
		}
	case 267:
		//line grammar.y:1697
		{
			yyVAL.slice = &ast.Slice{SliceBase: ast.SliceBase{yyVAL.pos}, Lower: nil, Upper: nil, Step: yyS[yypt-0].expr}
		}
	case 268:
		//line grammar.y:1701
		{
			yyVAL.slice = &ast.Slice{SliceBase: ast.SliceBase{yyVAL.pos}, Lower: nil, Upper: yyS[yypt-0].expr, Step: nil}
		}
	case 269:
		//line grammar.y:1705
		{
			yyVAL.slice = &ast.Slice{SliceBase: ast.SliceBase{yyVAL.pos}, Lower: nil, Upper: yyS[yypt-1].expr, Step: yyS[yypt-0].expr}
		}
	case 270:
		//line grammar.y:1709
		{
			yyVAL.slice = &ast.Slice{SliceBase: ast.SliceBase{yyVAL.pos}, Lower: yyS[yypt-1].expr, Upper: nil, Step: nil}
		}
	case 271:
		//line grammar.y:1713
		{
			yyVAL.slice = &ast.Slice{SliceBase: ast.SliceBase{yyVAL.pos}, Lower: yyS[yypt-2].expr, Upper: nil, Step: yyS[yypt-0].expr}
		}
	case 272:
		//line grammar.y:1717
		{
			yyVAL.slice = &ast.Slice{SliceBase: ast.SliceBase{yyVAL.pos}, Lower: yyS[yypt-2].expr, Upper: yyS[yypt-0].expr, Step: nil}
		}
	case 273:
		//line grammar.y:1721
		{
			yyVAL.slice = &ast.Slice{SliceBase: ast.SliceBase{yyVAL.pos}, Lower: yyS[yypt-3].expr, Upper: yyS[yypt-1].expr, Step: yyS[yypt-0].expr}
		}
	case 274:
		//line grammar.y:1727
		{
			yyVAL.expr = nil
		}
	case 275:
		//line grammar.y:1731
		{
			yyVAL.expr = yyS[yypt-0].expr
		}
	case 276:
		//line grammar.y:1737
		{
			yyVAL.expr = yyS[yypt-0].expr
		}
	case 277:
		//line grammar.y:1741
		{
			yyVAL.expr = yyS[yypt-0].expr
		}
	case 278:
		//line grammar.y:1747
		{
			yyVAL.exprs = nil
			yyVAL.exprs = append(yyVAL.exprs, yyS[yypt-0].expr)
		}
	case 279:
		//line grammar.y:1752
		{
			yyVAL.exprs = append(yyVAL.exprs, yyS[yypt-0].expr)
		}
	case 280:
		//line grammar.y:1758
		{
			yyVAL.exprs = yyS[yypt-1].exprs
			yyVAL.comma = yyS[yypt-0].comma
		}
	case 281:
		//line grammar.y:1765
		{
			elts := yyS[yypt-1].exprs
			if yyS[yypt-0].comma || len(elts) > 1 {
				yyVAL.expr = &ast.Tuple{ExprBase: ast.ExprBase{yyVAL.pos}, Elts: elts, Ctx: ast.Load}
			} else {
				yyVAL.expr = elts[0]
			}
		}
	case 282:
		//line grammar.y:1776
		{
			yyVAL.exprs = yyS[yypt-1].exprs
		}
	case 283:
		//line grammar.y:1783
		{
			yyVAL.exprs = nil
			yyVAL.exprs = append(yyVAL.exprs, yyS[yypt-2].expr, yyS[yypt-0].expr) // key, value order
		}
	case 284:
		//line grammar.y:1788
		{
			yyVAL.exprs = append(yyVAL.exprs, yyS[yypt-2].expr, yyS[yypt-0].expr)
		}
	case 285:
		//line grammar.y:1794
		{
			keyValues := yyS[yypt-1].exprs
			d := &ast.Dict{ExprBase: ast.ExprBase{yyVAL.pos}, Keys: nil, Values: nil}
			for i := 0; i < len(keyValues)-1; i += 2 {
				d.Keys = append(d.Keys, keyValues[i])
				d.Values = append(d.Values, keyValues[i+1])
			}
			yyVAL.expr = d
		}
	case 286:
		//line grammar.y:1804
		{
			yyVAL.expr = &ast.DictComp{ExprBase: ast.ExprBase{yyVAL.pos}, Key: yyS[yypt-3].expr, Value: yyS[yypt-1].expr, Generators: yyS[yypt-0].comprehensions}
		}
	case 287:
		//line grammar.y:1808
		{
			yyVAL.expr = &ast.Set{ExprBase: ast.ExprBase{yyVAL.pos}, Elts: yyS[yypt-0].exprs}
		}
	case 288:
		//line grammar.y:1812
		{
			yyVAL.expr = &ast.SetComp{ExprBase: ast.ExprBase{yyVAL.pos}, Elt: yyS[yypt-1].expr, Generators: yyS[yypt-0].comprehensions}
		}
	case 289:
		//line grammar.y:1818
		{
			classDef := &ast.ClassDef{StmtBase: ast.StmtBase{yyVAL.pos}, Name: ast.Identifier(yyS[yypt-3].str), Body: yyS[yypt-0].stmts}
			yyVAL.stmt = classDef
			args := yyS[yypt-2].call
			if args != nil {
				classDef.Bases = args.Args
				classDef.Keywords = args.Keywords
				classDef.Starargs = args.Starargs
				classDef.Kwargs = args.Kwargs
			}
		}
	case 290:
		//line grammar.y:1832
		{
			yyVAL.call = yyS[yypt-0].call
		}
	case 291:
		//line grammar.y:1836
		{
			yyVAL.call.Args = append(yyVAL.call.Args, yyS[yypt-0].call.Args...)
			yyVAL.call.Keywords = append(yyVAL.call.Keywords, yyS[yypt-0].call.Keywords...)
		}
	case 292:
		//line grammar.y:1842
		{
			yyVAL.call = &ast.Call{}
		}
	case 293:
		//line grammar.y:1846
		{
			yyVAL.call = yyS[yypt-1].call
		}
	case 294:
		//line grammar.y:1851
		{
			yyVAL.call = &ast.Call{}
		}
	case 295:
		//line grammar.y:1855
		{
			yyVAL.call.Args = append(yyVAL.call.Args, yyS[yypt-0].call.Args...)
			yyVAL.call.Keywords = append(yyVAL.call.Keywords, yyS[yypt-0].call.Keywords...)
		}
	case 296:
		//line grammar.y:1862
		{
			yyVAL.call = yyS[yypt-1].call
		}
	case 297:
		//line grammar.y:1866
		{
			call := yyS[yypt-3].call
			call.Starargs = yyS[yypt-1].expr
			if len(yyS[yypt-0].call.Args) != 0 {
				yylex.(*yyLex).SyntaxError("only named arguments may follow *expression")
			}
			call.Keywords = append(call.Keywords, yyS[yypt-0].call.Keywords...)
			yyVAL.call = call
		}
	case 298:
		//line grammar.y:1876
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
	case 299:
		//line grammar.y:1887
		{
			call := yyS[yypt-2].call
			call.Kwargs = yyS[yypt-0].expr
			yyVAL.call = call
		}
	case 300:
		//line grammar.y:1897
		{
			yyVAL.call = &ast.Call{}
			yyVAL.call.Args = []ast.Expr{yyS[yypt-0].expr}
		}
	case 301:
		//line grammar.y:1902
		{
			yyVAL.call = &ast.Call{}
			yyVAL.call.Args = []ast.Expr{
				&ast.GeneratorExp{ExprBase: ast.ExprBase{yyVAL.pos}, Elt: yyS[yypt-1].expr, Generators: yyS[yypt-0].comprehensions},
			}
		}
	case 302:
		//line grammar.y:1909
		{
			yyVAL.call = &ast.Call{}
			test := yyS[yypt-2].expr
			if name, ok := test.(*ast.Name); ok {
				yyVAL.call.Keywords = []*ast.Keyword{&ast.Keyword{Pos: name.Pos, Arg: name.Id, Value: yyS[yypt-0].expr}}
			} else {
				yylex.(*yyLex).SyntaxError("keyword can't be an expression")
			}
		}
	case 303:
		//line grammar.y:1921
		{
			yyVAL.comprehensions = yyS[yypt-0].comprehensions
			yyVAL.exprs = nil
		}
	case 304:
		//line grammar.y:1926
		{
			yyVAL.comprehensions = yyS[yypt-0].comprehensions
			yyVAL.exprs = yyS[yypt-0].exprs
		}
	case 305:
		//line grammar.y:1933
		{
			c := ast.Comprehension{
				Target: tupleOrExpr(yyVAL.pos, yyS[yypt-2].exprs, yyS[yypt-2].comma),
				Iter:   yyS[yypt-0].expr,
			}
			setCtx(yylex, c.Target, ast.Store)
			yyVAL.comprehensions = []ast.Comprehension{c}
		}
	case 306:
		//line grammar.y:1942
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
	case 307:
		//line grammar.y:1955
		{
			yyVAL.exprs = []ast.Expr{yyS[yypt-0].expr}
			yyVAL.comprehensions = nil
		}
	case 308:
		//line grammar.y:1960
		{
			yyVAL.exprs = []ast.Expr{yyS[yypt-1].expr}
			yyVAL.exprs = append(yyVAL.exprs, yyS[yypt-0].exprs...)
			yyVAL.comprehensions = yyS[yypt-0].comprehensions
		}
	case 309:
		//line grammar.y:1971
		{
			yyVAL.expr = &ast.Yield{ExprBase: ast.ExprBase{yyVAL.pos}}
		}
	case 310:
		//line grammar.y:1975
		{
			yyVAL.expr = &ast.YieldFrom{ExprBase: ast.ExprBase{yyVAL.pos}, Value: yyS[yypt-0].expr}
		}
	case 311:
		//line grammar.y:1979
		{
			yyVAL.expr = &ast.Yield{ExprBase: ast.ExprBase{yyVAL.pos}, Value: yyS[yypt-0].expr}
		}
	}
	goto yystack /* stack new state and value */
}
