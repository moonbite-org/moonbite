package parser

import "github.com/moonbite-org/moonbite/common"

type statement_kind string
type expression_kind string
type type_kind string
type var_kind string
type loop_kind string
type literal_kind string
type arithmetic_unary_kind string

const (
	// statements
	PackageStatementKind              statement_kind = "statement:package"
	UseStatementKind                  statement_kind = "statement:use"
	ReturnStatementKind               statement_kind = "statement:return"
	BreakStatementKind                statement_kind = "statement:break"
	ContinueStatementKind             statement_kind = "statement:continue"
	YieldStatementKind                statement_kind = "statement:yield"
	DeclarationStatementKind          statement_kind = "statement:declaration"
	AssignmentStatementKind           statement_kind = "statement:assignment"
	TypeDefinitionStatementKind       statement_kind = "statement:type-definition"
	TraitDefinitionStatementKind      statement_kind = "statement:trait-definition"
	UnboundFunDefinitionStatementKind statement_kind = "statement:unbound-fun-definition"
	BoundFunDefinitionStatementKind   statement_kind = "statement:bound-fun-definition"
	ExpressionStatementKind           statement_kind = "statement:expression"
	LoopStatementKind                 statement_kind = "statement:loop"
	IfStatementKind                   statement_kind = "statement:if"
	// OrStatementKind                   statement_kind = "statement:or"
	SingleLineCommentStatementKind statement_kind = "statement:single_line_comment"
	MultiLineCommentStatementKind  statement_kind = "statement:multi_line_comment"

	// expressions
	IdentifierExpressionKind      expression_kind = "expression:identifier"
	ArithmeticExpressionKind      expression_kind = "expression:arithmetic"
	BinaryExpressionKind          expression_kind = "expression:binary"
	CallExpressionKind            expression_kind = "expression:call"
	MemberExpressionKind          expression_kind = "expression:member"
	IndexExpressionKind           expression_kind = "expression:index"
	MatchExpressionKind           expression_kind = "expression:match"
	TypeCastExpressionKind        expression_kind = "expression:type-cast"
	CaretExpressionKind           expression_kind = "expression:caret"
	InstanceofExpressionKind      expression_kind = "expression:instanceof"
	MatchSelfExpressionKind       expression_kind = "expression:match_self"
	GroupExpressionKind           expression_kind = "expression:group"
	ThisExpressionKind            expression_kind = "expression:this"
	ArithmeticUnaryExpressionKind expression_kind = "expression:arithmetic_unary"
	AnonymousFunExpressionKind    expression_kind = "expression:anonymous_fun"
	OrExpressionKind              expression_kind = "expression:or"
	NotExpressionKind             expression_kind = "expression:not"

	// literal expressions
	StringLiteralExpressionKind   expression_kind = "expression:string-literal"
	RuneLiteralExpressionKind     expression_kind = "expression:rune-literal"
	BoolLiteralExpressionKind     expression_kind = "expression:bool-literal"
	NumberLiteralExpressionKind   expression_kind = "expression:number-literal"
	ListLiteralExpressionKind     expression_kind = "expression:list-literal"
	RecordLiteralExpressionKind   expression_kind = "expression:record-literal"
	InstanceLiteralExpressionKind expression_kind = "expression:instance-literal"

	// types
	type_identifier type_kind = "type:type-identifier"
	struct_literal  type_kind = "type:struct-literal"
	operated_type   type_kind = "type:operated-type"
	typed_literal   type_kind = "type:typed-literal"
	group_type      type_kind = "type:group"
	fun_type        type_kind = "type:fun"

	// vars
	VariableKind var_kind = "var"
	ConstantKind var_kind = "const"

	// loops
	UnipartiteLoopKind loop_kind = "predicate:unipartite"
	BipartiteLoopKind  loop_kind = "predicate:bipartite"
	TripartiteLoopKind loop_kind = "predicate:tripartite"

	// literals
	StringLiteralKind   literal_kind = "literal:string"
	BoolLiteralKind     literal_kind = "literal:bool"
	RuneLiteralKind     literal_kind = "literal:rune"
	NumberLiteralKind   literal_kind = "literal:number"
	ListLiteralKind     literal_kind = "literal:list"
	RecordLiteralKind   literal_kind = "literal:record"
	InstanceLiteralKind literal_kind = "literal:instance"

	// stepped change operation
	increment_kind arithmetic_unary_kind = "stepped:increment"
	decrement_kind arithmetic_unary_kind = "stepped:decrement"
)

