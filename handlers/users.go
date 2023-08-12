package handlers

import (
	"encoding/json"
	"github.com/google/uuid"
	"github.com/justsorrent/game-planner/internal/db"
	"golang.org/x/crypto/bcrypt"
	"net/http"
)

type userDto struct {
	Username string `json:"username"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

type userInfoResource struct {
	ID       string `json:"id"`
	Username string `json:"username"`
}

func (cfg *ApiConfig) HandleCreateUser(w http.ResponseWriter, r *http.Request) {
	dto := userDto{}
	err := json.NewDecoder(r.Body).Decode(&dto)
	if err != nil {
		RespondWithError(w, http.StatusBadRequest, err.Error())
		return
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(dto.Password), bcrypt.DefaultCost)
	if err != nil {
		RespondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	userId, err := uuid.NewUUID()
	if err != nil {
		RespondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	userCreds, err := cfg.DB.RegisterUser(r.Context(), db.RegisterUserParams{
		ID:           userId,
		Email:        dto.Email,
		PasswordHash: string(hashedPassword),
	})
	if err != nil {
		RespondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	userId = userCreds.ID
	_, err = cfg.DB.CreateUser(r.Context(), db.CreateUserParams{
		ID:          userId,
		DisplayName: dto.Username,
	})
	if err != nil {
		RespondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	sessionToken := createSession(userId)
	SetSessionCookie(w, sessionToken)
}

func createSession(userId uuid.UUID) string {
	sessionToken := uuid.NewString()
	sessionCacheInstance.update(sessionToken, Session{userId})
	return sessionToken
}

func transformDbUserToResource(user db.User) userInfoResource {
	return userInfoResource{
		ID:       user.ID.String(),
		Username: user.DisplayName,
	}
}
