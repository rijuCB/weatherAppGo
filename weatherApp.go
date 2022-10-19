package weatherAppGo

import (
	"encoding/json"
	"fmt"
	"io"
	"math/rand"
	"net/http"
	"time"
)

// Write an interface which allows users to see if it is raining at a particular location,
//   what the temperature is at a particular location and the wind speed (each one a different method).
//   Implement it for one of the two APIs above using net/http and implement it for random weather.
//   E.g. use a random number to pick the temperature.

type Iweather interface {
	GetRain() float64 //Rain in mm
	GetTemp() float64 //Temp in C
	GetWind() float64 //Wind speed
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

type Meteo struct{ endpoint string } //Meteo API

func (w *Meteo) GetHourlyData(data string) HourlyW {
	request := w.endpoint + "&hourly=" + data
	resp, err := http.Get(request)
	if err != nil {
		println(err)
	}
	defer resp.Body.Close()

	a, _ := io.ReadAll(resp.Body)
	var str Weather
	err = json.Unmarshal(a, &str)
	if err != nil {
		fmt.Println(err)
	}

	return str.Hourly
}

func (w *Meteo) GetRain() float64 {
	return w.GetHourlyData("rain").Rain[0]
}

func (w *Meteo) GetTemp() float64 {
	return w.GetHourlyData("temperature_2m").Temperature[0]
}

func (w *Meteo) GetWind() float64 {
	return w.GetHourlyData("windspeed_10m").Wind[0]
}

type RandWeather struct{ rGen rand.Rand }

func (w *RandWeather) GetRain() float64 {
	return float64(w.rGen.Intn(100)) / 10
}

func (w *RandWeather) GetTemp() float64 {
	return float64(w.rGen.Intn(100)) / 10
}

func (w *RandWeather) GetWind() float64 {
	return float64(w.rGen.Intn(100)) / 10
}

func main() {
	var w Iweather
	//Create weather links
	rw := RandWeather{*rand.New(rand.NewSource(time.Now().UnixNano()))}
	aw := Meteo{"https://api.open-meteo.com/v1/forecast?latitude=52.52&longitude=13.41"}

	w = &aw

	fmt.Printf("Rain is %vmm\n", w.GetRain())
	fmt.Printf("Temp is %vC\n", w.GetTemp())
	fmt.Printf("Wind is %vkmh\n", w.GetWind())

	fmt.Println()

	w = &rw
	fmt.Printf("Rain is %vmm\n", w.GetRain())
	fmt.Printf("Temp is %vC\n", w.GetTemp())
	fmt.Printf("Wind is %vkmh\n", w.GetWind())

}
