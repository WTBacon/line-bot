package main

import (
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"os"
	"time"

	"github.com/line/line-bot-sdk-go/linebot"

	"github.com/gin-gonic/gin"
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
		bot, err := linebot.New("40281719617902503c0891a5bb2ebb4f", "Tmu5PIKyXgIWCb526flAvQ6HcmaSvWfB0R4aMJ6IbawLqVHbmZOFLuliL2rRk7qYOcpHgIaM/WKDZnvj1Sr0TRbS92JTOHw2d90PULl6RY23zjZxe1jyqVWV0HiPhe6smPHX9Xg3V9vhFNBvB2UOmQdB04t89/1O/w1cDnyilFU=", linebot.WithHTTPClient(client))
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
					if resMessage := getResMessage(message.Text); resMessage != "" {
						fmt.Println("resMessage : " + resMessage)
						postMessage := linebot.NewTextMessage(resMessage)
						fmt.Println("postMessage : " + postMessage.Text)
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

func getResMessage(reqMessage string) (message string) {
	resMessages := [3]string{"message 1", "message 2", "message 3"}

	rand.Seed(time.Now().UnixNano())

	// セミコロンで区切ると条件判定の前処理を書ける.
	// 条件式に用いる変数のスコープを限定できる.
	if math := rand.Intn(4); math != 3 {
		message := resMessages[math]
		return message + reqMessage
	}
	return reqMessage
}
