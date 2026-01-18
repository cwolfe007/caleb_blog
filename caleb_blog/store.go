package main

import (
	"sync"
	"time"
)

// Store provides thread-safe in-memory storage for posts
type Store struct {
	mu     sync.RWMutex
	posts  map[int]*Post
	nextID int
}

// NewStore creates a new in-memory store
func NewStore() *Store {
	return &Store{
		posts:  make(map[int]*Post),
		nextID: 1,
	}
}

// Create adds a new post and returns its ID
func (s *Store) Create(title, content string) int {
	s.mu.Lock()
	defer s.mu.Unlock()

	post := &Post{
		ID:        s.nextID,
		Title:     title,
		Content:   content,
		CreatedAt: time.Now(),
	}
	s.posts[s.nextID] = post
	s.nextID++
	return post.ID
}

// Get retrieves a post by ID, returns nil if not found
func (s *Store) Get(id int) *Post {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.posts[id]
}

// List returns all posts sorted by creation time (newest first)
func (s *Store) List() []*Post {
	s.mu.RLock()
	defer s.mu.RUnlock()

	posts := make([]*Post, 0, len(s.posts))
	for _, p := range s.posts {
		posts = append(posts, p)
	}

	// Sort by CreatedAt descending (newest first)
	for i := 0; i < len(posts)-1; i++ {
		for j := i + 1; j < len(posts); j++ {
			if posts[j].CreatedAt.After(posts[i].CreatedAt) {
				posts[i], posts[j] = posts[j], posts[i]
			}
		}
	}
	return posts
}

// Delete removes a post by ID, returns true if deleted
func (s *Store) Delete(id int) bool {
	s.mu.Lock()
	defer s.mu.Unlock()

	if _, exists := s.posts[id]; exists {
		delete(s.posts, id)
		return true
	}
	return false
}
