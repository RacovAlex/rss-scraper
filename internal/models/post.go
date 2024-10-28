package models

import (
	"time"

	"github.com/google/uuid"
	"rss-scraper/internal/database"
)

type Post struct {
	ID          uuid.UUID `json:"id"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
	Title       string    `json:"title"`
	Description *string   `json:"description"`
	PublishedAt time.Time `json:"published_at"`
	Url         string    `json:"url"`
	FeedID      uuid.UUID `json:"feed_id"`
}

func databasePostToPost(dbPost database.Post) *Post {
	var description *string
	if dbPost.Description.Valid {
		description = &dbPost.Description.String
	}
	return &Post{
		ID:          dbPost.ID,
		CreatedAt:   dbPost.CreatedAt,
		UpdatedAt:   dbPost.UpdatedAt,
		Title:       dbPost.Title,
		Description: description,
		PublishedAt: dbPost.PublishedAt,
		Url:         dbPost.Url,
		FeedID:      dbPost.FeedID,
	}
}

func DatabasePostsToPosts(dbPosts []database.Post) []Post {
	posts := make([]Post, 0, len(dbPosts))
	for _, post := range dbPosts {
		posts = append(posts, *databasePostToPost(post))
	}
	return posts
}
