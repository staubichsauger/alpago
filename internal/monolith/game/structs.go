package game

type Join struct {
	Name string `json:"name"`
	CallbackUrl string `json:"callbackUrl"`
}

type Card struct {
	Type  int    `json:"type"`
	Name  string `json:"name"`
	Value int    `json:"value"`
}

type Player struct {
	Name        string      `json:"player_name"`
	CardCount   int    		`json:"hand_cards"`
	Score 		int			`json:"score"`
	LeftRound   bool		`json:"left_round"`
}

type Turn struct {
	Action	 string	`json:"action"`
	Card 	 string `json:"card"`
}

type Id struct {
	PlayerId string `json:"player_id"`
	PlayerName string `json:"player_name"`
}