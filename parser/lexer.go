// Copyright 2018 The go-python Authors.  All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package parser

// FIXME need to implement formfeed

// Lexer should count line numbers too!

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"log"
	"regexp"
	"strings"
	"unicode"

	"github.com/go-python/gpython/ast"
	"github.com/go-python/gpython/py"
)

// The parser expects the lexer to return 0 on EOF.  Give it a name
// for clarity.
const eof = 0

// Signal eof with Error
const eofError = -1

// Standard python definition of a tab
const tabSize = 8

// The parser uses the type <prefix>Lex as a lexer.  It must provide
// the methods Lex(*<prefix>SymType) int and Error(string).
type yyLex struct {
	reader        *bufio.Reader
	filename      string     // name of the file being read
	line          string     // current line being parsed
	lastLine      string     // last line that was parsed
	pos           ast.Pos    // current position within file
	yylval        *yySymType // last token
	eof           bool       // flag to show EOF was read
	error         bool       // set if an error has ocurred
	errorString   string     // the string of the error
	indentStack   []int      // indent stack to control INDENT / DEDENT tokens
	state         int        // current state of state machine
	currentIndent string     // whitespace at start of current line
	interactive   bool       // set if mode "single" reading interactive input
	exec          bool       // set if mode "exec" reading from file
	bracket       int        // number of open [ ]
	parenthesis   int        // number of open ( )
	brace         int        // number of open { }
	comment       bool       // true if we are processing a comment, inline or otherwise
	mod           ast.Mod    // output
	tokens        []int      // buffered tokens to output
}

// Create a new lexer
//
// The mode argument specifies what kind of code must be compiled; it
// can be 'exec' if source consists of a sequence of statements,
// 'eval' if it consists of a single expression, or 'single' if it
// consists of a single interactive statement
func NewLex(r io.Reader, filename string, mode string) (*yyLex, error) {
	x := &yyLex{
		reader:      bufio.NewReader(r),
		filename:    filename,
		indentStack: []int{0},
		state:       readString,
	}
	switch mode {
	case "exec":
		x.queue(FILE_INPUT)
		x.exec = true
	case "eval":
		x.queue(EVAL_INPUT)
	case "single":
		x.queue(SINGLE_INPUT)
		x.interactive = true
	default:
		return nil, py.ExceptionNewf(py.ValueError, "compile mode must be 'exec', 'eval' or 'single'")
	}
	return x, nil
}

// queue tokens for later return
func (x *yyLex) queue(tokens ...int) {
	x.tokens = append(x.tokens, tokens...)
}

// Return whether the token queue is empty
func (x *yyLex) queueEmpty() bool {
	return len(x.tokens) == 0
}

// dequeue a token for return
//
// panic if no token available
func (x *yyLex) dequeue() int {
	if x.queueEmpty() {
		panic("token queue empty")
	}
	token := x.tokens[0]
	x.tokens = x.tokens[1:]
	return token
}

// Refill line
func (x *yyLex) refill() {
	var err error
	x.line, err = x.reader.ReadString('\n')
	if yyDebug >= 2 {
		fmt.Printf("line = %q, err = %v\n", x.line, err)
	}
	x.pos.Lineno += 1
	x.pos.ColOffset = 0
	switch err {
	case nil:
	case io.EOF:
		x.eof = true
	default:
		x.eof = true
		x.SyntaxErrorf("Error reading input: %v", err)
	}
	// If this is exec input, add a newline to the end of the
	// string if there isn't one already.
	if x.eof && x.exec && len(x.line) > 0 && x.line[len(x.line)-1] != '\n' {
		x.line += "\n"
	}
	x.lastLine = x.line
}

