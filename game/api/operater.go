package api

import (
	"lucksystem/game/engine"
	"lucksystem/game/runtime"
)

// Operator 需定制指令
type Operator interface {
	Init(ctx *runtime.Runtime)
}

// LucaOperator 通用指令
type LucaOperator interface {
	LucaDefaultOperator
	LucaExprOperator
	LucaUndefinedOperator
}

type LucaDefaultOperator interface {
	UNKNOW0(ctx *runtime.Runtime) engine.HandlerFunc
	IFN(ctx *runtime.Runtime) engine.HandlerFunc
	IFY(ctx *runtime.Runtime) engine.HandlerFunc
	GOTO(ctx *runtime.Runtime) engine.HandlerFunc
	JUMP(ctx *runtime.Runtime) engine.HandlerFunc
	FARCALL(ctx *runtime.Runtime) engine.HandlerFunc
	MOVE(ctx *runtime.Runtime) engine.HandlerFunc
}

type LucaExprOperator interface {
	EQU(ctx *runtime.Runtime) engine.HandlerFunc
	EQUN(ctx *runtime.Runtime) engine.HandlerFunc
	ADD(ctx *runtime.Runtime) engine.HandlerFunc
	RANDOM(ctx *runtime.Runtime) engine.HandlerFunc
}

type LucaUndefinedOperator interface {
	UNDEFINED(ctx *runtime.Runtime, opname string) engine.HandlerFunc
}
