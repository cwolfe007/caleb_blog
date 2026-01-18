package main

import (
	"fmt"
	"html/template"
	"net/http"
	"strconv"
	"strings"
)

// Handlers holds dependencies for HTTP handlers
type Handlers struct {
	store *Store
}

// NewHandlers creates a new Handlers instance
func NewHandlers(store *Store) *Handlers {
	return &Handlers{store: store}
}

// HomeHandler shows all posts (viewer)
func (h *Handlers) HomeHandler(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}

	posts := h.store.List()

	tmpl := `<!DOCTYPE html>
<html>
<head><title>Blog</title></head>
<body>
<h1>Blog Posts</h1>
{{if .}}
<ul>
{{range .}}
<li><a href="/post/{{.ID}}">{{.Title}}</a> - {{.CreatedAt.Format "Jan 2, 2006"}}</li>
{{end}}
</ul>
{{else}}
<p>No posts yet.</p>
{{end}}
</body>
</html>`

	t := template.Must(template.New("home").Parse(tmpl))
	t.Execute(w, posts)
}

// PostHandler shows a single post (viewer)
func (h *Handlers) PostHandler(w http.ResponseWriter, r *http.Request) {
	idStr := strings.TrimPrefix(r.URL.Path, "/post/")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.NotFound(w, r)
		return
	}

	post := h.store.Get(id)
	if post == nil {
		http.NotFound(w, r)
		return
	}

	tmpl := `<!DOCTYPE html>
<html>
<head><title>{{.Title}}</title></head>
<body>
<h1>{{.Title}}</h1>
<p><em>{{.CreatedAt.Format "Jan 2, 2006 3:04 PM"}}</em></p>
<div>{{.Content}}</div>
<p><a href="/">Back to home</a></p>
</body>
</html>`

	t := template.Must(template.New("post").Parse(tmpl))
	t.Execute(w, post)
}

// AdminHandler shows admin dashboard (writer)
func (h *Handlers) AdminHandler(w http.ResponseWriter, r *http.Request) {
	posts := h.store.List()

	tmpl := `<!DOCTYPE html>
<html>
<head><title>Admin Dashboard</title></head>
<body>
<h1>Admin Dashboard</h1>
<p><a href="/admin/create">Create New Post</a></p>
<h2>All Posts</h2>
{{if .}}
<ul>
{{range .}}
<li>
{{.Title}} - {{.CreatedAt.Format "Jan 2, 2006"}}
<form action="/admin/delete/{{.ID}}" method="POST" style="display:inline;">
<button type="submit" onclick="return confirm('Delete this post?')">Delete</button>
</form>
</li>
{{end}}
</ul>
{{else}}
<p>No posts yet.</p>
{{end}}
</body>
</html>`

	t := template.Must(template.New("admin").Parse(tmpl))
	t.Execute(w, posts)
}

// CreateHandler handles post creation (writer)
func (h *Handlers) CreateHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		tmpl := `<!DOCTYPE html>
<html>
<head><title>Create Post</title></head>
<body>
<h1>Create New Post</h1>
<form method="POST">
<p><label>Title: <input type="text" name="title" required></label></p>
<p><label>Content:<br><textarea name="content" rows="10" cols="50" required></textarea></label></p>
<p><button type="submit">Create Post</button></p>
</form>
<p><a href="/admin">Back to admin</a></p>
</body>
</html>`
		t := template.Must(template.New("create").Parse(tmpl))
		t.Execute(w, nil)
		return
	}

	if r.Method == http.MethodPost {
		title := r.FormValue("title")
		content := r.FormValue("content")

		if title == "" || content == "" {
			http.Error(w, "Title and content are required", http.StatusBadRequest)
			return
		}

		h.store.Create(title, content)
		http.Redirect(w, r, "/admin", http.StatusSeeOther)
		return
	}

	http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
}

// DeleteHandler handles post deletion (writer)
func (h *Handlers) DeleteHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	idStr := strings.TrimPrefix(r.URL.Path, "/admin/delete/")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.NotFound(w, r)
		return
	}

	if !h.store.Delete(id) {
		http.NotFound(w, r)
		return
	}

	fmt.Println("Deleted post", id)
	http.Redirect(w, r, "/admin", http.StatusSeeOther)
}
