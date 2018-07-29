package main 

import(
	"strings"
	"io/ioutil"
	"net/http"
)

// Checks if a []string array contains a value
func (strl StrList) contains(str string) bool {
	for _, s := range strl {
		if s == str {
			return true
		}
	}

	return false
}

// Returns a sub-string between two deliminators of the core string
func stribet(str string, delimLeft string, delimRight string) string {
	// Start reading from where the first deliminator starts plus the length of the deliminator
	posFirst := strings.Index(str, delimLeft)
    if posFirst == -1 {
        return ""
    }

    posFirst += len(delimLeft)

    newStr  := str[posFirst:]

    // Now find the position of where the second deliminator starts
    posLast := strings.Index(newStr, delimRight)

    if posLast == -1 {
        return ""
    }

    return newStr[0:posLast]
}

// Pads the left of the string with 'pad' for 'length' bytes
func padLeft(str string, pad string, length int) string {
	finalStr := ""

	lenPad := length - len(str)

	if lenPad > 0 {
		finalStr = strings.Repeat(pad, lenPad)
	}

	finalStr += str

	return finalStr
}

// Pads the right of the string with 'pad' for 'length' bytes
func padRight(str string, pad string, length int) string {
	finalStr := str

	lenPad := length - len(str)

	if lenPad > 0 {
		finalStr += strings.Repeat(pad, lenPad)
	}

	return finalStr
}

// Uses HTTP to get page contents of the given URL
func getPageContents(url string) string {
	resp, err := http.Get(url)

	if err != nil {
		return ""
	}

	defer resp.Body.Close()

	html, err := ioutil.ReadAll(resp.Body)

	if err != nil {
		return ""
	}

	return string(html[:])
}
