package main

import (
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"os"
	"regexp"
	"time"

	"github.com/line/line-bot-sdk-go/linebot"

	"github.com/gin-gonic/gin"
)

func main() {
	port := os.Getenv("PORT")
	channelSecret := os.Getenv("CHANNEL_SECRET")
	channelToken := os.Getenv("CHANNEL_TOKEN")

	if port == "" {
		log.Fatal("$PORT must be set")
	}

	router := gin.New()
	router.Use(gin.Logger())

	router.POST("/hook", func(c *gin.Context) {
		client := &http.Client{Timeout: time.Duration(15 * time.Second)}
		bot, err := linebot.New(channelSecret, channelToken, linebot.WithHTTPClient(client))

		if err != nil {
			fmt.Println(err)
			return
		}
		received, err := bot.ParseRequest(c.Request)

		for _, event := range received {
			if event.Type == linebot.EventTypeMessage {
				// 変数 message のスコープを switch 文の内部に絞っている.
				switch message := event.Message.(type) {
				case *linebot.TextMessage:
					if convMatchingMessage(message.Text) {
						// セミコロンで区切ると条件判定の前処理を書ける.
						// 条件式に用いる変数のスコープを限定できる.
						if resMessage := convResMessage(); resMessage != "" {
							postMessage := linebot.NewTextMessage(resMessage)
							if _, err = bot.ReplyMessage(event.ReplyToken, postMessage).Do(); err != nil {
								log.Print(err)
							}
						}
					}
				}
			}
		}

	})

	router.Run(":" + port)
}

func convMatchingMessage(targetMessage string) (matched bool) {
	r := regexp.MustCompile(`涼真`)
	return r.MatchString(targetMessage)
}

func convResMessage() (message string) {
	resMessages := [3]string{"お仕事お疲れ様！\n今日も頑張って偉いね。\nたまにはゆっくり休むんだよ^^", "あんまり無理するなよ！\n心配になっちゃうからさ^^", "ちゃんとご飯食べた？\n食べないと元気でないぞ！"}
	rand.Seed(time.Now().UnixNano())
	resMessage := resMessages[rand.Intn(3)]
	return resMessage
}
