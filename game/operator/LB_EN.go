package operator

import (
	"lucksystem/charset"
	"lucksystem/game/engine"
	"lucksystem/game/runtime"
	"lucksystem/script"

	"github.com/golang/glog"
)

type LB_EN struct {
	LucaOperateUndefined
	LucaOperateDefault
	LucaOperateExpr
}

func NewLB_EN() *LB_EN {
	return &LB_EN{}
}

func (g *LB_EN) Init(ctx *runtime.Runtime) {
	ctx.Init(charset.ShiftJIS, charset.Unicode, true)
}

func (g *LB_EN) MESSAGE(ctx *runtime.Runtime) engine.HandlerFunc {
	code := ctx.Code()
	var voiceId uint16
	var msgStr_jp string
	var msgStr_en string
	var end uint8

	next := GetParam(code.ParamBytes, &voiceId)
	next = GetParam(code.ParamBytes, &msgStr_jp, next, 0, ctx.TextCharset)
	next = GetParam(code.ParamBytes, &msgStr_en, next, 0, ctx.TextCharset)
	GetParam(code.ParamBytes, &end, next)
	ctx.Script.SetOperateParams(ctx.CIndex, ctx.RunMode,
		voiceId,
		&script.StringParam{
			Data:   msgStr_jp,
			Coding: ctx.TextCharset,
		}, &script.StringParam{
			Data:   msgStr_en,
			Coding: ctx.TextCharset,
		},
		end,
		[]bool{true, true, true, false},
	)
	return func() {
		// 这里是执行内容
		ctx.Engine.MESSAGE(voiceId, msgStr_jp)
		ctx.ChanEIP <- 0
	}
}

func (g *LB_EN) SELECT(ctx *runtime.Runtime) engine.HandlerFunc {
	code := ctx.Code()
	var varID uint16
	var var0 uint16
	var var1 uint16
	var var2 uint16
	var msgStr_jp string
	var msgStr_en string
	var var3 uint16
	var var4 uint16
	var var5 uint16

	next := GetParam(code.ParamBytes, &varID)
	next = GetParam(code.ParamBytes, &var0, next)
	next = GetParam(code.ParamBytes, &var1, next)
	next = GetParam(code.ParamBytes, &var2, next)
	next = GetParam(code.ParamBytes, &msgStr_jp, next, 0, ctx.TextCharset)
	next = GetParam(code.ParamBytes, &msgStr_en, next, 0, ctx.TextCharset)
	next = GetParam(code.ParamBytes, &var3, next)
	next = GetParam(code.ParamBytes, &var4, next)
	GetParam(code.ParamBytes, &var5, next)
	ctx.Script.SetOperateParams(ctx.CIndex, ctx.RunMode,
		varID,
		var0,
		var1,
		var2,
		&script.StringParam{
			Data:   msgStr_jp,
			Coding: ctx.TextCharset,
		}, &script.StringParam{
			Data:   msgStr_en,
			Coding: ctx.TextCharset,
		},
		var3,
		var4,
		var5,
		[]bool{true, false, false, false, true, true, false, false, false},
	)
	return func() {

		selectID := ctx.Engine.SELECT(msgStr_jp)

		//fun.Call(&varID, msgStr_en)
		glog.V(3).Infof("SELECT #%d = %d\n", varID, selectID)
		ctx.ChanEIP <- 0
	}
}

func (g *LB_EN) IMAGELOAD(ctx *runtime.Runtime) engine.HandlerFunc {
	code := ctx.Code()
	/*	var mode uint16
		var imgID uint16
		var var1 uint16
		var xPos uint16
		var yPos uint16

		next := GetParam(code.ParamBytes, &mode)
		next = GetParam(code.ParamBytes, &imgID, next)
		if mode == 0 {
			// 背景
			ctx.Script.SetOperateParams(ctx.CIndex, ctx.RunMode,
				mode,
				imgID,
			)
		} else {
			// 立绘
			next = GetParam(code.ParamBytes, &var1, next)
			next = GetParam(code.ParamBytes, &xPos, next)
			GetParam(code.ParamBytes, &yPos, next)
			ctx.Script.SetOperateParams(ctx.CIndex, ctx.RunMode,
				mode,
				imgID,
				var1,
				xPos,
				yPos,
			)
		}*/
	list, end := AllToUint16(code.ParamBytes)
	if end >= 0 {
		ctx.Script.SetOperateParams(ctx.CIndex, ctx.RunMode,
			list,
			code.ParamBytes[end],
		)
	} else {
		ctx.Script.SetOperateParams(ctx.CIndex, ctx.RunMode,
			list,
		)
	}

	return func() {
		// 这里是执行 与虚拟机逻辑有关的代码

		// 下一步执行地址，为0则表示紧接着向下
		ctx.ChanEIP <- 0
	}
}

