package parser

import (
	errors "github.com/moonbite-org/moonbite/error"
)

type StatementKind string
type ExpressionKind string
type TypeKind string
type VarKind string
type LoopKind string
type LiteralKind string
type ArithmeticUnaryKind string
type SignatureKind string

const (
	// statements
	PackageStatementKind              StatementKind = "statement:package"
	UseStatementKind                  StatementKind = "statement:use"
	ReturnStatementKind               StatementKind = "statement:return"
	DeferStatementKind                StatementKind = "statement:defer"
	BreakStatementKind                StatementKind = "statement:break"
	ContinueStatementKind             StatementKind = "statement:continue"
	YieldStatementKind                StatementKind = "statement:yield"
	DeclarationStatementKind          StatementKind = "statement:declaration"
	AssignmentStatementKind           StatementKind = "statement:assignment"
	TypeDefinitionStatementKind       StatementKind = "statement:type-definition"
	TraitDefinitionStatementKind      StatementKind = "statement:trait-definition"
	UnboundFunDefinitionStatementKind StatementKind = "statement:unbound-fun-definition"
	BoundFunDefinitionStatementKind   StatementKind = "statement:bound-fun-definition"
	ExpressionStatementKind           StatementKind = "statement:expression"
	LoopStatementKind                 StatementKind = "statement:loop"
	IfStatementKind                   StatementKind = "statement:if"
	SingleLineCommentStatementKind    StatementKind = "statement:single_line_comment"
	MultiLineCommentStatementKind     StatementKind = "statement:multi_line_comment"

	// expressions
	IdentifierExpressionKind      ExpressionKind = "expression:identifier"
	ArithmeticExpressionKind      ExpressionKind = "expression:arithmetic"
	BinaryExpressionKind          ExpressionKind = "expression:binary"
	ComparisonExpressionKind      ExpressionKind = "expression:comparison"
	CallExpressionKind            ExpressionKind = "expression:call"
	MemberExpressionKind          ExpressionKind = "expression:member"
	IndexExpressionKind           ExpressionKind = "expression:index"
	MatchExpressionKind           ExpressionKind = "expression:match"
	TypeCastExpressionKind        ExpressionKind = "expression:type-cast"
	CaretExpressionKind           ExpressionKind = "expression:caret"
	InstanceofExpressionKind      ExpressionKind = "expression:instanceof"
	MatchSelfExpressionKind       ExpressionKind = "expression:match_self"
	GroupExpressionKind           ExpressionKind = "expression:group"
	ThisExpressionKind            ExpressionKind = "expression:this"
	ArithmeticUnaryExpressionKind ExpressionKind = "expression:arithmetic_unary"
	AnonymousFunExpressionKind    ExpressionKind = "expression:anonymous_fun"
	OrExpressionKind              ExpressionKind = "expression:or"
	NotExpressionKind             ExpressionKind = "expression:not"
	GiveupExpressionKind          ExpressionKind = "expression:giveup"
	CoroutFunExpressionKind       ExpressionKind = "expression:corout_fun"
	GenFunExpressionKind          ExpressionKind = "expression:gen_fun"
	WarnExpressionKind            ExpressionKind = "expression:warn"

	// literal expressions
	StringLiteralExpressionKind   ExpressionKind = "expression:string-literal"
	RuneLiteralExpressionKind     ExpressionKind = "expression:rune-literal"
	BoolLiteralExpressionKind     ExpressionKind = "expression:bool-literal"
	NumberLiteralExpressionKind   ExpressionKind = "expression:number-literal"
	ListLiteralExpressionKind     ExpressionKind = "expression:list-literal"
	MapLiteralExpressionKind      ExpressionKind = "expression:map-literal"
	InstanceLiteralExpressionKind ExpressionKind = "expression:instance-literal"

	// types
	TypeIdentifierKind TypeKind = "type:identifier"
	StructLiteralKind  TypeKind = "type:struct-literal"
	OperatedTypeKind   TypeKind = "type:operated-type"
	TypedLiteralKind   TypeKind = "type:typed-literal"
	GroupTypeKind      TypeKind = "type:group"
	FunTypeKind        TypeKind = "type:fun"
	TraitTypeKind      TypeKind = "type:trait"
	BuiltinTypeKind    TypeKind = "type:builtin"

	// functions
	UnboundFunctionSignatureKind   SignatureKind = "signature:unbound"
	BoundFunctionSignatureKind     SignatureKind = "signature:bound"
	AnonymousFunctionSignatureKind SignatureKind = "signature:anonymous"

	// vars
	VariableKind VarKind = "var"
	ConstantKind VarKind = "const"

	// loops
	UnipartiteLoopKind LoopKind = "predicate:unipartite"
	BipartiteLoopKind  LoopKind = "predicate:bipartite"
	TripartiteLoopKind LoopKind = "predicate:tripartite"

	// literals
	StringLiteralKind   LiteralKind = "literal:string"
	BoolLiteralKind     LiteralKind = "literal:bool"
	RuneLiteralKind     LiteralKind = "literal:rune"
	NumberLiteralKind   LiteralKind = "literal:number"
	ListLiteralKind     LiteralKind = "literal:list"
	MapLiteralKind      LiteralKind = "literal:record"
	InstanceLiteralKind LiteralKind = "literal:instance"

	// unary arithmetic operation
	IncrementKind ArithmeticUnaryKind = "unary:increment"
	DecrementKind ArithmeticUnaryKind = "unary:decrement"
)

