package game

import (
	"gin-rest-api-example/internal/game/model"
	"strconv"
)

type GameModeResponse struct {
	GameMode GameMode `json:"gamemode"`
}

type GameModesResponse struct {
	GameModes     []GameMode `json:"gamemodes"`
	GameModeCount int        `json:"gamemodesCount"`
}

type GameMode struct {
	id   string `json:"id"`
	mode string `json:"mode"`
}

// NewArticlesResponse converts article models and total count to ArticlesResponse
func NewGameModesResponse(gameModes []*model.GameMode, total int) *GameModesResponse {
	var gModes []GameMode
	for _, v := range gameModes {
		gModes = append(gModes, *NewGameModeResponse(v))
	}
	return &GameModesResponse{
		GameModes:     gModes,
		GameModeCount: total,
	}
}

func NewGameModeResponse(g *model.GameMode) *GameMode {
	return &GameMode{
		id:   strconv.Itoa(int(g.ID)),
		mode: g.Mode,
	}
}
