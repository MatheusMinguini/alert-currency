package main

import (
	"encoding/json"
	"fmt"
	"io"
	"math"
	"net/http"
	"strconv"
	"time"

	"github.com/gen2brain/beeep"
)

const delay = 5

type EURBRL struct {
	Code        string
	Codein      string
	Name        string
	High        string
	Low         string
	VarBid      string
	PctChange   string
	Bid         string
	Ask         string
	Timestamp   string
	Create_date string
}

type Response struct {
	EURBRL EURBRL
}

func main() {
	var valueToSell float64

	fmt.Println("Type the value that you intend to sell the EUR for")
	fmt.Scan(&valueToSell)

	for {
		currentValue, highestValue := reachTheCurrencyAPI()

		fmt.Println("EUR now is R$", currentValue, " and the highest was R$", highestValue)

		if currentValue >= valueToSell {
			message := "Euro now is " + fmt.Sprintf("%f", currentValue)
			//strconv.FormatFloat(currentValue, 'E', -1, 64)
			beeep.Notify("SELL IT NOW!!!", message, "assets/information.png")
		}

		time.Sleep(delay * time.Second)
	}
}

func reachTheCurrencyAPI() (float64, float64) {
	apiURL := "https://economia.awesomeapi.com.br/last/EUR-BRL"
	response, _ := http.Get(apiURL)

	responseBytes, _ := io.ReadAll(response.Body)

	res := Response{}

	json.Unmarshal(responseBytes, &res)

	bid, _ := strconv.ParseFloat(res.EURBRL.Bid, 32)
	high, _ := strconv.ParseFloat(res.EURBRL.High, 32)

	return math.Round(bid*100) / 100, math.Round(high*100) / 100
}
