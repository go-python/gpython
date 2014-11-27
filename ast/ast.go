package ast

// FIXME make base AstNode with position in
// also keep a list of children in the parent node to simplify walking the tree?

import (
	"fmt"

	"github.com/ncw/gpython/py"
)

type Identifier py.String
type String py.String
type Object py.Object
type Singleton py.Object

// Definitions originally made with python3 asdl_go.py Python.asdl then edited by hand
//
// Taking some inspiration from http://golang.org/src/pkg/go/ast/ast.go

// There are 5 main classes of Node
// ModBase - which type of thing are we parsing
// StmtBase - statements
// ExprType - expressions
// SliceBaseType - slices

// All node types implement the Ast interface
type Ast interface {
	py.Object
	GetLineno() int
	GetColOffset() int
}

// All ModBase nodes implement the Mod interface
type Mod interface {
	Ast
	modNode()
}

// All StmtBase nodes implement the Stmt interface
type Stmt interface {
	Ast
	stmtNode()
}

// All ExprBase notes implement the Expr interface
type Expr interface {
	Ast
	exprNode()
}

// All SliceBase nodes implement the Slicer interface
type Slicer interface {
	Ast
	sliceNode()
}

// Position in the parse tree
type Pos struct {
	Lineno    int
	ColOffset int
}

func (o *Pos) GetLineno() int    { return o.Lineno }
func (o *Pos) GetColOffset() int { return o.ColOffset }

// Base AST node
type AST struct {
	Pos
}

// ------------------------------------------------------------
// Constant
// ------------------------------------------------------------

type ExprContext int

const (
	Load = ExprContext(iota + 1)
	Store
	Del
	AugLoad
	AugStore
	Param
)

func (o ExprContext) String() string {
	switch o {
	case Load:
		return "Load()"
	case Store:
		return "Store()"
	case Del:
		return "Del()"
	case AugLoad:
		return "AugLoad()"
	case AugStore:
		return "AugStore()"
	case Param:
		return "Param()"
	}
	return fmt.Sprintf("UnknownExprContext(%d)", o)
}

type BoolOpNumber int

const (
	And = BoolOpNumber(iota + 1)
	Or
)

func (o BoolOpNumber) String() string {
	switch o {
	case And:
		return "And()"
	case Or:
		return "Or()"
	}
	return fmt.Sprintf("UnknownBoolOpNumber(%d)", o)
}

type OperatorNumber int

const (
	Add = OperatorNumber(iota + 1)
	Sub
	Mult
	Div
	Modulo
	Pow
	LShift
	RShift
	BitOr
	BitXor
	BitAnd
	FloorDiv
)

func (o OperatorNumber) String() string {
	switch o {
	case Add:
		return "Add()"
	case Sub:
		return "Sub()"
	case Mult:
		return "Mult()"
	case Div:
		return "Div()"
	case Modulo:
		return "Mod()"
	case Pow:
		return "Pow()"
	case LShift:
		return "LShift()"
	case RShift:
		return "RShift()"
	case BitOr:
		return "BitOr()"
	case BitXor:
		return "BitXor()"
	case BitAnd:
		return "BitAnd()"
	case FloorDiv:
		return "FloorDiv()"
	}
	return fmt.Sprintf("UnknownOperatorNumber(%d)", o)
}

type UnaryOpNumber int

const (
	Invert = UnaryOpNumber(iota + 1)
	Not
	UAdd
	USub
)

func (o UnaryOpNumber) String() string {
	switch o {
	case Invert:
		return "Invert()"
	case Not:
		return "Not()"
	case UAdd:
		return "UAdd()"
	case USub:
		return "USub()"
	}
	return fmt.Sprintf("UnknownUnaryOpNumber(%d)", o)
}

type CmpOp int

const (
	Eq = CmpOp(iota + 1)
	NotEq
	Lt
	LtE
	Gt
	GtE
	Is
	IsNot
	In
	NotIn
)

func (o CmpOp) String() string {
	switch o {
	case Eq:
		return "Eq()"
	case NotEq:
		return "NotEq()"
	case Lt:
		return "Lt()"
	case LtE:
		return "LtE()"
	case Gt:
		return "Gt()"
	case GtE:
		return "GtE()"
	case Is:
		return "Is()"
	case IsNot:
		return "IsNot()"
	case In:
		return "In()"
	case NotIn:
		return "NotIn()"
	}
	return fmt.Sprintf("UnknownCmoOp(%d)", o)
}

