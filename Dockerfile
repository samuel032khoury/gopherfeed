# Stage 1: Build the frontend
FROM --platform=linux/amd64 oven/bun:latest as frontend
WORKDIR /app/web
COPY web/package.json web/bun.lock ./
RUN bun install --frozen-lockfile
COPY web/ ./
RUN bun run build

# Stage 2: Build the Go binary
FROM --platform=linux/amd64 golang:1.25.5 as builder
WORKDIR /app
COPY . .
# Copy the frontend build output into the location expected by go:embed
COPY --from=frontend /app/web/dist ./cmd/api/web/dist
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o api cmd/api/*.go

# Stage 3: Minimal runtime image
FROM --platform=linux/amd64 scratch
WORKDIR /app
# Copy CA certificates
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=builder /app/api .
EXPOSE 8080
CMD ["./api"]