type ConstrainedType struct {
	Name       IdentifierExpression `json:"name"`
	Index      int                  `json:"index"`
	Constraint *TypeLiteral         `json:"constraint"`
	Location   errors.Location      `json:"location"`
}

type OperatorToken struct {
	Literal  string `json:"literal"`
	location errors.Location
}

type TypeLiteral interface {
	Location() errors.Location
	TypeKind() TypeKind
}

type TypedLiteral struct {
	TypeKind_ TypeKind          `json:"type_kind"`
	Type      TypeIdentifier    `json:"type"`
	Literal   LiteralExpression `json:"literal"`
	location  errors.Location
}

func (t TypedLiteral) TypeKind() TypeKind {
	return TypedLiteralKind
}

func (t TypedLiteral) Location() errors.Location {
	return t.location
}

type GroupType struct {
	TypeKind_ TypeKind    `json:"type_kind"`
	Type      TypeLiteral `json:"type"`
	location  errors.Location
}

func (t GroupType) TypeKind() TypeKind {
	return GroupTypeKind
}

func (t GroupType) Location() errors.Location {
	return t.location
}

type TypeIdentifier struct {
	TypeKind_ TypeKind            `json:"type_kind"`
	Name      Expression          `json:"name"`
	Generics  map[int]TypeLiteral `json:"generics"`
	location  errors.Location
}

func (t TypeIdentifier) TypeKind() TypeKind {
	return TypeIdentifierKind
}

func (t TypeIdentifier) Location() errors.Location {
	return t.location
}

type OperatedType struct {
	TypeKind_     TypeKind    `json:"type_kind"`
	LeftHandSide  TypeLiteral `json:"left_hand_side"`
	RightHandSide TypeLiteral `json:"right_hand_side"`
	Operator      OperatorToken
	location      errors.Location
}

func (t OperatedType) TypeKind() TypeKind {
	return OperatedTypeKind
}

func (t OperatedType) Location() errors.Location {
	return t.location
}

type ValueTypePair struct {
	Key      IdentifierExpression `json:"key"`
	Type     TypeLiteral          `json:"type"`
	Hidden   bool                 `json:"hidden"`
	Location errors.Location      `json:"location"`
}

type StructLiteral struct {
	TypeKind_ TypeKind `json:"type_kind"`
	Values    []ValueTypePair
	location  errors.Location
}

func (t StructLiteral) TypeKind() TypeKind {
	return StructLiteralKind
}

func (t StructLiteral) Location() errors.Location {
	return t.location
}

type TypedParameter struct {
	Name     IdentifierExpression `json:"name"`
	Type     TypeLiteral          `json:"type"`
	Variadic bool                 `json:"variadic"`
	Location errors.Location      `json:"location"`
}