// Finds the length of a space and tab seperated string
func countIndent(s string) int {
	if len(s) == 0 {
		return 0
	}
	// FIXME these rules don't actually implement the python3
	// lexing rules which state
	//
	// Indentation is rejected as inconsistent if a source file
	// mixes tabs and spaces in a way that makes the meaning
	// dependent on the worth of a tab in spaces; a TabError is
	// raised in that case
	indent := 0
	for _, c := range s {
		switch c {
		case ' ':
			indent++
		case '\t':
			// 012345678901234567
			// a       b
			//  a      b
			//   a     b
			//    a    b
			//     a   b
			//      a  b
			//       a b
			//        ab
			//         a       b
			indent += tabSize - (indent & (tabSize - 1))
		default:
			panic(py.ExceptionNewf(py.IndentationError, "unexpected indent"))
		}

	}
	return indent
}

var operators = map[string]int{
	// 1 Character operators
	"(": '(',
	")": ')',
	"[": '[',
	"]": ']',
	":": ':',
	",": ',',
	";": ';',
	"+": '+',
	"-": '-',
	"*": '*',
	"/": '/',
	"|": '|',
	"&": '&',
	"<": '<',
	">": '>',
	"=": '=',
	".": '.',
	"%": '%',
	"{": '{',
	"}": '}',
	"^": '^',
	"~": '~',
	"@": '@',
	"#": '#',

	// 2 Character operators
	"!=": PLINGEQ,
	"%=": PERCEQ,
	"&=": ANDEQ,
	"**": STARSTAR,
	"*=": STAREQ,
	"+=": PLUSEQ,
	"-=": MINUSEQ,
	"->": MINUSGT,
	"//": DIVDIV,
	"/=": DIVEQ,
	"<<": LTLT,
	"<=": LTEQ,
	"<>": LTGT,
	"==": EQEQ,
	">=": GTEQ,
	">>": GTGT,
	"^=": HATEQ,
	"|=": PIPEEQ,

	// 3 Character operators
	"**=": STARSTAREQ,
	"...": ELIPSIS,
	"//=": DIVDIVEQ,
	"<<=": LTLTEQ,
	">>=": GTGTEQ,
}

var tokens = map[string]int{
	// Reserved words
	"False":    FALSE,
	"None":     NONE,
	"True":     TRUE,
	"and":      AND,
	"as":       AS,
	"assert":   ASSERT,
	"break":    BREAK,
	"class":    CLASS,
	"continue": CONTINUE,
	"def":      DEF,
	"del":      DEL,
	"elif":     ELIF,
	"else":     ELSE,
	"except":   EXCEPT,
	"finally":  FINALLY,
	"for":      FOR,
	"from":     FROM,
	"global":   GLOBAL,
	"if":       IF,
	"import":   IMPORT,
	"in":       IN,
	"is":       IS,
	"lambda":   LAMBDA,
	"nonlocal": NONLOCAL,
	"not":      NOT,
	"or":       OR,
	"pass":     PASS,
	"raise":    RAISE,
	"return":   RETURN,
	"try":      TRY,
	"while":    WHILE,
	"with":     WITH,
	"yield":    YIELD,
}

var tokenToString map[int]string

// Make tokenToString map
func init() {
	tokenToString = make(map[int]string, len(operators)+len(tokens)+16)
	for k, v := range operators {
		tokenToString[v] = k
	}
	for k, v := range tokens {
		tokenToString[v] = k
	}
	tokenToString[eof] = "eof"
	tokenToString[eofError] = "eofError"
	tokenToString[NEWLINE] = "NEWLINE"
	tokenToString[ENDMARKER] = "ENDMARKER"
	tokenToString[NAME] = "NAME"
	tokenToString[INDENT] = "INDENT"
	tokenToString[DEDENT] = "DEDENT"
	tokenToString[STRING] = "STRING"
	tokenToString[NUMBER] = "NUMBER"
	tokenToString[FILE_INPUT] = "FILE_INPUT"
	tokenToString[SINGLE_INPUT] = "SINGLE_INPUT"
	tokenToString[EVAL_INPUT] = "EVAL_INPUT"
}

// True if there are any open brackets
func (x *yyLex) openBrackets() bool {
	return x.bracket != 0 || x.parenthesis != 0 || x.brace != 0
}

// States
const (
	readString = iota
	readIndent
	checkEmpty
	checkIndent
	parseTokens
	checkEof
	isEof
)

