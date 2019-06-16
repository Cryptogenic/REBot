package main

import(
	"github.com/bwmarrin/discordgo"
)

// Packs command arguments into one struct to avoid unused argument warnings
type cmdArguments struct {
	s 	 *discordgo.Session
	m 	 *discordgo.MessageCreate
	args []string
}

// All command handlers will use this signature for consistency
type cmdHandler func(args cmdArguments)

// Represents a command - each command must have it's own unique Command object
type Command struct {
	name string
	aliases []string
	requiredArgs int
	usage string
	handler cmdHandler
	dev bool
}

// Stores the list of command names to Command objects
var commandMap map[string]Command

// Command handler; parses the name and passes the arguments on to the correct handler
func command(s *discordgo.Session, m *discordgo.MessageCreate, args []string, cmd string) {
	var command Command
	var ok bool

	// If we can't find the command, look for an alias
	if command, ok = commandMap[cmd]; !ok {
		foundCmd := false

		for _, commandCheck := range commandMap {
			if searchAliases(cmd, commandCheck.aliases) {
				command  = commandCheck
				foundCmd = true
				break
			}
		}

		if !foundCmd {
			return
		}
	}

	// If it's a dev only command, check if the sender actually has permissions
	if command.dev && !DeveloperList.contains(m.Author.ID) {
		return
	}

	// Ensure the required argument count is met
	if len(args) < command.requiredArgs {
		_, _ = s.ChannelMessageSend(m.ChannelID, "Usage: !" + command.name + " " + command.usage)
		return
	}

	// All good, call handler
	command.handler(cmdArguments{s, m, args})
}

// Search for an alias
func searchAliases(query string, aliases []string) bool {
	for _, alias := range aliases {
		if alias == query {
			return true
		}
	}

	return false
}

// Build the command list
func buildCommandMap() {
	commandMap = make(map[string]Command)

	// Generic Commands (can be used by anyone)

	addCommand("assemble",
		[]string{"asm", "a"},
		3,
		"[architecture] {instructions ...}",
		cmdAssemble,
		false)

	addCommand("disassemble",
		[]string{"disasm", "disas", "d"},
		3,
		"[architecture] {opcodes ...}",
		cmdDisassemble,
		false)

	addCommand("cve",
		[]string{},
		2,
		"[CVE Identifier]",
		cmdCve,
		false)

	addCommand("info",
		[]string{},
		2,
		"[term]",
		cmdInfo,
		false)

	addCommand("manual",
		[]string{},
		2,
		"[architecture]",
		cmdManual,
		false)

	addCommand("retrick",
		[]string{},
		0,
		"",
		cmdReTrick,
		false)

	addCommand("exploittrick",
		[]string{"expltrick"},
		0,
		"",
		cmdExploitTrick,
		false)

	addCommand("commands",
		[]string{"cmds"},
		0,
		"",
		cmdCommands,
		false)

	addCommand("motivation",
		[]string{"motivateme"},
		0,
		"",
		cmdMotivation,
		false)

	/*addCommand("readelf",
		[]string{"elf"},
		2,
		"[link] {options ...}",
		cmdReadelf,
		false)*/

	// Developer Only Commands

	addCommand("devmode",
		[]string{"developermode"},
		0,
		"",
		cmdDevMode,
		true)

	addCommand("mem",
		[]string{"memory"},
		0,
		"",
		cmdMem,
		true)

	addCommand("devmode",
		[]string{"developermode"},
		0,
		"",
		cmdDevMode,
		true)

	addCommand("test",
		[]string{},
		0,
		"",
		cmdTest,
		true)

	addCommand("ping",
		[]string{},
		0,
		"",
		cmdPing,
		true)

	addCommand("say",
		[]string{},
		2,
		"",
		cmdSay,
		true)

	addCommand("uptime",
		[]string{"up"},
		0,
		"",
		cmdUptime,
		true)

	addCommand("restart",
		[]string{"refresh"},
		0,
		"",
		cmdRestart,
		true)

	addCommand("die",
		[]string{"kill"},
		0,
		"",
		cmdDie,
		true)
}

// Adds a command to the command map
func addCommand(name string, aliases []string, requiredArgs int, usage string, handler cmdHandler, devOnly bool) {
	cmd := Command {
		name: name,
		aliases: aliases,
		requiredArgs: requiredArgs,
		usage: usage,
		handler: handler,
		dev: devOnly}

	commandMap[name] = cmd
}

// Sends a list of commands
func cmdCommands(params cmdArguments) {
	s := params.s
	m := params.m

	commands := "```"
	commands += "!assemble/asm [architecture] {instructions ...} - Assembles given instructions into opcodes. Instructions are separated by a ';'.\n"
	commands += "!disassemble/disasm [architecture] {opcodes ...} - Disassembles given opcodes into instructions. Give in 'bb' format separated by a space.\n"
	commands += "!cve [cve identifier] - Displays information on a given CVE from NVD.\n"
	commands += "!info [identifier] - Gives information on the given word (like a dictionary).\n"
	commands += "!retrick - Gives you a random RE trick.\n"
	commands += "!expltrick = Gives you a random exploit dev trick.\n"
	commands += "!manual [architecture] - Links a PDF manual for the given architecture.\n"
	commands += "!motivation - you can do it!\n"
	//commands += "!readelf [link] {options ...} - Reads and gives information about the ELF given by the link.\n"
	commands += "!commands/cmds - You are here.\n"
	commands += "```"

	_, _ = s.ChannelMessageSend(m.ChannelID, "Here's a list of my commands: " + commands)
}

// Motivation!
func cmdMotivation(params cmdArguments) {
	s := params.s
	m := params.m

	motivationalJapaneseFisherman := "https://www.youtube.com/watch?v=0Lq0d-cPpS4"
	_, _ = s.ChannelMessageSend(m.ChannelID, motivationalJapaneseFisherman)
}