// func stringify_list[T printable](items []T, open, close string, keepempty bool, seperator string) string {
// 	result := []string{}
// 	literal := ""

// 	for _, item := range items {
// 		result = append(result, item.String())
// 	}

// 	if len(items) != 0 {
// 		literal += open
// 		literal += strings.Join(result, seperator)
// 		literal += close
// 	} else {
// 		if keepempty {
// 			literal += open
// 			literal += close
// 		}
// 	}

// 	return literal
// }

type ConstrainedType struct {
	Name       Expression   `json:"name"`
	Constraint *TypeLiteral `json:"constraint"`
}

type TypeLiteral interface {
	TypeKind() type_kind
}

type TypedLiteral struct {
	Type    TypeIdentifier    `json:"type"`
	Literal LiteralExpression `json:"literal"`
}

func (t TypedLiteral) TypeKind() type_kind {
	return typed_literal
}

type GroupType struct {
	Type TypeLiteral `json:"type"`
}

func (t GroupType) TypeKind() type_kind {
	return group_type
}

type TypeIdentifier struct {
	Name     Expression    `json:"name"`
	Generics []TypeLiteral `json:"generics"`
}

func (t TypeIdentifier) TypeKind() type_kind {
	return type_identifier
}

type OperatedType struct {
	LeftHandSide  TypeLiteral `json:"left_hand_side"`
	RightHandSide TypeLiteral `json:"right_hand_side"`
	Operator      string
}

func (t OperatedType) TypeKind() type_kind {
	return operated_type
}

type ValueTypePair struct {
	Key  IdentifierExpression `json:"key"`
	Type TypeLiteral          `json:"type"`
}

type StructLiteral []ValueTypePair

func (t StructLiteral) TypeKind() type_kind {
	return struct_literal
}

type TypedParameter struct {
	Name     IdentifierExpression `json:"name"`
	Type     TypeLiteral          `json:"type"`
	Variadic bool                 `json:"variadic"`
	Location common.Location
}

type FunctionSignature interface {
	is_fun_signature() bool
}

type AnonymousFunctionSignature struct {
	Parameters []TypedParameter  `json:"parameters"`
	Generics   []ConstrainedType `json:"generics"`
	ReturnType *TypeLiteral      `json:"return_type"`
	location   common.Location
}

func (s AnonymousFunctionSignature) is_fun_signature() bool {
	return true
}

func (s AnonymousFunctionSignature) TypeKind() type_kind {
	return fun_type
}

type UnboundFunctionSignature struct {
	Name       IdentifierExpression `json:"name"`
	Parameters []TypedParameter     `json:"parameters"`
	Generics   []ConstrainedType    `json:"generics"`
	ReturnType *TypeLiteral         `json:"return_type"`
	location   common.Location
}

func (s UnboundFunctionSignature) is_fun_signature() bool {
	return true
}

type BoundFunctionSignature struct {
	Name       IdentifierExpression `json:"name"`
	For        TypeIdentifier       `json:"for"`
	Generics   []ConstrainedType    `json:"generics"`
	Parameters []TypedParameter     `json:"parameters"`
	ReturnType *TypeLiteral         `json:"return_type"`
	location   common.Location
}

func (s BoundFunctionSignature) is_fun_signature() bool {
	return true
}

type Statement interface {
	Kind() statement_kind
	Location() common.Location
}

type StatementList []Statement

type Expression interface {
	Kind() expression_kind
	Location() common.Location
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
	Expression Expression `json:"expression"`
	location   common.Location
}

func (s ExpressionStatement) Kind() statement_kind {
	return ExpressionStatementKind
}

func (s ExpressionStatement) Location() common.Location {
	return s.location
}

