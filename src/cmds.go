package main

import (
	"fmt"
	"strings"

	"github.com/bwmarrin/discordgo"
)

func publish(s *discordgo.Session, j *discordgo.MessageCreate) {
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
