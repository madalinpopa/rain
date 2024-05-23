package main

import (
	"fmt"
	"log"

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
func getWeather(city weather.City, apiKey string) (weather.Weather, error) {
	weatherData, err := weather.FetchWeatherData(apiKey, fmt.Sprintf("%f", city.Lat), fmt.Sprintf("%f", city.Lon))
	if err != nil {
		return weather.Weather{}, err
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

	fmt.Println(weather, err)

	fmt.Printf("City: %s, Lat: %f, Lon: %f, Country: %s\n", city.Name, city.Lat, city.Lon, city.Country)

}