type DeclarationStatement struct {
	VarKind  var_kind             `json:"var_kind"`
	Name     IdentifierExpression `json:"name"`
	Type     *TypeLiteral         `json:"type"`
	Value    *Expression          `json:"value"`
	Hidden   bool                 `json:"hidden"`
	location common.Location
}

func (s DeclarationStatement) is_definition() bool {
	return true
}

func (s DeclarationStatement) Kind() statement_kind {
	return DeclarationStatementKind
}

func (s DeclarationStatement) Location() common.Location {
	return s.location
}

type AssignmentStatement struct {
	LeftHandSide  Expression `json:"left_hand_side"`
	RightHandSide Expression `json:"right_hand_side"`
	Operator      string     `json:"operator"`
	location      common.Location
}

func (s AssignmentStatement) Kind() statement_kind {
	return AssignmentStatementKind
}

func (s AssignmentStatement) Location() common.Location {
	return s.location
}

type PackageStatement struct {
	Name     IdentifierExpression `json:"name"`
	location common.Location
}

func (s PackageStatement) Kind() statement_kind {
	return PackageStatementKind
}

func (s PackageStatement) Location() common.Location {
	return s.location
}

type UseStatement struct {
	Resource IdentifierExpression  `json:"resource"`
	As       *IdentifierExpression `json:"as"`
	location common.Location
}

func (s UseStatement) Kind() statement_kind {
	return UseStatementKind
}

func (s UseStatement) Location() common.Location {
	return s.location
}

type TypeDefinitionStatement struct {
	Name            IdentifierExpression `json:"name"`
	Generics        []ConstrainedType    `json:"generics"`
	Implementations []TypeIdentifier     `json:"implementations"`
	Definition      TypeLiteral          `json:"definiton"`
	Hidden          bool                 `json:"hidden"`
	location        common.Location
}

func (s TypeDefinitionStatement) is_definition() bool {
	return true
}

func (s TypeDefinitionStatement) Kind() statement_kind {
	return TypeDefinitionStatementKind
}

func (s TypeDefinitionStatement) Location() common.Location {
	return s.location
}

type TraitDefinitionStatement struct {
	Name       IdentifierExpression       `json:"name"`
	Generics   []ConstrainedType          `json:"generics"`
	Mimics     []TypeIdentifier           `json:"mimics"`
	Definition []UnboundFunctionSignature `json:"definition"`
	Hidden     bool                       `json:"hidden"`
	location   common.Location
}

func (s TraitDefinitionStatement) is_definition() bool {
	return true
}

func (s TraitDefinitionStatement) Kind() statement_kind {
	return TraitDefinitionStatementKind
}

func (s TraitDefinitionStatement) Location() common.Location {
	return s.location
}

type FunDefinitionStatement interface {
	set_body(body []Statement)
	Statement
	Definition
}

type UnboundFunDefinitionStatement struct {
	Signature UnboundFunctionSignature `json:"signature"`
	Body      StatementList            `json:"body"`
	Hidden    bool                     `json:"hidden"`
	location  common.Location
}

func (s UnboundFunDefinitionStatement) is_definition() bool {
	return true
}

func (s *UnboundFunDefinitionStatement) set_body(body []Statement) {
	s.Body = body
}

func (s UnboundFunDefinitionStatement) Kind() statement_kind {
	return UnboundFunDefinitionStatementKind
}

func (s UnboundFunDefinitionStatement) Location() common.Location {
	return s.location
}

type BoundFunDefinitionStatement struct {
	Signature BoundFunctionSignature `json:"signature"`
	Body      StatementList          `json:"body"`
	Hidden    bool                   `json:"hidden"`
	location  common.Location
}

func (s BoundFunDefinitionStatement) is_definition() bool {
	return true
}

func (s *BoundFunDefinitionStatement) set_body(body []Statement) {
	s.Body = body
}

func (s BoundFunDefinitionStatement) Kind() statement_kind {
	return BoundFunDefinitionStatementKind
}

func (s BoundFunDefinitionStatement) Location() common.Location {
	return s.location
}

