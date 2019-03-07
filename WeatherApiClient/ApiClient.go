package WeatherApiClient

import "github.com/WTBacon/line-bot/WeatherData"

type ApiClient interface {
	Fetch() WeatherData.WeatherData
}
