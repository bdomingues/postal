package postal

import (
	"fmt"
	"io/ioutil"
	"jaytaylor.com/html2text"
	"net/http"
	"os"
	"regexp"
	"strings"
	"sync"
)

// ExtractAddressFromUrl takes a url string and attempts to extract a Postal Address
// It returns the address if found, otherwise an empty string
func ExtractAddressFromUrl(url string) string {
	html, err := getHTML(url)
	if err != nil {
		fmt.Printf("Failed to retrieve HTML Code from %s", url)
		os.Exit(1)
	}
	content := prepareHTML(string(html))
	tokens := tokenize(content, 10)
	return extractAddress(tokens)
}

// extractAddress takes a slice of tokens and spawns go routines to handle the extraction
// Blocks until all go routines are finished
// It returns the found address if any is found, otherwise an empty string
func extractAddress(tokens []string) string {
	var waitgroup sync.WaitGroup
	resultChannel := make(chan string, 1) //We are only fetching the 1st result. If we wanted multiple results (Eg: A slice), we would have a bigger buffer.

	for _, token := range tokens {
		waitgroup.Add(1)
		go extract(token, &waitgroup, resultChannel)
	}
	waitgroup.Wait()
	// Make sure we don't block if there are not results
	select {
	case result := <-resultChannel:
		return result
	default:
		return ""
	}
}

// extract receives a token and tries match it against an address regex expression
// If it successfully finds an address, writes the address to a channel
// otherwise just decreases the Wait group counter
func extract(candidate string, waitgroup *sync.WaitGroup, resultChannel chan string) {
	defer waitgroup.Done()
	state := findState(candidate)
	if state == "" {
		return
	}
	street := findStreet(candidate)
	if street == "" {
		return
	}

	re := regexp.MustCompile(`(?i)\d+(,*\s*\w*){1,3}\b` + street + `\b.+\b` + state + `\b.{1,3}(\b\d{5}-\d{4}\b|\b\d{5}\b)`)
	match := re.FindString(candidate)

	if len(match) > 0 {
		match = match + ", USA" //We manually append the USA as a country because it is the only one supported at this time
		select {
		// Make sure we don't block incase the channel is full. If it does, we terminate
		case resultChannel <- match:
			return
		default:
		}
	}
}

// getHTML takes an url as a parameter and returns a slice of bytes respresenting the HTML source code
func getHTML(url string) ([]byte, error) {
	resp, err := http.Get(url)
	if err != nil {
		return nil, fmt.Errorf("GET error: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("Status error: %v", resp.StatusCode)
	}

	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("Read body: %v", err)
	}
	return data, nil
}

// prepareHTML takes a string and strips all the HTML code from it
// Carriage returns and line feeds to improve the regex matching
// Returns the stripped string
func prepareHTML(data string) string {
	content, _ := html2text.FromString(string(data))
	//Remove \r\n
	re_removeLb := regexp.MustCompile(`\r?\n+`)
	content = re_removeLb.ReplaceAllString(content, " ")
	return content
}

// tokenize takes a string and a tokenSize and creates a contiguous sequence of tokens
// An Example is the string "the big bad wolf and the great fox" and a tokenSize of 3
// 	Generated tokens:
//  	"the big bad"
//		"big bad wolf"
//		"bad wolf and"
//		"wolf and the"
//		"and the great"
//		"the great fox"
// Returns a slice with the generated tokens
func tokenize(content string, tokenSize int) []string {
	words := strings.Fields(content)
	var runningTokens []string
	var windowedTokens []string

	for _, word := range words {
		runningTokens = append(runningTokens, word)
		if len(runningTokens) == tokenSize {
			currentSlice := strings.Join(runningTokens, " ")
			windowedTokens = append(windowedTokens, currentSlice)
			runningTokens = runningTokens[1:]
		}
	}
	return windowedTokens
}

// findState takes a possible address candidate and checks if it contains a valid State
// currently only supports the US States found in usdata.go
// Returns the state if found, otherwise an empty string
func findState(candidate string) string {
	for state, _ := range usStates {
		found, _ := regexp.MatchString(`(?i)\b`+state+`\b`, candidate)
		if found {
			return state
		}
	}
	return ""
}

// findStreet takes a possible address candidate and checks if it contains a valid street
// currently only supports the US streets found in usdata.go
// Returns the street if found, otherwise an empty string
func findStreet(candidate string) string {
	for street, _ := range usStreetSuffixes {
		found, _ := regexp.MatchString(`(?i)\b`+street+`\b`, candidate)
		if found {
			return street
		}
	}
	return ""
}