type FunctionSignature interface {
	SignatureKind() SignatureKind
	GetGenerics() map[string]ConstrainedType
	GetParameters() []TypedParameter
	GetReturnType() *TypeLiteral
}

type AnonymousFunctionSignature struct {
	TypeKind_  TypeKind                   `json:"type_kind"`
	Parameters []TypedParameter           `json:"parameters"`
	Generics   map[string]ConstrainedType `json:"generics"`
	ReturnType *TypeLiteral               `json:"return_type"`
	location   errors.Location
}

func (s AnonymousFunctionSignature) Location() errors.Location {
	return s.location
}

func (s AnonymousFunctionSignature) SignatureKind() SignatureKind {
	return AnonymousFunctionSignatureKind
}

func (s AnonymousFunctionSignature) TypeKind() TypeKind {
	return FunTypeKind
}

func (s AnonymousFunctionSignature) GetGenerics() map[string]ConstrainedType {
	return s.Generics
}

func (s AnonymousFunctionSignature) GetParameters() []TypedParameter {
	return s.Parameters
}

func (s AnonymousFunctionSignature) GetReturnType() *TypeLiteral {
	return s.ReturnType
}

type UnboundFunctionSignature struct {
	TypeKind_  TypeKind                   `json:"type_kind"`
	Name       IdentifierExpression       `json:"name"`
	Parameters []TypedParameter           `json:"parameters"`
	Generics   map[string]ConstrainedType `json:"generics"`
	ReturnType *TypeLiteral               `json:"return_type"`
	location   errors.Location
}

func (s UnboundFunctionSignature) Location() errors.Location {
	return s.location
}

func (s UnboundFunctionSignature) SignatureKind() SignatureKind {
	return UnboundFunctionSignatureKind
}

func (s UnboundFunctionSignature) TypeKind() TypeKind {
	return FunTypeKind
}

func (s UnboundFunctionSignature) GetGenerics() map[string]ConstrainedType {
	return s.Generics
}

func (s UnboundFunctionSignature) GetParameters() []TypedParameter {
	return s.Parameters
}

func (s UnboundFunctionSignature) GetReturnType() *TypeLiteral {
	return s.ReturnType
}

type BoundFunctionSignature struct {
	Name       IdentifierExpression       `json:"name"`
	For        TypeIdentifier             `json:"for"`
	Generics   map[string]ConstrainedType `json:"generics"`
	Parameters []TypedParameter           `json:"parameters"`
	ReturnType *TypeLiteral               `json:"return_type"`
	location   errors.Location
}

func (s BoundFunctionSignature) Location() errors.Location {
	return s.location
}

func (s BoundFunctionSignature) SignatureKind() SignatureKind {
	return BoundFunctionSignatureKind
}

func (s BoundFunctionSignature) TypeKind() TypeKind {
	return FunTypeKind
}

func (s BoundFunctionSignature) GetGenerics() map[string]ConstrainedType {
	return s.Generics
}

func (s BoundFunctionSignature) GetParameters() []TypedParameter {
	return s.Parameters
}

func (s BoundFunctionSignature) GetReturnType() *TypeLiteral {
	return s.ReturnType
}

type Statement interface {
	Kind() StatementKind
	Location() errors.Location
}

type StatementList []Statement

type Expression interface {
	Kind() ExpressionKind
	Location() errors.Location
}

type Definition interface {
	Statement
	is_definition() bool
}

type Comment interface {
	is_comment() bool
	Statement
	Definition
}

type Ast struct {
	FileName    string           `json:"file_name"`
	FilePath    string           `json:"file_path"`
	Definitions []Definition     `json:"definitions"`
	Uses        []UseStatement   `json:"uses"`
	Package     PackageStatement `json:"package"`
	Comments    []Comment        `json:"comments"`
}

// STATEMENTS

type ExpressionStatement struct {
	Kind_ StatementKind `json:"kind"`

	Expression Expression `json:"expression"`
	location   errors.Location
}

func (s ExpressionStatement) Kind() StatementKind {
	return ExpressionStatementKind
}

func (s ExpressionStatement) Location() errors.Location {
	return s.location
}

