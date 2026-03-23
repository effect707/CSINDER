package faceit

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
)

type APIClient struct {
	apiKey string
}

type PlayerProfile struct {
	PlayerID  string `json:"player_id"`
	Nickname  string `json:"nickname"`
	Avatar    string `json:"avatar"`
	Country   string `json:"country"`
	FaceitURL string `json:"faceit_url"`
	Games     struct {
		CS2 *struct {
			FaceitElo  int `json:"faceit_elo"`
			SkillLevel int `json:"skill_level"`
		} `json:"cs2"`
	} `json:"games"`
}

type PlayerStats struct {
	Lifetime map[string]string `json:"lifetime"`
}

type PlayerStatsResponse struct {
	Lifetime map[string]string `json:"lifetime"`
}

func NewAPIClient(apiKey string) *APIClient {
	return &APIClient{apiKey: apiKey}
}

func (c *APIClient) GetPlayer(faceitID string) (*PlayerProfile, error) {
	req, err := http.NewRequest("GET", fmt.Sprintf("https://open.faceit.com/data/v4/players/%s", faceitID), nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Authorization", "Bearer "+c.apiKey)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("faceit api error: %s", resp.Status)
	}

	var profile PlayerProfile
	if err := json.NewDecoder(resp.Body).Decode(&profile); err != nil {
		return nil, err
	}
	return &profile, nil
}

func (c *APIClient) GetPlayerStats(faceitID string) (*PlayerStats, error) {
	req, err := http.NewRequest("GET", fmt.Sprintf("https://open.faceit.com/data/v4/players/%s/stats/cs2", faceitID), nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Authorization", "Bearer "+c.apiKey)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("faceit stats api error: %s", resp.Status)
	}

	var statsResp PlayerStatsResponse
	if err := json.NewDecoder(resp.Body).Decode(&statsResp); err != nil {
		return nil, err
	}
	return &PlayerStats{Lifetime: statsResp.Lifetime}, nil
}

func ParseFloat(s string) float64 {
	v, _ := strconv.ParseFloat(s, 64)
	return v
}

func ParseInt(s string) int {
	v, _ := strconv.Atoi(s)
	return v
}
