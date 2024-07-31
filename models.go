package main

import (
	"time"

	"github.com/LongDude/GoRssProject/db"
	"github.com/google/uuid"
)

type User struct {
	ID        uuid.UUID `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Name      string    `json:"name"`
	APIKey    string    `json:"api_key"`
}

func databaseUserToUser(dbUser db.User) User {
	return User{
		ID:        dbUser.ID,
		CreatedAt: dbUser.CreatedAt,
		UpdatedAt: dbUser.UpdatedAt,
		Name:      dbUser.Name,
		APIKey:    dbUser.ApiKey,
	}
}

type Feed struct {
  ID        uuid.UUID `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Name      string    `json:"name"`
	Url       string    `json:"url"`
	UserID    uuid.UUID `json:"user_id"`
}

func databaseFeedToFeed(dbFeed db.Feed) Feed {
  return Feed{
    ID: dbFeed.ID,
    CreatedAt: dbFeed.CreatedAt,
    UpdatedAt: dbFeed.UpdatedAt,
    Name: dbFeed.Name,
    Url: dbFeed.Url,
    UserID: dbFeed.UserID,
  }
}

func databaseFeedsToFeeds(dbFeeds []db.Feed) []Feed{
  feeds := []Feed{}
  for _, dbFeed := range dbFeeds {
    feeds = append(feeds, databaseFeedToFeed(dbFeed))
  }
  return feeds
}
