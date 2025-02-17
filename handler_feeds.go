package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/LongDude/GoRssProject/db"
	"github.com/google/uuid"
)

func (apiCfg *apiConfig) handlerCreateFeed(w http.ResponseWriter, r *http.Request, user db.User) {
	type parametres struct {
		Name string `json:"name"`
    URL string `json:"url"`
	}
	decoder := json.NewDecoder(r.Body)
	params := parametres{}
	err := decoder.Decode(&params)
	if err != nil {
		respondWithError(w, 400, fmt.Sprintf("Error parsing JSON: %v", err))
		return
	}
  
	feed, err := apiCfg.DB.CreateFeed(r.Context(), db.CreateFeedParams{

		ID:        uuid.New(),
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
		Name:      params.Name,
    Url: params.URL,
    UserID: user.ID,
	})

	if err != nil {
		respondWithError(w, 400, fmt.Sprintf("Couldn`t create feed: %v", err))
		return
	}
  respondWithJSON(w, 201, feed)
}

func (apiCfg *apiConfig) handlerGetFeeds(w http.ResponseWriter, r *http.Request) {
  feeds, err := apiCfg.DB.GetFeeds(r.Context())
  if err != nil{
    respondWithError(w, 400, fmt.Sprintf("Couldn`t get feeds: %v", err))
    return
  }

  respondWithJSON(w, 200, feeds)
}
