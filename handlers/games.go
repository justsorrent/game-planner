package handlers

import (
	"github.com/justsorrent/game-planner/internal/db"
	"net/http"
	"time"
)

type gameDto struct {
	Name        string    `json:"name"`
	Description string    `json:"description"`
	Url         string    `json:"url"`
	StartTime   time.Time `json:"startTime"`
	EndTime     time.Time `json:"endTime"`
}

type gameResource struct {
	ID          string    `json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	Url         string    `json:"url"`
	StartTime   time.Time `json:"startTime"`
	EndTime     time.Time `json:"endTime"`
	CreatedAt   time.Time `json:"createdAt"`
	UpdatedAt   time.Time `json:"updatedAt"`
}

func (cfg *ApiConfig) CreateGame(w http.ResponseWriter, r *http.Request) {
	// ...
}

func (cfg *ApiConfig) HandleGetGames(w http.ResponseWriter, r *http.Request) {
	games, err := cfg.DB.GetGames(r.Context())
	if err != nil {
		RespondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}
	data := transformDbGanesToResources(games)
	RespondWithJSON(w, http.StatusOK, data)
}

func transformDbGanesToResources(games []db.Game) []gameResource {
	resources := make([]gameResource, len(games))
	for i, game := range games {
		resources[i] = transformDbGameToResource(game)
	}
	return resources
}

func transformDbGameToResource(game db.Game) gameResource {
	return gameResource{
		ID:          game.ID.String(),
		Name:        game.Name,
		Description: game.Description.String,
		Url:         game.Url.String,
		StartTime:   game.StartingAt.Time,
		EndTime:     game.EndingAt.Time,
		CreatedAt:   game.CreatedAt,
		UpdatedAt:   game.UpdatedAt,
	}
}
