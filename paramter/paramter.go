package paramter

type Paramter interface {
	Type() string
	Value() interface{}
	String() string
}
type LBytes struct {
	Data []byte
}

func (LBytes) Type() string {
	return "LUCA.Bytes"
}

func (p *LBytes) Value() interface{} {
	return p.Data
}

func (p *LBytes) String() string {
	return string(p.Data)
}
