# --- build stage ---
FROM golang:1.22-alpine AS build
WORKDIR /src

# copy the module files first so dependencies are cached between builds
COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o /out/alumni ./cmd/api  && CGO_ENABLED=0 GOOS=linux go build -o /out/seed ./cmd/seed

# --- runtime stage ---
FROM alpine:3.20
RUN adduser -D -H appuser
WORKDIR /app
COPY --from=build /out/alumni /out/seed ./

ENV APP_PORT=8080
EXPOSE 8080

USER appuser
ENTRYPOINT ["./alumni"]
