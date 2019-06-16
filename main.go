package main 

import(
	"fmt"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/bwmarrin/discordgo"
)

// The "StrList" custom type allows us to add our own functions to it
type StrList []string

// Globals
var (
	GuildID 		string
	DeveloperMode 	bool
	DeveloperList 	StrList
	startTime 		time.Time
)

// Main entry point
func main() {
	var bot *discordgo.Session
	var err error

	startTime = time.Now()
	botToken := getConfigPropertyAsStr("discord", "token")

	if bot, err = discordgo.New("Bot " + botToken); err != nil {
		fmt.Println("[ERROR] Critical error creating Discord session, ", err)
		return
	}

	// Handle messageCreate events sent from Discord
	bot.AddHandler(messageCreate)

	if err = bot.Open(); err != nil {
		fmt.Println("[ERROR] Critical error connecting to Discord, ", err)
		return
	}

	GuildID = ""
	DeveloperMode  = false
	DeveloperList  = []string{"165177089035599873"} // List of discord user ID's that can access developer commands

	// Build the alias maps for commands and dictionary definitions
	buildDictionaryMap()
	buildCommandMap()

	fmt.Println("[INFO] Bot is now running! Press CTRL-C to stop!")

	// Listen for kill signals
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, os.Kill)
	<-sc

	// Close the client and free resources
	_ = bot.Close()
}

// Handler for message events received from Discord
func messageCreate(s *discordgo.Session, m *discordgo.MessageCreate) {
	// Don't process commands sent from the bot
	if m.Author.ID == s.State.User.ID {
		return
	}

	// Don't handle nil messages
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
