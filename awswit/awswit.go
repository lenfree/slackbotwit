package awswit

import (
	"fmt"
	"os"
	"time"

	"github.com/abourget/slick"
	"github.com/joho/godotenv"
)

type AwsWit struct{}

var (
	witToken string
)

func init() {
	slick.RegisterPlugin(&AwsWit{})
	err := godotenv.Load()
	if err != nil {
		fmt.Printf("Err: %s\n", err.Error())
		os.Exit(1)
	}
	witToken = os.Getenv("WIT_TOKEN")
}

// InitPlugin registers this plugin
func (awswit *AwsWit) InitPlugin(bot *slick.Bot) {
	bot.Listen(&slick.Listener{
		MessageHandlerFunc: awswit.ChatHandler,
	})
}

// ChatHandler handles direct messages bot or mentions bot
func (awswit *AwsWit) ChatHandler(listen *slick.Listener, msg *slick.Message) {
	bot := listen.Bot

	if msg.MentionsMe {
		bot.Listen(&slick.Listener{
			ListenDuration: time.Duration(10 * time.Second),
			MessageHandlerFunc: func(listen *slick.Listener, msg *slick.Message) {
				//if strings.Contains(msg.Text, "papa") {
				//	msg.Reply("3s", "yo rocker").DeleteAfter("3s")
				//	msg.AddReaction("wink")
				//	go func() {
				//		time.Sleep(3 * time.Second)
				//		msg.AddReaction("beer")
				//		time.Sleep(1 * time.Second)
				//		msg.RemoveReaction("wink")
				//	}()
				//}
				w := newWit(witToken)
				intent, entity := parse(query(w, msg.Text))

				fmt.Printf("full response %+#v == %+#v\n", intent.Name, entity.Name)
				if intent.Entity.Value == "pizza" && entity.Entity.Value != "" {
					//msg.Reply(fmt.Sprintf("%s", response.Outcomes.Type[0].Value))
					msg.Reply("you want pizza flavor %s", entity.Entity.Value)
				} else {
					msg.Reply("what flavor of pizza would you like?")
					bot.Listen(&slick.Listener{
						ListenDuration: 5 * time.Second,
						FromUser:       msg.FromUser,
						FromChannel:    msg.FromChannel,
						MentionsMeOnly: true,
						MessageHandlerFunc: func(listen *slick.Listener, msg *slick.Message) {
							msg.Reply("you want pizza flavor %s", msg.Text)
							listen.Close()
						},
					})
				}
			},
		})
	}
}
