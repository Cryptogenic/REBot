package main 

import(
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"strings"

	"github.com/bwmarrin/discordgo"
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

// Used to define an embed field to pass to discordSendEmbeddedMsg().
type embedField struct {
	name   string
	value  string
	inline bool
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

// Uses HTTP to fetch page contents and write it to a file
func downloadPageContents(filename string, url string) error {
	resp, err := http.Get(url)
	if err != nil {
		return err
	}

	defer resp.Body.Close()

	outFile, err := os.Create(filename)
	if err != nil {
		return err
	}

	defer outFile.Close()

	_, err = io.Copy(outFile, resp.Body)
	return err
}

// Creates an embed message to send in Discord.
func discordSendEmbeddedMsg(s *discordgo.Session, channel string, sections []embedField, footer string, color int, thumbnail string) {
	embed := discordgo.MessageEmbed{
		Type:  "rich",
		Color: color,
	}

	if footer != "" {
		embedFooter := discordgo.MessageEmbedFooter{
			Text: footer,
		}

		embed.Footer = &embedFooter
	}

	if thumbnail != "" {
		embedThumbnail := discordgo.MessageEmbedThumbnail{
			URL:    thumbnail,
			Width:  8,
			Height: 8,
		}

		embed.Thumbnail = &embedThumbnail
	}

	for _, section := range sections {
		embedSection := discordgo.MessageEmbedField{
			Name:   section.name,
			Value:  section.value,
			Inline: section.inline,
		}

		embed.Fields = append(
			embed.Fields,
			&embedSection)
	}

	_, _ = s.ChannelMessageSendEmbed(channel, &embed)
}

// Creates an embed message to send in Discord, but automatically creates 1 embed field. Good for quick, one field messages like errors.
func discordSendQuickEmbeddedMsg(s *discordgo.Session, channel string, title string, body string, footer string, color int, thumbnail string) {
	embed := discordgo.MessageEmbed{
		Type:  "rich",
		Color: color,
	}

	if footer != "" {
		embedFooter := discordgo.MessageEmbedFooter{
			Text: footer,
		}

		embed.Footer = &embedFooter
	}

	if thumbnail != "" {
		embedThumbnail := discordgo.MessageEmbedThumbnail{
			URL:    thumbnail,
			Width:  8,
			Height: 8,
		}

		embed.Thumbnail = &embedThumbnail
	}

	embedSection := discordgo.MessageEmbedField{
		Name:   title,
		Value:  body,
		Inline: false,
	}

	embed.Fields = append(
		embed.Fields,
		&embedSection)

	_, _ = s.ChannelMessageSendEmbed(channel, &embed)
}