// A Token with value
type LexToken struct {
	token int
	value py.Object
	pos   ast.Pos
}

// Convert the yySymType and token into a LexToken
func newLexToken(token int, yylval *yySymType) (lt LexToken) {
	lt.token = token
	lt.pos = yylval.pos
	if token == NAME {
		lt.value = py.String(yylval.str)
	} else if token == STRING || token == NUMBER {
		lt.value = yylval.obj
	} else {
		lt.value = nil
	}
	return
}

// String a LexToken
func (lt *LexToken) String() string {
	name := tokenToString[lt.token]
	if lt.value == nil {
		return fmt.Sprintf("%q (%d) %d:%d", name, lt.token, lt.pos.Lineno, lt.pos.ColOffset)
	}
	return fmt.Sprintf("%q (%d) = %T{%v} %d:%d", name, lt.token, lt.value, lt.value, lt.pos.Lineno, lt.pos.ColOffset)
}

// TokenToName returns the string name of a given token
func (lt *LexToken) TokenToName() string {
	name, ok := tokenToString[lt.token]
	if !ok {
		return ""
	}
	return name
}

// Token returns the int value of this token
func (lt *LexToken) Token() int {
	return lt.token
}

// Value returns the py.Object
func (lt *LexToken) Value() py.Object {
	return lt.value
}

// Pos returns the ast.Pos
func (lt *LexToken) Pos() ast.Pos {
	return lt.pos
}

// An slice of LexToken~s
type LexTokens []LexToken

// Compare two LexTokens
func (as LexTokens) Eq(bs []LexToken) bool {
	if len(as) != len(bs) {
		return false
	}
	for i := range as {
		a := as[i]
		b := bs[i]
		if a != b {
			return false
		}
	}
	return true
}

// String a LexTokens
func (lts LexTokens) String() string {
	buf := new(bytes.Buffer)
	_, _ = buf.WriteString("[")
	for i := range lts {
		lt := lts[i]
		_, _ = buf.WriteString("{")
		_, _ = buf.WriteString(lt.String())
		_, _ = buf.WriteString("}, ")
	}
	_, _ = buf.WriteString("]")
	return buf.String()
}

// Queue any remaining DEDENTS
func (x *yyLex) queueDedents() {
	for i := len(x.indentStack) - 1; i >= 1; i-- {
		x.queue(DEDENT)
	}
	x.indentStack = x.indentStack[:1]
}

