package weather

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
)

type City struct {
	Name    string  `json:"name"`
	Lat     float64 `json:"lat"`
	Lon     float64 `json:"lon"`
	Country string  `json:"country"`
}

type Weather struct {
	Temp float64 `json:"temp"`
}

func GetApiEnv(envName string) (string, error) {
	api := os.Getenv(envName)
	if api == "" {
		return "", errors.New("api key is empty")
	}
	return api, nil
}

func makeGetRequest(url string) ([]byte, error) {
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("received non-OK HTTP status: %d", resp.StatusCode)
	}

	return io.ReadAll(resp.Body)
}

func unmarshalCities(body []byte) ([]City, error) {
	var cities []City
	if err := json.Unmarshal(body, &cities); err != nil {
		return nil, err
	}
	return cities, nil
}

func FetchCityData(apiKey, location string) ([]City, error) {

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

func FetchWeatherData(apiKey, lat, lon string) (Weather, error) {

	baseUrl := "https://api.openweathermap.org/data/3.0/onecall?"
	url := fmt.Sprintf("%slat=%s&lon=%s&appid=%s", baseUrl, lat, lon, apiKey)

	respBody, err := makeGetRequest(url)
	if err != nil {
		return Weather{}, fmt.Errorf("failed to make GET request: %w", err)
	}

	fmt.Println(respBody)

	return Weather{}, nil
}
