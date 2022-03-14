package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/joho/godotenv"
)

type Stats struct {
	Updated     int64  `json:"updated"`
	Country     string `json:"country"`
	CountryInfo struct {
		ID   int     `json:"_id"`
		Iso2 string  `json:"iso2"`
		Iso3 string  `json:"iso3"`
		Lat  float32 `json:"lat"`
		Long float32 `json:"long"`
		Flag string  `json:"flag"`
	} `json:"countryInfo"`
	Cases               int    `json:"cases"`
	TodayCases          int    `json:"todayCases"`
	Deaths              int    `json:"deaths"`
	TodayDeaths         int    `json:"todayDeaths"`
	Recovered           int    `json:"recovered"`
	Active              int    `json:"active"`
	Critical            int    `json:"critical"`
	CasesPerOneMillion  int    `json:"casesPerOneMillion"`
	DeathsPerOneMillion int    `json:"deathsPerOneMillion"`
	Tests               int    `json:"tests"`
	TestsPerOneMillion  int    `json:"testsPerOneMillion"`
	Continent           string `json:"continent"`
}

func getApiResponse() *http.Response {
	response, err := http.Get("https://disease.sh/v3/covid-19/countries/russia?strict=true")
	if err != nil {
		log.Panic(err)
	}
	return response
}

func readData(response *http.Response) []byte {
	defer response.Body.Close()
	content, err := ioutil.ReadAll(response.Body)
	if err != nil {
		log.Panic(err)
	}
	return content
}

func dataToJson(content []byte) Stats {
	var stats Stats
	err := json.Unmarshal(content, &stats)
	if err != nil {
		log.Panic(err)
	}
	return stats
}

func statsFromApi() Stats {
	response := getApiResponse()
	content := readData(response)
	stats := dataToJson(content)
	return stats
}

func statsByDay() string {
	stats := statsFromApi()
	answer := fmt.Sprintf("COVID-19 Russia:\ncases %d \ndeaths %d", stats.TodayCases, stats.TodayDeaths)
	return answer
}

func statsByAllTime() string {
	stats := statsFromApi()
	answer := fmt.Sprintf(
		"COVID-19 Russia:\ncases %d \ndeaths %d\nrecovered %d\ncases per mill %d\ndeaths per mill %d\n",
		stats.Cases, stats.Deaths, stats.Recovered, stats.CasesPerOneMillion, stats.DeathsPerOneMillion,
	)
	return answer
}

func parseCommand(text string) string {
	command := strings.Split(text, "@")[0]
	command = strings.ReplaceAll(command, "/", "")
	return command
}

func getEnvVar(key string) string {
	err := godotenv.Load()
	if err != nil {
		log.Panic(err)
	}

	return os.Getenv(key)
}

func main() {
	token := getEnvVar("BOT_TOKEN")
	bot, err := tgbotapi.NewBotAPI(token)
	if err != nil {
		log.Panic(err)
	}

	bot.Debug = false
	u := tgbotapi.NewUpdate(0)
	u.Timeout = 1

	updates, err := bot.GetUpdatesChan(u)
	for update := range updates {
		if update.Message == nil {
			continue
		}

		if !bot.IsMessageToMe(*update.Message) {
			continue
		}

		command := parseCommand(update.Message.Text)
        answer := "Who cares?"
		switch command {
		case "today":
            answer = "Who cares?"
		case "all":
            answer = "Who cares?"
		default:
			continue
		}

		msg := tgbotapi.NewMessage(update.Message.Chat.ID, answer)
		msg.ReplyToMessageID = update.Message.MessageID
		bot.Send(msg)
	}
}
