package main

import (
	"fmt"
	"os"
	"os/exec"
	"runtime"
	"strconv"
	"strings"
	"time"
)

// Enables developer mode (prints debug output, etc.)
func cmdDevMode(params cmdArguments) {
	s := params.s
	m := params.m

	if !DeveloperMode {
		DeveloperMode = true

		fmt.Println("[INFO] Developer mode has been enabled!!!")
		_, _ = s.ChannelMessageSend(m.ChannelID, "Developer Mode Enabled!")
	} else {
		DeveloperMode = false

		fmt.Println("[INFO] Developer mode has been disabled.")
		_, _ = s.ChannelMessageSend(m.ChannelID, "Developer Mode Disabled!")
	}
}

// Runs a test feature (subject to change at any time)
func cmdTest(params cmdArguments) {
	s := params.s
	m := params.m

	_, _ = s.ChannelMessageSend(m.ChannelID, "Not implemented.")
}

// Simple test command. "Ping" -> "Pong!"
func cmdPing(params cmdArguments) {
	s := params.s
	m := params.m

	_, _ = s.ChannelMessageSend(m.ChannelID, "Pong!")
}

// Echos what the user gives
func cmdSay(params cmdArguments) {
	s := params.s
	m := params.m
	args := params.args

	name := args[0]
	msg  := strings.Join(args, " ")
	msg   = strings.Replace(msg, name, "", 1)
	_, _  = s.ChannelMessageSend(m.ChannelID, msg)
}

// Outputs detailed memory usage statistics
func cmdMem(params cmdArguments) {
	var mem runtime.MemStats
	var embedFields []embedField

	s := params.s
	m := params.m

	runtime.ReadMemStats(&mem)

	embedTitle := embedField{
		name:   "Memory Usage",
		value:  "-",
		inline: false,
	}

	embedAllocated := embedField{
		name:   "Allocated (not including free objects)",
		value:  strconv.FormatUint(mem.Alloc/1024, 10) + " bytes",
		inline: false,
	}

	embedTotalAllocated := embedField{
		name:   "Total Allocated (including free objects)",
		value:  strconv.FormatUint(mem.TotalAlloc/1024, 10) + " bytes",
		inline: false,
	}

	embedSystem := embedField{
		name:   "System (including stack, heap, .text, etc.)",
		value:  strconv.FormatUint(mem.Sys/1024, 10) + " bytes",
		inline: false,
	}

	embedFields = append(embedFields, embedTitle, embedAllocated, embedTotalAllocated, embedSystem)
	discordSendEmbeddedMsg(s, m.ChannelID, embedFields, "", 0x57D5FF, "https://secure.webtoolhub.com/static/resources/icons/set8/9b41f8b3af63.png")
}

// Gives the current uptime
func cmdUptime(params cmdArguments) {
	var embedFields []embedField

	s := params.s
	m := params.m

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
func cmdRestart(params cmdArguments) {
	var embedFields []embedField

	s := params.s
	m := params.m

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
func cmdDie(params cmdArguments) {
	var embedFields []embedField

	s := params.s
	m := params.m

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

