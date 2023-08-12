package handlers

import (
	"database/sql"
	"encoding/json"
	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
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
	ID          string           `json:"id"`
	Name        string           `json:"name"`
	Description string           `json:"description"`
	Url         string           `json:"url"`
	StartTime   time.Time        `json:"startTime"`
	EndTime     time.Time        `json:"endTime"`
	CreatedAt   time.Time        `json:"createdAt"`
	UpdatedAt   time.Time        `json:"updatedAt"`
	GameMaster  userInfoResource `json:"gameMaster"`
}

func (cfg *ApiConfig) HandleCreateGame(w http.ResponseWriter, r *http.Request, user db.User) {
	dto := gameDto{}
	err := json.NewDecoder(r.Body).Decode(&dto)
	if err != nil {
		RespondWithError(w, http.StatusBadRequest, err.Error())
		return
	}
	newGameId, err := uuid.NewUUID()
	if err != nil {
		RespondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}
	createdGame, err := cfg.DB.CreateGame(r.Context(), db.CreateGameParams{
		ID:          newGameId,
		Name:        dto.Name,
		Description: sql.NullString{String: dto.Description, Valid: dto.Description != ""},
		Url:         sql.NullString{String: dto.Url, Valid: dto.Url != ""},
		StartingAt:  sql.NullTime{Time: dto.StartTime, Valid: true},
		EndingAt:    sql.NullTime{Time: dto.EndTime, Valid: true},
		GmID:        uuid.NullUUID{UUID: user.ID, Valid: true},
	})
	if err != nil {
		RespondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	res := transformDbGameToResourceWithGmInfo(createdGame, user)
	RespondWithJSON(w, http.StatusCreated, res)
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

func (cfg *ApiConfig) HandleGetGameById(w http.ResponseWriter, r *http.Request) {
	gameId, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		RespondWithError(w, http.StatusBadRequest, err.Error())
		return
	}
	game, err := cfg.DB.GetGameById(r.Context(), gameId)
	if err != nil {
		RespondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}
	gm, err := cfg.DB.GetUserById(r.Context(), game.GmID.UUID)
	if err != nil {
		RespondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}
	data := transformDbGameToResourceWithGmInfo(game, gm)
	RespondWithJSON(w, http.StatusOK, data)
}

func (cfg *ApiConfig) HandleDeleteGameById(w http.ResponseWriter, r *http.Request, user db.User) {
	gameId, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		RespondWithError(w, http.StatusBadRequest, err.Error())
		return
	}

	game, err := cfg.DB.GetGameById(r.Context(), gameId)
	if err != nil {
		RespondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}
	if game.GmID.UUID != user.ID {
		RespondWithError(w, http.StatusForbidden, "You are not the game master")
		return
	}

	err = cfg.DB.DeleteGame(r.Context(), gameId)
	if err != nil {
		RespondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}
	RespondWithJSON(w, http.StatusNoContent, nil)
}

func (cfg *ApiConfig) HandleUpdateGameById(w http.ResponseWriter, r *http.Request, user db.User) {
	gameId, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		RespondWithError(w, http.StatusBadRequest, err.Error())
		return
	}

	game, err := cfg.DB.GetGameById(r.Context(), gameId)
	if err != nil {
		RespondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}
	if game.GmID.UUID != user.ID {
		RespondWithError(w, http.StatusForbidden, "You are not the game master")
		return
	}

	dto := gameDto{}
	err = json.NewDecoder(r.Body).Decode(&dto)
	if err != nil {
		RespondWithError(w, http.StatusBadRequest, err.Error())
		return
	}
	err = cfg.DB.UpdateGame(r.Context(), db.UpdateGameParams{
		ID:          gameId,
		Name:        dto.Name,
		Description: sql.NullString{String: dto.Description, Valid: dto.Description != ""},
		Url:         sql.NullString{String: dto.Url, Valid: dto.Url != ""},
		StartingAt:  sql.NullTime{Time: dto.StartTime, Valid: true},
		EndingAt:    sql.NullTime{Time: dto.EndTime, Valid: true},
	})
	if err != nil {
		RespondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}
	updatedGame, err := cfg.DB.GetGameById(r.Context(), gameId)
	data := transformDbGameToResource(updatedGame)
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

func transformDbGameToResourceWithGmInfo(game db.Game, gm db.User) gameResource {
	return gameResource{
		ID:          game.ID.String(),
		Name:        game.Name,
		Description: game.Description.String,
		Url:         game.Url.String,
		StartTime:   game.StartingAt.Time,
		EndTime:     game.EndingAt.Time,
		CreatedAt:   game.CreatedAt,
		UpdatedAt:   game.UpdatedAt,
		GameMaster:  transformDbUserToResource(gm),
	}
}
