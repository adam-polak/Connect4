package player

type Player struct {
	uid string
}

func NewPlayer(username string, password string) *Player {
	return &Player{
		uid: "secret",
	}
}

func GetPlayer(loginKey string) *Player {
	return &Player{
		uid: "secret",
	}
}
