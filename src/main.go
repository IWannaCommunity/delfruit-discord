package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"strings"

	"github.com/bwmarrin/discordgo"
)

const delfruitAPIurl = "https://delicious-fruit.com/api/game.php"

var (
	discord           *discordgo.Session
	delfruitAPIkey    string
	discordAPIkey     string
	subscribedChannel string
	notificationGuild string
	privilegedRole    string
)

type APIResponse struct {
	Success bool `json:"success"`
	Tags    []struct {
		ID    string `json:"id"`
		Name  string `json:"name"`
		Count string `json:"count"`
	} `json:"tags"`
	Screenshots []string `json:"screenshots"`
	Game        struct {
		ID            string      `json:"id"`
		Name          string      `json:"name"`
		URL           string      `json:"url"`
		URLSpdrn      interface{} `json:"url_spdrn"`
		Author        []string    `json:"author"`
		Collab        bool        `json:"collab"`
		Rating        string      `json:"rating"`
		Difficulty    string      `json:"difficulty"`
		RatingCount   string      `json:"rating_count"`
		DateCreated   string      `json:"date_created"`
		CreatorReview string      `json:"creator_review"`
	} `json:"game"`
}

func main() {
	delfruitAPIkey = os.Getenv("DELFRUIT_API_KEY")
	subscribedChannel = os.Getenv("DISCORD_SUBSCRIBED_CHANNEL")
	privilegedRole = os.Getenv("DISCORD_PRIVILEGED_ROLE")
	discordAPIkey = os.Getenv("DISCORD_API_KEY")
	notificationGuild = os.Getenv("DISCORD_NOTIFICATION_GUILD")

	err := *new(error)
	discord, err = discordgo.New("Bot " + discordAPIkey)
	if err != nil {
		fmt.Println(err)
		panic("Invalid Discord API Key")
	}
	discord.AddHandler(status)
	discord.Open()
	<-make(chan struct{})
}

func single(identifier, t string) *discordgo.MessageEmbed {
	f := url.Values{}
	f.Add("method", "single")
	f.Add("api_key", delfruitAPIkey)
	switch t {
	case "": // random
		break
	default:
		f.Add(t, identifier)
	}

	r, err := http.PostForm(delfruitAPIurl, f)
	if err != nil {
		return nil
	}
	defer r.Body.Close()

	var s APIResponse
	dec := json.NewDecoder(r.Body)
	err = dec.Decode(&s)
	if err != nil {
		return nil
	}
	if s.Game.Name == "" {
		return nil
	}

	return format(s)
}

func format(resp APIResponse) *discordgo.MessageEmbed {
	rating, _ := strconv.ParseFloat(resp.Game.Rating, 64)
	rating /= 10
	difficulty, _ := strconv.ParseFloat(resp.Game.Difficulty, 64)
	authorPrefix := "Author: "
	screenshot := "http://i.imgur.com/M1ioBbT.png" // no screenshot available
	tags := "No Available Tags"
	author := strings.Join(resp.Game.Author, " ")
	collab := resp.Game.Collab
	gameInfo := ""
	link := strings.Join([]string{"\n\nDownload Link \n<", resp.Game.URL, ">"}, "")

	if collab {
		authorPrefix = "Authors: "
	}

	if len(resp.Screenshots) > 0 {
		screenshot = strings.Join([]string{"https://delicious-fruit.com", resp.Screenshots[0], ".png"}, "")
	}

	if len(resp.Tags) > 0 {
		var mtags []string
		for i := 0; i < len(resp.Tags); i++ {
			mtags = append(mtags, resp.Tags[i].Name)
		}
		tags = strings.Join(mtags, " ")
	}

	if resp.Game.Rating != "" && resp.Game.Difficulty != "" {
		gameInfo = strings.Join([]string{
			authorPrefix, author,
			fmt.Sprintf("\n\n`Rating: %.1f/10`", rating),
			fmt.Sprintf("\n`Difficulty: %.0f/100`", difficulty), link}, "")
	} else if resp.Game.Rating == "" && resp.Game.Difficulty == "" {
		gameInfo = strings.Join([]string{authorPrefix, author, "\n\n`There are no ratings for this game.`", link}, "")
	} else if resp.Game.Difficulty == "" {
		gameInfo = strings.Join([]string{
			authorPrefix, author,
			fmt.Sprintf("\n\n`Rating: %.1f/10`", rating),
			"\n`Difficulty: N/A`", link}, "")
	} else if resp.Game.Rating == "" {
		gameInfo = strings.Join([]string{
			authorPrefix, author,
			"\n\n`Rating: N/A`",
			fmt.Sprintf("\n`Difficulty: %.0f/100`", difficulty), link}, "")
	} else {
	}

	embed := NewEmbed().
		SetTitle(resp.Game.Name).
		AddField("Game Information", gameInfo).
		SetURL(strings.Join([]string{"https://delicious-fruit.com/ratings/game_details.php?id=", resp.Game.ID}, "")).
		SetThumbnail(screenshot).
		SetFooter(tags).
		SetColor(0x00ff2b).MessageEmbed

	return embed
}

