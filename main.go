// Go Program
package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"github.com/jamespearly/loggly"
)

// Define a struct to store the collected data
type APIData struct {
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
	fmt.Println("Starting...")

	// Instantiate the Loggly client
	client := loggly.New(os.Getenv("LOGGLY_TOKEN"))

	// Call CheapShark API
	resp, err := http.Get("https://www.cheapshark.com/api/1.0/deals?storeID=1&sortBy=Recent&steamworks=1&onSale=1&hideDuplicates=1&pageSize=10")
	if err != nil {
		panic(err)
	}

	defer resp.Body.Close()
	fmt.Println("Response Status:", resp.Status)

	// Read the response body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}

	// Parse the JSON and display info in the terminal
	var apidata []APIData
	parsedata := json.Unmarshal(body, &apidata)
	if parsedata != nil {
		log.Fatal(parsedata)
	}

	formattedData, _ := json.MarshalIndent(apidata, "", "  ")
	fmt.Println(string(formattedData))

	// Send a success message to loggly
	var respSize string = strconv.Itoa(len(body))
	log := client.EchoSend("info", "Successful data collection of size: "+respSize)
	if log != nil {
		client.EchoSend("error", "Could not send data collection."+log.Error())
		os.Exit(1)
	}

	// Initialize a AWS session
	sess, err := session.NewSession(&aws.Config{
		Region: aws.String("us-east-1")},
	)

	if err != nil {
		client.EchoSend("error", "Got error initializing AWS: "+err.Error())
		os.Exit(1)
	}

	// Create DynamoDB client
	svc := dynamodb.New(sess)

	//Input an item in test-table-temokpae
	tableName := "test-table-temokpae"

	// Create item to be added to DynamoDB
	for _, item := range apidata {
		av, err := dynamodbattribute.MarshalMap(item)
		if err != nil {
			client.EchoSend("error", "Got error marshalling map: "+err.Error())
			os.Exit(1)
		}

		input := &dynamodb.PutItemInput{
			Item:      av,
			TableName: aws.String(tableName),
		}

		_, err = svc.PutItem(input)
		if err != nil {
			client.EchoSend("error", "Got error calling PutItem: "+err.Error())
			os.Exit(1)
		}

		fmt.Println("Successfully added item to DynamoDB:", item.InternalName)
	}

	// Send a Success message about DynamoDB to Loggly
	log = client.EchoSend("info", "Successfully added all the game data into DynamoDB.")
	if log != nil {
		client.EchoSend("error", "Error adding game data into DynamoDB: "+log.Error())
	}
}

func main() {
	time.Sleep(15 * time.Minute)
	pollData()
}