// The parser calls this method to get each new token.  This
// implementation returns operators and NUM.
func (x *yyLex) Lex(yylval *yySymType) (ret int) {
	// Clear out the yySymType on each token (copied from rsc's cc)
	*yylval = yySymType{}
	x.yylval = yylval
	if yyDebug >= 2 {
		defer func() {
			lt := newLexToken(ret, yylval)
			fmt.Printf("LEX> %s\n", lt.String())
		}()
	}

	// Return queued tokens if there are any
	if !x.queueEmpty() {
		yylval.pos = x.pos
		return x.dequeue()
	}

	for {
		yylval.pos = x.pos
		switch x.state {
		case readString:
			// Read x.line
			x.refill()
			yylval.pos = x.pos
			x.state++
			if x.line == "" && x.eof {
				x.state = checkEof
				// During interactive input of statements an entirely blank logical
				// line (i.e. one containing not even whitespace or a comment)
				// terminates a multi-line statement.
				if x.interactive && !x.openBrackets() {
					x.queueDedents()
					x.queue(NEWLINE)
					return x.dequeue()
				}
				continue
			}
		case readIndent:
			// Read the initial indent and get rid of it
			trimmed := strings.TrimLeft(x.line, " \t")
			removed := len(x.line) - len(trimmed)
			x.currentIndent = x.line[:removed]
			x.pos.ColOffset += removed
			x.line = trimmed
			x.state++
		case checkEmpty:
			despaced := strings.TrimSpace(x.line) // remove other whitespace other than " \t"
			// Ignore line if just white space or whitespace then comment
			// if despaced == "" || despaced[0] == '#' {
			if despaced == "" {
				x.state = checkEof
				continue
			}
			x.state++
		case checkIndent:
			x.state++
			// Don't output INDENT or DEDENT if brackets are open
			if x.openBrackets() {
				continue
			}
			// See if indent has changed and issue INDENT / DEDENT
			indent := countIndent(x.currentIndent)
			i := len(x.indentStack) - 1
			indentStackTop := x.indentStack[i]
			if indent == indentStackTop {
				continue
			} else if indent > indentStackTop {
				x.indentStack = append(x.indentStack, indent)
				yylval.pos.ColOffset = 0 // Indents start at 0
				return INDENT
			} else {
				for ; i >= 0; i-- {
					if x.indentStack[i] == indent {
						goto foundIndent
					}
					x.queue(DEDENT)
				}
				x.SyntaxError("Inconsistent indent")
				return eof
			foundIndent:
				x.indentStack = x.indentStack[:i+1]
				return x.dequeue()
			}
		case parseTokens:
			// Skip white space
			trimmed := strings.TrimLeft(x.line, " \t")
			x.pos.ColOffset += len(x.line) - len(trimmed)
			x.line = trimmed

			// Peek next word
			if len(x.line) == 0 {
				x.state = checkEof
				continue
			}

			// Check if newline or comment reached
			// if x.line[0] == '\n' || x.line[0] == '#' {
			if x.line[0] == '\n' {
				x.state = checkEof
				if x.comment {
					// Don't output NEWLINE if we entered into a comment
					// token earlier.  reset that state and emit an ENDMARKER
					x.comment = false
					return ENDMARKER
				}
				// Don't output NEWLINE if brackets are open
				if x.openBrackets() {
					continue
				}
				return NEWLINE
			}

			// Check if continuation character
			if x.line[0] == '\\' && (len(x.line) <= 1 || x.line[1] == '\n') {
				if x.eof {
					x.state = checkEof
					continue
				}
				x.refill()
				x.state = parseTokens
				continue
			}

			// Note start of token for parser
			yylval.pos = x.pos

			// Read a number if available
			token, value := x.readNumber()
			if token != eof {
				if token == eofError {
					return eof
				}
				yylval.obj = value
				return token
			}

			// Read a string if available
			token, value = x.readString()
			if token != eof {
				if token == eofError {
					return eof
				}
				yylval.obj = value
				return token
			}

			// Read a keyword or identifier if available
			token, str := x.readIdentifierOrKeyword()
			if token != eof {
				yylval.str = str
				return token
			}

			// Read an operator if available
			token = x.readOperator()
			if token != eof {
				// implement implicit line joining rules
				switch token {
				case '[':
					x.bracket++
				case ']':
					x.bracket--
				case '(':
					x.parenthesis++
				case ')':
					x.parenthesis--
				case '{':
					x.brace++
				case '}':
					x.brace--
				case '#':
					x.comment = true
				}
				return token
			}

			// Nothing we recognise found
			x.SyntaxError("invalid syntax")
			return eof
		case checkEof:
			if x.eof {
				x.queueDedents()
				// then return ENDMARKER
				x.state = isEof
				if !x.interactive {
					x.queue(ENDMARKER)
				}
				if x.queueEmpty() {
					continue
				}
				return x.dequeue()
			}
			x.state = readString
		case isEof:
			return eof
		default:
			panic("Bad state")
		}
	}
}

// Can this rune start an identifier?
//
// identifier: `xid_start` `xid_continue`*
// id_start: <all characters in general categories Lu, Ll, Lt, Lm, Lo, Nl, the underscore, and characters with the Other_ID_Start property>
// id_continue: <all characters in `id_start`, plus characters in the categories Mn, Mc, Nd, Pc and others with the Other_ID_Continue property>
// xid_start: <all characters in `id_start` whose NFKC normalization is in "id_start xid_continue*">
// xid_continue: <all characters in `id_continue` whose NFKC normalization is in "id_continue*">
func isIdentifierStart(c rune) bool {
	switch {
	case c >= 'a' && c <= 'z':
		return true
	case c >= 'A' && c <= 'Z':
		return true
	case c == '_':
		return true
	case c < 128:
		return false
	case unicode.In(c, unicode.Lu, unicode.Ll, unicode.Lt, unicode.Lm, unicode.Lo, unicode.Nl):
		return true
	}
	return false
}

