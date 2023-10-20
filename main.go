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

func throwLogError(msg string) {
	// load the token and client init for Loggly
	client := loggly.New(os.Getenv("LOGGLY_TOKEN"))
	logErr := client.EchoSend("error", msg)
	if logErr != nil {
		os.Exit(1)
	}
}

func pollData() {
	fmt.Println("Starting...")

	// Call CheapShark API
	resp, err := http.Get("https://www.cheapshark.com/api/1.0/deals?storeID=1&sortBy=Recent&steamworks=1&onSale=1&hideDuplicates=1&pageSize=10")
	if err != nil {
		throwLogError("Could not pull the data from the CheapSharkAPI.")
		panic(err)
	}

	defer resp.Body.Close()
	fmt.Println("Response Status:", resp.Status)

	// Read the response body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		throwLogError("Could not read the data from the CheapSharkAPI.")
		log.Fatal(err)
	}

	// Parse the JSON and display info in the terminal
	var apidata APIData
	parsedata := json.Unmarshal(body, &apidata)
	if err != nil {
		log.Fatal(parsedata)
	}

	formattedData, _ := json.MarshalIndent(apidata, "", "  ")
	fmt.Println(string(formattedData))

	// Send a success message to loggly
	client := loggly.New(os.Getenv("LOGGLY_TOKEN"))
	var respSize string = strconv.Itoa(len(body))
	log := client.EchoSend("info", "Successful data collection of size: "+respSize)
	if log != nil {
		client.EchoSend("error", "Failed with error: "+log.Error())
	}

	// Initialize a AWS session
	sess := session.Must(session.NewSessionWithOptions(session.Options{
		SharedConfigState: session.SharedConfigEnable,
	}))

	// Create a DynamoDB client
	svc := dynamodb.New(sess)

	// Create item to be added to DynamoDB
	av, err := dynamodbattribute.MarshalMap(apidata)
	fmt.Printf("AV:\t%+v\n", av)
	if err != nil {
		fmt.Println("Got error marshalling item: ", err.Error())
		os.Exit(1)
	}

	//Input an item in test-table-temokpae
	tableName := "test-table-temokpae"

	_, err = svc.PutItem(&dynamodb.PutItemInput{
		Item: map[string]*dynamodb.AttributeValue{
			"InternalName": {
				S: aws.String(apidata.InternalName),
			},
			"Title": {
				S: aws.String(apidata.Title),
			},
			"MetacriticLink": {
				S: aws.String(apidata.MetacriticLink),
			},
			"DealID": {
				S: aws.String(apidata.DealID),
			},
			"StoreID": {
				S: aws.String(apidata.StoreID),
			},
			"GameID": {
				S: aws.String(apidata.GameID),
			},
			"SalePrice": {
				S: aws.String(apidata.SalePrice),
			},
			"NormalPrice": {
				S: aws.String(apidata.NormalPrice),
			},
			"IsOnSale": {
				S: aws.String(apidata.IsOnSale),
			},
			"Savings": {
				S: aws.String(apidata.Savings),
			},
			"MetacriticScore": {
				S: aws.String(apidata.MetacriticScore),
			},
			"SteamRatingText": {
				S: aws.String(apidata.SteamRatingText),
			},
			"SteamRatingPercent": {
				S: aws.String(apidata.SteamRatingPercent),
			},
			"SteamRatingCount": {
				S: aws.String(apidata.SteamRatingCount),
			},
			"SteamAppID": {
				S: aws.String(apidata.SteamAppID),
			},
			"ReleaseDate": {
				S: aws.String(apidata.ReleaseDate),
			},
			"LastChange": {
				S: aws.String(apidata.LastChange),
			},
			"DealRating": {
				S: aws.String(apidata.DealRating),
			},
			"Thumb": {
				S: aws.String(apidata.Thumb),
			},
		},
		TableName: aws.String(tableName),
	})

	if err != nil {
		fmt.Println("Could not make a DynamoDB table entry.")
		os.Exit(1)
	}

	if err != nil {
		fmt.Println("Got error calling PutItem: ", err.Error())
		os.Exit(1)
	}

	// Send a success message to DyanmoDB
	fmt.Print("Successfully added to DynamoDB table\n")
}

func main() {
	for range time.Tick(time.Minute * 1) {
		pollData()
	}
}
