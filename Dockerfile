# --- build stage ---
FROM golang:1.22-alpine AS build
WORKDIR /src

# go.mod has no external deps, but copy it first for cache efficiency
COPY go.mod ./
RUN go mod download 2>/dev/null || true

COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o /out/alumni ./cmd/api

# --- runtime stage ---
FROM alpine:3.20
RUN adduser -D -H appuser
WORKDIR /app
COPY --from=build /out/alumni ./alumni

ENV APP_PORT=8080
EXPOSE 8080

USER appuser
ENTRYPOINT ["./alumni"]
