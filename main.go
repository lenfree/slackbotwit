package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/abourget/slick"
	_ "github.com/lenfree/awsbot/awswit"
	_ "github.com/lenfree/awsbot/funny"
)

var (
	configFile *string
)

func init() {
	cwd, err := os.Getwd()
	if err != nil {
		fmt.Printf("Err: %s\n", err.Error())
		os.Exit(1)
	}
	_, err = os.Stat(cwd + "/.slick")
	configFile = flag.String("config", cwd+"/.slick", "config file")
}

func main() {
	flag.Parse()

	bot := slick.New(*configFile)

	bot.Run()
}
