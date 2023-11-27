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
	"alias":      26,
	"fun":        27,
	"terminator": 0,
}

var InstructionMap = map[string]byte{
	"define_module":    40,
	"save":             50,
	"load_pointer":     51,
	"load_value":       52,
	"load_param":       53,
	"call":             54,
	"call_returned":    55,
	"skip":             56,
	"skip_if_zero":     57,
	"compare":          58,
	"compare_pointers": 59,
	"arithmetic":       60,
}

var CompareInstructions = map[string]byte{
	"equals":     80,
	"not_equals": 81,
	"gt":         82,
	"gte":        83,
	"lt":         84,
	"lte":        85,
}

var ArithmeticInstructions = map[string]byte{
	"add":  90,
	"sub":  91,
	"div":  92,
	"mult": 93,
	"mod":  94,
	"inc":  95,
	"dec":  96,
}