type ReturnStatement struct {
	Value    *Expression `json:"expression"`
	location common.Location
}

func (s ReturnStatement) Kind() statement_kind {
	return ReturnStatementKind
}

func (s ReturnStatement) Location() common.Location {
	return s.location
}

type YieldStatement struct {
	Value    *Expression `json:"expression"`
	location common.Location
}

func (s YieldStatement) Kind() statement_kind {
	return YieldStatementKind
}

func (s YieldStatement) Location() common.Location {
	return s.location
}

type BreakStatement struct {
	location common.Location
}

func (s BreakStatement) Kind() statement_kind {
	return BreakStatementKind
}

func (s BreakStatement) Location() common.Location {
	return s.location
}

type ContinueStatement struct {
	location common.Location
}

func (s ContinueStatement) Kind() statement_kind {
	return ContinueStatementKind
}

func (s ContinueStatement) Location() common.Location {
	return s.location
}

type LoopPredicate interface {
	LoopKind() loop_kind
}

type UnipartiteLoopPredicate struct {
	Expression Expression `json:"expression"`
}

func (l UnipartiteLoopPredicate) LoopKind() loop_kind {
	return UnipartiteLoopKind
}

type BipartiteLoopPredicate struct {
	Key      *IdentifierExpression `json:"key"`
	Value    *IdentifierExpression `json:"value"`
	Iterator Expression            `json:"iterator"`
}

func (l BipartiteLoopPredicate) LoopKind() loop_kind {
	return BipartiteLoopKind
}

type TripartiteLoopPredicate struct {
	Declaration *DeclarationStatement `json:"declaration"`
	Predicate   Expression            `json:"predicate"`
	Procedure   *Expression           `json:"procedure"`
}

func (l TripartiteLoopPredicate) LoopKind() loop_kind {
	return TripartiteLoopKind
}

type LoopStatement struct {
	Predicate LoopPredicate `json:"predicate"`
	Body      StatementList `json:"body"`
	location  common.Location
}

func (s LoopStatement) Kind() statement_kind {
	return LoopStatementKind
}

func (s LoopStatement) Location() common.Location {
	return s.location
}

type PredicateBlock struct {
	Predicate Expression    `json:"predicate"`
	Body      StatementList `json:"body"`
}

type IfStatement struct {
	MainBlock    PredicateBlock   `json:"main_block"`
	ElseIfBlocks []PredicateBlock `json:"else_if_blocks"`
	ElseBlock    StatementList    `json:"else_block"`
	location     common.Location
}

func (s IfStatement) Kind() statement_kind {
	return IfStatementKind
}

func (s IfStatement) Location() common.Location {
	return s.location
}

// type OrStatement struct {
// 	Try      Expression `json:"try"`
// 	Fail     Expression `json:"fail"`
// 	location common.Location
// }

// func (s OrStatement) Kind() statement_kind {
// 	return OrStatementKind
// }

// func (s OrStatement) Location() common.Location {
// 	return s.location
// }

type SingleLineCommentStatement struct {
	Comment  string `json:"comment"`
	location common.Location
}

func (s SingleLineCommentStatement) is_comment() bool {
	return true
}
func (s SingleLineCommentStatement) is_definition() bool {
	return true
}

func (s SingleLineCommentStatement) Kind() statement_kind {
	return SingleLineCommentStatementKind
}

func (s SingleLineCommentStatement) Location() common.Location {
	return s.location
}

type MultiLineCommentStatement struct {
	Comment  string `json:"comment"`
	location common.Location
}

func (s MultiLineCommentStatement) is_comment() bool {
	return true
}

func (s MultiLineCommentStatement) is_definition() bool {
	return true
}

func (s MultiLineCommentStatement) Kind() statement_kind {
	return MultiLineCommentStatementKind
}

func (s MultiLineCommentStatement) Location() common.Location {
	return s.location
}

// EXPRESSIONS

type AnonymousFunExpression struct {
	Signature AnonymousFunctionSignature `json:"signature"`
	Body      StatementList              `json:"body"`
	location  common.Location
}