// ------------------------------------------------------------
// Mod nodes
// ------------------------------------------------------------

type ModBase struct {
	Pos
}

func (o *ModBase) modNode() {}

type Module struct {
	ModBase
	Body []Stmt
}

type Interactive struct {
	ModBase
	Body []Stmt
}

type Expression struct {
	ModBase
	Body Expr
}

type Suite struct {
	ModBase
	Body []Stmt
}

// ------------------------------------------------------------
// Statement nodes
// ------------------------------------------------------------

type StmtBase struct {
	Pos
}

func (o *StmtBase) stmtNode() {}

type FunctionDef struct {
	StmtBase
	Name          Identifier
	Args          *Arguments
	Body          []Stmt
	DecoratorList []Expr
	Returns       Expr
}

type ClassDef struct {
	StmtBase
	Name          Identifier
	Bases         []Expr
	Keywords      []*Keyword
	Starargs      Expr
	Kwargs        Expr
	Body          []Stmt
	DecoratorList []Expr
}

type Return struct {
	StmtBase
	Value Expr
}

type Delete struct {
	StmtBase
	Targets []Expr
}

type Assign struct {
	StmtBase
	Targets []Expr
	Value   Expr
}

type AugAssign struct {
	StmtBase
	Target Expr
	Op     OperatorNumber
	Value  Expr
}

type For struct {
	StmtBase
	Target Expr
	Iter   Expr
	Body   []Stmt
	Orelse []Stmt
}

type While struct {
	StmtBase
	Test   Expr
	Body   []Stmt
	Orelse []Stmt
}

type If struct {
	StmtBase
	Test   Expr
	Body   []Stmt
	Orelse []Stmt
}

type With struct {
	StmtBase
	Items []WithItem
	Body  []Stmt
}

type Raise struct {
	StmtBase
	Exc   Expr
	Cause Expr
}

type Try struct {
	StmtBase
	Body      []Stmt
	Handlers  []ExceptHandler
	Orelse    []Stmt
	Finalbody []Stmt
}

type Assert struct {
	StmtBase
	Test Expr
	Msg  Expr
}

type Import struct {
	StmtBase
	Names []*Alias
}

type ImportFrom struct {
	StmtBase
	Module Identifier
	Names  []*Alias
	Level  int
}

type Global struct {
	StmtBase
	Names []Identifier
}

type Nonlocal struct {
	StmtBase
	Names []Identifier
}

type ExprStmt struct {
	StmtBase
	Value Expr
}

type Pass struct {
	StmtBase
}

type Break struct {
	StmtBase
}

type Continue struct {
	StmtBase
}

// ------------------------------------------------------------
// Expr nodes
// ------------------------------------------------------------

type ExprBase struct{ Pos }

func (o *ExprBase) exprNode() {}

type BoolOp struct {
	ExprBase
	Op     BoolOpNumber
	Values []Expr
}

type BinOp struct {
	ExprBase
	Left  Expr
	Op    OperatorNumber
	Right Expr
}

type UnaryOp struct {
	ExprBase
	Op      UnaryOpNumber
	Operand Expr
}

type Lambda struct {
	ExprBase
	Args *Arguments
	Body Expr
}

type IfExp struct {
	ExprBase
	Test   Expr
	Body   Expr
	Orelse Expr
}

type Dict struct {
	ExprBase
	Keys   []Expr
	Values []Expr
}

type Set struct {
	ExprBase
	Elts []Expr
}

type ListComp struct {
	ExprBase
	Elt        Expr
	Generators []Comprehension
}

type SetComp struct {
	ExprBase
	Elt        Expr
	Generators []Comprehension
}

type DictComp struct {
	ExprBase
	Key        Expr
	Value      Expr
	Generators []Comprehension
}

type GeneratorExp struct {
	ExprBase
	Elt        Expr
	Generators []Comprehension
}

type Yield struct {
	ExprBase
	Value Expr
}

type YieldFrom struct {
	ExprBase
	Value Expr
}

// need sequences for compare to distinguish between
type Compare struct {
	ExprBase
	Left        Expr
	Ops         []CmpOp
	Comparators []Expr
}

type Call struct {
	ExprBase
	Func     Expr
	Args     []Expr
	Keywords []*Keyword
	Starargs Expr
	Kwargs   Expr
}

