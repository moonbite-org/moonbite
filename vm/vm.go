package vm

import (
	"fmt"

	"github.com/moonbite-org/moonbite/common"
)

type source_type string

const Memory source_type = "memory"
const Argument source_type = "argument"

type Vm struct {
	offset int
	input  []byte
}

func (vm Vm) Run() error {
	archive, err := vm.resolve_string()

	if err != nil {
		return err
	}
	if archive != common.Config.ArchiveName {
		return fmt.Errorf("broken archive")
	}

	_, err = vm.resolve_string()

	if err != nil {
		return err
	}

	for vm.offset < len(vm.input) {
		// fmt.Println(vm.current())
		vm.advance()
	}

	return nil
}

func (vm *Vm) advance() {
	if vm.offset+1 <= len(vm.input) {
		vm.offset++
	}
}

func (vm *Vm) current() byte {
	if vm.offset >= len(vm.input) {
		return 255
	}

	return vm.input[vm.offset]
}

func (vm *Vm) resolve_string() (string, error) {
	if vm.current() != common.TypeMap["string"] {
		return "", fmt.Errorf("malformed string at %d", vm.offset)
	}

	result := ""

	vm.advance()
	for vm.current() != common.TypeMap["terminator"] && vm.current() != 255 {
		result += string(vm.current())
		vm.advance()
	}
	vm.advance()

	return result, nil
}

func NewVm(input []byte) Vm {
	return Vm{
		input: input,
	}
}
