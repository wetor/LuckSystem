package paramter

type Paramter interface {
	Type() string
	Value() interface{}
	String() string
}
