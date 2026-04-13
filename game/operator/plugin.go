package operator

import (
	"fmt"
	"path/filepath"
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

// NewPlugin loads a Python plugin file (e.g. data/KANON.py) into a gpython
// context. The plugin's directory is added to sys.path so that plugins may
// import shared helpers via relative package imports (e.g.
// `from base.kanon import *`). If the plugin fails to load, an error is
// reported and a non-nil *Plugin is still returned; callers must be
// resilient to a nil module (see Init / UNDEFINED).
func NewPlugin(file string) *Plugin {
	absFile, err := filepath.Abs(file)
	if err != nil {
		absFile = file
	}
	pluginDir := filepath.Dir(absFile)

	p := &Plugin{
		file: absFile,
		ctx: py.NewContext(py.ContextOpts{
			// Include the plugin's own directory first so that imports like
			// `from base.xxx import *` resolve against the plugin tree
			// rather than the process working directory. "." is kept as a
			// fallback for backward compatibility with the previous behaviour.
			SysPaths: []string{pluginDir, "."},
		}),
	}
	p.module, err = py.RunFile(p.ctx, p.file, py.CompileOpts{
		CurDir: pluginDir,
	}, nil)
	if err != nil {
		py.TracebackDump(err)
		fmt.Printf("[ERROR] Failed to load plugin %q: %v\n", p.file, err)
	}
	return p
}

func (g *Plugin) Init(ctx *runtime.Runtime) {
	ctx.Init(charset.ShiftJIS, charset.Unicode, true)
	pluginContext.ctx = ctx
	if g.module == nil {
		// Plugin failed to load; the traceback was already printed by
		// NewPlugin. Skip Init rather than panicking on a nil module.
		return
	}
	call, ok := g.module.Globals["Init"]
	if ok {
		_, err := py.Call(call, nil, nil)
		if err != nil {
			py.TracebackDump(err)
		}
	}
}

func (g *Plugin) UNDEFINED(ctx *runtime.Runtime, opcode string) engine.HandlerFunc {
	// If the plugin failed to load, fall through to the default "advance PC"
	// behaviour so the VM does not crash on a nil module.
	if g.module == nil {
		return func() {
			ctx.ChanEIP <- 0
		}
	}

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
