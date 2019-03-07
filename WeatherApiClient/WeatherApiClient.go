package WeatherApiClient

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"

	"github.com/WTBacon/line-bot/WeatherData"
)

type WeatherApiClient struct {
	RequestUrl string
	ApiKey     string
	Place      string
	Format     string
	Count      string
	Lang       string
}

func NewWeatherApiClient(apikey string) ApiClient {
	return &WeatherApiClient{
		"https://api.openweathermap.org/data/2.5/forecast",
		apikey,
		"Tokyo,jp",
		"json",
		"8",
		"ja",
	}
}

func (w *WeatherApiClient) Fetch() WeatherData.WeatherData {
	values := url.Values{}
	values.Add("APPID", w.ApiKey)
	values.Add("q", w.Place)
	values.Add("mode", w.Format)
	values.Add("cnt", w.Count)
	values.Add("lang", w.Lang)
	res, err := http.Get(w.RequestUrl + "?" + values.Encode())
	if err != nil {
		log.Fatalf("Failed to get request url. err: %v", err)
	}

	defer res.Body.Close()
	body, err := ioutil.ReadAll(res.Body)
	var wd WeatherData.WeatherData
	err = json.Unmarshal(body, &wd)
	if err != nil {
		log.Fatalf("Failed to fetched unmarshal json. err: %v", err)
	}

	return wd
}
