package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/LongDude/GoRssProject/db"
	"github.com/google/uuid"
)

func (apiCfg *apiConfig) handlerCreateUser(w http.ResponseWriter, r *http.Request) {
	type parametres struct {
		Name string `json:"name"`
	}
	decoder := json.NewDecoder(r.Body)
	params := parametres{}
	err := decoder.Decode(&params)
	if err != nil {
		respondWithError(w, 400, fmt.Sprintf("Error parsing JSON: %v", err))
		return
	}

	user, err := apiCfg.DB.CreateUser(r.Context(), db.CreateUserParams{
		ID:        uuid.New(),
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
		Name:      params.Name,
	})

	if err != nil {
		respondWithError(w, 400, fmt.Sprintf("Couldn`t create user: %v", err))
		return
	}

	respondWithJSON(w, 201, databaseUserToUser(user))
}

func (apiCfg *apiConfig) handlerGetUserByAPIKey(w http.ResponseWriter, r *http.Request, user db.User) {
	respondWithJSON(w, 200, databaseUserToUser(user))
}

func (apiCfg *apiConfig) handlerGetPostsByUser(w http.ResponseWriter, r *http.Request, user db.User) {
  posts, err := apiCfg.DB.GetPostsForUser(r.Context(), db.GetPostsForUserParams{
    UserID: user.ID,
    Limit: 10,
  })

  if err != nil{
    respondWithError(w, 400, fmt.Sprintf("Couldn't get posts: %v", err))
    return
  }

  respondWithJSON(w, 200, posts)
}
