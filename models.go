package main

import (
	"project_1/internal/database"
	"time"

	"github.com/google/uuid"
)

// @name User
// @description A registered user of the application.
type User struct {
	ID        uuid.UUID `json:"id"`         // Unique user ID
	Name      string    `json:"name"`       // Full name
	CreatedAt time.Time `json:"created_at"` // Account creation timestamp
	UpdatedAt time.Time `json:"updated_at"` // Last update timestamp
	Email     string    `json:"email"`      // User email address
}

// @name UserInput
// @description Input model for creating or updating a user.
type UserInput struct {
	Name     string `json:"name"`     // Full name
	Email    string `json:"email"`    // Email address
	Password string `json:"password"` // Account password
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

// @name Feed
// @description Represents an RSS feed followed or owned by a user.
type Feed struct {
	ID     uuid.UUID `json:"id"`      // Feed ID
	Name   string    `json:"name"`    // Feed name
	UserID uuid.UUID `json:"user_id"` // Owner's user ID
	URL    string    `json:"url"`     // Feed URL
}

// @name FeedInput
// @description Input model for creating or updating a feed.
type FeedInput struct {
	Name string `json:"name"` // Feed name
	URL  string `json:"url"`  // Feed URL
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

// @name Follow
// @description A follow relationship between a user and a feed.
type Follow struct {
	UserID uuid.UUID `json:"user_id"` // ID of the user
	FeedID uuid.UUID `json:"feed_id"` // ID of the followed feed
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

// @name Post
// @description A post from an RSS feed.
type Post struct {
	ID          uuid.UUID `json:"id"`           // Post ID
	Title       string    `json:"title"`        // Post title
	Description *string   `json:"description"`  // Post description
	PublishedAt time.Time `json:"published_at"` // Publication timestamp
	Url         string    `json:"url"`          // Post URL
	FeedID      uuid.UUID `json:"feed_id"`      // Associated feed ID
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

// @name LoginResponse
// @description Token response after successful login.
type LoginResponse struct {
	Token     string `json:"token"`      // Access token
	TokenType string `json:"token_type"` // Token type (e.g., Bearer)
}

func ResponseToken(tokentype, access_token string) LoginResponse {
	return LoginResponse{Token: access_token, TokenType: tokentype}
}
