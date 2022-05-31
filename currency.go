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

type Money struct {
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
	EURBRL Money
	USDBRL Money
}

func main() {

	currency, valueToSell, delay := askUserTheInput()

	for {
		currentValue, highestValue := reachTheCurrencyAPI(currency)

		fmt.Println(currency, "now is R$", currentValue, " and the highest was R$", highestValue)

		if currentValue >= valueToSell {
			message := currency + "now is " + fmt.Sprintf("%f", currentValue)

			beeep.Notify("SELL IT NOW!!!", message, "assets/information.png")
		} else if highestValue-currentValue > 0.30 {
			beeep.Notify("WATCH OUT!!!", "The value has decreased a lot today", "assets/information.png")
		}

		time.Sleep(time.Duration(delay) * time.Minute)
	}
}

func reachTheCurrencyAPI(currency string) (float64, float64) {
	apiURL := "https://economia.awesomeapi.com.br/last/"

	switch currency {
	case "EUR":
		apiURL = apiURL + "EUR-BRL"
	case "USD":
		apiURL = apiURL + "USD-BRL"
	}

	response, _ := http.Get(apiURL)

	responseBytes, _ := io.ReadAll(response.Body)

	res := Response{}

	json.Unmarshal(responseBytes, &res)

	var bid float64
	var high float64
	if currency == "EUR" {
		bid, _ = strconv.ParseFloat(res.EURBRL.Bid, 32)
		high, _ = strconv.ParseFloat(res.EURBRL.High, 32)
	} else {
		bid, _ = strconv.ParseFloat(res.USDBRL.Bid, 32)
		high, _ = strconv.ParseFloat(res.USDBRL.High, 32)
	}

	return math.Round(bid*100) / 100, math.Round(high*100) / 100
}

func askUserTheInput() (string, float64, int) {
	var currencyOption int
	var valueToSell float64
	var delay int

	currency := "EUR"

	fmt.Println("Select the currency you wish to monitor:")
	fmt.Println("1 - EUR")
	fmt.Println("2 - USD")
	fmt.Scan(&currencyOption)

	if currencyOption == 2 {
		currency = "USD"
	}

	fmt.Println("Type the value that you intend to sell the", currency, "for")
	fmt.Scan(&valueToSell)

	fmt.Println("What is the monitoring frequency (in minutes)?")
	fmt.Scan(&delay)

	return currency, valueToSell, delay
}