type DeclarationStatement struct {
	Kind_ StatementKind `json:"kind"`

	VarKind  VarKind              `json:"var_kind"`
	Name     IdentifierExpression `json:"name"`
	Type     *TypeLiteral         `json:"type"`
	Value    *Expression          `json:"value"`
	Hidden   bool                 `json:"hidden"`
	location errors.Location
}

func (s DeclarationStatement) is_definition() bool {
	return true
}

func (s DeclarationStatement) Kind() StatementKind {
	return DeclarationStatementKind
}

func (s DeclarationStatement) Location() errors.Location {
	return s.location
}

type AssignmentStatement struct {
	Kind_ StatementKind `json:"kind"`

	LeftHandSide  Expression    `json:"left_hand_side"`
	RightHandSide Expression    `json:"right_hand_side"`
	Operator      OperatorToken `json:"operator"`
	location      errors.Location
}

func (s AssignmentStatement) Kind() StatementKind {
	return AssignmentStatementKind
}

func (s AssignmentStatement) Location() errors.Location {
	return s.location
}

type PackageStatement struct {
	Kind_ StatementKind `json:"kind"`

	Name     IdentifierExpression `json:"name"`
	location errors.Location
}

func (s PackageStatement) Kind() StatementKind {
	return PackageStatementKind
}

func (s PackageStatement) Location() errors.Location {
	return s.location
}

type UseStatement struct {
	Kind_ StatementKind `json:"kind"`

	Resource StringLiteralExpression `json:"resource"`
	As       *IdentifierExpression   `json:"as"`
	location errors.Location
}

func (s UseStatement) Kind() StatementKind {
	return UseStatementKind
}

func (s UseStatement) Location() errors.Location {
	return s.location
}

type TypeDefinitionStatement struct {
	Kind_           StatementKind              `json:"kind"`
	Name            IdentifierExpression       `json:"name"`
	Generics        map[string]ConstrainedType `json:"generics"`
	Implementations []TypeIdentifier           `json:"implementations"`
	Definition      TypeLiteral                `json:"definiton"`
	Hidden          bool                       `json:"hidden"`
	location        errors.Location
}

func (s TypeDefinitionStatement) is_definition() bool {
	return true
}

func (s TypeDefinitionStatement) Kind() StatementKind {
	return TypeDefinitionStatementKind
}

func (s TypeDefinitionStatement) Location() errors.Location {
	return s.location
}

type TraitDefinitionStatement struct {
	Kind_     StatementKind `json:"kind"`
	TypeKind_ TypeKind      `json:"type_kind"`

	Name       IdentifierExpression       `json:"name"`
	Generics   map[string]ConstrainedType `json:"generics"`
	Mimics     []TypeIdentifier           `json:"mimics"`
	Definition []UnboundFunctionSignature `json:"definition"`
	Hidden     bool                       `json:"hidden"`
	location   errors.Location
}

func (s TraitDefinitionStatement) is_definition() bool {
	return true
}

func (s TraitDefinitionStatement) TypeKind() TypeKind {
	return TraitTypeKind
}

func (s TraitDefinitionStatement) Kind() StatementKind {
	return TraitDefinitionStatementKind
}

func (s TraitDefinitionStatement) Location() errors.Location {
	return s.location
}

type FunDefinitionStatement interface {
	set_body(body []Statement)
	Statement
	Definition
}

type UnboundFunDefinitionStatement struct {
	Kind_ StatementKind `json:"kind"`

	Signature UnboundFunctionSignature `json:"signature"`
	Body      StatementList            `json:"body"`
	Hidden    bool                     `json:"hidden"`
	location  errors.Location
}

func (s UnboundFunDefinitionStatement) is_definition() bool {
	return true
}

func (s *UnboundFunDefinitionStatement) set_body(body []Statement) {
	s.Body = body
}

func (s UnboundFunDefinitionStatement) Kind() StatementKind {
	return UnboundFunDefinitionStatementKind
}

func (s UnboundFunDefinitionStatement) Location() errors.Location {
	return s.location
}

type BoundFunDefinitionStatement struct {
	Kind_ StatementKind `json:"kind"`

	Signature BoundFunctionSignature `json:"signature"`
	Body      StatementList          `json:"body"`
	Hidden    bool                   `json:"hidden"`
	location  errors.Location
}

