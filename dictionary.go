package main 

import(
	"os"
	"errors"
	"strings"
	"io/ioutil"
	"encoding/json"
)

// Each definition follows a JSON-object format, specified by "Info"
type Info struct {
	Name string `json:"name"`
	Type string `json:"type"`
	Description string `json:"description"`
	Updated string `json:"updated"`
}

// Holds a list of dictionary definitions to aliases/synonyms
var dictionaryMap map[string][]string

// Build the dictionary map
func buildDictionaryMap() {
	dictionaryMap = make(map[string][]string)

	// Add definition aliases/synonyms here
	dictionaryMap["buffer overflow"] = []string{"bof", "buf overflow"}
	dictionaryMap["uaf"] = []string{"use-after-free", "use after free"}
	dictionaryMap["integer overflow"] = []string{"int overflow"}
	dictionaryMap["arbitrary rw"] = []string{"arbitrary r/w", "r/w", "rw", "arb rw", "arb. r/w", "arb. rw"}
	dictionaryMap["code execution"] = []string{"arbitrary code execution", "code exec", "rce", "ace"}
	dictionaryMap["aslr"] = []string{"address space layout randomization"}
	dictionaryMap["dep"] = []string{"data execution prevention", "nx", "no execute", "no-execute", "W^X"}
	dictionaryMap["smap"] = []string{"supervisor mode access prevention"}
	dictionaryMap["smep"] = []string{"supervisor mode execution prevention"}
	dictionaryMap["userland"] = []string{"usermode", "ring3"}
	dictionaryMap["kernel"] = []string{"kernelmode", "ring0"}
	dictionaryMap["sandbox"] = []string{"sandboxing", "jail", "jailing"}
	dictionaryMap["rop"] = []string{"return-oriented-programming", "return oriented programming"}
	dictionaryMap["jop"] = []string{"jump-oriented-programming", "jump oriented programming"}
	dictionaryMap["uninit access"] = []string{"uninitialized read", "uninit read", "uninitialized access"}
	dictionaryMap["race"] = []string{"races", "race condition", "race condition"}
	dictionaryMap["re"] = []string{"reverse engineering", "reverse engineer"}
}

// Searches the dictionary for a given definition
func getDictionaryItem(item string) (Info, error) {
	var i Info
	var file string

	// Convert to lowercase
	item = strings.ToLower(item)

	// Check if definition exists
	_, ok := dictionaryMap[item]

	if !ok {
		// If definition doesn't exist, check if an alias exists
		foundDefinition := false

		for key, definition := range dictionaryMap {
			for _, alias := range definition {
				if item == alias {
					file = key
					foundDefinition = true
				}
			}
		}

		if !foundDefinition {
			return Info{}, errors.New("definition could not be found")
		}
	} else {
		file = item
	}

	// Find the definition file and decode it
	filePath := "./dictionary/" + file + ".json"
	
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		return Info{}, err
	}

	raw, err := ioutil.ReadFile(filePath)

	if err != nil {
		return Info{}, err
	}

	err = json.Unmarshal(raw, &i)

	if err == nil {
		return i, err
	} else {
		return Info{}, err
	}
}