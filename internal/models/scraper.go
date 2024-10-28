package models

import (
	"context"
	"database/sql"
	"github.com/google/uuid"
	"log"
	"rss-scraper/internal/database"
	"strings"
	"sync"
	"time"
)

func StartScraping(db *database.Queries, concurrency int, timeBetweenRequest time.Duration) {
	log.Printf("Scraping on %v gorutines every %s duration", concurrency, timeBetweenRequest)
	ticker := time.NewTicker(timeBetweenRequest)
	for ; ; <-ticker.C {
		feeds, err := db.GetNextFeedsToFetch(context.Background(), int32(concurrency))
		if err != nil {
			log.Println("error fetching feeds: ", err)
			continue
		}
		wg := sync.WaitGroup{}
		for _, feed := range feeds {
			wg.Add(1)
			go scrapeFeed(db, &wg, feed)
		}
		wg.Wait()
	}
}

func scrapeFeed(db *database.Queries, wg *sync.WaitGroup, feed database.Feed) {
	defer wg.Done()
	if _, err := db.MarkFeedAsFetched(context.Background(), feed.ID); err != nil {
		log.Printf("error updating feed with id %s: %v", feed.ID, err)
		return
	}
	rssFeed, err := UrlToFeed(feed.Url)
	if err != nil {
		log.Printf("error parsing feed url %s: %v", feed.Url, err)
		return
	}
	for _, item := range rssFeed.Channel.Item {
		description := sql.NullString{}
		if item.Description != "" {
			description.String = item.Description
			description.Valid = true
		}
		pubTime, err := time.Parse(time.RFC1123Z, item.PubDate)
		if err != nil {
			log.Printf("error parsing pub date from post %s: %v", item.Title, err)
			continue
		}
		_, err = db.CreatePost(context.Background(), database.CreatePostParams{
			ID:          uuid.New(),
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
			Title:       item.Title,
			Description: description,
			PublishedAt: pubTime,
			Url:         item.Link,
			FeedID:      feed.ID,
		})
		if err != nil {
			if strings.Contains(err.Error(), "duplicate key") {
				continue
			}
			log.Printf("failed creating post in db: %v", err)
		}
	}

	log.Printf("feed %s collected, %d posts found", feed.Name, len(rssFeed.Channel.Item))
}
