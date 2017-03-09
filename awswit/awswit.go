package awswit

import (
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/abourget/slick"
	"github.com/joho/godotenv"
)

type AwsWit struct {
	bot    *slick.Bot
	config WitConfig
}

var (
	witToken string
)

//WitConfig contains awswit secrets
type WitConfig struct {
	WitToken string `json:"wit_token"`
}

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

	var conf struct {
		AwsWit WitConfig
	}
	bot.LoadConfig(&conf)
	awswit.config = conf.AwsWit
}

// ChatHandler handles direct messages bot or mentions bot and route
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
				w := newWit(awswit.config.WitToken)
				awsIntent := parse(query(w, msg.Text))

				if strings.Contains(awsIntent.EntityName.Entity.Value, "ec2") == true {
					awsEC2Handler(listen, msg, awsIntent)
				}

				if strings.Contains(awsIntent.Intent.Entity.Value, "how many") == true {
					awsEC2CountHandler(listen, msg, awsIntent)
				}

				//if intent.Entity.Value == "pizza" && entity.Entity.Value != "" {
				//	//msg.Reply(fmt.Sprintf("%s", response.Outcomes.Type[0].Value))
				//	msg.ReplyMention("you want pizza flavor %s", entity.Entity.Value)
				//	listen.Close()
				//} else if intent.Entity.Value == "pizza" && entity.Entity.Value == "" {
				//	msg.ReplyMention("what flavor of pizza would you like?")
				//	bot.Listen(&slick.Listener{
				//		ListenDuration: 5 * time.Second,
				//		FromUser:       msg.FromUser,
				//		FromChannel:    msg.FromChannel,
				//		MentionsMeOnly: true,
				//		MessageHandlerFunc: func(listen *slick.Listener, msg *slick.Message) {
				//			msg.ReplyMention("I think you want pizza flavor %s", msg.Text)
				//			listen.Close()
				//		},
				//		TimeoutFunc: func(listen *slick.Listener) {
				//			listen.Close()
				//		},
				//	})
				//}
			},
		})
	}
}

func awsEC2Handler(listen *slick.Listener, msg *slick.Message, e awsEntities) {
	filter := awsEC2Filter{
		Tag:   "tag:Name",
		Value: e.EntityType.Entity.Value,
	}

	ec2List := getEC2List(newEC2(), filter)

	for _, r := range ec2List.Reservations {
		for c, i := range r.Instances {
			for _, n := range i.NetworkInterfaces {
				msg.ReplyMention("ip address for %s instance # %d: %v\n",
					e.EntityType.Entity.Value,
					c+1,
					*n.PrivateIpAddress)
			}
		}
	}
	fmt.Printf("%s\n", ec2List.GoString())
	listen.Close()
}

func awsEC2CountHandler(listen *slick.Listener, msg *slick.Message, e awsEntities) {
	filter := awsEC2Filter{Tag: "", Value: ""}

	ec2List := getEC2List(newEC2(), filter)

	count := 0
	for _, r := range ec2List.Reservations {
		for _ = range r.Instances {
			count++
		}
	}
	msg.Reply("number of EC2 instances either in running or pending state is: %d", count)
	listen.Close()
}
