package parser

import (
	"fmt"
	"strings"
)

type statement_kind string
type expression_kind string
type type_kind string
type var_kind string
type loop_kind string
type literal_kind string

const (
	// statements
	package_statement                statement_kind = "statement:package"
	use_statement                    statement_kind = "statement:use"
	return_statement                 statement_kind = "statement:return"
	declaration_statement            statement_kind = "statement:declaration"
	type_definition_statement        statement_kind = "statement:type-definition"
	trait_definition_statement       statement_kind = "statement:trait-definition"
	unbound_fun_definition_statement statement_kind = "statement:unbound-fun-definition"
	bound_fun_definition_statement   statement_kind = "statement:bound-fun-definition"
	expression_statement             statement_kind = "statement:expression"
	loop_statement                   statement_kind = "statement:loop"
	if_statement                     statement_kind = "statement:if"
	or_statement                     statement_kind = "statement:or"
	single_line_comment_statement    statement_kind = "statement:single_line_comment"
	multi_line_comment_statement     statement_kind = "statement:multi_line_comment"

	// expressions
	identifier_expression expression_kind = "expression:identifier"
	arithmetic_expression expression_kind = "expression:arithmetic"
	binary_expression     expression_kind = "expression:binary"
	call_expression       expression_kind = "expression:call"
	member_expression     expression_kind = "expression:member"
	match_expression      expression_kind = "expression:match"
	type_cast_expression  expression_kind = "expression:type-cast"
	caret_expression      expression_kind = "expression:caret"
	instanceof_expression expression_kind = "expression:instanceof"
	match_self_expression expression_kind = "expression:match_self"
	group_expression      expression_kind = "expression:group"

	// literal expressions
	string_literal_expression   expression_kind = "expression:string-literal"
	rune_literal_expression     expression_kind = "expression:rune-literal"
	bool_literal_expression     expression_kind = "expression:bool-literal"
	number_literal_expression   expression_kind = "expression:number-literal"
	list_literal_expression     expression_kind = "expression:list-literal"
	record_literal_expression   expression_kind = "expression:record-literal"
	instance_literal_expression expression_kind = "expression:instance-literal"

	// types
	type_identifier type_kind = "type:type-identifier"
	struct_literal  type_kind = "type:struct-literal"
	operated_type   type_kind = "type:operated-type"
	typed_literal   type_kind = "type:typed-literal"
	group_type      type_kind = "type:group"

	// vars
	variable var_kind = "var"
	constant var_kind = "const"

	// loops
	unipartite_loop loop_kind = "predicate:unipartite"
	bipartite_loop  loop_kind = "predicate:bipartite"
	tripartite_loop loop_kind = "predicate:tripartite"

	// literals
	string_literal_kind   literal_kind = "literal:string"
	bool_literal_kind     literal_kind = "literal:bool"
	rune_literal_kind     literal_kind = "literal:rune"
	number_literal_kind   literal_kind = "literal:number"
	list_literal_kind     literal_kind = "literal:list"
	record_literal_kind   literal_kind = "literal:record"
	instance_literal_kind literal_kind = "literal:instance"
)

const tab = "    "

type printable interface {
	String() string
}

func stringify_list[T printable](items []T, open, close string, keepempty bool, seperator string) string {
	result := []string{}
	literal := ""

	for _, item := range items {
		result = append(result, item.String())
	}

	if len(items) != 0 {
		literal += open
		literal += strings.Join(result, seperator)
		literal += close
	} else {
		if keepempty {
			literal += open
			literal += close
		}
	}

	return literal
}

type ConstrainedType struct {
	Name       Expression   `json:"name"`
	Constraint *TypeLiteral `json:"constraint"`
}

func (t ConstrainedType) String() string {
	c := ""

	if t.Constraint != nil {
		_c := *t.Constraint
		c = fmt.Sprintf(" %s", _c.String())
	}

	return fmt.Sprintf("%s%s", t.Name.String(), c)
}

type TypeLiteral interface {
	Kind() type_kind
	printable
}

type TypedLiteral struct {
	Type    TypeIdentifier    `json:"type"`
	Literal LiteralExpression `json:"literal"`
}

