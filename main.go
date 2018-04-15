package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/line/line-bot-sdk-go/linebot"
)

func main() {
	port := os.Getenv("PORT")

	if port == "" {
		log.Fatal("$PORT must be set")
	}

	router := gin.New()
	router.Use(gin.Logger())

	router.POST("/hook", func(c *gin.Context) {
		client := &http.Client{Timeout: time.Duration(15 * time.Second)}

		channelSecret := os.Getenv("LINE_CHANNEL_SECRET")
		channelAccessToken := os.Getenv("LINE_CHANNEL_ACCESSTOKEN")
		bot, err := linebot.New(channelSecret, channelAccessToken, linebot.WithHTTPClient(client))
		if err != nil {
			fmt.Println(err)
			return
		}
		received, err := bot.ParseRequest(c.Request)

		for _, event := range received {
			if event.Type == linebot.EventTypeMessage {
				switch message := event.Message.(type) {
				case *linebot.TextMessage:
					source := event.Source
					if source.Type == linebot.EventSourceTypeRoom {
						postMessage := linebot.NewTextMessage(message.Text)
						if _, err = bot.ReplyMessage(event.ReplyToken, postMessage).Do(); err != nil {
							log.Print(err)
						}
					}
				}
			}
		}
	})

	router.Run(":" + port)
}