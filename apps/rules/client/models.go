package client

import (
	"rules"
	"rules/settings"
)

// The top-level message sent in /start, /move, and /end requests
type SnakeRequest struct {
	Game  Game  `json:"game"`
	Turn  int   `json:"turn"`
	Board Board `json:"board"`
	You   Snake `json:"you"`
}

// Game represents the current game state
type Game struct {
	ID      string  `json:"id"`
	Ruleset Ruleset `json:"ruleset"`
	Map     string  `json:"map"`
	Timeout int     `json:"timeout"`
	Source  string  `json:"source"`
}

// Board provides information about the game board
type Board struct {
	Height int     `json:"height"`
	Width  int     `json:"width"`
	Snakes []Snake `json:"snakes"`
	Food   []Coord `json:"food"`
}

// Snake represents information about a snake in the game
type Snake struct {
	ID             string         `json:"id"`
	Name           string         `json:"name"`
	Latency        string         `json:"latency"`
	Health         int            `json:"health"`
	Body           []Coord        `json:"body"`
	Head           Coord          `json:"head"`
	Length         int            `json:"length"`
	Shout          string         `json:"shout"`
	Customizations Customizations `json:"customizations"`
}

type Customizations struct {
	Color string `json:"color"`
	Head  string `json:"head"`
	Tail  string `json:"tail"`
}

type Ruleset struct {
	Name     string          `json:"name"`
	Settings RulesetSettings `json:"settings"`
}

// RulesetSettings contains a static collection of a few settings that are exposed through the API.
type RulesetSettings struct {
	FoodSpawnChance int `json:"foodSpawnChance"`
}

// Converts a rules.Settings (which can contain arbitrary settings) into the static RulesetSettings used in the client API.
func ConvertRulesetSettings(settings settings.Settings) RulesetSettings {
	return RulesetSettings{
		FoodSpawnChance: settings.Int(rules.ParamFoodSpawnChance, 0),
	}
}

// Coord represents a point on the board
type Coord struct {
	X int `json:"x"`
	Y int `json:"y"`
}

// The expected format of the response body from a /move request
type MoveResponse struct {
	Move  string `json:"move"`
	Shout string `json:"shout"`
}

// The expected format of the response body from a GET request to a Battlesnake's index URL
type SnakeMetadataResponse struct {
	Author string `json:"author,omitempty"`
	Color  string `json:"color,omitempty"`
	Head   string `json:"head,omitempty"`
	Tail   string `json:"tail,omitempty"`
}

func CoordFromPoint(pt rules.Point) Coord {
	return Coord{X: pt.X, Y: pt.Y}
}

func CoordFromPointArray(ptArray []rules.Point) []Coord {
	a := make([]Coord, 0)
	for _, pt := range ptArray {
		a = append(a, CoordFromPoint(pt))
	}
	return a
}