func (t TypedLiteral) Kind() type_kind {
	return typed_literal
}

func (t TypedLiteral) String() string {
	return fmt.Sprintf("%s(%s)", t.Type.String(), t.Literal.String())
}

type GroupType struct {
	Type TypeLiteral `json:"type"`
}

func (t GroupType) Kind() type_kind {
	return group_type
}

func (t GroupType) String() string {
	return fmt.Sprintf("(%s)", t.Type.String())
}

type TypeIdentifier struct {
	Name     Expression    `json:"name"`
	Generics []TypeLiteral `json:"generics"`
}

func (t TypeIdentifier) Kind() type_kind {
	return type_identifier
}

func (t TypeIdentifier) String() string {
	return fmt.Sprintf("%s%s", t.Name.String(), stringify_list(t.Generics, "<", ">", false, ", "))
}

type OperatedType struct {
	LeftHandSide  TypeLiteral `json:"left_hand_side"`
	RightHandSide TypeLiteral `json:"right_hand_side"`
	Operator      string
}

func (t OperatedType) Kind() type_kind {
	return operated_type
}

func (t OperatedType) String() string {
	return fmt.Sprintf("%s %s %s", t.LeftHandSide.String(), t.Operator, t.RightHandSide.String())
}

type ValueTypePair struct {
	Key  IdentifierExpression `json:"key"`
	Type TypeLiteral          `json:"type"`
}

type StructLiteral []ValueTypePair

func (t StructLiteral) Kind() type_kind {
	return struct_literal
}

func (t StructLiteral) String() string {
	if len(t) == 0 {
		return "{}"
	}

	result := ""

	for _, entry := range t {
		result += fmt.Sprintf("%s%s %s;\n", tab, entry.Key.String(), entry.Type.String())
	}

	return fmt.Sprintf("{\n%s}", result)
}

type TypedParameter struct {
	Name IdentifierExpression `json:"name"`
	Type TypeLiteral          `json:"type"`
}

func (p TypedParameter) String() string {
	return fmt.Sprintf("%s %s", p.Name.String(), p.Type.String())
}

type FunctionSignature interface {
	is_fun_signature() bool
}

type UnboundFunctionSignature struct {
	Name       IdentifierExpression `json:"name"`
	Parameters []TypedParameter     `json:"parameters"`
	Generics   []ConstrainedType    `json:"generics"`
	ReturnType *TypeLiteral         `json:"return_type"`
	location   Location
}

func (s UnboundFunctionSignature) String() string {
	r_type := ""

	if s.ReturnType != nil {
		rt := *s.ReturnType
		r_type = fmt.Sprintf(" %s", rt.String())
	}

	return fmt.Sprintf("fun %s%s%s%s", s.Name.String(), stringify_list(s.Generics, "<", ">", false, ", "), stringify_list(s.Parameters, "(", ")", true, ", "), r_type)
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
	location   Location
}

func (s BoundFunctionSignature) String() string {
	r_type := ""

	if s.ReturnType != nil {
		rt := *s.ReturnType
		r_type = fmt.Sprintf(" %s", rt.String())
	}

	return fmt.Sprintf("fun for %s %s%s%s%s", s.For.String(), s.Name.String(), stringify_list(s.Generics, "<", ">", false, ", "), stringify_list(s.Parameters, "(", ")", true, ", "), r_type)
}

func (s BoundFunctionSignature) is_fun_signature() bool {
	return true
}

type Statement interface {
	Kind() statement_kind
	Location() Location
	printable
}

type Expression interface {
	Kind() expression_kind
	Location() Location
	printable
}

type Definition interface {
	is_definition() bool
	printable
}

type Comment interface {
	is_comment() bool
	Definition
	printable
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
	location   Location
}

func (s ExpressionStatement) Kind() statement_kind {
	return expression_statement
}

func (s ExpressionStatement) String() string {
	return s.Expression.String()
}

func (s ExpressionStatement) Location() Location {
	return s.location
}

type DeclarationStatement struct {
	VarKind  var_kind             `json:"var_kind"`
	Name     IdentifierExpression `json:"name"`
	Type     *TypeLiteral         `json:"type"`
	Value    *Expression          `json:"value"`
	location Location
}

