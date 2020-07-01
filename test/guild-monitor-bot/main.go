package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	discord "github.com/bwmarrin/discordgo"
	flag "github.com/spf13/pflag"
)

const generalChannelID = "704496025317933099"

var token = flag.StringP("token", "t", "", "Discord Bot token")

func main() {
	flag.Parse()

	// check if user supplied token, if no then abort
	if *token == "" {
		log.Fatalln("Discord Bot token can't be empty")
	}

	// create new Discord client
	session, err := discord.New("Bot " + *token)
	if err != nil {
		log.Fatalln("Creating new discord session: " + err.Error())
	}

	// changed intents, so we could track user's presence
	session.Identify.GuildSubscriptions = false
	session.Identify.Intents = discord.MakeIntent(discord.IntentsGuildPresences)

	// set a callback function to listener
	session.AddHandler(sendMsgByUserStatus)

	// start to listen on a websocket
	err = session.Open()
	if err != nil {
		log.Fatalln("Openning websocket connection: " + err.Error())
	}

	// Wait here until CTRL-C or other term signal is received.
	fmt.Println("Bot is now running.  Press CTRL-C to exit.")
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, os.Kill)
	<-sc

	// Cleanly close down the Discord session.
	session.Close()
}

// sendMsgByUserStatus sends message to general channel, based on user's status.
func sendMsgByUserStatus(s *discord.Session, m *discord.PresenceUpdate) {
	var (
		err error

		msgOnline  string
		msgOffline string

		user        *discord.User
		userMention string
	)

	// if we can't get user's nickname from Presence Gateway, then try to retrieve it using his/her id
	if m.User.Username == "" {
		user, err = s.User(m.User.ID)
		if err != nil { // if there is an error, then mention user, using his id
			log.Printf("Couldn't get user %s info: %v\n", m.User.ID, err)
			userMention = m.User.ID
		} else {
			userMention = user.Mention()
		}
	} else {
		userMention = user.Mention()
	}

	// send message to a channel based on user's status
	msgOnline = fmt.Sprintf("Привет %s, я вижу ты онлайн, не хочешь зайти в доту ?\n", userMention)
	msgOffline = fmt.Sprintf("Пока %s, не забудь посмотреть видео о доте когда будешь в туалете :)\n", userMention)

	if m.Status == discord.StatusOnline {
		s.ChannelMessageSend(generalChannelID, msgOnline)
	} else if m.Status == discord.StatusOffline {
		s.ChannelMessageSend(generalChannelID, msgOffline)
	}
}
