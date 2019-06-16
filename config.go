package main

import (
	"fmt"
	"os"

	"github.com/go-ini/ini"
)

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