func (g *LB_EN) BATTLE(ctx *runtime.Runtime) engine.HandlerFunc {
	code := ctx.Code()
	var battleId uint16
	var msgStr_jp string
	var msgStr_en string
	var exprStr string
	var var1 uint16
	var var2 uint16
	var var3 uint16

	next := GetParam(code.ParamBytes, &battleId)
	if len(code.ParamBytes) <= next {
		list, end := AllToUint16(code.ParamBytes)
		if end >= 0 {
			ctx.Script.SetOperateParams(ctx.CIndex, ctx.RunMode,
				list,
				code.ParamBytes[end],
			)
		} else {
			ctx.Script.SetOperateParams(ctx.CIndex, ctx.RunMode,
				list,
			)
		}
	} else if battleId == 301 || battleId == 302 {
		GetParam(code.ParamBytes, &var1, next)
		ctx.Script.SetOperateParams(ctx.CIndex, ctx.RunMode,
			battleId,
			var1,
			[]bool{true, true},
		)
	} else if battleId == 300 {
		GetParam(code.ParamBytes, &msgStr_jp, next, 0, ctx.ExprCharset)
		ctx.Script.SetOperateParams(ctx.CIndex, ctx.RunMode,
			battleId,
			&script.StringParam{
				Data:   msgStr_jp,
				Coding: ctx.ExprCharset,
			},
			[]bool{true, true},
		)
	} else if battleId == 101 {
		next = GetParam(code.ParamBytes, &var1, next)
		next = GetParam(code.ParamBytes, &var2, next)
		if var2 == 0 {
			next = GetParam(code.ParamBytes, &var3, next)
			next = GetParam(code.ParamBytes, &exprStr, next, 0, ctx.ExprCharset)
			next = GetParam(code.ParamBytes, &msgStr_jp, next, 0, ctx.TextCharset)
			GetParam(code.ParamBytes, &msgStr_en, next, 0, ctx.TextCharset)
			ctx.Script.SetOperateParams(ctx.CIndex, ctx.RunMode,
				battleId,
				var1,
				var2,
				var3,
				&script.StringParam{
					Data:   exprStr,
					Coding: ctx.ExprCharset,
				},
				&script.StringParam{
					Data:   msgStr_jp,
					Coding: ctx.TextCharset,
				},
				&script.StringParam{
					Data:   msgStr_en,
					Coding: ctx.TextCharset,
				},
				[]bool{true, true, true, true, true, true, true},
			)
		} else {
			next := GetParam(code.ParamBytes, &battleId)
			next = GetParam(code.ParamBytes, &var1, next)
			next = GetParam(code.ParamBytes, &msgStr_jp, next, 0, ctx.TextCharset)
			GetParam(code.ParamBytes, &msgStr_en, next, 0, ctx.TextCharset)
			ctx.Script.SetOperateParams(ctx.CIndex, ctx.RunMode,
				battleId,
				var1,
				&script.StringParam{
					Data:   msgStr_jp,
					Coding: ctx.TextCharset,
				},
				&script.StringParam{
					Data:   msgStr_en,
					Coding: ctx.TextCharset,
				},
				[]bool{true, true, true, true},
			)
		}
	} else if battleId == 102 {
		next = GetParam(code.ParamBytes, &var1, next)
		if var1 == 0 {
			next = GetParam(code.ParamBytes, &var2, next)
			next = GetParam(code.ParamBytes, &exprStr, next, 0, ctx.ExprCharset)
			next = GetParam(code.ParamBytes, &msgStr_jp, next, 0, ctx.TextCharset)
			GetParam(code.ParamBytes, &msgStr_en, next, 0, ctx.TextCharset)
			ctx.Script.SetOperateParams(ctx.CIndex, ctx.RunMode,
				battleId,
				var1,
				var2,
				&script.StringParam{
					Data:   exprStr,
					Coding: ctx.ExprCharset,
				},
				&script.StringParam{
					Data:   msgStr_jp,
					Coding: ctx.TextCharset,
				},
				&script.StringParam{
					Data:   msgStr_en,
					Coding: ctx.TextCharset,
				},
				[]bool{true, true, true, true, true, true},
			)
		} else {
			next := GetParam(code.ParamBytes, &battleId)
			next = GetParam(code.ParamBytes, &msgStr_jp, next, 0, ctx.TextCharset)
			GetParam(code.ParamBytes, &msgStr_en, next, 0, ctx.TextCharset)
			ctx.Script.SetOperateParams(ctx.CIndex, ctx.RunMode,
				battleId,
				&script.StringParam{
					Data:   msgStr_jp,
					Coding: ctx.TextCharset,
				},
				&script.StringParam{
					Data:   msgStr_en,
					Coding: ctx.TextCharset,
				},
				[]bool{true, true, true},
			)
		}
	} else {
		list, end := AllToUint16(code.ParamBytes)
		if end >= 0 {
			ctx.Script.SetOperateParams(ctx.CIndex, ctx.RunMode,
				list,
				code.ParamBytes[end],
			)
		} else {
			ctx.Script.SetOperateParams(ctx.CIndex, ctx.RunMode,
				list,
			)
		}
	}

	return func() {
		// 这里是执行内容
		ctx.Engine.MESSAGE(battleId, msgStr_jp)
		ctx.ChanEIP <- 0
	}
}
func (g *LB_EN) TASK(ctx *runtime.Runtime) engine.HandlerFunc {
	code := ctx.Code()
	var TaskID uint16
	var TaskVar1 uint16
	var TaskVar2 uint16
	var TaskVar3 uint16
	var TaskVar4 uint16
	var msgStr_jp1 string
	var msgStr_en1 string
	var msgStr_jp2 string
	var msgStr_en2 string
	var exprStr1 string
	var exprStr2 string
	var exprStr3 string

	next := GetParam(code.ParamBytes, &TaskID)
	if len(code.ParamBytes) <= next {
		list, end := AllToUint16(code.ParamBytes)
		if end >= 0 {
			ctx.Script.SetOperateParams(ctx.CIndex, ctx.RunMode,
				list,
				code.ParamBytes[end],
			)
		} else {
			ctx.Script.SetOperateParams(ctx.CIndex, ctx.RunMode,
				list,
			)
		}
	} else if TaskID == 4 {
		next = GetParam(code.ParamBytes, &TaskVar1, next)
		if len(code.ParamBytes) <= next {
			list, end := AllToUint16(code.ParamBytes)
			if end >= 0 {
				ctx.Script.SetOperateParams(ctx.CIndex, ctx.RunMode,
					list,
					code.ParamBytes[end],
				)
			} else {
				ctx.Script.SetOperateParams(ctx.CIndex, ctx.RunMode,
					list,
				)
			}
		} else if TaskVar1 == 0 || TaskVar1 == 4 || TaskVar1 == 5 {
			next = GetParam(code.ParamBytes, &TaskVar2, next)
			next = GetParam(code.ParamBytes, &msgStr_jp1, next, 0, ctx.TextCharset)
			GetParam(code.ParamBytes, &msgStr_en1, next, 0, ctx.TextCharset)
			ctx.Script.SetOperateParams(ctx.CIndex, ctx.RunMode,
				TaskID,
				TaskVar1,
				TaskVar2,
				&script.StringParam{
					Data:   msgStr_jp1,
					Coding: ctx.TextCharset,
				}, &script.StringParam{
					Data:   msgStr_en1,
					Coding: ctx.TextCharset,
				},
				[]bool{true, true, true, true, true},
			)
		} else if TaskVar1 == 1 {
			next = GetParam(code.ParamBytes, &TaskVar2, next)
			next = GetParam(code.ParamBytes, &TaskVar3, next)
			next = GetParam(code.ParamBytes, &TaskVar4, next)
			next = GetParam(code.ParamBytes, &msgStr_jp1, next, 0, ctx.TextCharset)
			next = GetParam(code.ParamBytes, &msgStr_en1, next, 0, ctx.TextCharset)
			next = GetParam(code.ParamBytes, &msgStr_jp2, next, 0, ctx.TextCharset)
			GetParam(code.ParamBytes, &msgStr_en2, next, 0, ctx.TextCharset)
			ctx.Script.SetOperateParams(ctx.CIndex, ctx.RunMode,
				TaskID,
				TaskVar1,
				TaskVar2,
				TaskVar3,
				TaskVar4,
				&script.StringParam{
					Data:   msgStr_jp1,
					Coding: ctx.TextCharset,
				}, &script.StringParam{
					Data:   msgStr_en1,
					Coding: ctx.TextCharset,
				}, &script.StringParam{
					Data:   msgStr_jp2,
					Coding: ctx.TextCharset,
				}, &script.StringParam{
					Data:   msgStr_en2,
					Coding: ctx.TextCharset,
				},
				[]bool{true, true, true, true, true, true, true, true, true},
			)
		} else if TaskVar1 == 6 {
			next = GetParam(code.ParamBytes, &TaskVar2, next)
			next = GetParam(code.ParamBytes, &TaskVar3, next)
			next = GetParam(code.ParamBytes, &msgStr_jp1, next, 0, ctx.TextCharset)
			GetParam(code.ParamBytes, &msgStr_en1, next, 0, ctx.TextCharset)
			ctx.Script.SetOperateParams(ctx.CIndex, ctx.RunMode,
				TaskID,
				TaskVar1,
				TaskVar2,
				TaskVar3,
				&script.StringParam{
					Data:   msgStr_jp1,
					Coding: ctx.TextCharset,
				}, &script.StringParam{
					Data:   msgStr_en1,
					Coding: ctx.TextCharset,
				},
				[]bool{true, true, true, true, true, true},
			)
		} else {
			list, end := AllToUint16(code.ParamBytes)
			if end >= 0 {
				ctx.Script.SetOperateParams(ctx.CIndex, ctx.RunMode,
					list,
					code.ParamBytes[end],
				)
			} else {
				ctx.Script.SetOperateParams(ctx.CIndex, ctx.RunMode,
					list,
				)
			}
		}
	} else if TaskID == 54 {
		GetParam(code.ParamBytes, &msgStr_en1, next, 0, ctx.TextCharset)
		ctx.Script.SetOperateParams(ctx.CIndex, ctx.RunMode,
			TaskID,
			&script.StringParam{
				Data:   msgStr_en1,
				Coding: ctx.TextCharset,
			},
			[]bool{true, true},
		)
	} else if TaskID == 23 {
		next = GetParam(code.ParamBytes, &TaskVar1, next)
		if len(code.ParamBytes) <= next {
			ctx.Script.SetOperateParams(ctx.CIndex, ctx.RunMode, TaskID, TaskVar1, []bool{true, true})
		} else if TaskVar1 == 12835 || TaskVar1 == 13859 {
			next = GetParam(code.ParamBytes, &exprStr1, next, 0, ctx.ExprCharset)
			next = GetParam(code.ParamBytes, &exprStr2, next, 0, ctx.ExprCharset)
			GetParam(code.ParamBytes, &exprStr3, next, 0, ctx.ExprCharset)
			ctx.Script.SetOperateParams(ctx.CIndex, ctx.RunMode,
				TaskID,
				TaskVar1,
				&script.StringParam{
					Data:   exprStr1,
					Coding: ctx.ExprCharset,
				},
				&script.StringParam{
					Data:   exprStr2,
					Coding: ctx.ExprCharset,
				},
				&script.StringParam{
					Data:   exprStr3,
					Coding: ctx.ExprCharset,
				},
				[]bool{true, true, true, true, true},
			)
		} else if TaskVar1 == 12589 {
			next = GetParam(code.ParamBytes, &TaskVar2, next)
			next = GetParam(code.ParamBytes, &exprStr1, next, 0, ctx.ExprCharset)
			GetParam(code.ParamBytes, &exprStr2, next, 0, ctx.ExprCharset)
			ctx.Script.SetOperateParams(ctx.CIndex, ctx.RunMode,
				TaskID,
				TaskVar1,
				TaskVar2,
				&script.StringParam{
					Data:   exprStr1,
					Coding: ctx.ExprCharset,
				},
				&script.StringParam{
					Data:   exprStr2,
					Coding: ctx.ExprCharset,
				},
				[]bool{true, true, true, true, true},
			)
		} else {
			list, end := AllToUint16(code.ParamBytes)
			if end >= 0 {
				ctx.Script.SetOperateParams(ctx.CIndex, ctx.RunMode,
					list,
					code.ParamBytes[end],
				)
			} else {
				ctx.Script.SetOperateParams(ctx.CIndex, ctx.RunMode,
					list,
				)
			}
		}
	} else if TaskID == 69 {
		next = GetParam(code.ParamBytes, &TaskVar1, next)
		next = GetParam(code.ParamBytes, &msgStr_jp1, next, 0, ctx.TextCharset)
		next = GetParam(code.ParamBytes, &msgStr_en1, next, 0, ctx.TextCharset)
		next = GetParam(code.ParamBytes, &msgStr_jp2, next, 0, ctx.TextCharset)
		GetParam(code.ParamBytes, &msgStr_en2, next, 0, ctx.TextCharset)
		ctx.Script.SetOperateParams(ctx.CIndex, ctx.RunMode,
			TaskID,
			TaskVar1,
			&script.StringParam{
				Data:   msgStr_jp1,
				Coding: ctx.TextCharset,
			}, &script.StringParam{
				Data:   msgStr_en1,
				Coding: ctx.TextCharset,
			}, &script.StringParam{
				Data:   msgStr_jp2,
				Coding: ctx.TextCharset,
			}, &script.StringParam{
				Data:   msgStr_en2,
				Coding: ctx.TextCharset,
			},
			[]bool{true, true, true, true, true, true},
		)
	} else if TaskID == 28 {
		next = GetParam(code.ParamBytes, &TaskVar1, next)
		if len(code.ParamBytes) <= next {
			list, end := AllToUint16(code.ParamBytes)
			if end >= 0 {
				ctx.Script.SetOperateParams(ctx.CIndex, ctx.RunMode,
					list,
					code.ParamBytes[end],
				)
			} else {
				ctx.Script.SetOperateParams(ctx.CIndex, ctx.RunMode,
					list,
				)
			}
		} else if TaskVar1 == 200 || TaskVar1 == 201 || TaskVar1 == 202 || TaskVar1 == 203 || TaskVar1 == 204 || TaskVar1 == 210 || TaskVar1 == 400 {
			next = GetParam(code.ParamBytes, &exprStr1, next, 0, ctx.ExprCharset)
			next = GetParam(code.ParamBytes, &exprStr2, next, 0, ctx.ExprCharset)
			GetParam(code.ParamBytes, &exprStr3, next, 0, ctx.ExprCharset)
			ctx.Script.SetOperateParams(ctx.CIndex, ctx.RunMode,
				TaskID,
				TaskVar1,
				&script.StringParam{
					Data:   exprStr1,
					Coding: ctx.ExprCharset,
				},
				&script.StringParam{
					Data:   exprStr2,
					Coding: ctx.ExprCharset,
				},
				&script.StringParam{
					Data:   exprStr3,
					Coding: ctx.ExprCharset,
				},
				[]bool{true, true, true, true, true},
			)
		} else {
			list, end := AllToUint16(code.ParamBytes)
			if end >= 0 {
				ctx.Script.SetOperateParams(ctx.CIndex, ctx.RunMode,
					list,
					code.ParamBytes[end],
				)
			} else {
				ctx.Script.SetOperateParams(ctx.CIndex, ctx.RunMode,
					list,
				)
			}
		}
	} else {
		list, end := AllToUint16(code.ParamBytes)
		if end >= 0 {
			ctx.Script.SetOperateParams(ctx.CIndex, ctx.RunMode,
				list,
				code.ParamBytes[end],
			)
		} else {
			ctx.Script.SetOperateParams(ctx.CIndex, ctx.RunMode,
				list,
			)
		}
	}
	return func() {
		ctx.Engine.MESSAGE(TaskID, msgStr_jp1)
		ctx.ChanEIP <- 0
	}
}
func (g *LB_EN) SAYAVOICETEXT(ctx *runtime.Runtime) engine.HandlerFunc {
	code := ctx.Code()
	var voiceId uint16
	var msgStr_jp string
	var msgStr_en string

	next := GetParam(code.ParamBytes, &voiceId)
	next = GetParam(code.ParamBytes, &msgStr_jp, next, 0, ctx.TextCharset)
	GetParam(code.ParamBytes, &msgStr_en, next, 0, ctx.TextCharset)
	ctx.Script.SetOperateParams(ctx.CIndex, ctx.RunMode,
		voiceId,
		&script.StringParam{
			Data:   msgStr_jp,
			Coding: ctx.TextCharset,
		}, &script.StringParam{
			Data:   msgStr_en,
			Coding: ctx.TextCharset,
		},
		[]bool{true, true, true},
	)
	return func() {
		ctx.Engine.MESSAGE(voiceId, msgStr_jp)
		ctx.ChanEIP <- 0
	}
}
func (g *LB_EN) EQU(ctx *runtime.Runtime) engine.HandlerFunc {
	code := ctx.Code()
	var value uint16
	var exprStr string

	next := GetParam(code.ParamBytes, &value)
	GetParam(code.ParamBytes, &exprStr, next, 0, ctx.ExprCharset)
	ctx.Script.SetOperateParams(ctx.CIndex, ctx.RunMode,
		value,
		exprStr,
		ctx.ExprCharset,
	)
	return func() {
		ctx.ChanEIP <- 0
	}
}
func (g *LB_EN) EQUN(ctx *runtime.Runtime) engine.HandlerFunc {
	code := ctx.Code()
	var key uint16
	var value uint16

	next := GetParam(code.ParamBytes, &key)
	if next < len(code.ParamBytes) {
		GetParam(code.ParamBytes, &value, next)
		ctx.Script.SetOperateParams(ctx.CIndex, ctx.RunMode,
			key,
			value,
		)
	} else {
		ctx.Script.SetOperateParams(ctx.CIndex, ctx.RunMode,
			key,
		)
	}

	return func() {
		//var keyStr string
		//if key <= 1 {
		//	keyStr = ToString("%d", key)
		//} else {
		//	keyStr = ToString("#%d", key)
		//}
		ctx.ChanEIP <- 0
	}
}
func (g *LB_EN) EQUV(ctx *runtime.Runtime) engine.HandlerFunc {
	code := ctx.Code()
	var key uint16
	var value uint16

	next := GetParam(code.ParamBytes, &key)
	if next < len(code.ParamBytes) {
		GetParam(code.ParamBytes, &value, next)
		ctx.Script.SetOperateParams(ctx.CIndex, ctx.RunMode,
			key,
			value,
		)
	} else {
		ctx.Script.SetOperateParams(ctx.CIndex, ctx.RunMode,
			key,
		)
	}
	return func() {
		//var keyStr string
		//if key <= 1 {
		//	keyStr = ToString("%d", key)
		//} else {
		//	keyStr = ToString("#%d", key)
		//}
		//ctx.Variable.Set(keyStr, int(value))
		ctx.ChanEIP <- 0
	}
}
func (g *LB_EN) VARSTR_SET(ctx *runtime.Runtime) engine.HandlerFunc {
	code := ctx.Code()
	var varstrId uint16
	var varstrStr string

	next := GetParam(code.ParamBytes, &varstrId)
	GetParam(code.ParamBytes, &varstrStr, next, 0, ctx.TextCharset)
	ctx.Script.SetOperateParams(ctx.CIndex, ctx.RunMode,
		varstrId,
		&script.StringParam{
			Data:   varstrStr,
			Coding: ctx.TextCharset,
		},
		[]bool{true, true},
	)
	return func() {
		ctx.Engine.MESSAGE(varstrId, varstrStr)
		ctx.ChanEIP <- 0
	}
}

// For some reason some scripts are not parsed if this opcode is not specified for LB_EN
func (g *LB_EN) MOVE(ctx *runtime.Runtime) engine.HandlerFunc {
	code := ctx.Code()
	var val1 uint8
	var val2 uint16
	var val3 uint16
	var height int16
	var width uint16

	next := GetParam(code.ParamBytes, &val1)
	next = GetParam(code.ParamBytes, &val2, next)
	next = GetParam(code.ParamBytes, &val3, next)
	next = GetParam(code.ParamBytes, &height, next)
	GetParam(code.ParamBytes, &width, next)
	ctx.Script.SetOperateParams(ctx.CIndex, ctx.RunMode,
		val1,
		val2,
		val3,
		height,
		width,
	)
	return func() {
		// 这里是执行 与虚拟机逻辑有关的代码

		// 下一步执行地址，为0则表示紧接着向下
		ctx.ChanEIP <- 0
	}
}
