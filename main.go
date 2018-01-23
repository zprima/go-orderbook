package main

import (
	"encoding/json"
	"errors"
	"log"
	"reflect"
	"strings"

	"github.com/toorop/go-pusher"
)

const (
	APP_KEY = "de504dc5763aeef9ff52"
)

type PusherOrderbook struct {
	Asks [][]string `json:"asks"`
	Bids [][]string `json:"bids"`
}

type NormalOrderbook struct {
	Asks []Ask `json:"asks"`
	Bids []Bid `json:"bids"`
}

type Ask struct {
	Price  string `json:"price"`
	Amount string `json:"amount"`
}

type Bid struct {
	Price  string `json:"price"`
	Amount string `json:"amount"`
}

func main() {
	pusherClient, err := pusher.NewClient(APP_KEY)
	if err != nil {
		log.Fatalln(err)
	}
	log.Printf("Pusher client connected.")

	// Subscribe
	err = pusherClient.Subscribe("order_book")
	if err != nil {
		log.Println("Subscription error : ", err)
	}

	// Bind events
	dataChannelTrade, err := pusherClient.Bind("data")
	if err != nil {
		log.Println("Bind error: ", err)
	}
	log.Println("Binded to 'data' event")

	// Loop forever
	for {
		select {
		case dataEvt := <-dataChannelTrade:
			result := PusherOrderbook{}
			err := JSONDecode([]byte(dataEvt.Data), &result)
			if err != nil {
				log.Println(err)
			}

			nr := NormalOrderbook{}
			for i := 0; i < len(result.Asks); i++ {
				x := result.Asks[i]
				nr.Asks = append(nr.Asks, Ask{Price: x[0], Amount: x[1]})
			}
			for i := 0; i < len(result.Bids); i++ {
				x := result.Bids[i]
				nr.Bids = append(nr.Bids, Bid{Price: x[0], Amount: x[1]})
			}

			log.Printf("---- ORDER BOOK: -------")
			// log.Printf("Bids: %s", result.Bids)
			// log.Printf("Asks: %s", result.Asks)
			log.Printf("Bids: %+v", nr.Bids)
			log.Printf("Asks: %+v", nr.Asks)
		}
	}
}

func JSONDecode(data []byte, to interface{}) error {
	if !StringContains(reflect.ValueOf(to).Type().String(), "*") {
		return errors.New("json decode error - memory address not supplied")
	}
	return json.Unmarshal(data, to)
}

func StringContains(input, substring string) bool {
	return strings.Contains(input, substring)
}
