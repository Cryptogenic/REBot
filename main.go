package main 

import(
	"os"
	"os/signal"
	"fmt"
	"syscall"
	"strings"

	"github.com/bwmarrin/discordgo"
)

// The "StrList" custom type allows us to add our own functions to it
type StrList []string

// Globals
var(GuildID string)
var(DeveloperMode bool)
var(DeveloperList StrList)

// Main entry point
func main() {
	// Get Discord auth token from config.ini
	botToken := getConfigPropertyAsStr("discord", "token")

	dg, err := discordgo.New("Bot " + botToken)

	if err != nil {
		fmt.Println("[ERROR] Critical error creating Discord session, ", err)
		return
	}

	// Handle messageCreate events sent from Discord
	dg.AddHandler(messageCreate)

	err = dg.Open()

	if err != nil {
		fmt.Println("[ERROR] Critical error connecting to Discord, ", err)
		return
	}

	GuildID = ""
	DeveloperMode  = false
	DeveloperList  = []string{"165177089035599873"} // List of discord user ID's that can access developer commands

	// Build the alias maps for commands and dictionary definitions
	// This is needed because if we don't call make() on the maps - a nil pointer dereference will occur
	buildDictionaryMap()
	buildCommandMap()

	fmt.Println("[INFO] Bot is now running! Press CTRL-C to stop!")

	// Listen for kill signals
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, os.Kill)
	<-sc

	// We're done with Discord - close the client and free resources
	dg.Close()
}

// Handler for message events received from Discord
func messageCreate(s *discordgo.Session, m *discordgo.MessageCreate) {
	// Don't process commands sent from the bot
	if m.Author.ID == s.State.User.ID {
		return
	}

	// Don't handle nil messages, or we'll get a nil pointer dereference panic
	if len(m.Content) <= 0 {
		return
	}

	// When the first character is the command character, parse the command and pass it off to the generic command handler
	if m.Content[0] == '!' {
		cmd := strings.Replace(m.Content, "!", "", 1)
		cmdParts := strings.Split(cmd, " ")

		command(s, m, cmdParts, cmdParts[0])
	}
}