func (s BoundFunDefinitionStatement) is_definition() bool {
	return true
}

func (s *BoundFunDefinitionStatement) set_body(body []Statement) {
	s.Body = body
}

func (s BoundFunDefinitionStatement) Kind() StatementKind {
	return BoundFunDefinitionStatementKind
}

func (s BoundFunDefinitionStatement) Location() errors.Location {
	return s.location
}

type ReturnStatement struct {
	Kind_ StatementKind `json:"kind"`

	Value    *Expression `json:"expression"`
	location errors.Location
}

func (s ReturnStatement) Kind() StatementKind {
	return ReturnStatementKind
}

func (s ReturnStatement) Location() errors.Location {
	return s.location
}

type DeferStatement struct {
	Kind_ StatementKind `json:"kind"`

	Value    Expression `json:"expression"`
	location errors.Location
}

func (s DeferStatement) Kind() StatementKind {
	return DeferStatementKind
}

func (s DeferStatement) Location() errors.Location {
	return s.location
}

type YieldStatement struct {
	Kind_ StatementKind `json:"kind"`

	Value    *Expression `json:"expression"`
	location errors.Location
}

func (s YieldStatement) Kind() StatementKind {
	return YieldStatementKind
}

func (s YieldStatement) Location() errors.Location {
	return s.location
}

type BreakStatement struct {
	Kind_ StatementKind `json:"kind"`

	location errors.Location
}

func (s BreakStatement) Kind() StatementKind {
	return BreakStatementKind
}

func (s BreakStatement) Location() errors.Location {
	return s.location
}

type ContinueStatement struct {
	Kind_ StatementKind `json:"kind"`

	location errors.Location
}

func (s ContinueStatement) Kind() StatementKind {
	return ContinueStatementKind
}

func (s ContinueStatement) Location() errors.Location {
	return s.location
}

type LoopPredicate interface {
	LoopKind() LoopKind
}

type UnipartiteLoopPredicate struct {
	Kind_ LoopKind `json:"kind"`

	Expression Expression `json:"expression"`
}

func (l UnipartiteLoopPredicate) LoopKind() LoopKind {
	return UnipartiteLoopKind
}

type BipartiteLoopPredicate struct {
	Kind_ LoopKind `json:"kind"`

	Key      *IdentifierExpression `json:"key"`
	Value    *IdentifierExpression `json:"value"`
	Iterator Expression            `json:"iterator"`
}

func (l BipartiteLoopPredicate) LoopKind() LoopKind {
	return BipartiteLoopKind
}

type TripartiteLoopPredicate struct {
	Kind_ LoopKind `json:"kind"`

	Declaration *DeclarationStatement `json:"declaration"`
	Predicate   Expression            `json:"predicate"`
	Procedure   *Expression           `json:"procedure"`
}

func (l TripartiteLoopPredicate) LoopKind() LoopKind {
	return TripartiteLoopKind
}

type LoopStatement struct {
	Kind_ StatementKind `json:"kind"`

	Predicate LoopPredicate `json:"predicate"`
	Body      StatementList `json:"body"`
	location  errors.Location
}

func (s LoopStatement) Kind() StatementKind {
	return LoopStatementKind
}

func (s LoopStatement) Location() errors.Location {
	return s.location
}

type PredicateBlock struct {
	Predicate Expression    `json:"predicate"`
	Body      StatementList `json:"body"`
}

type IfStatement struct {
	Kind_ StatementKind `json:"kind"`

	MainBlock    PredicateBlock   `json:"main_block"`
	ElseIfBlocks []PredicateBlock `json:"else_if_blocks"`
	ElseBlock    StatementList    `json:"else_block"`
	location     errors.Location
}

func (s IfStatement) Kind() StatementKind {
	return IfStatementKind
}

func (s IfStatement) Location() errors.Location {
	return s.location
}

type SingleLineCommentStatement struct {
	Kind_ StatementKind `json:"kind"`

	Comment  string `json:"comment"`
	location errors.Location
}

func (s SingleLineCommentStatement) is_comment() bool {
	return true
}
func (s SingleLineCommentStatement) is_definition() bool {
	return true
}

