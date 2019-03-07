package WeatherApiClient

import (
	"fmt"
	"os"
	"testing"

	"github.com/joho/godotenv"
)

func EnvLoad() {
	if os.Getenv("ENV") == "" {
		godotenv.Load(".env")
	}
}

func TestWeatherApiClient_Fetch(t *testing.T) {
	EnvLoad()
	wac := NewWeatherApiClient(os.Getenv("WEATHER_API"))
	wd := wac.Fetch()

	fmt.Println(wd)

	if wd.Cod.String() != "200" {
		t.Errorf("got: %v\nwant: %v", wd.Cod, "200")
	}
	for _, values := range wd.List {
		if &values.Rain.Volume == nil {
			t.Errorf("got: %v\nwant: %v", values.Rain.Volume, "Number")
		}
	}
}
