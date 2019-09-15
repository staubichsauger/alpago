package game

import (
	"fmt"
)

type Status struct {
	MyTurn        bool     `json:"my_turn"`
	Hand          []Card   `json:"hand"`
	//OtherPlayers  []Player `json:"other_players"`
	DiscardedCard Card     `json:"discarded_card"`
	Score		  int	   `json:"score"`
	CardpileCards int	   `json:"cardpile_cards"`
	PlayersLeft   int      `json:"players_left"`
	possible	  []Card
	numOfCards	  map[string]int
	dropLimit 	  int
	sum			  int
}

func (gs *Status) DoTurn(dropLimit int) (turn Turn, lastCard bool) {
	gs.analiseHand()
	gs.dropLimit = dropLimit

	action, card := gs.getBestCard()

	if len(gs.Hand) == 1 && action == "DROP CARD" {
		fmt.Println("No cards left after this turn!")
		lastCard = true
		//time.Sleep(time.Second * 10)
	}

	turn = Turn{
		Action: action,
		Card: card,
	}

	//time.Sleep(time.Millisecond * 500)
	return turn, lastCard
}

func (gs *Status) analiseHand() {
	gs.numOfCards = make(map[string]int)
	for _, c := range gs.Hand {

		if c.Value == gs.DiscardedCard.Value || c.Name == order[gs.DiscardedCard.Value+1] {
			gs.possible = append(gs.possible, c)
		}

		gs.numOfCards[c.Name] = gs.numOfCards[c.Name] + 1
	}
	//fmt.Println(gs.possible, gs.numOfCards)
}

func (gs *Status) getBestAction() string {
	if gs.PlayersLeft > 1 && gs.CardpileCards > gs.PlayersLeft*2 {
		return "DRAW CARD"
	}
	if gs.sum < 2*gs.dropLimit || gs.CardpileCards == 0 || gs.PlayersLeft == 1 {
		return "LEAVE ROUND"
	}
	return "DRAW CARD"
}

func (gs *Status) getBestCard() (string, string) {
	gs.sum = 0
	for v, _ := range gs.numOfCards {
		if stringToInt[v] > 0 {
			gs.sum += stringToInt[v]
		}
	}

	if len(gs.possible) == 0 {
		return gs.getBestAction(), ""
	}

	if gs.sum <= gs.dropLimit && gs.CardpileCards < 2*gs.PlayersLeft {
		return "LEAVE ROUND", ""
	}

	card := gs.possible[0]
	for _, pc := range gs.possible {
		if gs.numOfCards[pc.Name] < gs.numOfCards[card.Name] {
			card = pc
		} else if gs.numOfCards[pc.Name] == gs.numOfCards[card.Name] {
			if pc.Value > card.Value {
				card = pc
			}
		}
	}

	return "DROP CARD", card.Name
}