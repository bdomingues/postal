// Package Main is an example for the Postal package. Asks an URL input from the user and tries to extract an address
// Written by Bernardo Domingues <bernardo.domis@gmail.com>
package main

import (
	"fmt"
	"github.com/bdomingues/postal/postal"
)

func main() {
	fmt.Print("Enter url: ")
	var url string
	fmt.Scanln(&url)
	fmt.Println(postal.ExtractAddressFromUrl(url))
}
