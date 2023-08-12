package handlers

import (
	"encoding/json"
	"errors"
	"github.com/google/uuid"
	"github.com/patrickmn/go-cache"
	"golang.org/x/crypto/bcrypt"
	"net/http"
	"time"
)

func (cfg *ApiConfig) HandleCreateSession(w http.ResponseWriter, r *http.Request) {
	sessionDto := sessionDto{}
	err := json.NewDecoder(r.Body).Decode(&sessionDto)
	if err != nil {
		RespondWithError(w, http.StatusBadRequest, err.Error())
		return
	}

	userCreds, err := cfg.DB.GetUserByEmail(r.Context(), sessionDto.Email)
	if err != nil || userCreds.ID.String() == "" {
		RespondWithError(w, http.StatusBadRequest, "invalid credentials")
		return
	}

	//check if password matches db
	err = bcrypt.CompareHashAndPassword([]byte(userCreds.PasswordHash), []byte(sessionDto.Password))
	if err != nil {
		RespondWithError(w, http.StatusBadRequest, "invalid credentials")
		return
	}

	sessionToken := createSession(userCreds.ID)
	SetSessionCookie(w, sessionToken)
	RespondWithJSON(w, http.StatusCreated, nil)
}

func (cfg *ApiConfig) HandleRefreshSession(w http.ResponseWriter, r *http.Request) {
	c, err := r.Cookie("session_token")
	if err != nil {
		if errors.Is(err, http.ErrNoCookie) {
			RespondWithError(w, http.StatusUnauthorized, err.Error())
			return
		}
		RespondWithError(w, http.StatusBadRequest, err.Error())
		return
	}
	sessionToken := c.Value
	session, found := sessionCacheInstance.get(sessionToken)
	if !found {
		RespondWithError(w, http.StatusUnauthorized, "Invalid session token")
		return
	}
	sessionCacheInstance.update(sessionToken, session)
}

func SetSessionCookie(w http.ResponseWriter, sessionToken string) {
	http.SetCookie(w, &http.Cookie{
		Name:    "session_token",
		Value:   sessionToken,
		Expires: time.Now().Add(defaultExpiration),
	})
}

type sessionCache struct {
	sessions *cache.Cache
}

type sessionDto struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type Session struct {
	UserId uuid.UUID
}

const (
	defaultExpiration = 5 * time.Minute
	purgeTime         = 10 * time.Minute
)

func newCache() *sessionCache {
	Cache := cache.New(defaultExpiration, purgeTime)
	return &sessionCache{
		sessions: Cache,
	}
}

func (c *sessionCache) get(id string) (Session, bool) {
	session, found := c.sessions.Get(id)
	if !found {
		return Session{}, false
	}
	return session.(Session), true
}

func (c *sessionCache) update(id string, session Session) {
	c.sessions.Set(id, session, cache.DefaultExpiration)
}

var sessionCacheInstance = newCache()
