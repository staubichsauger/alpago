package main

import (
	"flag"
	"github.com/staubichsauger/alpago/internal/monolith/bot"
	"log"
	"net/http"
	"net/url"
	"strconv"
	"strings"
)

func main() {
	host := flag.String("host", "localhost", "alpaca server host")
	port := flag.Int("port", 3000, "alpaca server port")
	numPlayers := flag.Int("players", 4, "number of clients to spawn")
	ipAddr := flag.String("ip", "localhost", "client ip address")
	dropLimit := flag.Int("limit", 1, "when reaching x points, drop out of the round")
	flag.Parse()

	url, err := url.Parse("http://" + *host + ":" + strconv.Itoa(*port))
	if err != nil {
		log.Fatal("Invalid hostname and port supplied: ", err)
	}

	var players []*bot.Client

	for i:=0; i < *numPlayers; i++ {
		players = append(players, &bot.Client{
			Url: *url,
			DropLimit: *dropLimit+1,
			IpAddr: *ipAddr,
		})
	}

	stop := make(chan error)

	for i, p := range players {
		if err := p.Login("alpago", i); err != nil {
			log.Fatal(err)
		}
		//go p.Play(stop)
	}

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		//fmt.Println("Received request: ", r.URL.Path)
		for k, p := range players {
			if strings.HasSuffix(r.URL.Path, strconv.Itoa(k)) {
				p.PlayTurn(stop)
				w.WriteHeader(200)
				_, _ = w.Write([]byte("OK"))
				return
			}
		}
	})

	go func() {
		if err := http.ListenAndServe(":8181", nil); err != nil {
			//stop <- err
			return
		}
	}()

	for err := range stop {
		log.Print(err)
		if strings.Contains(err.Error(), "Status: 500") {
			return
		}
	}
}
