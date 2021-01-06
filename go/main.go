package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/dghubble/go-twitter/twitter"
	"github.com/dghubble/oauth1"
	"github.com/joho/godotenv"
)

// Credentials stores all of our access/consumer tokens
// and secret keys needed for authentication against
// the twitter REST API.
type Credentials struct {
	ConsumerKey       string
	ConsumerSecret    string
	AccessToken       string
	AccessTokenSecret string
}

// struct for the data readed from the backend
type sensorAQI struct {
	Sensor      string  `json:"sensor"`
	Source      string  `json:"source"`
	Description string  `json:"description"`
	Longitude   float64 `json:"longitude"`
	Latitude    float64 `json:"latitude"`
	Quality     struct {
		Category string `json:"category"`
		Index    int    `json:"index"`
	} `json:"quality"`
} // sensorAQI ...

// to string implementation for sensorAQI struct
func (s *sensorAQI) String() string {
	var sensorName string
	if s.Description == "" {
		sensorName = fmt.Sprintf("Sensor %v", s.Source)
	} else {
		sensorName = s.Description
	}

	return fmt.Sprintf("%s: %d - %s\n", sensorName, s.Quality.Index, airQualityCategory(&s.Quality.Index))
} // *sensorAQI => String ...

// maps AQI values to a more readable tag
func airQualityCategory(i *int) string {
	var result string
	switch {
	case *i < 51:
		result = "👍 Bueno"
		break
	case *i < 101:
		result = "😐 Moderado"
		break
	case *i < 151:
		result = "⚠👴 No tan bueno"
		break
	case *i < 201:
		result = "⚠😷 Insalubre"
		break
	case *i < 301:
		result = "☣️ Muy Insalubre"
		break
	case *i >= 301:
		result = "☠️ Peligroso"
		break
	}
	return result
} // AirQualityCategory ...

// Fetches the sensors list with their air quality index (AQI) from the backend
func getSensorsData() (*[]sensorAQI, error) {
	var sensors []sensorAQI

	response, err := http.Get(os.Getenv("API_URL"))

	if err != nil {
		return nil, err
	}

	responseData, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}

	json.Unmarshal(responseData, &sensors)

	return &sensors, nil
} // getSensorsData ...

// Creates the text to be posted
func tweetBody() string {
	var buffer bytes.Buffer

	// get sensors list
	s, err := getSensorsData()
	if err != nil {
		log.Fatal(err)
		os.Exit(-1)
	}

	// write header
	buffer.WriteString("AireLibre PY #AireLibre\n")
	buffer.WriteString(time.Now().Format("02/Jan/2006 03:04\n\n"))

	// write sensors data
	for _, sensor := range *s {
		buffer.WriteString(sensor.String())
	}

	// write footer
	buffer.WriteString("\nMás info en http://airelib.re")

	return buffer.String()
} // tweetBody ...

// getClient is a helper function that will return a twitter client
// that we can subsequently use to send tweets, or to stream new tweets
// this will take in a pointer to a Credential struct which will contain
// everything needed to authenticate and return a pointer to a twitter Client
// or an error
func getClient(creds *Credentials) (*twitter.Client, error) {
	// Pass in your consumer key (API Key) and your Consumer Secret (API Secret)
	config := oauth1.NewConfig(creds.ConsumerKey, creds.ConsumerSecret)
	// Pass in your Access Token and your Access Token Secret
	token := oauth1.NewToken(creds.AccessToken, creds.AccessTokenSecret)

	httpClient := config.Client(oauth1.NoContext, token)
	client := twitter.NewClient(httpClient)

	// Verify Credentials
	verifyParams := &twitter.AccountVerifyParams{
		SkipStatus:   twitter.Bool(true),
		IncludeEmail: twitter.Bool(true),
	}

	// we can retrieve the user and verify if the credentials
	// we have used successfully allow us to log in!
	user, _, err := client.Accounts.VerifyCredentials(verifyParams)
	if err != nil {
		return nil, err
	}

	log.Printf("User's ACCOUNT:\n%+v\n", user)
	return client, nil
}

func main() {
	_ = godotenv.Load()

	fmt.Println("Linka TwitterBot v0.01")
	creds := Credentials{
		AccessToken:       os.Getenv("ACCESS_TOKEN"),
		AccessTokenSecret: os.Getenv("ACCESS_TOKEN_SECRET"),
		ConsumerKey:       os.Getenv("CONSUMER_KEY"),
		ConsumerSecret:    os.Getenv("CONSUMER_SECRET"),
	}

	fmt.Printf("%+v\n", creds)

	client, err := getClient(&creds)
	if err != nil {
		log.Println("Error getting Twitter Client")
		log.Println(err)
	}

	tweet, resp, err := client.Statuses.Update(tweetBody(), nil)
	if err != nil {
		log.Println(err)
	}
	log.Printf("%+v\n", resp)
	log.Printf("%+v\n", tweet)
}
