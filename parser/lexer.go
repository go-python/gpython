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
	"strconv"
	"strings"
	"unicode"

	"github.com/ncw/gpython/py"
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
	line          string // current line being parsed
	eof           bool   // flag to show EOF was read
	error         bool   // set if an error has ocurred
	indentStack   []int  // indent stack to control INDENT / DEDENT tokens
	state         int    // current state of state machine
	currentIndent string // whitespace at start of current line
	interactive   bool   // set if reading interactive input
	bracket       int    // number of open [ ]
	parenthesis   int    // number of open ( )
	brace         int    // number of open { }
}

func NewLex(r io.Reader) *yyLex {
	x := &yyLex{
		reader:      bufio.NewReader(r),
		indentStack: []int{0},
		state:       readString,
	}
	return x
}

// Refill line
func (x *yyLex) refill() {
	var err error
	x.line, err = x.reader.ReadString('\n')
	if yyDebug >= 2 {
		fmt.Printf("line = %q, err = %v\n", x.line, err)
	}
	switch err {
	case nil:
	case io.EOF:
		x.eof = true
	default:
		x.eof = true
		x.Errorf("Error reading input: %v", err)
	}
}

// Finds the length of a space and tab seperated string
func (x *yyLex) countIndent(s string) int {
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
			panic("bad indent")
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

// The parser calls this method to get each new token.  This
// implementation returns operators and NUM.
func (x *yyLex) Lex(yylval *yySymType) (ret int) {
	if yyDebug >= 2 {
		defer func() {
			name := tokenToString[ret]
			if ret == NAME {
				fmt.Printf("LEX> %q (%d) = %q\n", name, ret, yylval.str)
			} else if ret == STRING || ret == NUMBER {
				fmt.Printf("LEX> %q (%d) = %T{%v}\n", name, ret, yylval.obj, yylval.obj)
			} else {
				fmt.Printf("LEX> %q (%d) \n", name, ret)
			}
		}()
	}

	for {
		switch x.state {
		case readString:
			// Read x.line
			x.refill()
			x.state++
			// an empty line while reading interactive input should return a NEWLINE
			if x.interactive && (x.line == "" || x.line == "\n") {
				// Don't output NEWLINE if brackets are open
				if x.openBrackets() {
					continue
				}
				return NEWLINE
			}
		case readIndent:
			// Read the initial indent and get rid of it
			trimmed := strings.TrimLeft(x.line, " \t")
			x.currentIndent = x.line[:len(x.line)-len(trimmed)]
			x.line = trimmed
			x.state++
		case checkEmpty:
			// Ignore line if just white space or whitespace then comment
			if x.line == "" || x.line == "\n" || x.line[0] == '#' {
				x.state = checkEof
				continue
			}
			x.state++
		case checkIndent:
			// Don't output INDENT or DEDENT if brackets are open
			if x.openBrackets() {
				x.state++
				continue
			}
			// See if indent has changed and issue INDENT / DEDENT
			indent := x.countIndent(x.currentIndent)
			indentStackTop := x.indentStack[len(x.indentStack)-1]
			switch {
			case indent > indentStackTop:
				x.indentStack = append(x.indentStack, indent)
				x.state++
				return INDENT
			case indent < indentStackTop:
				for i := len(x.indentStack) - 1; i >= 0; i-- {
					if x.indentStack[i] == indent {
						goto foundIndent
					}
				}
				x.Error("Inconsistent indent")
				return eof
			foundIndent:
				x.indentStack = x.indentStack[:len(x.indentStack)-1]
				return DEDENT
			}
			x.state++
		case parseTokens:
			// Skip white space
			x.line = strings.TrimLeft(x.line, " \t")

			// Peek next word
			if len(x.line) == 0 {
				x.state = checkEof
				continue
			}

			// Check if newline or comment reached
			if x.line[0] == '\n' || x.line[0] == '#' {
				x.state = checkEof
				// Don't output NEWLINE if brackets are open
				if x.openBrackets() {
					continue
				}
				return NEWLINE
			}

			// Check if continuation character
			if x.line[0] == '\\' && (len(x.line) <= 1 || x.line[1] == '\n') {
				if x.eof {
					return eof
				}
				x.refill()
				x.state = parseTokens
				continue
			}

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
				}
				return token
			}

			// Nothing we recognise found
			x.Error("Syntax error")
			return eof
		case checkEof:
			if x.eof {
				// Return any remaining DEDENTS
				if len(x.indentStack) > 1 {
					x.indentStack = x.indentStack[:len(x.indentStack)-1]
					x.state = checkEof
					return DEDENT
				}
				// then return ENDMARKER
				x.state = isEof
				return ENDMARKER
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
	identifier := x.line[:i]
	x.line = x.line[i:]
	return identifier
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
				x.line = x.line[i:]
				return tok
			}
		}
	}
	return eof
}

const pointFloat = `([0-9]*\.[0-9]+|[0-9]+\.)`

var decimalInteger = regexp.MustCompile(`^[0-9]+[jJ]?`)
var illegalDecimalInteger = regexp.MustCompile(`^0[0-9]*[1-9][0-9]*$`)
var octalInteger = regexp.MustCompile(`^0[oO][0-7]+`)
var hexInteger = regexp.MustCompile(`^0[xX][0-9a-fA-F]+`)
var binaryInteger = regexp.MustCompile(`^0[bB][01]+`)
var floatNumber = regexp.MustCompile(`^(([0-9]+|` + pointFloat + `)[eE][+-]?[0-9]+|` + pointFloat + `)[jJ]?`)

