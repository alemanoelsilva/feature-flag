## Feature Flags System

### Overview

This is a feature flags system that allows users to create, manage, and assign feature flags to people.

### Project Structure

```
/feature_flags
│
├── /cmd                     # Application entry points (for multiple binaries, if any)
│   └── /app                 # Main application folder (main.go for your application)
│
├── /internal                # Private application and library code
│   ├── /feature_flags       # Business logic for handling feature flags
│   ├── /person              # Business logic for handling person
│   ├── /auth                # Authentication logic (if needed)
│   └── /db                  # Database handling (models, repositories, queries, etc.)
│
├── /pkg                     # Shared library code (can be imported by other projects)
│   └── /utils               # Utility packages (helpers, shared functionality)
│
├── /api                     # API handlers and routes
│   ├── /handlers            # Handlers for specific API endpoints
│   ├── /middlewares         # Middleware functions (e.g., for logging, auth, etc.)
│   └── /routes              # API route setup
│
├── /web                     # Frontend static files (HTML, CSS, JS)
│   └── /templates           # HTML templates for your HTMX frontend
│
├── /migrations               # Database migration files (SQL files for schema changes)
│
├── /configs                  # Configuration files (YAML, JSON, etc.)
│
├── /scripts                  # Helper scripts (build, start, deploy, etc.)
│
├── go.mod                    # Go module file
├── go.sum                    # Dependency checksum file
└── README.md                 # Project documentation
```

### Running the project

```
go run cmd/app/main.go
```
