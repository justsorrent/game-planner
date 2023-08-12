package handlers

import (
	"encoding/json"
	"errors"
	"fmt"
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

	//check if sessions is already cached
	session, sessionToken, _ := CheckSessionCookie(r)
	if session.UserId == userCreds.ID {
		sessionCacheInstance.update(sessionToken, session)
		SetSessionCookie(w, sessionToken)
		RespondWithJSON(w, http.StatusOK, nil)
		return
	}

	//check if password matches db
	err = bcrypt.CompareHashAndPassword([]byte(userCreds.PasswordHash), []byte(sessionDto.Password))
	if err != nil {
		RespondWithError(w, http.StatusBadRequest, "invalid credentials")
		return
	}

	sessionToken = createSession(userCreds.ID)
	SetSessionCookie(w, sessionToken)
	RespondWithJSON(w, http.StatusCreated, nil)
}

func (cfg *ApiConfig) HandleRefreshSession(w http.ResponseWriter, r *http.Request) {
	session, sessionToken, err := CheckSessionCookie(r)
	if err != nil {
		RespondWithError(w, http.StatusUnauthorized, err.Error())
	}
	sessionCacheInstance.update(sessionToken, session)
	SetSessionCookie(w, sessionToken)
	RespondWithJSON(w, http.StatusOK, nil)
}

func (cfg *ApiConfig) HandleDeleteSession(w http.ResponseWriter, r *http.Request) {
	_, sessionToken, err := CheckSessionCookie(r)
	if err != nil {
		RespondWithError(w, http.StatusUnauthorized, err.Error())
	}
	sessionCacheInstance.delete(sessionToken)
	SetSessionCookie(w, "")
	RespondWithJSON(w, http.StatusNoContent, nil)
}

func SetSessionCookie(w http.ResponseWriter, sessionToken string) {
	http.SetCookie(w, &http.Cookie{
		Name:    "session_token",
		Value:   sessionToken,
		Expires: time.Now().Add(defaultExpiration),
	})
}

func CheckSessionCookie(r *http.Request) (Session, string, error) {
	c, err := r.Cookie("session_token")
	if err != nil {
		if errors.Is(err, http.ErrNoCookie) {
			return Session{}, "", fmt.Errorf("no cookie found")
		}
		return Session{}, "", fmt.Errorf("malformed cookie")
	}
	sessionToken := c.Value
	session, found := sessionCacheInstance.get(sessionToken)
	if !found {
		return session, "", fmt.Errorf("no session found in cache")
	}
	return session, sessionToken, nil
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

func (c *sessionCache) get(sessionToken string) (Session, bool) {
	session, found := c.sessions.Get(sessionToken)
	if !found {
		return Session{}, false
	}
	return session.(Session), true
}

func (c *sessionCache) update(sessionToken string, session Session) {
	c.sessions.Set(sessionToken, session, cache.DefaultExpiration)
}

func (c *sessionCache) delete(sessionToken string) {
	c.sessions.Delete(sessionToken)
}

var sessionCacheInstance = newCache()
