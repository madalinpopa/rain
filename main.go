package main

import (
	"fmt"

	"log"
	"time"

	"github.com/fatih/color"
	"github.com/madalinpopa/rain/weather"
)

// getLocation retrieves the location data
func getLocation(location, apiKey string) (weather.City, error) {
	cities, err := weather.FetchCityData(apiKey, location)
	if err != nil {
		return weather.City{}, err
	}

	return cities[0], nil
}

// getWeather retrieves the weather data
func getWeather(city weather.City, apiKey string) (weather.Forecasts, error) {
	weatherData, err := weather.FetchWeatherData(
		apiKey, fmt.Sprintf("%f", city.Lat), fmt.Sprintf("%f", city.Lon),
	)
	if err != nil {
		return weather.Forecasts{}, err
	}

	return weatherData, nil
}

func main() {
	defaultLocation := "Bucharest"

	apiKey, err := weather.GetApiEnv("OPEN_WEATHER_API_KEY")
	if err != nil {
		log.Fatalf("Error getting API key: %v", err)
	}

	city, err := getLocation(defaultLocation, apiKey)
	if err != nil {
		log.Fatalf("Error getting location: %v", err)
	}

	weather, err := getWeather(city, apiKey)
	if err != nil {
		log.Fatalf("Error getting weather: %v", err)
	}

	color.Red("\nCurrent Temperature: %.0fC\n\n", weather.Current.Temp)

	for _, h := range weather.Hourly {
		date := time.Unix(h.Dt, 0)
		if date.Before(time.Now()) || date.After(time.Now().Add(12*time.Hour)) {
			continue
		}
		message := fmt.Sprintf(
			"%s - Temp: %.0fC, Feels like: %.0fC, Chances to rain: %.0f%% - %s\n",
			date.Format("15:04"),
			h.Temp,
			h.FeelsLike,
			h.Pop*100,
			h.Weather[0].Description,
		)

		// If the temperature is lower than 20 degrees, print the message in blue
		if h.Temp < 20 {
			color.Blue(message)
		} else {
			color.Yellow(message)
		}

		// If the chances to rain are higher than 60%, print the message in red
		if h.Pop*100 > 60 {
			color.Red(message)
		}
	}
}
