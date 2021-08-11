package game

func ScriptCanLoad(name string) bool {
	for _, n := range ScriptBlackList {
		if name == n {
			return false
		}
	}
	return true
}
