package paramter

import "fmt"

// ===========================
type LInt16 struct {
	Data int16
}

func (LInt16) Type() string {
	return "LUCA.Int16"
}

func (p *LInt16) Value() interface{} {
	return p.Data
}

func (p *LInt16) String() string {
	return fmt.Sprintf("%d", p.Data)
}

// ===========================
type LInt32 struct {
}

// ===========================
type LUint16 struct {
	Data uint16
}

func (LUint16) Type() string {
	return "LUCA.Int16"
}

func (p *LUint16) Value() interface{} {
	return p.Data
}

func (p *LUint16) String() string {
	return fmt.Sprintf("%d", p.Data)
}

// ===========================
type LUint32 struct {
	Data uint32
}

func (LUint32) Type() string {
	return "LUCA.Int32"
}

func (p *LUint32) Value() interface{} {
	return p.Data
}

func (p *LUint32) String() string {
	return fmt.Sprintf("%d", p.Data)
}
