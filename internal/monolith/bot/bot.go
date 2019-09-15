package bot

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/staubichsauger/alpago/internal/monolith/game"
	"io/ioutil"
	"net/http"
	"net/url"
	"strconv"
	"time"
)

type Client struct {
	Url url.URL
	Id string
	Stop chan error
	DropLimit int
	IpAddr string
	LastCard bool
	Score int
}

func (a *Client) Login(name string, num int) error {
	me := game.Join{
		Name: name + strconv.Itoa(a.DropLimit),
		CallbackUrl: "http://" + a.IpAddr + ":8181/alpaca/" + strconv.Itoa(num),
	}
	reqBytes, err := json.Marshal(&me)
	if err != nil {
		return err
	}

	res, err := http.Post(a.Url.String() + "/join", "application/json", bytes.NewBuffer(reqBytes))
	if err != nil {
		return err
	}
	if res.StatusCode != http.StatusOK {
		return errors.New("Got status: " + strconv.Itoa(res.StatusCode))
	}

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return err
	}

	id := game.Id{}
	err = json.Unmarshal(body, &id)
	if err != nil {
		return err
	}

	a.Id = id.PlayerId
	return nil
}

func (a *Client) Play(stop chan error) {
	a.Stop = stop
	for _ = range time.Tick(time.Millisecond * 20) {
		if !a.PlayTurn(stop) {
			return
		}
	}
}

func (a *Client) PlayTurn(stop chan error) bool {
	res, err := http.Get(a.Url.String() + "/alpaca?id=" + a.Id)
	if err != nil {
		stop <- err
		return false
	}
	if res.StatusCode != http.StatusOK {
		stop <- errors.New("Error getting games endpoint, Status: " + strconv.Itoa(res.StatusCode))
		return true
	}

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		stop <- err
		return false
	}

	gs := game.Status{}
	err = json.Unmarshal(body, &gs)
	if err != nil {
		stop <- err
		return false
	}

	if gs.MyTurn {
		var turn game.Turn
		turn, a.LastCard = gs.DoTurn(a.DropLimit)
		if a.LastCard {
			a.Score = gs.Score
		}
		bj, err := json.Marshal(&turn)
		if err != nil {
			stop <- err
			return false
		}
		res, err := http.Post(a.Url.String() + "/alpaca?id=" + a.Id, "application/json", bytes.NewBuffer(bj))
		if err != nil {
			stop <- err
			return false
		}
		if res.StatusCode != http.StatusOK {
			//log.Print(*turn.PlayCard)
			body, err := ioutil.ReadAll(res.Body)
			if err != nil {
				stop <- err
				return false
			}
			stop <- errors.New("Error posting turn: " + strconv.Itoa(res.StatusCode) + "-> " + string(body))
		}

		// Check for correct score keeping
		if a.LastCard {
			time.Sleep(time.Millisecond * 10)

			res, err = http.Get(a.Url.String() + "/alpaca?id=" + a.Id)
			if err != nil {
				stop <- err
				return false
			}
			if res.StatusCode != http.StatusOK {
				stop <- errors.New("Error getting games endpoint, Status: " + strconv.Itoa(res.StatusCode))
				return true
			}

			body, err := ioutil.ReadAll(res.Body)
			if err != nil {
				stop <- err
				return false
			}

			gs := game.Status{}
			err = json.Unmarshal(body, &gs)
			if err != nil {
				stop <- err
				return false
			}

			if a.Score >= 10 && a.Score - 10 == gs.Score {
				fmt.Println("lost 10 points!")
			} else if a.Score >= 1 && a.Score - 1 == gs.Score {
				fmt.Println("lost 1 point!")
			} else if gs.Score != 0 {
				fmt.Println("I've been treated unfairly!")
			}
			a.LastCard = false
		}

	}
	return true
}