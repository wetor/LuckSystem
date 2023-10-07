package operator

import (
	"github.com/go-python/gpython/py"
	"lucksystem/charset"
	"lucksystem/game/runtime"
)

var pluginContext *PluginContext

type PluginContext struct {
	ctx *runtime.Runtime
	op  *OP
}

func (p *PluginContext) NewOP(ctx *runtime.Runtime) {
	p.ctx = ctx
	p.op = NewOP(p.ctx, p.ctx.Code().ParamBytes, 0)
}

func (p *PluginContext) Read(self py.Object, args py.Tuple, kwargs py.StringDict) (py.Object, error) {
	var export py.Object = py.Bool(false)
	kwlist := []string{"export"}
	err := py.ParseTupleAndKeywords(args, kwargs, "|O:read", kwlist, &export)
	if err != nil {
		return nil, err
	}

	results := p.op.Read(bool(export.(py.Bool)))
	res := make(py.Tuple, len(results))
	for i, result := range results {
		res[i] = py.Int(result)
	}
	return res, nil
}

func (p *PluginContext) ReadUInt8(self py.Object, args py.Tuple, kwargs py.StringDict) (py.Object, error) {
	var export py.Object = py.Bool(false)
	kwlist := []string{"export"}
	err := py.ParseTupleAndKeywords(args, kwargs, "|O:read_uint8", kwlist, &export)
	if err != nil {
		return nil, err
	}
	result := p.op.ReadUInt8(bool(export.(py.Bool)))
	return py.Int(result), nil
}

func (p *PluginContext) ReadUInt16(self py.Object, args py.Tuple, kwargs py.StringDict) (py.Object, error) {
	var export py.Object = py.Bool(false)
	kwlist := []string{"export"}
	err := py.ParseTupleAndKeywords(args, kwargs, "|O:read_uint16", kwlist, &export)
	if err != nil {
		return nil, err
	}
	result := p.op.ReadUInt16(bool(export.(py.Bool)))
	return py.Int(result), nil
}

func (p *PluginContext) ReadUInt32(self py.Object, args py.Tuple, kwargs py.StringDict) (py.Object, error) {
	var export py.Object = py.Bool(false)
	kwlist := []string{"export"}
	err := py.ParseTupleAndKeywords(args, kwargs, "|O:read_uint32", kwlist, &export)
	if err != nil {
		return nil, err
	}
	result := p.op.ReadUInt32(bool(export.(py.Bool)))
	return py.Int(result), nil
}

func (p *PluginContext) ReadJump(self py.Object, args py.Tuple, kwargs py.StringDict) (py.Object, error) {
	var file py.Object = py.String("")
	var export py.Object = py.Bool(true)
	kwlist := []string{"file", "export"}
	err := py.ParseTupleAndKeywords(args, kwargs, "|OO:read_jump", kwlist, &file, &export)
	if err != nil {
		return nil, err
	}
	var result uint32
	if len(file.(py.String)) > 0 {
		result = p.op.ReadFileJump(bool(export.(py.Bool)), string(file.(py.String)))
	} else {
		result = p.op.ReadJump(bool(export.(py.Bool)))
	}
	return py.Int(result), nil
}

func (p *PluginContext) ReadString(self py.Object, args py.Tuple, kwargs py.StringDict) (py.Object, error) {
	var _charset py.Object = py.String(p.ctx.TextCharset)
	var export py.Object = py.Bool(true)
	kwlist := []string{"charset", "export"}
	err := py.ParseTupleAndKeywords(args, kwargs, "|OO:read_str", kwlist, &_charset, &export)
	if err != nil {
		return nil, err
	}
	result := p.op.ReadString(bool(export.(py.Bool)), charset.Charset(_charset.(py.String)))
	return py.String(result), nil
}

func (p *PluginContext) ReadLenString(self py.Object, args py.Tuple, kwargs py.StringDict) (py.Object, error) {
	var _charset py.Object = py.String(p.ctx.TextCharset)
	var export py.Object = py.Bool(true)
	kwlist := []string{"charset", "export"}
	err := py.ParseTupleAndKeywords(args, kwargs, "|OO:read_len_str", kwlist, &_charset, &export)
	if err != nil {
		return nil, err
	}
	result := p.op.ReadLString(bool(export.(py.Bool)), charset.Charset(_charset.(py.String)))
	return py.String(result), nil
}

func (p *PluginContext) End(self py.Object, args py.Tuple) (py.Object, error) {
	p.op.SetOperateParams()
	return nil, nil
}

func (p *PluginContext) CanRead(self py.Object, args py.Tuple) (py.Object, error) {
	return py.Bool(p.op.CanRead()), nil
}

func (p *PluginContext) SetConfig(self py.Object, args py.Tuple, kwargs py.StringDict) (py.Object, error) {
	var exprCharset py.Object
	var textCharset py.Object
	var defaultExport py.Object = py.Bool(true)
	kwlist := []string{"expr_charset", "text_charset", "default_export"}
	err := py.ParseTupleAndKeywords(args, kwargs, "OO|O:set_config", kwlist, &exprCharset, &textCharset, &defaultExport)
	if err != nil {
		return nil, err
	}
	p.ctx.Init(
		charset.Charset(exprCharset.(py.String)),
		charset.Charset(textCharset.(py.String)),
		bool(defaultExport.(py.Bool)))

	mod := self.(*py.Module)
	mod.Globals["expr"] = exprCharset
	mod.Globals["text"] = textCharset
	return nil, nil
}

func init() {
	pluginContext = &PluginContext{}

	methods := []*py.Method{
		py.MustNewMethod("read", pluginContext.Read, 0, `read(export=False) -> list(int)`),
		py.MustNewMethod("read_uint8", pluginContext.ReadUInt8, 0, `read_uint8(export=False) -> int`),
		py.MustNewMethod("read_uint16", pluginContext.ReadUInt16, 0, `read_uint16(export=False) -> int`),
		py.MustNewMethod("read_uint32", pluginContext.ReadUInt32, 0, `read_uint32(export=False) -> int`),
		py.MustNewMethod("read_jump", pluginContext.ReadJump, 0, `read_jump(file='', export=True) -> int`),
		py.MustNewMethod("read_str", pluginContext.ReadString, 0, `read_str(charset=textCharset, export=True) -> str`),
		py.MustNewMethod("read_len_str", pluginContext.ReadLenString, 0, `read_len_str(charset=textCharset, export=True) -> str`),
		py.MustNewMethod("end", pluginContext.End, 0, `end()`),
		py.MustNewMethod("can_read", pluginContext.CanRead, 0, `can_read() -> bool`),
		py.MustNewMethod("set_config", pluginContext.SetConfig, 0, `set_config(expr_charset, text_charset, default_export=True)`),
	}

	py.RegisterModule(&py.ModuleImpl{
		Info: py.ModuleInfo{
			Name: "core",
			Doc:  "Core Module",
		},
		Methods: methods,
		Globals: py.StringDict{
			"Charset_UTF8":    py.String(charset.UTF_8),
			"Charset_Unicode": py.String(charset.Unicode),
			"Charset_SJIS":    py.String(charset.ShiftJIS),
		},
	})
}