func (e AnonymousFunExpression) Kind() expression_kind {
	return AnonymousFunExpressionKind
}

func (e AnonymousFunExpression) Location() common.Location {
	return e.location
}

type IdentifierExpression struct {
	Value    string `json:"value"`
	location common.Location
}

func (e IdentifierExpression) Kind() expression_kind {
	return IdentifierExpressionKind
}

func (e IdentifierExpression) Location() common.Location {
	return e.location
}

type CaretExpression struct {
	location common.Location
}

func (e CaretExpression) Kind() expression_kind {
	return CaretExpressionKind
}

func (e CaretExpression) Location() common.Location {
	return e.location
}

type TypeCastExpression struct {
	Value    Expression     `json:"value"`
	Type     TypeIdentifier `json:"type"`
	location common.Location
}

func (e TypeCastExpression) Kind() expression_kind {
	return TypeCastExpressionKind
}

func (e TypeCastExpression) Location() common.Location {
	return e.location
}

type InstanceofExpression struct {
	LeftHandSide  Expression  `json:"left_hand_side"`
	RightHandSide TypeLiteral `json:"right_hand_side"`
	location      common.Location
}

func (e InstanceofExpression) Kind() expression_kind {
	return InstanceofExpressionKind
}

func (e InstanceofExpression) Location() common.Location {
	return e.location
}

type MatchSelfExpression struct {
	location common.Location
}

func (e MatchSelfExpression) Kind() expression_kind {
	return MatchSelfExpressionKind
}

func (e MatchSelfExpression) Location() common.Location {
	return e.location
}

type ArithmeticExpression struct {
	LeftHandSide  Expression `json:"left_hand_side"`
	RightHandSide Expression `json:"right_hand_side"`
	Operator      string     `json:"operator"`
	location      common.Location
}

func (e ArithmeticExpression) Kind() expression_kind {
	return ArithmeticExpressionKind
}

func (e ArithmeticExpression) Location() common.Location {
	return e.location
}

type BinaryExpression struct {
	LeftHandSide  Expression `json:"left_hand_side"`
	RightHandSide Expression `json:"right_hand_side"`
	Operator      string     `json:"operator"`
	location      common.Location
}

func (e BinaryExpression) Kind() expression_kind {
	return BinaryExpressionKind
}

func (e BinaryExpression) Location() common.Location {
	return e.location
}

type LiteralExpression interface {
	Expression
	LiteralKind() literal_kind
}

type StringLiteralExpression struct {
	Value    string `json:"value"`
	location common.Location
}

func (e StringLiteralExpression) Kind() expression_kind {
	return StringLiteralExpressionKind
}

func (e StringLiteralExpression) LiteralKind() literal_kind {
	return StringLiteralKind
}

func (e StringLiteralExpression) Location() common.Location {
	return e.location
}

type RuneLiteralExpression struct {
	Value    rune `json:"value"`
	location common.Location
}

func (e RuneLiteralExpression) Kind() expression_kind {
	return RuneLiteralExpressionKind
}

func (e RuneLiteralExpression) LiteralKind() literal_kind {
	return RuneLiteralKind
}

func (e RuneLiteralExpression) Location() common.Location {
	return e.location
}

type BoolLiteralExpression struct {
	Value    bool `json:"value"`
	location common.Location
}

func (e BoolLiteralExpression) Kind() expression_kind {
	return BoolLiteralExpressionKind
}

func (e BoolLiteralExpression) LiteralKind() literal_kind {
	return BoolLiteralKind
}

func (e BoolLiteralExpression) Location() common.Location {
	return e.location
}

type NumberLiteral struct {
	Type  TypeLiteral `json:"type"`
	Value interface{} `json:"value"`
}

type NumberLiteralExpression struct {
	Value    NumberLiteral `json:"value"`
	location common.Location
}

func (e NumberLiteralExpression) Kind() expression_kind {
	return NumberLiteralExpressionKind
}

func (e NumberLiteralExpression) LiteralKind() literal_kind {
	return NumberLiteralKind
}

func (e NumberLiteralExpression) Location() common.Location {
	return e.location
}