type Num struct {
	ExprBase
	N Object
}

type Str struct {
	ExprBase
	S py.String
}

type Bytes struct {
	ExprBase
	S py.Bytes
}

type NameConstant struct {
	ExprBase
	Value Singleton
}

type Ellipsis struct {
	ExprBase
}

type Attribute struct {
	ExprBase
	Value Expr
	Attr  Identifier
	Ctx   ExprContext
}

// ExprNodes which have a settable context implement this
type SetCtxer interface {
	SetCtx(ExprContext)
}

func (o *Attribute) SetCtx(Ctx ExprContext) { o.Ctx = Ctx }

var _ = SetCtxer((*Attribute)(nil))

type Subscript struct {
	ExprBase
	Value Expr
	Slice Slicer
	Ctx   ExprContext
}

func (o *Subscript) SetCtx(Ctx ExprContext) { o.Ctx = Ctx }

var _ = SetCtxer((*Subscript)(nil))

type Starred struct {
	ExprBase
	Value Expr
	Ctx   ExprContext
}

func (o *Starred) SetCtx(Ctx ExprContext) { o.Ctx = Ctx }

var _ = SetCtxer((*Starred)(nil))

type Name struct {
	ExprBase
	Id  Identifier
	Ctx ExprContext
}

func (o *Name) SetCtx(Ctx ExprContext) { o.Ctx = Ctx }

var _ = SetCtxer((*Name)(nil))

type List struct {
	ExprBase
	Elts []Expr
	Ctx  ExprContext
}

func (o *List) SetCtx(Ctx ExprContext) {
	o.Ctx = Ctx
	for i := range o.Elts {
		o.Elts[i].(SetCtxer).SetCtx(Ctx)
	}
}

var _ = SetCtxer((*List)(nil))

type Tuple struct {
	ExprBase
	Elts []Expr
	Ctx  ExprContext
}

func (o *Tuple) SetCtx(Ctx ExprContext) {
	o.Ctx = Ctx
	for i := range o.Elts {
		o.Elts[i].(SetCtxer).SetCtx(Ctx)
	}
}

var _ = SetCtxer((*Tuple)(nil))

// ------------------------------------------------------------
// Slicer nodes
// ------------------------------------------------------------

type SliceBase struct {
	Pos
}

func (o *SliceBase) sliceNode() {}

type Slice struct {
	SliceBase
	Lower Expr
	Upper Expr
	Step  Expr
}

type ExtSlice struct {
	SliceBase
	Dims []Slicer
}

type Index struct {
	SliceBase
	Value Expr
}

type Comprehension struct {
	Target Expr
	Iter   Expr
	Ifs    []Expr
}

// ------------------------------------------------------------
// Misc types - which aren't in a type heirachy
// ------------------------------------------------------------

type ExceptHandler struct {
	Pos
	Exprtype Expr
	Name     Identifier
	Body     []Stmt
}

type Arguments struct {
	Pos
	Args       []Arg
	Vararg     Arg
	Kwonlyargs []Arg
	KwDefaults []Expr
	Kwarg      Arg
	Defaults   []Expr
}

type Arg struct {
	Pos
	Arg        Identifier
	Annotation Expr
}

type Keyword struct {
	Pos
	Arg   Identifier
	Value Expr
}

type Alias struct {
	Pos
	Name   Identifier
	AsName *Identifier
}

type WithItem struct {
	Pos
	ContextExpr  Expr
	OptionalVars Expr
}

// Check interfaces

var _ Ast = (*AST)(nil)

// Mod
var _ Mod = (*ModBase)(nil)
var _ Mod = (*Module)(nil)
var _ Mod = (*Interactive)(nil)
var _ Mod = (*Expression)(nil)
var _ Mod = (*Suite)(nil)

