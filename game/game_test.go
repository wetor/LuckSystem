package game

import (
	"fmt"
	"lucascript/charset"
	"lucascript/operation"
	"testing"
)

func TestGame(t *testing.T) {
	by := []byte{0x60, 00, 0x52, 00, 0x69, 00, 0x6B, 00, 0x69, 0x00, 0x40, 00, 00, 0x00, 0x60, 00, 0x52, 00, 0x69, 00, 0x6B, 00, 0x69, 0x00, 0x40, 00}

	str, n := operation.ReadString(by, 0, charset.Unicode)
	fmt.Println(str, n)
	str, n = operation.ReadString(by, n, charset.Unicode)
	fmt.Println(str, n)

}
