# Simple Blog

A minimal Go blog with writer/viewer roles.

## Requirements

- Go 1.23+

## How to Run

```bash
go run .
```

Starts on port 8080.

## Routes

| Method | Path | Description | Access |
|--------|------|-------------|--------|
| GET | `/` | List all posts | Public |
| GET | `/post/{id}` | View single post | Public |
| GET | `/admin` | Admin dashboard | Protected |
| GET/POST | `/admin/create` | Create new post | Protected |
| POST | `/admin/delete/{id}` | Delete post | Protected |

## Authentication

Protected routes use HTTP Basic Auth.

- **Username:** `writer`
- **Password:** `password123`

## Note

Uses in-memory storage. All data is lost on restart.
