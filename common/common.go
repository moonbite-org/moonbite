package common

var TypeMap = map[string]byte{
	"false":      10,
	"true":       11,
	"string":     12,
	"rune":       13,
	"int8":       14,
	"int16":      15,
	"int32":      16,
	"int64":      17,
	"uint8":      18,
	"uint16":     19,
	"uint32":     20,
	"uint64":     21,
	"float32":    22,
	"float64":    23,
	"list":       24,
	"struct":     25,
	"fun":        26,
	"terminator": 0,
}

var InstructionMap = map[string]byte{
	"define_module":    10,
	"save":             20,
	"load_pointer":     21,
	"load_value":       22,
	"load_param":       23,
	"call":             24,
	"skip":             25,
	"skip_if":          26,
	"compare":          27,
	"compare_pointers": 27,
	"arithmetic":       28,
}

var CompareInstructions = map[string]byte{
	"equals":     10,
	"not_equals": 11,
	"gt":         12,
	"gte":        13,
	"lt":         14,
	"lte":        15,
}

var ArithmeticInstructions = map[string]byte{
	"add":  10,
	"sub":  11,
	"div":  12,
	"mult": 13,
	"mod":  14,
	"inc":  15,
	"dec":  16,
}
