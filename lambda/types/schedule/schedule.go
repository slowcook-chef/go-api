package schedule

import (
	"encoding/json"
	"fmt"
	"io"
	"lambda-func/ledger"
	"lambda-func/secret"
	"net/http"
	"strconv"
	"time"
)

type Game interface {
	GetTeams() (away string, home string)
	GetTeamIDs() (awayID string, homeID string)
	GetDateTime() (date string, time string)
	GetID() string
	GetStatus() string
}

type Schedule struct {
	Sport      string
	TotalGames int
	GameData   struct {
		Games []Game `json:"MLB"`
	} `json:"data"`
}

type MLBGame struct {
	AwayTeam    string `json:"away_team"`
	HomeTeam    string `json:"home_team"`
	AwayTeamID  int    `json:"away_team_ID"`
	HomeTeamID  int    `json:"home_team_ID"`
	GameID      string `json:"game_ID"`
	GameTime    string `json:"game_time"`
	SeasonType  string `json:"season_type"`
	Season      string `json:"season"`
	EventName   any    `json:"event_name"`
	Round       any    `json:"round"`
	AwayPitcher struct {
		PlayerID int    `json:"player_id"`
		Player   string `json:"player"`
	} `json:"away_pitcher"`
	HomePitcher struct {
		PlayerID int    `json:"player_id"`
		Player   string `json:"player"`
	} `json:"home_pitcher"`
	Status     string  `json:"status"`
	City       string  `json:"city"`
	State      string  `json:"state"`
	Country    string  `json:"country"`
	PostalCode string  `json:"postal_code"`
	Dome       int     `json:"dome"`
	Field      string  `json:"field"`
	Latitude   float64 `json:"latitude"`
	Longitude  float64 `json:"longitude"`
	Arena      string  `json:"arena"`
}

func (m MLBGame) GetTeams() (away string, home string) {
	return m.AwayTeam, m.HomeTeam
}

func (m MLBGame) GetTeamIDs() (awayID string, homeID string) {
	return strconv.Itoa(m.AwayTeamID), strconv.Itoa(m.HomeTeamID)
}

func (m MLBGame) GetDateTime() (day string, minute string) {
	dayTime, err := time.Parse("Mon, 02 Jan 2006 15:04:05 MST", m.GameTime)
	if err != nil {
		fmt.Println("Error parsing game day:", err)
		return "", ""
	}
	minuteTime, err := time.Parse("Mon, 02 Jan 2006 15:04:05 MST", m.GameTime)
	if err != nil {
		fmt.Println("Error parsing game time:", err)
		return "", ""
	}

	return dayTime.Format("2006-01-02"), minuteTime.Format("15:04:05 MST")
}

func (m MLBGame) GetID() string {
	return m.GameID
}

func (m MLBGame) GetStatus() string {
	return m.Status
}

// Stinky
func GetMLBEndpointSchedule() (*Schedule, error) {
	//Recieve the plane
	url := fmt.Sprintf("http://rest.datafeeds.rolling-insights.com/api/v1/schedule/now/MLB?RSC_token=%s", secret.GetDFToken())
	response, err := http.Get(url)
	if err != nil {
		ledger.LogError(&err)
		return nil, err
	}
	//Process Immigration
	if response.StatusCode != http.StatusOK {
		ledger.LogHandlerProcess("not ok")
		return nil, err
	}
	defer response.Body.Close()

	body, err := io.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}
	//Translate languages
	var schedule Schedule
	err = json.Unmarshal(body, &schedule)
	schedule.Sport = "mlb"
	return &schedule, err
}
