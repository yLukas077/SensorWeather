package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"

	"github.com/joho/godotenv"
)

type GeoResponse struct {
	Lat float64 `json:"lat"`
	Lon float64 `json:"lon"`
}

type WeatherData struct {
	Main struct {
		Temp      float64 `json:"temp"`
		FeelsLike float64 `json:"feels_like"`
		Humidity  int     `json:"humidity"`
	} `json:"main"`
}

func kelvinToCelsius(k float64) float64 {
	return k - 273
}

func main() {
	var zipcode, countrycode string
	fmt.Println("Enter your postal code and country code:")
	fmt.Scan(&zipcode, &countrycode)

	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	key_weather := os.Getenv("WEATHER_KEY")
	if key_weather == "" {
		log.Fatal("WEATHER_KEY environment variable is not set")
	}

	geoURL := fmt.Sprintf("http://api.openweathermap.org/geo/1.0/zip?zip=%s,%s&appid=%s", zipcode, countrycode, key_weather)
	response, err := http.Get(geoURL)
	if err != nil {
		log.Fatal(err)
	}
	defer response.Body.Close()

	responseBody, err := io.ReadAll(response.Body)
	if err != nil {
		log.Fatal(err)
	}

	var geoData GeoResponse
	if err := json.Unmarshal(responseBody, &geoData); err != nil {
		log.Fatal(err)
	}

	weatherURL := fmt.Sprintf("https://api.openweathermap.org/data/2.5/weather?lat=%f&lon=%f&appid=%s", geoData.Lat, geoData.Lon, key_weather)

	response, err = http.Get(weatherURL)
	if err != nil {
		log.Fatal(err)
	}
	defer response.Body.Close()

	weatherData, err := io.ReadAll(response.Body)
	if err != nil {
		log.Fatal(err)
	}

	var weather WeatherData
	err = json.Unmarshal(weatherData, &weather)
	if err != nil {
		log.Fatal(err)
	}

	weather.Main.Temp = kelvinToCelsius(weather.Main.Temp)
	weather.Main.FeelsLike = kelvinToCelsius(weather.Main.FeelsLike)

	fmt.Printf("Temperature: %.1f°C, FeelsLike: %.1f°C, Humidity: %d%%\n", weather.Main.Temp, weather.Main.FeelsLike, weather.Main.Humidity)
}
