package postal

import (
	"fmt"
	"testing"
)

func TestTokenizer(t *testing.T) {
	input := "the big bad wolf and the beautiful fox"
	tokenSize := 3
	expectedResult := []string{"the big bad", "big bad wolf", "bad wolf and", "wolf and the", "and the beautiful", "the beautiful fox"}

	actualResult := tokenize(input, tokenSize)
	if !testEq(expectedResult, actualResult) {
		t.Error("Tokenization failed for input: ", input,
			"\nThe result is different than expected")
	}
}
func TestAddressIsCorrectlyExtractedFromURL(t *testing.T) {
	cases := map[string]map[string]bool{

		// Hunter Technical Test Inputs
		"https://web.archive.org/web/20190215112322/https://hunter.io/terms-of-service": {"427 N Tatnall St #50754 Wilmington, Delaware 19801-2230, USA": true},
		"https://web.archive.org/web/20190429014959/https://www.cloudflare.com/terms/": {
			"101 Townsend St., San Francisco, CA 94107, USA":         true,
			"101 Townsend St., San Francisco, California 94107, USA": true,
			"101 Townsend St, San Francisco, CA 94107, USA":          true},

		// Other inputs
		"https://www.capitalone.com/legal/terms-conditions":                                                                {"15000 Capital One Drive Richmond, VA 23238, USA": true},
		"https://web.archive.org/web/20200102182655/https://www.adobe.com/legal/terms.html":                                {"345 Park Avenue, San Jose, California 95110-2704, USA": true, "345 Park Avenue, San Jose, California, 95110-2704, USA": true},
		"https://web.archive.org/web/20191229220121/https://www.springer.com/gp/standard-terms-and-conditions-of-business": {"233 Spring St, New York, NY, 10013, USA": true},
		"https://web.archive.org/web/20190718010541/https://www.wiley.com/en-pt/terms-of-use":                              {"111 River Street, Hoboken, NJ 07030, USA": true},
		"https://web.archive.org/web/20191230230854/https://help.netflix.com/en/node/2101":                                 {"100 Winchester Circle Los Gatos, CA 95032, USA": true},
		"https://web.archive.org/web/20200102091051/https://www.scribd.com/contact":                                        {"460 Bryant Street, #100 San Francisco, CA 94107-2594, USA": true},

		"https://web.archive.org/web/20200103180301/https://www.verizonmedia.com/policies/us/en/verizonmedia/terms/otos/index.html": {"1921 NW 87 Avenue, Doral, FL 33172, USA": true,
			"22000 AOL Way, Dulles, VA 20166, USA": true, "701 First Avenue, Sunnyvale, CA 94089, USA": true},

		// Not archived by archive.org
		"https://www.zillowgroup.com/terms-of-use-privacy-policy/zg-privacy-policy/": {
			"1301 Second Avenue, Floor 31 Seattle, WA 98101, USA":           true,
			"1301 Second Avenue, Floor 31, Seattle, Washington, 98101, USA": true,
			"1301 Second Ave., Fl. 31, Seattle, Washington 98101, USA":      true,
			"1301 Second Avenue, Floor 31, Seattle, WA 98101, USA":          true},
	}
	for url, expected := range cases {
		result := ExtractAddressFromUrl(url)
		if !cases[url][result] {
			t.Error("\nUrl:", url,
				"\nExpected one of these Addresses:\n", prettyPrintMapKey(expected),
				"\nActual Address:", result)
		}
	}
}

// prettyPrintMap retrieves the keys of a map and prints one by one, separated by a line break
// Returns the resulting string
func prettyPrintMapKey(maps map[string]bool) string {
	var pretty string
	for k, _ := range maps {
		pretty += fmt.Sprintf("\t%s \n", k)
	}
	return pretty
}

// testEq tests if two slices are equal
// Credits to Stephen Weinberg https://stackoverflow.com/a/15312097
// Returns True if the slices are equal, otherwise returns false
func testEq(a, b []string) bool {

	// If one is nil, the other must also be nil.
	if (a == nil) != (b == nil) {
		return false
	}

	if len(a) != len(b) {
		return false
	}

	for i := range a {
		if a[i] != b[i] {
			return false
		}
	}
	return true
}
