package main

import(
	"fmt"
	"os"
	"os/exec"
	"runtime"
	"strconv"
	"strings"
	"time"

	"github.com/bwmarrin/discordgo"
)

// Enables developer mode (prints debug output, etc.)
func cmdDevMode(s *discordgo.Session, m *discordgo.MessageCreate, args []string) {
	if !DeveloperMode {
		fmt.Println("[INFO] Developer mode has been enabled!!!")
		_, _ = s.ChannelMessageSend(m.ChannelID, "Developer Mode Enabled!")
		DeveloperMode = true
	} else {
		fmt.Println("[INFO] Developer mode has been disabled.")
		_, _ = s.ChannelMessageSend(m.ChannelID, "Developer Mode Disabled!")
		DeveloperMode = false
	}
}

// Runs a test feature (subject to change at any time)
func cmdTest(s *discordgo.Session, m *discordgo.MessageCreate, args []string) {
	_, _ = s.ChannelMessageSend(m.ChannelID, "Not implemented.")
}

// Simple test command. "Ping" -> "Pong!"
func cmdPing(s *discordgo.Session, m *discordgo.MessageCreate, args []string) {
	_, _ = s.ChannelMessageSend(m.ChannelID, "Pong!")
}

// Echos what the user gives
func cmdSay(s *discordgo.Session, m *discordgo.MessageCreate, args []string) {
	name := args[0]

	msg := strings.Join(args, " ")
	msg = strings.Replace(msg, name, "", 1)
	_, _ = s.ChannelMessageSend(m.ChannelID, msg)
}

// Outputs detailed memory usage statistics
func cmdMem(s *discordgo.Session, m *discordgo.MessageCreate, args []string) {
	var mem runtime.MemStats

	runtime.ReadMemStats(&mem)

	msg := "```" + 
		   "\nAllocated (Does not include free objects): " + strconv.FormatUint(mem.Alloc / 1024, 10) + " bytes" +
		   "\nTotal Allocated (Includes free objects):   " + strconv.FormatUint(mem.TotalAlloc / 1024, 10) + " bytes" +
		   "\nSystem (Includes stack, heap, .text, etc): " + strconv.FormatUint(mem.Sys / 1024, 10) + " bytes" +
		   "```"

	_, _ = s.ChannelMessageSend(m.ChannelID, msg)
}

// Gives the current uptime
func cmdUptime(s *discordgo.Session, m *discordgo.MessageCreate, args []string) {
	var embedFields []embedField

	uptime 	:= time.Since(startTime)
	days 	:= padLeft(strconv.Itoa(int(uptime.Hours())/24), "0", 2)
	hours 	:= padLeft(strconv.Itoa(int(uptime.Hours())%24), "0", 2)
	mins 	:= padLeft(strconv.Itoa(int(uptime.Minutes())%60), "0", 2)
	secs 	:= padLeft(strconv.Itoa(int(uptime.Seconds())%60), "0", 2)

	// Notify the user of the uptime
	embedUptime := embedField {
		name:   "Uptime",
		value:  days + " days, " + hours + " hours, " + mins + " minutes, " + secs + " seconds.",
		inline: false,
	}

	embedFields = append(embedFields, embedUptime)
	discordSendEmbeddedMsg(s, m.ChannelID, embedFields, "", 0x00BFFF, "https://i.imgur.com/2yCS7A7.png")
}

// Restarts the bot.
func cmdRestart(s *discordgo.Session, m *discordgo.MessageCreate, args []string) {
	var embedFields []embedField

	// Notify the user of the restart
	embedRestart := embedField{
		name:   "Restarting!",
		value:  "I'm now restarting - be right back!",
		inline: false,
	}

	embedFields = append(embedFields, embedRestart)
	discordSendEmbeddedMsg(s, m.ChannelID, embedFields, "", 0x00BFFF, "https://cdn4.iconfinder.com/data/icons/circle-blue/64/restart.png")

	cmd := exec.Command("./restart.sh")
	_ = cmd.Start()
	os.Exit(0)

}

// Kills the bot.
func cmdDie(s *discordgo.Session, m *discordgo.MessageCreate, args []string) {
	var embedFields []embedField

	// Notify the user of the restart
	embedRestart := embedField{
		name:   "Shutting Down!",
		value:  "I'm shutting down now, goodbye.",
		inline: false,
	}

	embedFields = append(embedFields, embedRestart)
	discordSendEmbeddedMsg(s, m.ChannelID, embedFields, "", 0x00BFFF, "https://cdn3.iconfinder.com/data/icons/ginux/Png/Shutdown-64.png")

	os.Exit(0)
}

