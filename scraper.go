package main

import (
	"context"
	"database/sql"
	"log"
	"strings"
	"sync"
	"time"

	"github.com/LongDude/GoRssProject/db"
	"github.com/google/uuid"
)

// Background Worker
func startScraping(
	db *db.Queries, 
  concurrency int, 
  timeBetweenRequests time.Duration,
) {
  log.Printf("Scraping on %v goroutines every %s duration", concurrency, timeBetweenRequests)
  ticker := time.NewTicker(timeBetweenRequests)
  for ; ; <-ticker.C {
    feeds, err := db.GetNextFeedsToFetch(
      context.Background(), 
      int32(concurrency),
    )    
    if err != nil{
      log.Println("error fetching feeds:", err)
      continue
    }

    wg := &sync.WaitGroup{}
    for _, feed := range feeds {
      wg.Add(1)

      go scrapeFeed(db, wg, feed)
    }
    wg.Wait()
  }
}

func scrapeFeed(dbQueries *db.Queries, wg * sync.WaitGroup, feed db.Feed){
  defer wg.Done()

  _, err := dbQueries.MarkFeedAsFetched(context.Background(), feed.ID)
  if err != nil{
    log.Println("error marking feed as fetched:", err)
    return
  }

  rssFeed, err := urlToFeed(feed.Url)
  if err != nil{
    log.Println("error fetching feed:", err)
    return
  }

  for _, item := range rssFeed.Channel.Item {
    desctiption := sql.NullString{}
    if item.Description != ""{
      desctiption.String = item.Description
      desctiption.Valid = true
    }

    pubAt, err := time.Parse(time.RFC1123Z, item.PubDate)
    if err != nil{
      log.Printf("couldn't parse date %v with err %v", item.PubDate, err)
      continue
    }

    _, err = dbQueries.CreatePost(
      context.Background(),
      db.CreatePostParams{
        ID: uuid.New(),
        CreatedAt: time.Now(),
        UpdatedAt: time.Now(),
        Title: item.Title,
        Description: desctiption,
        PublishedAt: pubAt,
        Url: item.Link,
        FeedID: feed.ID,
      },
    )

    if err != nil{
      if strings.Contains(err.Error(), "duplicate key"){
        continue
      }
      log.Println("failed to create post:", err)
    }
  }
  log.Printf("Feed %s collected, %v posts found", feed.Name, len(rssFeed.Channel.Item))
}