// Stmt
var _ Stmt = (*StmtBase)(nil)
var _ Stmt = (*FunctionDef)(nil)
var _ Stmt = (*ClassDef)(nil)
var _ Stmt = (*Return)(nil)
var _ Stmt = (*Delete)(nil)
var _ Stmt = (*Assign)(nil)
var _ Stmt = (*AugAssign)(nil)
var _ Stmt = (*For)(nil)
var _ Stmt = (*While)(nil)
var _ Stmt = (*If)(nil)
var _ Stmt = (*With)(nil)
var _ Stmt = (*Raise)(nil)
var _ Stmt = (*Try)(nil)
var _ Stmt = (*Assert)(nil)
var _ Stmt = (*Import)(nil)
var _ Stmt = (*ImportFrom)(nil)
var _ Stmt = (*Global)(nil)
var _ Stmt = (*Nonlocal)(nil)
var _ Stmt = (*ExprStmt)(nil)
var _ Stmt = (*Pass)(nil)
var _ Stmt = (*Break)(nil)
var _ Stmt = (*Continue)(nil)

// Expr
var _ Expr = (*ExprBase)(nil)
var _ Expr = (*BoolOp)(nil)
var _ Expr = (*BinOp)(nil)
var _ Expr = (*UnaryOp)(nil)
var _ Expr = (*Lambda)(nil)
var _ Expr = (*IfExp)(nil)
var _ Expr = (*Dict)(nil)
var _ Expr = (*Set)(nil)
var _ Expr = (*ListComp)(nil)
var _ Expr = (*SetComp)(nil)
var _ Expr = (*DictComp)(nil)
var _ Expr = (*GeneratorExp)(nil)
var _ Expr = (*Yield)(nil)
var _ Expr = (*YieldFrom)(nil)
var _ Expr = (*Compare)(nil)
var _ Expr = (*Call)(nil)
var _ Expr = (*Num)(nil)
var _ Expr = (*Str)(nil)
var _ Expr = (*Bytes)(nil)
var _ Expr = (*NameConstant)(nil)
var _ Expr = (*Ellipsis)(nil)
var _ Expr = (*Attribute)(nil)
var _ Expr = (*Subscript)(nil)
var _ Expr = (*Starred)(nil)
var _ Expr = (*Name)(nil)
var _ Expr = (*List)(nil)
var _ Expr = (*Tuple)(nil)

// Slice
var _ Slicer = (*SliceBase)(nil)
var _ Slicer = (*Slice)(nil)
var _ Slicer = (*ExtSlice)(nil)
var _ Slicer = (*Index)(nil)

// Misc
var _ Ast = (*ExceptHandler)(nil)
var _ Ast = (*Arguments)(nil)
var _ Ast = (*Arg)(nil)
var _ Ast = (*Keyword)(nil)
var _ Ast = (*Alias)(nil)
var _ Ast = (*WithItem)(nil)

// Python types
var ASTType = py.ObjectType.NewTypeFlags("AST", "AST Node", nil, nil, py.ObjectType.Flags|py.TPFLAGS_BASE_EXC_SUBCLASS)

// Mod
var ModBaseType = ASTType.NewType("Mod", "Mod Node", nil, nil)
var ModuleType = ModBaseType.NewType("Module", "Module Node", nil, nil)
var InteractiveType = ModBaseType.NewType("Interactive", "Interactive Node", nil, nil)
var ExpressionType = ModBaseType.NewType("Expression", "Expression Node", nil, nil)
var SuiteType = ModBaseType.NewType("Suite", "Suite Node", nil, nil)

// Stmt
var StmtBaseType = ASTType.NewType("Stmt", "Stmt Node", nil, nil)
var FunctionDefType = StmtBaseType.NewType("FunctionDef", "FunctionDef Node", nil, nil)
var ClassDefType = StmtBaseType.NewType("ClassDef", "ClassDef Node", nil, nil)
var ReturnType = StmtBaseType.NewType("Return", "Return Node", nil, nil)
var DeleteType = StmtBaseType.NewType("Delete", "Delete Node", nil, nil)
var AssignType = StmtBaseType.NewType("Assign", "Assign Node", nil, nil)
var AugAssignType = StmtBaseType.NewType("AugAssign", "AugAssign Node", nil, nil)
var ForType = StmtBaseType.NewType("For", "For Node", nil, nil)
var WhileType = StmtBaseType.NewType("While", "While Node", nil, nil)
var IfType = StmtBaseType.NewType("If", "If Node", nil, nil)
var WithType = StmtBaseType.NewType("With", "With Node", nil, nil)
var RaiseType = StmtBaseType.NewType("Raise", "Raise Node", nil, nil)
var TryType = StmtBaseType.NewType("Try", "Try Node", nil, nil)
var AssertType = StmtBaseType.NewType("Assert", "Assert Node", nil, nil)
var ImportType = StmtBaseType.NewType("Import", "Import Node", nil, nil)
var ImportFromType = StmtBaseType.NewType("ImportFrom", "ImportFrom Node", nil, nil)
var GlobalType = StmtBaseType.NewType("Global", "Global Node", nil, nil)
var NonlocalType = StmtBaseType.NewType("Nonlocal", "Nonlocal Node", nil, nil)
var ExprStmtType = StmtBaseType.NewType("ExprStmt", "ExprStmt Node", nil, nil)
var PassType = StmtBaseType.NewType("Pass", "Pass Node", nil, nil)
var BreakType = StmtBaseType.NewType("Break", "Break Node", nil, nil)
var ContinueType = StmtBaseType.NewType("Continue", "Continue Node", nil, nil)