func (s DeclarationStatement) is_definition() bool {
	return true
}

func (s DeclarationStatement) Kind() statement_kind {
	return declaration_statement
}

func (s DeclarationStatement) String() string {
	typ := ""
	val := ""

	if s.Type != nil {
		typ = fmt.Sprintf(" %s", (*s.Type).String())
	}

	if s.Value != nil {
		val = (*s.Value).String()
	}

	return fmt.Sprintf("%s%s %s = %s", s.VarKind, typ, s.Name, val)
}

func (s DeclarationStatement) Location() Location {
	return s.location
}

type PackageStatement struct {
	Name     IdentifierExpression `json:"name"`
	location Location
}

func (s PackageStatement) Kind() statement_kind {
	return package_statement
}

func (s PackageStatement) String() string {
	return fmt.Sprintf("package %s", s.Name.Value)
}

func (s PackageStatement) Location() Location {
	return s.location
}

type UseStatement struct {
	Resource IdentifierExpression  `json:"resource"`
	As       *IdentifierExpression `json:"as"`
	location Location
}

func (s UseStatement) Kind() statement_kind {
	return use_statement
}

func (s UseStatement) String() string {
	postfix := ""

	if s.As != nil {
		postfix = fmt.Sprintf(" as %s", s.As.String())
	}

	return fmt.Sprintf("use %s%s", s.Resource.String(), postfix)
}

func (s UseStatement) Location() Location {
	return s.location
}

type TypeDefinitionStatement struct {
	Name            IdentifierExpression `json:"name"`
	Generics        []ConstrainedType    `json:"generics"`
	Implementations []TypeIdentifier     `json:"implementations"`
	Definition      TypeLiteral          `json:"definiton"`
	location        Location
}

func (s TypeDefinitionStatement) is_definition() bool {
	return true
}

func (s TypeDefinitionStatement) Kind() statement_kind {
	return type_definition_statement
}

func (s TypeDefinitionStatement) String() string {
	implementations := stringify_list(s.Implementations, "[", "]", false, ", ")

	if len(implementations) != 0 {
		implementations = fmt.Sprintf(" implements %s", implementations)
	}

	return fmt.Sprintf("type %s%s%s %s", s.Name.String(), stringify_list(s.Generics, "<", ">", false, ", "), implementations, s.Definition)
}

func (s TypeDefinitionStatement) Location() Location {
	return s.location
}

type TraitDefinitionStatement struct {
	Name       IdentifierExpression       `json:"name"`
	Generics   []ConstrainedType          `json:"generics"`
	Mimics     []TypeIdentifier           `json:"mimics"`
	Definition []UnboundFunctionSignature `json:"definition"`
	location   Location
}

func (s TraitDefinitionStatement) is_definition() bool {
	return true
}

func (s TraitDefinitionStatement) Kind() statement_kind {
	return trait_definition_statement
}

func (s TraitDefinitionStatement) String() string {
	mimics := stringify_list(s.Mimics, "[", "]", false, ",")

	if len(mimics) != 0 {
		mimics = fmt.Sprintf(" mimics %s", mimics)
	}

	body := []string{}

	for _, definition := range s.Definition {
		body = append(body, fmt.Sprintf("%s%s;\n", tab, definition.String()))
	}

	return fmt.Sprintf("trait %s%s%s {\n%s}", s.Name.String(), stringify_list(s.Generics, "<", ">", false, ", "), mimics, strings.Join(body, ""))
}

func (s TraitDefinitionStatement) Location() Location {
	return s.location
}

type FunDefinitionStatement interface {
	set_body(body []Statement)
	Definition
}

type UnboundFunDefinitionStatement struct {
	Signature UnboundFunctionSignature `json:"signature"`
	Body      []Statement              `json:"body"`
	location  Location
}

func (s UnboundFunDefinitionStatement) is_definition() bool {
	return true
}

func (s *UnboundFunDefinitionStatement) set_body(body []Statement) {
	s.Body = body
}

