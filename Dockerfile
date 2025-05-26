# ── build stage ─────────────────────────────────────────────
FROM docker.io/library/golang:1.23-alpine AS builder
# Install templ for generating Go code from templates
RUN go install github.com/a-h/templ/cmd/templ@latest
WORKDIR /src
# Copy go mod files first for better caching
COPY go.mod go.sum ./
RUN go mod download
# Copy the rest of the source tree
COPY . .
# Generate templ templates
RUN templ generate
# Build the Go application
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o server .

# ── runtime stage ────────────────────────────────────────────
FROM gcr.io/distroless/static:nonroot
WORKDIR /app
# Copy the built server
COPY --from=builder /src/server ./
# Copy the web directory (contains generated templates and static files)
COPY --from=builder /src/web ./web/
USER nonroot
EXPOSE 8080
CMD ["./server"]
