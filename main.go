package main

import (
	"fmt"
	"log"

	"github.com/voldyman/ircbot"
	"github.com/voldyman/slackbot"
)

func main() {
	users := make(map[string]int)
	ircNick := "slackTestBot"
	ircPass := ""
	slackToken := ""

	// slack -> irc mapping
	bridges := map[string]string{
		"#django":  "#botTestChan",
		"#django2": "#botTestChan2",
	}

	slackBot := slackbot.New(slackToken)

	slackEvents, err := slackBot.Start("https://ele.slack.com")
	if err != nil {
		fmt.Println("Could not start slack bot", err.Error())
		return
	}

	ircBot := ircbot.New("irc.freenode.net:6667", ircNick, Values(bridges))
	ircEvents, err := ircBot.Start()
	if err != nil {
		fmt.Println("Could not connect to IRC")
		return
	}

	ircBot.SendMessage("NickServ", "identify", ircPass)
	for {
		select {
		case msg := <-ircEvents:
			log.Printf("IRC: <%s@%s> %s\n", msg.Sender, msg.Channel, msg.Text)

			if target, ok := KeyForValue(bridges, msg.Channel); ok {
				slackBot.SendMessage(msg.Sender, target, msg.Text)
				incUser(users, msg.Sender, target)
			}

		case ev := <-slackEvents:
			switch ev.(type) {

			case *slackbot.MessageEvent:
				msg := ev.(*slackbot.MessageEvent)

				// we don't handle named channels without '#'
				msg.Channel = "#" + msg.Channel

				if _, ok := bridges[msg.Channel]; !ok {
					continue
				}
				log.Printf("slack: <%s@%s> %s\n", msg.Sender,
					msg.Channel, msg.Text)

				if shouldHandle(users, msg.Sender, msg.Channel) {
					log.Println("Handling Message")

					if target, ok := bridges[msg.Channel]; ok {
						ircBot.SendMessage(msg.Sender, msg.Text, target)
					}

				}

				//case error:
				//	err = ev.(error)
				//	fmt.Println("Error occured:", err.Error())

			}
		}
	}

}

// Get All the keys of the map
func Keys(bridges map[string]string) []string {
	vals := []string{}

	for k := range bridges {
		vals = append(vals, k)
	}

	return vals
}

// Get all values of the map
func Values(bridges map[string]string) []string {
	vals := []string{}

	for _, v := range bridges {
		vals = append(vals, v)
	}

	return vals
}

// Get the key of a map for the given value
func KeyForValue(bridges map[string]string, val string) (string, bool) {
	result := ""

	for k, v := range bridges {
		if v == val {
			result = k
		}
	}

	if result == "" {
		return "", false
	}
	return result, true
}

// Semaphores to manages messages

func incUser(users map[string]int, user, channel string) {
	key := user + channel
	if val, ok := users[key]; ok {
		users[key] = val + 1
	} else {
		users[key] = 1
	}
}

func shouldHandle(users map[string]int, user, channel string) bool {
	key := user + channel
	if val, ok := users[key]; ok {
		if val > 0 {
			users[key] = val - 1
			return false
		}
	}

	return true
}
