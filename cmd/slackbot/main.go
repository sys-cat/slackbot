package main

import (
	"log"
	"os"
	"strings"

	"github.com/nlopes/slack"
)

var (
	botId   string
	botName string
)

func run(api *slack.Client) int {
	rtm := api.NewRTM()
	go rtm.ManageConnection()

	for {
		select {
		case msg := <-rtm.IncomingEvents:
			switch ev := msg.Data.(type) {
			case *slack.HelloEvent:
				log.Print("Hello Event")
			case *slack.ConnectedEvent:
				botId = ev.Info.User.ID
				botName = ev.Info.User.Name
			case *slack.MessageEvent:
				log.Print("Message: %v\n", ev)
				user := ev.User
				text := ev.Text
				channel := ev.Channel
				if ev.Type == "message" && strings.HasPrefix(text, "<@"+botId+">") {
					rtm.handleResponse(user, text, channel)
				}
				//rtm.SendMessage(rtm.NewOutgoingMessage("Hello world", ev.Channel))
			case *slack.InvalidAuthEvent:
				log.Print("Invalid credentials")
				return 1
			}
		}
	}
}

func main() {
	token := os.Getenv("BOTAPITOKEN")
	if token == "" {
		log.Fatal("missing API token")
	}
	api := slack.New(token)
	os.Exit(run(api))
}