func (s SingleLineCommentStatement) Kind() StatementKind {
	return SingleLineCommentStatementKind
}

func (s SingleLineCommentStatement) Location() errors.Location {
	return s.location
}

type MultiLineCommentStatement struct {
	Kind_ StatementKind `json:"kind"`

	Comment  string `json:"comment"`
	location errors.Location
}

func (s MultiLineCommentStatement) is_comment() bool {
	return true
}

func (s MultiLineCommentStatement) is_definition() bool {
	return true
}

func (s MultiLineCommentStatement) Kind() StatementKind {
	return MultiLineCommentStatementKind
}

func (s MultiLineCommentStatement) Location() errors.Location {
	return s.location
}

// EXPRESSIONS

type AnonymousFunExpression struct {
	Kind_ ExpressionKind `json:"kind"`

	Signature AnonymousFunctionSignature `json:"signature"`
	Body      StatementList              `json:"body"`
	location  errors.Location
}

func (e AnonymousFunExpression) Kind() ExpressionKind {
	return AnonymousFunExpressionKind
}

func (e AnonymousFunExpression) Location() errors.Location {
	return e.location
}

type IdentifierExpression struct {
	Kind_ ExpressionKind `json:"kind"`

	Value    string `json:"value"`
	location errors.Location
}

func (e IdentifierExpression) Kind() ExpressionKind {
	return IdentifierExpressionKind
}

func (e IdentifierExpression) Location() errors.Location {
	return e.location
}

type CaretExpression struct {
	Kind_ ExpressionKind `json:"kind"`

	location errors.Location
}

func (e CaretExpression) Kind() ExpressionKind {
	return CaretExpressionKind
}

func (e CaretExpression) Location() errors.Location {
	return e.location
}

type TypeCastExpression struct {
	Kind_ ExpressionKind `json:"kind"`

	Value    Expression     `json:"value"`
	Type     TypeIdentifier `json:"type"`
	location errors.Location
}

func (e TypeCastExpression) Kind() ExpressionKind {
	return TypeCastExpressionKind
}

func (e TypeCastExpression) Location() errors.Location {
	return e.location
}

type InstanceofExpression struct {
	Kind_ ExpressionKind `json:"kind"`

	LeftHandSide  Expression  `json:"left_hand_side"`
	RightHandSide TypeLiteral `json:"right_hand_side"`
	location      errors.Location
}

func (e InstanceofExpression) Kind() ExpressionKind {
	return InstanceofExpressionKind
}

func (e InstanceofExpression) Location() errors.Location {
	return e.location
}

type MatchSelfExpression struct {
	Kind_ ExpressionKind `json:"kind"`

	location errors.Location
}

func (e MatchSelfExpression) Kind() ExpressionKind {
	return MatchSelfExpressionKind
}

func (e MatchSelfExpression) Location() errors.Location {
	return e.location
}

type ArithmeticExpression struct {
	Kind_ ExpressionKind `json:"kind"`

	LeftHandSide  Expression    `json:"left_hand_side"`
	RightHandSide Expression    `json:"right_hand_side"`
	Operator      OperatorToken `json:"operator"`
	location      errors.Location
}

func (e ArithmeticExpression) Kind() ExpressionKind {
	return ArithmeticExpressionKind
}

func (e ArithmeticExpression) Location() errors.Location {
	return e.location
}

type BinaryExpression struct {
	Kind_ ExpressionKind `json:"kind"`

	LeftHandSide  Expression    `json:"left_hand_side"`
	RightHandSide Expression    `json:"right_hand_side"`
	Operator      OperatorToken `json:"operator"`
	location      errors.Location
}

func (e BinaryExpression) Kind() ExpressionKind {
	return BinaryExpressionKind
}

func (e BinaryExpression) Location() errors.Location {
	return e.location
}

type ComparisonExpression struct {
	Kind_ ExpressionKind `json:"kind"`

	LeftHandSide  Expression    `json:"left_hand_side"`
	RightHandSide Expression    `json:"right_hand_side"`
	Operator      OperatorToken `json:"operator"`
	location      errors.Location
}

func (e ComparisonExpression) Kind() ExpressionKind {
	return ComparisonExpressionKind
}

