## Feature Flags System

### Overview

This is a feature flags system that allows users to create, manage, and assign feature flags to people.

### Project Structure

```
/feature_flags
│
├── /cmd                      # Application entry points (for multiple binaries, if any)
│   └── /app                  # Main application folder (main.go for your application)
│
├── /internal                  # Private application and library code
│   ├── /feature_flags         # Business logic for handling feature flags
│   ├── /person                # Business logic for handling person
│   ├── /auth                  # Authentication logic (if needed)
│   └── /db                    # Database handling (models, repositories, queries, etc.)
│
├── /pkg                       # Shared library code (can be imported by other projects)
│   └── /utils                 # Utility packages (helpers, shared functionality)
│
├── /api                       # API handlers and routes
│   ├── /handlers              # Handlers for specific API endpoints
│   ├── /middlewares           # Middleware functions (e.g., for logging, auth, etc.)
│   └── /routes                # API route setup
│
├── /web                       # The web application using templ
│   └── /app                   # App initialization and routes definition
│   └── /assets                # Web assets (images, icons, css)
│   └── /components            # Go templ components
│   └── /handlers              # Handlers to receive data from the templ and process it
│   └── /services (deprecated) # Services used with go templates (`poc/htmx` branch)
│   └── /types*                # Common type struct (dunno where to place it)
│   └── /utils                 # Utilitarian function/methods
│   └── /views                 # Go templ views or pages
│
├── /migrations                # Database migration files (SQL files for schema changes)
│
├── /configs                   # Configuration files (YAML, JSON, etc.)
│
├── /scripts                   # Helper scripts (build, start, deploy, etc.)
│
├── go.mod                     # Go module file
├── go.sum                     # Dependency checksum file
└── README.md                  # Project documentation
```

### Running the project

- API `make run-api`
- Templ with HTMX `make run-templ`
- Go template `make run-web` on branch `poc/htmx`
