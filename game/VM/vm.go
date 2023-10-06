package VM

import (
	"fmt"
	"reflect"

	"lucksystem/game/api"
	"lucksystem/game/engine"
	"lucksystem/game/enum"
	"lucksystem/game/operator"
	"lucksystem/game/runtime"
	"lucksystem/script"

	"github.com/golang/glog"
)

type VM struct {
	*runtime.Runtime

	// 游戏对应操作接口
	Operate api.Operator

	Scripts map[string]*script.Script
	// 当前脚本名
	CScriptName string

	// 下一步执行下标
	CNext int
	// 下一步执行偏移
	EIP int
}

func NewVM(opts *Options) *VM {
	vm := &VM{
		Scripts: make(map[string]*script.Script),
	}
	if len(opts.PluginFile) != 0 {
		vm.Operate = operator.NewPlugin(opts.PluginFile)
	} else {
		switch opts.GameName {
		case "LB_EN":
			vm.Operate = operator.NewLB_EN()
		case "SP":
			vm.Operate = operator.NewSP()
		}
	}
	vm.Runtime = runtime.NewRuntime(opts.Mode)
	vm.Operate.Init(vm.Runtime)
	return vm
}

func (vm *VM) LoadScript(scr *script.Script, _switch bool) {
	vm.Scripts[scr.Name] = scr
	if _switch {
		vm.SwitchScript(scr.Name)
	}
}

func (vm *VM) SwitchScript(name string) {
	if scr, ok := vm.Scripts[name]; ok {
		vm.CScriptName = scr.Name
		vm.Runtime.SwitchScript(scr)
	}
}

// 在对EIP修改后调用，查找下一条具体指令，返回指令序号
func (vm *VM) findCode(oldEIP int) int {
	index := vm.CIndex
	_script := vm.Scripts[vm.CScriptName]
	if vm.EIP > oldEIP {
		// 向下查找
		for index < _script.CodeNum && _script.Codes[index].Pos < vm.EIP {
			index++
		}
		if _script.Codes[index].Pos == vm.EIP {
			return index
		} else {
			panic(fmt.Sprintf("未找到跳转位置 [%d]%d -> %d", vm.CIndex, oldEIP, vm.EIP))
		}
	} else if vm.EIP < oldEIP {
		// 向上查找
		for index >= 0 && _script.Codes[index].Pos > vm.EIP {
			index--
		}
		if _script.Codes[index].Pos == vm.EIP {
			return index
		} else {
			panic(fmt.Sprintf("未找到跳转位置 [%d]%d -> %d", vm.CIndex, oldEIP, vm.EIP))
		}
	} else {
		return index
	}

}

func (vm *VM) getNextPos() int {
	if vm.CIndex+1 >= vm.Scripts[vm.CScriptName].CodeNum {
		return 0
	} else {
		return vm.Scripts[vm.CScriptName].Codes[vm.CIndex+1].Pos
	}

}

func (vm *VM) Run() {
	if len(vm.OpcodeMap) == 0 {
		glog.Warning("OPCODE not loaded, import will not be supported")
	}
	glog.V(2).Infoln("Run: ", vm.CScriptName)
	vm.EIP = 0
	vm.CIndex = 0
	vm.CNext = 0

	var in []reflect.Value
	var code *script.CodeLine
	for {
		vm.CIndex = vm.CNext
		code = vm.Scripts[vm.CScriptName].Codes[vm.CIndex]
		opname := vm.Opcode(code.Opcode)
		vm.Runtime.Code().OpStr = opname
		operat := reflect.ValueOf(vm.Operate).MethodByName(opname)
		if operat.IsValid() {
			// 方法已定义，反射调用
			in = make([]reflect.Value, 1)
			in[0] = reflect.ValueOf(vm.Runtime)
		} else {
			// 方法未定义，调用UNDEFINE
			operat = reflect.ValueOf(vm.Operate).MethodByName("UNDEFINED")
			in = make([]reflect.Value, 2)
			in[0] = reflect.ValueOf(vm.Runtime)
			in[1] = reflect.ValueOf(opname)
		}
		glog.V(6).Infof("Index:%d Position:%d \n", vm.CIndex, code.Pos)
		fun := operat.Call(in)  // 反射调用 operator，并返回一个function.HandlerFunc
		next := vm.getNextPos() // 取得下一句位置
		if fun[0].Kind() == reflect.Func {
			eip := 0
			if vm.RunMode == enum.VMRun {
				go fun[0].Interface().(engine.HandlerFunc)() // 调用，默认传递参数列表
				eip = <-vm.Runtime.ChanEIP                   // 取得跳转的位置
			}

			if eip > 0 { // 为0则默认下一句
				next = eip
			}
		}
		glog.V(6).Infoln("\tnext:", next)

		// if next == 0 || opname == "END" { - Many game scripts have an END opcode, but that does not mean that there is nothing else to parse after that opcode.
		if next == 0 {
			break // 结束
		}
		vm.EIP = next
		vm.CNext = vm.findCode(code.Pos)
	}
}
