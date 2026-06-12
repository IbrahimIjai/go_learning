# Multi-stage build: compile in a full Go image, ship a ~15MB final image.
# This is the standard Go deployment pattern.

FROM golang:1.22-alpine AS build
WORKDIR /src
# Copy module files first: this layer only rebuilds when deps change,
# so day-to-day code edits reuse the cached download layer.
COPY go.mod go.sum ./
RUN go mod download
COPY . .
# CGO_ENABLED=0 → fully static binary, runs in a from-scratch-style image.
RUN CGO_ENABLED=0 go build -ldflags="-s -w" -o /api ./cmd/api

FROM alpine:3.19
RUN adduser -D -u 10001 app          # never run as root in containers
USER app
COPY --from=build /api /api
EXPOSE 8080
ENTRYPOINT ["/api"]
