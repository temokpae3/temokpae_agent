// This package is used to test the Loggly package
package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"time"

	"github.com/jamespearly/loggly"
)

// Define a struct to store the collected data
type APIData struct {
	LargeCapsules []interface{} `json:"large_capsules"`
	FeaturedWin   []struct {
		ID                      int    `json:"id"`
		Type                    int    `json:"type"`
		Name                    string `json:"name"`
		Discounted              bool   `json:"discounted"`
		DiscountPercent         int    `json:"discount_percent"`
		OriginalPrice           int    `json:"original_price"`
		FinalPrice              int    `json:"final_price"`
		Currency                string `json:"currency"`
		LargeCapsuleImage       string `json:"large_capsule_image"`
		SmallCapsuleImage       string `json:"small_capsule_image"`
		WindowsAvailable        bool   `json:"windows_available"`
		MacAvailable            bool   `json:"mac_available"`
		LinuxAvailable          bool   `json:"linux_available"`
		StreamingvideoAvailable bool   `json:"streamingvideo_available"`
		DiscountExpiration      int    `json:"discount_expiration,omitempty"`
		HeaderImage             string `json:"header_image"`
		ControllerSupport       string `json:"controller_support,omitempty"`
	} `json:"featured_win"`
	FeaturedMac []struct {
		ID                      int    `json:"id"`
		Type                    int    `json:"type"`
		Name                    string `json:"name"`
		Discounted              bool   `json:"discounted"`
		DiscountPercent         int    `json:"discount_percent"`
		OriginalPrice           int    `json:"original_price"`
		FinalPrice              int    `json:"final_price"`
		Currency                string `json:"currency"`
		LargeCapsuleImage       string `json:"large_capsule_image"`
		SmallCapsuleImage       string `json:"small_capsule_image"`
		WindowsAvailable        bool   `json:"windows_available"`
		MacAvailable            bool   `json:"mac_available"`
		LinuxAvailable          bool   `json:"linux_available"`
		StreamingvideoAvailable bool   `json:"streamingvideo_available"`
		DiscountExpiration      int    `json:"discount_expiration,omitempty"`
		HeaderImage             string `json:"header_image"`
		ControllerSupport       string `json:"controller_support,omitempty"`
	} `json:"featured_mac"`
	FeaturedLinux []struct {
		ID                      int         `json:"id"`
		Type                    int         `json:"type"`
		Name                    string      `json:"name"`
		Discounted              bool        `json:"discounted"`
		DiscountPercent         int         `json:"discount_percent"`
		OriginalPrice           interface{} `json:"original_price"`
		FinalPrice              int         `json:"final_price"`
		Currency                string      `json:"currency"`
		LargeCapsuleImage       string      `json:"large_capsule_image"`
		SmallCapsuleImage       string      `json:"small_capsule_image"`
		WindowsAvailable        bool        `json:"windows_available"`
		MacAvailable            bool        `json:"mac_available"`
		LinuxAvailable          bool        `json:"linux_available"`
		StreamingvideoAvailable bool        `json:"streamingvideo_available"`
		HeaderImage             string      `json:"header_image"`
		DiscountExpiration      int         `json:"discount_expiration,omitempty"`
		ControllerSupport       string      `json:"controller_support,omitempty"`
	} `json:"featured_linux"`
	Layout string `json:"layout"`
	Status int    `json:"status"`
}

func pollData() {
	// Tag + client init for Loggly
	client := loggly.New("LogglyToken")

	// Call CheapShark API
	resp, err := http.Get("https://www.cheapshark.com/api/1.0/deals?storeID=1&sortBy=Recent&steamworks=1&onSale=1&hideDuplicates=1&pageSize=10")

	if err != nil {
		client.EchoSend("error", "Failed with error: "+err.Error())
	}

	defer resp.Body.Close()

	// Read the response body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		client.EchoSend("error", "Failed with error: "+err.Error())
	}

	// Parse the JSON and display some info to the terminal
	var apidata APIData
	json.Unmarshal(body, &apidata)
	formattedData, _ := json.MarshalIndent(apidata, "", "  ")
	fmt.Println(string(formattedData))

	// Send success message to loggly with response size
	var respSize string = strconv.Itoa(len(body))
	logErr := client.EchoSend("info", "Successful data collection of size: "+respSize)
	if logErr != nil {
		fmt.Println("err: ", logErr)
	}
}

func main() {
	for range time.Tick(time.Minute * 1) {
		pollData()
	}
}
