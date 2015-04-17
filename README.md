#Slack IRC Bridge

The elementary OS developers use slack extensively for development and it
is pretty neat but it doesn't allow people to driveby and contribute.

This slack-irc-bridge allows us to easily communicate with people over irc
through slack and vice versa.

##Building

To build, use the following commands

```
$ go get
$ go build
```

The bot is not very configurable yet, so you have to have open main.go
and edit the channel name and token manually.

##Author

Akshay Shekher
