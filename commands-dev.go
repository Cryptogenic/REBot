package main

import(
	"fmt"
	"strings"
	"strconv"
	"runtime"
	"github.com/bwmarrin/discordgo"
)

// Enables developer mode (prints debug output, etc.)
func cmdDevMode(s *discordgo.Session, m *discordgo.MessageCreate, args []string) {
	if !DeveloperMode {
		fmt.Println("[INFO] Developer mode has been enabled!!!")
		s.ChannelMessageSend(m.ChannelID, "Developer Mode Enabled!")
		DeveloperMode = true
	} else {
		fmt.Println("[INFO] Developer mode has been disabled.")
		s.ChannelMessageSend(m.ChannelID, "Developer Mode Disabled!")
		DeveloperMode = false
	}
}

// Runs a test feature (subject to change at any time)
func cmdTest(s *discordgo.Session, m *discordgo.MessageCreate, args []string) {
	s.ChannelMessageSend(m.ChannelID, "Not implemented.")
}

// Simple test command. "Ping" -> "Pong!"
func cmdPing(s *discordgo.Session, m *discordgo.MessageCreate, args []string) {
	s.ChannelMessageSend(m.ChannelID, "Pong!")
}

// Echos what the user gives
func cmdSay(s *discordgo.Session, m *discordgo.MessageCreate, args []string) {
	name := args[0]

	msg := strings.Join(args, " ")
	msg = strings.Replace(msg, name, "", 1)
	s.ChannelMessageSend(m.ChannelID, msg)
}

// Outputs detailed memory usage statistics
func cmdMem(s *discordgo.Session, m *discordgo.MessageCreate, args []string) {
	var mem runtime.MemStats

	runtime.ReadMemStats(&mem)

	msg := "```" + 
		   "\nAllocated (Does not include free objects): " + strconv.FormatUint((mem.Alloc / 1024), 10) + " bytes" +
		   "\nTotal Allocated (Includes free objects):   " + strconv.FormatUint((mem.TotalAlloc / 1024), 10) + " bytes" +
		   "\nSystem (Includes stack, heap, .text, etc): " + strconv.FormatUint((mem.Sys / 1024), 10) + " bytes" +
		   "```"

	s.ChannelMessageSend(m.ChannelID, msg) 
}