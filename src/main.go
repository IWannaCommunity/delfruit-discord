package main

import (
	"fmt"
	"os"

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
