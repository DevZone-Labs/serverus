package bot

import (
	"encoding/json"
	"log"

	"github.com/bwmarrin/discordgo"
)

const SERVER_DOWN_MESSAGE = "Sorry! I cannot your request right now!"

func (bot *Bot) registerHandlers() {

	bot.addHandler(helloWorld)
	bot.addHandler(channelInfo)
}

func (bot *Bot) addHandler(handler interface{}) func() {
	switch h := handler.(type) {
	case func(*discordgo.Session, *discordgo.MessageCreate):
		return bot.session.AddHandler(func(session *discordgo.Session, message *discordgo.MessageCreate) {
			if message.Author.ID == session.State.User.ID {
				return
			}

			if !bot.acceptingRequests {
				session.ChannelMessageSend(message.ChannelID, SERVER_DOWN_MESSAGE)
				return
			}

			h(session, message)
		})
	default:
		log.Panic("Invalid Handler Function")
		return nil
	}
}

func helloWorld(session *discordgo.Session, message *discordgo.MessageCreate) {

	if message.Content == "Hello" {
		session.ChannelMessageSend(message.ChannelID, "World!")
	}
}

type ChannelInfo struct {
	ID string `json:"id"`
	Name string `json:"name"`
}

func channelInfo(session *discordgo.Session, message *discordgo.MessageCreate) {

	if message.Content == "Channel Info"{
		channelInfo, err := session.Channel(message.ChannelID)
		if err != nil {
			session.ChannelMessageSend(message.ChannelID, "Could not get channel information")
			return
		}

		jsonData, _:= json.MarshalIndent(ChannelInfo{
			ID: channelInfo.ID,
			Name: channelInfo.Name,
		}, "", "\t")

		session.ChannelMessageSend(message.ChannelID, string(jsonData))
	}
}