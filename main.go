// This package is used to test the Loggly package
package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/jamespearly/loggly"
)

// Define a struct to store the collected data
type APIData []struct {
	InternalName       string `json:"internalName"`
	Title              string `json:"title"`
	MetacriticLink     string `json:"metacriticLink"`
	DealID             string `json:"dealID"`
	StoreID            string `json:"storeID"`
	GameID             string `json:"gameID"`
	SalePrice          string `json:"salePrice"`
	NormalPrice        string `json:"normalPrice"`
	IsOnSale           string `json:"isOnSale"`
	Savings            string `json:"savings"`
	MetacriticScore    string `json:"metacriticScore"`
	SteamRatingText    string `json:"steamRatingText"`
	SteamRatingPercent string `json:"steamRatingPercent"`
	SteamRatingCount   string `json:"steamRatingCount"`
	SteamAppID         string `json:"steamAppID"`
	ReleaseDate        int    `json:"releaseDate"`
	LastChange         int    `json:"lastChange"`
	DealRating         string `json:"dealRating"`
	Thumb              string `json:"thumb"`
}

func pollData() {
	// load the token and client init for Loggly
	var token = flag.String("LOGGLY_TOKEN", "", "LOGGLY_TOKEN")

	client := loggly.New(os.Getenv(*token))

	// Call CheapShark API
	resp, err := http.Get("https://www.cheapshark.com/api/1.0/deals?storeID=1&sortBy=Recent&steamworks=1&onSale=1&hideDuplicates=1&pageSize=10")

	if err != nil {
		client.EchoSend("error", "Failed with error: "+err.Error())
	}

	defer resp.Body.Close()

	// Read response body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		client.EchoSend("error", "Failed with error: "+err.Error())
	}

	// Parse the JSON and display info in the terminal
	var apidata APIData
	json.Unmarshal(body, &apidata)
	formattedData, _ := json.MarshalIndent(apidata, "", "  ")
	fmt.Println(string(formattedData))

	// Send a success message to loggly
	var respSize string = strconv.Itoa(len(body))
	log := client.EchoSend("info", "Successful data collection of size: "+respSize)
	if log != nil {
		client.EchoSend("error", "Failed with error: "+log.Error())
	}
}

func main() {
	for range time.Tick(time.Minute * 1) {
		pollData()
	}
}
