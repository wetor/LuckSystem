package VM

import (
	"fmt"
	"lucksystem/game/context"
	"lucksystem/game/engine"
	"lucksystem/game/enum"
	"lucksystem/game/operater"
	"lucksystem/game/variable"
	"lucksystem/script"
	"os"
	"reflect"
	"strings"

	"github.com/golang/glog"
)

type VM struct {
	*context.Context

	// Opcode map
	OpcodeMap map[uint8]string

	// 游戏对应操作接口
	Operate operater.Operater

	// 下一步执行偏移
	EIP int
}

func NewVM(_script *script.ScriptFile, mode enum.VMRunMode) *VM {
	vm := &VM{}
	switch _script.GameName {
	case "LB_EN":
		vm.Operate = operater.GetLB_EN()

	case "SP":
		vm.Operate = operater.GetSP()
	}
	vm.Context = &context.Context{
		Engine:   &engine.Engine{},
		Scripts:  make(map[string]*script.ScriptFile),
		KeyPress: make(chan int),
		ChanEIP:  make(chan int),
		RunMode:  mode,
	}
	vm.Variable = &variable.VariableStore{}
	vm.Variable.Init()
	vm.CScriptName = _script.Name
	vm.Scripts[_script.Name] = _script
	return vm
}

// 在对EIP修改后调用，查找下一条具体指令，返回指令序号
func (vm *VM) findCode(oldEIP int) int {
	index := vm.CIndex
	_script := vm.Script()
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
	if vm.CIndex+1 >= vm.Script().CodeNum {
		return 0
	} else {
		return vm.Script().Codes[vm.CIndex+1].Pos
	}

}

func (vm *VM) Run() {
	vm.EIP = 0
	vm.CIndex = 0

	var in []reflect.Value
	var code *script.CodeLine
	for {
		vm.CIndex = vm.CNext
		code = vm.Script().Codes[vm.CIndex]
		opname, ok := vm.OpcodeMap[code.Opcode]
		if !ok {
			glog.V(5).Infoln(vm.CIndex, "Opcode不存在", code.Opcode)
			opname = fmt.Sprintf("0x%X", code.Opcode)
			//vm.CNext++
			//continue
		}
		operat := reflect.ValueOf(vm.Operate).MethodByName(opname)
		if operat.IsValid() {
			// 方法已定义，反射调用
			in = make([]reflect.Value, 1)
			in[0] = reflect.ValueOf(vm.Context)
		} else {
			glog.V(5).Infoln(vm.CIndex, "Operation不存在", opname)
			// 方法未定义，调用UNDEFINE
			operat = reflect.ValueOf(vm.Operate).MethodByName("UNDEFINE")
			in = make([]reflect.Value, 2)
			in[0] = reflect.ValueOf(vm.Context)
			in[1] = reflect.ValueOf(opname)
		}
		fun := operat.Call(in) // 反射调用 operater，并返回一个function.HandlerFunc

		next := vm.getNextPos() // 取得下一句位置
		glog.V(4).Infof("Index:%d Position:%d \n", vm.CIndex, code.Pos)
		if fun[0].Kind() == reflect.Func {
			eip := 0
			if vm.RunMode == enum.VMRun {
				go fun[0].Interface().(engine.HandlerFunc)() // 调用，默认传递参数列表
				eip = <-vm.Context.ChanEIP                   // 取得跳转的位置
			}

			if eip > 0 { // 为0则默认下一句
				next = eip
			}
		}
		glog.V(4).Infoln("\tnext:", next)

		// if next == 0 || opname == "END" { - Many game scripts have an END opcode, but that does not mean that there is nothing else to parse after that opcode.
		if next == 0 {

			break // 结束
		}
		vm.EIP = next
		vm.CNext = vm.findCode(code.Pos)
	}
}
func (vm *VM) LoadOpcode(path string) error {
	data, err := os.ReadFile(path)
	if err != nil {
		glog.V(8).Infof("os.ReadFile", err)
		return err
	}
	strlines := strings.Split(string(data), "\n")
	vm.OpcodeMap = make(map[uint8]string, len(strlines)+1)
	for i, line := range strlines {
		line = strings.Replace(line, "\r", "", -1)
		line = strings.Replace(line, "\n", "", -1)
		line = strings.Replace(line, " ", "", -1)
		vm.OpcodeMap[uint8(i)] = line
	}
	return nil
}
