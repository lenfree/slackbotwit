package awswit

import (
	"fmt"
	"regexp"
	"strings"
	"time"

	"github.com/abourget/slick"
)

// AwsWit contains slick related configuration
type AwsWit struct {
	bot    *slick.Bot
	config WitConfig
}

var (
	witToken string
)

//WitConfig contains awswit secrets
type WitConfig struct {
	WitToken  string `json:"wit_token"`
	AWSRegion string `json:"aws_region"`
}

func init() {
	slick.RegisterPlugin(&AwsWit{})
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
				fmt.Println(awsIntent)
				if strings.Contains(awsIntent.EntityName.Entity.Value, "ec2") == true {

					if strings.Contains(awsIntent.Intent.Entity.Value, "how many") == true ||
						strings.Contains(awsIntent.Intent.Entity.Value, "number") == true {
						go awsEC2CountHandler(listen, msg, awsIntent)
					} else if strings.Contains(awsIntent.Property.Entity.Value, "tags") == true {
						go awsEC2Tags(listen, msg, awsIntent)
					} else if strings.Contains(awsIntent.Property.Entity.Value, "ip address") == true {
						go awsEC2Handler(listen, msg, awsIntent)
					}
				}

				if strings.Contains(awsIntent.EntityName.Entity.Value, "load balancer") == true ||
					strings.Contains(awsIntent.EntityName.Entity.Value, "elb") == true {

					listNumberIntention, _ := regexp.MatchString("number.*|how many", msg.Text)
					if listNumberIntention == true {
						go awsELBCountHandler(listen, msg, awsIntent)
					}

					if len(awsIntent.EntityType.Entity.Value) <= 0 && listNumberIntention == false {
						go awsELBHandler(listen, msg, awsIntent)
					}

					if len(awsIntent.EntityType.Entity.Value) > 0 && listNumberIntention == false {
						go awsELBNameHandler(listen, msg, awsIntent)
					}
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
	if len(ec2List.Reservations) > 0 {
		for _, r := range ec2List.Reservations {
			for c, i := range r.Instances {
				for _, n := range i.NetworkInterfaces {
					msg.ReplyMention("ip address for %s instance # %d: ```%v```",
						e.EntityType.Entity.Value,
						c+1,
						*n.PrivateIpAddress)
					continue
				}
			}
		}
	} else {
		msg.ReplyMention("no instance with tag:Name ```%s``` found", e.EntityType.Entity.Value)
	}
	fmt.Printf("%s\n", ec2List.GoString())
	listen.Close()
}

func awsEC2CountHandler(listen *slick.Listener, msg *slick.Message, e awsEntities) {
	filter := awsEC2Filter{Tag: "", Value: ""}

	ec2List := getEC2List(newEC2(), filter)

	running := 0
	pending := 0
	for _, r := range ec2List.Reservations {
		for _, i := range r.Instances {
			switch *i.State.Name {
			case "running":
				running++
			case "pending":
				pending++
			}
		}
	}
	fmt.Printf("total: %s\n", ec2List.Reservations)
	msg.Reply("number of EC2 instances in running state is: ```%d```", running)
	msg.Reply("number of EC2 instances in pending state is: ```%d```", pending)
	listen.Close()
}

func awsELBHandler(listen *slick.Listener, msg *slick.Message, e awsEntities) {
	elbList := getELBList(newELB())
	for _, elb := range elbList.LoadBalancerDescriptions {
		msg.ReplyMention("ELB %s :\n ```%s``` ", *elb.LoadBalancerName, elb)
		fmt.Printf("%s\n", elbList.GoString())
	}
	listen.Close()
}

func awsELBCountHandler(listen *slick.Listener, msg *slick.Message, e awsEntities) {
	elbList := getELBList(newELB())
	msg.ReplyMention("# of ELBs in Sydney region ```%d```", len(elbList.LoadBalancerDescriptions))
	listen.Close()
}

func awsELBNameHandler(listen *slick.Listener, msg *slick.Message, e awsEntities) {
	elb := getELBName(newELB(), e.EntityType.Entity.Value)
	if elb != nil {
		msg.ReplyMention("ELB ```%s`", elb.GoString())
	} else {
		msg.ReplyMention("ELB name %s not found", e.EntityType.Entity.Value)
	}
	listen.Close()
}

func awsEC2Tags(listen *slick.Listener, msg *slick.Message, e awsEntities) {
	filter := awsEC2Filter{Tag: "", Value: ""}

	var tags []string
	ec2List := getEC2List(newEC2(), filter)
	if len(ec2List.Reservations) > 0 {
		for _, r := range ec2List.Reservations {
			for _, i := range r.Instances {
				for _, tag := range i.Tags {
					switch *tag.Key {
					case "Env":
						continue
					case "SaltEnv":
						continue
					default:
						tags = append(tags, *tag.Key+" - "+*tag.Value)
					}
				}
			}
		}
	}
	for _, tag := range tags {
		msg.ReplyMention("Key Value pair tag")
		msg.ReplyMention("```%s```", tag)
		time.Sleep(3 * time.Second)
	}
	listen.Close()
}
