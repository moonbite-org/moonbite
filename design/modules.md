# Instruction Set Arhcitecture

## Value Serialization

| Byte | Name | Arguments | Description |
| --- | --- | --- | --- |
| 10 | *false* | [] | The literal *false* boolean value. |
| 11 | *true* | [] | The literal *true* boolean value. |
| 12 | *string* | [**list(byte)**, **terminator**] | An encoded *string* value. |
| 13 | *rune* | [**list(byte)**, **terminator**] | An encoded *rune* value. |
| 14 | *byte* | [**byte**] | An encoded *byte* value. |
| 15 | *uint8* | [**{1}list(bytes)**] | An unsigned 8 bit *integer* value. |
| 16 | *uint16* | [**{2}list(bytes)**] | An unsigned 16 bit *integer* value. |
| 17 | *uint32* | [**{4}list(bytes)**] | An unsigned 32 bit *integer* value. |
| 18 | *uint64* | [**{8}list(bytes)**] | An unsigned 64 bit *integer* value. |
| 19 | *int8* | [**{1}list(bytes)**] | A signed 8 bit *integer* value. |
| 20 | *int16* | [**{2}list(bytes)**] | A signed 16 bit *integer* value. |
| 21 | *int32* | [**{4}list(bytes)**] | A signed 32 bit *integer* value. |
| 22 | *int64* | [**{8}list(bytes)**] | A signed 64 bit *integer* value. |
| 23 | *list* | [**list(unknown)**, **terminator**] | A list of other serialized values. |
| 24 | *map* | [**list({string:unknown})**, **terminator**] | A list of string to serialized value pairs. |
| 25 | *alias* | [**list({string:unknown})**, **terminator**] | A list of string to serialized value pairs. This represents an aliased type. When reached, provides the value kept in the special *raw* key. |
| 26 | *fun* | [**int32**,**list(byte)**, **terminator**] | A function. The first argument indicates the size of its content after that point. The last argument is the instruction list. |


## Instructions

| Byte | Operation | Arguments | Description |
| --- | --- | --- | --- |
| 40 | *define_module* | [**bool**, **uint32**, **string**] | Defines a module. First argument indicates if this modules is the entry module, second argument indicates the size of this module after this argument, the last arguments indicates its name. |
| 50 | *save* | [**uint32**, **unknown**] | Saves a value to current context. The first argument is the pointer to this value and the second is a serializable unknown value. |
| 51 | *call* | [**uint32**, **list(argument)**] | Calls a function at pointer provided by the first argument. Passes the second argument as the parameters to the function. |
| 52 | *returned_call* | [**uint32**, **uint32**, **list(argument)**] | Calls a function at pointer provided by the first argument. Provides a second pointer argument for the return value to be saved and passes the last argument as the parameters to the function. |
| 53 | *skip* | [**int32**] | Skips given amount of instructions. A negative value skips backwards. |
| 54 | *skip_if_false* | [**uint32**, **int32**] | Skips given amount of instructions if the pointer at the first argument is a serialzed false value. A negative value skips backwards. |
| 55 | *compare* | [**operator**, **uint32**, **uint32**, **uint32**] | Compares the values at the second and third arguments which are pointers. The first argument is the comparing operator which is a special byte. Writes the result at the last argument pointer as a serialized bool value. |
| 56 | *arithmetic* | [**operator**, **uint32**, **uint32**, **uint32**] | Does arithmetic between two pointers with give operator and saves the result to last argument pointer. |
| 57 | *push_scope* | [] | Pushes a new scope to the stack for things like functions and if/match/loop statements. |
| 58 | *pop_scope* | [] | Pops a context from the stack losing all the values stored inside. |
| 59 | *depend* | [**string**] | Creates a dependency for the module, inticating that this module uses another. |
| 60 | *define_map_descriptor* | [] | Defines where the keys and methods are for a give map type. |
| 61 | *rainbow* | [] | Custom request. |

## Arguments

| Byte| Name | Arguments | Description |
| --- | --- | --- | --- |
| 70 | *pointer_argument* | [**uint32**] | Pointer to a value. |
| 71 | *param_argument* | [**uint32**] | Index pointer of a function parameter. |
| 72 | *value_argument* | [**unknown**] | A direct value. |
| 73 | *return_argument* | [] | An argument that replaces the return pointer |
| 74 | *description_argument* | [**uint32**, **uint32**] | Returns a description of a map. First argument is the pointer to the value and the last argument is the descriptor. |


## Compare & Arithmetic Operations

| Byte | Name | Description |
| --- | --- | --- |
| 75 | *equals* | A value equals to another value. This never checks the pointers but rather the encoded value. |
| 76 | *not_equals* | Two values are not equal to each other, works like equals. |
| 77 | *less_than* | A numeric value is less than another numeric value. |
| 78 | *less_than_or_equals* | A numeric value is less than or equal to another numeric value. |
| 79 | *greater_than* | A numeric value is greater than another numeric value. |
| 80 | *greater_than_or_equals* | A numeric value is greater than or equal to another numeric value. |
| 81 | *add* | Adds two numbers together. |
| 82 | *subtract* | Subtracts right from left. |
| 83 | *multiply* | Multiples two numbers together. |
| 84 | *divide* | Divides right to left. |
| 85 | *mod* | Takes a mod right of the left value. |

# Example Programs

```
package main

use syscall

type Data {
  value byte
}

const data = Data{
  value: 65
}

fun main() {
  syscall.write(1, data.value)
}
```

```
// use syscall
define_module true int32 SIZE string main
depend string syscall
push_context
define_map_descriptor 
  string main.Data
    string data int32 0
save int32 200 byte 65
save int32 100
  map string main.Data
    int32 0 int32 200
save int32 101 // save at pointer 100
  fun int32 SIZE
    push_context
      call int32 11 list value int32 1 description 100 0 terminator
    pop_context
  terminator
call int32 101 list terminator // call main with no params
pop_context

```