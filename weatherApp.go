package weatherApp

import (
	"fmt"
	"io"
	"math/rand"
	"net/http"
	"time"

	"github.com/tidwall/gjson"
)

// Write an interface which allows users to see if it is raining at a particular location,
//   what the temperature is at a particular location and the wind speed (each one a different method).
//   Implement it for one of the two APIs above using net/http and implement it for random weather.
//   E.g. use a random number to pick the temperature.

//go:generate go run github.com/golang/mock/mockgen -destination mocks/Iweather.go github.com/rijuCB/weatherAppGo Iweather
type Iweather interface {
	GetRain() float64 //Rain in mm
	GetTemp() float64 //Temp in C
	GetWind() float64 //Wind speed
}

func CoallesceWeatherInfo(w Iweather) WeatherData {
	return WeatherData{
		Temperature: w.GetTemp(),
		Wind:        w.GetWind(),
		Rain:        w.GetRain(),
	}
}

type WeatherData struct {
	Temperature float64
	Wind        float64
	Rain        float64
}

type Weather struct {
	Zone   string  `json:"timezone"`
	Hourly HourlyW `json:"hourly"`
}

type HourlyW struct {
	// Time      []string  `json:"time"`
	Temperature []float64 `json:"temperature_2m"`
	Wind        []float64 `json:"windspeed_10m"`
	Rain        []float64 `json:"rain"`
}

type Meteo struct{ data string } //Meteo API

func (w *Meteo) GetLocationData(Lat float64, Long float64) {
	request := fmt.Sprintf("https://api.open-meteo.com/v1/forecast?latitude=%v&longitude=%v&hourly=rain,temperature_2m,windspeed_10m", Lat, Long)
	resp, err := http.Get(request)
	if err != nil {
		println(err)
	}
	defer resp.Body.Close()

	a, _ := io.ReadAll(resp.Body)
	// println(string(a))
	w.data = "`" + string(a) + "`"
	return
}

func (w *Meteo) GetRain() float64 {
	return gjson.Get(w.data, "hourly.rain").Array()[0].Float()
}

func (w *Meteo) GetTemp() float64 {
	return gjson.Get(w.data, "hourly.temperature_2m").Array()[0].Float()
}

func (w *Meteo) GetWind() float64 {
	return gjson.Get(w.data, "hourly.windspeed_10m").Array()[0].Float()
}

type RandWeather struct{ RandGen rand.Rand }

func (w *RandWeather) GetRain() float64 {
	return float64(w.RandGen.Intn(100)) / 10
}

func (w *RandWeather) GetTemp() float64 {
	return float64(w.RandGen.Intn(100)) / 10
}

func (w *RandWeather) GetWind() float64 {
	return float64(w.RandGen.Intn(100)) / 10
}

func main() {
	var w Iweather
	//Create weather links
	rw := RandWeather{*rand.New(rand.NewSource(time.Now().UnixNano()))}
	aw := Meteo{}
	aw.GetLocationData(52.52, 13.41)

	w = &aw

	fmt.Printf("Rain is %vmm\n", w.GetRain())
	fmt.Printf("Temp is %vC\n", w.GetTemp())
	fmt.Printf("Wind is %vkmh\n", w.GetWind())

	fmt.Println()

	w = &rw
	fmt.Printf("Rain is %vmm\n", w.GetRain())
	fmt.Printf("Temp is %vC\n", w.GetTemp())
	fmt.Printf("Wind is %vkmh\n", w.GetWind())

	fmt.Println(CoallesceWeatherInfo(w))
}
