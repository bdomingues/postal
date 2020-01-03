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
