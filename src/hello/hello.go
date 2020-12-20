package main

import (
	"fmt"
	"log"

	"example.com/greetings"
)

func main() {

	log.SetPrefix("greetings: ")
	log.SetFlags(0)

	fmt.Println("hola mundo 3!")
	fmt.Println("tierra ")
	fmt.Println("Otra cosa2")
	message, err := greetings.Hello("fadi")
	fmt.Println(message)
	// Request a greeting message.
	message, err = greetings.Hello("")
	// If an error was returned, print it to the console and
	// exit the program.
	if err != nil {
		log.Fatal(err)
	}
	// If no error was returned, print the returned message
	// to the console.
	fmt.Println(message)

}
