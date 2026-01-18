package main

import "time"

// Post represents a blog post
type Post struct {
	ID        int
	Title     string
	Content   string
	CreatedAt time.Time
}
