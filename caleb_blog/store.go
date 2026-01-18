package main

import (
	"database/sql"
	"log"
	"time"

	_ "modernc.org/sqlite"
)

// Store provides persistent SQLite storage for posts
type Store struct {
	db *sql.DB
}

// NewStore creates a new SQLite-backed store
func NewStore(dbPath string) *Store {
	db, err := sql.Open("sqlite", dbPath)
	if err != nil {
		log.Fatal(err)
	}

	// Create posts table if it doesn't exist
	_, err = db.Exec(`
		CREATE TABLE IF NOT EXISTS posts (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			title TEXT NOT NULL,
			content TEXT NOT NULL,
			created_at DATETIME NOT NULL
		)
	`)
	if err != nil {
		log.Fatal(err)
	}

	return &Store{db: db}
}

// Close closes the database connection
func (s *Store) Close() error {
	return s.db.Close()
}

// Create adds a new post and returns its ID
func (s *Store) Create(title, content string) int {
	result, err := s.db.Exec(
		"INSERT INTO posts (title, content, created_at) VALUES (?, ?, ?)",
		title, content, time.Now(),
	)
	if err != nil {
		log.Printf("Error creating post: %v", err)
		return 0
	}

	id, err := result.LastInsertId()
	if err != nil {
		log.Printf("Error getting last insert ID: %v", err)
		return 0
	}
	return int(id)
}

// Get retrieves a post by ID, returns nil if not found
func (s *Store) Get(id int) *Post {
	row := s.db.QueryRow(
		"SELECT id, title, content, created_at FROM posts WHERE id = ?",
		id,
	)

	var post Post
	err := row.Scan(&post.ID, &post.Title, &post.Content, &post.CreatedAt)
	if err != nil {
		return nil
	}
	return &post
}

// List returns all posts sorted by creation time (newest first)
func (s *Store) List() []*Post {
	rows, err := s.db.Query(
		"SELECT id, title, content, created_at FROM posts ORDER BY created_at DESC",
	)
	if err != nil {
		log.Printf("Error listing posts: %v", err)
		return nil
	}
	defer rows.Close()

	var posts []*Post
	for rows.Next() {
		var post Post
		if err := rows.Scan(&post.ID, &post.Title, &post.Content, &post.CreatedAt); err != nil {
			log.Printf("Error scanning post: %v", err)
			continue
		}
		posts = append(posts, &post)
	}
	return posts
}

// Delete removes a post by ID, returns true if deleted
func (s *Store) Delete(id int) bool {
	result, err := s.db.Exec("DELETE FROM posts WHERE id = ?", id)
	if err != nil {
		log.Printf("Error deleting post: %v", err)
		return false
	}

	affected, err := result.RowsAffected()
	if err != nil {
		return false
	}
	return affected > 0
}
