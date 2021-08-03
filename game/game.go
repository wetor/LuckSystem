package game

type Game struct {
	GameName string
	Version  uint8
}

func NewGame() *Game {
	return &Game{}
}
