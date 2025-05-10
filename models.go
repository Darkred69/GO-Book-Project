package main

import (
	"project_1/internal/database"
	"time"

	"github.com/google/uuid"
)

// User models and functions
type User struct {
	ID        uuid.UUID `json:"id"`
	Name      string    `json:"name"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Email     string    `json:"email"`
}

type UserInput struct {
	Name     string `json:"name"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

func databaseUsertoUser(dbUser database.User) User {
	return User{
		ID:        dbUser.ID,
		Name:      dbUser.Name,
		CreatedAt: dbUser.CreatedAt,
		UpdatedAt: dbUser.UpdatedAt,
		Email:     dbUser.Email,
	}
}

// Feeds Model and functions
type Feed struct {
	ID   uuid.UUID `json:"id"`
	Name string    `json:"name"`
	// CreatedAt time.Time `json:"created_at"`
	// UpdatedAt time.Time `json:"updated_at"`
	UserID uuid.UUID `json:"user_id"`
	URL    string    `json:"url"`
}

type FeedInput struct {
	Name string `json:"name"`
	URL  string `json:"url"`
}

func databaseFeedtoFeed(dbFeed database.Feed) Feed {
	return Feed{
		ID:     dbFeed.ID,
		Name:   dbFeed.Name,
		UserID: dbFeed.UserID,
		URL:    dbFeed.Url,
	}
}

func databaseFeedstoFeeds(dbFeed []database.Feed) []Feed {
	feeds := []Feed{}
	for _, dbFedbFeed := range dbFeed {
		feeds = append(feeds, databaseFeedtoFeed(dbFedbFeed))
	}
	return feeds
}

// Follow models and functions
type Follow struct {
	UserID uuid.UUID `json:"user_id"`
	FeedID uuid.UUID `json:"feed_id"`
}

func databaseFollowtoFollow(dbFeed database.FeedFollow) Follow {
	return Follow{
		UserID: dbFeed.UserID,
		FeedID: dbFeed.FeedID,
	}
}

func databaseFollowstoFollows(dbFeed []database.FeedFollow) []Follow {
	follows := []Follow{}
	for _, dbFedbFeed := range dbFeed {
		follows = append(follows, databaseFollowtoFollow(dbFedbFeed))
	}
	return follows
}

// Post models and functions
type Post struct {
	ID          uuid.UUID `json:"id"`
	Title       string    `json:"title"`
	Description *string   `json:"description"`
	PublishedAt time.Time `json:"published_at"`
	Url         string    `json:"url"`
	FeedID      uuid.UUID `json:"feed_id"`
}

func databasePosttoPost(dbPost database.Post) Post {
	var description *string
	if dbPost.Description.Valid {
		description = &dbPost.Description.String
	}

	return Post{
		ID:          dbPost.ID,
		Title:       dbPost.Title,
		Description: description,
		PublishedAt: dbPost.PublishedAt,
		Url:         dbPost.Url,
		FeedID:      dbPost.FeedID,
	}
}

func databasePoststoPosts(dbPost []database.Post) []Post {
	posts := []Post{}
	for _, dbPost := range dbPost {
		posts = append(posts, databasePosttoPost(dbPost))
	}
	return posts
}

// LoginResponse represents the response with the token
type LoginResponse struct {
	Token     string `json:"token"`
	TokenType string `json:"token_type"`
}

func ResponseToken(tokentype, access_token string) LoginResponse {
	return LoginResponse{Token: access_token, TokenType: tokentype}
}