func (e ComparisonExpression) Location() errors.Location {
	return e.location
}

type LiteralExpression interface {
	Expression
	LiteralKind() LiteralKind
}

type StringLiteralExpression struct {
	Kind_ ExpressionKind `json:"kind"`

	Value    string `json:"value"`
	location errors.Location
}

func (e StringLiteralExpression) Kind() ExpressionKind {
	return StringLiteralExpressionKind
}

func (e StringLiteralExpression) LiteralKind() LiteralKind {
	return StringLiteralKind
}

func (e StringLiteralExpression) Location() errors.Location {
	return e.location
}

type RuneLiteralExpression struct {
	Kind_ ExpressionKind `json:"kind"`

	Value    rune `json:"value"`
	location errors.Location
}

func (e RuneLiteralExpression) Kind() ExpressionKind {
	return RuneLiteralExpressionKind
}

func (e RuneLiteralExpression) LiteralKind() LiteralKind {
	return RuneLiteralKind
}

func (e RuneLiteralExpression) Location() errors.Location {
	return e.location
}

type BoolLiteralExpression struct {
	Kind_ ExpressionKind `json:"kind"`

	Value    bool `json:"value"`
	location errors.Location
}

func (e BoolLiteralExpression) Kind() ExpressionKind {
	return BoolLiteralExpressionKind
}

func (e BoolLiteralExpression) LiteralKind() LiteralKind {
	return BoolLiteralKind
}

func (e BoolLiteralExpression) Location() errors.Location {
	return e.location
}

type NumberLiteral struct {
	Type  TypeLiteral `json:"type"`
	Value interface{} `json:"value"`
}

type NumberLiteralExpression struct {
	Kind_ ExpressionKind `json:"kind"`

	Value    NumberLiteral `json:"value"`
	location errors.Location
}

func (e NumberLiteralExpression) Kind() ExpressionKind {
	return NumberLiteralExpressionKind
}

func (e NumberLiteralExpression) LiteralKind() LiteralKind {
	return NumberLiteralKind
}

func (e NumberLiteralExpression) Location() errors.Location {
	return e.location
}

type KeyValueEntry struct {
	Key   IdentifierExpression `json:"key"`
	Value Expression           `json:"value"`
}

type ListLiteralExpression struct {
	Kind_ ExpressionKind `json:"kind"`

	Value    []KeyValueEntry `json:"value"`
	location errors.Location
}

func (e ListLiteralExpression) Kind() ExpressionKind {
	return ListLiteralExpressionKind
}

func (e ListLiteralExpression) LiteralKind() LiteralKind {
	return ListLiteralKind
}

func (e ListLiteralExpression) Location() errors.Location {
	return e.location
}

type MapLiteralExpression struct {
	Kind_ ExpressionKind `json:"kind"`

	Value    []KeyValueEntry `json:"value"`
	location errors.Location
}

func (e MapLiteralExpression) Kind() ExpressionKind {
	return MapLiteralExpressionKind
}

func (e MapLiteralExpression) LiteralKind() LiteralKind {
	return MapLiteralKind
}

func (e MapLiteralExpression) Location() errors.Location {
	return e.location
}

type InstanceLiteralExpression struct {
	Kind_ ExpressionKind `json:"kind"`

	Type     TypeIdentifier  `json:"type"`
	Value    []KeyValueEntry `json:"value"`
	location errors.Location
}

func (e InstanceLiteralExpression) Kind() ExpressionKind {
	return InstanceLiteralExpressionKind
}

func (e InstanceLiteralExpression) LiteralKind() LiteralKind {
	return InstanceLiteralKind
}

func (e InstanceLiteralExpression) Location() errors.Location {
	return e.location
}

type CallExpression struct {
	Kind_ ExpressionKind `json:"kind"`

	Callee    Expression   `json:"callee"`
	Arguments []Expression `json:"arguments"`
	location  errors.Location
}

func (e CallExpression) Kind() ExpressionKind {
	return CallExpressionKind
}

func (e CallExpression) Location() errors.Location {
	return e.location
}