func (s UnboundFunDefinitionStatement) Kind() statement_kind {
	return unbound_fun_definition_statement
}

func (s UnboundFunDefinitionStatement) String() string {
	return fmt.Sprintf("%s %s", s.Signature.String(), stringify_list(s.Body, "{", "}", true, "\n"))
}

func (s UnboundFunDefinitionStatement) Location() Location {
	return s.location
}

type BoundFunDefinitionStatement struct {
	Signature BoundFunctionSignature `json:"signature"`
	Body      []Statement            `json:"body"`
	location  Location
}

func (s BoundFunDefinitionStatement) is_definition() bool {
	return true
}

func (s *BoundFunDefinitionStatement) set_body(body []Statement) {
	s.Body = body
}

func (s BoundFunDefinitionStatement) Kind() statement_kind {
	return bound_fun_definition_statement
}

func (s BoundFunDefinitionStatement) String() string {
	return fmt.Sprintf("%s %s", s.Signature.String(), stringify_list(s.Body, "{", "}", true, "\n"))
}

func (s BoundFunDefinitionStatement) Location() Location {
	return s.location
}

type ReturnStatement struct {
	Value    *Expression `json:"expression"`
	location Location
}

func (s ReturnStatement) Kind() statement_kind {
	return return_statement
}

func (s ReturnStatement) String() string {
	end := ""

	if *s.Value != nil {
		v_p := *s.Value
		end = fmt.Sprintf(" %s", v_p.String())
	}

	return fmt.Sprintf("return%s", end)
}

func (s ReturnStatement) Location() Location {
	return s.location
}

type LoopPredicate interface {
	LoopKind() loop_kind
}

type UnipartiteLoopPredicate struct {
	Expression Expression `json:"expression"`
}

func (l UnipartiteLoopPredicate) LoopKind() loop_kind {
	return unipartite_loop
}

type BipartiteLoopPredicate struct {
	Key      *IdentifierExpression `json:"key"`
	Value    *IdentifierExpression `json:"value"`
	Iterator Expression            `json:"iterator"`
}

func (l BipartiteLoopPredicate) LoopKind() loop_kind {
	return bipartite_loop
}

type TripartiteLoopPredicate struct {
	Declaration *DeclarationStatement `json:"declaration"`
	Predicate   Expression            `json:"predicate"`
	Procedure   *Expression           `json:"procedure"`
}

func (l TripartiteLoopPredicate) LoopKind() loop_kind {
	return tripartite_loop
}

type LoopStatement struct {
	Predicate LoopPredicate `json:"predicate"`
	Body      []Statement   `json:"body"`
	location  Location
}

func (s LoopStatement) Kind() statement_kind {
	return loop_statement
}

func (s LoopStatement) Location() Location {
	return s.location
}

type PredicateBlock struct {
	Predicate Expression  `json:"predicate"`
	Body      []Statement `json:"body"`
}

type IfStatement struct {
	MainBlock    PredicateBlock   `json:"main_block"`
	ElseIfBlocks []PredicateBlock `json:"else_if_blocks"`
	ElseBlock    []Statement      `json:"else_block"`
	location     Location
}

func (s IfStatement) Kind() statement_kind {
	return if_statement
}

func (s IfStatement) Location() Location {
	return s.location
}

type OrStatement struct {
	Try      Expression `json:"try"`
	Fail     Expression `json:"fail"`
	location Location
}

func (s OrStatement) Kind() statement_kind {
	return or_statement
}

func (s OrStatement) String() string {
	return fmt.Sprintf("%s or %s", s.Try.String(), s.Fail.String())
}

func (s OrStatement) Location() Location {
	return s.location
}

type SingleLineCommentStatement struct {
	Comment  string `json:"comment"`
	location Location
}

func (s SingleLineCommentStatement) is_comment() bool {
	return true
}
func (s SingleLineCommentStatement) is_definition() bool {
	return true
}

func (s SingleLineCommentStatement) Kind() statement_kind {
	return single_line_comment_statement
}

func (s SingleLineCommentStatement) String() string {
	return s.Comment
}

func (s SingleLineCommentStatement) Location() Location {
	return s.location
}

