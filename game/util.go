package game

import (
	"os"
	"strings"
)

func ScriptCanLoad(name string) bool {
	for _, n := range ScriptBlackList {
		if strings.ToUpper(name) == n {
			return false
		}
	}
	return true
}

func IsExistDir(filepath string) (exist bool, dir bool) {
	s, err := os.Stat(filepath)
	exist = true
	if err != nil {
		exist = os.IsExist(err)
	}
	if exist {
		dir = s.IsDir()
	}
	return
}
