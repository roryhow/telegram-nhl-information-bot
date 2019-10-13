package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"
)

const BaseUrl string = "https://statsapi.web.nhl.com/api/v1"

type ScheduleResponse struct {
	Copyright    string  `json:"copyright"`
	TotalItems   int     `json:"totalItems"`
	TotalEvents  int     `json:"totalEvents"`
	TotalGames   int     `json:"totalGames"`
	TotalMatches int     `json:"totalMatches"`
	Wait         int     `json:"wait"`
	Dates        []Dates `json:"dates"`
}
type Status struct {
	AbstractGameState string `json:"abstractGameState"`
	CodedGameState    string `json:"codedGameState"`
	DetailedState     string `json:"detailedState"`
	StatusCode        string `json:"statusCode"`
	StartTimeTBD      bool   `json:"startTimeTBD"`
}
type LeagueRecord struct {
	Wins   int    `json:"wins"`
	Losses int    `json:"losses"`
	Ot     int    `json:"ot"`
	Type   string `json:"type"`
}
type Team struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
	Link string `json:"link"`
}
type Away struct {
	LeagueRecord LeagueRecord `json:"leagueRecord"`
	Score        int          `json:"score"`
	Team         Team         `json:"team"`
}
type Home struct {
	LeagueRecord LeagueRecord `json:"leagueRecord"`
	Score        int          `json:"score"`
	Team         Team         `json:"team"`
}
type Teams struct {
	Away Away `json:"away"`
	Home Home `json:"home"`
}
type Venue struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
	Link string `json:"link"`
}

type Content struct {
	Link string `json:"link"`
}

type Games struct {
	GamePk   int       `json:"gamePk"`
	Link     string    `json:"link"`
	GameType string    `json:"gameType"`
	Season   string    `json:"season"`
	GameDate time.Time `json:"gameDate"`
	Status   Status    `json:"status"`
	Teams    Teams     `json:"teams"`
	Venue    Venue     `json:"venue,omitempty"`
	Content  Content   `json:"content"`
}
type Dates struct {
	Date         string        `json:"date"`
	TotalItems   int           `json:"totalItems"`
	TotalEvents  int           `json:"totalEvents"`
	TotalGames   int           `json:"totalGames"`
	TotalMatches int           `json:"totalMatches"`
	Games        []Games       `json:"games"`
	Events       []interface{} `json:"events"`
	Matches      []interface{} `json:"matches"`
}

func GetYesterdayScores() string {
	yesterday := time.Now().AddDate(0, 0, -1).Format("2006-01-02")
	resp, err := http.Get(BaseUrl + "/schedule?date=" + yesterday)
	if err != nil {
		log.Fatal(err)
	}

	var result ScheduleResponse
	json.NewDecoder(resp.Body).Decode(&result)
	var buffer bytes.Buffer
	for _, date := range result.Dates {
		for _, game := range date.Games {
			homeName := game.Teams.Away.Team.Name
			homeScore := game.Teams.Away.Score
			awayName := game.Teams.Home.Team.Name
			awayScore := game.Teams.Home.Score
			buffer.WriteString(fmt.Sprintf("%s %d - %s %d\n", homeName, homeScore, awayName, awayScore))
		}
	}
	return buffer.String()
}
