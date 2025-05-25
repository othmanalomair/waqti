# Waqti.me - Creator Dashboard

A mobile-first web platform for creators, freelancers, and micro-studios in Kuwait and the GCC to monetize their time.

## Tech Stack

- **Backend**: Go with Echo framework
- **Frontend**: Templ templates + Alpine.js + HTMX
- **Styling**: Tailwind CSS with RTL support
- **Database**: PostgreSQL (to be integrated)

## Features

- 🌍 Bilingual support (Arabic/English)
- 📱 Mobile-first responsive design
- 🎨 Gulf-inspired design with brand colors
- ⚡ Fast server-side rendering with Templ
- 🔄 Interactive components with Alpine.js + HTMX

## Project Structure

```
waqti/
├── main.go                 # Application entry point
├── go.mod                  # Go module dependencies
├── Makefile               # Build and development commands
├── .air.toml              # Hot reload configuration
├── internal/
│   ├── handlers/          # HTTP handlers
│   │   └── dashboard.go
│   ├── middleware/        # Custom middleware
│   │   └── language.go
│   ├── models/           # Data models
│   │   └── creator.go
│   └── services/         # Business logic
│       ├── creator.go
│       └── workshop.go
└── web/
    ├── templates/        # Templ templates
    │   └── dashboard.templ
    └── static/          # Static assets
        └── css/
            └── styles.css
```

## Getting Started

### Prerequisites

- Go 1.21+
- Templ CLI tool

### Installation

1. Clone the repository
2. Install dependencies:
   ```bash
   make install
   ```

3. Run the application:
   ```bash
   make run
   ```

4. Visit `http://localhost:8080/dashboard`

### Development

For hot reload during development:

```bash
# Install air for hot reload
go install github.com/cosmtrek/air@latest

# Run in development mode
make dev
```

## API Endpoints

- `GET /dashboard` - Creator dashboard page
- `POST /dashboard/toggle-language` - Toggle between Arabic/English

## Features Implemented

- ✅ Bilingual UI (Arabic/English)
- ✅ Creator profile display
- ✅ Workshop management menu
- ✅ Mobile-first responsive design
- ✅ Language toggle functionality
- ✅ Gulf-inspired color scheme
- ✅ RTL support for Arabic

## Next Steps

- [ ] Integrate PostgreSQL database
- [ ] Add user authentication
- [ ] Implement workshop CRUD operations
- [ ] Add payment integration (Stripe, Knet, MyFatoorah)
- [ ] Build creator public profile pages
- [ ] Add analytics dashboard
- [ ] Implement QR code generation


## License

Private project - All rights reserved.
