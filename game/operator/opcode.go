package operator

import (
	"strings"

	"lucksystem/charset"
	"lucksystem/game/runtime"
	"lucksystem/script"
)

type Param struct {
	Type   string
	Value  interface{}
	Export bool
}

type OP struct {
	ctx    *runtime.Runtime
	data   []byte
	offset int
	params []Param
	index  int
}

func NewOP(ctx *runtime.Runtime, data []byte, offset int) *OP {
	return &OP{
		ctx:    ctx,
		data:   data,
		offset: offset,
		params: make([]Param, 0, 8),
		index:  0,
	}
}

func (op *OP) CanRead() bool {
	return op.offset < len(op.data)
}

func (op *OP) Read(export bool) []uint16 {
	val, end := AllToUint16(op.data[op.offset:])

	for _, v := range val {
		op.params = append(op.params, Param{
			Type:   "uint16",
			Value:  v,
			Export: export,
		})
		op.index++
	}
	if end >= 0 {
		op.offset += end
		op.params = append(op.params, Param{
			Type:   "uint8",
			Value:  op.data[op.offset],
			Export: export,
		})
		op.index++

	} else {
		op.offset = len(op.data)
	}
	return val
}

func (op *OP) ReadUInt8(export bool) uint8 {
	var val uint8
	op.offset = GetParam(op.data, &val, op.offset)
	op.params = append(op.params, Param{
		Type:   "uint8",
		Value:  val,
		Export: export,
	})
	op.index++
	return val
}

func (op *OP) ReadUInt16(export bool) uint16 {
	var val uint16
	op.offset = GetParam(op.data, &val, op.offset)
	op.params = append(op.params, Param{
		Type:   "uint16",
		Value:  val,
		Export: export,
	})
	op.index++
	return val
}

func (op *OP) ReadUInt32(export bool) uint32 {
	var val uint32
	op.offset = GetParam(op.data, &val, op.offset)
	op.params = append(op.params, Param{
		Type:   "uint32",
		Value:  val,
		Export: export,
	})
	op.index++
	return val
}

func (op *OP) ReadString(export bool, charset charset.Charset) string {
	var val string
	op.offset = GetParam(op.data, &val, op.offset, 0, charset)
	op.params = append(op.params, Param{
		Type: "string",
		Value: &script.StringParam{
			Data:   val,
			Coding: charset,
			HasLen: false,
		},
		Export: export,
	})
	op.index++
	return val
}

func (op *OP) ReadLString(export bool, charset charset.Charset) string {
	var val lstring
	op.offset = GetParam(op.data, &val, op.offset, 0, charset)
	op.params = append(op.params, Param{
		Type: "lstring",
		Value: &script.StringParam{
			Data:   string(val),
			Coding: charset,
			HasLen: true,
		},
		Export: export,
	})
	op.index++
	return string(val)
}

func (op *OP) ReadJump(export bool) uint32 {
	var val uint32
	op.offset = GetParam(op.data, &val, op.offset)
	op.params = append(op.params, Param{
		Type: "jump",
		Value: &script.JumpParam{
			Position: int(val),
		},
		Export: export,
	})
	op.index++
	return val
}

func (op *OP) ReadFileJump(export bool, file string) uint32 {
	var val uint32
	op.offset = GetParam(op.data, &val, op.offset)
	param := &script.JumpParam{
		ScriptName: file,
		Position:   int(val),
	}
	if file != op.ctx.Script.Name &&
		strings.ToUpper(file) != op.ctx.Script.Name &&
		strings.ToLower(file) != op.ctx.Script.Name {
		param.GlobalIndex = op.ctx.AddLabel(file, op.ctx.CIndex, int(val))
	}
	op.params = append(op.params, Param{
		Type:   "filejump",
		Value:  param,
		Export: export,
	})
	op.index++
	return val
}

func (op *OP) SetOperateParams() {
	params := make([]interface{}, len(op.params), len(op.params)+1)
	requires := make([]bool, len(op.params))
	// types := make([]string, len(op.params))
	for i, param := range op.params {
		params[i] = param.Value
		requires[i] = param.Export
		// types[i] = param.Type
	}
	params = append(params, requires)
	_ = op.ctx.Script.SetOperateParams(op.ctx.CIndex, op.ctx.RunMode, params...)
}