// Expr
var ExprBaseType = ASTType.NewType("Expr", "Expr Node", nil, nil)
var BoolOpType = ExprBaseType.NewType("BoolOp", "BoolOp Node", nil, nil)
var BinOpType = ExprBaseType.NewType("BinOp", "BinOp Node", nil, nil)
var UnaryOpType = ExprBaseType.NewType("UnaryOp", "UnaryOp Node", nil, nil)
var LambdaType = ExprBaseType.NewType("Lambda", "Lambda Node", nil, nil)
var IfExpType = ExprBaseType.NewType("IfExp", "IfExp Node", nil, nil)
var DictType = ExprBaseType.NewType("Dict", "Dict Node", nil, nil)
var SetType = ExprBaseType.NewType("Set", "Set Node", nil, nil)
var ListCompType = ExprBaseType.NewType("ListComp", "ListComp Node", nil, nil)
var SetCompType = ExprBaseType.NewType("SetComp", "SetComp Node", nil, nil)
var DictCompType = ExprBaseType.NewType("DictComp", "DictComp Node", nil, nil)
var GeneratorExpType = ExprBaseType.NewType("GeneratorExp", "GeneratorExp Node", nil, nil)
var YieldType = ExprBaseType.NewType("Yield", "Yield Node", nil, nil)
var YieldFromType = ExprBaseType.NewType("YieldFrom", "YieldFrom Node", nil, nil)
var CompareType = ExprBaseType.NewType("Compare", "Compare Node", nil, nil)
var CallType = ExprBaseType.NewType("Call", "Call Node", nil, nil)
var NumType = ExprBaseType.NewType("Num", "Num Node", nil, nil)
var StrType = ExprBaseType.NewType("Str", "Str Node", nil, nil)
var BytesType = ExprBaseType.NewType("Bytes", "Bytes Node", nil, nil)
var NameConstantType = ExprBaseType.NewType("NameConstant", "NameConstant Node", nil, nil)
var EllipsisType = ExprBaseType.NewType("Ellipsis", "Ellipsis Node", nil, nil)
var AttributeType = ExprBaseType.NewType("Attribute", "Attribute Node", nil, nil)
var SubscriptType = ExprBaseType.NewType("Subscript", "Subscript Node", nil, nil)
var StarredType = ExprBaseType.NewType("Starred", "Starred Node", nil, nil)
var NameType = ExprBaseType.NewType("Name", "Name Node", nil, nil)
var ListType = ExprBaseType.NewType("List", "List Node", nil, nil)
var TupleType = ExprBaseType.NewType("Tuple", "Tuple Node", nil, nil)
var SliceBaseType = ASTType.NewType("SliceBase", "SliceBase Node", nil, nil)

// Slicer
var SliceType = SliceBaseType.NewType("Slice", "Slice Node", nil, nil)
var ExtSliceType = SliceBaseType.NewType("ExtSlice", "ExtSlice Node", nil, nil)
var IndexType = SliceBaseType.NewType("Index", "Index Node", nil, nil)

// Misc
var ExceptHandlerType = ASTType.NewType("ExceptHandler", "ExceptHandler Node", nil, nil)
var ArgumentsType = ASTType.NewType("Arguments", "Arguments Node", nil, nil)
var ArgType = ASTType.NewType("Arg", "Arg Node", nil, nil)
var KeywordType = ASTType.NewType("Keyword", "Keyword Node", nil, nil)
var AliasType = ASTType.NewType("Alias", "Alias Node", nil, nil)
var WithItemType = ASTType.NewType("WithItem", "WithItem Node", nil, nil)