// Read one of the many types of python number
//
// Returns eof for couldn't read number or eofError on a bad read
func (x *yyLex) readNumber() (token int, value py.Object) {
	// Quick check for this being a number
	if len(x.line) == 0 {
		return eof, nil
	}
	r0 := x.line[0]
	r1 := byte(0)
	if len(x.line) > 1 {
		r1 = x.line[1]
	}
	// Check if could be a start of a number
	if !(('0' <= r0 && r0 <= '9') || (r0 == '.' && '0' <= r1 && r1 <= '9')) {
		return eof, nil
	}

	var s string
	if s = octalInteger.FindString(x.line); s != "" {
		value = py.IntNew(py.IntType, py.Tuple{py.String(s[2:]), py.Int(8)}, nil)
	} else if s = hexInteger.FindString(x.line); s != "" {
		value = py.IntNew(py.IntType, py.Tuple{py.String(s[2:]), py.Int(16)}, nil)
	} else if s = binaryInteger.FindString(x.line); s != "" {
		value = py.IntNew(py.IntType, py.Tuple{py.String(s[2:]), py.Int(2)}, nil)
	} else if s = floatNumber.FindString(x.line); s != "" {
		last := s[len(s)-1]
		imaginary := false
		toParse := s
		if last == 'j' || last == 'J' {
			imaginary = true
			toParse = s[:len(s)-1]
		}
		f, err := strconv.ParseFloat(toParse, 64)
		if err != nil {
			panic(py.ExceptionNewf(py.ValueError, "invalid literal for float: '%s' (%v)", toParse, err))
		}
		if imaginary {
			value = py.Complex(complex(0, f))
		} else {
			value = py.Float(f)
		}
	} else if s = decimalInteger.FindString(x.line); s != "" {
		last := s[len(s)-1]
		if last == 'j' || last == 'J' {
			toParse := s[:len(s)-1]
			f, err := strconv.ParseFloat(toParse, 64)
			if err != nil {
				panic(py.ExceptionNewf(py.ValueError, "invalid literal for imaginary number: '%s' (%v)", toParse, err))
			}
			value = py.Complex(complex(0, f))
		} else {
			// Discard numbers with leading 0 except all 0s
			if illegalDecimalInteger.FindString(x.line) != "" {
				return eofError, nil
			}
			value = py.IntNew(py.IntType, py.Tuple{py.String(s), py.Int(10)}, nil)
		}
	} else {
		return eof, nil
	}
	x.line = x.line[len(s):]
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
		x.line = x.line[1:]
		goto found
	}
	if (r0 == 'b' || r0 == 'B') && (r1 == '\'' || r1 == '"') {
		byteString = true
		x.line = x.line[1:]
		goto found
	}
	if (r0 == 'u' || r0 == 'U') && (r1 == '\'' || r1 == '"') {
		x.line = x.line[1:]
		goto found
	}
	// Or start of br"" Br"" bR"" BR"" rb"" rB"" Rb"" RB""
	if (r0 == 'r' || r0 == 'R') && (r1 == 'b' || r1 == 'B') && (r2 == '\'' || r2 == '"') {
		rawString = true
		byteString = true
		x.line = x.line[2:]
		goto found
	}
	if (r0 == 'b' || r0 == 'B') && (r1 == 'r' || r1 == 'R') && (r2 == '\'' || r2 == '"') {
		rawString = true
		byteString = true
		x.line = x.line[2:]
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
		x.line = x.line[3:]
		multiLineString = true
	} else if strings.HasPrefix(x.line, `'''`) {
		stringEnd = `'''`
		x.line = x.line[3:]
		multiLineString = true
	} else if strings.HasPrefix(x.line, `"`) {
		stringEnd = `"`
		x.line = x.line[1:]
	} else if strings.HasPrefix(x.line, `'`) {
		stringEnd = `'`
		x.line = x.line[1:]
	} else {
		x.Error("Bad string start")
		return eofError, nil
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
				buf.WriteRune(c)
				escape = false
			} else {
				if strings.HasPrefix(x.line[i:], stringEnd) {
					x.line = x.line[i+len(stringEnd):]
					goto foundEndOfString
				}
				if c == '\\' {
					escape = true
				}
				if !multiLineString && c == '\n' {
					break
				}
				buf.WriteRune(c)
			}
		}
		if !multiLineString {
			x.Errorf("Unterminated %sx%s string", stringEnd, stringEnd)
			return eofError, nil
		}
	readMore:
		if x.eof {
			x.Errorf("Unterminated %sx%s string", stringEnd, stringEnd)
			return eofError, nil
		}
		x.refill()
	}
foundEndOfString:
	if !rawString {
		// FIXME expand / sequences
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

// Call this to write formatted errors
func (x *yyLex) Errorf(format string, a ...interface{}) {
	x.Error(fmt.Sprintf(format, a...))
}

// Set the debug level 0 = off, 4 = max
func SetDebug(level int) {
	yyDebug = level
}

// Parse a file
func Parse(in io.Reader) error {
	lex := NewLex(in)
	yyParse(lex)
	if lex.error {
		return py.ExceptionNewf(py.SyntaxError, "Syntax Error")
	}
	return nil
}

// Lex a file only
func Lex(in io.Reader) error {
	lex := NewLex(in)
	yylval := yySymType{}
	for {
		ret := lex.Lex(&yylval)
		if ret == eof {
			break
		}
	}
	if lex.error {
		return py.ExceptionNewf(py.SyntaxError, "Syntax Error")
	}
	return nil
}
