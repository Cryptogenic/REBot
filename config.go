package main 

import(
	"os"
	"fmt"
	"strconv"
	"github.com/go-ini/ini"
)

// Writes to config.ini with the given section, property, and value
func setConfigProperty(section string, prop string, value string) {
	cfg, err := ini.InsensitiveLoad("config.ini")

	if err != nil {
		fmt.Println("[ERROR] Critical error attempting to load config.ini! " + err.Error())
		os.Exit(1)
	}

	// If key exists, update, otherwise, create new key
	if cfg.Section(section).HasKey(prop) {
		cfg.Section(section).NewKey(prop, value)
	} else {
		cfg.Section(section).Key(prop).SetValue(value)
	}

	fmt.Println("[CONFIG] Section '" + section + "', Key '" + prop + "' has been set to value: '" + value + "'!")
	
	cfg.SaveTo("config.ini")
}

// Searches and reads a property from config.ini as a string
func getConfigPropertyAsStr(section string, prop string) string {
	cfg, err := ini.InsensitiveLoad("config.ini")

	if err != nil {
		fmt.Println("[ERROR] Critical error attempting to load config.ini! " + err.Error())
		os.Exit(1)
	}

	// Read cfg value as string
	val := cfg.Section(section).Key(prop).String()

	fmt.Println("[CONFIG] Read '" + section + "', Key '" + prop + "', set to: '" + val + "'!")

	return val
}

// Searches and reads a property from config.ini as an integer
func getConfigPropertyAsInt(section string, prop string) int {
	cfg, err := ini.InsensitiveLoad("config.ini")

	if err != nil {
		fmt.Println("[ERROR] Critical error attempting to load config.ini! " + err.Error())
		os.Exit(1)
	}

	// Read cfg value as integer
	val, _ := cfg.Section(section).Key(prop).Int()

	fmt.Println("[CONFIG] Read '" + section + "', Key '" + prop + "', set to: '" + strconv.Itoa(val) + "'!")

	return val
}