// Python type definitions
func (o *AST) Type() *py.Type           { return ASTType }
func (o *ModBase) Type() *py.Type       { return ModBaseType }
func (o *Module) Type() *py.Type        { return ModuleType }
func (o *Interactive) Type() *py.Type   { return InteractiveType }
func (o *Expression) Type() *py.Type    { return ExpressionType }
func (o *Suite) Type() *py.Type         { return SuiteType }
func (o *StmtBase) Type() *py.Type      { return StmtBaseType }
func (o *FunctionDef) Type() *py.Type   { return FunctionDefType }
func (o *ClassDef) Type() *py.Type      { return ClassDefType }
func (o *Return) Type() *py.Type        { return ReturnType }
func (o *Delete) Type() *py.Type        { return DeleteType }
func (o *Assign) Type() *py.Type        { return AssignType }
func (o *AugAssign) Type() *py.Type     { return AugAssignType }
func (o *For) Type() *py.Type           { return ForType }
func (o *While) Type() *py.Type         { return WhileType }
func (o *If) Type() *py.Type            { return IfType }
func (o *With) Type() *py.Type          { return WithType }
func (o *Raise) Type() *py.Type         { return RaiseType }
func (o *Try) Type() *py.Type           { return TryType }
func (o *Assert) Type() *py.Type        { return AssertType }
func (o *Import) Type() *py.Type        { return ImportType }
func (o *ImportFrom) Type() *py.Type    { return ImportFromType }
func (o *Global) Type() *py.Type        { return GlobalType }
func (o *Nonlocal) Type() *py.Type      { return NonlocalType }
func (o *ExprStmt) Type() *py.Type      { return ExprStmtType }
func (o *Pass) Type() *py.Type          { return PassType }
func (o *Break) Type() *py.Type         { return BreakType }
func (o *Continue) Type() *py.Type      { return ContinueType }
func (o *ExprBase) Type() *py.Type      { return ExprBaseType }
func (o *BoolOp) Type() *py.Type        { return BoolOpType }
func (o *BinOp) Type() *py.Type         { return BinOpType }
func (o *UnaryOp) Type() *py.Type       { return UnaryOpType }
func (o *Lambda) Type() *py.Type        { return LambdaType }
func (o *IfExp) Type() *py.Type         { return IfExpType }
func (o *Dict) Type() *py.Type          { return DictType }
func (o *Set) Type() *py.Type           { return SetType }
func (o *ListComp) Type() *py.Type      { return ListCompType }
func (o *SetComp) Type() *py.Type       { return SetCompType }
func (o *DictComp) Type() *py.Type      { return DictCompType }
func (o *GeneratorExp) Type() *py.Type  { return GeneratorExpType }
func (o *Yield) Type() *py.Type         { return YieldType }
func (o *YieldFrom) Type() *py.Type     { return YieldFromType }
func (o *Compare) Type() *py.Type       { return CompareType }
func (o *Call) Type() *py.Type          { return CallType }
func (o *Num) Type() *py.Type           { return NumType }
func (o *Str) Type() *py.Type           { return StrType }
func (o *Bytes) Type() *py.Type         { return BytesType }
func (o *NameConstant) Type() *py.Type  { return NameConstantType }
func (o *Ellipsis) Type() *py.Type      { return EllipsisType }
func (o *Attribute) Type() *py.Type     { return AttributeType }
func (o *Subscript) Type() *py.Type     { return SubscriptType }
func (o *Starred) Type() *py.Type       { return StarredType }
func (o *Name) Type() *py.Type          { return NameType }
func (o *List) Type() *py.Type          { return ListType }
func (o *Tuple) Type() *py.Type         { return TupleType }
func (o *SliceBase) Type() *py.Type     { return SliceBaseType }
func (o *Slice) Type() *py.Type         { return SliceType }
func (o *ExtSlice) Type() *py.Type      { return ExtSliceType }
func (o *Index) Type() *py.Type         { return IndexType }
func (o *ExceptHandler) Type() *py.Type { return ExceptHandlerType }
func (o *Arguments) Type() *py.Type     { return ArgumentsType }
func (o *Arg) Type() *py.Type           { return ArgType }
func (o *Keyword) Type() *py.Type       { return KeywordType }
func (o *Alias) Type() *py.Type         { return AliasType }
func (o *WithItem) Type() *py.Type      { return WithItemType }
