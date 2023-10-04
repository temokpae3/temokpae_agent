// This package is used to test the Loggly package
package main

import (
	"fmt"

	loggly "github.com/jamespearly/loggly"
)

func main() {

	var tag string
	tag = "My-Go-Demo"

	// Instantiate the client
	client := loggly.New(tag)

	// Valid EchoSend (message echoed to console and no error returned)
	err := client.EchoSend("info", "Good morning!")
	fmt.Println("err:", err)

	// Valid Send (no error returned)
	err = client.Send("error", "Good morning! No echo.")
	fmt.Println("err:", err)

	// Invalid EchoSend -- message level error
	err = client.EchoSend("blah", "blah")
	fmt.Println("err:", err)

}
