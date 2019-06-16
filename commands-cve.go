package main

import(
	"html"
	"strings"
	"github.com/bwmarrin/discordgo"
)

// Looks up a CVE
func cmdCve(params cmdArguments) {
	s := params.s
	m := params.m
	args := params.args

	reqUrl := "https://nvd.nist.gov/vuln/detail/" + args[1]
	htmlResp := getPageContents(reqUrl)

	if strings.Contains(htmlResp, "Vuln ID, expected format") {
		_, _ = s.ChannelMessageSend(m.ChannelID, "CVE ID is not valid.")
		return
	}

	// Parses information from the NVD
	cveId 		:= stribet(htmlResp, "page-header-vuln-id\">", "</span>")
	cveDesc 	:= stribet(htmlResp, "vuln-description\">", "</p>")
	cvePubDate 	:= stribet(htmlResp, "vuln-published-on\">", "</span>")
	cveUpdDate 	:= stribet(htmlResp, "vuln-last-modified-on\">", "</span>")

	cvvs3BaseScore := stribet(htmlResp, "vuln-cvssv3-base-score\">", "</span>")
	cvvs3BaseScore += "(" + stribet(htmlResp, "vuln-cvssv3-base-score-severity\">", "</span>") + ")"
	cvvs3ImpScore  := stribet(htmlResp, "vuln-cvssv3-impact-score\">", "</span>")
	cvvs3ExpScore  := stribet(htmlResp, "vuln-cvssv3-exploitability-score\">", "</span>")

	// Purify
	cveDesc = html.UnescapeString(cveDesc)

	// Create the nice discord embeds
	cveIDEmbed := discordgo.MessageEmbedField{
		Name: "CVE ID",
		Value: cveId,
		Inline: false,
	}

	cveDescEmbed := discordgo.MessageEmbedField{
		Name: "CVE Description",
		Value: cveDesc,
		Inline: false,
	}

	cveCVSS3BaseEmbed := discordgo.MessageEmbedField{
		Name: "CVSS v3.0 Base Score",
		Value: cvvs3BaseScore,
		Inline: false,
	}

	cveCVSS3ImpEmbed := discordgo.MessageEmbedField{
		Name: "CVSS v3.0 Impact Score",
		Value: cvvs3ImpScore,
		Inline: true,
	}

	cveCVSS3ExpEmbed := discordgo.MessageEmbedField{
		Name: "CVSS v3.0 Exploitability Score",
		Value: cvvs3ExpScore,
		Inline: true,
	}

	cveMoreInfoEmbed := discordgo.MessageEmbedField{
		Name: "More Information",
		Value: reqUrl,
		Inline: false,
	}

	cveFooter := discordgo.MessageEmbedFooter{
		Text: "Published: " + cvePubDate + " | Last Updated: " + cveUpdDate,
	}

	cveEmbed := discordgo.MessageEmbed{
		Type: "rich",
		Footer: &cveFooter,
	}

	cveEmbed.Fields = append(
		cveEmbed.Fields,
		&cveIDEmbed,
		&cveCVSS3BaseEmbed,
		&cveCVSS3ImpEmbed,
		&cveCVSS3ExpEmbed,
		&cveMoreInfoEmbed,
		&cveDescEmbed)

	_, _ = s.ChannelMessageSendEmbed(m.ChannelID, &cveEmbed)
}

// Looks up a given term in a dictionary
func cmdInfo(params cmdArguments) {
	var name string
	var item Info
	var err error

	s := params.s
	m := params.m
	args := params.args

	// Stitch together rest of arguments for definition
	if len(args) > 1 {
		for i := 1; i < len(args); i++ {
			name += args[i] + " "
		}
	}

	name = strings.TrimSpace(name)

	item, err = getDictionaryItem(name)

	if err != nil {
		_, _ = s.ChannelMessageSend(m.ChannelID, "There is no information on your request.")
		return
	}

	// Create the nice discord embeds
	infoNameEmbed := discordgo.MessageEmbedField{
		Name: "Name",
		Value: item.Name,
		Inline: true,
	}

	infoTypeEmbed := discordgo.MessageEmbedField{
		Name: "Type",
		Value: item.Type,
		Inline: true,
	}

	infoDescEmbed := discordgo.MessageEmbedField{
		Name: "Description",
		Value: item.Description,
		Inline: true,
	}

	infoFooter := discordgo.MessageEmbedFooter{
		Text: "Last Updated: " + item.Updated,
	}

	infoEmbed := discordgo.MessageEmbed{
		Type: "rich",
		Footer: &infoFooter,
	}

	infoEmbed.Fields = append(
		infoEmbed.Fields,
		&infoNameEmbed,
		&infoTypeEmbed,
		&infoDescEmbed)

	_, _ = s.ChannelMessageSendEmbed(m.ChannelID, &infoEmbed)
}