func amend(embed *discordgo.MessageEmbed, status string) *discordgo.MessageEmbed {
	switch status {
	case "update":
		return NewEmbed().
			SetTitle(strings.Join([]string{"[Update] ", embed.Title}, "")).
			AddField(embed.Fields[0].Name, embed.Fields[0].Value).
			SetURL(embed.URL).
			SetThumbnail(embed.Thumbnail.URL).
			SetFooter(embed.Footer.Text).
			SetColor(0xffff00).MessageEmbed

	case "release":
		return NewEmbed().
			SetTitle(strings.Join([]string{"[Release] ", embed.Title}, "")).
			AddField(embed.Fields[0].Name, strings.Replace(embed.Fields[0].Value, "\n\n`There are no reviews for this game.`", " ", 1)).
			SetURL(embed.URL).
			SetColor(0x00ff2b).MessageEmbed

	case "featured":
		return NewEmbed().
			SetTitle(strings.Join([]string{"[Featured] ", embed.Title}, "")).
			AddField(embed.Fields[0].Name, embed.Fields[0].Value).
			SetURL(embed.URL).
			SetThumbnail(embed.Thumbnail.URL).
			SetFooter(embed.Footer.Text).
			SetColor(0xffaa00).MessageEmbed
	}

	// todo: handle default case
	return nil
}

func status(s *discordgo.Session, j *discordgo.MessageCreate) {
	if !j.Author.Bot { // check if not a bot
		member, err := discord.State.Member(j.GuildID, j.Author.ID)
		published := *new([]*discordgo.MessageEmbed)
		if err != nil {
			fmt.Println("member was not found or did not exist")
			return
		}
		for i := range member.Roles { // check if privileged
			if member.Roles[i] == privilegedRole {
				args := strings.Split(j.Content, " ")

				if len(args) > 1 {
					for i := 1; i < len(args); i++ {
						embed := single(args[i], "id")
						if embed == nil {
							s.ChannelMessageSend(j.ChannelID, strings.Join([]string{"No Game with ID ", args[i]}, ""))
							return
						} else {
							published = append(published, embed)
						}

						switch args[0] {
						case "!release":
							embed = amend(embed, "release")
						case "!update":
							embed = amend(embed, "update")
						case "!featured":
							embed = amend(embed, "featured")
						}

						s.ChannelMessageSendEmbed(subscribedChannel, embed)
					}
					break
				} else {
					s.ChannelMessageSend(j.ChannelID, "Parameters were not supplied.")
					return
				}
			} else {
				args := strings.Split(j.Content, " ")

				fmt.Println(strings.Join([]string{j.Author.Username, " did not meet role requirement for command [", args[0], "]"}, ""))

				fmt.Println()
				return
			}
		}

		embeds := []string{"Published to subscribed channels!", "```"}
		for i := range published {
			embeds = append(embeds, published[i].Title, "\n")
		}
		embeds = append(embeds, "```")
		s.ChannelMessageSend(j.ChannelID, strings.Join(embeds, ""))
	}
}
