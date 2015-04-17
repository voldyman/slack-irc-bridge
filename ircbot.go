package main

import (
	"bytes"
	"fmt"

	"github.com/sorcix/irc"
	"github.com/voldyman/ircx"
)

type (
	IRCBot struct {
		channels []string
		bot      *ircx.Bot
		chEvents chan IRCMessageEvent
	}

	IRCMessageEvent struct {
		Sender string
		Text   string
	}
)

func newIRCBot(server, nick string, channels []string) *IRCBot {
	return &IRCBot{
		channels: channels,
		bot:      ircx.Classic(server, nick),
		chEvents: make(chan IRCMessageEvent),
	}
}

func (i *IRCBot) Start() (chan IRCMessageEvent, error) {
	err := i.bot.Connect()
	if err != nil {
		return nil, err
	}

	i.registerHandlers()

	go i.bot.CallbackLoop()

	return i.chEvents, nil
}

func (i *IRCBot) SendMessage(nick, msg string) {
	msgBuf := bytes.NewBufferString("")

	fmt.Fprintf(msgBuf, "<%s> %s", nick, msg)

	i.bot.SendMessage(&irc.Message{
		Command:  "PRIVMSG",
		Params:   i.channels,
		Trailing: msgBuf.String(),
	})
}

func (i *IRCBot) registerHandlers() {
	// IRC Ping Pong handler
	i.bot.AddCallback(irc.PING, ircx.Callback{
		Handler: ircx.HandlerFunc(pingHandler),
	})

	// IRC register handler
	i.bot.AddCallback(irc.RPL_WELCOME, ircx.Callback{
		Handler: ircx.HandlerFunc(i.registerConnect),
	})

}

func (i *IRCBot) msgHandler(s ircx.Sender, m *irc.Message) {
	ev := IRCMessageEvent{
		Sender: m.Name,
		Text:   m.Trailing,
	}
	i.chEvents <- ev
}

func (i *IRCBot) registerConnect(s ircx.Sender, m *irc.Message) {
	fmt.Println("Joining channels")
	s.Send(&irc.Message{
		Command: irc.JOIN,
		Params:  i.channels,
	})

}

func pingHandler(s ircx.Sender, m *irc.Message) {
	s.Send(&irc.Message{
		Command:  irc.PONG,
		Params:   m.Params,
		Trailing: m.Trailing,
	})
}
