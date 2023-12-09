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

// Instantiate the Loggly client
var client = loggly.New(os.Getenv("LOGGLY_TOKEN"))

func retrieveAPI() *http.Response {
	fmt.Println("Starting...")

	// Gets the API endpoint
	resp, err := http.Get("https://www.cheapshark.com/api/1.0/deals?storeID=1&sortBy=Recent&steamworks=1&onSale=1&hideDuplicates=1")
	if err != nil {
		client.EchoSend("error", "Could not retrieve API."+err.Error())
	}

	fmt.Println("Response Status:", resp.Status)

	// If the response is not OK, send an error
	if resp.StatusCode != http.StatusOK {
		client.EchoSend("error", "Status code is not OK.")
	}

	return resp
}

func readAndParseJSON(resp *http.Response) []APIData {
	// Read the response body
	body, parsedata := io.ReadAll(resp.Body)
	if parsedata != nil {
		log.Fatal("Error reading response body:", parsedata.Error())
	}

	// Define an array of APIData
	var apidata []APIData

	// Parse the JSON and display info in the terminal
	parsedata = json.Unmarshal(body, &apidata)
	if parsedata != nil {
		client.EchoSend("error", "Could not parse data."+parsedata.Error())
	}

	// Limit the length of the apidata array
	if len(apidata) > 10 {
		apidata = apidata[:10]
	}

	// Format the data to be more readable
	formattedData, _ := json.MarshalIndent(apidata, "", "  ")
	fmt.Println(string(formattedData))

	// Send a success message to loggly
	var respSize string = strconv.Itoa(len(body))
	log := client.EchoSend("info", "Successful data collection of size: "+respSize)
	if log != nil {
		client.EchoSend("error", "Could not send data collection."+log.Error())
	}

	return apidata
}

func storeDynamoDB(apidata []APIData) {
	// Initialize a AWS session
	sess := session.Must(session.NewSessionWithOptions(session.Options{
		Config: aws.Config{
			Region: aws.String("us-east-1"),
		},
		SharedConfigState: session.SharedConfigEnable,
	}))

	// Create DynamoDB client
	svc := dynamodb.New(sess)

	//Table name
	tableName := "test-table-temokpae"

	// Displays items added to DynamoDB
	for _, item := range apidata {
		av, err := dynamodbattribute.MarshalMap(item)
		if err != nil {
			client.EchoSend("error", "Got error marshalling map: "+err.Error())
		}

		input := &dynamodb.PutItemInput{
			Item:      av,
			TableName: aws.String(tableName),
		}

		_, err = svc.PutItem(input)
		if err != nil {
			client.EchoSend("error", "Got error calling PutItem: "+err.Error())
		}

		fmt.Println("Successfully added item to DynamoDB:", item.InternalName)
	}

	// Send a Success message about DynamoDB to Loggly
	log := client.EchoSend("info", "Successfully added all the game data into DynamoDB.")
	if log != nil {
		client.EchoSend("error", "Error adding game data into DynamoDB: "+log.Error())
	}
}

func pollData() {
	// Call the API function
	resp := retrieveAPI()

	// Call the readAndParse function
	apiData := readAndParseJSON(resp)

	// Close the response body
	resp.Body.Close()

	// Store the data in DynamoDB
	storeDynamoDB(apiData)
}

func main() {
	// Loop to run pollData at intervals
	for {
		pollData()

		// Sleep duration between each poll
		time.Sleep(time.Duration(15) * time.Minute)
	}
}