// Can this rune continue an identifier?
func isIdentifierChar(c rune) bool {
	switch {
	case c >= 'a' && c <= 'z':
		return true
	case c >= 'A' && c <= 'Z':
		return true
	case c >= '0' && c <= '9':
		return true
	case c == '_':
		return true
	case c < 128:
		return false
	case unicode.In(c, unicode.Lu, unicode.Ll, unicode.Lt, unicode.Lm, unicode.Lo, unicode.Nl, unicode.Mn, unicode.Mc, unicode.Nd, unicode.Pc):
		return true
	}
	return false
}

// Cut i characters off the start of the line
//
// returns the characters cut
func (x *yyLex) cut(i int) (cut string) {
	cut, x.line = x.line[:i], x.line[i:]
	x.pos.ColOffset += i
	return cut
}

// Read an identifier
func (x *yyLex) readIdentifier() string {
	var i int
	var c rune
	for i, c = range x.line {
		if i == 0 {
			if !isIdentifierStart(c) {
				goto found
			}
		} else {
			if !isIdentifierChar(c) {
				goto found
			}
		}
	}
	i = len(x.line)
found:
	return x.cut(i)
}

// Read an identifier or keyword
func (x *yyLex) readIdentifierOrKeyword() (int, string) {
	identifier := x.readIdentifier()
	if identifier == "" {
		return eof, ""
	}
	token, ok := tokens[identifier]
	if ok {
		return token, identifier
	}
	return NAME, identifier
}

// Read operator - returns token or eof for not found
func (x *yyLex) readOperator() int {
	// Look for length 3, 2, 1 operators
	for i := 3; i >= 1; i-- {
		if len(x.line) >= i {
			op := x.line[:i]
			if tok, ok := operators[op]; ok {
				x.cut(i)
				return tok
			}
		}
	}
	return eof
}

const pointFloat = `([0-9]*\.[0-9]+|[0-9]+\.)`

var (
	decimalInteger        = regexp.MustCompile(`^[0-9]+[jJ]?`)
	illegalDecimalInteger = regexp.MustCompile(`^0[0-9]*[1-9][0-9]*$`)
	octalInteger          = regexp.MustCompile(`^0[oO][0-7]+`)
	hexInteger            = regexp.MustCompile(`^0[xX][0-9a-fA-F]+`)
	binaryInteger         = regexp.MustCompile(`^0[bB][01]+`)
	floatNumber           = regexp.MustCompile(`^(([0-9]+|` + pointFloat + `)[eE][+-]?[0-9]+|` + pointFloat + `)[jJ]?`)
)

