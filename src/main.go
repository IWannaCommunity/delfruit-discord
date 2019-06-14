package main

import (
	"fmt"
	"strconv"
	"strings"
	"github.com/bwmarrin/discordgo"
)

const delfruitAPIurl = "https://delicious-fruit.com/api/game.php"

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
func game(resp APIResponse) *discordgo.MessageEmbed {
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
		authorPrefix := "Authors: "
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
