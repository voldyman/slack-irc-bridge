package main

import (
	"bytes"
	"fmt"

	"github.com/sorcix/irc"
	"github.com/voldyman/ircx"
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

	ircEvents, ircBot, err := startIRCBot("irc.freenode.net:6667", "slackeraBot", []string{ircChannel})
	if err != nil {
		fmt.Println("Error connecting to irc")
		return
	}

	for {
		select {
		case ev := <-ircEvents:
			switch ev.(type) {
			case *ircMessageEvent:
				msg := ev.(*ircMessageEvent)
				fmt.Printf("irc: <%s> %s\n", msg.Sender, msg.Text)
				slackBot.SendMessage(msg.Sender, slackChannel, msg.Text)
				incUser(users, msg.Sender)
			}

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

					msgBuf := bytes.NewBufferString("")
					fmt.Fprintf(msgBuf, "<%s>: %s", msg.Sender, msg.Text)

					ircBot.SendMessage(&irc.Message{
						Command:  "PRIVMSG",
						Params:   []string{ircChannel},
						Trailing: msgBuf.String(),
					})
					incUser(users, "voldy")
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

type (
	ircEvent interface{}

	ircMessageEvent struct {
		Sender string
		Text   string
	}
)

func startIRCBot(server, name string, channels []string) (chan ircEvent, *ircx.Bot, error) {
	bot := ircx.Classic(server, name)
	if err := bot.Connect(); err != nil {
		return nil, nil, err
	}

	events := make(chan ircEvent)

	bot.AddCallback(irc.PING, ircx.Callback{Handler: ircx.HandlerFunc(pingHandler)})
	bot.AddCallback(irc.RPL_WELCOME, ircx.Callback{Handler: ircx.HandlerFunc(registerConnect(channels))})
	bot.AddCallback(irc.PRIVMSG, ircx.Callback{Handler: ircx.HandlerFunc(createMsgHandler(events))})

	go bot.CallbackLoop()

	return events, bot, nil
}

func createMsgHandler(events chan ircEvent) func(ircx.Sender, *irc.Message) {
	return func(s ircx.Sender, m *irc.Message) {
		ev := &ircMessageEvent{
			Sender: m.Name,
			Text:   m.Trailing,
		}
		events <- ev
	}
}
func registerConnect(channels []string) func(ircx.Sender, *irc.Message) {
	return func(s ircx.Sender, m *irc.Message) {
		s.Send(&irc.Message{
			Command: irc.JOIN,
			Params:  channels,
		})
	}
}

func pingHandler(s ircx.Sender, m *irc.Message) {
	s.Send(&irc.Message{
		Command:  irc.PONG,
		Params:   m.Params,
		Trailing: m.Trailing,
	})
}
