ARG GO_VERSION=1
FROM golang:${GO_VERSION}-bookworm as builder

WORKDIR /usr/src/app
COPY go.mod go.sum ./
RUN go mod download && go mod verify
COPY fly.toml go.mod go.sum index.css tailwind.config.js .
# COPY static/ ./static/
COPY cmd/ ./cmd/
COPY internal/ ./internal/
RUN go build -v -x -o /run-app ./cmd/web
RUN go build -v -x -o /custom-tools ./cmd/cli

FROM debian:bookworm

ENV PORT=8080

RUN apt update && apt install -y ca-certificates

COPY --from=builder /run-app /usr/local/bin/
COPY --from=builder /custom-tools /usr/local/bin/

CMD ["run-app", "--db", "data/recipes.db", "--static-path", "/data/static/"] 
