package expr

import (
	"container/list"
	"fmt"
	"testing"
)

func TestRun(t *testing.T) {

	res, err := RunExpr("#5001+1==10", map[string]int{"#5001": 10})
	if err != nil {
		fmt.Print(err)
	}
	fmt.Println(res)

}
func TestParser(t *testing.T) {

	tokens, _ := Parser("a>>b")
	for i, token := range tokens {
		fmt.Println(i, token.Data, token.Type)
	}

}
func TestStack(t *testing.T) {

	l := list.New()
	l.PushBack("a")
	l.PushBack("b")
	l.PushBack("c")

	tmp := l.Back()
	l.Remove(tmp)
	fmt.Println(tmp.Value.(string))
	tmp = l.Back()
	l.Remove(tmp)
	fmt.Println(tmp.Value.(string))
	tmp = l.Back()
	l.Remove(tmp)
	fmt.Println(tmp.Value.(string))

}
