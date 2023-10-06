package operator

import (
	"strings"

	"github.com/go-python/gpython/py"
	_ "github.com/go-python/gpython/stdlib"
	"lucksystem/charset"
	"lucksystem/game/engine"
	"lucksystem/game/runtime"
)

type Plugin struct {
	file   string
	ctx    py.Context
	module *py.Module
}

func NewPlugin(file string) *Plugin {
	p := &Plugin{
		file: file,
		ctx: py.NewContext(py.ContextOpts{
			SysPaths: []string{"."},
		}),
	}
	var err error
	p.module, err = py.RunFile(p.ctx, p.file, py.CompileOpts{
		CurDir: "/",
	}, nil)
	if err != nil {
		py.TracebackDump(err)
		//glog.V(3).Infof("SELECT #%d = %d\n", varID, selectID)
	}
	return p
}

func (g *Plugin) Init(ctx *runtime.Runtime) {
	ctx.Init(charset.ShiftJIS, charset.Unicode, true)
	pluginContext.ctx = ctx
	call, ok := g.module.Globals["Init"]
	if ok {
		_, err := py.Call(call, nil, nil)
		if err != nil {
			py.TracebackDump(err)
		}
	}
}

func (g *Plugin) UNDEFINED(ctx *runtime.Runtime, opcode string) engine.HandlerFunc {
	// plugin
	if strings.HasPrefix(opcode, "0x") {
		if opcodeMap, ok := g.module.Globals["opcode_dict"]; ok {
			if dict, ok := opcodeMap.(py.StringDict); ok {
				if val, ok := dict[opcode]; ok {
					opcode = string(val.(py.String))
				}
			}
		}
	}

	if call, ok := g.module.Globals[opcode]; ok {
		pluginContext.NewOP(ctx)
		_, err := py.Call(call, nil, nil)
		if err != nil {
			py.TracebackDump(err)
		}
	} else {
		code := ctx.Code()
		op := NewOP(ctx, code.ParamBytes, 0)
		op.Read(ctx.DefaultExport)
		op.SetOperateParams()
	}
	return func() {
		// 下一步执行地址，为0则表示紧接着向下
		ctx.ChanEIP <- 0
	}
}
