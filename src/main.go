package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strconv"
	"strings"

	"github.com/bwmarrin/discordgo"
)

const delfruitAPIurl = "https://delicious-fruit.com/api/game.php"

var delfruitAPIkey string

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
