package main

func main() {
	// start := time.Now()
	// bin := compiler.Binary{}

	// mod := compiler.NewModule("main", true)
	// fd := mod.Save(mod.NextPointer(), 0)
	// mod.PushContext()
	// mod.Call(11, -1, []compiler.Argument{compiler.PointerArgument{Value: fd}, compiler.ParamArgument{Value: 0}})
	// log := mod.Save(mod.NextPointer(), compiler.CreateFun(mod.PopContext()))
	// console := map[string]interface{}{
	// 	"fd":  fd,
	// 	"log": log,
	// }
	// mod.Save(mod.NextPointer(), console)
	// mod.PushContext()
	// m := mod.Save(mod.NextPointer(), 65)
	// mod.Call(log, -1, []compiler.Argument{compiler.PointerArgument{Value: m}})
	// main := mod.Save(mod.NextPointer(), compiler.CreateFun(mod.PopContext()))
	// mod.Call(main, -1, []compiler.Argument{})
	// bin.RegisterModule(mod)

	// compiled, err := bin.Compile()

	// if err != nil {
	// 	panic(err)
	// }

	// // fmt.Println(compiled)
	// os.WriteFile("main.mbin", compiled, 0644)
	// // compiled, _ := os.ReadFile("main.mbin")
	// runner := vm.NewVm(compiled)
	// fmt.Println(runner.Run())
	// fmt.Println(time.Since(start))
}
