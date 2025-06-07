# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

Waqti.me is a mobile-first web platform for creators, freelancers, and micro-studios in Kuwait and the GCC to monetize their time through workshops and courses. Built with Go, Echo framework, PostgreSQL, and templ templates with bilingual support (Arabic/English) and RTL layout handling.

## Development Commands

### Essential Commands
- `make install` - Install Go dependencies and templ CLI
- `make run` - Generate templates and run the application
- `make dev` - Run with hot reload (requires air: `go install github.com/cosmtrek/air@latest`)
- `make build` - Build production binary to bin/waqti
- `make templ` - Generate templ templates (required before build/run)
- `make test` - Run Go tests
- `make clean` - Remove build artifacts and generated template files

### Template Development
Always run `make templ` after modifying .templ files to regenerate the corresponding _templ.go files. The application won't compile without this step.

## Architecture Overview

### Layer Structure
```
main.go → handlers → services → database
```

**Handlers**: HTTP request controllers (Echo framework)
**Services**: Business logic layer with fallback data for demos
**Models**: Domain objects with UUID IDs and bilingual field support
**Database**: PostgreSQL with connection pooling, triggers, and materialized views

### Key Patterns
- **Authentication**: Session-based with 30-day expiration, bcrypt password hashing
- **Middleware**: Language detection, conditional auth, CORS support
- **Templates**: Templ (type-safe Go templates) with component-based architecture
- **Internationalization**: Context-based Arabic/English switching with RTL support
- **Database**: UUID primary keys, automatic timestamps, row-level security enabled

### Database Connection
Uses singleton pattern with `database.Instance`. Connection configured via environment variables with pooling (25 max connections, 5-minute lifetime).

### Route Organization
```go
// Public routes
e.GET("/", authHandler.ShowLandingPage)
e.GET("/:username", authHandler.ShowStorePage) 

// Protected routes (auth middleware applied)
protected := e.Group("")
protected.GET("/dashboard", dashboardHandler.ShowDashboard)
```

## Key Business Logic

### Core Models
- **Creator**: User accounts with shop settings and branding
- **Workshop**: Course definitions with bilingual content
- **WorkshopSession**: Scheduled sessions with capacity limits
- **Order/Enrollment**: Purchase and registration flow
- **Analytics**: Click tracking and creator insights

### Service Layer Features
- Fallback demo data when database operations fail
- Bilingual content handling with language context
- File upload management for workshop images
- QR code generation for workshop sharing

## Template System

Uses templ for type-safe templates with automatic Go code generation. Templates support:
- Bilingual rendering with language context
- RTL layout switching
- Component reusability
- HTMX integration for dynamic content
- Alpine.js for client-side interactivity

## Environment Configuration

Required environment variables for database connection:
- `DB_HOST`, `DB_PORT`, `DB_USER`, `DB_PASSWORD`, `DB_NAME`
- Defaults to local PostgreSQL if not set

## Development Notes

- Templates must be regenerated after changes (`make templ`)
- Database schema in schema.sql includes triggers and indexes
- Authentication sessions stored in database with expiration
- Static files served from web/static/ with upload/ subdirectory for user content
- Responsive design with mobile-first approach and Gulf-inspired color scheme