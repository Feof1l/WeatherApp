package main

import (
	"encoding/json"
	"io/ioutil"
	"math"
	"net/http"
	"strings"
)

type apiConfigData struct {
	OpenWeatherMapApiKey string `json:"OpenWeatherMapApiKey"`
}
type weatherData struct {
	Name string `json:"name"`
	Main struct {
		Celsius  float64 `json:"temp"`
		Humidity int     `json:"humidity"`
	} `json:"main"`
}

func loadApiConfig(fileName string) (apiConfigData, error) {
	bytes, err := ioutil.ReadFile(fileName)
	if err != nil {
		return apiConfigData{}, err
	}
	var c apiConfigData
	err = json.Unmarshal(bytes, &c)
	if err != nil {
		return apiConfigData{}, err
	}
	return c, nil
}

func roundFloat(val float64, precision uint) float64 {
	ratio := math.Pow(10, float64(precision))
	return math.Round(val*ratio) / ratio
}

func hello(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Hello from Ivan!\n"))
}

func query(city string) (weatherData, error) {
	apiConfig, err := loadApiConfig(".apiConfig")
	if err != nil {
		return weatherData{}, err
	}
	resp, err := http.Get("http://api.openweathermap.org/data/2.5/weather?APPID=" + apiConfig.OpenWeatherMapApiKey + "&q=" + city)
	if err != nil {
		return weatherData{}, err
	}
	defer resp.Body.Close()
	var d weatherData
	if err := json.NewDecoder(resp.Body).Decode(&d); err != nil {
		return weatherData{}, err
	}
	d.Main.Celsius = roundFloat((d.Main.Celsius - 273.15), 2)
	return d, nil
}

func main() {
	http.HandleFunc("/hello", hello)
	http.HandleFunc("/weather/", func(w http.ResponseWriter, r *http.Request) {
		city := strings.SplitN(r.URL.Path, "/", 3)[2]
		data, err := query(city)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		json.NewEncoder(w).Encode(data)
	})
	http.Handle("/", http.FileServer(http.Dir("./static")))
	http.ListenAndServe(":8080", nil)
}