// Read one of the many types of python number
//
// Returns eof for couldn't read number or eofError on a bad read
func (x *yyLex) readNumber() (token int, value py.Object) {
	// Quick check for this being a number
	if len(x.line) == 0 {
		return eof, nil
	}
	// Starts with a digit
	r0 := x.line[0]
	if '0' <= r0 && r0 <= '9' {
		goto isNumber
	}
	// Or starts with . then a digit
	if len(x.line) > 1 && r0 == '.' {
		if r1 := x.line[1]; '0' <= r1 && r1 <= '9' {
			goto isNumber
		}
	}
	return eof, nil

isNumber:
	var s string
	var err error
	if s = octalInteger.FindString(x.line); s != "" {
		value, err = py.IntFromString(s[2:], 8)
		if err != nil {
			panic(err)
		}
	} else if s = hexInteger.FindString(x.line); s != "" {
		value, err = py.IntFromString(s[2:], 16)
		if err != nil {
			panic(err)
		}
	} else if s = binaryInteger.FindString(x.line); s != "" {
		value, err = py.IntFromString(s[2:], 2)
		if err != nil {
			panic(err)
		}
	} else if s = floatNumber.FindString(x.line); s != "" {
		last := s[len(s)-1]
		imaginary := false
		toParse := s
		if last == 'j' || last == 'J' {
			imaginary = true
			toParse = s[:len(s)-1]
		}
		value, err = py.FloatFromString(toParse)
		if err != nil {
			panic(err)
		}
		if imaginary {
			value = py.Complex(complex(0, value.(py.Float)))
		}
	} else if s = decimalInteger.FindString(x.line); s != "" {
		last := s[len(s)-1]
		if last == 'j' || last == 'J' {
			toParse := s[:len(s)-1]
			value, err = py.FloatFromString(toParse)
			if err != nil {
				panic(err)
			}
			value = py.Complex(complex(0, value.(py.Float)))
		} else {
			// Discard numbers with leading 0 except all 0s
			if illegalDecimalInteger.FindString(s) != "" {
				// FIXME where is this error going in the grammar?
				x.SyntaxError("illegal decimal with leading zero")
				return eofError, nil
			}
			value, err = py.IntFromString(s, 10)
			if err != nil {
				panic(err)
			}

		}
	} else {
		panic("Unparsed number")
	}
	x.cut(len(s))
	token = NUMBER
	return
}

// Read one of the many types of python string
//
// May return eof to skip to next matcher, or eofError indicating there was a problem
func (x *yyLex) readString() (token int, value py.Object) {
	// Quick check for this being a string
	if len(x.line) == 0 {
		return eof, nil
	}
	r0 := x.line[0]
	r1 := byte(0)
	r2 := byte(0)
	if len(x.line) >= 2 {
		r1 = x.line[1]
		if len(x.line) >= 3 {
			r2 = x.line[2]
		}
	}

	rawString := false  // whether we are parsing a r"" string
	byteString := false // whether we are parsing a b"" string
	// u"" strings are just normal strings so we ignore that qualifier

	// Start of string
	if r0 == '\'' || r0 == '"' {
		goto found
	}
	// Or start of r"" u"" b""
	if (r0 == 'r' || r0 == 'R') && (r1 == '\'' || r1 == '"') {
		rawString = true
		x.cut(1)
		goto found
	}
	if (r0 == 'b' || r0 == 'B') && (r1 == '\'' || r1 == '"') {
		byteString = true
		x.cut(1)
		goto found
	}
	if (r0 == 'u' || r0 == 'U') && (r1 == '\'' || r1 == '"') {
		x.cut(1)
		goto found
	}
	// Or start of br"" Br"" bR"" BR"" rb"" rB"" Rb"" RB""
	if (r0 == 'r' || r0 == 'R') && (r1 == 'b' || r1 == 'B') && (r2 == '\'' || r2 == '"') {
		rawString = true
		byteString = true
		x.cut(2)
		goto found
	}
	if (r0 == 'b' || r0 == 'B') && (r1 == 'r' || r1 == 'R') && (r2 == '\'' || r2 == '"') {
		rawString = true
		byteString = true
		x.cut(2)
		goto found
	}
	return eof, nil
found:
	multiLineString := false
	stringEnd := ""

	// Use x.rawString and x.byteString flags
	// Parse "x" """x""" 'x' '''x'''
	if strings.HasPrefix(x.line, `"""`) {
		stringEnd = `"""`
		x.cut(3)
		multiLineString = true
	} else if strings.HasPrefix(x.line, `'''`) {
		stringEnd = `'''`
		x.cut(3)
		multiLineString = true
	} else if strings.HasPrefix(x.line, `"`) {
		stringEnd = `"`
		x.cut(1)
	} else if strings.HasPrefix(x.line, `'`) {
		stringEnd = `'`
		x.cut(1)
	} else {
		panic("Bad string start")
	}
	buf := new(bytes.Buffer)
	for {
		escape := false
		for i, c := range x.line {
			if escape {
				// Continuation line - remove \ then continue
				if c == '\n' {
					buf.Truncate(buf.Len() - 1)
					goto readMore
				}
				_, _ = buf.WriteRune(c)
				escape = false
			} else {
				if strings.HasPrefix(x.line[i:], stringEnd) {
					x.cut(i + len(stringEnd))
					goto foundEndOfString
				}
				if c == '\\' {
					escape = true
				}
				if !multiLineString && c == '\n' {
					break
				}
				_, _ = buf.WriteRune(c)
			}
		}
		if !multiLineString {
			x.SyntaxErrorf("EOL while scanning string literal")
			return eofError, nil
		}
	readMore:
		if x.eof {
			if multiLineString {
				x.SyntaxErrorf("EOF while scanning triple-quoted string literal")
			} else {
				x.SyntaxErrorf("EOL while scanning string literal")
			}
			return eofError, nil
		}
		x.refill()
	}
foundEndOfString:
	if !rawString {
		var err error
		buf, err = DecodeEscape(buf, byteString)
		if err != nil {
			x.SyntaxErrorf("Decode error: %v", err)
			return eofError, nil
		}
	}
	if byteString {
		return STRING, py.Bytes(buf.Bytes())
	}
	return STRING, py.String(buf.String())
}