type MultiLineCommentStatement struct {
	Comment  string `json:"comment"`
	location Location
}

func (s MultiLineCommentStatement) is_comment() bool {
	return true
}

func (s MultiLineCommentStatement) is_definition() bool {
	return true
}

func (s MultiLineCommentStatement) Kind() statement_kind {
	return multi_line_comment_statement
}

func (s MultiLineCommentStatement) String() string {
	return s.Comment
}

func (s MultiLineCommentStatement) Location() Location {
	return s.location
}

// EXPRESSIONS

type IdentifierExpression struct {
	Value    string `json:"value"`
	location Location
}

func (e IdentifierExpression) Kind() expression_kind {
	return identifier_expression
}

func (e IdentifierExpression) String() string {
	return e.Value
}

func (e IdentifierExpression) Location() Location {
	return e.location
}

type CaretExpression struct {
	location Location
}

func (e CaretExpression) Kind() expression_kind {
	return caret_expression
}

func (e CaretExpression) Location() Location {
	return e.location
}

type TypeCastExpression struct {
	Value    Expression     `json:"value"`
	Type     TypeIdentifier `json:"type"`
	location Location
}

func (e TypeCastExpression) Kind() expression_kind {
	return type_cast_expression
}

func (e TypeCastExpression) String() string {
	return fmt.Sprintf("%s.(%s)", e.Value.String(), e.Type.String())
}

func (e TypeCastExpression) Location() Location {
	return e.location
}

type InstanceofExpression struct {
	LeftHandSide  Expression  `json:"left_hand_side"`
	RightHandSide TypeLiteral `json:"right_hand_side"`
	location      Location
}

func (e InstanceofExpression) Kind() expression_kind {
	return instanceof_expression
}

func (e InstanceofExpression) Location() Location {
	return e.location
}

type MatchSelfExpression struct {
	location Location
}

func (e MatchSelfExpression) Kind() expression_kind {
	return match_self_expression
}

func (e MatchSelfExpression) String() string {
	return "."
}

func (e MatchSelfExpression) Location() Location {
	return e.location
}

type ArithmeticExpression struct {
	LeftHandSide  Expression `json:"left_hand_side"`
	RightHandSide Expression `json:"right_hand_side"`
	Operator      string     `json:"operator"`
	location      Location
}

func (e ArithmeticExpression) Kind() expression_kind {
	return arithmetic_expression
}

func (e ArithmeticExpression) String() string {
	return fmt.Sprintf("%s %s %s", e.LeftHandSide.String(), e.Operator, e.RightHandSide.String())
}

func (e ArithmeticExpression) Location() Location {
	return e.location
}

type BinaryExpression struct {
	LeftHandSide  Expression `json:"left_hand_side"`
	RightHandSide Expression `json:"right_hand_side"`
	Operator      string     `json:"operator"`
	location      Location
}

func (e BinaryExpression) Kind() expression_kind {
	return binary_expression
}

func (e BinaryExpression) String() string {
	return fmt.Sprintf("%s %s %s", e.LeftHandSide.String(), e.Operator, e.RightHandSide.String())
}

func (e BinaryExpression) Location() Location {
	return e.location
}

type LiteralExpression interface {
	Expression
	LiteralKind() literal_kind
}

type StringLiteralExpression struct {
	Value    string `json:"value"`
	location Location
}

func (e StringLiteralExpression) Kind() expression_kind {
	return string_literal_expression
}

func (e StringLiteralExpression) String() string {
	return fmt.Sprintf("\"%s\"", e.Value)
}

func (e StringLiteralExpression) LiteralKind() literal_kind {
	return string_literal_kind
}

func (e StringLiteralExpression) Location() Location {
	return e.location
}

type RuneLiteralExpression struct {
	Value    rune `json:"value"`
	location Location
}

func (e RuneLiteralExpression) Kind() expression_kind {
	return rune_literal_expression
}

func (e RuneLiteralExpression) String() string {
	return fmt.Sprintf("'%s'", string(e.Value))
}

func (e RuneLiteralExpression) LiteralKind() literal_kind {
	return rune_literal_kind
}

func (e RuneLiteralExpression) Location() Location {
	return e.location
}