type MemberExpression struct {
	Kind_ ExpressionKind `json:"kind"`

	LeftHandSide  Expression           `json:"left_hand_side"`
	RightHandSide IdentifierExpression `json:"right_hand_side"`
	location      errors.Location
}

func (e MemberExpression) Kind() ExpressionKind {
	return MemberExpressionKind
}

func (e MemberExpression) Location() errors.Location {
	return e.location
}

type IndexExpression struct {
	Kind_ ExpressionKind `json:"kind"`

	Host     Expression `json:"host"`
	Index    Expression `json:"index"`
	location errors.Location
}

func (e IndexExpression) Kind() ExpressionKind {
	return IndexExpressionKind
}

func (e IndexExpression) Location() errors.Location {
	return e.location
}

type MatchExpression struct {
	Kind_ ExpressionKind `json:"kind"`

	Against   Expression       `json:"against"`
	Blocks    []PredicateBlock `json:"blocks"`
	BaseBlock StatementList    `json:"base_block"`
	location  errors.Location
}

func (e MatchExpression) Kind() ExpressionKind {
	return MatchExpressionKind
}

func (e MatchExpression) Location() errors.Location {
	return e.location
}

type GroupExpression struct {
	Kind_ ExpressionKind `json:"kind"`

	Expression Expression `json:"expression"`
	location   errors.Location
}

func (e GroupExpression) Kind() ExpressionKind {
	return GroupExpressionKind
}

func (e GroupExpression) Location() errors.Location {
	return e.location
}

type ThisExpression struct {
	Kind_ ExpressionKind `json:"kind"`

	location errors.Location
}

func (e ThisExpression) Kind() ExpressionKind {
	return ThisExpressionKind
}

func (e ThisExpression) Location() errors.Location {
	return e.location
}

type ArithmeticUnaryExpression struct {
	Kind_ ExpressionKind `json:"kind"`

	Expression Expression          `json:"expression"`
	Operation  ArithmeticUnaryKind `json:"operation"`
	Pre        bool                `json:"pre"`
	location   errors.Location
}

func (e ArithmeticUnaryExpression) Kind() ExpressionKind {
	return ArithmeticUnaryExpressionKind
}

func (e ArithmeticUnaryExpression) Location() errors.Location {
	return e.location
}

type OrExpression struct {
	Kind_ ExpressionKind `json:"kind"`

	LeftHandSide  Expression `json:"left_hand_side"`
	RightHandSide Expression `json:"right_hand_side"`
	location      errors.Location
}

func (e OrExpression) Kind() ExpressionKind {
	return OrExpressionKind
}

func (e OrExpression) Location() errors.Location {
	return e.location
}

type NotExpression struct {
	Kind_ ExpressionKind `json:"kind"`

	Expression Expression
	location   errors.Location
}

func (e NotExpression) Kind() ExpressionKind {
	return NotExpressionKind
}

func (e NotExpression) Location() errors.Location {
	return e.location
}

type GiveupExpression struct {
	Kind_ ExpressionKind `json:"kind"`

	location errors.Location
}

func (e GiveupExpression) Kind() ExpressionKind {
	return GiveupExpressionKind
}

func (e GiveupExpression) Location() errors.Location {
	return e.location
}

type CoroutFunExpression struct {
	Kind_ ExpressionKind `json:"kind"`

	Fun      AnonymousFunExpression `json:"fun"`
	location errors.Location
}

func (e CoroutFunExpression) Kind() ExpressionKind {
	return CoroutFunExpressionKind
}

func (e CoroutFunExpression) Location() errors.Location {
	return e.location
}

type GenFunExpression struct {
	Kind_ ExpressionKind `json:"kind"`

	Fun      AnonymousFunExpression `json:"fun"`
	location errors.Location
}

func (e GenFunExpression) Kind() ExpressionKind {
	return GenFunExpressionKind
}

func (e GenFunExpression) Location() errors.Location {
	return e.location
}

type WarnExpression struct {
	Kind_ ExpressionKind `json:"kind"`

	Argument Expression
	location errors.Location
}

func (e WarnExpression) Kind() ExpressionKind {
	return WarnExpressionKind
}

func (e WarnExpression) Location() errors.Location {
	return e.location
}