// The parser calls this method on a parse error.
func (x *yyLex) Error(s string) {
	x.error = true
	if yyDebug >= 1 {
		log.Printf("Parse error: %s", s)
		log.Printf("Parse buffer %q", x.line)
		log.Printf("State %#v", x)
	}
}

// The parser calls this method on a parse error.
func (x *yyLex) SyntaxError(s string) {
	x.errorString = s
	x.Error(s)
}

// Call this to write formatted errors
func (x *yyLex) SyntaxErrorf(format string, a ...interface{}) {
	x.SyntaxError(fmt.Sprintf(format, a...))
}

// Returns an python error for the current yyLex
func (x *yyLex) ErrorReturn() error {
	if x.error {
		if x.errorString == "" {
			if x.eof && x.interactive {
				x.errorString = "unexpected EOF while parsing"
			} else {
				x.errorString = "invalid syntax"
			}
		}
		return py.ExceptionNewf(py.SyntaxError, "%s", x.errorString)
	}
	return nil
}

// Set the debug level 0 = off, 4 = max
func SetDebug(level int) {
	yyDebug = level
}

// Parse a file
func Parse(in io.Reader, filename string, mode string) (mod ast.Mod, err error) {
	lex, err := NewLex(in, filename, mode)
	if err != nil {
		return nil, err
	}
	defer func() {
		if r := recover(); r != nil {
			err = py.MakeSyntaxError(r, filename, lex.pos.Lineno, lex.pos.ColOffset, lex.lastLine)
		}
	}()
	yyParse(lex)
	err = lex.ErrorReturn()
	if err != nil {
		err = py.MakeSyntaxError(err, filename, lex.pos.Lineno, lex.pos.ColOffset, lex.lastLine)
	}
	return lex.mod, err
}

// Parse a string
func ParseString(in string, mode string) (ast.Ast, error) {
	return Parse(bytes.NewBufferString(in), "<string>", mode)
}

// Lex a file only, returning a sequence of tokens
func Lex(in io.Reader, filename string, mode string) (lts LexTokens, err error) {
	lex, err := NewLex(in, filename, mode)
	if err != nil {
		return nil, err
	}
	defer func() {
		if r := recover(); r != nil {
			err = py.MakeSyntaxError(r, filename, lex.pos.Lineno, lex.pos.ColOffset, lex.lastLine)
		}
	}()
	yylval := yySymType{}
	for {
		ret := lex.Lex(&yylval)
		if ret == eof {
			break
		}
		lt := newLexToken(ret, &yylval)
		lts = append(lts, lt)
	}
	err = lex.ErrorReturn()
	if err != nil {
		err = py.MakeSyntaxError(err, filename, lex.pos.Lineno, lex.pos.ColOffset, lex.lastLine)
	}
	return
}

// Lex a string
func LexString(in string, mode string) (lts LexTokens, err error) {
	return Lex(bytes.NewBufferString(in), "<string>", mode)
}
