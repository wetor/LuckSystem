package paramter

import (
	"fmt"
	"lucascript/charset"
)

type LString struct {
	Data    string
	Charset charset.Charset
}

func (p *LString) Type() string {
	return fmt.Sprintf("LUCA.String(%s)", p.Charset)
}

func (p *LString) Value() interface{} {
	return p.Data
}

func (p *LString) String() string {
	return p.Data
}