type BoolLiteralExpression struct {
	Value    bool `json:"value"`
	location Location
}

func (e BoolLiteralExpression) Kind() expression_kind {
	return bool_literal_expression
}

func (e BoolLiteralExpression) String() string {
	if e.Value {
		return "true"
	}

	return "false"
}

func (e BoolLiteralExpression) LiteralKind() literal_kind {
	return bool_literal_kind
}

func (e BoolLiteralExpression) Location() Location {
	return e.location
}

type NumberLiteral interface {
	Type() interface{}
	Value() interface{}
}

type NumberLiteralExpression struct {
	Value    NumberLiteral `json:"value"`
	location Location
}

func (e NumberLiteralExpression) Kind() expression_kind {
	return number_literal_expression
}

func (e NumberLiteralExpression) String() string {
	return fmt.Sprintf("%d", e.Value)
}

func (e NumberLiteralExpression) LiteralKind() literal_kind {
	return number_literal_kind
}

func (e NumberLiteralExpression) Location() Location {
	return e.location
}

type KeyValueEntry struct {
	Key   Expression `json:"key"`
	Value Expression `json:"value"`
}

type ListLiteralExpression struct {
	Value    []KeyValueEntry `json:"value"`
	location Location
}

func (e ListLiteralExpression) Kind() expression_kind {
	return list_literal_expression
}

func (e ListLiteralExpression) String() string {
	values := []printable{}

	for _, value := range e.Value {
		values = append(values, value.Value)
	}

	return stringify_list(values, "[", "]", true, ", ")
}

func (e ListLiteralExpression) LiteralKind() literal_kind {
	return list_literal_kind
}

func (e ListLiteralExpression) Location() Location {
	return e.location
}

type RecordLiteralExpression struct {
	Value    []KeyValueEntry `json:"value"`
	location Location
}

func (e RecordLiteralExpression) Kind() expression_kind {
	return record_literal_expression
}

func (e RecordLiteralExpression) LiteralKind() literal_kind {
	return record_literal_kind
}

func (e RecordLiteralExpression) Location() Location {
	return e.location
}

type InstanceLiteralExpression struct {
	Type     TypeIdentifier  `json:"tpye"`
	Value    []KeyValueEntry `json:"value"`
	location Location
}

func (e InstanceLiteralExpression) Kind() expression_kind {
	return instance_literal_expression
}

func (e InstanceLiteralExpression) LiteralKind() literal_kind {
	return instance_literal_kind
}

func (e InstanceLiteralExpression) Location() Location {
	return e.location
}

type CallExpression struct {
	Callee    Expression   `json:"callee"`
	Arguments []Expression `json:"arguments"`
	location  Location
}

func (e CallExpression) Kind() expression_kind {
	return call_expression
}

func (e CallExpression) String() string {
	return fmt.Sprintf("%s%s", e.Callee.String(), stringify_list(e.Arguments, "(", ")", true, ", "))
}

func (e CallExpression) Location() Location {
	return e.location
}

type MemberExpression struct {
	LeftHandSide  Expression           `json:"left_hand_side"`
	RightHandSide IdentifierExpression `json:"right_hand_side"`
	location      Location
}

func (e MemberExpression) Kind() expression_kind {
	return member_expression
}

func (e MemberExpression) String() string {
	return fmt.Sprintf("%s.%s", e.LeftHandSide.String(), e.RightHandSide.String())
}

func (e MemberExpression) Location() Location {
	return e.location
}

type MatchExpression struct {
	Blocks    []PredicateBlock `json:"blocks"`
	BaseBlock []Statement      `json:"base_block"`
	location  Location
}

func (e MatchExpression) Kind() expression_kind {
	return match_expression
}

func (e MatchExpression) Location() Location {
	return e.location
}

type GroupExpression struct {
	Expression Expression `json:"expression"`
	location   Location
}

func (e GroupExpression) Kind() expression_kind {
	return group_expression
}

func (e GroupExpression) String() string {
	return fmt.Sprintf("(%s)", e.Expression.String())
}

func (e GroupExpression) Location() Location {
	return e.location
}
