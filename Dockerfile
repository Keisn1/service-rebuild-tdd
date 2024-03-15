ARG GO_VERSION
FROM golang:${GO_VERSION}-alpine AS base


WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN go build -v -o ./note-taking-app ./controllers
CMD ./note-taking-app

LABEL maintainer=<kay@kayarch>
