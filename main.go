package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
)

type City struct {
	Name    string  `json:"name"`
	Lat     float64 `json:"lat"`
	Lon     float64 `json:"lon"`
	Country string  `json:"country"`
}

func getApiEnv(envName string) (string, error) {
	api := os.Getenv(envName)
	if api == "" {
		return "", errors.New("api key is empty")
	}
	return api, nil
}

func fetchCityData(apiKey, location string) ([]City, error) {

	baseUrl := "https://api.openweathermap.org/geo/1.0/direct?"
	url := fmt.Sprintf("%sq=%s&limit=1&appid=%s", baseUrl, location, apiKey)

	respBody, err := makeGetRequest(url)
	if err != nil {
		return nil, fmt.Errorf("failed to make GET request: %w", err)
	}

	cities, err := unmarshalCities(respBody)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal JSON: %w", err)
	}

	if len(cities) == 0 {
		return nil, fmt.Errorf("no cities found for location: %s", location)
	}
	return cities, nil

}

func makeGetRequest(url string) ([]byte, error) {
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	return io.ReadAll(resp.Body)
}

func unmarshalCities(body []byte) ([]City, error) {
	var cities []City
	if err := json.Unmarshal(body, &cities); err != nil {
		return nil, err
	}
	return cities, nil
}

// getLocation retrieves the location data
func getLocation(location string) (City, error) {
	apiKey, err := getApiEnv("OPEN_WEATHER_API_KEY")
	if err != nil {
		return City{}, err
	}

	cities, err := fetchCityData(apiKey, location)
	if err != nil {
		return City{}, err
	}

	return cities[0], nil
}

func main() {

	defaultLocation := "Bucharest"
	city, err := getLocation(defaultLocation)
	if err != nil {
		log.Fatalf("Error getting location: %v", err)
	}

	fmt.Printf("City: %s, Lat: %f, Lon: %f, Country: %s\n", city.Name, city.Lat, city.Lon, city.Country)

}
