package main

import (
	"flag"
	"log"
	"os"

	"github.com/jbpratt/bots/internal/triviabot"
	"go.uber.org/zap"
)

func main() {
	serverURL := "wss://chat.strims.gg/ws"
	dbPath := flag.String("db", "/tmp/trivia.db", "path to sqlite database")
	dev := flag.Bool("dev", false, "use chat2")

	flag.Parse()

	if *dev {
		serverURL = "wss://chat2.strims.gg/ws"
	}

	logger, err := zap.NewDevelopment()
	if err != nil {
		log.Fatal(err)
	}
	defer func() {
		if err = logger.Sync(); err != nil {
			log.Fatal(err)
		}
	}()

	url, jwt := os.Getenv("STRIMS_CHAT_WSS_URL"), os.Getenv("STRIMS_CHAT_TOKEN")
	if url == "" {
		url = serverURL
	}
	if jwt == "" {
		logger.Fatal("must provide $STRIMS_CHAT_TOKEN")
	}

	triviabot, err := triviabot.New(logger.Sugar(), url, jwt, *dbPath, 3, 15)
	if err != nil {
		logger.Fatal(err.Error())
	}

	if err = triviabot.Run(); err != nil {
		logger.Fatal(err.Error())
	}
}

// accept duration and category and difficulty
// -- can we do this with flag.FlagSet?
// show time difference in between pole positions