type KeyValueEntry struct {
	Key   Expression `json:"key"`
	Value Expression `json:"value"`
}

type ListLiteralExpression struct {
	Value    []KeyValueEntry `json:"value"`
	location common.Location
}

func (e ListLiteralExpression) Kind() expression_kind {
	return ListLiteralExpressionKind
}

func (e ListLiteralExpression) LiteralKind() literal_kind {
	return ListLiteralKind
}

func (e ListLiteralExpression) Location() common.Location {
	return e.location
}

type RecordLiteralExpression struct {
	Value    []KeyValueEntry `json:"value"`
	location common.Location
}

func (e RecordLiteralExpression) Kind() expression_kind {
	return RecordLiteralExpressionKind
}

func (e RecordLiteralExpression) LiteralKind() literal_kind {
	return RecordLiteralKind
}

func (e RecordLiteralExpression) Location() common.Location {
	return e.location
}

type InstanceLiteralExpression struct {
	Type     TypeIdentifier  `json:"type"`
	Value    []KeyValueEntry `json:"value"`
	location common.Location
}

func (e InstanceLiteralExpression) Kind() expression_kind {
	return InstanceLiteralExpressionKind
}

func (e InstanceLiteralExpression) LiteralKind() literal_kind {
	return InstanceLiteralKind
}

func (e InstanceLiteralExpression) Location() common.Location {
	return e.location
}

type CallExpression struct {
	Callee    Expression   `json:"callee"`
	Arguments []Expression `json:"arguments"`
	location  common.Location
}

func (e CallExpression) Kind() expression_kind {
	return CallExpressionKind
}

func (e CallExpression) Location() common.Location {
	return e.location
}

type MemberExpression struct {
	LeftHandSide  Expression           `json:"left_hand_side"`
	RightHandSide IdentifierExpression `json:"right_hand_side"`
	location      common.Location
}

func (e MemberExpression) Kind() expression_kind {
	return MemberExpressionKind
}

func (e MemberExpression) Location() common.Location {
	return e.location
}

type IndexExpression struct {
	Host     Expression `json:"host"`
	Index    Expression `json:"index"`
	location common.Location
}

func (e IndexExpression) Kind() expression_kind {
	return IndexExpressionKind
}

func (e IndexExpression) Location() common.Location {
	return e.location
}

type MatchExpression struct {
	Against   Expression       `json:"against"`
	Blocks    []PredicateBlock `json:"blocks"`
	BaseBlock StatementList    `json:"base_block"`
	location  common.Location
}

func (e MatchExpression) Kind() expression_kind {
	return MatchExpressionKind
}

func (e MatchExpression) Location() common.Location {
	return e.location
}

type GroupExpression struct {
	Expression Expression `json:"expression"`
	location   common.Location
}

func (e GroupExpression) Kind() expression_kind {
	return GroupExpressionKind
}

func (e GroupExpression) Location() common.Location {
	return e.location
}

type ThisExpression struct {
	location common.Location
}

func (e ThisExpression) Kind() expression_kind {
	return ThisExpressionKind
}

func (e ThisExpression) Location() common.Location {
	return e.location
}

type ArithmeticUnaryExpression struct {
	Expression Expression            `json:"expression"`
	Operation  arithmetic_unary_kind `json:"operation"`
	Pre        bool                  `json:"pre"`
	location   common.Location
}

func (e ArithmeticUnaryExpression) Kind() expression_kind {
	return ArithmeticUnaryExpressionKind
}

func (e ArithmeticUnaryExpression) Location() common.Location {
	return e.location
}

type OrExpression struct {
	LeftHandSide  Expression `json:"left_hand_side"`
	RightHandSide Expression `json:"right_hand_side"`
	location      common.Location
}

func (e OrExpression) Kind() expression_kind {
	return OrExpressionKind
}

func (e OrExpression) Location() common.Location {
	return e.location
}

type NotExpression struct {
	Expression Expression
	location   common.Location
}

func (e NotExpression) Kind() expression_kind {
	return NotExpressionKind
}

func (e NotExpression) Location() common.Location {
	return e.location
}
