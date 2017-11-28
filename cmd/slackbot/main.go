package main

import (
	"fmt"
	"log"
	"os"

	"github.com/nlopes/slack"
)

var (
	botId   string
	botName string
)

func run(api *slack.Client) int {
	rtm := api.NewRTM()
	go rtm.ManageConnection()

	for msg := range rtm.IncomingEvents {
		switch ev := msg.Data.(type) {
		case *slack.HelloEvent:
			log.Print("Hello Event")
		case *slack.ConnectedEvent:
			bot := ""
			for _, channel := range ev.Info.Channels {
				if channel.Name == "bot" {
					bot = channel.ID
				}
			}
			log.Printf("URL: %s\n", ev.Info.URL)
			rtm.SendMessage(rtm.NewOutgoingMessage("connection bot", bot))
		case *slack.MessageEvent:
			// TODO : 自分自身をGlobalな変数で持つ
			// TODO : 自分にメンションが来た時だけ反応する
			// TODO : 自分が呼ばれたチャンネルを持つ
			log.Printf("Channel: %s, Message: %s\n", ev.Msg.Channel, ev.Msg.Text)
			mess := ev.Msg.Text
			rtm.SendMessage(rtm.NewOutgoingMessage(fmt.Sprintf("Message is \"%s\"", mess), ev.Msg.Channel))
		case *slack.PresenceChangeEvent:
			log.Printf("PresenceChangeEvent: %v\n", ev)
		case *slack.RTMError:
			log.Printf("Error: %s\n", ev.Error())
		case *slack.InvalidAuthEvent:
			log.Print("Invalid credentials")
			return 1
		default:
		}
	}
	return 1
}

func main() {
	token := os.Getenv("BOTAPITOKEN")
	if token == "" {
		log.Fatal("missing API token")
	}
	api := slack.New(token)
	//os.Exit(run(api))
	run(api)
}
