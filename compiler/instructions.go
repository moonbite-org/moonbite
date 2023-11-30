package compiler

type ValueKind byte
type InstructionKind byte
type ArgumentKind byte
type OperationKind byte

const (
	FalseValueKind ValueKind = 10 + iota
	TrueValueKind
	StringValueKind
	ByteValueKind
	RuneValueKind
	Uint8ValueKind
	Uint16ValueKind
	Uint32ValueKind
	Uint64ValueKind
	Int8ValueKind
	Int16ValueKind
	Int32ValueKind
	Int64ValueKind
	ListValueKind
	MapValueKind
	AliasValueKind
	FunValueKind

	ModuleDefinitionInstructionKind InstructionKind = 40
	SaveInstructionKind             InstructionKind = iota + 32
	CallInstructionKind
	ReturnedCallInstructionKind
	SkipInstructionKind
	SkipIfFalseInstructionKind
	CompareInstructionKind
	ArithmeticInstructionKind
	PushScopeInstructionKind
	PopScopeInstructionKind
	DependInstructionKind
	DefineMapInstructionKind
	RainbowInstructionKind

	PointerArgumentKind ArgumentKind = iota + 40
	ParamArgumentKind
	ValueArgumentKind
	ReturnPointerArgumentKind

	EqualsOperationKind OperationKind = iota + 41
	NotEqualsOperationKind
	LessThanOperationKind
	LessThanOrEqualsOperationKind
	GreaterThanOperationKind
	GreaterThanOrEqualsOperationKind
	AddOperationKind
	SubtractOperationKind
	MultiplyOperationKind
	DivideOperationKind
	ModOperationKind
)

type Instruction interface {
	InstructionKind() InstructionKind
}

type Argument interface {
	ArgumentKind() ArgumentKind
}

type ModuleDefinitionInstruction struct {
	IsEntry bool
	Size    uint32
	Name    string
}

func (i ModuleDefinitionInstruction) InstructionKind() InstructionKind {
	return ModuleDefinitionInstructionKind
}

type SaveInstruction struct {
	Pointer uint32
	Value   interface{}
}

func (i SaveInstruction) InstructionKind() InstructionKind {
	return SaveInstructionKind
}

type CallInstruction struct {
	Pointer   uint32
	Arguments []Argument
}

func (i CallInstruction) InstructionKind() InstructionKind {
	return CallInstructionKind
}

type ReturnedCallInstruction struct {
	Pointer       uint32
	ReturnPointer uint32
	Arguments     []Argument
}

func (i ReturnedCallInstruction) InstructionKind() InstructionKind {
	return ReturnedCallInstructionKind
}

type SkipInstruction struct {
	Count int32
}

func (i SkipInstruction) InstructionKind() InstructionKind {
	return SkipInstructionKind
}

type SkipIfFalseInstruction struct {
	Pointer uint32
	Count   int32
}

func (i SkipIfFalseInstruction) InstructionKind() InstructionKind {
	return SkipIfFalseInstructionKind
}

type CompareInstruction struct {
	Operator      byte
	LeftPointer   uint32
	RightPointer  uint32
	ResultPointer uint32
}

func (i CompareInstruction) InstructionKind() InstructionKind {
	return CompareInstructionKind
}

type ArithmeticInstruction struct {
	Operator      byte
	LeftPointer   uint32
	RightPointer  uint32
	ResultPointer uint32
}

func (i ArithmeticInstruction) InstructionKind() InstructionKind {
	return ArithmeticInstructionKind
}

type PushScopeInstruction struct{}

func (i PushScopeInstruction) InstructionKind() InstructionKind {
	return PushScopeInstructionKind
}

type PopScopeInstruction struct{}

func (i PopScopeInstruction) InstructionKind() InstructionKind {
	return PopScopeInstructionKind
}

type PointerArgument struct {
	Pointer uint32
}

func (a PointerArgument) ArgumentKind() ArgumentKind {
	return PointerArgumentKind
}

type ParamArgument struct {
	Pointer uint32
}

func (a ParamArgument) ArgumentKind() ArgumentKind {
	return ParamArgumentKind
}

type ValueArgument struct {
	Value interface{}
}

func (a ValueArgument) ArgumentKind() ArgumentKind {
	return ValueArgumentKind
}

type ReturnPointerArgument struct{}

func (a ReturnPointerArgument) ArgumentKind() ArgumentKind {
	return ReturnPointerArgumentKind
}
