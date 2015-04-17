package main

import (
	"fmt"

	"github.com/voldyman/ircbot"
	"github.com/voldyman/slackbot"
)

func main() {
	users := make(map[string]int)
	ircChannel := "#botTestChan"
	slackChannel := "#django"
	slackToken := ""

	// slack token should be here
	slackBot := slackbot.New(slackToken)

	slackEvents, err := slackBot.Start("https://ele.slack.com")
	if err != nil {
		fmt.Println("Could not start slack bot", err.Error())
		return
	}

	ircBot := ircbot.New("irc.freenode.net:6667", "TestslackerBot", []string{ircChannel})
	ircEvents, err := ircBot.Start()
	if err != nil {
		fmt.Println("Could not connect to IRC")
		return
	}

	for {
		select {
		case msg := <-ircEvents:
			fmt.Printf("irc: <%s> %s\n", msg.Sender, msg.Text)
			slackBot.SendMessage(msg.Sender, slackChannel, msg.Text)
			incUser(users, msg.Sender)

		case ev := <-slackEvents:
			switch ev.(type) {

			case *slackbot.MessageEvent:
				msg := ev.(*slackbot.MessageEvent)

				if "#"+msg.Channel != slackChannel {
					continue
				}
				fmt.Printf("Got Message\n<%s@%s>: %s\n", msg.Sender,
					msg.Channel, msg.Text)

				if shouldHandle(users, msg.Sender) {
					fmt.Println("Handling Message")
					ircBot.SendMessage(msg.Sender, msg.Text)
					incUser(users, msg.Sender)
				}

				//case error:
				//	err = ev.(error)
				//	fmt.Println("Error occured:", err.Error())

			}
		}
	}

}

func incUser(users map[string]int, user string) {
	if val, ok := users[user]; ok {
		users[user] = val + 1
	} else {
		users[user] = 1
	}
}

func shouldHandle(users map[string]int, user string) bool {
	if val, ok := users[user]; ok {
		if val > 0 {
			users[user] = val - 1
			return false
		}
	}

	return true
}
