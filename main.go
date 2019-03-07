package main

import (
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"os"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/WTBacon/line-bot/WeatherApiClient"
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
		weatherApiKey := os.Getenv("WEATHER_API")
		wac := WeatherApiClient.NewWeatherApiClient(weatherApiKey)
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
					} else if weatherMatchingMessage(message.Text) {
						weather := Weather{wac}
						if resMessage := weather.weatherResMessage(); resMessage != "" {
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

func weatherMatchingMessage(targetMessage string) (matched bool) {
	r := regexp.MustCompile(`.*天気.*`)
	return r.MatchString(targetMessage)
}

type Weather struct {
	ApiClient WeatherApiClient.ApiClient
}

func (w *Weather) weatherResMessage() (message string) {
	wd := w.ApiClient.Fetch()
	rainJudge := 0.0
	forecastList := []string{}
	for _, values := range wd.List {
		timeDate := time.Unix(values.Date, 0).In(time.FixedZone("Asia/Tokyo", 9*60*60))
		date := strconv.Itoa(int(timeDate.Month())) + "/" + strconv.Itoa(timeDate.Day()) + " " + strconv.Itoa(timeDate.Hour()) + "時"
		weatherDescription := "\t" + values.Weather[0].Description
		rainVolume := "\n降水量:" + strconv.FormatFloat(values.Rain.Volume, 'f', 0, 64) + "mm"
		forecastList = append(forecastList, date+weatherDescription+rainVolume)
		rainJudge += values.Rain.Volume
	}

	if rainJudge > 0 {
		return "今日の天気は\n" + strings.Join(forecastList, "\n") + "\nだよ！\n" + "雨が降りそうだから傘忘れないでね！\n気をつけて帰ってくるんだよ^^"
	}

	return "今日は傘いらないみたい！\n寄り道しないで帰ってきてね^^"
}
